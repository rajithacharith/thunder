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

package crypto

import (
	gocrypto "crypto"
	"crypto/tls"
)

// KeyRef identifies a cryptographic key by its ID.
type KeyRef struct {
	KeyID string
}

// Algorithm represents a cryptographic algorithm identifier.
type Algorithm string

const (
	// AlgorithmRS256 represents RSA PKCS1v15 signature with SHA-256 (JWA, RFC 7518).
	AlgorithmRS256 Algorithm = "RS256"
	// AlgorithmRS512 represents RSA PKCS1v15 signature with SHA-512 (JWA, RFC 7518).
	AlgorithmRS512 Algorithm = "RS512"
	// AlgorithmPS256 represents RSA-PSS signature with SHA-256 (JWA, RFC 7518).
	AlgorithmPS256 Algorithm = "PS256"
	// AlgorithmES256 represents ECDSA signature with SHA-256 (JWA, RFC 7518).
	AlgorithmES256 Algorithm = "ES256"
	// AlgorithmES384 represents ECDSA signature with SHA-384 (JWA, RFC 7518).
	AlgorithmES384 Algorithm = "ES384"
	// AlgorithmES512 represents ECDSA signature with SHA-512 (JWA, RFC 7518).
	AlgorithmES512 Algorithm = "ES512"
	// AlgorithmEdDSA represents EdDSA signature algorithm (JWA, RFC 7518).
	AlgorithmEdDSA Algorithm = "EdDSA"
	// AlgorithmRSAOAEP256 represents RSA-OAEP key encryption with SHA-256 (JWA, RFC 7518).
	AlgorithmRSAOAEP256 Algorithm = "RSA-OAEP-256"
	// AlgorithmECDHES represents ECDH-ES direct key agreement (JWA, RFC 7518).
	AlgorithmECDHES Algorithm = "ECDH-ES"
	// AlgorithmECDHESA128KW represents ECDH-ES with AES-128 key wrap (JWA, RFC 7518).
	AlgorithmECDHESA128KW Algorithm = "ECDH-ES+A128KW"
	// AlgorithmECDHESA256KW represents ECDH-ES with AES-256 key wrap (JWA, RFC 7518).
	AlgorithmECDHESA256KW Algorithm = "ECDH-ES+A256KW"
	// AlgorithmAESGCM represents AES-GCM symmetric authenticated encryption.
	AlgorithmAESGCM Algorithm = "AES-GCM"
)

// PublicKeyFilter specifies criteria for filtering public keys in GetPublicKeys.
type PublicKeyFilter struct {
	KeyID     string
	Algorithm Algorithm
}

// PublicKeyInfo describes a public key returned by GetPublicKeys.
type PublicKeyInfo struct {
	KeyID      string
	Algorithm  Algorithm
	PublicKey  gocrypto.PublicKey
	Thumbprint string
}

// TLSMaterial holds the TLS certificate material for a key reference.
type TLSMaterial struct {
	Certificate tls.Certificate
}

// AlgorithmParams carries the algorithm and any algorithm-specific inputs for a crypto operation.
// The relevant algorithm-specific ContentEncryptionAlgorithm must be set to the content encryption
// algorithm (e.g. "A128GCM") for the following operations:
//   - RSA-OAEP-256 encrypt (determines CEK size)
//   - ECDH-ES encrypt/decrypt (determines CEK size and is used as the KDF algorithm identifier)
//   - ECDH-ES+A128KW / ECDH-ES+A256KW encrypt (determines CEK size)
//
// ECDH-ES+A128KW and ECDH-ES+A256KW decrypt do not require ContentEncryptionAlgorithm because
// the KDF uses the alg value (e.g. "ECDH-ES+A128KW") directly.
type AlgorithmParams struct {
	Algorithm  Algorithm
	RSAOAEP256 RSAOAEP256Params
	ECDHES     ECDHESParams
}

// RSAOAEP256Params carries RSA-OAEP-256-specific inputs.
type RSAOAEP256Params struct {
	ContentEncryptionAlgorithm Algorithm
}

// ECDHESParams carries ECDH-ES-specific inputs.
// For ECDH-ES decrypt, EPK must be populated with the ephemeral public key from the JWE header.
type ECDHESParams struct {
	EPK                        gocrypto.PublicKey
	ContentEncryptionAlgorithm Algorithm
}

// CryptoDetails carries algorithm-specific outputs from an Encrypt operation.
// EPK is the generated ephemeral public key for ECDH-ES variants, to be embedded in the JWE header.
// CEK is the content encryption key generated or derived during key establishment.
// Nil CryptoDetails is returned for algorithms that produce no extra output (e.g. AES-GCM).
// For RSA-OAEP-256 and ECDH-ES variants, both EPK (where applicable) and CEK are populated.
// CEK is nil for AES-GCM; EPK is nil for RSA-OAEP-256 (no ephemeral key is generated).
type CryptoDetails struct {
	EPK gocrypto.PublicKey // ECDH-ES variants only; nil for RSA-OAEP-256 and AES-GCM
	CEK []byte             // Generated or derived CEK; nil for AES-GCM
}
