package command

import (
	"fmt"

	"github.com/fancyinnovations/fancyspaces/storage/pkg/protocol"
)

type Handler func(ctx *ConnCtx, msg *protocol.Message, cmd *protocol.Command) (*protocol.Response, error)

type Service struct {
	handlers map[uint16]Handler
}

func NewService() *Service {
	return &Service{
		handlers: make(map[uint16]Handler),
	}
}

func (s *Service) RegisterHandler(commandID uint16, handler Handler) {
	if _, exists := s.handlers[commandID]; exists {
		return
	}

	s.handlers[commandID] = handler
}

func (s *Service) RegisterHandlers(handlers map[uint16]Handler) {
	for commandID, handler := range handlers {
		s.RegisterHandler(commandID, handler)
	}
}

func (s *Service) Handle(ctx *ConnCtx, msg *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
	if handler, exists := s.handlers[cmd.ID]; exists {
		return handler(ctx, msg, cmd)
	}

	return &protocol.Response{
		Code:    protocol.StatusCommandNotFound,
		Payload: []byte(fmt.Sprintf("command with ID %d not found", cmd.ID)),
	}, nil
}
