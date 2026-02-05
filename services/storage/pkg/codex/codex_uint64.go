package codex

import "encoding/binary"

// EncodeUint64Into encodes an uint64 into a byte slice.
// | Type (1 byte) | uint64 Data (8 bytes) |
func EncodeUint64Into(val uint64, target []byte) []byte {
	totalLength := 1 + 8

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	target[0] = byte(TypeUint64)
	binary.BigEndian.PutUint64(target[1:9], val)

	return target
}

// EncodeUint64 encodes an uint32 into a new byte slice.
// | Type (1 byte) | uint64 Data (8 bytes) |
func EncodeUint64(val uint64) []byte {
	return EncodeUint64Into(val, nil)
}

// DecodeUint64 decodes an uint64 from a byte slice.
// | Type (1 byte) | uint64 Data (8 bytes) |
func DecodeUint64(data []byte) (uint64, error) {
	if len(data) < 9 {
		return 0, ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeUint64) {
		return 0, ErrInvalidType
	}

	return binary.BigEndian.Uint64(data[1:9]), nil
}
