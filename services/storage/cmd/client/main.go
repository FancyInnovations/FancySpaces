package main

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"net"
	"strconv"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/storage/internal/protocol"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8091")
	if err != nil {
		slog.Error("Could not connect to server", slog.String("address", "localhost:8091"), slog.Any("error", err))
		return
	}
	defer conn.Close()

	cmd := &protocol.Command{
		ID:             42,
		DatabaseName:   "fancyspaces",
		CollectionName: "spaces",
		Payload:        []byte("HELLO WOLRD!"),
	}
	if err := SendCmd(conn, cmd); err != nil {
		slog.Error("Could not send command to server", slog.Any("error", err))
		return
	}

	ReadResponse(conn)
}

func SendCmd(conn net.Conn, cmd *protocol.Command) error {
	msg := protocol.Message{
		ProtocolVersion: byte(protocol.ProtocolVersion1),
		Flags:           0x00,
		Type:            byte(protocol.MessageTypeCommand),
		Payload:         EncodeCommand(cmd),
	}

	data := protocol.V1.EncodeMessage(&msg)
	if err := protocol.V1.WriteFrame(conn,
		data); err != nil {
		slog.Warn("Failed to write cmd", sloki.WrapError(err))
	}
	return nil
}

func ReadResponse(conn net.Conn) {
	frame, err := protocol.V1.ReadFrame(conn)
	if err != nil {
		slog.Error("Could not read frame from server", slog.Any("error", err))
		return
	}
	msg, err := protocol.V1.DecodeMessage(frame)
	if err != nil {
		slog.Error("Could not read message from server", slog.Any("error", err))
		return
	}
	resp, err := DecodeResponse(msg.Payload)
	if err != nil {
		slog.Error("Could not decode response from server", slog.Any("error", err))
		return
	}

	slog.Info(
		"Received response from server",
		slog.String("message_type", strconv.Itoa(int(msg.Type))),
		slog.String("response_code", strconv.Itoa(int(resp.Code))),
		slog.String("response_payload", string(resp.Payload)),
	)
}

func EncodeCommand(cmd *protocol.Command) []byte {
	dbNameLen := len(cmd.DatabaseName)
	collectionNameLen := len(cmd.CollectionName)
	payloadLen := len(cmd.Payload)

	totalLength := 2 + 2 + dbNameLen + 2 + collectionNameLen + 4 + payloadLen
	data := make([]byte, totalLength)

	binary.BigEndian.PutUint16(data[0:2], cmd.ID)
	binary.BigEndian.PutUint16(data[2:4], uint16(dbNameLen))
	copy(data[4:4+dbNameLen], []byte(cmd.DatabaseName))
	collectionNameStart := 4 + dbNameLen
	binary.BigEndian.PutUint16(data[collectionNameStart:collectionNameStart+2], uint16(collectionNameLen))
	copy(data[collectionNameStart+2:collectionNameStart+2+collectionNameLen], []byte(cmd.CollectionName))
	payloadStart := collectionNameStart + 2 + collectionNameLen
	binary.BigEndian.PutUint32(data[payloadStart:payloadStart+4], uint32(payloadLen))
	copy(data[payloadStart+4:], cmd.Payload)

	return data
}

func DecodeResponse(data []byte) (*protocol.Response, error) {
	if len(data) < 6 {
		return nil, fmt.Errorf("data too short to be a valid response")
	}

	resp := &protocol.Response{}
	resp.Code = binary.BigEndian.Uint16(data[0:2])

	payloadLength := int(binary.BigEndian.Uint32(data[2:6]))
	if len(data) < 6+payloadLength {
		return nil, fmt.Errorf("data too short for declared payload length")
	}

	resp.Payload = data[6 : 6+payloadLength]

	return resp, nil
}
