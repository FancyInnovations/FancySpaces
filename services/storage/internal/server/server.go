package server

import (
	"errors"
	"io"
	"log/slog"
	"net"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/storage/internal/protocol"
)

type Server struct {
	addr     string
	listener net.Listener
}

func New(addr string) *Server {
	return &Server{
		addr: addr,
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

	for {
		frame, err := protocol.V1.ReadFrame(conn)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				slog.Info("Connection closed by client")
				return
			}

			if errors.Is(err, io.EOF) {
				slog.Info("Connection closed by client (EOF)")
				return
			}

			if errors.Is(err, protocol.ErrFrameLengthInvalid) {
				s.writeResponse(conn, &protocol.Response{
					Code:    protocol.StatusInvalidMessage,
					Payload: []byte(err.Error()),
				})
			} else {
				// could because by connection closed by client
				slog.Error("Failed to read frame", sloki.WrapError(err))
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

		slog.Info(
			"Received command",
			slog.Int("CommandID", int(cmd.ID)),
			slog.String("Database", cmd.DatabaseName),
			slog.String("Collection", cmd.CollectionName),
			slog.String("Payload", string(cmd.Payload)),
		)

		// TODO: process command and send response / error back to client
		s.writeResponse(conn, &protocol.Response{
			Code:    protocol.StatusOK,
			Payload: []byte("Command received successfully"),
		})
	}
}

func (s *Server) writeResponse(conn net.Conn, resp *protocol.Response) {
	msg := protocol.Message{
		ProtocolVersion: byte(protocol.ProtocolVersion1),
		Flags:           0x00,
		Type:            byte(protocol.MessageTypeResponse),
		Payload:         protocol.V1.EncodeResponse(resp),
	}

	data := protocol.V1.EncodeMessage(&msg)
	if err := protocol.V1.WriteFrame(conn,
		data); err != nil {
		slog.Warn("Failed to write response", sloki.WrapError(err))
	}
}
