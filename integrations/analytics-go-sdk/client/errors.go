package client

import "errors"

var (
	ErrNoProjectID          = errors.New("no project ID provided")
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
)
