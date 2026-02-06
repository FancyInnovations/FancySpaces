package codex

import (
	"encoding/binary"
	"errors"
	"testing"
)

func TestEncodeStringInto(t *testing.T) {
	tests := []struct {
		name          string
		val           string
		initialTarget []byte
		expectReuse   bool
	}{
		{
			name: "empty string",
			val:  "",
		},
		{
			name: "ascii string",
			val:  "abc",
		},
		{
			name: "utf8 string",
			val:  "€",
		},
		{
			name:          "reuse target slice",
			val:           "hello",
			initialTarget: make([]byte, 32),
			expectReuse:   true,
		},
		{
			name:          "insufficient capacity allocates",
			val:           "hello",
			initialTarget: make([]byte, 0, 2),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := tt.initialTarget
			out := EncodeStringInto(tt.val, tt.initialTarget)

			if out[0] != byte(TypeString) {
				t.Fatalf("unexpected type byte: %d", out[0])
			}

			length := int(binary.BigEndian.Uint32(out[1:5]))
			if length != len(tt.val) {
				t.Fatalf("unexpected length: got %d want %d", length, len(tt.val))
			}

			if string(out[5:]) != tt.val {
				t.Fatalf("unexpected payload: got %q want %q", string(out[5:]), tt.val)
			}

			if tt.expectReuse && orig != nil {
				if &out[0] != &orig[0] {
					t.Fatalf("expected slice reuse")
				}
			}
		})
	}
}

func TestEncodeString(t *testing.T) {
	out := EncodeString("test")

	if out[0] != byte(TypeString) {
		t.Fatalf("unexpected type byte")
	}

	length := binary.BigEndian.Uint32(out[1:5])
	if length != 4 {
		t.Fatalf("unexpected length: %d", length)
	}

	if string(out[5:]) != "test" {
		t.Fatalf("unexpected payload")
	}
}

func TestDecodeString(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    string
		wantErr error
	}{
		{
			name:    "payload too short (empty)",
			data:    []byte{},
			wantErr: ErrPayloadTooShort,
		},
		{
			name:    "payload too short (header)",
			data:    []byte{byte(TypeString), 0x00},
			wantErr: ErrPayloadTooShort,
		},
		{
			name:    "invalid type",
			data:    []byte{0xff, 0x00, 0x00, 0x00, 0x00},
			wantErr: ErrInvalidType,
		},
		{
			name: "declared length exceeds payload",
			data: []byte{
				byte(TypeString),
				0x00, 0x00, 0x00, 0x05,
				'a', 'b',
			},
			wantErr: ErrPayloadTooShort,
		},
		{
			name: "empty string",
			data: []byte{
				byte(TypeString),
				0x00, 0x00, 0x00, 0x00,
			},
			want: "",
		},
		{
			name: "ascii string",
			data: []byte{
				byte(TypeString),
				0x00, 0x00, 0x00, 0x03,
				'a', 'b', 'c',
			},
			want: "abc",
		},
		{
			name: "utf8 string",
			data: []byte{
				byte(TypeString),
				0x00, 0x00, 0x00, 0x03,
				0xe2, 0x82, 0xac,
			},
			want: "€",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeString(tt.data)

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

			if got != tt.want {
				t.Fatalf("unexpected value: got %q want %q", got, tt.want)
			}
		})
	}
}
