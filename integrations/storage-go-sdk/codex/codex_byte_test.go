package codex

import (
	"errors"
	"testing"
)

func TestEncodeByte(t *testing.T) {
	tests := []struct {
		name string
		val  byte
	}{
		{"zero", 0x00},
		{"one", 0x01},
		{"max", 0xFF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := EncodeByte(tt.val)

			if len(out) != 2 {
				t.Fatalf("expected length 2, got %d", len(out))
			}

			if out[0] != byte(TypeByte) {
				t.Fatalf("expected type %d, got %d", TypeByte, out[0])
			}

			if out[1] != tt.val {
				t.Fatalf("expected value %d, got %d", tt.val, out[1])
			}
		})
	}
}

func TestEncodeByteInto(t *testing.T) {
	tests := []struct {
		name          string
		val           byte
		target        []byte
		expectRealloc bool
	}{
		{
			name:          "nil target",
			val:           0x01,
			target:        nil,
			expectRealloc: true,
		},
		{
			name:          "insufficient capacity",
			val:           0x02,
			target:        make([]byte, 0, 1),
			expectRealloc: true,
		},
		{
			name:          "sufficient capacity",
			val:           0xFF,
			target:        make([]byte, 0, 8),
			expectRealloc: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origCap := cap(tt.target)
			out := EncodeByteInto(tt.val, tt.target)

			if len(out) != 2 {
				t.Fatalf("expected length 2, got %d", len(out))
			}

			if out[0] != byte(TypeByte) {
				t.Fatalf("expected type %d, got %d", TypeByte, out[0])
			}

			if out[1] != tt.val {
				t.Fatalf("expected value %d, got %d", tt.val, out[1])
			}

			if !tt.expectRealloc && cap(out) != origCap {
				t.Fatalf("unexpected reallocation: cap changed from %d to %d", origCap, cap(out))
			}
		})
	}
}

func TestDecodeByte(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    byte
		wantErr error
	}{
		{
			name: "valid zero",
			data: EncodeByte(0x00),
			want: 0x00,
		},
		{
			name: "valid non-zero",
			data: EncodeByte(0xAB),
			want: 0xAB,
		},
		{
			name: "trailing bytes ignored",
			data: []byte{
				byte(TypeByte),
				0x42,
				0xFF,
				0xEE,
			},
			want: 0x42,
		},
		{
			name:    "payload too short",
			data:    []byte{byte(TypeByte)},
			wantErr: ErrPayloadTooShort,
		},
		{
			name: "invalid type",
			data: []byte{
				0xFF,
				0x01,
			},
			wantErr: ErrInvalidType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeByte(tt.data)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, got)
			}
		})
	}
}
