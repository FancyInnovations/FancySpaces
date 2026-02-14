package kvcmds

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/OliverSchlueter/goutils/problems"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/commonresponses"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/protocol"
	"github.com/fancyinnovations/fancyspaces/storage/internal/command"
	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
	"github.com/fancyinnovations/fancyspaces/storage/internal/engine/kvengine"
)

// handleExists implements the server side of the protocol.ServerCommandKVExists command over TCP.
// Payload format: | Key Length (2 bytes) | Key (variable) |
func (c *Commands) handleExists(_ *command.ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
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

	key := string(data[2 : 2+keyLen])

	exists := kve.Exists(key)

	if !exists {
		return &protocol.Response{
			Code:    protocol.StatusNotFound,
			Payload: *commonresponses.EmptyPayload,
		}, nil
	}

	return &protocol.Response{
		Code:    protocol.StatusOK,
		Payload: *commonresponses.EmptyPayload,
	}, nil
}

// existsRequestHTTP is the request format for the exists command over HTTP.
type existsRequestHTTP struct {
	Key string `json:"key"`
}

// existsResponseHTTP is the response format for the exists command over HTTP.
type existsResponseHTTP struct {
	Exists bool `json:"exists"`
}

// handleExistsHTTP implements the server side of the protocol.ServerCommandKVExists command over HTTP.
func (c *Commands) handleExistsHTTP(w http.ResponseWriter, r *http.Request, _ *database.Database, _ *database.Collection, kve *kvengine.Engine) {
	var req existsRequestHTTP
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode exists request", sloki.WrapError(err))
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}

	exists := kve.Exists(req.Key)

	resp := existsResponseHTTP{
		Exists: exists,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
