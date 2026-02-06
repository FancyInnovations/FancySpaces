package codex

import (
	"bytes"
	"errors"
	"testing"
)

func TestEncodeInt64Into(t *testing.T) {
	tests := []struct {
		name          string
		val           int64
		initialTarget []byte
		expectReuse   bool
		expected      []byte
	}{
		{
			name:          "nil target allocates new slice",
			val:           42,
			initialTarget: nil,
			expectReuse:   false,
			expected: []byte{
				byte(TypeInt64),
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x2a,
			},
		},
		{
			name:          "insufficient capacity allocates new slice",
			val:           -1,
			initialTarget: make([]byte, 0, 8), // cap < 9
			expectReuse:   false,
			expected: []byte{
				byte(TypeInt64),
				0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff,
			},
		},
		{
			name:          "sufficient capacity reuses slice",
			val:           0x0102030405060708,
			initialTarget: make([]byte, 16), // cap >= 9
			expectReuse:   true,
			expected: []byte{
				byte(TypeInt64),
				0x01, 0x02, 0x03, 0x04,
				0x05, 0x06, 0x07, 0x08,
			},
		},
		{
			name:          "minimum int64 value",
			val:           -9223372036854775808,
			initialTarget: make([]byte, 9),
			expectReuse:   true,
			expected: []byte{
				byte(TypeInt64),
				0x80, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			name:          "maximum int64 value",
			val:           9223372036854775807,
			initialTarget: make([]byte, 9),
			expectReuse:   true,
			expected: []byte{
				byte(TypeInt64),
				0x7f, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := tt.initialTarget
			out := EncodeInt64Into(tt.val, tt.initialTarget)

			if !bytes.Equal(out, tt.expected) {
				t.Fatalf("unexpected output: got %v, want %v", out, tt.expected)
			}

			if len(out) != 9 {
				t.Fatalf("unexpected length: got %d, want 9", len(out))
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

func TestEncodeInt64(t *testing.T) {
	val := int64(-9223372036854775808)
	out := EncodeInt64(val)

	expected := []byte{
		byte(TypeInt64),
		0x80, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	if !bytes.Equal(out, expected) {
		t.Fatalf("unexpected output: got %v, want %v", out, expected)
	}
}

func TestDecodeInt64(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    int64
		wantErr error
	}{
		{
			name:    "payload too short (empty)",
			data:    []byte{},
			wantErr: ErrPayloadTooShort,
		},
		{
			name:    "payload too short (partial)",
			data:    []byte{byte(TypeInt64), 0x00, 0x00, 0x00},
			wantErr: ErrPayloadTooShort,
		},
		{
			name:    "invalid type",
			data:    []byte{0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
			wantErr: ErrInvalidType,
		},
		{
			name: "valid positive value",
			data: []byte{
				byte(TypeInt64),
				0x12, 0x34, 0x56, 0x78,
				0x9a, 0xbc, 0xde, 0xf0,
			},
			want: 0x123456789abcdef0,
		},
		{
			name: "valid negative value",
			data: []byte{
				byte(TypeInt64),
				0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xfe,
			},
			want: -2,
		},
		{
			name: "minimum int64 value",
			data: []byte{
				byte(TypeInt64),
				0x80, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
			want: -9223372036854775808,
		},
		{
			name: "maximum int64 value",
			data: []byte{
				byte(TypeInt64),
				0x7f, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff,
			},
			want: 9223372036854775807,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeInt64(tt.data)

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
