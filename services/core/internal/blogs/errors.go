package blogs

import "errors"

var (
	ErrArticleNotFound = errors.New("article not found")
	ErrTitleTooLong    = errors.New("title exceeds maximum length of 256 characters")
	ErrSummaryTooLong  = errors.New("summary exceeds maximum size of 1KB")
	ErrContentTooLong  = errors.New("content exceeds maximum size of 10MB")
)
