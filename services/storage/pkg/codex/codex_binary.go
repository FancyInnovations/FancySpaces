package codex

import "encoding/binary"

// EncodeBinary encodes a byte slice into a new byte slice.
// | Type (1 byte) | Length (4 bytes) | Binary Data (N bytes) |
func EncodeBinary(val []byte) []byte {
	return EncodeBinaryInto(val, nil)
}

// DecodeBinary decodes a byte slice from a byte slice.
// | Type (1 byte) | Length (4 bytes) | Binary Data (N bytes) |
func DecodeBinary(data []byte) ([]byte, error) {
	if len(data) < 5 {
		return nil, ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeBinary) {
		return nil, ErrInvalidType
	}

	binLen := int(binary.BigEndian.Uint16(data[1:5]))
	if len(data) < 2+binLen {
		return nil, ErrPayloadTooShort
	}

	return data[5 : 5+binLen], nil
}

// EncodeBinaryInto encodes a byte slice into a byte slice.
// | Type (1 byte) | Length (4 bytes) | Binary Data (N bytes) |
func EncodeBinaryInto(val []byte, target []byte) []byte {
	valLength := len(val)
	totalLength := 1 + 2 + len(val)

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	target[0] = byte(TypeBinary)
	binary.BigEndian.PutUint16(target[1:5], uint16(valLength))
	copy(target[5:], val)

	return target
}
