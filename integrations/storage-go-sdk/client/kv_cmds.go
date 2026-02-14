package client

import (
	"encoding/binary"

	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/codex"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/commonresponses"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/protocol"
)

// KVSet implements the client side of the protocol.ServerCommandKVSet command.
func (c *Client) KVSet(db, coll string, key string, value *codex.Value) error {
	data := codex.EncodeValue(value)

	totalLen := 2 + len(key) + len(data)
	payload := make([]byte, totalLen)

	// Key
	binary.BigEndian.PutUint16(payload[0:2], uint16(len(key)))
	copy(payload[2:2+len(key)], []byte(key))

	// Value
	copy(payload[2+len(key):], data)

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandKVSet,
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

// KVSetTTL implements the client side of the protocol.ServerCommandKVSetTTL command.
func (c *Client) KVSetTTL(db, coll string, key string, value *codex.Value, expiresAt uint64) error {
	data := codex.EncodeValue(value)

	totalLen := 2 + len(key) + len(data) + 8
	payload := make([]byte, totalLen)

	// Key
	binary.BigEndian.PutUint16(payload[0:2], uint16(len(key)))
	copy(payload[2:2+len(key)], []byte(key))

	// Value
	copy(payload[2+len(key):2+len(key)+len(data)], data)

	// TTL
	binary.BigEndian.PutUint64(payload[2+len(key)+len(data):], expiresAt)

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandKVSetTTL,
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

// KVDelete implements the client side of the protocol.ServerCommandKVDelete command.
func (c *Client) KVDelete(db, coll string, key string) error {
	totalLen := 2 + len(key)
	payload := make([]byte, totalLen)

	// Key
	binary.BigEndian.PutUint16(payload[0:2], uint16(len(key)))
	copy(payload[2:2+len(key)], []byte(key))

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandKVDelete,
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

// KVDeleteMultiple implements the client side of the protocol.ServerCommandKVDeleteMultiple command.
func (c *Client) KVDeleteMultiple(db, coll string, keys []string) error {
	keyVals := codex.NewStringListValue(keys)

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandKVDeleteMultiple,
		DatabaseName:   db,
		CollectionName: coll,
		Payload:        codex.EncodeValue(keyVals),
	})
	if err != nil {
		return err
	}

	if resp.Code != protocol.StatusOK {
		return ErrUnexpectedStatusCode
	}

	return nil
}

// KVDeleteAll implements the client side of the protocol.ServerCommandKVDeleteAll command.
func (c *Client) KVDeleteAll(db, coll string) error {
	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandKVDeleteAll,
		DatabaseName:   db,
		CollectionName: coll,
		Payload:        *commonresponses.EmptyPayload,
	})
	if err != nil {
		return err
	}

	if resp.Code != protocol.StatusOK {
		return ErrUnexpectedStatusCode
	}

	return nil
}

// KVExists implements the client side of the protocol.ServerCommandKVExists command.
func (c *Client) KVExists(db, coll string, key string) (bool, error) {
	totalLen := 2 + len(key)
	payload := make([]byte, totalLen)

	// Key
	binary.BigEndian.PutUint16(payload[0:2], uint16(len(key)))
	copy(payload[2:2+len(key)], []byte(key))

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandKVExists,
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

// KVGet implements the client side of the protocol.ServerCommandKVGet command.
func (c *Client) KVGet(db, coll string, key string) (*codex.Value, error) {
	totalLen := 2 + len(key)
	payload := make([]byte, totalLen)

	// Key
	binary.BigEndian.PutUint16(payload[0:2], uint16(len(key)))
	copy(payload[2:2+len(key)], []byte(key))

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandKVGet,
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

// KVGetMultiple implements the client side of the protocol.ServerCommandKVGetMultiple command.
func (c *Client) KVGetMultiple(db, coll string, keys []string) (map[string]*codex.Value, error) {
	keyVals := codex.NewStringListValue(keys)

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandKVGetMultiple,
		DatabaseName:   db,
		CollectionName: coll,
		Payload:        codex.EncodeValue(keyVals),
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

	if val.Type != codex.TypeMap {
		return nil, ErrUnexpectedDataType
	}

	return val.AsMap(), nil
}

// KVGetAll implements the client side of the protocol.ServerCommandKVGetAll command.
func (c *Client) KVGetAll(db, coll string) (map[string]*codex.Value, error) {
	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandKVGetAll,
		DatabaseName:   db,
		CollectionName: coll,
		Payload:        *commonresponses.EmptyPayload,
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

	if val.Type != codex.TypeMap {
		return nil, ErrUnexpectedDataType
	}

	return val.AsMap(), nil
}

// KVGetTTL implements the client side of the protocol.ServerCommandKVGetTTL command.
func (c *Client) KVGetTTL(db, coll string, key string) (int64, error) {
	totalLen := 2 + len(key)
	payload := make([]byte, totalLen)

	// Key
	binary.BigEndian.PutUint16(payload[0:2], uint16(len(key)))
	copy(payload[2:2+len(key)], []byte(key))

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandKVGetTTL,
		DatabaseName:   db,
		CollectionName: coll,
		Payload:        payload,
	})
	if err != nil {
		return 0, err
	}

	if resp.Code != protocol.StatusOK {
		if resp.Code == protocol.StatusNotFound {
			return 0, ErrKeyNotFound
		}

		return 0, ErrUnexpectedStatusCode
	}

	if len(resp.Payload) != 8 {
		return 0, ErrUnexpectedDataType
	}

	ttl := binary.BigEndian.Uint64(resp.Payload)
	return int64(ttl), nil
}

// KVGetMultipleTTL implements the client side of the protocol.ServerCommandKVGetMultipleTTL command.
func (c *Client) KVGetMultipleTTL(db, coll string, keys []string) (map[string]int64, error) {
	keyVals := codex.NewStringListValue(keys)

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandKVGetMultipleTTL,
		DatabaseName:   db,
		CollectionName: coll,
		Payload:        codex.EncodeValue(keyVals),
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

	if val.Type != codex.TypeMap {
		return nil, ErrUnexpectedDataType
	}

	result := make(map[string]int64)
	for key, ttlVal := range val.AsMap() {
		if ttlVal.Type != codex.TypeInt64 {
			return nil, ErrUnexpectedDataType
		}
		result[key] = ttlVal.AsInt64()
	}

	return result, nil
}

func (c *Client) KVGetAllTTL(db, coll string) (map[string]int64, error) {
	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandKVGetAllTTL,
		DatabaseName:   db,
		CollectionName: coll,
		Payload:        *commonresponses.EmptyPayload,
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

	if val.Type != codex.TypeMap {
		return nil, ErrUnexpectedDataType
	}

	result := make(map[string]int64)
	for key, ttlVal := range val.AsMap() {
		if ttlVal.Type != codex.TypeInt64 {
			return nil, ErrUnexpectedDataType
		}
		result[key] = ttlVal.AsInt64()
	}

	return result, nil
}

// KVKeys implements the client side of the protocol.ServerCommandKVKeys command.
func (c *Client) KVKeys(db, coll string) ([]string, error) {
	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandKVKeys,
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

// KVCount implements the client side of the protocol.ServerCommandKVCount command.
func (c *Client) KVCount(db, coll string) (uint32, error) {
	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandKVCount,
		DatabaseName:   db,
		CollectionName: coll,
		Payload:        *commonresponses.EmptyPayload,
	})
	if err != nil {
		return 0, err
	}

	if resp.Code != protocol.StatusOK {
		return 0, ErrUnexpectedStatusCode
	}

	count, err := codex.DecodeUint32(resp.Payload)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// KVSize implements the client side of the protocol.ServerCommandKVSize command.
func (c *Client) KVSize(db, coll string) (uint64, error) {
	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandKVSize,
		DatabaseName:   db,
		CollectionName: coll,
		Payload:        *commonresponses.EmptyPayload,
	})
	if err != nil {
		return 0, err
	}

	if resp.Code != protocol.StatusOK {
		return 0, ErrUnexpectedStatusCode
	}

	size, err := codex.DecodeUint64(resp.Payload)
	if err != nil {
		return 0, err
	}

	return size, nil
}
