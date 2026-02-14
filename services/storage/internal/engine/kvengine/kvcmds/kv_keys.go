package kvcmds

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/codex"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/commonresponses"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/protocol"
	"github.com/fancyinnovations/fancyspaces/storage/internal/command"
	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
	"github.com/fancyinnovations/fancyspaces/storage/internal/engine/kvengine"
)

// handleKeys implements the server side of the protocol.ServerCommandKVKeys command over TCP.
// Payload format: empty
func (c *Commands) handleKeys(_ *command.ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
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

	keys := kve.Keys()

	// Encode keys as a codex array
	val := codex.NewStringListValue(keys)

	return &protocol.Response{
		Code:    protocol.StatusOK,
		Payload: codex.EncodeValue(val),
	}, nil
}

// keysResponseHTTP is the response format for the keys command over HTTP.
type keysResponseHTTP struct {
	Keys []string `json:"keys"`
}

// handleKeysHTTP implements the server side of the protocol.ServerCommandKVKeys command over HTTP.
func (c *Commands) handleKeysHTTP(w http.ResponseWriter, _ *http.Request, _ *database.Database, _ *database.Collection, kve *kvengine.Engine) {
	keys := kve.Keys()

	resp := keysResponseHTTP{
		Keys: keys,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
