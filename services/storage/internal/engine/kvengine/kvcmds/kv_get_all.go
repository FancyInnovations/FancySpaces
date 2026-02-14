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

// handleGetAll implements the server side of the protocol.ServerCommandKVGetAll command over TCP.
// Payload format: empty
func (c *Commands) handleGetAll(_ *command.ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
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

	all := kve.GetAll()

	return &protocol.Response{
		Code:    protocol.StatusOK,
		Payload: codex.EncodeMap(all),
	}, nil
}

// getAllResponseHTTP is the response format for the get all command over HTTP.
type getAllResponseHTTP struct {
	Values map[string]any `json:"values"`
}

// handleGetAllHTTP implements the server side of the protocol.ServerCommandKVGetAll command over HTTP.
func (c *Commands) handleGetAllHTTP(w http.ResponseWriter, _ *http.Request, _ *database.Database, _ *database.Collection, kve *kvengine.Engine) {
	all := kve.GetAll()

	values := make(map[string]any, len(all))
	for k, v := range all {
		data, err := v.ToAny()
		if err != nil {
			slog.Error("Failed to convert value to any",
				slog.String("key", k),
				sloki.WrapError(err),
			)
			continue
		}
		values[k] = data
	}

	resp := getAllResponseHTTP{
		Values: values,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
