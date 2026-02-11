package codex

import (
	"fmt"
	"reflect"
	"time"
)

// MarshalToMap converts a struct into a map[string]*Value.
// The struct fields can be tagged with `json:"fieldName"` to specify the corresponding key in the map.
// If a field is not tagged, it will use the field name as the key.
// Nested structs and maps/lists are supported.
func MarshalToMap(s any) (map[string]*Value, error) {
	switch s.(type) {
	case time.Time, *time.Time:
		// Special handling for time.Time values
		return timeToValue(reflect.ValueOf(s).Interface().(time.Time)), nil
	}

	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct or pointer to struct")
	}

	out := make(map[string]*Value)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Only exportable fields can be marshaled
		if !fieldType.IsExported() {
			continue
		}

		if tag := fieldType.Tag.Get("json"); tag != "" && tag != "-" {
			fieldType.Name = tag
		}

		value, err := toValue(field)
		if err != nil {
			return nil, fmt.Errorf("field %s: %w", fieldType.Name, err)
		}

		out[fieldType.Name] = value
	}

	return out, nil
}

// Marshal converts a struct into a byte slice.
// The struct fields can be tagged with `json:"fieldName"` to specify the corresponding key in the map.
// If a field is not tagged, it will use the field name as the key.
// Nested structs and maps/lists are supported.
func Marshal(s any) ([]byte, error) {
	m, err := MarshalToMap(s)
	if err != nil {
		return nil, err
	}

	return EncodeMap(m), nil
}

func toValue(v reflect.Value) (*Value, error) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return NewEmptyValue(), nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Bool:
		return NewValue(v.Bool())
	case reflect.Uint8:
		return NewValue(byte(v.Uint()))
	case reflect.Uint16:
		return NewValue(uint16(v.Uint()))
	case reflect.Uint32, reflect.Uint:
		return NewValue(uint32(v.Uint()))
	case reflect.Uint64:
		return NewValue(v.Uint())
	case reflect.Int16:
		return NewValue(int16(v.Int()))
	case reflect.Int32, reflect.Int:
		return NewValue(int32(v.Int()))
	case reflect.Int64:
		return NewValue(v.Int())
	case reflect.Float32:
		return NewValue(float32(v.Float()))
	case reflect.Float64:
		return NewValue(v.Float())
	case reflect.String:
		return NewValue(v.String())
	case reflect.Slice, reflect.Array:
		l := v.Len()
		values := make([]*Value, l)
		for i := 0; i < l; i++ {
			elemVal, err := toValue(v.Index(i))
			if err != nil {
				return nil, fmt.Errorf("index %d: %w", i, err)
			}
			values[i] = elemVal
		}
		return &Value{Type: TypeList, data: values}, nil
	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return nil, fmt.Errorf("map key must be string")
		}
		iter := v.MapRange()
		m := make(map[string]*Value)
		for iter.Next() {
			key := iter.Key().String()
			val, err := toValue(iter.Value())
			if err != nil {
				return nil, fmt.Errorf("map key %s: %w", key, err)
			}
			m[key] = val
		}
		return &Value{Type: TypeMap, data: m}, nil
	case reflect.Struct:
		// Nested struct as map
		m, err := MarshalToMap(v.Interface())
		if err != nil {
			return nil, err
		}
		return &Value{Type: TypeMap, data: m}, nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", v.Kind())
	}
}

// UnmarshalFromMap populates a struct from a map[string]*Value.
// The struct fields can be tagged with `json:"fieldName"` to specify the corresponding key in the map.
// If a field is not tagged, it will match the map key with the same name as the field.
// Nested structs and maps/lists are supported.
func UnmarshalFromMap(m map[string]*Value, target any) error {
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct")
	}

	structVal := val.Elem()
	structType := structVal.Type()

	for i := 0; i < structVal.NumField(); i++ {
		field := structVal.Field(i)
		fieldType := structType.Field(i)
		if !fieldType.IsExported() {
			continue
		}

		fieldName := fieldType.Name

		if tag := fieldType.Tag.Get("json"); tag != "" && tag != "-" {
			fieldName = tag
		}

		value, ok := m[fieldName]
		if !ok || value.IsEmpty() {
			continue
		}

		// Special handling for time.Time fields
		if field.Type() == reflect.TypeOf(time.Time{}) || (field.Type().Kind() == reflect.Ptr && field.Type().Elem() == reflect.TypeOf(time.Time{})) {
			t, err := valueToTime(value)
			if err != nil {
				return fmt.Errorf("field %s: %w", fieldName, err)
			}
			if field.Type().Kind() == reflect.Ptr {
				field.Set(reflect.ValueOf(&t))
			} else {
				field.Set(reflect.ValueOf(t))
			}
			continue
		}

		if err := setFieldValue(field, value); err != nil {
			return fmt.Errorf("field %s: %w", fieldName, err)
		}
	}

	return nil
}

// Unmarshal populates a struct from a byte slice.
// The struct fields can be tagged with `json:"fieldName"` to specify the corresponding key in the map.
// If a field is not tagged, it will match the map key with the same name as the field.
// Nested structs and maps/lists are supported.
func Unmarshal(data []byte, target any) error {
	m, err := DecodeMap(data)
	if err != nil {
		return err
	}

	return UnmarshalFromMap(m, target)
}

func setFieldValue(field reflect.Value, v *Value) error {
	if !field.CanSet() {
		return fmt.Errorf("cannot set field")
	}

	if field.Kind() == reflect.Ptr {
		ptr := reflect.New(field.Type().Elem())
		if err := setFieldValue(ptr.Elem(), v); err != nil {
			return err
		}
		field.Set(ptr)
		return nil
	}

	switch v.Type {
	case TypeBoolean:
		if field.Kind() == reflect.Bool {
			field.SetBool(v.AsBoolean())
			return nil
		}
	case TypeByte:
		if field.Kind() == reflect.Uint8 {
			field.SetUint(uint64(v.AsByte()))
			return nil
		}
	case TypeUint16:
		if field.Kind() == reflect.Uint16 {
			field.SetUint(uint64(v.AsUint16()))
			return nil
		}
	case TypeUint32:
		if field.Kind() == reflect.Uint32 || field.Kind() == reflect.Uint {
			field.SetUint(uint64(v.AsUint32()))
			return nil
		}
	case TypeUint64:
		if field.Kind() == reflect.Uint64 {
			field.SetUint(v.AsUint64())
			return nil
		}
	case TypeInt16:
		if field.Kind() == reflect.Int16 {
			field.SetInt(int64(v.AsInt16()))
			return nil
		}
	case TypeInt32:
		if field.Kind() == reflect.Int32 || field.Kind() == reflect.Int {
			field.SetInt(int64(v.AsInt32()))
			return nil
		}
	case TypeInt64:
		if field.Kind() == reflect.Int64 {
			field.SetInt(v.AsInt64())
			return nil
		}
	case TypeFloat32:
		if field.Kind() == reflect.Float32 {
			field.SetFloat(float64(v.AsFloat32()))
			return nil
		}
	case TypeFloat64:
		if field.Kind() == reflect.Float64 {
			field.SetFloat(v.AsFloat64())
			return nil
		}
	case TypeString:
		if field.Kind() == reflect.String {
			field.SetString(v.AsString())
			return nil
		}
	case TypeList:
		list := v.AsList()
		slice := reflect.MakeSlice(field.Type(), len(list), len(list))
		for i, item := range list {
			if err := setFieldValue(slice.Index(i), item); err != nil {
				return fmt.Errorf("index %d: %w", i, err)
			}
		}
		field.Set(slice)
		return nil
	case TypeMap:
		m := v.AsMap()
		switch field.Kind() {
		case reflect.Map:
			if field.IsNil() {
				field.Set(reflect.MakeMap(field.Type()))
			}
			for k, val := range m {
				mapVal := reflect.New(field.Type().Elem()).Elem()
				if err := setFieldValue(mapVal, val); err != nil {
					return fmt.Errorf("map key %s: %w", k, err)
				}
				field.SetMapIndex(reflect.ValueOf(k), mapVal)
			}
			return nil
		case reflect.Struct:
			// Nested struct
			return UnmarshalFromMap(m, field.Addr().Interface())
		default:
			return fmt.Errorf("unsupported field type for map: %s", field.Kind())
		}
	}

	return fmt.Errorf("unsupported conversion: ValueType=%d FieldType=%s", v.Type, field.Kind())
}

// time helpers

func timeToValue(t time.Time) map[string]*Value {
	return map[string]*Value{
		"millis": {Type: TypeInt64, data: t.UnixMilli()},
	}
}

func valueToTime(v *Value) (time.Time, error) {
	if v.Type != TypeMap {
		return time.Time{}, fmt.Errorf("expected map for time value, got %d", v.Type)
	}

	m := v.AsMap()
	millisVal, ok := m["millis"]
	if !ok || millisVal.Type != TypeInt64 {
		return time.Time{}, fmt.Errorf("time value must have 'millis' field of type int64")
	}

	millis := millisVal.AsInt64()
	return time.UnixMilli(millis), nil
}
