package objectcmds

import (
	"github.com/fancyinnovations/fancyspaces/storage/internal/command"
	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
	"github.com/fancyinnovations/fancyspaces/storage/internal/engine"
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
		//protocol.ServerCommandObjectPut:         c.handleObjectPut,
		//protocol.ServerCommandObjectGet:         c.handleObjectGet,
		//protocol.ServerCommandObjectGetMetadata: c.handleObjectGetMetadata,
		//protocol.ServerCommandObjectDelete:      c.handleObjectDelete,
		//protocol.ServerCommandObjectExists:      c.handleObjectExists,
		//protocol.ServerCommandObjectList:        c.handleObjectList,
		//protocol.ServerCommandObjectCopy:        c.handleObjectCopy,
		//protocol.ServerCommandObjectMove:        c.handleObjectMove,
		//protocol.ServerCommandObjectRename:      c.handleObjectRename,
	}
}
