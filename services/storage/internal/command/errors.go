package command

import "errors"

var (
	ErrCommandAlreadyRegistered = errors.New("handler for command ID is already registered")
)
