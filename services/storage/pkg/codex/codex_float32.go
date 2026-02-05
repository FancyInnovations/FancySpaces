package codex

import (
	"encoding/binary"
	"math"
)

// EncodeFloat32Into encodes a float32 into a byte slice.
// | Type (1 byte) | Float32 Data (4 bytes) |
func EncodeFloat32Into(val float32, target []byte) []byte {
	totalLength := 1 + 4

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	target[0] = byte(TypeFloat32)

	binary.BigEndian.PutUint32(target[1:5], math.Float32bits(val))

	return target
}

// EncodeFloat32 encodes a float32 into a new byte slice.
// | Type (1 byte) | Float32 Data (4 bytes) |
func EncodeFloat32(val float32) []byte {
	return EncodeFloat32Into(val, nil)
}

// DecodeFloat32 decodes a float32 from a byte slice.
// | Type (1 byte) | Float32 Data (4 bytes) |
func DecodeFloat32(data []byte) (float32, error) {
	if len(data) < 5 {
		return 0, ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeFloat32) {
		return 0, ErrInvalidType
	}

	bits := binary.BigEndian.Uint32(data[1:5])

	return math.Float32frombits(bits), nil
}
