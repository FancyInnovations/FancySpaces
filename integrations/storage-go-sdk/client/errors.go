package client

import "errors"

var (
	ErrProtocolVersionNotSupported = errors.New("protocol version not supported")
	ErrClientNotConnected          = errors.New("client not connected")
	ErrUnexpectedStatusCode        = errors.New("unexpected status code")
	ErrInvalidCredentials          = errors.New("invalid credentials")
	ErrUnexpectedDataType          = errors.New("unexpected data type")
	ErrInvalidPayloadLength        = errors.New("invalid payload length")
	ErrCommandTimeout              = errors.New("command timed out")
	ErrDatabaseNotFound            = errors.New("database not found")
	ErrCollectionNotFound          = errors.New("collection not found")
	ErrKeyNotFound                 = errors.New("key not found")
	ErrInvalidEngine               = errors.New("invalid engine")
)
