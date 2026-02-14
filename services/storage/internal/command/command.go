package command

import (
	"fmt"

	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/protocol"
)

type Handler func(ctx *ConnCtx, msg *protocol.Message, cmd *protocol.Command) (*protocol.Response, error)
type Middleware func(Handler) Handler

type Service struct {
	handlers    map[uint16]Handler
	middlewares []Middleware
}

func NewService() *Service {
	return &Service{
		handlers:    make(map[uint16]Handler),
		middlewares: []Middleware{},
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

func (s *Service) RegisterMiddleware(mw Middleware) {
	s.middlewares = append(s.middlewares, mw)
}

func (s *Service) RegisterMiddlewares(mws []Middleware) {
	for _, mw := range mws {
		s.RegisterMiddleware(mw)
	}
}

func (s *Service) Handle(ctx *ConnCtx, msg *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
	handler, exists := s.handlers[cmd.ID]
	if !exists {
		return &protocol.Response{
			Code:    protocol.StatusCommandNotFound,
			Payload: []byte(fmt.Sprintf("command with ID %d not found", cmd.ID)),
		}, nil
	}

	for i := len(s.middlewares) - 1; i >= 0; i-- {
		handler = s.middlewares[i](handler)
	}

	return handler(ctx, msg, cmd)
}
