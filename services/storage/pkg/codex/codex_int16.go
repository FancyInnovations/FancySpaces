package codex

import "encoding/binary"

// EncodeInt16Into encodes an int16 into a byte slice.
// | Type (1 byte) | int16 Data (2 bytes) |
func EncodeInt16Into(val int16, target []byte) []byte {
	totalLength := 1 + 2

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	target[0] = byte(TypeInt16)
	binary.BigEndian.PutUint16(target[1:3], uint16(val))

	return target
}

// EncodeInt16 encodes an int16 into a new byte slice.
// | Type (1 byte) | int16 Data (2 bytes) |
func EncodeInt16(val int16) []byte {
	return EncodeInt16Into(val, nil)
}

// DecodeInt16 decodes an int16 from a byte slice.
// | Type (1 byte) | int16 Data (2 bytes) |
func DecodeInt16(data []byte) (int16, error) {
	if len(data) < 3 {
		return 0, ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeInt16) {
		return 0, ErrInvalidType
	}

	return int16(binary.BigEndian.Uint16(data[1:3])), nil
}
