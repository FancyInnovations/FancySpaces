package objectcmds

import (
	"encoding/binary"
	"errors"
	"log/slog"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/storage/internal/auth"
	"github.com/fancyinnovations/fancyspaces/storage/internal/command"
	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
	"github.com/fancyinnovations/fancyspaces/storage/internal/engine/objectengine"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/codex"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/commonresponses"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/protocol"
)

// handleGet processes a get command for an object engine. The payload is expected to be in the format:
// Payload format: | Key Length (2 bytes) | Key (variable) |
// Response payload will be the binary data associated with the key, encoded using codex.TypeBinary.
// If the key is not found, a protocol.StatusNotFound response will be returned with an empty payload.
func (c *Commands) handleGet(ctx *command.ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
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

	data := cmd.Payload
	if len(data) < 2 {
		return &protocol.Response{
			Code:    protocol.StatusBadRequest,
			Payload: []byte("invalid payload length"),
		}, nil
	}

	keyLen := binary.BigEndian.Uint16(data[0:2])
	if len(data) < int(2+keyLen) {
		return &protocol.Response{
			Code:    protocol.StatusBadRequest,
			Payload: []byte("invalid payload length for key"),
		}, nil
	}

	key := string(data[2 : 2+keyLen])

	binData, err := obje.Get(key)
	if err != nil {
		if errors.Is(err, objectengine.ErrKeyNotFound) {
			return &protocol.Response{
				Code:    protocol.StatusNotFound,
				Payload: *commonresponses.EmptyPayload,
			}, nil
		}

		slog.Error("Failed to get object",
			slog.String("database", cmd.DatabaseName),
			slog.String("collection", cmd.CollectionName),
			slog.String("key", key),
			sloki.WrapError(err),
		)
		return commonresponses.InternalServerError, nil
	}

	return &protocol.Response{
		Code:    protocol.StatusOK,
		Payload: codex.EncodeBinary(binData),
	}, nil
}
