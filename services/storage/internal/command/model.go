package command

import (
	"context"
	"net"
)

type ConnCtx struct {
	ID   string
	Conn net.Conn
	Ctx  context.Context
}
