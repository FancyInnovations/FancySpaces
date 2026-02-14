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

// handleCountTCP implements the server side of the protocol.ServerCommandKVCount command.
// Payload format: empty
func (c *Commands) handleCountTCP(_ *command.ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
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

	count := kve.Count()

	return &protocol.Response{
		Code:    protocol.StatusOK,
		Payload: codex.EncodeUint32(count),
	}, nil
}

// countResponseHTTP is the response format for the count command over HTTP.
type countResponseHTTP struct {
	Count uint32 `json:"count"`
}

// handleCountHTTP implements the server side of the protocol.ServerCommandKVCount command.
func (c *Commands) handleCountHTTP(w http.ResponseWriter, _ *http.Request, _ *database.Database, _ *database.Collection, kve *kvengine.Engine) {
	count := kve.Count()

	resp := countResponseHTTP{
		Count: count,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=10") // 10s
	json.NewEncoder(w).Encode(resp)
}
