package engine

import (
	"context"
	"log/slog"
	"sync"

	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
	"github.com/fancyinnovations/fancyspaces/storage/internal/engine/broker"
	"github.com/fancyinnovations/fancyspaces/storage/internal/engine/kv"
)

type Service struct {
	dbStore *database.Store

	engines             map[string]*Entry
	enginesMu           sync.RWMutex
	sendBrokerMessage   func(db, coll, connID, subject string, msgs [][]byte)
	isConnectionHealthy func(connID string) bool
}

type Configuration struct {
	DatabaseStore       *database.Store
	SendBrokerMessage   func(db, coll, connID, subject string, msgs [][]byte)
	IsConnectionHealthy func(connID string) bool
}

func NewService(cfg Configuration) *Service {
	return &Service{
		dbStore:             cfg.DatabaseStore,
		engines:             make(map[string]*Entry),
		sendBrokerMessage:   cfg.SendBrokerMessage,
		isConnectionHealthy: cfg.IsConnectionHealthy,
	}
}

func (s *Service) LoadEngines() error {
	ctx := context.Background()

	// TODO Add a function to load all collections of all databases in one call

	dbs, err := s.dbStore.GetAllDatabases(ctx)
	if err != nil {
		return err
	}

	colls := make([]*database.Collection, 0)
	for _, db := range dbs {
		dbColls, err := s.dbStore.GetAllCollections(ctx, db)
		if err != nil {
			return err
		}
		colls = append(colls, dbColls...)
	}

	for _, coll := range colls {
		var e any

		switch coll.Engine {
		case database.EngineKeyValue:
			e = kv.NewEngine(kv.Configuration{
				DisableTTL: coll.KVSettings != nil && coll.KVSettings.DisableTTL,
			})
		case database.EngineBroker:
			e = broker.NewBroker(broker.Configuration{
				PublishCallback: func(sub *broker.Subscriber, subject string, msgs [][]byte) {
					s.sendBrokerMessage(coll.Database, coll.Name, sub.ID, subject, msgs)
				},
				IsClientHealthy: s.isConnectionHealthy,
			})
		}

		if e == nil {
			slog.Warn(
				"Unknown engine type for collection",
				slog.String("database", coll.Database),
				slog.String("collection", coll.Name),
				slog.String("engine", string(coll.Engine)),
			)
			continue
		}

		entry := &Entry{
			Database:   coll.Database,
			Collection: coll.Name,
			Type:       coll.Engine,
			engine:     e,
		}
		s.enginesMu.Lock()
		s.engines[toKey(coll.Database, coll.Name)] = entry
		s.enginesMu.Unlock()
	}

	return nil
}

func (s *Service) GetEngine(db, coll string) (*Entry, error) {
	s.enginesMu.RLock()
	defer s.enginesMu.RUnlock()

	entry, exists := s.engines[toKey(db, coll)]
	if !exists {
		return nil, database.ErrCollectionNotFound
	}

	return entry, nil
}

func toKey(database, collection string) string {
	return database + "_" + collection
}
