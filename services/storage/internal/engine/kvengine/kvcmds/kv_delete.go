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

// handleDelete implements the server side of the protocol.ServerCommandKVDelete command over TCP.
// Payload format: | Key Length (2 bytes) | Key (variable) |
func (c *Commands) handleDelete(_ *command.ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
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

	kve.Delete(key)

	return &protocol.Response{
		Code:    protocol.StatusOK,
		Payload: *commonresponses.EmptyPayload,
	}, nil
}

// deleteRequestHTTP is the request format for the delete command over HTTP.
type deleteRequestHTTP struct {
	Key string `json:"key"`
}

// handleDeleteHTTP implements the server side of the protocol.ServerCommandKVDelete command over HTTP.
func (c *Commands) handleDeleteHTTP(w http.ResponseWriter, r *http.Request, _ *database.Database, _ *database.Collection, kve *kvengine.Engine) {
	var req deleteRequestHTTP
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode delete request", sloki.WrapError(err))
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}

	kve.Delete(req.Key)

	w.WriteHeader(http.StatusOK)
}
