package codex

// EncodeBoolInto encodes a bool into a byte slice.
// | Type (1 byte) | Bool Data (1 byte, 0x01=true) |
func EncodeBoolInto(val bool, target []byte) []byte {
	var byteVal byte
	if val {
		byteVal = 1
	} else {
		byteVal = 0
	}
	return EncodeByteInto(byteVal, target)
}

// EncodeBool encodes a bool into a new byte slice.
// | Type (1 byte) | Bool Data (1 byte, 0x01=true) |
func EncodeBool(val bool) []byte {
	return EncodeBoolInto(val, nil)
}

// DecodeBool decodes a bool from a byte slice.
// | Type (1 byte) | Bool Data (1 byte, 0x01=true) |
func DecodeBool(data []byte) (bool, error) {
	byteVal, err := DecodeByte(data)
	if err != nil {
		return false, err
	}

	return byteVal == 1, nil
}
