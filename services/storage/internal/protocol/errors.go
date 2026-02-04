package protocol

import "errors"

var (
	ErrFrameLengthInvalid = errors.New("invalid frame length")
	ErrMagicNumberInvalid = errors.New("invalid magic number")
	ErrPayloadTooShort    = errors.New("payload too short")
	ErrEmptyPayload       = errors.New("empty payload")

	ErrInvalidProtocolVersion = errors.New("invalid protocol version")
	ErrUnknownMessageType     = errors.New("unknown message type")
)
