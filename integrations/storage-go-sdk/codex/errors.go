package codex

import "errors"

var (
	ErrPayloadTooShort = errors.New("payload too short")
	ErrInvalidType     = errors.New("invalid type")
)
