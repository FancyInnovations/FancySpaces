package dbcmds

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

// handleCollectionGet processes a command to retrieve collection information.
// Payload format: | DB name length (2 bytes) | DB name (variable) | Collection name length (2 bytes) | Collection name (variable) |
func (c *Commands) handleCollectionGet(ctx *command.ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
	u := auth.UserFromContext(ctx.Ctx)
	if u == nil || !u.Verified || !u.IsActive {
		return commonresponses.Unauthorized, nil
	}

	if len(cmd.Payload) < 2 {
		return &protocol.Response{
			Code:    protocol.StatusBadRequest,
			Payload: []byte("invalid payload length"),
		}, nil
	}

	dbNameLen := int(binary.BigEndian.Uint16(cmd.Payload[:2]))
	if len(cmd.Payload) < 2+dbNameLen {
		return &protocol.Response{
			Code:    protocol.StatusBadRequest,
			Payload: []byte("invalid payload length"),
		}, nil
	}

	dbName := string(cmd.Payload[2 : 2+dbNameLen])

	db, err := c.dbStore.GetDatabase(ctx.Ctx, dbName)
	if err != nil {
		if errors.Is(err, database.ErrDatabaseNotFound) {
			return commonresponses.DatabaseNotFound, nil
		}

		slog.Error("Failed to get database",
			slog.String("database", dbName),
			sloki.WrapError(err),
		)
		return commonresponses.InternalServerError, nil
	}

	if !u.IsAdmin() && !db.HasPermission(u.ID, database.PermissionLevelReadOnly) {
		return commonresponses.Forbidden, nil
	}

	if len(cmd.Payload) < 2+dbNameLen+2 {
		return &protocol.Response{
			Code:    protocol.StatusBadRequest,
			Payload: []byte("invalid payload length"),
		}, nil
	}

	collNameLen := int(binary.BigEndian.Uint16(cmd.Payload[2+dbNameLen : 2+dbNameLen+2]))
	if len(cmd.Payload) < 2+dbNameLen+2+collNameLen {
		return &protocol.Response{
			Code:    protocol.StatusBadRequest,
			Payload: []byte("invalid payload length"),
		}, nil
	}

	collName := string(cmd.Payload[2+dbNameLen+2 : 2+dbNameLen+2+collNameLen])

	coll, err := c.dbStore.GetCollection(ctx.Ctx, db, collName)
	if err != nil {
		if errors.Is(err, database.ErrCollectionNotFound) {
			return commonresponses.CollectionNotFound, nil
		}

		slog.Error("Failed to get collection",
			slog.String("database", dbName),
			slog.String("collection", collName),
			sloki.WrapError(err),
		)
		return commonresponses.InternalServerError, nil
	}

	data, err := codex.Marshal(coll)
	if err != nil {
		slog.Error("Failed to marshal collection",
			slog.String("database", cmd.DatabaseName),
			slog.String("collection", collName),
			sloki.WrapError(err),
		)
		return commonresponses.InternalServerError, nil
	}

	return &protocol.Response{
		Code:    protocol.StatusOK,
		Payload: codex.EncodeBinary(data),
	}, nil
}
