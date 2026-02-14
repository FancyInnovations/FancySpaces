package kvcmds

import (
	"encoding/binary"
	"errors"
	"log/slog"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/codex"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/commonresponses"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/protocol"
	"github.com/fancyinnovations/fancyspaces/storage/internal/command"
	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
)

// handleGetMultiple handles the protocol.ServerCommandKVGetMultiple command, which retrieves the TTL for a given keys from a key-value collection.
// Payload format: encoded list of strings (keys), see codex.EncodeListInto
func (c *Commands) handleGetMultipleTTL(ctx *command.ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
	e, err := c.engineService.GetEngine(cmd.DatabaseName, cmd.CollectionName)
	if err != nil {
		if errors.Is(err, database.ErrCollectionNotFound) {
			return commonresponses.CollectionNotFound, nil
		}

		slog.Error("Failed to get engine",
			slog.String("database", cmd.DatabaseName),
			slog.String("collection", cmd.CollectionName),
			sloki.WrapError(err),
		)
		return commonresponses.InternalServerError, nil
	}

	if e.Type != database.EngineKeyValue {
		return commonresponses.CommandNotAllowed, nil
	}
	kve := e.AsKeyValueEngine()

	data := cmd.Payload
	if len(data) < 2 {
		return &protocol.Response{
			Code:    protocol.StatusBadRequest,
			Payload: []byte("invalid payload length"),
		}, nil
	}

	keyLen := int(binary.BigEndian.Uint16(data[0:2]))
	if len(data) < 2+keyLen {
		return &protocol.Response{
			Code:    protocol.StatusBadRequest,
			Payload: []byte("invalid payload length for key"),
		}, nil
	}

	keyVals, err := codex.DecodeList(data)
	if err != nil {
		return &protocol.Response{
			Code:    protocol.StatusBadRequest,
			Payload: []byte("invalid payload format"),
		}, nil
	}

	keys := make([]string, len(keyVals))
	for i, keyVal := range keyVals {
		if keyVal.Type != codex.TypeString {
			return &protocol.Response{
				Code:    protocol.StatusBadRequest,
				Payload: []byte("invalid payload format: all keys must be strings"),
			}, nil
		}

		keys[i] = keyVal.AsString()
	}

	ttls := kve.GetMultipleTTL(keys)

	respMap := make(map[string]*codex.Value, len(ttls))
	for key, ttl := range ttls {
		ttlVal, err := codex.NewValue(ttl)
		if err != nil {
			slog.Error("Failed to encode TTL value",
				slog.String("key", key),
				sloki.WrapError(err),
			)
			return &protocol.Response{
				Code:    protocol.StatusInternalServerError,
				Payload: []byte("failed to encode TTL value"),
			}, nil
		}

		respMap[key] = ttlVal
	}

	return &protocol.Response{
		Code:    protocol.StatusOK,
		Payload: codex.EncodeMap(respMap),
	}, nil
}
