package kvcmds

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

// handleSet handles the protocol.CommandKVSet command, which sets a value for a given key in the key-value engine.
// Payload format: | Key Length (2 bytes) | Key (variable) | Value (codex-encoded) |
func (c *Commands) handleSet(ctx *command.ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
	u := auth.UserFromContext(ctx.Ctx)
	if u == nil || !u.Verified || !u.IsActive {
		return commonresponses.Unauthorized, nil
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

	value, err := codex.DecodeValue(data[2+keyLen:])
	if err != nil {
		return &protocol.Response{
			Code:    protocol.StatusBadRequest,
			Payload: []byte("invalid value encoding: " + err.Error()),
		}, nil
	}

	kve.Set(key, value)

	return &protocol.Response{
		Code:    protocol.StatusOK,
		Payload: *commonresponses.EmptyPayload,
	}, nil
}
