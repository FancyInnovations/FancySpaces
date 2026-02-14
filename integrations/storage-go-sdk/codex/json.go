package codex

import (
	"encoding/base64"
	"fmt"
)

// ToAny converts the Value to a Go interface{} that can be easily marshaled to JSON.
func (v *Value) ToAny() (any, error) {
	switch v.Type {
	case TypeEmpty:
		return nil, nil

	case TypeBoolean:
		return v.AsBoolean(), nil

	case TypeByte:
		return v.AsByte(), nil

	case TypeUint16:
		return v.AsUint16(), nil

	case TypeUint32:
		return v.AsUint32(), nil

	case TypeUint64:
		return v.AsUint64(), nil

	case TypeInt16:
		return v.AsInt16(), nil

	case TypeInt32:
		return v.AsInt32(), nil

	case TypeInt64:
		return v.AsInt64(), nil

	case TypeFloat32:
		return v.AsFloat32(), nil

	case TypeFloat64:
		return v.AsFloat64(), nil

	case TypeString:
		return v.AsString(), nil

	case TypeBinary:
		// JSON does not support raw bytes — encode as base64 string
		return base64.StdEncoding.EncodeToString(v.AsBinary()), nil

	case TypeList:
		raw := v.data

		switch list := raw.(type) {
		case []*Value:
			result := make([]any, len(list))
			for i, item := range list {
				val, err := item.ToAny()
				if err != nil {
					return nil, err
				}
				result[i] = val
			}
			return result, nil

		case []Value:
			result := make([]any, len(list))
			for i := range list {
				val, err := list[i].ToAny()
				if err != nil {
					return nil, err
				}
				result[i] = val
			}
			return result, nil

		default:
			return nil, fmt.Errorf("invalid internal list type: %T", raw)
		}

	case TypeMap:
		raw := v.data

		m, ok := raw.(map[string]*Value)
		if !ok {
			return nil, fmt.Errorf("invalid internal map type: %T", raw)
		}

		result := make(map[string]any, len(m))
		for k, val := range m {
			converted, err := val.ToAny()
			if err != nil {
				return nil, err
			}
			result[k] = converted
		}

		return result, nil

	default:
		return nil, fmt.Errorf("unsupported value type: %v", v.Type)
	}
}
