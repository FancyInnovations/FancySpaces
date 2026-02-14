package codex

import (
	"bytes"
	"errors"
	"testing"
)

func TestEncodeInt32Into(t *testing.T) {
	tests := []struct {
		name          string
		val           int32
		initialTarget []byte
		expectReuse   bool
		expected      []byte
	}{
		{
			name:          "nil target allocates new slice",
			val:           42,
			initialTarget: nil,
			expectReuse:   false,
			expected:      []byte{byte(TypeInt32), 0x00, 0x00, 0x00, 0x2a},
		},
		{
			name:          "insufficient capacity allocates new slice",
			val:           -1,
			initialTarget: make([]byte, 0, 4), // cap < 5
			expectReuse:   false,
			expected:      []byte{byte(TypeInt32), 0xff, 0xff, 0xff, 0xff},
		},
		{
			name:          "sufficient capacity reuses slice",
			val:           0x01020304,
			initialTarget: make([]byte, 10), // cap >= 5
			expectReuse:   true,
			expected:      []byte{byte(TypeInt32), 0x01, 0x02, 0x03, 0x04},
		},
		{
			name:          "minimum int32 value",
			val:           -2147483648,
			initialTarget: make([]byte, 5),
			expectReuse:   true,
			expected:      []byte{byte(TypeInt32), 0x80, 0x00, 0x00, 0x00},
		},
		{
			name:          "maximum int32 value",
			val:           2147483647,
			initialTarget: make([]byte, 5),
			expectReuse:   true,
			expected:      []byte{byte(TypeInt32), 0x7f, 0xff, 0xff, 0xff},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := tt.initialTarget
			out := EncodeInt32Into(tt.val, tt.initialTarget)

			if !bytes.Equal(out, tt.expected) {
				t.Fatalf("unexpected output: got %v, want %v", out, tt.expected)
			}

			if len(out) != 5 {
				t.Fatalf("unexpected length: got %d, want 5", len(out))
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
		})
	}
}

func TestEncodeInt32(t *testing.T) {
	val := int32(-2147483648)
	out := EncodeInt32(val)

	expected := []byte{byte(TypeInt32), 0x80, 0x00, 0x00, 0x00}
	if !bytes.Equal(out, expected) {
		t.Fatalf("unexpected output: got %v, want %v", out, expected)
	}
}

func TestDecodeInt32(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    int32
		wantErr error
	}{
		{
			name:    "payload too short (empty)",
			data:    []byte{},
			wantErr: ErrPayloadTooShort,
		},
		{
			name:    "payload too short (partial)",
			data:    []byte{byte(TypeInt32), 0x00, 0x00},
			wantErr: ErrPayloadTooShort,
		},
		{
			name:    "invalid type",
			data:    []byte{0xff, 0x00, 0x00, 0x00, 0x01},
			wantErr: ErrInvalidType,
		},
		{
			name: "valid positive value",
			data: []byte{
				byte(TypeInt32),
				0x12, 0x34, 0x56, 0x78,
			},
			want: 0x12345678,
		},
		{
			name: "valid negative value",
			data: []byte{
				byte(TypeInt32),
				0xff, 0xff, 0xff, 0xfe,
			},
			want: -2,
		},
		{
			name: "minimum int32 value",
			data: []byte{
				byte(TypeInt32),
				0x80, 0x00, 0x00, 0x00,
			},
			want: -2147483648,
		},
		{
			name: "maximum int32 value",
			data: []byte{
				byte(TypeInt32),
				0x7f, 0xff, 0xff, 0xff,
			},
			want: 2147483647,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeInt32(tt.data)

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
