package server

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"time"

	"github.com/OliverSchlueter/goutils/idgen"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/storage/internal/command"
	"github.com/fancyinnovations/fancyspaces/storage/internal/protocol"
)

type Server struct {
	addr        string
	listener    net.Listener
	cmdService  *command.Service
	connections map[string]*command.ConnCtx
}

type Configuration struct {
	Addr       string
	CmdService *command.Service
}

func New(cfg Configuration) *Server {
	return &Server{
		addr:        cfg.Addr,
		cmdService:  cfg.CmdService,
		connections: make(map[string]*command.ConnCtx),
	}
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.listener = ln

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	ctx := &command.ConnCtx{
		ID:   idgen.GenerateID(16),
		Conn: conn,
		Ctx:  context.Background(),
	}
	s.connections[ctx.ID] = ctx
	defer delete(s.connections, ctx.ID)

	slog.Info("New connection established", slog.String("ConnID", ctx.ID), slog.String("RemoteAddr", conn.RemoteAddr().String()))

	for {
		startTime := time.Now()

		frame, err := protocol.V1.ReadFrame(conn)
		if err != nil {
			if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
				slog.Info("Connection closed by client", slog.String("ConnID", ctx.ID))
				return
			}

			if errors.Is(err, protocol.ErrFrameLengthInvalid) {
				s.writeResponse(conn, &protocol.Response{
					Code:    protocol.StatusInvalidMessage,
					Payload: []byte(err.Error()),
				})
			} else {
				slog.Error("Failed to read frame", slog.String("ConnID", ctx.ID), sloki.WrapError(err))
			}

			continue
		}

		msg, err := protocol.V1.DecodeMessage(frame)
		if err != nil {
			s.writeResponse(conn, &protocol.Response{
				Code:    protocol.StatusInvalidMessage,
				Payload: []byte(err.Error()),
			})

			continue
		}

		if msg.Type != byte(protocol.MessageTypeCommand) {
			s.writeResponse(conn, &protocol.Response{
				Code:    protocol.StatusInvalidMessage,
				Payload: []byte("Only command messages are allowed"),
			})

			continue
		}

		cmd, err := protocol.V1.DecodeCommand(msg)
		if err != nil {
			s.writeResponse(conn, &protocol.Response{
				Code:    protocol.StatusInvalidMessage,
				Payload: []byte(err.Error()),
			})

			continue
		}

		resp, err := s.cmdService.Handle(ctx, msg, cmd)
		if err != nil {
			slog.Warn("Command handler returned error", sloki.WrapError(err))
			s.writeResponse(conn, &protocol.Response{
				Code:    protocol.StatusInternalServerError,
				Payload: []byte("Error while processing command"),
			})

			continue
		}

		s.writeResponse(conn, resp)

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
	}
}

func (s *Server) writeResponse(conn net.Conn, resp *protocol.Response) {
	// TODO: use message pool and just update the payload

	msg := protocol.Message{
		ProtocolVersion: byte(protocol.ProtocolVersion1),
		Flags:           0x00,
		Type:            byte(protocol.MessageTypeResponse),
		Payload:         protocol.V1.EncodeResponse(resp),
	}

	data := protocol.V1.EncodeMessage(&msg)
	if err := protocol.V1.WriteFrame(conn, data); err != nil {
		slog.Warn("Failed to write response", sloki.WrapError(err))
	}
}

func (s *Server) broadcastMessage(msg *protocol.Message) {
	data := protocol.V1.EncodeMessage(msg)

	for _, ctx := range s.connections {
		if err := protocol.V1.WriteFrame(ctx.Conn, data); err != nil {
			slog.Warn("Failed to broadcast message to connection", slog.String("ConnID", ctx.ID), sloki.WrapError(err))
		}
	}
}
