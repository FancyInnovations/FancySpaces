package server

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/OliverSchlueter/goutils/idgen"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/storage/internal/command"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/commonresponses"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/protocol"
)

type Server struct {
	addr          string
	listener      net.Listener
	cmdService    *command.Service
	connections   map[string]*command.ConnCtx
	connectionsMu sync.Mutex
}

type Configuration struct {
	Addr string
}

func New(cfg Configuration) *Server {
	return &Server{
		addr:        cfg.Addr,
		connections: make(map[string]*command.ConnCtx),
	}
}

func (s *Server) SetCommandService(cmdService *command.Service) {
	if s.cmdService != nil {
		slog.Warn("Command service is already set, ignoring new service")
		return
	}
	s.cmdService = cmdService
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.listener = ln

	go s.cleanupInactiveConnections()

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go s.handleConnection(conn)
	}
}

// handleConnection manages the lifecycle of a single client connection.
// It reads messages in a loop until the connection is closed.
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	ctx := &command.ConnCtx{
		ID:   idgen.GenerateID(16),
		Conn: conn,
		Ctx:  context.Background(),
	}
	s.connectionsMu.Lock()
	s.connections[ctx.ID] = ctx
	s.connectionsMu.Unlock()

	defer func() {
		s.connectionsMu.Lock()
		delete(s.connections, ctx.ID)
		s.connectionsMu.Unlock()
	}()

	slog.Info("New connection established", slog.String("ConnID", ctx.ID), slog.String("RemoteAddr", conn.RemoteAddr().String()))

	for {
		if s.handleMessage(ctx) {
			break
		}
	}
}

// handleMessage reads and processes a single message from the connection.
// It returns true if the connection should be closed.
func (s *Server) handleMessage(ctx *command.ConnCtx) bool {
	conn := ctx.Conn

	frameBuf := protocol.GetRequestBufferFromPool()
	defer protocol.PutRequestBufferToPool(frameBuf)

	frame, err := protocol.V1.ReadFrameInto(conn, frameBuf)
	if err != nil {
		if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
			slog.Info("Connection closed by client", slog.String("ConnID", ctx.ID))
			return true
		}

		if errors.Is(err, protocol.ErrFrameLengthInvalid) {
			s.writeResponse(conn, &protocol.Response{
				Code:    protocol.StatusInvalidMessage,
				Payload: []byte(err.Error()),
			})
		} else {
			slog.Error("Failed to read frame", slog.String("ConnID", ctx.ID), sloki.WrapError(err))
		}

		return false
	}

	startTime := time.Now()

	msg := protocol.GetMessageFromPool()
	defer protocol.PutMessageToPool(msg)
	if err := protocol.V1.DecodeMessageInto(frame, msg); err != nil {
		s.writeResponse(conn, &protocol.Response{
			Code:    protocol.StatusInvalidMessage,
			Payload: []byte(err.Error()),
		})

		return false
	}

	if msg.Type != byte(protocol.MessageTypeCommand) {
		s.writeResponse(conn, &protocol.Response{
			Code:    protocol.StatusInvalidMessage,
			Payload: []byte("Only command messages are allowed"),
		})

		return false
	}

	cmd := protocol.GetCommandFromPool()
	defer protocol.PutCommandToPool(cmd)
	if err := protocol.V1.DecodeCommandInto(msg, cmd); err != nil {
		s.writeResponse(conn, &protocol.Response{
			Code:    protocol.StatusInvalidMessage,
			Payload: []byte(err.Error()),
		})

		return false
	}

	resp, err := s.cmdService.Handle(ctx, msg, cmd)
	if err != nil {
		slog.Warn("Command handler returned error", sloki.WrapError(err))
		s.writeResponse(conn, commonresponses.InternalServerError)

		return false
	}

	s.writeResponse(conn, resp)

	// Update last activity timestamp for cleanup purposes
	ctx.LastActivity = startTime.UnixMilli()

	elapsedTime := time.Since(startTime)

	slog.Info(
		"Processed command",
		slog.String("ConnID", ctx.ID),
		slog.Int("CommandID", int(cmd.ID)),
		slog.String("Database", cmd.DatabaseName),
		slog.String("Collection", cmd.CollectionName),
		slog.String("Payload", string(cmd.Payload)),
		slog.Duration("Duration", elapsedTime),
	)

	return false
}

func (s *Server) writeResponse(conn net.Conn, resp *protocol.Response) {
	msg := protocol.GetMessageFromPool()
	defer protocol.PutMessageToPool(msg)

	payloadBuf := protocol.GetResponseBufferFromPool()
	defer protocol.PutResponseBufferToPool(payloadBuf)

	msg.ProtocolVersion = byte(protocol.ProtocolVersion1)
	msg.Flags = 0x00
	msg.Type = byte(protocol.MessageTypeResponse)
	msg.Payload = protocol.V1.EncodeResponseInto(resp, payloadBuf)

	msgDataBuf := protocol.GetResponseBufferFromPool()
	defer protocol.PutResponseBufferToPool(msgDataBuf)

	msgDataBuf = protocol.V1.EncodeMessageInto(msg, msgDataBuf)
	if err := protocol.V1.WriteFrame(conn, msgDataBuf); err != nil {
		slog.Warn("Failed to write response", sloki.WrapError(err))
	}
}

func (s *Server) broadcastMessage(msg *protocol.Message) {
	msgDataBuf := protocol.GetResponseBufferFromPool()
	defer protocol.PutResponseBufferToPool(msgDataBuf)

	protocol.V1.EncodeMessageInto(msg, msgDataBuf)

	for _, ctx := range s.connections {
		if err := protocol.V1.WriteFrame(ctx.Conn, msgDataBuf); err != nil {
			slog.Warn("Failed to broadcast message to connection", slog.String("ConnID", ctx.ID), sloki.WrapError(err))
		}
	}
}

func (s *Server) SendBrokerMessage(connID, subject string, msgs [][]byte) {

}
