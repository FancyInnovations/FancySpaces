package engine

import (
	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
	"github.com/fancyinnovations/fancyspaces/storage/internal/engine/brokerengine"
	"github.com/fancyinnovations/fancyspaces/storage/internal/engine/kvengine"
	"github.com/fancyinnovations/fancyspaces/storage/internal/engine/objectengine"
)

type Entry struct {
	Database   string
	Collection string
	Type       database.Engine
	engine     any
}

func (e *Entry) AsKeyValueEngine() *kvengine.Engine {
	if e.Type != database.EngineKeyValue {
		return nil
	}

	return e.engine.(*kvengine.Engine)
}

func (e *Entry) AsObjectEngine() *objectengine.Bucket {
	if e.Type != database.EngineObject {
		return nil
	}

	return e.engine.(*objectengine.Bucket)
}

func (e *Entry) AsBrokerEngine() *brokerengine.Broker {
	if e.Type != database.EngineBroker {
		return nil
	}

	return e.engine.(*brokerengine.Broker)
}
