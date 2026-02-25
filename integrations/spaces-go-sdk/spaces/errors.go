package spaces

import (
	"errors"
)

var (
	ErrSlugTooLong        = errors.New("slug exceeds maximum length of 20 characters")
	ErrSlugTooShort       = errors.New("slug must be at least 3 characters long")
	ErrSlugInvalidFormat  = errors.New("slug must consist of lowercase letters, numbers, and hyphens, and cannot start or end with a hyphen")
	ErrTitleTooLong       = errors.New("title exceeds maximum length of 100 characters")
	ErrTitleTooShort      = errors.New("title must be at least 3 characters long")
	ErrSummaryTooLong     = errors.New("summary exceeds maximum length of 300 characters")
	ErrDescriptionTooLong = errors.New("description exceeds maximum length of 10000 characters")
	ErrSpaceNotFound      = errors.New("space not found")
)
