package command

import "github.com/fancyinnovations/fancyspaces/storage/internal/protocol"

func handlePing(_ *protocol.Message, _ *protocol.Command) (*protocol.Response, error) {
	return &protocol.Response{
		Code:    protocol.StatusOK,
		Payload: []byte("pong"),
	}, nil
}
