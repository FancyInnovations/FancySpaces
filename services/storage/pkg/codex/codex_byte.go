package codex

// EncodeByteInto encodes a byte into a byte slice.
// | Type (1 byte) | Byte Data (1 byte) |
func EncodeByteInto(val byte, target []byte) []byte {
	totalLength := 2

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	target[0] = byte(TypeByte)
	target[1] = val

	return target
}

// EncodeByte encodes a byte into a new byte slice.
// | Type (1 byte) | Byte Data (1 byte) |
func EncodeByte(val byte) []byte {
	return EncodeByteInto(val, nil)
}

// DecodeByte decodes a byte from a byte slice.
// | Type (1 byte) | Byte Data (1 byte) |
func DecodeByte(data []byte) (byte, error) {
	if len(data) < 2 {
		return 0, ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeByte) {
		return 0, ErrInvalidType
	}

	return data[1], nil
}
