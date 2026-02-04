package command

import (
	"fmt"

	"github.com/fancyinnovations/fancyspaces/storage/internal/protocol"
)

type Handler func(ctx *ConnCtx, msg *protocol.Message, cmd *protocol.Command) (*protocol.Response, error)

type Service struct {
	handlers map[uint16]Handler
}

func NewService() *Service {
	handlers := make(map[uint16]Handler)

	// Register system command handlers
	handlers[protocol.CommandPing] = handlePing
	handlers[protocol.CommandSupportedProtocolVersions] = handleSupportedProtocolVersions
	handlers[protocol.CommandLogin] = handleLogin
	handlers[protocol.CommandAuthStatus] = handleAuthStatus

	return &Service{
		handlers: handlers,
	}
}

func (s *Service) RegisterHandler(commandID uint16, handler Handler) error {
	if _, exists := s.handlers[commandID]; exists {
		return ErrCommandAlreadyRegistered
	}

	s.handlers[commandID] = handler
	return nil
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
