package kvcmds

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/commonresponses"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/protocol"
	"github.com/fancyinnovations/fancyspaces/storage/internal/command"
	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
	"github.com/fancyinnovations/fancyspaces/storage/internal/engine/kvengine"
)

// handleDeleteAll implements the server side of the protocol.ServerCommandKVDeleteAll command over TCP.
// Payload format: empty
func (c *Commands) handleDeleteAll(_ *command.ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
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

	kve.DeleteAll()

	return commonresponses.OK, nil
}

// handleDeleteAllHTTP implements the server side of the protocol.ServerCommandKVDeleteAll command over HTTP.
func (c *Commands) handleDeleteAllHTTP(w http.ResponseWriter, _ *http.Request, _ *database.Database, _ *database.Collection, kve *kvengine.Engine) {
	kve.DeleteAll()

	w.WriteHeader(http.StatusOK)
}
