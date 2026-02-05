package client

import (
	"encoding/binary"
	"time"

	"github.com/fancyinnovations/fancyspaces/storage/pkg/codex"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/protocol"
)

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

// KVSet sets a key-value pair in the specified collection.
// The value can be of any type that codex (codex.ValueType) supports.
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

// KVSetTTL sets a key-value pair in the specified collection with a time-to-live (TTL).
// The value can be of any type that codex (codex.ValueType) supports.
// After the TTL expires, the key-value pair will be automatically deleted from the collection.
func (c *Client) KVSetTTL(db, coll string, key string, value any, ttlMillis uint64) error {
	val, err := codex.NewValue(value)
	if err != nil {
		return err
	}
	data := codex.EncodeValue(val)

	totalLen := 2 + len(key) + len(data) + 8
	payload := make([]byte, totalLen)

	// Key
	binary.BigEndian.PutUint16(payload[0:2], uint16(len(key)))
	copy(payload[2:2+len(key)], []byte(key))

	// Value
	copy(payload[2+len(key):2+len(key)+len(data)], data)

	// TTL
	ttlNanos := ttlMillis * 1_000_000 // Convert milliseconds to nanoseconds
	expiresAt := time.Now().UnixNano() + int64(ttlNanos)
	binary.BigEndian.PutUint64(payload[2+len(key)+len(data):], uint64(expiresAt))

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.CommandKVSetTTL,
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

// KVDelete deletes the key-value pair associated with the specified key from the collection.
func (c *Client) KVDelete(db, coll string, key string) error {
	totalLen := 2 + len(key)
	payload := make([]byte, totalLen)

	// Key
	binary.BigEndian.PutUint16(payload[0:2], uint16(len(key)))
	copy(payload[2:2+len(key)], []byte(key))

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.CommandKVDelete,
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

// KVExists checks if a key exists in the specified collection.
// It returns true if the key exists, false if it does not exist, and an error if there was an issue checking.
func (c *Client) KVExists(db, coll string, key string) (bool, error) {
	totalLen := 2 + len(key)
	payload := make([]byte, totalLen)

	// Key
	binary.BigEndian.PutUint16(payload[0:2], uint16(len(key)))
	copy(payload[2:2+len(key)], []byte(key))

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.CommandKVExists,
		DatabaseName:   db,
		CollectionName: coll,
		Payload:        payload,
	})
	if err != nil {
		return false, err
	}

	switch resp.Code {
	case protocol.StatusOK:
		return true, nil
	case protocol.StatusNotFound:
		return false, nil
	default:
		return false, ErrUnexpectedStatusCode
	}
}

// KVKeys retrieves all keys in the specified collection.
// It returns a slice of strings representing the keys, or an error if there was an issue retrieving the keys.
func (c *Client) KVKeys(db, coll string) ([]string, error) {
	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.CommandKVKeys,
		DatabaseName:   db,
		CollectionName: coll,
		Payload:        []byte{},
	})
	if err != nil {
		return nil, err
	}

	if resp.Code != protocol.StatusOK {
		return nil, ErrUnexpectedStatusCode
	}

	val, err := codex.DecodeValue(resp.Payload)
	if err != nil {
		return nil, err
	}

	if val.Type != codex.TypeList {
		return nil, ErrUnexpectedDataType
	}

	keys := make([]string, len(val.AsList()))
	for i, item := range val.AsList() {
		keys[i] = item.AsString()
	}

	return keys, nil
}
