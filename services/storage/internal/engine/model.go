package engine

import (
	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
	"github.com/fancyinnovations/fancyspaces/storage/internal/engine/broker"
	"github.com/fancyinnovations/fancyspaces/storage/internal/engine/kv"
)

type Entry struct {
	Database   string
	Collection string
	Type       database.Engine
	engine     any
}

func (e *Entry) AsKeyValueEngine() *kv.Engine {
	if e.Type != database.EngineKeyValue {
		return nil
	}

	return e.engine.(*kv.Engine)
}

func (e *Entry) AsBrokerEngine() *broker.Broker {
	if e.Type != database.EngineBroker {
		return nil
	}

	return e.engine.(*broker.Broker)
}
