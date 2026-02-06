package codex

import (
	"encoding/binary"
	"errors"
	"testing"
)

func TestEncodeListInto(t *testing.T) {
	v1 := &Value{Type: TypeInt32, data: int32(1)}
	v2 := &Value{Type: TypeInt32, data: int32(2)}
	empty := &Value{Type: TypeEmpty}

	tests := []struct {
		name          string
		vals          []*Value
		dst           []byte
		expectNil     bool
		expectReuse   bool
		expectedType  byte
		expectedCount uint16
	}{
		{
			name:          "empty list",
			vals:          nil,
			dst:           nil,
			expectReuse:   false,
			expectedType:  byte(TypeEmpty),
			expectedCount: 0,
		},
		{
			name:          "single item list",
			vals:          []*Value{v1},
			dst:           nil,
			expectReuse:   false,
			expectedType:  byte(TypeInt32),
			expectedCount: 1,
		},
		{
			name:          "multiple items same type",
			vals:          []*Value{v1, v2},
			dst:           make([]byte, 64),
			expectReuse:   true,
			expectedType:  byte(TypeInt32),
			expectedCount: 2,
		},
		{
			name:      "mixed types is invalid",
			vals:      []*Value{v1, &Value{Type: TypeInt64, data: int64(1)}},
			expectNil: true,
		},
		{
			name:          "empty values allowed",
			vals:          []*Value{v1, empty},
			dst:           nil,
			expectReuse:   false,
			expectedType:  byte(TypeInt32),
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := tt.dst
			out := EncodeListInto(tt.vals, tt.dst)

			if tt.expectNil {
				if out != nil {
					t.Fatalf("expected nil output, got %v", out)
				}
				return
			}

			if out == nil {
				t.Fatalf("unexpected nil output")
			}

			if out[0] != byte(TypeList) {
				t.Fatalf("unexpected type byte: got %d", out[0])
			}

			if out[1] != tt.expectedType {
				t.Fatalf("unexpected item type: got %d, want %d", out[1], tt.expectedType)
			}

			count := binary.BigEndian.Uint16(out[2:4])
			if count != tt.expectedCount {
				t.Fatalf("unexpected count: got %d, want %d", count, tt.expectedCount)
			}

			payloadLen := binary.BigEndian.Uint32(out[4:8])
			if int(payloadLen) != len(out)-8 {
				t.Fatalf("payload length mismatch: header=%d actual=%d", payloadLen, len(out)-8)
			}

			if tt.expectReuse && orig != nil {
				if &out[0] != &orig[0] {
					t.Fatalf("expected slice reuse but got new allocation")
				}
			}
		})
	}
}

func TestEncodeList(t *testing.T) {
	vals := []*Value{
		{Type: TypeInt32, data: int32(123)},
	}

	out := EncodeList(vals)
	if out == nil {
		t.Fatalf("unexpected nil output")
	}

	if out[0] != byte(TypeList) {
		t.Fatalf("unexpected type byte: %d", out[0])
	}
}

func TestDecodeList(t *testing.T) {
	validVals := []*Value{
		{Type: TypeInt32, data: int32(1)},
		{Type: TypeInt32, data: int32(2)},
	}

	validEncoded := EncodeList(validVals)

	tests := []struct {
		name    string
		data    []byte
		wantLen int
		wantErr error
	}{
		{
			name:    "payload too short (header)",
			data:    []byte{byte(TypeList)},
			wantErr: ErrPayloadTooShort,
		},
		{
			name:    "invalid list type",
			data:    []byte{0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			wantErr: ErrInvalidType,
		},
		{
			name: "payload too short (items)",
			data: func() []byte {
				b := make([]byte, 8)
				b[0] = byte(TypeList)
				b[1] = byte(TypeInt32)
				binary.BigEndian.PutUint16(b[2:4], 1)
				binary.BigEndian.PutUint32(b[4:8], 10)
				return b
			}(),
			wantErr: ErrPayloadTooShort,
		},
		{
			name:    "valid list",
			data:    validEncoded,
			wantLen: 2,
		},
		{
			name: "invalid item type",
			data: func() []byte {
				b := EncodeList(validVals)
				b[1] = byte(TypeInt64) // force mismatch
				return b
			}(),
			wantErr: ErrInvalidType,
		},
		{
			name: "decode value error bubbles up",
			data: func() []byte {
				b := EncodeList(validVals)
				return b[:len(b)-1] // truncate item payload
			}(),
			wantErr: ErrPayloadTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := DecodeList(tt.data)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("unexpected error: got %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(items) != tt.wantLen {
				t.Fatalf("unexpected item count: got %d, want %d", len(items), tt.wantLen)
			}
		})
	}
}
