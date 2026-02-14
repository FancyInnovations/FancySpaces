package codex

import (
	"encoding/binary"
	"errors"
	"math"
	"testing"
)

func TestEncodeFloat64(t *testing.T) {
	tests := []struct {
		name string
		val  float64
	}{
		{
			name: "zero",
			val:  0,
		},
		{
			name: "negative zero",
			val:  math.Copysign(0, -1),
		},
		{
			name: "one",
			val:  1,
		},
		{
			name: "negative",
			val:  -123.456,
		},
		{
			name: "max",
			val:  math.MaxFloat64,
		},
		{
			name: "min",
			val:  math.SmallestNonzeroFloat64,
		},
		{
			name: "positive infinity",
			val:  math.Inf(1),
		},
		{
			name: "negative infinity",
			val:  math.Inf(-1),
		},
		{
			name: "nan",
			val:  math.NaN(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := EncodeFloat64(tt.val)

			if len(out) != 9 {
				t.Fatalf("expected length 9, got %d", len(out))
			}

			if out[0] != byte(TypeFloat64) {
				t.Fatalf("expected type %d, got %d", TypeFloat64, out[0])
			}

			wantBits := math.Float64bits(tt.val)
			gotBits := binary.BigEndian.Uint64(out[1:9])

			if wantBits != gotBits {
				t.Fatalf("bit mismatch: expected %016x, got %016x", wantBits, gotBits)
			}
		})
	}
}

func TestEncodeFloat64Into(t *testing.T) {
	tests := []struct {
		name          string
		val           float64
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
			target:        make([]byte, 0, 4),
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
			out := EncodeFloat64Into(tt.val, tt.target)

			if len(out) != 9 {
				t.Fatalf("expected length 9, got %d", len(out))
			}

			if out[0] != byte(TypeFloat64) {
				t.Fatalf("expected type %d, got %d", TypeFloat64, out[0])
			}

			wantBits := math.Float64bits(tt.val)
			gotBits := binary.BigEndian.Uint64(out[1:9])

			if wantBits != gotBits {
				t.Fatalf("bit mismatch: expected %016x, got %016x", wantBits, gotBits)
			}

			if !tt.expectRealloc && cap(out) != origCap {
				t.Fatalf("unexpected reallocation: cap changed from %d to %d", origCap, cap(out))
			}
		})
	}
}

func TestDecodeFloat64(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    float64
		wantErr error
	}{
		{
			name:    "valid normal value",
			data:    EncodeFloat64(3.1415),
			want:    3.1415,
			wantErr: nil,
		},
		{
			name:    "valid negative zero",
			data:    EncodeFloat64(math.Copysign(0, -1)),
			want:    math.Copysign(0, -1),
			wantErr: nil,
		},
		{
			name:    "valid nan",
			data:    EncodeFloat64(math.NaN()),
			want:    math.NaN(),
			wantErr: nil,
		},
		{
			name:    "trailing bytes ignored",
			data:    append(EncodeFloat64(1.5), 0xFF, 0xEE),
			want:    1.5,
			wantErr: nil,
		},
		{
			name:    "payload too short",
			data:    []byte{byte(TypeFloat64)},
			want:    0,
			wantErr: ErrPayloadTooShort,
		},
		{
			name:    "invalid type",
			data:    []byte{0xFF, 0, 0, 0, 0, 0, 0, 0, 0},
			want:    0,
			wantErr: ErrInvalidType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeFloat64(tt.data)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if math.IsNaN(tt.want) {
				if !math.IsNaN(got) {
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
