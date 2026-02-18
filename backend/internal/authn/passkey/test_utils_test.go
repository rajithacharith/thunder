package passkey

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/require"
)

// generateTestKeys generates an ECDSA key pair and returns the private key
// and the COSE encoded public key bytes.
func generateTestKeys(t *testing.T) (*ecdsa.PrivateKey, []byte) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	publicKey := privateKey.PublicKey

	// COSE Key Map
	// 1 (kty): 2 (EC2)
	// 3 (alg): -7 (ES256)
	// -1 (crv): 1 (P-256)
	// -2 (x): x-coord
	// -3 (y): y-coord

	// Manually encode to CBOR bytes to avoid using a library
	// Map(5 items) -> 0xa5
	// 1 (0x01) -> 2 (0x02)
	// 3 (0x03) -> -7 (0x26)
	// -1 (0x20) -> 1 (0x01)
	// -2 (0x21) -> bytes(32)
	// -3 (0x22) -> bytes(32)

	xBytes := publicKey.X.Bytes()
	yBytes := publicKey.Y.Bytes()

	// Pad X and Y to 32 bytes if needed
	if len(xBytes) < 32 {
		pad := make([]byte, 32-len(xBytes))
		xBytes = append(pad, xBytes...)
	}
	if len(yBytes) < 32 {
		pad := make([]byte, 32-len(yBytes))
		yBytes = append(pad, yBytes...)
	}

	coseKey := []byte{}
	coseKey = append(coseKey, 0xa5) // map(5)

	// kty: EC2
	coseKey = append(coseKey, 0x01, 0x02)

	// alg: ES256
	coseKey = append(coseKey, 0x03, 0x26) // 0x26 is -7 (0x20 | 6) -> 001 00110 = -1-6 = -7. Wait.
	// Major type 1 (negative): -1 - val.
	// -7 = -1 - 6. So val is 6.
	// 001 00110 -> 0x26. Correct.

	// crv: P-256
	coseKey = append(coseKey, 0x20, 0x01) // key -1 (0x20), val 1 (0x01)

	// x
	coseKey = append(coseKey, 0x21)       // key -2
	coseKey = append(coseKey, 0x58, 0x20) // bytes(32) -> 0x58 (header) 0x20 (len)
	coseKey = append(coseKey, xBytes...)

	// y
	coseKey = append(coseKey, 0x22)       // key -3
	coseKey = append(coseKey, 0x58, 0x20) // bytes(32)
	coseKey = append(coseKey, yBytes...)

	return privateKey, coseKey
}

// createTestAuthData creates a valid authenticator data byte array.
func createTestAuthData(rpID string) []byte {
	rpIDHash := sha256.Sum256([]byte(rpID))

	authData := make([]byte, 0, 37)
	authData = append(authData, rpIDHash[:]...) // 32 bytes

	// Flags: User Present (0x01) + User Verified (0x04) = 0x05
	authData = append(authData, 0x05)

	// Sign Count: 0
	authData = append(authData, 0x00, 0x00, 0x00, 0x00)

	return authData
}

// signData signs the data using the private key.
func signData(t *testing.T, privateKey *ecdsa.PrivateKey, data []byte) []byte {
	hash := sha256.Sum256(data)
	signature, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	require.NoError(t, err)
	return signature
}
