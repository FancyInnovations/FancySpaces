package client

import (
	"encoding/binary"
	"fmt"

	"github.com/fancyinnovations/fancyspaces/storage/pkg/codex"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/protocol"
)

// ObjPut stores an object in the specified database and collection with the given key and binary data.
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

// ObjGet retrieves the binary data associated with the given key from the specified database and collection.
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

// ObjGetMetadata retrieves the metadata (size and checksum) for the object associated with the given key from the specified database and collection.
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
	if len(data) != 8+4 {
		fmt.Printf("invalid metadata payload length: expected 12, got %d\n", len(data))
		return nil, ErrInvalidPayloadLength
	}

	return &ObjectMetadata{
		Size:     int64(binary.BigEndian.Uint64(data[0:8])),
		Checksum: binary.BigEndian.Uint32(data[8:12]),
	}, nil
}

// ObjDelete removes the object associated with the given key from the specified database and collection.
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
