package codex

import (
	"fmt"
	"log/slog"
)

type Value struct {
	Type ValueType
	data any
}

type ValueType uint8

const (
	TypeEmpty   ValueType = 0
	TypeBoolean ValueType = 1
	TypeByte    ValueType = 2
	TypeUint16  ValueType = 3
	TypeUint32  ValueType = 4
	TypeUint64  ValueType = 5
	TypeInt16   ValueType = 6
	TypeInt32   ValueType = 7
	TypeInt64   ValueType = 8
	TypeFloat32 ValueType = 9
	TypeFloat64 ValueType = 10
	TypeBinary  ValueType = 11
	TypeString  ValueType = 12
	TypeList    ValueType = 13
	TypeMap     ValueType = 14
)

var EmptyValue = &Value{Type: TypeEmpty, data: nil}

func (v *Value) IsEmpty() bool {
	return v.Type == TypeEmpty
}

func (v *Value) IsNumeric() bool {
	switch v.Type {
	case TypeByte, TypeUint16, TypeUint32, TypeUint64,
		TypeInt16, TypeInt32, TypeInt64,
		TypeFloat32, TypeFloat64:
		return true
	default:
		return false
	}
}

func (v *Value) AsBoolean() bool {
	if v.Type != TypeBoolean {
		slog.Warn("Value is not a boolean", slog.Any("value_type", v.Type))
	}

	return v.data.(bool)
}

func (v *Value) AsByte() byte {
	if v.Type != TypeByte {
		slog.Warn("Value is not a byte", slog.Any("value_type", v.Type))
	}
	return v.data.(byte)
}

func (v *Value) AsUint8() uint8 {
	if v.Type != TypeByte {
		slog.Warn("Value is not a byte (required for AsUint8)", slog.Any("value_type", v.Type))
	}
	return uint8(v.data.(byte))
}

func (v *Value) AsUint16() uint16 {
	if v.Type != TypeUint16 {
		slog.Warn("Value is not a uint16", slog.Any("value_type", v.Type))
	}
	return v.data.(uint16)
}

func (v *Value) AsUint32() uint32 {
	if v.Type != TypeUint32 {
		slog.Warn("Value is not a uint32", slog.Any("value_type", v.Type))
	}
	return v.data.(uint32)
}

func (v *Value) AsUint() uint {
	if v.Type != TypeUint32 {
		slog.Warn("Value is not a uint32 (required for AsUint)", slog.Any("value_type", v.Type))
	}
	return uint(v.data.(uint32))
}

func (v *Value) AsUint64() uint64 {
	if v.Type != TypeUint64 {
		slog.Warn("Value is not a uint64", slog.Any("value_type", v.Type))
	}
	return v.data.(uint64)
}

func (v *Value) AsInt16() int16 {
	if v.Type != TypeInt16 {
		slog.Warn("Value is not an int16", slog.Any("value_type", v.Type))
	}
	return v.data.(int16)
}

func (v *Value) AsInt32() int32 {
	if v.Type != TypeInt32 {
		slog.Warn("Value is not an int32", slog.Any("value_type", v.Type))
	}
	return v.data.(int32)
}

func (v *Value) AsInt() int {
	if v.Type != TypeInt32 {
		slog.Warn("Value is not an int32 (required for AsInt)", slog.Any("value_type", v.Type))
	}
	return int(v.data.(int32))
}

func (v *Value) AsInt64() int64 {
	if v.Type != TypeInt64 {
		slog.Warn("Value is not an int64", slog.Any("value_type", v.Type))
	}
	return v.data.(int64)
}

func (v *Value) AsFloat32() float32 {
	if v.Type != TypeFloat32 {
		slog.Warn("Value is not a float32", slog.Any("value_type", v.Type))
	}
	return v.data.(float32)
}

func (v *Value) AsFloat64() float64 {
	if v.Type != TypeFloat64 {
		slog.Warn("Value is not a float64", slog.Any("value_type", v.Type))
	}
	return v.data.(float64)
}

func (v *Value) AsBinary() []byte {
	if v.Type != TypeBinary {
		slog.Warn("Value is not binary data", slog.Any("value_type", v.Type))
	}
	return v.data.([]byte)
}

func (v *Value) AsString() string {
	if v.Type != TypeString {
		slog.Warn("Value is not a string", slog.Any("value_type", v.Type))
	}
	return v.data.(string)
}

func (v *Value) AsList() []*Value {
	if v.Type != TypeList {
		slog.Warn("Value is not a list", slog.Any("value_type", v.Type))
	}
	return v.data.([]*Value)
}

func (v *Value) AsMap() map[string]*Value {
	if v.Type != TypeMap {
		slog.Warn("Value is not a map", slog.Any("value_type", v.Type))
	}
	return v.data.(map[string]*Value)
}

func NewValue(data any) (*Value, error) {
	switch v := data.(type) {
	case nil:
		return &Value{Type: TypeEmpty, data: nil}, nil
	case bool:
		return &Value{Type: TypeBoolean, data: v}, nil
	case byte:
		return &Value{Type: TypeByte, data: v}, nil
	case uint16:
		return &Value{Type: TypeUint16, data: v}, nil
	case uint32:
		return &Value{Type: TypeUint32, data: v}, nil
	case uint64:
		return &Value{Type: TypeUint64, data: v}, nil
	case int16:
		return &Value{Type: TypeInt16, data: v}, nil
	case int32:
		return &Value{Type: TypeInt32, data: v}, nil
	case int:
		return &Value{Type: TypeInt32, data: int32(v)}, nil
	case int64:
		return &Value{Type: TypeInt64, data: v}, nil
	case float32:
		return &Value{Type: TypeFloat32, data: v}, nil
	case float64:
		return &Value{Type: TypeFloat64, data: v}, nil
	case []byte:
		return &Value{Type: TypeBinary, data: v}, nil
	case string:
		return &Value{Type: TypeString, data: v}, nil
	case []any:
		itemType := TypeEmpty
		values := make([]Value, len(v))
		for i, item := range v {
			value, err := NewValue(item)
			if err != nil {
				return nil, err
			}

			// validate that all items in the list have the same type
			if i == 0 {
				itemType = value.Type
			} else if value.Type != itemType {
				return nil, ErrInvalidType
			}

			values[i] = *value
		}
		return &Value{Type: TypeList, data: values}, nil
	case map[string]any:
		itemType := TypeEmpty
		values := make(map[string]Value)
		for key, item := range v {
			value, err := NewValue(item)
			if err != nil {
				return nil, err
			}

			// validate that all values in the map have the same type
			if itemType == TypeEmpty {
				itemType = value.Type
			} else if value.Type != itemType {
				return nil, ErrInvalidType
			}

			values[key] = *value
		}
		return &Value{Type: TypeMap, data: values}, nil
	default:
		return nil, fmt.Errorf("unsupported value type: %T", data)
	}
}

func NewEmptyValue() *Value {
	return &Value{Type: TypeEmpty, data: nil}
}

func NewStringListValue(items []string) *Value {
	values := make([]*Value, len(items))
	for i, item := range items {
		values[i] = &Value{Type: TypeString, data: item}
	}
	return &Value{Type: TypeList, data: values}
}
