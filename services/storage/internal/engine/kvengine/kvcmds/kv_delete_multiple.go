package kvcmds

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/OliverSchlueter/goutils/problems"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/codex"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/commonresponses"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/protocol"
	"github.com/fancyinnovations/fancyspaces/storage/internal/command"
	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
	"github.com/fancyinnovations/fancyspaces/storage/internal/engine/kvengine"
)

// handleDeleteMultiple implements the server side of the protocol.ServerCommandKVDeleteMultiple command over TCP.
// Payload format: encoded list of strings (keys), see codex.EncodeListInto
func (c *Commands) handleDeleteMultiple(_ *command.ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
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
	for i, kv := range keyVals {
		keys[i] = kv.AsString()
	}

	kve.DeleteMultiple(keys)

	return commonresponses.OK, nil
}

// deleteMultipleRequestHTTP is the request format for the delete multiple command over HTTP.
type deleteMultipleRequestHTTP struct {
	Keys []string `json:"keys"`
}

// handleDeleteMultipleHTTP implements the server side of the protocol.ServerCommandKVDeleteMultiple command over HTTP.
func (c *Commands) handleDeleteMultipleHTTP(w http.ResponseWriter, r *http.Request, _ *database.Database, _ *database.Collection, kve *kvengine.Engine) {
	var req deleteMultipleRequestHTTP
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}

	if len(req.Keys) == 0 {
		problems.ValidationError("keys", "At least one key is required").WriteToHTTP(w)
		return
	}

	kve.DeleteMultiple(req.Keys)

	w.WriteHeader(http.StatusOK)
}
