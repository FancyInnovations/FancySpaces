package client

import (
	"encoding/binary"

	"github.com/fancyinnovations/fancyspaces/storage/pkg/codex"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/protocol"
)

// KVSet sets a key-value pair in the specified collection. The value can be of any type that codex (codex.ValueType) supports.
func (c *Client) KVSet(db, coll string, key string, value any) error {
	val, err := codex.NewValue(value)
	if err != nil {
		return err
	}
	data := codex.EncodeValue(val)

	totalLen := 2 + len(key) + len(data)
	payload := make([]byte, totalLen)

	// Key
	binary.BigEndian.PutUint16(payload[0:2], uint16(len(key)))
	copy(payload[2:2+len(key)], []byte(key))

	// Value
	copy(payload[2+len(key):], data)

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.CommandKVSet,
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

// KVGet retrieves the value associated with the specified key from the collection.
// It returns a codex.Value, which can be of any type supported by codex.
func (c *Client) KVGet(db, coll string, key string) (*codex.Value, error) {
	totalLen := 2 + len(key)
	payload := make([]byte, totalLen)

	// Key
	binary.BigEndian.PutUint16(payload[0:2], uint16(len(key)))
	copy(payload[2:2+len(key)], []byte(key))

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.CommandKVGet,
		DatabaseName:   db,
		CollectionName: coll,
		Payload:        payload,
	})
	if err != nil {
		return nil, err
	}

	if resp.Code != protocol.StatusOK {
		if resp.Code == protocol.StatusNotFound {
			return nil, ErrKeyNotFound
		}

		return nil, ErrUnexpectedStatusCode
	}

	val, err := codex.DecodeValue(resp.Payload)
	if err != nil {
		return nil, err
	}

	return val, nil
}
