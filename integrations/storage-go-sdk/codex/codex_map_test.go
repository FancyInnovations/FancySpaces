package codex

import (
	"encoding/binary"
	"errors"
	"testing"
)

func TestEncodeMapInto(t *testing.T) {
	v1 := &Value{Type: TypeInt32, data: int32(1)}
	v2 := &Value{Type: TypeInt32, data: int32(2)}
	empty := &Value{Type: TypeEmpty}

	tests := []struct {
		name          string
		vals          map[string]*Value
		dst           []byte
		expectNil     bool
		expectReuse   bool
		expectedCount uint32
	}{
		{
			name:          "empty map",
			vals:          map[string]*Value{},
			dst:           nil,
			expectedCount: 0,
		},
		{
			name: "single entry",
			vals: map[string]*Value{
				"a": v1,
			},
			dst:           nil,
			expectedCount: 1,
		},
		{
			name: "multiple entries same type",
			vals: map[string]*Value{
				"a": v1,
				"b": v2,
			},
			dst:           make([]byte, 128),
			expectReuse:   true,
			expectedCount: 2,
		},
		{
			name: "mixed value types is valid",
			vals: map[string]*Value{
				"a": v1,
				"b": {Type: TypeInt64, data: int64(1)},
			},
			expectNil:     false,
			expectedCount: 2,
		},
		{
			name: "empty value allowed",
			vals: map[string]*Value{
				"a": v1,
				"b": empty,
			},
			dst:           nil,
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := tt.dst
			out := EncodeMapInto(tt.vals, tt.dst)

			if tt.expectNil {
				if out != nil {
					t.Fatalf("expected nil output, got %v", out)
				}
				return
			}

			if out == nil {
				t.Fatalf("unexpected nil output")
			}

			if out[0] != byte(TypeMap) {
				t.Fatalf("unexpected type byte: got %d", out[0])
			}

			count := binary.BigEndian.Uint32(out[1:5])
			if count != tt.expectedCount {
				t.Fatalf("unexpected count: got %d want %d", count, tt.expectedCount)
			}

			if tt.expectReuse && orig != nil {
				if &out[0] != &orig[0] {
					t.Fatalf("expected slice reuse but got new allocation")
				}
			}
		})
	}
}

func TestEncodeMap(t *testing.T) {
	vals := map[string]*Value{
		"foo": {Type: TypeInt32, data: int32(42)},
	}

	out := EncodeMap(vals)
	if out == nil {
		t.Fatalf("unexpected nil output")
	}

	if out[0] != byte(TypeMap) {
		t.Fatalf("unexpected type byte: %d", out[0])
	}
}

func TestDecodeMap(t *testing.T) {
	validVals := map[string]*Value{
		"a": {Type: TypeInt32, data: int32(1)},
		"b": {Type: TypeInt32, data: int32(2)},
	}

	validEncoded := EncodeMap(validVals)

	tests := []struct {
		name    string
		data    []byte
		wantLen int
		wantErr error
	}{
		{
			name:    "payload too short (header)",
			data:    []byte{byte(TypeMap)},
			wantErr: ErrPayloadTooShort,
		},
		{
			name:    "invalid map type",
			data:    []byte{0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			wantErr: ErrInvalidType,
		},
		{
			name: "payload too short (items)",
			data: func() []byte {
				b := make([]byte, 8)
				b[0] = byte(TypeMap)
				b[1] = byte(TypeInt32)
				binary.BigEndian.PutUint16(b[2:4], 1)
				binary.BigEndian.PutUint32(b[4:8], 10)
				return b
			}(),
			wantErr: ErrPayloadTooShort,
		},
		{
			name:    "valid map",
			data:    validEncoded,
			wantLen: 2,
		},
		{
			name: "invalid value type",
			data: func() []byte {
				b := EncodeMap(validVals)
				b[0] = byte(TypeInt64) // force mismatch
				return b
			}(),
			wantErr: ErrInvalidType,
		},
		{
			name: "decode string error bubbles up",
			data: func() []byte {
				b := EncodeMap(validVals)
				return b[:len(b)-1] // truncate key/value payload
			}(),
			wantErr: ErrPayloadTooShort,
		},
		{
			name: "decode value error bubbles up",
			data: func() []byte {
				b := EncodeMap(validVals)
				return b[:len(b)-1]
			}(),
			wantErr: ErrPayloadTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vals, err := DecodeMap(tt.data)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("unexpected error: got %v want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(vals) != tt.wantLen {
				t.Fatalf("unexpected map size: got %d want %d", len(vals), tt.wantLen)
			}

			// spot-check decoded values
			for k, v := range vals {
				if v == nil {
					t.Fatalf("nil value for key %q", k)
				}
			}
		})
	}
}
