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
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/asgardeo/thunder/internal/system/jose/jws"
)

// epkToMap converts an ephemeral public key to a JWK-like map representation.
func epkToMap(pub crypto.PublicKey) map[string]interface{} {
	ecdhPub, ok := pub.(*ecdh.PublicKey)
	if !ok {
		return nil
	}

	raw := ecdhPub.Bytes()
	var crv string
	var x, y []byte

	switch len(raw) {
	case 65:
		crv = jws.P256
		x = raw[1:33]
		y = raw[33:]
	case 97:
		crv = jws.P384
		x = raw[1:49]
		y = raw[49:]
	case 133:
		crv = jws.P521
		x = raw[1:67]
		y = raw[67:]
	default:
		return nil
	}

	return map[string]interface{}{
		"kty": "EC",
		"crv": crv,
		"x":   base64.RawURLEncoding.EncodeToString(x),
		"y":   base64.RawURLEncoding.EncodeToString(y),
	}
}

// encryptContent encrypts the payload using the content encryption key (CEK).
func encryptContent(payload []byte, cek []byte, enc ContentEncAlgorithm, aad []byte) ([]byte, []byte, []byte, error) {
	block, err := aes.NewCipher(cek)
	if err != nil {
		return nil, nil, nil, err
	}

	var gcm cipher.AEAD
	switch enc {
	case A128GCM, A192GCM, A256GCM:
		gcm, err = cipher.NewGCM(block)
	default:
		return nil, nil, nil, fmt.Errorf("unsupported encryption algorithm: %s", enc)
	}

	if err != nil {
		return nil, nil, nil, err
	}

	iv := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(iv); err != nil {
		return nil, nil, nil, err
	}

	ciphertextWithTag := gcm.Seal(nil, iv, payload, aad)
	tagSize := gcm.Overhead()
	ciphertext := ciphertextWithTag[:len(ciphertextWithTag)-tagSize]
	tag := ciphertextWithTag[len(ciphertextWithTag)-tagSize:]

	return iv, ciphertext, tag, nil
}

// decryptContent decrypts the ciphertext using the content encryption key (CEK).
func decryptContent(ciphertext, iv, tag []byte, cek []byte, enc ContentEncAlgorithm, aad []byte) ([]byte, error) {
	block, err := aes.NewCipher(cek)
	if err != nil {
		return nil, err
	}

	var gcm cipher.AEAD
	switch enc {
	case A128GCM, A192GCM, A256GCM:
		gcm, err = cipher.NewGCM(block)
	default:
		return nil, fmt.Errorf("unsupported encryption algorithm: %s", enc)
	}

	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, iv, append(ciphertext, tag...), aad)
}

// DecodeJWE decodes a JWE compact serialization into its five parts.
func DecodeJWE(jweToken string) (header map[string]interface{}, headerBase64 string,
	encryptedKey, iv, ciphertext, tag []byte, err error) {
	parts := strings.Split(jweToken, ".")
	if len(parts) != 5 {
		return nil, "", nil, nil, nil, nil, errors.New("invalid JWE format")
	}

	headerBase64 = parts[0]
	headerBytes, err := base64.RawURLEncoding.DecodeString(headerBase64)
	if err != nil {
		return nil, "", nil, nil, nil, nil, fmt.Errorf("failed to decode header: %w", err)
	}

	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, "", nil, nil, nil, nil, fmt.Errorf("failed to unmarshal header: %w", err)
	}

	encryptedKey, err = base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, "", nil, nil, nil, nil, fmt.Errorf("failed to decode encrypted key: %w", err)
	}

	iv, err = base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, "", nil, nil, nil, nil, fmt.Errorf("failed to decode IV: %w", err)
	}

	ciphertext, err = base64.RawURLEncoding.DecodeString(parts[3])
	if err != nil {
		return nil, "", nil, nil, nil, nil, fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	tag, err = base64.RawURLEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, "", nil, nil, nil, nil, fmt.Errorf("failed to decode tag: %w", err)
	}

	return header, headerBase64, encryptedKey, iv, ciphertext, tag, nil
}
