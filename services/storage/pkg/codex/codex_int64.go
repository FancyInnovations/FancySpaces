package codex

import "encoding/binary"

// EncodeInt64Into encodes an int64 into a byte slice.
// | Type (1 byte) | int64 Data (8 bytes) |
func EncodeInt64Into(val int64, target []byte) []byte {
	totalLength := 1 + 8

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	target[0] = byte(TypeInt64)
	binary.BigEndian.PutUint64(target[1:9], uint64(val))

	return target
}

// EncodeInt64 encodes an int64 into a new byte slice.
// | Type (1 byte) | int64 Data (8 bytes) |
func EncodeInt64(val int64) []byte {
	return EncodeInt64Into(val, nil)
}

// DecodeInt64 decodes an int64 from a byte slice.
// | Type (1 byte) | int64 Data (8 bytes) |
func DecodeInt64(data []byte) (int64, error) {
	if len(data) < 9 {
		return 0, ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeInt64) {
		return 0, ErrInvalidType
	}

	return int64(binary.BigEndian.Uint64(data[1:9])), nil
}
