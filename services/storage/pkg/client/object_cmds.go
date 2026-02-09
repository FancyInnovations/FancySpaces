package client

import (
	"encoding/binary"

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
