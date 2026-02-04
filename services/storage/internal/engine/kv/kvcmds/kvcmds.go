package kvcmds

import (
	"github.com/fancyinnovations/fancyspaces/storage/internal/command"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/protocol"
)

func Commands() map[uint16]command.Handler {
	return map[uint16]command.Handler{
		protocol.CommandKVGet: handleGet,
	}
}

func handleGet(ctx *command.ConnCtx, msg *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {

	return &protocol.Response{
		Code:    protocol.StatusOK,
		Payload: make([]byte, 0),
	}, nil
}
