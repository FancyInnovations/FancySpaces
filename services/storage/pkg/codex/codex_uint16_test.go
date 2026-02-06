package codex

import (
	"bytes"
	"errors"
	"testing"
)

func TestEncodeUint16Into(t *testing.T) {
	tests := []struct {
		name          string
		val           uint16
		initialTarget []byte
		expectReuse   bool
		expected      []byte
	}{
		{
			name:          "nil target allocates new slice",
			val:           42,
			initialTarget: nil,
			expectReuse:   false,
			expected:      []byte{byte(TypeUint16), 0x00, 0x2a},
		},
		{
			name:          "insufficient capacity allocates new slice",
			val:           123,
			initialTarget: make([]byte, 0, 2), // cap < 3
			expectReuse:   false,
			expected:      []byte{byte(TypeUint16), 0x00, 0x7b},
		},
		{
			name:          "sufficient capacity reuses slice",
			val:           256,
			initialTarget: make([]byte, 5), // cap >= 3
			expectReuse:   true,
			expected:      []byte{byte(TypeUint16), 0x01, 0x00},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := tt.initialTarget
			out := EncodeUint16Into(tt.val, tt.initialTarget)

			if !bytes.Equal(out, tt.expected) {
				t.Fatalf("unexpected output: got %v, want %v", out, tt.expected)
			}

			if tt.expectReuse && orig != nil {
				if &out[0] != &orig[0] {
					t.Fatalf("expected slice reuse but got new allocation")
				}
			}

			if !tt.expectReuse && orig != nil && len(orig) > 0 {
				if &out[0] == &orig[0] {
					t.Fatalf("expected new allocation but slice was reused")
				}
			}

			if len(out) != 3 {
				t.Fatalf("unexpected length: got %d, want 3", len(out))
			}
		})
	}
}

func TestEncodeUint16(t *testing.T) {
	val := uint16(32768)
	out := EncodeUint16(val)

	expected := []byte{byte(TypeUint16), 0x80, 0x00}
	if !bytes.Equal(out, expected) {
		t.Fatalf("unexpected output: got %v, want %v", out, expected)
	}
}

func TestDecodeUint16(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    uint16
		wantErr error
	}{
		{
			name:    "payload too short (empty)",
			data:    []byte{},
			wantErr: ErrPayloadTooShort,
		},
		{
			name:    "payload too short (partial)",
			data:    []byte{byte(TypeUint16), 0x00},
			wantErr: ErrPayloadTooShort,
		},
		{
			name:    "invalid type",
			data:    []byte{0xff, 0x00, 0x01},
			wantErr: ErrInvalidType,
		},
		{
			name: "valid value",
			data: []byte{
				byte(TypeUint16),
				0x12, 0x34,
			},
			want: 0x1234,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeUint16(tt.data)

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

			if got != tt.want {
				t.Fatalf("unexpected value: got %d, want %d", got, tt.want)
			}
		})
	}
}
