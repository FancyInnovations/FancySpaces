package codex

import "encoding/binary"

// EncodeListInto encodes a list of Values into a byte slice.
// | Type (1 byte) | Item type (1 byte) | Count (2 bytes) | Payload length (4 bytes) | Items ... |
func EncodeListInto(vals []*Value, dst []byte) []byte {
	count := len(vals)
	itemType := TypeEmpty
	if count > 0 {
		itemType = vals[0].Type
	}

	totalLength := 1 + 1 + 2 + 4
	var itemPayload []byte
	for _, val := range vals {
		if val.Type != itemType && val.Type != TypeEmpty {
			return nil
		}
		itemPayload = append(itemPayload, EncodeValue(val)...)
	}
	totalLength += len(itemPayload)

	// ensure capacity
	if cap(dst) < totalLength {
		dst = make([]byte, totalLength)
	} else {
		dst = dst[:totalLength]
	}

	dst[0] = byte(TypeList)
	dst[1] = byte(itemType)
	binary.BigEndian.PutUint16(dst[2:4], uint16(count))
	binary.BigEndian.PutUint32(dst[4:8], uint32(len(itemPayload)))
	copy(dst[8:], itemPayload)

	return dst

}

// EncodeList encodes a list of Values into a new byte slice.
// | Type (1 byte) | Item type (1 byte) | Count (2 bytes) | Payload length (4 bytes) | Items ... |
func EncodeList(val []*Value) []byte {
	return EncodeListInto(val, nil)
}

// DecodeList decodes a list of Values from a byte slice.
// | Type (1 byte) | Item type (1 byte) | Count (2 bytes) | Payload length (4 bytes) | Items ... |
func DecodeList(data []byte) ([]*Value, error) {
	if len(data) < 8 {
		return nil, ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeList) {
		return nil, ErrInvalidType
	}

	itemType := data[1]
	count := int(binary.BigEndian.Uint16(data[2:4]))
	payloadLen := int(binary.BigEndian.Uint32(data[4:8]))

	if len(data) < 8+payloadLen {
		return nil, ErrPayloadTooShort
	}

	items := make([]*Value, 0, count)
	offset := 8
	for i := 0; i < count; i++ {
		itemData := data[offset:]
		item, err := DecodeValue(itemData)
		if err != nil {
			return nil, err
		}
		if item.Type != ValueType(itemType) && item.Type != TypeEmpty {
			return nil, ErrInvalidType
		}
		items = append(items, item)
		offset += len(EncodeValue(item)) // TODO: This is inefficient. We should track the length of the encoded item instead of re-encoding it to get the length.
	}

	return items, nil
}
