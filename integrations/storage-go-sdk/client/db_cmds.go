package client

import (
	"encoding/binary"

	"github.com/fancyinnovations/fancyspaces/storage/pkg/codex"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/protocol"
)

func (c *Client) DBDatabaseGet(dbName string) (*DatabaseDatabase, error) {
	totalLen := 2 + len(dbName)
	payload := make([]byte, totalLen)

	// Database name
	binary.BigEndian.PutUint16(payload[0:2], uint16(len(dbName)))
	copy(payload[2:2+len(dbName)], []byte(dbName))

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandDBDatabaseGet,
		DatabaseName:   "",
		CollectionName: "",
		Payload:        payload,
	})
	if err != nil {
		return nil, err
	}

	if resp.Code != protocol.StatusOK {
		if resp.Code == protocol.StatusNotFound {
			return nil, ErrDatabaseNotFound
		}

		return nil, ErrUnexpectedStatusCode
	}

	data, err := codex.DecodeValue(resp.Payload)
	if err != nil {
		return nil, err
	}

	var db DatabaseDatabase
	if err := codex.Unmarshal(data.AsBinary(), &db); err != nil {
		return nil, err
	}

	return &db, nil
}

func (c *Client) DBCollectionGet(dbName, collName string) (*DatabaseCollection, error) {
	totalLen := 2 + len(dbName) + 2 + len(collName)
	payload := make([]byte, totalLen)

	// Database name
	binary.BigEndian.PutUint16(payload[0:2], uint16(len(dbName)))
	copy(payload[2:2+len(dbName)], []byte(dbName))

	// Collection name
	offset := 2 + len(dbName)
	binary.BigEndian.PutUint16(payload[offset:offset+2], uint16(len(collName)))
	copy(payload[offset+2:offset+2+len(collName)], []byte(collName))

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandDBCollectionGet,
		DatabaseName:   "",
		CollectionName: "",
		Payload:        payload,
	})
	if err != nil {
		return nil, err
	}

	if resp.Code != protocol.StatusOK {
		if resp.Code == protocol.StatusNotFound {
			return nil, ErrCollectionNotFound
		}

		return nil, ErrUnexpectedStatusCode
	}

	data, err := codex.DecodeValue(resp.Payload)
	if err != nil {
		return nil, err
	}

	var coll DatabaseCollection
	if err := codex.Unmarshal(data.AsBinary(), &coll); err != nil {
		return nil, err
	}

	return &coll, nil
}
