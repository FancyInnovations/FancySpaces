package codex

import "encoding/binary"

// EncodeUint16Into encodes an uint16 into a byte slice.
// | Type (1 byte) | uint16 Data (2 bytes) |
func EncodeUint16Into(val uint16, target []byte) []byte {
	totalLength := 1 + 2

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	target[0] = byte(TypeUint16)
	binary.BigEndian.PutUint16(target[1:3], val)

	return target
}

// EncodeUint16 encodes an uint16 into a new byte slice.
// | Type (1 byte) | uint16 Data (2 bytes) |
func EncodeUint16(val uint16) []byte {
	return EncodeUint16Into(val, nil)
}

// DecodeUint16 decodes an uint16 from a byte slice.
// | Type (1 byte) | uint16 Data (2 bytes) |
func DecodeUint16(data []byte) (uint16, error) {
	if len(data) < 3 {
		return 0, ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeUint16) {
		return 0, ErrInvalidType
	}

	return binary.BigEndian.Uint16(data[1:3]), nil
}
