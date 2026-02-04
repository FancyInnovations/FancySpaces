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

// WriteFrame writes a length-prefixed frame to the connection.
// Uses the same format as ReadFrame.
func (p *ProtoV1) WriteFrame(conn net.Conn, data []byte) error {
	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(len(data)))

	_, err := conn.Write(length)
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}

// DecodeMessage decodes a Message from the given byte slice.
// | Magic Number (1 byte)     |
// | Protocol Version (1 byte) |
// | Flags (1 byte)            |
// | Type (1 byte)             |
// | Payload length (4 bytes)  |
// | Payload (variable length) |
func (p *ProtoV1) DecodeMessage(data []byte) (*Message, error) {
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
	if msg.Type > byte(MessageTypeResponse) {
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

// EncodeMessage encodes a Message into a byte slice.
// Uses the same format as DecodeMessage.
func (p *ProtoV1) EncodeMessage(msg *Message) []byte {
	payloadLength := len(msg.Payload)
	totalLength := 1 + 1 + 1 + 1 + 4 + payloadLength
	data := make([]byte, totalLength)

	data[0] = magicNumber
	data[1] = msg.ProtocolVersion
	data[2] = msg.Flags
	data[3] = msg.Type
	binary.BigEndian.PutUint32(data[4:8], uint32(payloadLength))
	copy(data[8:], msg.Payload)

	return data
}

// DecodeCommand decodes a Command from the given Message's payload.
// | CMD ID (2 bytes)                |
// | DB name length (2 byte)         |
// | DB name (variable)              |
// | Collection name length (2 byte) |
// | Collection name (variable)      |
// | Payload length (4 byte)         |
// | Payload (variable)              |
func (p *ProtoV1) DecodeCommand(msg *Message) (*Command, error) {
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

// EncodeResponse encodes a Response into a byte slice.
// | Status Code (2 bytes)     |
// | Payload length (4 bytes)  |
// | Payload (variable length) |
func (p *ProtoV1) EncodeResponse(resp *Response) []byte {
	payloadLength := len(resp.Payload)
	totalLength := 2 + 4 + payloadLength
	data := make([]byte, totalLength)

	binary.BigEndian.PutUint16(data[0:2], resp.Code)
	binary.BigEndian.PutUint32(data[2:6], uint32(payloadLength))
	copy(data[6:], resp.Payload)

	return data
}
