package objectcmds

import (
	"encoding/binary"
	"errors"
	"log/slog"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/storage/internal/auth"
	"github.com/fancyinnovations/fancyspaces/storage/internal/command"
	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/commonresponses"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/protocol"
)

// handleSize returns the total size of objects in an object engine in bytes.
// Payload format: empty
// Response payload format: | Size (8 bytes) |
func (c *Commands) handleSize(ctx *command.ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
	u := auth.UserFromContext(ctx.Ctx)
	if u == nil || !u.Verified || !u.IsActive {
		return commonresponses.Unauthorized, nil
	}

	db, err := c.dbStore.GetDatabase(ctx.Ctx, cmd.DatabaseName)
	if err != nil {
		if errors.Is(err, database.ErrDatabaseNotFound) {
			return commonresponses.DatabaseNotFound, nil
		}

		slog.Error("Failed to get database",
			slog.String("database", cmd.DatabaseName),
			sloki.WrapError(err),
		)
		return commonresponses.InternalServerError, nil
	}

	if !u.IsAdmin() && !db.HasPermission(u.ID, database.PermissionLevelReadOnly) {
		return commonresponses.Forbidden, nil
	}

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

	if e.Type != database.EngineObject {
		return commonresponses.CommandNotAllowed, nil
	}
	obje := e.AsObjectEngine()

	size := obje.Size()

	respPayload := make([]byte, 8)
	binary.BigEndian.PutUint64(respPayload, size)

	return &protocol.Response{
		Code:    protocol.StatusOK,
		Payload: respPayload,
	}, nil
}
