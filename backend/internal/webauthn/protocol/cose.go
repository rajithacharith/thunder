package protocol

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"fmt"
	"math"
	"math/big"
)

const (
	// COSE Key Types
	coseKtyOKP = 1
	coseKtyEC2 = 2
	coseKtyRSA = 3

	// COSE Key Parameters
	coseKeyKty = 1
	coseKeyAlg = 3

	// EC2 Parameters
	coseKeyCrv = -1
	coseKeyX   = -2
	coseKeyY   = -3

	// RSA Parameters
	coseKeyN = -1
	coseKeyE = -2

	// Curves
	coseCrvP256 = 1
)

// ParseCOSEKey parses a COSE Key from raw bytes into a crypto.PublicKey and its algorithm.
func ParseCOSEKey(data []byte) (crypto.PublicKey, int64, error) {
	// Use our internal CBOR decoder
	val, _, err := UnmarshalNext(data)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to unmarshal COSE key: %w", err)
	}

	m, ok := val.(map[interface{}]interface{})
	if !ok {
		return nil, 0, fmt.Errorf("COSE key is not a map")
	}

	// Get Key Type (kty)
	ktyVal, ok := m[int64(coseKeyKty)]
	if !ok {
		// Try int key just in case decoder returns int
		ktyVal, ok = m[coseKeyKty]
		if !ok {
			return nil, 0, fmt.Errorf("missing kty")
		}
	}

	kty, ok := toInt64(ktyVal)
	if !ok {
		return nil, 0, fmt.Errorf("invalid kty type")
	}

	// Get Algorithm (alg)
	algVal, ok := m[int64(coseKeyAlg)]
	var alg int64
	if ok {
		alg, _ = toInt64(algVal)
	}

	var key crypto.PublicKey
	var parseErr error

	switch kty {
	case coseKtyEC2:
		key, parseErr = parseEC2Key(m)
	case coseKtyRSA:
		key, parseErr = parseRSAKey(m)
	default:
		return nil, 0, fmt.Errorf("unsupported kty: %d", kty)
	}

	return key, alg, parseErr
}

func parseEC2Key(m map[interface{}]interface{}) (crypto.PublicKey, error) {
	// Check curve
	crvVal, ok := m[int64(coseKeyCrv)]
	if !ok {
		return nil, fmt.Errorf("missing crv")
	}
	crv, ok := toInt64(crvVal)
	if !ok {
		return nil, fmt.Errorf("invalid crv type")
	}

	if crv != coseCrvP256 {
		return nil, fmt.Errorf("unsupported curve: %d", crv)
	}

	// Get X and Y coordinates
	xVal, ok := m[int64(coseKeyX)]
	if !ok {
		return nil, fmt.Errorf("missing x coordinate")
	}
	xBytes, ok := xVal.([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid x coordinate type")
	}

	yVal, ok := m[int64(coseKeyY)]
	if !ok {
		return nil, fmt.Errorf("missing y coordinate")
	}
	yBytes, ok := yVal.([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid y coordinate type")
	}

	key := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int).SetBytes(xBytes),
		Y:     new(big.Int).SetBytes(yBytes),
	}
	return key, nil
}

func parseRSAKey(m map[interface{}]interface{}) (crypto.PublicKey, error) {
	// Get Modulus (n)
	nVal, ok := m[int64(coseKeyN)]
	if !ok {
		return nil, fmt.Errorf("missing n")
	}
	nBytes, ok := nVal.([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid n type")
	}

	// Get Exponent (e)
	eVal, ok := m[int64(coseKeyE)]
	if !ok {
		return nil, fmt.Errorf("missing e")
	}

	var e int
	if eInt, ok := toInt64(eVal); ok {
		e = int(eInt)
	} else if eBytes, ok := eVal.([]byte); ok {
		e = int(new(big.Int).SetBytes(eBytes).Int64())
	} else {
		return nil, fmt.Errorf("invalid e type")
	}

	key := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: e,
	}
	return key, nil
}

func toInt64(i interface{}) (int64, bool) {
	switch v := i.(type) {
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	case uint:
		// Go's uint can be 32 or 64 bits.
		// If it fits in int64 (positive), it's safe.
		// MaxInt64 is 1<<63 - 1.
		// If uint is 64 bit and > MaxInt64, it overflows int64.
		if uint64(v) > uint64(math.MaxInt64) {
			return 0, false
		}
		return int64(v), true
	case uint8:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint64:
		if v > math.MaxInt64 {
			return 0, false
		}
		return int64(v), true
	default:
		return 0, false
	}
}
