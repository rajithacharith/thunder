/*
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package jwe

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/asgardeo/thunder/internal/system/config"
	syscrypto "github.com/asgardeo/thunder/internal/system/crypto"
	"github.com/asgardeo/thunder/internal/system/crypto/pki"
	cryptoruntime "github.com/asgardeo/thunder/internal/system/crypto/runtime"
	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
	"github.com/asgardeo/thunder/internal/system/jose/jws"
	"github.com/asgardeo/thunder/internal/system/log"
)

// JWEServiceInterface defines the interface for JWE operations.
type JWEServiceInterface interface {
	Encrypt(ctx context.Context, payload []byte, recipientKeyID string,
		alg KeyEncAlgorithm, enc ContentEncAlgorithm) (string, *serviceerror.ServiceError)
	Decrypt(ctx context.Context, jweToken string) ([]byte, *serviceerror.ServiceError)
}

// jweService implements the JWEServiceInterface.
type jweService struct {
	kid            string
	cryptoProvider syscrypto.RuntimeCryptoProvider
	logger         *log.Logger
}

// newJWEService creates a new JWE service instance.
func newJWEService(pkiService pki.PKIServiceInterface) (JWEServiceInterface, error) {
	preferredKid := config.GetThunderRuntime().Config.JWT.PreferredKeyID
	kid := pkiService.GetCertThumbprint(preferredKid)
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, "JWEService"))

	return &jweService{
		kid:            kid,
		cryptoProvider: cryptoruntime.GetRuntimeCryptoService(),
		logger:         logger,
	}, nil
}

// Encrypt encrypts the payload using the recipient key identified by recipientKeyID.
func (js *jweService) Encrypt(ctx context.Context, payload []byte, recipientKeyID string,
	alg KeyEncAlgorithm, enc ContentEncAlgorithm) (string, *serviceerror.ServiceError) {
	switch enc {
	case A128GCM, A192GCM, A256GCM:
	default:
		return "", &ErrorUnsupportedEncryptionAlgorithm
	}

	keyRef := syscrypto.KeyRef{KeyID: recipientKeyID}
	params := syscrypto.AlgorithmParams{Algorithm: syscrypto.Algorithm(alg)}
	switch alg {
	case RSAOAEP256:
		params.RSAOAEP256 = syscrypto.RSAOAEP256Params{
			ContentEncryptionAlgorithm: syscrypto.Algorithm(enc),
		}
	case ECDHES, ECDHESA128KW, ECDHESA256KW:
		params.ECDHES = syscrypto.ECDHESParams{
			ContentEncryptionAlgorithm: syscrypto.Algorithm(enc),
		}
	}

	encryptedKey, details, err := js.cryptoProvider.Encrypt(ctx, keyRef, params, nil)
	if err != nil {
		js.logger.Error("Failed to establish key", log.Error(err))
		return "", &serviceerror.InternalServerError
	}
	if details == nil || details.CEK == nil {
		return "", &ErrorUnsupportedJWEAlgorithm
	}

	// Build JWE header.
	header := map[string]interface{}{
		"alg": string(alg),
		"enc": string(enc),
		"typ": "JWE",
		"kid": recipientKeyID,
	}
	if details.EPK != nil {
		header["epk"] = epkToMap(details.EPK)
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		js.logger.Error("Failed to marshal JWE header: " + err.Error())
		return "", &serviceerror.InternalServerError
	}
	headerBase64 := base64.RawURLEncoding.EncodeToString(headerJSON)

	// Encrypt content.
	iv, ciphertext, tag, err := encryptContent(payload, details.CEK, enc, []byte(headerBase64))
	if err != nil {
		js.logger.Error("Failed to encrypt content: " + err.Error())
		return "", &serviceerror.InternalServerError
	}

	// Build compact serialization.
	return fmt.Sprintf("%s.%s.%s.%s.%s",
		headerBase64,
		base64.RawURLEncoding.EncodeToString(encryptedKey),
		base64.RawURLEncoding.EncodeToString(iv),
		base64.RawURLEncoding.EncodeToString(ciphertext),
		base64.RawURLEncoding.EncodeToString(tag),
	), nil
}

// Decrypt decrypts the JWE compact serialization using the server's private key via RuntimeCryptoProvider.
func (js *jweService) Decrypt(ctx context.Context, jweToken string) ([]byte, *serviceerror.ServiceError) {
	header, headerBase64, encryptedKey, iv, ciphertext, tag, err := DecodeJWE(jweToken)
	if err != nil {
		js.logger.Debug("Failed to decode JWE: " + err.Error())
		return nil, &ErrorDecodingJWE
	}

	algStr, ok := header["alg"].(string)
	if !ok {
		return nil, &ErrorUnsupportedJWEAlgorithm
	}
	encStr, ok := header["enc"].(string)
	if !ok {
		return nil, &ErrorUnsupportedEncryptionAlgorithm
	}

	kidStr, _ := header["kid"].(string)
	if kidStr == "" {
		kidStr = js.kid
	}
	keyRef := syscrypto.KeyRef{KeyID: kidStr}
	params := syscrypto.AlgorithmParams{Algorithm: syscrypto.Algorithm(algStr)}
	alg := KeyEncAlgorithm(algStr)
	if alg == ECDHES {
		params.ECDHES.ContentEncryptionAlgorithm = syscrypto.Algorithm(encStr)
	}
	switch alg {
	case ECDHES, ECDHESA128KW, ECDHESA256KW:
		if epkMap, ok := header["epk"].(map[string]interface{}); ok {
			ephemeralPub, epkErr := jws.JWKToECPublicKey(epkMap)
			if epkErr != nil {
				js.logger.Error("Failed to extract EPK from JWE header: " + epkErr.Error())
				return nil, &ErrorJWEDecryptionFailed
			}
			params.ECDHES.EPK = ephemeralPub
		}
	}

	cek, err := js.cryptoProvider.Decrypt(ctx, keyRef, params, encryptedKey)
	if err != nil {
		js.logger.Error("Failed to decrypt CEK: " + err.Error())
		return nil, &ErrorJWEDecryptionFailed
	}

	payload, err := decryptContent(ciphertext, iv, tag, cek, ContentEncAlgorithm(encStr), []byte(headerBase64))
	if err != nil {
		js.logger.Error("Failed to decrypt content: " + err.Error())
		return nil, &ErrorJWEDecryptionFailed
	}

	return payload, nil
}
