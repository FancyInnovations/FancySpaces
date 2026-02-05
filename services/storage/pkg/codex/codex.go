package codex

import (
	"encoding/binary"
	"math"
)

// EncodeValueInto encodes a Value into a byte slice.
// It uses the appropriate encoding function based on the ValueType.
// Returns the encoded byte slice.
func EncodeValueInto(val *Value, target []byte) []byte {
	switch val.Type {
	case TypeEmpty:
		return EncodeEmptyInto(target)
	case TypeBoolean:
		return EncodeBoolInto(val.AsBoolean(), target)
	case TypeByte:
		return EncodeByteInto(val.AsByte(), target)
	case TypeUint16:
		return EncodeUint16Into(val.AsUint16(), target)
	case TypeUint32:
		return EncodeUint32Into(val.AsUint32(), target)
	case TypeUint64:
		return EncodeUint64Into(val.AsUint64(), target)
	case TypeInt16:
		return EncodeInt16Into(val.AsInt16(), target)
	case TypeInt32:
		return EncodeInt32Into(val.AsInt32(), target)
	case TypeFloat64:
		return EncodeFloat64Into(val.AsFloat64(), target)
	case TypeBinary:
		return EncodeBinaryInto(val.AsBinary(), target)
	case TypeString:
		return EncodeStringInto(val.AsString(), target)
	case TypeList:
		return EncodeListInto(val.AsList(), target)
	case TypeMap:
		return EncodeMapInto(val.AsMap(), target)
	default:
		return nil
	}
}

// EncodeValue encodes a Value into a new byte slice.
// It uses the appropriate encoding function based on the ValueType.
// Returns the encoded byte slice.
func EncodeValue(val *Value) []byte {
	return EncodeValueInto(val, nil)
}

// DecodeValue decodes a Value from a byte slice.
// It uses the appropriate decoding function based on the ValueType.
// Returns the decoded Value.
func DecodeValue(data []byte) (*Value, error) {
	if len(data) == 0 {
		return nil, ErrPayloadTooShort
	}

	typeByte := data[0]
	switch ValueType(typeByte) {
	case TypeEmpty:
		return &Value{Type: TypeEmpty, data: nil}, nil
	case TypeBoolean:
		val, err := DecodeBool(data)
		if err != nil {
			return nil, err
		}
		return &Value{Type: TypeBoolean, data: val}, nil
	case TypeByte:
		val, err := DecodeByte(data)
		if err != nil {
			return nil, err
		}
		return &Value{Type: TypeByte, data: val}, nil
	case TypeUint16:
		val, err := DecodeUint16(data)
		if err != nil {
			return nil, err
		}
		return &Value{Type: TypeUint16, data: val}, nil
	case TypeUint32:
		val, err := DecodeUint32(data)
		if err != nil {
			return nil, err
		}
		return &Value{Type: TypeUint32, data: val}, nil
	case TypeUint64:
		val, err := DecodeUint64(data)
		if err != nil {
			return nil, err
		}
		return &Value{Type: TypeUint64, data: val}, nil
	case TypeInt16:
		val, err := DecodeInt16(data)
		if err != nil {
			return nil, err
		}
		return &Value{Type: TypeInt16, data: val}, nil
	case TypeInt32:
		val, err := DecodeInt32(data)
		if err != nil {
			return nil, err
		}
		return &Value{Type: TypeInt32, data: val}, nil
	case TypeFloat64:
		val, err := DecodeFloat64(data)
		if err != nil {
			return nil, err
		}
		return &Value{Type: TypeFloat64, data: val}, nil
	case TypeBinary:
		val, err := DecodeBinary(data)
		if err != nil {
			return nil, err
		}
		return &Value{Type: TypeBinary, data: val}, nil
	case TypeString:
		val, err := DecodeString(data)
		if err != nil {
			return nil, err
		}
		return &Value{Type: TypeString, data: val}, nil
	case TypeList:
		val, err := DecodeList(data)
		if err != nil {
			return nil, err
		}
		return &Value{Type: TypeList, data: val}, nil
	case TypeMap:
		val, err := DecodeMap(data)
		if err != nil {
			return nil, err
		}
		return &Value{Type: TypeMap, data: val}, nil
	default:
		return nil, ErrInvalidType
	}
}

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

// EncodeUint16Into encodes an uint16 into a byte slice.
// | Type (1 byte) | uint16 Data (2 bytes) |
func EncodeUint16Into(val uint16, target []byte) []byte {
	totalLength := 1 + 2

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	target[0] = byte(TypeUint16)
	binary.BigEndian.PutUint16(target[1:3], val)

	return target
}

// EncodeUint16 encodes an uint16 into a new byte slice.
// | Type (1 byte) | uint16 Data (2 bytes) |
func EncodeUint16(val uint16) []byte {
	return EncodeUint16Into(val, nil)
}

// DecodeUint16 decodes an uint16 from a byte slice.
// | Type (1 byte) | uint16 Data (2 bytes) |
func DecodeUint16(data []byte) (uint16, error) {
	if len(data) < 3 {
		return 0, ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeUint16) {
		return 0, ErrInvalidType
	}

	return binary.BigEndian.Uint16(data[1:3]), nil
}

// EncodeUint32Into encodes an uint32 into a byte slice.
// | Type (1 byte) | uint32 Data (4 bytes) |
func EncodeUint32Into(val uint32, target []byte) []byte {
	totalLength := 1 + 4

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	target[0] = byte(TypeUint32)
	binary.BigEndian.PutUint32(target[1:5], val)

	return target
}

// EncodeUint32 encodes an uint32 into a new byte slice.
// | Type (1 byte) | uint32 Data (4 bytes) |
func EncodeUint32(val uint32) []byte {
	return EncodeUint32Into(val, nil)
}

// DecodeUint32 decodes an uint32 from a byte slice.
// | Type (1 byte) | uint32 Data (4 bytes) |
func DecodeUint32(data []byte) (uint32, error) {
	if len(data) < 5 {
		return 0, ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeUint32) {
		return 0, ErrInvalidType
	}

	return binary.BigEndian.Uint32(data[1:5]), nil
}

// EncodeUint64Into encodes an uint64 into a byte slice.
// | Type (1 byte) | uint64 Data (8 bytes) |
func EncodeUint64Into(val uint64, target []byte) []byte {
	totalLength := 1 + 8

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	target[0] = byte(TypeUint64)
	binary.BigEndian.PutUint64(target[1:9], val)

	return target
}

// EncodeUint64 encodes an uint32 into a new byte slice.
// | Type (1 byte) | uint64 Data (8 bytes) |
func EncodeUint64(val uint64) []byte {
	return EncodeUint64Into(val, nil)
}

// DecodeUint64 decodes an uint64 from a byte slice.
// | Type (1 byte) | uint64 Data (8 bytes) |
func DecodeUint64(data []byte) (uint64, error) {
	if len(data) < 9 {
		return 0, ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeUint64) {
		return 0, ErrInvalidType
	}

	return binary.BigEndian.Uint64(data[1:9]), nil
}

// EncodeInt16Into encodes an int16 into a byte slice.
// | Type (1 byte) | int16 Data (2 bytes) |
func EncodeInt16Into(val int16, target []byte) []byte {
	totalLength := 1 + 2

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	target[0] = byte(TypeInt16)
	binary.BigEndian.PutUint16(target[1:3], uint16(val))

	return target
}

// EncodeInt16 encodes an int16 into a new byte slice.
// | Type (1 byte) | int16 Data (2 bytes) |
func EncodeInt16(val int16) []byte {
	return EncodeInt16Into(val, nil)
}

// DecodeInt16 decodes an int16 from a byte slice.
// | Type (1 byte) | int16 Data (2 bytes) |
func DecodeInt16(data []byte) (int16, error) {
	if len(data) < 3 {
		return 0, ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeInt16) {
		return 0, ErrInvalidType
	}

	return int16(binary.BigEndian.Uint16(data[1:3])), nil
}

// EncodeInt32Into encodes an int32 into a byte slice.
// | Type (1 byte) | int32 Data (4 bytes) |
func EncodeInt32Into(val int32, target []byte) []byte {
	totalLength := 1 + 4

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	target[0] = byte(TypeInt32)
	binary.BigEndian.PutUint32(target[1:5], uint32(val))

	return target
}

// EncodeInt32 encodes an int32 into a new byte slice.
// | Type (1 byte) | int32 Data (4 bytes) |
func EncodeInt32(val int32) []byte {
	return EncodeInt32Into(val, nil)
}

// DecodeInt32 decodes an int32 from a byte slice.
// | Type (1 byte) | int32 Data (4 bytes) |
func DecodeInt32(data []byte) (int32, error) {
	if len(data) < 5 {
		return 0, ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeInt32) {
		return 0, ErrInvalidType
	}

	return int32(binary.BigEndian.Uint32(data[1:5])), nil
}

// EncodeInt64Into encodes an int64 into a byte slice.
// | Type (1 byte) | int64 Data (8 bytes) |
func EncodeInt64Into(val int64, target []byte) []byte {
	totalLength := 1 + 8

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	target[0] = byte(TypeInt64)
	binary.BigEndian.PutUint64(target[1:9], uint64(val))

	return target
}

// EncodeInt64 encodes an int64 into a new byte slice.
// | Type (1 byte) | int64 Data (8 bytes) |
func EncodeInt64(val int64) []byte {
	return EncodeInt64Into(val, nil)
}

// DecodeInt64 decodes an int64 from a byte slice.
// | Type (1 byte) | int64 Data (8 bytes) |
func DecodeInt64(data []byte) (int64, error) {
	if len(data) < 9 {
		return 0, ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeInt64) {
		return 0, ErrInvalidType
	}

	return int64(binary.BigEndian.Uint64(data[1:9])), nil
}

// EncodeFloat32Into encodes a float32 into a byte slice.
// | Type (1 byte) | Float32 Data (4 bytes) |
func EncodeFloat32Into(val float32, target []byte) []byte {
	totalLength := 1 + 4

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	target[0] = byte(TypeFloat32)

	binary.BigEndian.PutUint32(target[1:5], math.Float32bits(val))

	return target
}

// EncodeFloat32 encodes a float32 into a new byte slice.
// | Type (1 byte) | Float32 Data (4 bytes) |
func EncodeFloat32(val float32) []byte {
	return EncodeFloat32Into(val, nil)
}

// DecodeFloat32 decodes a float32 from a byte slice.
// | Type (1 byte) | Float32 Data (4 bytes) |
func DecodeFloat32(data []byte) (float32, error) {
	if len(data) < 5 {
		return 0, ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeFloat32) {
		return 0, ErrInvalidType
	}

	bits := binary.BigEndian.Uint32(data[1:5])

	return math.Float32frombits(bits), nil
}

// EncodeFloat64Into encodes a float64 into a byte slice.
// | Type (1 byte) | Float64 Data (8 bytes) |
func EncodeFloat64Into(val float64, target []byte) []byte {
	totalLength := 1 + 8

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	target[0] = byte(TypeFloat64)

	binary.BigEndian.PutUint64(target[1:9], math.Float64bits(val))

	return target
}

// EncodeFloat64 encodes a float64 into a new byte slice.
// | Type (1 byte) | Float64 Data (8 bytes) |
func EncodeFloat64(val float64) []byte {
	return EncodeFloat64Into(val, nil)
}

// DecodeFloat64 decodes a float64 from a byte slice.
// | Type (1 byte) | Float64 Data (8 bytes) |
func DecodeFloat64(data []byte) (float64, error) {
	if len(data) < 9 {
		return 0, ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeFloat64) {
		return 0, ErrInvalidType
	}

	bits := binary.BigEndian.Uint64(data[1:9])

	return math.Float64frombits(bits), nil
}

// EncodeBinaryInto encodes a byte slice into a byte slice.
// | Type (1 byte) | Length (4 bytes) | Binary Data (N bytes) |
func EncodeBinaryInto(val []byte, target []byte) []byte {
	valLength := len(val)
	totalLength := 1 + 2 + len(val)

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	target[0] = byte(TypeBinary)
	binary.BigEndian.PutUint16(target[1:5], uint16(valLength))
	copy(target[5:], val)

	return target
}

// EncodeBinary encodes a byte slice into a new byte slice.
// | Type (1 byte) | Length (4 bytes) | Binary Data (N bytes) |
func EncodeBinary(val []byte) []byte {
	return EncodeBinaryInto(val, nil)
}

// DecodeBinary decodes a byte slice from a byte slice.
// | Type (1 byte) | Length (4 bytes) | Binary Data (N bytes) |
func DecodeBinary(data []byte) ([]byte, error) {
	if len(data) < 5 {
		return nil, ErrPayloadTooShort
	}

	typeByte := data[0]
	if typeByte != byte(TypeBinary) {
		return nil, ErrInvalidType
	}

	binLen := int(binary.BigEndian.Uint16(data[1:5]))
	if len(data) < 2+binLen {
		return nil, ErrPayloadTooShort
	}

	return data[5 : 5+binLen], nil
}

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
