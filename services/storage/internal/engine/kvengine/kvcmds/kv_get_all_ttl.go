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

// handleGetAllTTL processes a command to get the TTL of all keys in a key-value collection.
// Payload format: empty
func (c *Commands) handleGetAllTTL(ctx *command.ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
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

	ttls := kve.GetAllTTL()

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
