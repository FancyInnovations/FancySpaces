package kvcmds

import (
	"encoding/binary"
	"errors"
	"log/slog"
	"time"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/storage/internal/command"
	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/codex"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/commonresponses"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/protocol"
)

// handleSetTTL handles the protocol.ServerCommandKVSetTTL command, which sets a value with a TTL for a given key in the key-value engine.
// Payload format: | Key Length (2 bytes) | Key (variable) | Value (codex-encoded) | Expired At (8 bytes, unix nanos) |
func (c *Commands) handleSetTTL(ctx *command.ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
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
	if len(data) < 2+keyLen+8 {
		return &protocol.Response{
			Code:    protocol.StatusBadRequest,
			Payload: []byte("invalid payload length for key and TTL"),
		}, nil
	}

	key := string(data[2 : 2+keyLen])

	value, err := codex.DecodeValue(data[2+keyLen : len(data)-8])
	if err != nil {
		return &protocol.Response{
			Code:    protocol.StatusBadRequest,
			Payload: []byte("invalid value encoding: " + err.Error()),
		}, nil
	}

	expiresAt := int64(binary.BigEndian.Uint64(data[len(data)-8:]))
	if expiresAt < 0 || expiresAt <= time.Now().UnixNano() {
		return &protocol.Response{
			Code:    protocol.StatusBadRequest,
			Payload: []byte("invalid TTL: expiresAt must be a future timestamp in nanoseconds"),
		}, nil
	}

	kve.SetWithTTL(key, value, expiresAt)

	return &protocol.Response{
		Code:    protocol.StatusOK,
		Payload: *commonresponses.EmptyPayload,
	}, nil
}
