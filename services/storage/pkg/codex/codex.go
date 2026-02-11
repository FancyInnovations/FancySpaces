package codex

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
	case TypeInt64:
		return EncodeInt64Into(val.AsInt64(), target)
	case TypeFloat32:
		return EncodeFloat32Into(val.AsFloat32(), target)
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
	case TypeInt64:
		val, err := DecodeInt64(data)
		if err != nil {
			return nil, err
		}
		return &Value{Type: TypeInt64, data: val}, nil
	case TypeFloat32:
		val, err := DecodeFloat32(data)
		if err != nil {
			return nil, err
		}
		return &Value{Type: TypeFloat32, data: val}, nil
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

func SizeOfValue(val *Value) uint64 {
	switch val.Type {
	case TypeEmpty:
		return SizeOfEmpty()
	case TypeBoolean:
		return SizeOfBool()
	case TypeByte:
		return SizeOfByte()
	case TypeUint16:
		return SizeOfUint16()
	case TypeUint32:
		return SizeOfUint32()
	case TypeUint64:
		return SizeOfUint64()
	case TypeInt16:
		return SizeOfInt16()
	case TypeInt32:
		return SizeOfInt32()
	case TypeInt64:
		return SizeOfInt64()
	case TypeFloat32:
		return SizeOfFloat32()
	case TypeFloat64:
		return SizeOfFloat64()
	case TypeBinary:
		return SizeOfBinary(val.AsBinary())
	case TypeString:
		return SizeOfString(val.AsString())
	case TypeList:
		return SizeOfList(val.AsList())
	case TypeMap:
		return SizeOfMap(val.AsMap())
	default:
		return 0
	}
}
