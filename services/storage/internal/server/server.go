package server

import (
	"net"

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
			continue // TODO write error
		}

		msg, err := protocol.V1.ReadMessage(frame)
		if err != nil {
			continue // TODO write error
		}

		cmd, err := protocol.V1.ReadCommand(msg)
		if err != nil {
			continue // TODO write error
		}

		_ = cmd

		// TODO: process command and send response / error back to client
	}
}
