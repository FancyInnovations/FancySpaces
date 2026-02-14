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

// handleGetMultiple implements the server side of the protocol.ServerCommandKVGetMultiple command over TCP.
// Payload format: encoded list of strings (keys), see codex.EncodeListInto
func (c *Commands) handleGetMultiple(_ *command.ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
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

// getMultipleRequestHTTP is the request format for the get multiple command over HTTP.
type getMultipleRequestHTTP struct {
	Keys []string `json:"keys"`
}

// getMultipleResponseHTTP is the response format for the get multiple command over HTTP.
type getMultipleResponseHTTP struct {
	Values map[string]any `json:"values"`
}

// handleGetMultipleHTTP implements the server side of the protocol.ServerCommandKVGetMultiple command over HTTP.
func (c *Commands) handleGetMultipleHTTP(w http.ResponseWriter, r *http.Request, _ *database.Database, _ *database.Collection, kve *kvengine.Engine) {
	var req getMultipleRequestHTTP
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}

	if len(req.Keys) == 0 {
		problems.ValidationError("keys", "At least one key is required").WriteToHTTP(w)
		return
	}

	vals := kve.GetMultiple(req.Keys)

	values := make(map[string]any, len(vals))
	for k, v := range vals {
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

	resp := getMultipleResponseHTTP{
		Values: values,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
