package codex

import (
	"encoding/binary"
	"errors"
	"math"
	"testing"
)

func TestEncodeFloat32(t *testing.T) {
	tests := []struct {
		name string
		val  float32
	}{
		{"zero", 0},
		{"negative zero", float32(math.Copysign(0, -1))},
		{"one", 1},
		{"negative", -123.456},
		{"max", math.MaxFloat32},
		{"min", math.SmallestNonzeroFloat32},
		{"positive infinity", float32(math.Inf(1))},
		{"negative infinity", float32(math.Inf(-1))},
		{"nan", float32(math.NaN())},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := EncodeFloat32(tt.val)

			if len(out) != 5 {
				t.Fatalf("expected length 5, got %d", len(out))
			}

			if out[0] != byte(TypeFloat32) {
				t.Fatalf("expected type %d, got %d", TypeFloat32, out[0])
			}

			wantBits := math.Float32bits(tt.val)
			gotBits := binary.BigEndian.Uint32(out[1:5])

			if wantBits != gotBits {
				t.Fatalf("bit mismatch: expected %08x, got %08x", wantBits, gotBits)
			}
		})
	}
}

func TestEncodeFloat32Into(t *testing.T) {
	tests := []struct {
		name          string
		val           float32
		target        []byte
		expectRealloc bool
	}{
		{
			name:          "nil target",
			val:           1.23,
			target:        nil,
			expectRealloc: true,
		},
		{
			name:          "insufficient capacity",
			val:           -4.56,
			target:        make([]byte, 0, 3),
			expectRealloc: true,
		},
		{
			name:          "sufficient capacity",
			val:           7.89,
			target:        make([]byte, 0, 16),
			expectRealloc: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origCap := cap(tt.target)
			out := EncodeFloat32Into(tt.val, tt.target)

			if len(out) != 5 {
				t.Fatalf("expected length 5, got %d", len(out))
			}

			if out[0] != byte(TypeFloat32) {
				t.Fatalf("expected type %d, got %d", TypeFloat32, out[0])
			}

			wantBits := math.Float32bits(tt.val)
			gotBits := binary.BigEndian.Uint32(out[1:5])

			if wantBits != gotBits {
				t.Fatalf("bit mismatch: expected %08x, got %08x", wantBits, gotBits)
			}

			if !tt.expectRealloc && cap(out) != origCap {
				t.Fatalf("unexpected reallocation: cap changed from %d to %d", origCap, cap(out))
			}
		})
	}
}

func TestDecodeFloat32(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    float32
		wantErr error
	}{
		{
			name: "valid normal value",
			data: EncodeFloat32(3.14),
			want: 3.14,
		},
		{
			name: "valid negative zero",
			data: EncodeFloat32(float32(math.Copysign(0, -1))),
			want: float32(math.Copysign(0, -1)),
		},
		{
			name: "valid nan",
			data: EncodeFloat32(float32(math.NaN())),
			want: float32(math.NaN()),
		},
		{
			name: "trailing bytes ignored",
			data: append(EncodeFloat32(1.5), 0xFF, 0xEE),
			want: 1.5,
		},
		{
			name:    "payload too short",
			data:    []byte{byte(TypeFloat32)},
			wantErr: ErrPayloadTooShort,
		},
		{
			name: "invalid type",
			data: []byte{
				0xFF,
				0x00, 0x00, 0x00, 0x00,
			},
			wantErr: ErrInvalidType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeFloat32(tt.data)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if math.IsNaN(float64(tt.want)) {
				if !math.IsNaN(float64(got)) {
					t.Fatalf("expected NaN, got %v", got)
				}
				return
			}

			if got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
