package protocol

import (
	"encoding/binary"
	"net"
)

var V1 = &ProtoV1{}

type ProtoV1 struct {
}

// ReadFrame reads a length-prefixed frame from the connection.
// | Payload length (4 bytes)  |
// | Payload (variable length) |
func (p *ProtoV1) ReadFrame(conn net.Conn) ([]byte, error) {
	length := make([]byte, 4)
	_, err := conn.Read(length)
	if err != nil {
		return nil, err
	}

	frameLength := int(binary.BigEndian.Uint32(length))
	if frameLength <= 0 {
		return nil, ErrFrameLengthInvalid
	}

	frame := make([]byte, frameLength)
	_, err = conn.Read(frame)
	if err != nil {
		return nil, err
	}

	return frame, nil
}

// ReadMessage decodes a Message from the given byte slice.
// | Magic Number (1 byte)     |
// | Protocol Version (1 byte) |
// | Flags (1 byte)            |
// | Type (1 byte)             |
// | Payload length (4 bytes)  |
// | Payload (variable length) |
func (p *ProtoV1) ReadMessage(data []byte) (*Message, error) {
	if len(data) == 0 {
		return nil, ErrPayloadTooShort
	}

	if data[0] != magicNumber {
		return nil, ErrMagicNumberInvalid
	}

	if len(data) < 8 {
		return nil, ErrPayloadTooShort
	}

	msg := &Message{
		ProtocolVersion: data[1],
		Flags:           data[2],
		Type:            data[3],
	}
	if msg.ProtocolVersion != byte(ProtocolVersion1) {
		return nil, ErrInvalidProtocolVersion
	}
	if msg.Type > byte(MessageTypeError) {
		return nil, ErrUnknownMessageType
	}

	payloadLength := int(binary.BigEndian.Uint32(data[4:8]))
	if payloadLength == 0 {
		return nil, ErrEmptyPayload
	}
	if len(data) < 8+payloadLength {
		return nil, ErrPayloadTooShort
	}

	msg.Payload = data[8 : 8+payloadLength]

	return msg, nil
}

// ReadCommand decodes a Command from the given Message's payload.
// | CMD ID (2 bytes)                |
// | DB name length (2 byte)         |
// | DB name (variable)              |
// | Collection name length (2 byte) |
// | Collection name (variable)      |
// | Payload length (4 byte)         |
// | Payload (variable)              |
func (p *ProtoV1) ReadCommand(msg *Message) (*Command, error) {
	data := msg.Payload

	if len(data) < 10 {
		return nil, ErrPayloadTooShort
	}

	cmd := &Command{}
	cmd.ID = binary.BigEndian.Uint16(data[0:2])

	dbNameLen := int(binary.BigEndian.Uint16(data[2:4]))
	if len(data) < 4+dbNameLen+2 {
		return nil, ErrPayloadTooShort
	}
	cmd.DatabaseName = string(data[4 : 4+dbNameLen])

	collectionNameStart := 4 + dbNameLen
	collectionNameLen := int(binary.BigEndian.Uint16(data[collectionNameStart : collectionNameStart+2]))
	if len(data) < collectionNameStart+2+collectionNameLen+4 {
		return nil, ErrPayloadTooShort
	}
	cmd.CollectionName = string(data[collectionNameStart+2 : collectionNameStart+2+collectionNameLen])

	payloadStart := collectionNameStart + 2 + collectionNameLen
	payloadLen := int(binary.BigEndian.Uint32(data[payloadStart : payloadStart+4]))
	if len(data) < payloadStart+4+payloadLen {
		return nil, ErrPayloadTooShort
	}
	cmd.Payload = data[payloadStart+4 : payloadStart+4+payloadLen]

	return cmd, nil
}
