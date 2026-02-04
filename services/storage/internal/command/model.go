package command

import (
	"context"
	"net"
)

type ConnCtx struct {
	Conn net.Conn
	Ctx  context.Context
}
