package spaces

import (
	"errors"
)

var (
	ErrSlugTooLong        = errors.New("slug exceeds maximum length of 20 characters")
	ErrSlugTooShort       = errors.New("slug must be at least 3 characters long")
	ErrTitleTooLong       = errors.New("title exceeds maximum length of 100 characters")
	ErrTitleTooShort      = errors.New("title must be at least 3 characters long")
	ErrDescriptionTooLong = errors.New("description exceeds maximum length of 10000 characters")
	ErrSpaceNotFound      = errors.New("space not found")
)
