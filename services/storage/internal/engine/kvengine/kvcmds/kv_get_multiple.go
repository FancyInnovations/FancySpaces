package kvcmds

import (
	"errors"
	"log/slog"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/codex"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/commonresponses"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/protocol"
	"github.com/fancyinnovations/fancyspaces/storage/internal/command"
	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
)

// handleGetMultiple handles the protocol.ServerCommandKVGetMultiple command, which retrieves the values for a given keys from a key-value collection.
// Payload format: encoded list of strings (keys), see codex.EncodeListInto
func (c *Commands) handleGetMultiple(ctx *command.ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
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
				Payload: []byte("invalid payload format: all items must be strings"),
			}, nil
		}

		keys[i] = keyVal.AsString()
	}

	values := kve.GetMultiple(keys)

	return &protocol.Response{
		Code:    protocol.StatusOK,
		Payload: codex.EncodeMap(values),
	}, nil
}
