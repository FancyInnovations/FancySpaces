package command

import (
	"context"
	"net"
)

// ConnCtx encapsulates the context for a client connection
type ConnCtx struct {
	// ID is a unique identifier for the connection, used for tracking and logging.
	ID string

	// Conn is the underlying network connection to the client.
	Conn net.Conn

	// Ctx is the context for managing the connection's lifecycle.
	Ctx context.Context

	// LastActivity tracks the last time a message was received from the client.
	LastActivity int64
}
