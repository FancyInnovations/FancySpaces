package codex

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"
)

func TestEncodeBinary(t *testing.T) {
	tests := []struct {
		name string
		val  []byte
	}{
		{"empty", []byte{}},
		{"small", []byte{0x01, 0x02, 0x03}},
		{"large", bytes.Repeat([]byte{0xAB}, 100_000)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := EncodeBinary(tt.val)

			expectedLen := 1 + 4 + len(tt.val)
			if len(out) != expectedLen {
				t.Fatalf("expected length %d, got %d", expectedLen, len(out))
			}

			if out[0] != byte(TypeBinary) {
				t.Fatalf("wrong type byte: %d", out[0])
			}

			gotLen := binary.BigEndian.Uint32(out[1:5])
			if int(gotLen) != len(tt.val) {
				t.Fatalf("expected encoded length %d, got %d", len(tt.val), gotLen)
			}

			if !bytes.Equal(out[5:], tt.val) {
				t.Fatalf("payload mismatch")
			}
		})
	}
}

func TestEncodeBinaryInto(t *testing.T) {
	tests := []struct {
		name          string
		val           []byte
		target        []byte
		expectRealloc bool
	}{
		{
			name:          "nil target",
			val:           []byte{0x01},
			target:        nil,
			expectRealloc: true,
		},
		{
			name:          "insufficient capacity",
			val:           []byte{0x01, 0x02},
			target:        make([]byte, 0, 2),
			expectRealloc: true,
		},
		{
			name:          "sufficient capacity",
			val:           []byte{0x01, 0x02},
			target:        make([]byte, 0, 64),
			expectRealloc: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origCap := cap(tt.target)
			out := EncodeBinaryInto(tt.val, tt.target)

			if out[0] != byte(TypeBinary) {
				t.Fatalf("wrong type byte")
			}

			gotLen := binary.BigEndian.Uint32(out[1:5])
			if int(gotLen) != len(tt.val) {
				t.Fatalf("expected length %d, got %d", len(tt.val), gotLen)
			}

			if !bytes.Equal(out[5:], tt.val) {
				t.Fatalf("payload mismatch")
			}

			if !tt.expectRealloc && cap(out) != origCap {
				t.Fatalf("unexpected reallocation")
			}
		})
	}
}

func TestDecodeBinary(t *testing.T) {
	validPayload := []byte{0x10, 0x20, 0x30}
	validEncoded := EncodeBinary(validPayload)

	tests := []struct {
		name    string
		data    []byte
		want    []byte
		wantErr error
	}{
		{
			name:    "too short",
			data:    []byte{0x01, 0x02},
			wantErr: ErrPayloadTooShort,
		},
		{
			name: "invalid type",
			data: func() []byte {
				b := make([]byte, len(validEncoded))
				copy(b, validEncoded)
				b[0] = 0xFF
				return b
			}(),
			wantErr: ErrInvalidType,
		},
		{
			name: "length exceeds payload",
			data: func() []byte {
				b := make([]byte, 5)
				b[0] = byte(TypeBinary)
				binary.BigEndian.PutUint32(b[1:5], 10)
				return b
			}(),
			wantErr: ErrPayloadTooShort,
		},
		{
			name: "valid",
			data: validEncoded,
			want: validPayload,
		},
		{
			name: "trailing bytes ignored",
			data: append(validEncoded, 0xFF, 0xEE),
			want: validPayload,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeBinary(tt.data)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !bytes.Equal(got, tt.want) {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
