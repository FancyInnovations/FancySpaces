package codex

// EncodeEmptyInto encodes an empty value into a byte slice.
// | Type (1 byte) |
func EncodeEmptyInto(target []byte) []byte {
	totalLength := 1

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	target[0] = byte(TypeEmpty)

	return target
}

// EncodeEmpty encodes an empty value into a new byte slice.
// | Type (1 byte) |
func EncodeEmpty() []byte {
	return EncodeEmptyInto(nil)
}
