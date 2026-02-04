package kv

type Value struct {
	Type ValueType
	data any
}

type ValueType uint8

const (
	TypeString ValueType = iota
	TypeNumber
	TypeList
	TypeSet
	TypeMap
	TypeBinary
)

func (v Value) AsString() string {
	if v.Type != TypeString {
		return ""
	}
	return v.data.(string)
}

func (v Value) AsNumber() float64 {
	if v.Type != TypeNumber {
		return 0
	}
	return v.data.(float64)
}

func (v Value) AsList() []Value {
	if v.Type != TypeList {
		return nil
	}
	return v.data.([]Value)
}

func (v Value) AsSet() map[string]struct{} {
	if v.Type != TypeSet {
		return nil
	}
	return v.data.(map[string]struct{})
}

func (v Value) AsMap() map[string]Value {
	if v.Type != TypeMap {
		return nil
	}
	return v.data.(map[string]Value)
}

func (v Value) AsBinary() []byte {
	if v.Type != TypeBinary {
		return nil
	}
	return v.data.([]byte)
}

func NewStringValue(s string) Value {
	return Value{
		Type: TypeString,
		data: s,
	}
}

func NewNumberValue(n float64) Value {
	return Value{
		Type: TypeNumber,
		data: n,
	}
}

func NewListValue(l []Value) Value {
	return Value{
		Type: TypeList,
		data: l,
	}
}

func NewSetValue(s map[string]struct{}) Value {
	return Value{
		Type: TypeSet,
		data: s,
	}
}

func NewMapValue(m map[string]Value) Value {
	return Value{
		Type: TypeMap,
		data: m,
	}
}

func NewBinaryValue(b []byte) Value {
	return Value{
		Type: TypeBinary,
		data: b,
	}
}
