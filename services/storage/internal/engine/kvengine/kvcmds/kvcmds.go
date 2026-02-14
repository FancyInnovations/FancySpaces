package kvcmds

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/OliverSchlueter/goutils/problems"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/protocol"
	"github.com/fancyinnovations/fancyspaces/storage/internal/auth"
	"github.com/fancyinnovations/fancyspaces/storage/internal/command"
	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
	"github.com/fancyinnovations/fancyspaces/storage/internal/engine"
)

type Commands struct {
	dbStore       *database.Store
	engineService *engine.Service
	userFromCtx   func(ctx context.Context) *auth.User
}

type Configuration struct {
	DatabaseStore *database.Store
	EngineService *engine.Service
	UserFromCtx   func(ctx context.Context) *auth.User
}

func New(cfg Configuration) *Commands {
	return &Commands{
		dbStore:       cfg.DatabaseStore,
		engineService: cfg.EngineService,
		userFromCtx:   cfg.UserFromCtx,
	}
}

func (c *Commands) Get() map[uint16]command.Handler {
	return map[uint16]command.Handler{
		protocol.ServerCommandKVSet:            c.handleSet,
		protocol.ServerCommandKVSetTTL:         c.handleSetTTL,
		protocol.ServerCommandKVDelete:         c.handleDelete,
		protocol.ServerCommandKVDeleteMultiple: c.handleDelete,
		protocol.ServerCommandKVDeleteAll:      c.handleDeleteAll,
		protocol.ServerCommandKVExists:         c.handleExists,
		protocol.ServerCommandKVGet:            c.handleGet,
		protocol.ServerCommandKVGetMultiple:    c.handleGetMultiple,
		protocol.ServerCommandKVGetAll:         c.handleGetAll,
		protocol.ServerCommandKVGetTTL:         c.handleGetTTL,
		protocol.ServerCommandKVGetMultipleTTL: c.handleGetMultipleTTL,
		protocol.ServerCommandKVGetAllTTL:      c.handleGetAllTTL,
		protocol.ServerCommandKVKeys:           c.handleKeys,
		protocol.ServerCommandKVCount:          c.handleCountTCP,
		protocol.ServerCommandKVSize:           c.handleSize,
	}
}

func (c *Commands) RegisterHTTP(prefix string, mux *http.ServeMux) {
	mux.HandleFunc(prefix+"/databases/{db_id}/collections/{coll_name}/kv/{cmd_id}", c.handleCommandHTTP)
}

func (c *Commands) handleCommandHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		problems.MethodNotAllowed(r.Method, []string{http.MethodPost}).WriteToHTTP(w)
		return
	}

	u := c.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}

	dbid := r.PathValue("db_id")
	if dbid == "" {
		problems.ValidationError("db_id", "Database ID is required").WriteToHTTP(w)
		return
	}

	db, err := c.dbStore.GetDatabase(r.Context(), dbid)
	if err != nil {
		if errors.Is(err, database.ErrDatabaseNotFound) {
			problems.NotFound("Database", dbid).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get database by id", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}
	if !db.IsMember(u.ID) && !u.IsAdmin() {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	collname := r.PathValue("coll_name")
	if collname == "" {
		problems.ValidationError("coll_name", "Collection name is required").WriteToHTTP(w)
		return
	}

	coll, err := c.dbStore.GetCollection(r.Context(), db, collname)
	if err != nil {
		if errors.Is(err, database.ErrCollectionNotFound) {
			problems.NotFound("Collection", collname).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get collection by name", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	e, err := c.engineService.GetEngine(db.Name, coll.Name)
	if err != nil {
		if errors.Is(err, database.ErrCollectionNotFound) {
			problems.NotFound("Collection", collname).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get engine",
			slog.String("database", db.Name),
			slog.String("collection", coll.Name),
			sloki.WrapError(err),
		)
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	if e.Type != database.EngineKeyValue {
		problems.ValidationError("Engine", "Collection is not a key-value collection").WriteToHTTP(w)
		return
	}

	kve := e.AsKeyValueEngine()

	cmdIDStr := r.PathValue("cmd_id")
	if cmdIDStr == "" {
		problems.ValidationError("cmd_id", "Command ID is required").WriteToHTTP(w)
		return
	}
	cmdID, err := strconv.ParseUint(cmdIDStr, 10, 16)
	if err != nil {
		problems.ValidationError("cmd_id", "Invalid Command ID").WriteToHTTP(w)
		return
	}

	switch uint16(cmdID) {
	case protocol.ServerCommandKVDelete:
		c.handleDeleteHTTP(w, r, db, coll, kve)
	case protocol.ServerCommandKVDeleteMultiple:
		c.handleDeleteMultipleHTTP(w, r, db, coll, kve)
	case protocol.ServerCommandKVDeleteAll:
		c.handleDeleteAllHTTP(w, r, db, coll, kve)
	case protocol.ServerCommandKVExists:
		c.handleExistsHTTP(w, r, db, coll, kve)
	case protocol.ServerCommandKVGet:
		c.handleGetHTTP(w, r, db, coll, kve)
	case protocol.ServerCommandKVGetMultiple:
		c.handleGetMultipleHTTP(w, r, db, coll, kve)
	case protocol.ServerCommandKVGetAll:
		c.handleGetAllHTTP(w, r, db, coll, kve)
	case protocol.ServerCommandKVKeys:
		c.handleKeysHTTP(w, r, db, coll, kve)
	case protocol.ServerCommandKVCount:
		c.handleCountHTTP(w, r, db, coll, kve)
	case protocol.ServerCommandKVSize:
		c.handleSizeHTTP(w, r, db, coll, kve)
	default:
		problems.NotFound("Command", cmdIDStr).WriteToHTTP(w)
	}
}
