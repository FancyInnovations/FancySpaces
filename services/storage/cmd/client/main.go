package main

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"net"
	"strconv"

	"github.com/fancyinnovations/fancyspaces/storage/internal/protocol"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	conn, err := net.Dial("tcp", "localhost:8091")
	if err != nil {
		slog.Error("Could not connect to server", slog.String("address", "localhost:8091"), slog.Any("error", err))
		return
	}
	defer conn.Close()

	if err := Ping(conn); err != nil {
		slog.Error("Ping failed", slog.Any("error", err))
		return
	}
	slog.Info("Ping successful")

	if err := Login(conn, "oliver", "hello"); err != nil {
		slog.Error("Login failed", slog.Any("error", err))
		return
	}
	slog.Info("Login successful")

	isAuth, err := IsAuthenticated(conn)
	if err != nil {
		slog.Error("Could not check authentication status", slog.Any("error", err))
		return
	}
	slog.Info("Authentication status", slog.Bool("is_authenticated", isAuth))
}

func SendCmd(conn net.Conn, cmd *protocol.Command) error {
	msg := protocol.Message{
		ProtocolVersion: byte(protocol.ProtocolVersion1),
		Flags:           0x00,
		Type:            byte(protocol.MessageTypeCommand),
		Payload:         EncodeCommand(cmd),
	}

	data := protocol.V1.EncodeMessage(&msg)
	if err := protocol.V1.WriteFrame(conn, data); err != nil {
		return err
	}

	slog.Debug(
		"Sent command to server",
		slog.String("command_id", strconv.Itoa(int(cmd.ID))),
		slog.String("command_payload", string(cmd.Payload)),
		slog.String("message_size", strconv.Itoa(len(data))),
	)

	return nil
}

func ReadResponse(conn net.Conn) (*protocol.Response, error) {
	frame, err := protocol.V1.ReadFrame(conn)
	if err != nil {
		return nil, err
	}
	msg, err := protocol.V1.DecodeMessage(frame)
	if err != nil {
		return nil, err
	}
	resp, err := DecodeResponse(msg.Payload)
	if err != nil {
		return nil, err
	}

	slog.Debug(
		"Received response from server",
		slog.String("message_type", strconv.Itoa(int(msg.Type))),
		slog.String("response_code", strconv.Itoa(int(resp.Code))),
		slog.String("response_payload", string(resp.Payload)),
	)

	return resp, nil
}

func Ping(conn net.Conn) error {
	cmd := &protocol.Command{
		ID:             protocol.CommandPing,
		DatabaseName:   "",
		CollectionName: "",
		Payload:        []byte{},
	}
	if err := SendCmd(conn, cmd); err != nil {
		return err
	}
	resp, err := ReadResponse(conn)
	if err != nil {
		return err
	}

	if resp.Code != protocol.StatusOK {
		return fmt.Errorf("unexpected response code: %d", resp.Code)
	}

	return nil
}

func Login(conn net.Conn, username, password string) error {
	loginPayload := []byte{
		0x01, // type: password
	}
	loginPayload = append(loginPayload, 0, byte(len(username))) // username length (2 bytes)
	loginPayload = append(loginPayload, []byte(username)...)    // username
	loginPayload = append(loginPayload, 0, byte(len(password))) // password length (2 bytes)
	loginPayload = append(loginPayload, []byte(password)...)    // password

	cmd := &protocol.Command{
		ID:             protocol.CommandLogin,
		DatabaseName:   "",
		CollectionName: "",
		Payload:        loginPayload,
	}
	if err := SendCmd(conn, cmd); err != nil {
		return err
	}

	resp, err := ReadResponse(conn)
	if err != nil {
		return err
	}

	if resp.Code != protocol.StatusOK {
		if resp.Code == protocol.StatusInvalidCredentials {
			return fmt.Errorf("invalid credentials")
		}

		return fmt.Errorf("unexpected response code: %d", resp.Code)
	}

	return nil
}

func IsAuthenticated(conn net.Conn) (bool, error) {
	cmd := &protocol.Command{
		ID:             protocol.CommandAuthStatus,
		DatabaseName:   "",
		CollectionName: "",
		Payload:        []byte{},
	}
	if err := SendCmd(conn, cmd); err != nil {
		return false, err
	}
	resp, err := ReadResponse(conn)
	if err != nil {
		return false, err
	}

	if resp.Code == protocol.StatusOK {
		return true, nil
	}

	if resp.Code == protocol.StatusUnauthorized {
		return false, nil
	}

	return false, fmt.Errorf("unexpected response code: %d", resp.Code)
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
