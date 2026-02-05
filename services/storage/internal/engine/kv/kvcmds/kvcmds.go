package kvcmds

import (
	"encoding/binary"
	"errors"
	"log/slog"
	"time"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/storage/internal/auth"
	"github.com/fancyinnovations/fancyspaces/storage/internal/command"
	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
	"github.com/fancyinnovations/fancyspaces/storage/internal/engine"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/codex"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/commonresponses"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/protocol"
)

type Commands struct {
	dbStore       *database.Store
	engineService *engine.Service
}

type Configuration struct {
	DatabaseStore *database.Store
	EngineService *engine.Service
}

func New(cfg Configuration) *Commands {
	return &Commands{
		dbStore:       cfg.DatabaseStore,
		engineService: cfg.EngineService,
	}
}

func (c *Commands) Get() map[uint16]command.Handler {
	return map[uint16]command.Handler{
		protocol.CommandKVGet:    c.handleGet,
		protocol.CommandKVSet:    c.handleSet,
		protocol.CommandKVSetTTL: c.handleSetTTL,
		protocol.CommandKVDelete: c.handleDelete,
	}
}

// handleGet handles the protocol.CommandKVGet command, which retrieves the value for a given key from the key-value engine.
// Payload format: | Key Length (2 bytes) | Key (variable) |
func (c *Commands) handleGet(ctx *command.ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
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

	val := kve.Get(key)
	if val == nil {
		return &protocol.Response{
			Code:    protocol.StatusNotFound,
			Payload: *commonresponses.EmptyPayload,
		}, nil
	}

	return &protocol.Response{
		Code:    protocol.StatusOK,
		Payload: codex.EncodeValue(val),
	}, nil
}

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

// handleSetTTL handles the protocol.CommandKVSetTTL command, which sets a value with a TTL for a given key in the key-value engine.
// Payload format: | Key Length (2 bytes) | Key (variable) | Value (codex-encoded) | Expired At (8 bytes, unix nanos) |
func (c *Commands) handleSetTTL(ctx *command.ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
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

// handleDelete handles the protocol.CommandKVDelete command, which deletes a key-value pair for a given key in the key-value engine.
// Payload format: | Key Length (2 bytes) | Key (variable) |
func (c *Commands) handleDelete(ctx *command.ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
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

	kve.Delete(key)

	return &protocol.Response{
		Code:    protocol.StatusOK,
		Payload: *commonresponses.EmptyPayload,
	}, nil
}
