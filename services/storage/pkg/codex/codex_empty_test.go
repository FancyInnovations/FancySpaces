package codex

import "testing"

func TestEncodeEmpty(t *testing.T) {
	out := EncodeEmpty()

	if len(out) != 1 {
		t.Fatalf("expected length 1, got %d", len(out))
	}

	if out[0] != byte(TypeEmpty) {
		t.Fatalf("expected type %d, got %d", TypeEmpty, out[0])
	}
}

func TestEncodeEmptyInto(t *testing.T) {
	tests := []struct {
		name          string
		target        []byte
		expectRealloc bool
	}{
		{
			name:          "nil target",
			target:        nil,
			expectRealloc: true,
		},
		{
			name:          "insufficient capacity",
			target:        make([]byte, 0, 0),
			expectRealloc: true,
		},
		{
			name:          "sufficient capacity",
			target:        make([]byte, 0, 4),
			expectRealloc: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origCap := cap(tt.target)
			out := EncodeEmptyInto(tt.target)

			if len(out) != 1 {
				t.Fatalf("expected length 1, got %d", len(out))
			}

			if out[0] != byte(TypeEmpty) {
				t.Fatalf("expected type %d, got %d", TypeEmpty, out[0])
			}

			if !tt.expectRealloc && cap(out) != origCap {
				t.Fatalf("unexpected reallocation: cap changed from %d to %d", origCap, cap(out))
			}
		})
	}
}
