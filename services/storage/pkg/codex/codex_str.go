package codex

import "encoding/binary"

// EncodeStringInto encodes a string into a byte slice.
// | Type (1 byte) | Length (4 bytes) | String Data (N bytes) |
func EncodeStringInto(val string, target []byte) []byte {
	valLength := len(val)
	totalLength := 1 + 4 + len(val)

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	target[0] = byte(TypeString)
	binary.BigEndian.PutUint16(target[1:5], uint16(valLength))
	copy(target[5:], val)

	return target
}

// EncodeString encodes a string into a new byte slice.
// | Type (1 byte) | Length (4 bytes) | String Data (N bytes) |
func EncodeString(val string) []byte {
	return EncodeStringInto(val, nil)
}

// DecodeString decodes a string from a byte slice.
// | Type (1 byte) | Length (4 bytes) | String Data (N bytes) |
func DecodeString(data []byte) (string, error) {
	if len(data) < 5 {
		return "", ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeString) {
		return "", ErrInvalidType
	}

	strLen := int(binary.BigEndian.Uint16(data[1:5]))
	if len(data) < 2+strLen {
		return "", ErrPayloadTooShort
	}

	return string(data[5 : 5+strLen]), nil
}
