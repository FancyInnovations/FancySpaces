package codex

import "encoding/binary"

// EncodeMapInto encodes a map of string keys to Values into a byte slice.
// | Type (1 byte) | Count (4 bytes) | Key length (2 bytes) | Key (N bytes) | Val length (4 bytes) | Val (M bytes) | ... |
func EncodeMapInto(vals map[string]*Value, dst []byte) []byte {
	count := len(vals)

	totalLength := 1 + 4

	var itemPayload []byte
	for key, val := range vals {
		// Key
		keyBytes := []byte(key)
		tmp2 := make([]byte, 2)
		binary.BigEndian.PutUint16(tmp2, uint16(len(keyBytes)))
		itemPayload = append(itemPayload, tmp2...)
		itemPayload = append(itemPayload, keyBytes...)

		// Value
		valBytes := EncodeValue(val)
		tmp4 := make([]byte, 4)
		binary.BigEndian.PutUint32(tmp4, uint32(len(valBytes)))
		itemPayload = append(itemPayload, tmp4...)
		itemPayload = append(itemPayload, valBytes...)
	}
	totalLength += len(itemPayload)

	// ensure capacity
	if cap(dst) < totalLength {
		dst = make([]byte, totalLength)
	} else {
		dst = dst[:totalLength]
	}

	dst[0] = byte(TypeMap)
	binary.BigEndian.PutUint32(dst[1:5], uint32(count))

	if len(itemPayload) > 0 {
		copy(dst[5:], itemPayload)
	}

	return dst
}

// EncodeMap encodes a map of string keys to Values into a new byte slice.
// | Type (1 byte) | Count (4 bytes) | Key length (2 bytes) | Key (N bytes) | Val length (4 bytes) | Val (M bytes) | ... |
func EncodeMap(vals map[string]*Value) []byte {
	return EncodeMapInto(vals, nil)
}

// DecodeMap decodes a byte slice into a map of string keys to Values.
func DecodeMap(data []byte) (map[string]*Value, error) {
	if len(data) < 1 {
		return nil, ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeMap) {
		return nil, ErrInvalidType
	}

	if len(data) < 5 { // header (1 byte) + count (4 bytes)
		return nil, ErrPayloadTooShort
	}

	count := int(binary.BigEndian.Uint32(data[1:5]))
	if len(data) < 5+count*(2+4) { // header (5 bytes) + min key (2 bytes) + min value (4 bytes)
		return nil, ErrPayloadTooShort
	}

	vals := make(map[string]*Value)
	if count == 0 {
		return vals, nil
	}

	offset := 5
	for i := 0; i < count; i++ {
		if len(data[offset:]) < 2 {
			return nil, ErrPayloadTooShort
		}

		// Key
		keyLen := int(binary.BigEndian.Uint16(data[offset : offset+2]))
		offset += 2

		if len(data[offset:]) < keyLen {
			return nil, ErrPayloadTooShort
		}
		key := string(data[offset : offset+keyLen])
		offset += keyLen

		// Value
		if len(data[offset:]) < 4 {
			return nil, ErrPayloadTooShort
		}
		valLen := int(binary.BigEndian.Uint32(data[offset : offset+4]))
		offset += 4

		if len(data[offset:]) < valLen {
			return nil, ErrPayloadTooShort
		}

		val, err := DecodeValue(data[offset : offset+valLen])
		if err != nil {
			return nil, err
		}
		offset += valLen

		vals[key] = val
	}

	return vals, nil
}

func SizeOfMap(vals map[string]*Value) uint64 {
	size := uint64(1 + 4) // Type (1 byte) + Count (4 bytes)

	for key, val := range vals {
		size += 2                // Key length (2 bytes)
		size += uint64(len(key)) // Key (N bytes)
		size += 4                // Val length (4 bytes)
		size += SizeOfValue(val) // Val (M bytes)
	}

	return size
}
