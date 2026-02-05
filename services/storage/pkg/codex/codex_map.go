package codex

import "encoding/binary"

// EncodeMapInto encodes a map of string keys to Values into a byte slice.
// | Type (1 byte) | Val type (1 byte) | Count (2 bytes) | Payload length (4 bytes) | Key length (2 bytes) | Key (N bytes) | Val (M bytes) | ... |
func EncodeMapInto(vals map[string]*Value, dst []byte) []byte {
	count := len(vals)
	valType := TypeEmpty
	if count > 0 {
		for _, val := range vals {
			valType = val.Type
			break
		}
	}

	totalLength := 1 + 1 + 2 + 4
	var itemPayload []byte
	for key, val := range vals {
		if val.Type != valType && val.Type != TypeEmpty {
			return nil
		}

		itemPayload = append(itemPayload, EncodeString(key)...)
		itemPayload = append(itemPayload, EncodeValue(val)...)
	}
	totalLength += len(itemPayload)

	// ensure capacity
	if cap(dst) < totalLength {
		dst = make([]byte, totalLength)
	} else {
		dst = dst[:totalLength]
	}

	dst[0] = byte(TypeMap)
	dst[1] = byte(valType)
	binary.BigEndian.PutUint16(dst[2:4], uint16(count))
	binary.BigEndian.PutUint32(dst[4:8], uint32(len(itemPayload)))
	copy(dst[8:], itemPayload)

	return dst
}

// EncodeMap encodes a map of string keys to Values into a new byte slice.
// | Type (1 byte) | Val type (1 byte) | Count (2 bytes) | Payload length (4 bytes) | Key length (2 bytes) | Key (N bytes) | Val (M bytes) | ... |
func EncodeMap(vals map[string]*Value) []byte {
	return EncodeMapInto(vals, nil)
}

func DecodeMap(data []byte) (map[string]*Value, error) {
	if len(data) < 8 {
		return nil, ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeMap) {
		return nil, ErrInvalidType
	}

	valType := data[1]
	count := int(binary.BigEndian.Uint16(data[2:4]))
	payloadLen := int(binary.BigEndian.Uint32(data[4:8]))

	if len(data) < 8+payloadLen {
		return nil, ErrPayloadTooShort
	}

	vals := make(map[string]*Value)
	offset := 8
	for i := 0; i < count; i++ {
		keyData := data[offset:]
		key, err := DecodeString(keyData)
		if err != nil {
			return nil, err
		}
		offset += len(EncodeString(key))

		valData := data[offset:]
		val, err := DecodeValue(valData)
		if err != nil {
			return nil, err
		}
		if val.Type != ValueType(valType) && val.Type != TypeEmpty {
			return nil, ErrInvalidType
		}
		vals[key] = val
		offset += len(EncodeValue(val)) // TODO: This is inefficient. We should track the length of the encoded value instead of re-encoding it to get the length.
	}

	return vals, nil
}
