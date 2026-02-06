package codex

import "encoding/binary"

// EncodeStringInto encodes a string into a byte slice.
// | Type (1 byte) | Length (4 bytes) | String Data (N bytes) |
func EncodeStringInto(val string, target []byte) []byte {
	valLength := len(val)
	totalLength := 1 + 4 + valLength

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	target[0] = byte(TypeString)
	binary.BigEndian.PutUint32(target[1:5], uint32(valLength))
	copy(target[5:], val)

	return target
}

// EncodeString encodes a string into a new byte slice.
func EncodeString(val string) []byte {
	return EncodeStringInto(val, nil)
}

// DecodeString decodes a string from a byte slice.
// | Type (1 byte) | Length (4 bytes) | String Data (N bytes) |
func DecodeString(data []byte) (string, error) {
	if len(data) < 5 {
		return "", ErrPayloadTooShort
	}

	if data[0] != byte(TypeString) {
		return "", ErrInvalidType
	}

	strLen := int(binary.BigEndian.Uint32(data[1:5]))
	if strLen < 0 || len(data) < 5+strLen {
		return "", ErrPayloadTooShort
	}

	return string(data[5 : 5+strLen]), nil
}
