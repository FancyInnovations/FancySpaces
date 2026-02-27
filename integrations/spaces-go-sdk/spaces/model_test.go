package spaces

import (
	"testing"
)

func TestSpace_Validate(t *testing.T) {
	tests := []struct {
		name    string
		space   Space
		wantErr error
	}{
		{
			name: "valid space",
			space: Space{
				Slug:        "validslug",
				Title:       "Valid Title",
				Description: "A valid description.",
			},
			wantErr: nil,
		},
		{
			name: "slug too short",
			space: Space{
				Slug:        "ab",
				Title:       "Valid Title",
				Description: "A valid description.",
			},
			wantErr: ErrSlugTooShort,
		},
		{
			name: "slug too long",
			space: Space{
				Slug:        "thisslugiswaytoolongtobevalid",
				Title:       "Valid Title",
				Description: "A valid description.",
			},
			wantErr: ErrSlugTooLong,
		},
		{
			name: "title too short",
			space: Space{
				Slug:        "validslug",
				Title:       "ab",
				Description: "A valid description.",
			},
			wantErr: ErrTitleTooShort,
		},
		{
			name: "title too long",
			space: Space{
				Slug:        "validslug",
				Title:       string(make([]byte, 101)),
				Description: "A valid description.",
			},
			wantErr: ErrTitleTooLong,
		},
		{
			name: "description too long",
			space: Space{
				Slug:        "validslug",
				Title:       "Valid Title",
				Description: string(make([]byte, 10001)),
			},
			wantErr: ErrDescriptionTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.space.Validate()
			if err != tt.wantErr {
				t.Errorf("got error %v, want %v", err, tt.wantErr)
			}
		})
	}
}
