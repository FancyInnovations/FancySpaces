package codex

import "encoding/binary"

// EncodeListInto encodes a list of Values into a byte slice.
// | Type (1 byte) | Count (4 bytes) | | Item length (4 bytes) | Item data ... | ...
func EncodeListInto(vals []*Value, dst []byte) []byte {
	count := len(vals)

	totalLength := 1 + 4
	var itemPayload []byte
	for _, val := range vals {
		encodeValue := EncodeValue(val)

		// Item length (4 bytes)
		tmp4 := make([]byte, 4)
		binary.BigEndian.PutUint32(tmp4, uint32(len(encodeValue)))
		itemPayload = append(itemPayload, tmp4...)

		// Item data
		itemPayload = append(itemPayload, encodeValue...)
	}
	totalLength += len(itemPayload)

	// ensure capacity
	if cap(dst) < totalLength {
		dst = make([]byte, totalLength)
	} else {
		dst = dst[:totalLength]
	}

	dst[0] = byte(TypeList)
	binary.BigEndian.PutUint32(dst[1:5], uint32(count))
	copy(dst[5:], itemPayload)

	return dst

}

// EncodeList encodes a list of Values into a new byte slice.
// | Type (1 byte) | Count (4 bytes) | | Item length (4 bytes) | Item data ... | ...
func EncodeList(val []*Value) []byte {
	return EncodeListInto(val, nil)
}

// DecodeList decodes a list of Values from a byte slice.
// | Type (1 byte) | Count (4 bytes) | | Item length (4 bytes) | Item data ... | ...
func DecodeList(data []byte) ([]*Value, error) {
	if len(data) < 1 {
		return nil, ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeList) {
		return nil, ErrInvalidType
	}

	if len(data) < 5 {
		return nil, ErrPayloadTooShort
	}

	count := int(binary.BigEndian.Uint32(data[1:5]))
	if len(data) < 5+count*4 { // Header (5 bytes) + count * item length (4 bytes)
		return nil, ErrPayloadTooShort
	}

	items := make([]*Value, 0, count)
	offset := 5
	for i := 0; i < count; i++ {
		if len(data) < offset+4 { // Need at least 4 bytes for item length
			return nil, ErrPayloadTooShort
		}

		itemLength := int(binary.BigEndian.Uint32(data[offset : offset+4]))
		if len(data) < offset+itemLength {
			return nil, ErrPayloadTooShort
		}
		offset += 4 // Move past the item length

		itemData := data[offset : offset+itemLength]
		val, err := DecodeValue(itemData)
		if err != nil {
			return nil, err
		}
		items = append(items, val)

		offset += itemLength // Move to the next item
	}

	return items, nil
}
