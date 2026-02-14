package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"testing"
	"time"
)

// readerConn is a simple net.Conn wrapper that reads from a bytes.Reader.
type readerConn struct {
	r *bytes.Reader
}

func (c *readerConn) Read(p []byte) (int, error)         { return c.r.Read(p) }
func (c *readerConn) Write(p []byte) (int, error)        { return len(p), nil }
func (c *readerConn) Close() error                       { return nil }
func (c *readerConn) LocalAddr() net.Addr                { return &net.TCPAddr{} }
func (c *readerConn) RemoteAddr() net.Addr               { return &net.TCPAddr{} }
func (c *readerConn) SetDeadline(t time.Time) error      { return nil }
func (c *readerConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *readerConn) SetWriteDeadline(t time.Time) error { return nil }

// errConn injects an error on a specific Read call (1-based).
type errConn struct {
	r           *bytes.Reader
	errOn       int
	callCount   int
	errToReturn error
}

func (c *errConn) Read(p []byte) (int, error) {
	c.callCount++
	if c.callCount == c.errOn {
		if c.errToReturn != nil {
			return 0, c.errToReturn
		}
		return 0, errors.New("injected read error")
	}
	return c.r.Read(p)
}
func (c *errConn) Write(p []byte) (int, error)        { return len(p), nil }
func (c *errConn) Close() error                       { return nil }
func (c *errConn) LocalAddr() net.Addr                { return &net.TCPAddr{} }
func (c *errConn) RemoteAddr() net.Addr               { return &net.TCPAddr{} }
func (c *errConn) SetDeadline(t time.Time) error      { return nil }
func (c *errConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *errConn) SetWriteDeadline(t time.Time) error { return nil }

func TestReadFrame_Success(t *testing.T) {
	payload := []byte("hello")
	header := []byte{0, 0, 0, byte(len(payload))}
	data := append(header, payload...)
	c := &readerConn{r: bytes.NewReader(data)}

	frame, err := V1.ReadFrame(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(frame, payload) {
		t.Fatalf("frame mismatch: got %v want %v", frame, payload)
	}
}

func TestReadFrame_ErrorOnFirstRead(t *testing.T) {
	c := &errConn{
		r:           bytes.NewReader(nil),
		errOn:       1,
		errToReturn: errors.New("first read failed"),
	}

	_, err := V1.ReadFrame(c)
	if err == nil {
		t.Fatal("expected error on first read, got nil")
	}
}

func TestReadFrame_ErrorOnSecondRead(t *testing.T) {
	payloadLen := 5
	header := []byte{0, 0, 0, byte(payloadLen)}
	// Provide header so first read succeeds; second read will error.
	c := &errConn{
		r:           bytes.NewReader(header),
		errOn:       2,
		errToReturn: errors.New("second read failed"),
	}

	_, err := V1.ReadFrame(c)
	if err == nil {
		t.Fatal("expected error on second read, got nil")
	}
}

func TestReadFrame_ZeroLength(t *testing.T) {
	header := []byte{0, 0, 0, 0}
	c := &readerConn{r: bytes.NewReader(header)}

	_, err := V1.ReadFrame(c)
	if err == nil {
		t.Fatal("expected error for zero-length frame, got nil")
	}
	if !errors.Is(err, ErrFrameLengthInvalid) {
		t.Fatalf("expected ErrFrameLengthInvalid, got %v", err)
	}
}

func TestDecodeMessage_Empty(t *testing.T) {
	msg, err := V1.DecodeMessage([]byte{})
	if msg != nil {
		t.Fatalf("expected nil message for empty input, got %+v", msg)
	}
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
	if !errors.Is(err, ErrPayloadTooShort) {
		t.Fatalf("expected ErrPayloadTooShort, got %v", err)
	}
}

func TestDecodeMessage_WrongMagic(t *testing.T) {
	// wrong magic byte
	data := []byte{magicNumber + 1, 1, 2, 3, 0, 0, 0, 0}
	msg, err := V1.DecodeMessage(data)
	if msg != nil {
		t.Fatalf("expected nil message for wrong magic, got %+v", msg)
	}
	if err == nil {
		t.Fatal("expected error for wrong magic, got nil")
	}
	if !errors.Is(err, ErrMagicNumberInvalid) {
		t.Fatalf("expected ErrMagicNumberInvalid, got %v", err)
	}
}

func TestDecodeMessage_TooShort(t *testing.T) {
	// correct magic but less than 8 bytes
	data := []byte{magicNumber, 1, 2, 3}
	msg, err := V1.DecodeMessage(data)
	if msg != nil {
		t.Fatalf("expected nil message for too-short data, got %+v", msg)
	}
	if err == nil {
		t.Fatal("expected error for too-short data, got nil")
	}
	if !errors.Is(err, ErrPayloadTooShort) {
		t.Fatalf("expected ErrPayloadTooShort, got %v", err)
	}
}

func TestDecodeMessage_InvalidProtocolVersion(t *testing.T) {
	// protocol version is not 0x01
	payload := []byte{0xAA}
	payloadLen := len(payload)
	header := []byte{
		magicNumber,
		0x02, // invalid protocol version
		0x00,
		byte(MessageTypeCommand),
		byte(payloadLen >> 24), byte(payloadLen >> 16), byte(payloadLen >> 8), byte(payloadLen),
	}
	data := append(header, payload...)

	msg, err := V1.DecodeMessage(data)
	if msg != nil {
		t.Fatalf("expected nil message for invalid protocol version, got %+v", msg)
	}
	if err == nil {
		t.Fatal("expected error for invalid protocol version, got nil")
	}
	if !errors.Is(err, ErrInvalidProtocolVersion) {
		t.Fatalf("expected ErrInvalidProtocolVersion, got %v", err)
	}
}

func TestDecodeMessage_UnknownMessageType(t *testing.T) {
	// type is outside allowed range (1-3)
	payload := []byte{0xAA}
	payloadLen := len(payload)
	header := []byte{
		magicNumber,
		byte(ProtocolVersion1),
		0x00,
		0x04, // unknown type (> MessageTypeError)
		byte(payloadLen >> 24), byte(payloadLen >> 16), byte(payloadLen >> 8), byte(payloadLen),
	}
	data := append(header, payload...)

	msg, err := V1.DecodeMessage(data)
	if msg != nil {
		t.Fatalf("expected nil message for unknown message type, got %+v", msg)
	}
	if err == nil {
		t.Fatal("expected error for unknown message type, got nil")
	}
	if !errors.Is(err, ErrUnknownMessageType) {
		t.Fatalf("expected ErrUnknownMessageType, got %v", err)
	}
}

func TestDecodeMessage_ZeroPayload(t *testing.T) {
	// header with zero payload length
	data := []byte{
		magicNumber, // magic
		0x01,        // protocol version
		0x00,        // flags
		0x01,        // type
		0, 0, 0, 0,  // payload length = 0
	}
	msg, err := V1.DecodeMessage(data)
	if msg != nil {
		t.Fatalf("expected nil message for zero payload, got %+v", msg)
	}
	if err == nil {
		t.Fatal("expected error for zero payload, got nil")
	}
	if !errors.Is(err, ErrEmptyPayload) {
		t.Fatalf("expected ErrEmptyPayload, got %v", err)
	}
}

func TestDecodeMessage_InsufficientPayload(t *testing.T) {
	// header says 5 bytes payload but only provide 3
	payloadLen := 5
	data := []byte{
		magicNumber,
		0x01,
		0x02,
		0x01,
		byte(payloadLen >> 24), byte(payloadLen >> 16), byte(payloadLen >> 8), byte(payloadLen),
		0xAA, 0xBB, 0xCC, // only 3 bytes instead of 5
	}
	msg, err := V1.DecodeMessage(data)
	if msg != nil {
		t.Fatalf("expected nil message for insufficient payload, got %+v", msg)
	}
	if err == nil {
		t.Fatal("expected error for insufficient payload, got nil")
	}
	if !errors.Is(err, ErrPayloadTooShort) {
		t.Fatalf("expected ErrPayloadTooShort, got %v", err)
	}
}

func TestDecodeMessage_Success(t *testing.T) {
	payload := []byte("hello")
	payloadLen := len(payload)
	header := []byte{
		magicNumber,
		0x01,
		0x00,
		0x02,
		byte(payloadLen >> 24), byte(payloadLen >> 16), byte(payloadLen >> 8), byte(payloadLen),
	}
	data := append(header, payload...)

	msg, err := V1.DecodeMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg == nil {
		t.Fatalf("expected non-nil message, got nil")
	}
	if msg.ProtocolVersion != 0x01 || msg.Flags != 0x00 || msg.Type != 0x02 {
		t.Fatalf("header mismatch: %+v", msg)
	}
	if string(msg.Payload) != string(payload) {
		t.Fatalf("payload mismatch: got %v want %v", msg.Payload, payload)
	}
}

func TestDecodeCommand_TooShort(t *testing.T) {
	// less than 10 bytes overall (now less than 14 with ReqID)
	data := make([]byte, 13)
	msg := &Message{
		ProtocolVersion: 0x01,
		Flags:           0x00,
		Type:            0x01,
		Payload:         data,
	}
	cmd, err := V1.DecodeCommand(msg)
	if cmd != nil {
		t.Fatalf("expected nil cmd for too-short data, got %+v", cmd)
	}
	if err == nil {
		t.Fatal("expected error for too-short data, got nil")
	}
	if !errors.Is(err, ErrPayloadTooShort) {
		t.Fatalf("expected ErrPayloadTooShort, got %v", err)
	}
}

func TestDecodeCommand_InsufficientDBName(t *testing.T) {
	var b []byte

	// ReqID
	tmp4 := make([]byte, 4)
	binary.BigEndian.PutUint32(tmp4, 0x01020304)
	b = append(b, tmp4...)

	// CMD ID
	tmp2 := make([]byte, 2)
	binary.BigEndian.PutUint16(tmp2, 0x1234)
	b = append(b, tmp2...)

	// dbNameLen = 5, provide only 3 bytes
	tmp2 = make([]byte, 2)
	binary.BigEndian.PutUint16(tmp2, 5)
	b = append(b, tmp2...)
	b = append(b, []byte("abc")...)

	msg := &Message{
		ProtocolVersion: 0x01,
		Flags:           0x00,
		Type:            0x01,
		Payload:         b,
	}
	cmd, err := V1.DecodeCommand(msg)
	if cmd != nil {
		t.Fatalf("expected nil cmd for insufficient DB name, got %+v", cmd)
	}
	if err == nil {
		t.Fatal("expected error for insufficient DB name, got nil")
	}
	if !errors.Is(err, ErrPayloadTooShort) {
		t.Fatalf("expected ErrPayloadTooShort, got %v", err)
	}
}

func TestDecodeCommand_InsufficientCollectionName(t *testing.T) {
	var b []byte

	// ReqID
	tmp4 := make([]byte, 4)
	binary.BigEndian.PutUint32(tmp4, 0x01020304)
	b = append(b, tmp4...)

	// CMD ID
	tmp2 := make([]byte, 2)
	binary.BigEndian.PutUint16(tmp2, 0x0102)
	b = append(b, tmp2...)

	// DB name
	tmp2 = make([]byte, 2)
	binary.BigEndian.PutUint16(tmp2, 2)
	b = append(b, tmp2...)
	b = append(b, []byte("db")...)

	// Collection name len = 4, only 2 bytes provided
	tmp2 = make([]byte, 2)
	binary.BigEndian.PutUint16(tmp2, 4)
	b = append(b, tmp2...)
	b = append(b, []byte("co")...) // insufficient

	msg := &Message{
		ProtocolVersion: 0x01,
		Flags:           0x00,
		Type:            0x01,
		Payload:         b,
	}
	cmd, err := V1.DecodeCommand(msg)
	if cmd != nil {
		t.Fatalf("expected nil cmd for insufficient collection name, got %+v", cmd)
	}
	if err == nil {
		t.Fatal("expected error for insufficient collection name, got nil")
	}
	if !errors.Is(err, ErrPayloadTooShort) {
		t.Fatalf("expected ErrPayloadTooShort, got %v", err)
	}
}

func TestDecodeCommand_InsufficientPayload(t *testing.T) {
	var b []byte

	// ReqID
	tmp4 := make([]byte, 4)
	binary.BigEndian.PutUint32(tmp4, 0x01020304)
	b = append(b, tmp4...)

	// CMD ID
	tmp2 := make([]byte, 2)
	binary.BigEndian.PutUint16(tmp2, 0x0A0B)
	b = append(b, tmp2...)

	// DB name
	tmp2 = make([]byte, 2)
	binary.BigEndian.PutUint16(tmp2, 2)
	b = append(b, tmp2...)
	b = append(b, []byte("db")...)

	// Collection name
	tmp2 = make([]byte, 2)
	binary.BigEndian.PutUint16(tmp2, 3)
	b = append(b, tmp2...)
	b = append(b, []byte("col")...)

	// Payload length = 10, provide only 3 bytes
	tmp4 = make([]byte, 4)
	binary.BigEndian.PutUint32(tmp4, 10)
	b = append(b, tmp4...)
	b = append(b, []byte{0xAA, 0xBB, 0xCC}...)

	msg := &Message{
		ProtocolVersion: 0x01,
		Flags:           0x00,
		Type:            0x01,
		Payload:         b,
	}
	cmd, err := V1.DecodeCommand(msg)
	if cmd != nil {
		t.Fatalf("expected nil cmd for insufficient payload, got %+v", cmd)
	}
	if err == nil {
		t.Fatal("expected error for insufficient payload, got nil")
	}
	if !errors.Is(err, ErrPayloadTooShort) {
		t.Fatalf("expected ErrPayloadTooShort, got %v", err)
	}
}

func TestDecodeCommand_Success(t *testing.T) {
	var b []byte

	reqID := uint32(0x01020304)

	// ReqID
	tmp4 := make([]byte, 4)
	binary.BigEndian.PutUint32(tmp4, reqID)
	b = append(b, tmp4...)

	// CMD ID
	tmp2 := make([]byte, 2)
	binary.BigEndian.PutUint16(tmp2, 0xDEAD)
	b = append(b, tmp2...)

	// DB name
	tmp2 = make([]byte, 2)
	binary.BigEndian.PutUint16(tmp2, 2)
	b = append(b, tmp2...)
	b = append(b, []byte("db")...)

	// Collection name
	tmp2 = make([]byte, 2)
	binary.BigEndian.PutUint16(tmp2, 3)
	b = append(b, tmp2...)
	b = append(b, []byte("col")...)

	// Payload
	payload := []byte("payload-data")
	tmp4 = make([]byte, 4)
	binary.BigEndian.PutUint32(tmp4, uint32(len(payload)))
	b = append(b, tmp4...)
	b = append(b, payload...)

	msg := &Message{
		ProtocolVersion: 0x01,
		Flags:           0x00,
		Type:            0x01,
		Payload:         b,
	}

	cmd, err := V1.DecodeCommand(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd == nil {
		t.Fatal("expected non-nil cmd, got nil")
	}
	if cmd.ReqID != reqID {
		t.Fatalf("unexpected ReqID: got 0x%X, want 0x%X", cmd.ReqID, reqID)
	}
	if cmd.ID != 0xDEAD {
		t.Fatalf("unexpected cmd ID: got 0x%X", cmd.ID)
	}
	if cmd.DatabaseName != "db" {
		t.Fatalf("unexpected DB name: got %q", cmd.DatabaseName)
	}
	if cmd.CollectionName != "col" {
		t.Fatalf("unexpected collection name: got %q", cmd.CollectionName)
	}
	if !bytes.Equal(cmd.Payload, payload) {
		t.Fatalf("payload mismatch: got %v want %v", cmd.Payload, payload)
	}
}
