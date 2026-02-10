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
	"github.com/fancyinnovations/fancyspaces/storage/pkg/commonresponses"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/protocol"
)

// handleGetMetadata processes a get metadata command for an object engine.
// Payload format: | Key Length (2 bytes) | Key (variable) |
// Response payload format: | Size (8 bytes) | CRC32 (4 bytes) | CreatedAt (8 bytes) | ModifiedAt (8 bytes) |
func (c *Commands) handleGetMetadata(ctx *command.ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
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

	omd, err := obje.GetMeta(key)
	if err != nil {
		if errors.Is(err, objectengine.ErrKeyNotFound) {
			return &protocol.Response{
				Code:    protocol.StatusNotFound,
				Payload: *commonresponses.EmptyPayload,
			}, nil
		}

		slog.Error("Failed to get object metadata",
			slog.String("database", cmd.DatabaseName),
			slog.String("collection", cmd.CollectionName),
			slog.String("key", key),
			sloki.WrapError(err),
		)
		return commonresponses.InternalServerError, nil
	}

	totalLen := 8 + 4 + 8 + 8
	payload := make([]byte, totalLen)

	// Size
	binary.BigEndian.PutUint64(payload[0:8], uint64(omd.Size))

	// CRC32
	binary.BigEndian.PutUint32(payload[8:12], omd.Checksum)

	// CreatedAt
	binary.BigEndian.PutUint64(payload[12:20], uint64(omd.CreatedAt))

	// ModifiedAt
	binary.BigEndian.PutUint64(payload[20:28], uint64(omd.ModifiedAt))

	return &protocol.Response{
		Code:    protocol.StatusOK,
		Payload: payload,
	}, nil
}
