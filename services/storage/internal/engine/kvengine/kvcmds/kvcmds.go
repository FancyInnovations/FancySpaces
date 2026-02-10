package kvcmds

import (
	"github.com/fancyinnovations/fancyspaces/storage/internal/command"
	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
	"github.com/fancyinnovations/fancyspaces/storage/internal/engine"
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
		protocol.ServerCommandKVSet:            c.handleSet,
		protocol.ServerCommandKVSetTTL:         c.handleSetTTL,
		protocol.ServerCommandKVDelete:         c.handleDelete,
		protocol.ServerCommandKVDeleteMultiple: c.handleDelete,
		protocol.ServerCommandKVDeleteAll:      c.handleDeleteAll,
		protocol.ServerCommandKVExists:         c.handleExists,
		protocol.ServerCommandKVGet:            c.handleGet,
		protocol.ServerCommandKVGetMultiple:    c.handleGetMultiple,
		protocol.ServerCommandKVGetAll:         c.handleGetAll,
		protocol.ServerCommandKVKeys:           c.handleKeys,
		protocol.ServerCommandKVCount:          c.handleCount,
	}
}
