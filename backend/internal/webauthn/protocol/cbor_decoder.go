// Package protocol implements WebAuthn protocol structures and parsing logic.
package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// CBOR Major Types
const (
	majorTypeUnsignedInt = 0
	majorTypeNegativeInt = 1
	majorTypeByteString  = 2
	majorTypeTextString  = 3
	majorTypeArray       = 4
	majorTypeMap         = 5
	majorTypeTag         = 6
	majorTypeSimple      = 7
)

// Simple Values
const (
	simpleValueFalse     = 20
	simpleValueTrue      = 21
	simpleValueNull      = 22
	simpleValueUndefined = 23
)

// cborDecoder is a simple CBOR decoder implementation.
type cborDecoder struct {
	data []byte
	pos  int
}

func newCBORDecoder(data []byte) *cborDecoder {
	return &cborDecoder{data: data, pos: 0}
}

func (d *cborDecoder) decode() (interface{}, error) {
	if d.pos >= len(d.data) {
		return nil, io.EOF
	}

	initialByte := d.data[d.pos]
	d.pos++
	majorType := initialByte >> 5
	additionalInfo := initialByte & 0x1f

	switch majorType {
	case majorTypeUnsignedInt:
		return d.decodeUnsignedInt(additionalInfo)
	case majorTypeNegativeInt:
		return d.decodeNegativeInt(additionalInfo)
	case majorTypeByteString:
		return d.decodeByteString(additionalInfo)
	case majorTypeTextString:
		return d.decodeTextString(additionalInfo)
	case majorTypeArray:
		return d.decodeArray(additionalInfo)
	case majorTypeMap:
		return d.decodeMap(additionalInfo)
	case majorTypeTag:
		return d.decodeTag(additionalInfo)
	case majorTypeSimple:
		return d.decodeSimple(additionalInfo)
	default:
		return nil, fmt.Errorf("unknown major type: %d", majorType)
	}
}

func (d *cborDecoder) decodeUnsignedInt(info byte) (interface{}, error) {
	val, err := d.readLength(info)
	if err != nil {
		return nil, err
	}
	if val > math.MaxInt64 {
		return nil, fmt.Errorf("unsigned integer too large for int64")
	}
	return int64(val), nil
}

func (d *cborDecoder) decodeNegativeInt(info byte) (interface{}, error) {
	val, err := d.readLength(info)
	if err != nil {
		return nil, err
	}
	if val > math.MaxInt64 {
		return nil, fmt.Errorf("negative integer too large for int64")
	}
	return -1 - int64(val), nil
}

func (d *cborDecoder) decodeByteString(info byte) (interface{}, error) {
	length, err := d.readLength(info)
	if err != nil {
		return nil, err
	}
	if length > uint64(math.MaxInt) {
		return nil, fmt.Errorf("byte string length too large")
	}
	if d.pos > len(d.data) {
		return nil, io.ErrUnexpectedEOF
	}
	remaining := len(d.data) - d.pos
	if length > uint64(remaining) {
		return nil, fmt.Errorf("byte string length exceeds data size")
	}
	bytes := make([]byte, length)
	copy(bytes, d.data[d.pos:d.pos+int(length)])
	d.pos += int(length)
	return bytes, nil
}

func (d *cborDecoder) decodeTextString(info byte) (interface{}, error) {
	length, err := d.readLength(info)
	if err != nil {
		return nil, err
	}
	if length > uint64(math.MaxInt) {
		return nil, fmt.Errorf("text string length too large")
	}
	if d.pos > len(d.data) {
		return nil, io.ErrUnexpectedEOF
	}
	remaining := len(d.data) - d.pos
	if length > uint64(remaining) {
		return nil, fmt.Errorf("text string length exceeds data size")
	}
	str := string(d.data[d.pos : d.pos+int(length)])
	d.pos += int(length)
	return str, nil
}

func (d *cborDecoder) decodeArray(info byte) (interface{}, error) {
	length, err := d.readLength(info)
	if err != nil {
		return nil, err
	}
	if length > uint64(math.MaxInt) {
		return nil, fmt.Errorf("array length too large")
	}
	if d.pos > len(d.data) {
		return nil, io.ErrUnexpectedEOF
	}
	remaining := len(d.data) - d.pos
	// Sanity check: each element takes at least 1 byte.
	if length > uint64(remaining) {
		return nil, fmt.Errorf("array length exceeds data size")
	}
	// Additional sanity check for max array size to prevent OOM
	if length > 65536 {
		return nil, fmt.Errorf("array length too large")
	}
	arr := make([]interface{}, int(length))
	for i := 0; i < int(length); i++ {
		elem, err := d.decode()
		if err != nil {
			return nil, err
		}
		arr[i] = elem
	}
	return arr, nil
}

func (d *cborDecoder) decodeMap(info byte) (interface{}, error) {
	length, err := d.readLength(info)
	if err != nil {
		return nil, err
	}
	// Map length check (length is number of pairs)
	if length > 65536 {
		return nil, fmt.Errorf("map length too large")
	}

	m := make(map[interface{}]interface{})
	for i := 0; i < int(length); i++ {
		key, err := d.decode()
		if err != nil {
			return nil, err
		}
		val, err := d.decode()
		if err != nil {
			return nil, err
		}
		m[key] = val
	}
	return m, nil
}

func (d *cborDecoder) decodeTag(info byte) (interface{}, error) {
	// Skip tag and decode content
	_, err := d.readLength(info)
	if err != nil {
		return nil, err
	}
	return d.decode()
}

func (d *cborDecoder) decodeSimple(info byte) (interface{}, error) {
	switch info {
	case simpleValueFalse:
		return false, nil
	case simpleValueTrue:
		return true, nil
	case simpleValueNull:
		return nil, nil
	case simpleValueUndefined:
		return nil, nil
	case 24: // Ignore simple value (1 byte)
		if d.pos >= len(d.data) {
			return nil, io.ErrUnexpectedEOF
		}
		d.pos++
		return nil, nil
	case 25: // Half-precision float (2 bytes)
		if d.pos+2 > len(d.data) {
			return nil, io.ErrUnexpectedEOF
		}
		// Not implemented, return float64(0)
		d.pos += 2
		return float64(0), nil
	case 26: // Single-precision float (4 bytes)
		if d.pos+4 > len(d.data) {
			return nil, io.ErrUnexpectedEOF
		}
		bits := binary.BigEndian.Uint32(d.data[d.pos:])
		d.pos += 4
		return math.Float32frombits(bits), nil
	case 27: // Double-precision float (8 bytes)
		if d.pos+8 > len(d.data) {
			return nil, io.ErrUnexpectedEOF
		}
		bits := binary.BigEndian.Uint64(d.data[d.pos:])
		d.pos += 8
		return math.Float64frombits(bits), nil
	default:
		return nil, fmt.Errorf("unsupported simple value: %d", info)
	}
}

func (d *cborDecoder) readLength(additionalInfo byte) (uint64, error) {
	if additionalInfo < 24 {
		return uint64(additionalInfo), nil
	}
	switch additionalInfo {
	case 24:
		if d.pos >= len(d.data) {
			return 0, io.ErrUnexpectedEOF
		}
		val := uint64(d.data[d.pos])
		d.pos++
		return val, nil
	case 25:
		if d.pos+2 > len(d.data) {
			return 0, io.ErrUnexpectedEOF
		}
		val := uint64(binary.BigEndian.Uint16(d.data[d.pos:]))
		d.pos += 2
		return val, nil
	case 26:
		if d.pos+4 > len(d.data) {
			return 0, io.ErrUnexpectedEOF
		}
		val := uint64(binary.BigEndian.Uint32(d.data[d.pos:]))
		d.pos += 4
		return val, nil
	case 27:
		if d.pos+8 > len(d.data) {
			return 0, io.ErrUnexpectedEOF
		}
		val := binary.BigEndian.Uint64(d.data[d.pos:])
		d.pos += 8
		return val, nil
	default:
		return 0, fmt.Errorf("invalid length encoding: %d", additionalInfo)
	}
}

// UnmarshalAttestationObject unmarshals CBOR data into AttestationObject struct.
func UnmarshalAttestationObject(data []byte, v *AttestationObject) error {
	decoder := newCBORDecoder(data)
	decoded, err := decoder.decode()
	if err != nil {
		return err
	}

	m, ok := decoded.(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf("attestation object must be a map")
	}

	// Populate struct fields
	if val, ok := m["fmt"]; ok {
		if s, ok := val.(string); ok {
			v.Format = s
		}
	}

	if val, ok := m["authData"]; ok {
		if b, ok := val.([]byte); ok {
			v.AuthData = b
		}
	}

	if val, ok := m["attStmt"]; ok {
		if stmtMap, ok := val.(map[interface{}]interface{}); ok {
			v.AttStatement = convertMapKeysToString(stmtMap)
		}
	}

	return nil
}

func convertMapKeysToString(m map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		if s, ok := k.(string); ok {
			// Recursively convert maps if needed
			if nestedMap, ok := v.(map[interface{}]interface{}); ok {
				result[s] = convertMapKeysToString(nestedMap)
			} else {
				result[s] = v
			}
		}
	}
	return result
}

// UnmarshalNext decodes the next CBOR data item and returns it along with the number of bytes consumed.
func UnmarshalNext(data []byte) (interface{}, int, error) {
	decoder := newCBORDecoder(data)
	val, err := decoder.decode()
	if err != nil {
		return nil, 0, err
	}
	return val, decoder.pos, nil
}
