package codex

import "encoding/binary"

// EncodeInt32Into encodes an int32 into a byte slice.
// | Type (1 byte) | int32 Data (4 bytes) |
func EncodeInt32Into(val int32, target []byte) []byte {
	totalLength := 1 + 4

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	target[0] = byte(TypeInt32)
	binary.BigEndian.PutUint32(target[1:5], uint32(val))

	return target
}

// EncodeInt32 encodes an int32 into a new byte slice.
// | Type (1 byte) | int32 Data (4 bytes) |
func EncodeInt32(val int32) []byte {
	return EncodeInt32Into(val, nil)
}

// DecodeInt32 decodes an int32 from a byte slice.
// | Type (1 byte) | int32 Data (4 bytes) |
func DecodeInt32(data []byte) (int32, error) {
	if len(data) < 5 {
		return 0, ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeInt32) {
		return 0, ErrInvalidType
	}

	return int32(binary.BigEndian.Uint32(data[1:5])), nil
}
