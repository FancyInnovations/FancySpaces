package codex

import "encoding/binary"

// EncodeUint32Into encodes an uint32 into a byte slice.
// | Type (1 byte) | uint32 Data (4 bytes) |
func EncodeUint32Into(val uint32, target []byte) []byte {
	totalLength := 1 + 4

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	target[0] = byte(TypeUint32)
	binary.BigEndian.PutUint32(target[1:5], val)

	return target
}

// EncodeUint32 encodes an uint32 into a new byte slice.
// | Type (1 byte) | uint32 Data (4 bytes) |
func EncodeUint32(val uint32) []byte {
	return EncodeUint32Into(val, nil)
}

// DecodeUint32 decodes an uint32 from a byte slice.
// | Type (1 byte) | uint32 Data (4 bytes) |
func DecodeUint32(data []byte) (uint32, error) {
	if len(data) < 5 {
		return 0, ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeUint32) {
		return 0, ErrInvalidType
	}

	return binary.BigEndian.Uint32(data[1:5]), nil
}

func SizeOfUint32() uint64 {
	return 5
}
