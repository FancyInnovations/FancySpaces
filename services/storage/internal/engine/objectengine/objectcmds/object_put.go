package objectcmds

import (
	"encoding/binary"
	"errors"
	"log/slog"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/storage/internal/auth"
	"github.com/fancyinnovations/fancyspaces/storage/internal/command"
	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/codex"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/commonresponses"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/protocol"
)

// handlePut processes a put command for an object engine. The payload is expected to be in the format:
// Payload format: | Key Length (2 bytes) | Key (variable) | Data (codex.TypeBinary) |
func (c *Commands) handlePut(ctx *command.ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
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

	if !u.IsAdmin() && !db.HasPermission(u.ID, database.PermissionLevelReadWrite) {
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

	binVal, err := codex.DecodeBinary(data[2+keyLen:])
	if err != nil {
		return &protocol.Response{
			Code:    protocol.StatusBadRequest,
			Payload: []byte("invalid binary value encoding: " + err.Error()),
		}, nil
	}

	if err := obje.Put(key, binVal); err != nil {
		slog.Error("Failed to put object",
			slog.String("database", cmd.DatabaseName),
			slog.String("collection", cmd.CollectionName),
			slog.String("key", key),
			sloki.WrapError(err),
		)

		return commonresponses.InternalServerError, nil
	}

	return commonresponses.OK, nil
}
