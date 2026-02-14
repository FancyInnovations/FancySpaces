package codex

import (
	"encoding/binary"
	"math"
)

// EncodeFloat64Into encodes a float64 into a byte slice.
// | Type (1 byte) | Float64 Data (8 bytes) |
func EncodeFloat64Into(val float64, target []byte) []byte {
	totalLength := 1 + 8

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	target[0] = byte(TypeFloat64)

	binary.BigEndian.PutUint64(target[1:9], math.Float64bits(val))

	return target
}

// EncodeFloat64 encodes a float64 into a new byte slice.
// | Type (1 byte) | Float64 Data (8 bytes) |
func EncodeFloat64(val float64) []byte {
	return EncodeFloat64Into(val, nil)
}

// DecodeFloat64 decodes a float64 from a byte slice.
// | Type (1 byte) | Float64 Data (8 bytes) |
func DecodeFloat64(data []byte) (float64, error) {
	if len(data) < 9 {
		return 0, ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeFloat64) {
		return 0, ErrInvalidType
	}

	bits := binary.BigEndian.Uint64(data[1:9])

	return math.Float64frombits(bits), nil
}

func SizeOfFloat64() uint64 {
	return 9
}
