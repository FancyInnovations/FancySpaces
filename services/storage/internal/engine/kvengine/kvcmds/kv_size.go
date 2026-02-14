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

// handleSize implements the server side of the protocol.ServerCommandKVSize command over TCP.
// Payload format: empty
func (c *Commands) handleSize(_ *command.ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
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

	size := kve.Size()

	return &protocol.Response{
		Code:    protocol.StatusOK,
		Payload: codex.EncodeUint64(size),
	}, nil
}

// sizeResponseHTTP is the response format for the size command over HTTP.
type sizeResponseHTTP struct {
	Size uint64 `json:"size"`
}

// handleSizeHTTP implements the server side of the protocol.ServerCommandKVSize command over HTTP.
func (c *Commands) handleSizeHTTP(w http.ResponseWriter, _ *http.Request, _ *database.Database, _ *database.Collection, kve *kvengine.Engine) {
	size := kve.Size()

	resp := sizeResponseHTTP{
		Size: size,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60") // 60s
	json.NewEncoder(w).Encode(resp)
}
