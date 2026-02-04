package client

import "errors"

var (
	ErrProtocolVersionNotSupported = errors.New("protocol version not supported")
	ErrClientNotConnected          = errors.New("client not connected")
	ErrUnexpectedStatusCode        = errors.New("unexpected status code")
	ErrInvalidCredentials          = errors.New("invalid credentials")
)
