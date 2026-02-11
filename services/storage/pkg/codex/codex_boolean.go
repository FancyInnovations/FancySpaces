package codex

// EncodeBoolInto encodes a bool into a byte slice.
// | Type (1 byte) | Bool Data (1 byte, 0x01=true) |
func EncodeBoolInto(val bool, target []byte) []byte {
	totalLength := 2

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	target[0] = byte(TypeBoolean)
	if val {
		target[1] = 1
	} else {
		target[1] = 0
	}

	return target
}

// EncodeBool encodes a bool into a new byte slice.
// | Type (1 byte) | Bool Data (1 byte, 0x01=true) |
func EncodeBool(val bool) []byte {
	return EncodeBoolInto(val, nil)
}

// DecodeBool decodes a bool from a byte slice.
// | Type (1 byte) | Bool Data (1 byte, 0x01=true) |
func DecodeBool(data []byte) (bool, error) {
	if len(data) < 2 {
		return false, ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeBoolean) {
		return false, ErrInvalidType
	}

	return data[1] == 1, nil
}

func SizeOfBool() uint64 {
	return 2
}
