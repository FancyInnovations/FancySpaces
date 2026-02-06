package codex

import (
	"errors"
	"testing"
)

func TestEncodeBool(t *testing.T) {
	tests := []struct {
		name     string
		val      bool
		wantByte byte
	}{
		{"true", true, 1},
		{"false", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := EncodeBool(tt.val)

			if len(out) != 2 {
				t.Fatalf("expected length 2, got %d", len(out))
			}

			if out[0] != byte(TypeByte) {
				t.Fatalf("expected type %d, got %d", TypeByte, out[0])
			}

			if out[1] != tt.wantByte {
				t.Fatalf("expected byte value %d, got %d", tt.wantByte, out[1])
			}
		})
	}
}

func TestEncodeBoolInto(t *testing.T) {
	tests := []struct {
		name          string
		val           bool
		target        []byte
		expectRealloc bool
	}{
		{
			name:          "nil target",
			val:           true,
			target:        nil,
			expectRealloc: true,
		},
		{
			name:          "insufficient capacity",
			val:           false,
			target:        make([]byte, 0, 1),
			expectRealloc: true,
		},
		{
			name:          "sufficient capacity",
			val:           true,
			target:        make([]byte, 0, 8),
			expectRealloc: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origCap := cap(tt.target)
			out := EncodeBoolInto(tt.val, tt.target)

			if len(out) != 2 {
				t.Fatalf("expected length 2, got %d", len(out))
			}

			if out[0] != byte(TypeByte) {
				t.Fatalf("expected type %d, got %d", TypeByte, out[0])
			}

			expectedByte := byte(0)
			if tt.val {
				expectedByte = 1
			}

			if out[1] != expectedByte {
				t.Fatalf("expected byte %d, got %d", expectedByte, out[1])
			}

			if !tt.expectRealloc && cap(out) != origCap {
				t.Fatalf("unexpected reallocation: cap changed from %d to %d", origCap, cap(out))
			}
		})
	}
}

func TestDecodeBool(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    bool
		wantErr error
	}{
		{
			name: "true",
			data: EncodeBool(true),
			want: true,
		},
		{
			name: "false",
			data: EncodeBool(false),
			want: false,
		},
		{
			name: "non-1 byte is false",
			data: []byte{
				byte(TypeByte),
				0xFF,
			},
			want: false,
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
			got, err := DecodeBool(tt.data)

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
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
