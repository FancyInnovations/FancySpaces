package client

import (
	"encoding/binary"
	"fmt"

	"github.com/fancyinnovations/fancyspaces/storage/pkg/codex"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/protocol"
)

// ObjPut implements the client side of the protocol.ServerCommandObjectPut command.
func (c *Client) ObjPut(db, coll string, key string, data []byte) error {
	encodedData := codex.EncodeBinary(data)

	totalLen := 2 + len(key) + len(encodedData)
	payload := make([]byte, totalLen)

	// Key
	binary.BigEndian.PutUint16(payload[0:2], uint16(len(key)))
	copy(payload[2:2+len(key)], []byte(key))

	// Value
	copy(payload[2+len(key):], encodedData)

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandObjectPut,
		DatabaseName:   db,
		CollectionName: coll,
		Payload:        payload,
	})
	if err != nil {
		return err
	}

	if resp.Code != protocol.StatusOK {
		return ErrUnexpectedStatusCode
	}

	return nil
}

// ObjGet implements the client side of the protocol.ServerCommandObjectGet command.
func (c *Client) ObjGet(db, coll string, key string) ([]byte, error) {
	payload := make([]byte, 2+len(key))
	binary.BigEndian.PutUint16(payload[0:2], uint16(len(key)))
	copy(payload[2:], []byte(key))

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandObjectGet,
		DatabaseName:   db,
		CollectionName: coll,
		Payload:        payload,
	})
	if err != nil {
		return nil, err
	}

	if resp.Code == protocol.StatusNotFound {
		return nil, ErrKeyNotFound
	}

	if resp.Code != protocol.StatusOK {
		return nil, ErrUnexpectedStatusCode
	}

	return codex.DecodeBinary(resp.Payload)
}

// ObjGetMetadata implements the client side of the protocol.ServerCommandObjectGetMetadata command.
func (c *Client) ObjGetMetadata(db, coll string, key string) (*ObjectMetadata, error) {
	payload := make([]byte, 2+len(key))
	binary.BigEndian.PutUint16(payload[0:2], uint16(len(key)))
	copy(payload[2:], []byte(key))

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandObjectGetMetadata,
		DatabaseName:   db,
		CollectionName: coll,
		Payload:        payload,
	})
	if err != nil {
		return nil, err
	}

	if resp.Code == protocol.StatusNotFound {
		return nil, ErrKeyNotFound
	}

	if resp.Code != protocol.StatusOK {
		return nil, ErrUnexpectedStatusCode
	}

	data := resp.Payload
	if len(data) != 8+4+8+8 {
		fmt.Printf("invalid metadata payload length: expected 12, got %d\n", len(data))
		return nil, ErrInvalidPayloadLength
	}

	return &ObjectMetadata{
		Key:        key,
		Size:       uint32(binary.BigEndian.Uint64(data[0:8])),
		Checksum:   binary.BigEndian.Uint32(data[8:12]),
		CreatedAt:  int64(binary.BigEndian.Uint64(data[12:20])),
		ModifiedAt: int64(binary.BigEndian.Uint64(data[20:28])),
	}, nil
}

// ObjDelete implements the client side of the protocol.ServerCommandObjectDelete command.
func (c *Client) ObjDelete(db, coll string, key string) error {
	payload := make([]byte, 2+len(key))
	binary.BigEndian.PutUint16(payload[0:2], uint16(len(key)))
	copy(payload[2:], []byte(key))

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandObjectDelete,
		DatabaseName:   db,
		CollectionName: coll,
		Payload:        payload,
	})
	if err != nil {
		return err
	}

	if resp.Code == protocol.StatusNotFound {
		return ErrKeyNotFound
	}

	if resp.Code != protocol.StatusOK {
		return ErrUnexpectedStatusCode
	}

	return nil
}
