package protocol

import (
	"encoding/binary"
	"io"
	"net"
)

const maxFrameSize = 16 * 1024 * 1024 // 16 MB

var V1 = &ProtoV1{
	Version: ProtocolVersion1,
}

type ProtoV1 struct {
	Version ProtocolVersion
}

// ReadFrameInto reads a length-prefixed frame from the connection into the provided target buffer.
// | Payload length (4 bytes)  |
// | Payload (variable length) |
func (p *ProtoV1) ReadFrameInto(conn net.Conn, target []byte) ([]byte, error) {
	var lenBuf [4]byte
	if _, err := io.ReadFull(conn, lenBuf[:]); err != nil {
		return nil, err
	}

	frameLength := int(binary.BigEndian.Uint32(lenBuf[:]))
	if frameLength <= 0 || frameLength > maxFrameSize {
		return nil, ErrFrameLengthInvalid
	}

	if cap(target) < frameLength {
		target = make([]byte, frameLength)
	} else {
		target = target[:frameLength]
	}

	if _, err := io.ReadFull(conn, target); err != nil {
		return nil, err
	}

	return target, nil
}

// ReadFrame reads a length-prefixed frame from the connection and returns it as a new byte slice.
func (p *ProtoV1) ReadFrame(conn net.Conn) ([]byte, error) {
	return p.ReadFrameInto(conn, nil)
}

// WriteFrame writes a length-prefixed frame to the connection.
// Uses the same format as ReadFrame.
func (p *ProtoV1) WriteFrame(conn net.Conn, data []byte) error {
	var lenBuf [4]byte
	binary.BigEndian.PutUint32(lenBuf[:], uint32(len(data)))

	_, err := conn.Write(append(lenBuf[:], data...))
	return err
}

// DecodeMessageInto decodes a Message from the given byte slice into the provided Message struct.
// | Magic Number (1 byte)     |
// | Protocol Version (1 byte) |
// | Flags (1 byte)            |
// | Type (1 byte)             |
// | Payload length (4 bytes)  |
// | Payload (variable length) |
func (p *ProtoV1) DecodeMessageInto(data []byte, target *Message) error {
	if len(data) == 0 {
		return ErrPayloadTooShort
	}

	if data[0] != magicNumber {
		return ErrMagicNumberInvalid
	}

	if len(data) < 8 {
		return ErrPayloadTooShort
	}

	if data[1] != byte(p.Version) {
		return ErrInvalidProtocolVersion
	}

	target.ProtocolVersion = data[1]
	target.Flags = data[2]
	target.Type = data[3]

	if target.ProtocolVersion != byte(ProtocolVersion1) {
		return ErrInvalidProtocolVersion
	}
	if target.Type > byte(MessageTypeResponse) {
		return ErrUnknownMessageType
	}

	payloadLength := int(binary.BigEndian.Uint32(data[4:8]))
	if payloadLength == 0 {
		return ErrEmptyPayload
	}
	if len(data) < 8+payloadLength {
		return ErrPayloadTooShort
	}

	target.Payload = data[8 : 8+payloadLength]

	return nil
}

// DecodeMessage decodes a Message from the given byte slice into a new Message struct.
func (p *ProtoV1) DecodeMessage(data []byte) (*Message, error) {
	msg := &Message{}
	err := p.DecodeMessageInto(data, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// EncodeMessageInto encodes a Message into the given byte slice buffer.
// Uses the same format as DecodeMessageInto.
func (p *ProtoV1) EncodeMessageInto(msg *Message, target []byte) []byte {
	payloadLength := len(msg.Payload)
	totalLength := 1 + 1 + 1 + 1 + 4 + payloadLength

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	target[0] = magicNumber
	target[1] = msg.ProtocolVersion
	target[2] = msg.Flags
	target[3] = msg.Type
	binary.BigEndian.PutUint32(target[4:8], uint32(payloadLength))
	copy(target[8:], msg.Payload)

	return target
}

// EncodeMessage encodes a Message into a new byte slice.
func (p *ProtoV1) EncodeMessage(msg *Message) []byte {
	return p.EncodeMessageInto(msg, nil)
}

// EncodeCommandInto encodes a Command into the given byte slice buffer.
// | CMD ID (2 bytes)                |
// | DB name length (2 byte)         |
// | DB name (variable)              |
// | Collection name length (2 byte) |
// | Collection name (variable)      |
// | Payload length (4 byte)         |
// | Payload (variable)              |
func (p *ProtoV1) EncodeCommandInto(cmd *Command, target []byte) []byte {
	dbNameLen := len(cmd.DatabaseName)
	collectionNameLen := len(cmd.CollectionName)
	payloadLen := len(cmd.Payload)

	totalLength := 2 + 2 + dbNameLen + 2 + collectionNameLen + 4 + payloadLen

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	binary.BigEndian.PutUint16(target[0:2], cmd.ID)
	binary.BigEndian.PutUint16(target[2:4], uint16(dbNameLen))
	copy(target[4:4+dbNameLen], []byte(cmd.DatabaseName))
	collectionNameStart := 4 + dbNameLen
	binary.BigEndian.PutUint16(target[collectionNameStart:collectionNameStart+2], uint16(collectionNameLen))
	copy(target[collectionNameStart+2:collectionNameStart+2+collectionNameLen], []byte(cmd.CollectionName))
	payloadStart := collectionNameStart + 2 + collectionNameLen
	binary.BigEndian.PutUint32(target[payloadStart:payloadStart+4], uint32(payloadLen))
	copy(target[payloadStart+4:payloadStart+4+payloadLen], cmd.Payload)

	return target
}

// EncodeCommand encodes a Command into a new byte slice.
func (p *ProtoV1) EncodeCommand(cmd *Command) []byte {
	return p.EncodeCommandInto(cmd, nil)
}

// DecodeCommandInto decodes a Command from the given Message's payload into the provided Command struct.
func (p *ProtoV1) DecodeCommandInto(msg *Message, target *Command) error {
	data := msg.Payload

	if len(data) < 10 {
		return ErrPayloadTooShort
	}

	target.ID = binary.BigEndian.Uint16(data[0:2])

	dbNameLen := int(binary.BigEndian.Uint16(data[2:4]))
	if len(data) < 4+dbNameLen+2 {
		return ErrPayloadTooShort
	}
	target.DatabaseName = string(data[4 : 4+dbNameLen])

	collectionNameStart := 4 + dbNameLen
	collectionNameLen := int(binary.BigEndian.Uint16(data[collectionNameStart : collectionNameStart+2]))
	if len(data) < collectionNameStart+2+collectionNameLen+4 {
		return ErrPayloadTooShort
	}
	target.CollectionName = string(data[collectionNameStart+2 : collectionNameStart+2+collectionNameLen])

	payloadStart := collectionNameStart + 2 + collectionNameLen
	payloadLen := int(binary.BigEndian.Uint32(data[payloadStart : payloadStart+4]))
	if len(data) < payloadStart+4+payloadLen {
		return ErrPayloadTooShort
	}
	target.Payload = data[payloadStart+4 : payloadStart+4+payloadLen]

	return nil
}

// DecodeCommand decodes a Command from the given Message's payload into a new Command struct.
func (p *ProtoV1) DecodeCommand(msg *Message) (*Command, error) {
	cmd := &Command{}
	err := p.DecodeCommandInto(msg, cmd)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

// EncodeResponseInto encodes a Response into a byte slice buffer.
// | Status Code (2 bytes)     |
// | Payload length (4 bytes)  |
// | Payload (variable length) |
func (p *ProtoV1) EncodeResponseInto(resp *Response, target []byte) []byte {
	payloadLength := len(resp.Payload)
	totalLength := 2 + 4 + payloadLength

	// ensure capacity
	if cap(target) < totalLength {
		target = make([]byte, totalLength)
	} else {
		target = target[:totalLength]
	}

	binary.BigEndian.PutUint16(target[0:2], resp.Code)
	binary.BigEndian.PutUint32(target[2:6], uint32(payloadLength))
	copy(target[6:], resp.Payload)

	return target
}

// EncodeResponse encodes a Response into a new byte slice.
func (p *ProtoV1) EncodeResponse(resp *Response) []byte {
	return p.EncodeResponseInto(resp, nil)
}

// DecodeResponseInto decodes a Response from the given Message's payload into the provided Response struct.
func (p *ProtoV1) DecodeResponseInto(msg *Message, target *Response) error {
	data := msg.Payload

	if len(data) < 6 {
		return ErrPayloadTooShort
	}

	target.Code = binary.BigEndian.Uint16(data[0:2])

	payloadLength := int(binary.BigEndian.Uint32(data[2:6]))
	if len(data) < 6+payloadLength {
		return ErrPayloadTooShort
	}

	target.Payload = data[6 : 6+payloadLength]

	return nil
}

// DecodeResponse decodes a Response from the given Message's payload into a new Response struct.
func (p *ProtoV1) DecodeResponse(msg *Message) (*Response, error) {
	resp := &Response{}
	err := p.DecodeResponseInto(msg, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
