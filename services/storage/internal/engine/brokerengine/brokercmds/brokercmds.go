package brokercmds

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
		protocol.ServerCommandBrokerSubscribe:      c.handleSubscribe,
		protocol.ServerCommandBrokerSubscribeQueue: c.handleSubscribeQueue,
		protocol.ServerCommandBrokerUnsubscribe:    c.handleUnsubscribe,
		protocol.ServerCommandBrokerPublish:        c.handlePublish,
	}
}
