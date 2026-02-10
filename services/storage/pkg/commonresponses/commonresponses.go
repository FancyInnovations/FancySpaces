package commonresponses

import "github.com/fancyinnovations/fancyspaces/storage/pkg/protocol"

var EmptyPayload = &[]byte{}

var (
	OK = &protocol.Response{
		Code:    protocol.StatusOK,
		Payload: *EmptyPayload,
	}

	Unauthorized = &protocol.Response{
		Code:    protocol.StatusUnauthorized,
		Payload: *EmptyPayload,
	}

	Forbidden = &protocol.Response{
		Code:    protocol.StatusForbidden,
		Payload: *EmptyPayload,
	}

	DatabaseNotFound = &protocol.Response{
		Code:    protocol.StatusDatabaseNotFound,
		Payload: *EmptyPayload,
	}

	CollectionNotFound = &protocol.Response{
		Code:    protocol.StatusCollectionNotFound,
		Payload: *EmptyPayload,
	}

	CommandNotAllowed = &protocol.Response{
		Code:    protocol.StatusCommandNotAllowed,
		Payload: *EmptyPayload,
	}

	InternalServerError = &protocol.Response{
		Code:    protocol.StatusInternalServerError,
		Payload: *EmptyPayload,
	}
)
