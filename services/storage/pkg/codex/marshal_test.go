package codex

import (
	"reflect"
	"testing"
)

type Nested struct {
	Flag bool `json:"flag"`
	Name string
}

type Sample struct {
	Bool      bool              `json:"bool"`
	Int       int               `json:"int"`
	String    string            `json:"string"`
	Ptr       *string           `json:"ptr"`
	List      []int             `json:"list"`
	Map       map[string]string `json:"map"`
	Nested    Nested            `json:"nested"`
	PtrNested *Nested           `json:"ptr_nested"`
}

func TestMarshalToMap(t *testing.T) {
	str := "hello"

	tests := []struct {
		name    string
		input   any
		wantErr bool
	}{
		{
			name: "basic struct",
			input: Sample{
				Bool:   true,
				Int:    42,
				String: "test",
				Ptr:    &str,
				List:   []int{1, 2, 3},
				Map:    map[string]string{"a": "b"},
				Nested: Nested{
					Flag: true,
					Name: "nested",
				},
				PtrNested: &Nested{
					Flag: false,
					Name: "ptr",
				},
			},
		},
		{
			name:    "non-struct input",
			input:   123,
			wantErr: true,
		},
		{
			name:    "nil pointer",
			input:   (*Sample)(nil),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := MarshalToMap(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if out == nil {
				t.Fatalf("expected non-nil map")
			}

			// spot-check a few important fields
			if v, ok := out["bool"]; !ok || v.AsBoolean() != true {
				t.Errorf("bool field not marshaled correctly")
			}

			if v, ok := out["int"]; !ok || v.AsInt() != 42 {
				t.Errorf("int field not marshaled correctly")
			}

			if v, ok := out["string"]; !ok || v.AsString() != "test" {
				t.Errorf("string field not marshaled correctly")
			}

			if v, ok := out["nested"]; !ok || v.Type != TypeMap {
				t.Errorf("nested struct not marshaled as map")
			}
		})
	}
}

func TestUnmarshalFromMap(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]*Value
		target  any
		want    any
		wantErr bool
	}{
		{
			name: "basic unmarshal",
			input: map[string]*Value{
				"bool":   mustValue(true),
				"int":    mustValue(int32(10)),
				"string": mustValue("hello"),
				"list": &Value{
					Type: TypeList,
					data: []*Value{
						mustValue(int32(1)),
						mustValue(int32(2)),
					},
				},
				"map": &Value{
					Type: TypeMap,
					data: map[string]*Value{
						"a": mustValue("b"),
					},
				},
				"nested": &Value{
					Type: TypeMap,
					data: map[string]*Value{
						"flag": mustValue(true),
						"Name": mustValue("nested"),
					},
				},
			},
			target: &Sample{},
			want: &Sample{
				Bool:   true,
				Int:    10,
				String: "hello",
				List:   []int{1, 2},
				Map:    map[string]string{"a": "b"},
				Nested: Nested{
					Flag: true,
					Name: "nested",
				},
			},
		},
		{
			name:    "non-pointer target",
			input:   map[string]*Value{},
			target:  Sample{},
			wantErr: true,
		},
		{
			name:    "pointer to non-struct",
			input:   map[string]*Value{},
			target:  new(int),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalFromMap(tt.input, tt.target)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(tt.target, tt.want) {
				t.Errorf("result mismatch\nwant: %#v\ngot:  %#v", tt.want, tt.target)
			}
		})
	}
}

func TestMarshalToMapUnmarshalFromMap_RoundTrip(t *testing.T) {
	orig := Sample{
		Bool:   true,
		Int:    99,
		String: "roundtrip",
		List:   []int{7, 8, 9},
		Map:    map[string]string{"k": "v"},
		Nested: Nested{
			Flag: true,
			Name: "inner",
		},
	}

	m, err := MarshalToMap(orig)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var out Sample
	if err := UnmarshalFromMap(m, &out); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if !reflect.DeepEqual(orig, out) {
		t.Errorf("round-trip mismatch\nwant: %#v\ngot:  %#v", orig, out)
	}
}

func TestUnmarshalFromMap_MissingFields(t *testing.T) {
	input := map[string]*Value{
		"bool": mustValue(true),
		"int":  mustValue(int32(5)),
	}

	var target Sample
	if err := UnmarshalFromMap(input, &target); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if target.Bool != true {
		t.Errorf("expected Bool=true, got %v", target.Bool)
	}
	if target.Int != 5 {
		t.Errorf("expected Int=5, got %d", target.Int)
	}
	if target.String != "" {
		t.Errorf("expected String=\"\", got %q", target.String)
	}
	if target.List != nil {
		t.Errorf("expected List=nil, got %v", target.List)
	}
	if target.Map != nil {
		t.Errorf("expected Map=nil, got %v", target.Map)
	}
	if (target.Nested != Nested{}) {
		t.Errorf("expected Nested zero value, got %v", target.Nested)
	}
	if target.PtrNested != nil {
		t.Errorf("expected PtrNested=nil, got %v", target.PtrNested)
	}
}

func TestMarshalUnmarshal(t *testing.T) {
	orig := Sample{
		Bool:   false,
		Int:    -1,
		String: "marshal test",
		Ptr:    nil,
		List:   []int{10, 20},
		Map:    map[string]string{"x": "y"},
		Nested: Nested{
			Flag: false,
			Name: "nested marshal",
		},
	}

	data, err := Marshal(orig)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var out Sample
	if err := Unmarshal(data, &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if !reflect.DeepEqual(orig, out) {
		t.Errorf("Marshal/Unmarshal mismatch\nwant: %#v\ngot:  %#v", orig, out)
	}
}

// helper
func mustValue(v any) *Value {
	val, err := NewValue(v)
	if err != nil {
		panic(err)
	}
	return val
}
