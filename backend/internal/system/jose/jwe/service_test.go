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
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"testing"

	"github.com/asgardeo/thunder/internal/system/config"
	syscrypto "github.com/asgardeo/thunder/internal/system/crypto"
	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
	"github.com/asgardeo/thunder/internal/system/log"
	"github.com/asgardeo/thunder/tests/mocks/crypto/cryptomock"
	"github.com/asgardeo/thunder/tests/mocks/crypto/pki/pkimock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type JWEServiceTestSuite struct {
	suite.Suite
	jweService *jweService
	pkiMock    *pkimock.PKIServiceInterfaceMock
	cryptoMock *cryptomock.RuntimeCryptoProviderMock
}

func TestJWEServiceSuite(t *testing.T) {
	suite.Run(t, new(JWEServiceTestSuite))
}

func (suite *JWEServiceTestSuite) SetupTest() {
	config.ResetThunderRuntime()
	suite.pkiMock = pkimock.NewPKIServiceInterfaceMock(suite.T())
	suite.cryptoMock = cryptomock.NewRuntimeCryptoProviderMock(suite.T())

	testConfig := &config.Config{
		JWT: config.JWTConfig{
			PreferredKeyID: "test-kid",
		},
	}
	_ = config.InitializeThunderRuntime("", testConfig)
}

func (suite *JWEServiceTestSuite) TestEncryptDecrypt_RSA() {
	suite.jweService = &jweService{
		kid:            "test-kid",
		cryptoProvider: suite.cryptoMock,
		logger:         log.GetLogger(),
	}

	payload := []byte("Hello, RSA JWE!")
	testCases := []struct {
		enc    ContentEncAlgorithm
		cekLen int
	}{
		{A128GCM, 16},
		{A192GCM, 24},
		{A256GCM, 32},
	}

	xorBytes := func(in []byte) []byte {
		out := make([]byte, len(in))
		for i, b := range in {
			out[i] = b ^ 0xFF
		}
		return out
	}

	for _, tc := range testCases {
		fixedCEK := make([]byte, tc.cekLen)
		_, _ = rand.Read(fixedCEK)

		suite.cryptoMock.EXPECT().Encrypt(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(xorBytes(fixedCEK), &syscrypto.CryptoDetails{CEK: fixedCEK}, nil).Once()

		suite.cryptoMock.EXPECT().Decrypt(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, _ syscrypto.KeyRef,
				_ syscrypto.AlgorithmParams, wrapped []byte,
			) ([]byte, error) {
				return xorBytes(wrapped), nil
			}).Once()

		token, sErr := suite.jweService.Encrypt(context.Background(), payload, "test-kid", RSAOAEP256, tc.enc)
		assert.Nil(suite.T(), sErr)
		decrypted, sErr := suite.jweService.Decrypt(context.Background(), token)
		assert.Nil(suite.T(), sErr)
		assert.Equal(suite.T(), payload, decrypted)
	}
}

func (suite *JWEServiceTestSuite) TestEncryptDecrypt_ECDHES_Direct() {
	suite.jweService = &jweService{
		kid:            "test-kid",
		cryptoProvider: suite.cryptoMock,
		logger:         log.GetLogger(),
	}

	payload := []byte("Hello, ECDH-ES JWE!")

	// A real ephemeral public key is required so that epkToMap / JWKToECPublicKey round-trips.
	ephKey, err := ecdh.P256().GenerateKey(rand.Reader)
	assert.NoError(suite.T(), err)
	fakeEPK := ephKey.PublicKey()

	testCases := []struct {
		enc    ContentEncAlgorithm
		cekLen int
	}{
		{A128GCM, 16},
		{A192GCM, 24},
		{A256GCM, 32},
	}

	for _, tc := range testCases {
		fixedCEK := make([]byte, tc.cekLen)
		_, _ = rand.Read(fixedCEK)

		suite.cryptoMock.EXPECT().Encrypt(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, &syscrypto.CryptoDetails{EPK: fakeEPK, CEK: fixedCEK}, nil).Once()
		suite.cryptoMock.EXPECT().Decrypt(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(fixedCEK, nil).Once()

		token, sErr := suite.jweService.Encrypt(context.Background(), payload, "test-kid", ECDHES, tc.enc)
		assert.Nil(suite.T(), sErr)
		decrypted, sErr := suite.jweService.Decrypt(context.Background(), token)
		assert.Nil(suite.T(), sErr)
		assert.Equal(suite.T(), payload, decrypted)
	}
}

func (suite *JWEServiceTestSuite) TestEncryptDecrypt_ECDHESKW() {
	suite.jweService = &jweService{
		kid:            "test-kid",
		cryptoProvider: suite.cryptoMock,
		logger:         log.GetLogger(),
	}

	payload := []byte("Hello, ECDH-ES+KW JWE!")

	ephKey, err := ecdh.P256().GenerateKey(rand.Reader)
	assert.NoError(suite.T(), err)
	fakeEPK := ephKey.PublicKey()

	testCases := []struct {
		alg    KeyEncAlgorithm
		enc    ContentEncAlgorithm
		cekLen int
	}{
		{ECDHESA128KW, A128GCM, 16},
		{ECDHESA256KW, A256GCM, 32},
	}

	for _, tc := range testCases {
		xorAABytes := func(in []byte) []byte {
			out := make([]byte, len(in))
			for i, b := range in {
				out[i] = b ^ 0xAA
			}
			return out
		}

		fixedCEK := make([]byte, tc.cekLen)
		_, _ = rand.Read(fixedCEK)

		suite.cryptoMock.EXPECT().Encrypt(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(xorAABytes(fixedCEK), &syscrypto.CryptoDetails{EPK: fakeEPK, CEK: fixedCEK}, nil).Once()

		suite.cryptoMock.EXPECT().Decrypt(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, _ syscrypto.KeyRef,
				_ syscrypto.AlgorithmParams, wrapped []byte,
			) ([]byte, error) {
				return xorAABytes(wrapped), nil
			}).Once()

		token, sErr := suite.jweService.Encrypt(context.Background(), payload, "test-kid", tc.alg, tc.enc)
		assert.Nil(suite.T(), sErr)
		decrypted, sErr := suite.jweService.Decrypt(context.Background(), token)
		assert.Nil(suite.T(), sErr)
		assert.Equal(suite.T(), payload, decrypted)
	}
}

func (suite *JWEServiceTestSuite) TestEncrypt_Errors() {
	suite.jweService = &jweService{
		kid:            "test-kid",
		cryptoProvider: suite.cryptoMock,
		logger:         log.GetLogger(),
	}

	// Unsupported content encryption algorithm.
	_, sErr := suite.jweService.Encrypt(context.Background(), []byte("p"), "key-id", RSAOAEP256, "INVALID")
	assert.NotNil(suite.T(), sErr)
	assert.Equal(suite.T(), ErrorUnsupportedEncryptionAlgorithm, *sErr)

	// Provider returns any error (including unsupported algorithm) → server error.
	suite.cryptoMock.EXPECT().Encrypt(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil, errors.New("unsupported algorithm: INVALID_ALG")).Once()
	_, sErr = suite.jweService.Encrypt(context.Background(), []byte("p"), "key-id", "INVALID_ALG", A128GCM)
	assert.NotNil(suite.T(), sErr)
	assert.Equal(suite.T(), serviceerror.InternalServerError, *sErr)

	// Provider returns nil error but nil details (missing CEK) → unsupported algorithm error.
	suite.cryptoMock.EXPECT().Encrypt(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil, nil).Once()
	_, sErr = suite.jweService.Encrypt(context.Background(), []byte("p"), "key-id", RSAOAEP256, A128GCM)
	assert.NotNil(suite.T(), sErr)
	assert.Equal(suite.T(), ErrorUnsupportedJWEAlgorithm, *sErr)

	// RuntimeCryptoProvider returns error for ECDH-ES.
	suite.cryptoMock.EXPECT().Encrypt(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil, errors.New("crypto error")).Once()
	_, sErr = suite.jweService.Encrypt(context.Background(), []byte("p"), "key-id", ECDHES, A128GCM)
	assert.NotNil(suite.T(), sErr)
	assert.Equal(suite.T(), serviceerror.InternalServerError, *sErr)

	// RuntimeCryptoProvider returns error for ECDH-ES+KW.
	suite.cryptoMock.EXPECT().Encrypt(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil, errors.New("crypto error")).Once()
	_, sErr = suite.jweService.Encrypt(context.Background(), []byte("p"), "key-id", ECDHESA128KW, A128GCM)
	assert.NotNil(suite.T(), sErr)
	assert.Equal(suite.T(), serviceerror.InternalServerError, *sErr)

	// RuntimeCryptoProvider returns error for RSA-OAEP-256.
	suite.cryptoMock.EXPECT().Encrypt(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil, errors.New("crypto error")).Once()
	_, sErr = suite.jweService.Encrypt(context.Background(), []byte("p"), "key-id", RSAOAEP256, A128GCM)
	assert.NotNil(suite.T(), sErr)
	assert.Equal(suite.T(), serviceerror.InternalServerError, *sErr)
}

func (suite *JWEServiceTestSuite) TestDecrypt_Errors() {
	suite.jweService = &jweService{
		kid:            "test-kid",
		cryptoProvider: suite.cryptoMock,
		logger:         log.GetLogger(),
	}

	// Invalid JWE format.
	_, sErr := suite.jweService.Decrypt(context.Background(), "invalid.jwe")
	assert.NotNil(suite.T(), sErr)
	assert.Equal(suite.T(), ErrorDecodingJWE, *sErr)

	// Build a valid RSA JWE token, then have the crypto provider fail on decrypt.
	fixedCEK1 := make([]byte, 16)
	suite.cryptoMock.EXPECT().Encrypt(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]byte("wrapped-cek"), &syscrypto.CryptoDetails{CEK: fixedCEK1}, nil).Once()
	token, _ := suite.jweService.Encrypt(context.Background(), []byte("data"), "test-kid", RSAOAEP256, A128GCM)

	suite.cryptoMock.EXPECT().Decrypt(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("key not found")).Once()
	_, sErr = suite.jweService.Decrypt(context.Background(), token)
	assert.NotNil(suite.T(), sErr)
	assert.Equal(suite.T(), ErrorJWEDecryptionFailed, *sErr)

	// Crypto provider returns wrong CEK → decryptContent authentication tag mismatch.
	fixedCEK2 := make([]byte, 16)
	_, _ = rand.Read(fixedCEK2)
	suite.cryptoMock.EXPECT().Encrypt(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]byte("wrapped-cek"), &syscrypto.CryptoDetails{CEK: fixedCEK2}, nil).Once()
	validToken, _ := suite.jweService.Encrypt(context.Background(), []byte("data"), "test-kid", RSAOAEP256, A128GCM)

	wrongCEK := make([]byte, 16)
	suite.cryptoMock.EXPECT().Decrypt(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(wrongCEK, nil).Once()
	_, sErr = suite.jweService.Decrypt(context.Background(), validToken)
	assert.NotNil(suite.T(), sErr)
	assert.Equal(suite.T(), ErrorJWEDecryptionFailed, *sErr)
}

func (suite *JWEServiceTestSuite) TestDecrypt_EdgeCases() {
	suite.jweService = &jweService{
		kid:            "test-kid",
		cryptoProvider: suite.cryptoMock,
		logger:         log.GetLogger(),
	}

	// Malformed JWE (wrong number of parts).
	_, sErr := suite.jweService.Decrypt(context.Background(), "malformed.jwe")
	assert.NotNil(suite.T(), sErr)
	assert.Equal(suite.T(), ErrorDecodingJWE, *sErr)

	// Invalid base64 data in header.
	_, sErr = suite.jweService.Decrypt(context.Background(), "invalid-base64.key.iv.ciphertext.tag")
	assert.NotNil(suite.T(), sErr)
	assert.Equal(suite.T(), ErrorDecodingJWE, *sErr)

	// Invalid JSON in header.
	invalidHeader := base64.RawURLEncoding.EncodeToString([]byte("{invalid json"))
	_, sErr = suite.jweService.Decrypt(context.Background(), invalidHeader+".key.iv.ciphertext.tag")
	assert.NotNil(suite.T(), sErr)
	assert.Equal(suite.T(), ErrorDecodingJWE, *sErr)

	// Missing alg field in header.
	headerMissingAlg := base64.RawURLEncoding.EncodeToString([]byte(`{"enc":"A128GCM"}`))
	_, sErr = suite.jweService.Decrypt(context.Background(), headerMissingAlg+".key.iv.ciphertext.tag")
	assert.NotNil(suite.T(), sErr)
	assert.Equal(suite.T(), ErrorUnsupportedJWEAlgorithm, *sErr)

	// ECDH-ES header missing epk field → provider call with nil EPK → returns error.
	suite.cryptoMock.EXPECT().Decrypt(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("EPK required")).Once()
	headerMissingEPK := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"ECDH-ES","enc":"A128GCM"}`))
	_, sErr = suite.jweService.Decrypt(context.Background(), headerMissingEPK+".key.iv.ciphertext.tag")
	assert.NotNil(suite.T(), sErr)
	assert.Equal(suite.T(), ErrorJWEDecryptionFailed, *sErr)
}

func (suite *JWEServiceTestSuite) TestInitialize() {
	suite.pkiMock.EXPECT().GetCertThumbprint(mock.Anything).Return("test-kid")

	service, err := Initialize(suite.pkiMock)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), service)
}
