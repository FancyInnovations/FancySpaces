package fake

import (
	"context"
	"sync"

	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
)

type DB struct {
	dbs   map[string]database.Database
	colls map[string]map[string]database.Collection
	mu    sync.RWMutex
}

func NewDatabaseDB() *DB {
	return &DB{
		dbs:   make(map[string]database.Database),
		colls: make(map[string]map[string]database.Collection),
	}
}

func (dbdb *DB) GetDatabase(_ context.Context, name string) (*database.Database, error) {
	dbdb.mu.RLock()
	defer dbdb.mu.RUnlock()

	if db, exists := dbdb.dbs[name]; exists {
		return &db, nil
	}

	return nil, database.ErrDatabaseNotFound
}

func (dbdb *DB) GetAllDatabases(_ context.Context) ([]*database.Database, error) {
	dbdb.mu.RLock()
	defer dbdb.mu.RUnlock()

	var dbs []*database.Database
	for _, db := range dbdb.dbs {
		dbCopy := db
		dbs = append(dbs, &dbCopy)
	}

	return dbs, nil
}

func (dbdb *DB) CreateDatabase(_ context.Context, db database.Database) error {
	dbdb.mu.Lock()
	defer dbdb.mu.Unlock()

	if _, exists := dbdb.dbs[db.Name]; exists {
		return database.ErrDatabaseAlreadyExists
	}

	dbdb.dbs[db.Name] = db
	return nil
}

func (dbdb *DB) UpdateDatabase(_ context.Context, db database.Database) error {
	dbdb.mu.Lock()
	defer dbdb.mu.Unlock()

	if _, exists := dbdb.dbs[db.Name]; !exists {
		return database.ErrDatabaseNotFound
	}

	dbdb.dbs[db.Name] = db
	return nil
}

func (dbdb *DB) DeleteDatabase(_ context.Context, name string) error {
	dbdb.mu.Lock()
	defer dbdb.mu.Unlock()

	if _, exists := dbdb.dbs[name]; !exists {
		return database.ErrDatabaseNotFound
	}

	delete(dbdb.dbs, name)
	return nil
}

func (dbdb *DB) GetCollection(_ context.Context, db string, name string) (*database.Collection, error) {
	dbdb.mu.RLock()
	defer dbdb.mu.RUnlock()

	colls, exists := dbdb.colls[db]
	if !exists {
		return nil, database.ErrCollectionNotFound
	}

	coll, exists := colls[name]
	if !exists {
		return nil, database.ErrCollectionNotFound
	}

	return &coll, nil
}

func (dbdb *DB) GetAllCollections(_ context.Context, db string) ([]*database.Collection, error) {
	dbdb.mu.RLock()
	defer dbdb.mu.RUnlock()

	collsMap, exists := dbdb.colls[db]
	if !exists {
		return nil, nil
	}

	var colls []*database.Collection
	for _, coll := range collsMap {
		collCopy := coll
		colls = append(colls, &collCopy)
	}

	return colls, nil
}

func (dbdb *DB) CreateCollection(_ context.Context, collection database.Collection) error {
	dbdb.mu.Lock()
	defer dbdb.mu.Unlock()

	if _, exists := dbdb.colls[collection.Database]; !exists {
		dbdb.colls[collection.Database] = make(map[string]database.Collection)
	}

	if _, exists := dbdb.colls[collection.Database][collection.Name]; exists {
		return database.ErrCollectionAlreadyExists
	}

	dbdb.colls[collection.Database][collection.Name] = collection
	return nil
}

func (dbdb *DB) UpdateCollection(_ context.Context, collection database.Collection) error {
	dbdb.mu.Lock()
	defer dbdb.mu.Unlock()

	colls, exists := dbdb.colls[collection.Database]
	if !exists {
		return database.ErrCollectionNotFound
	}

	if _, exists := colls[collection.Name]; !exists {
		return database.ErrCollectionNotFound
	}

	colls[collection.Name] = collection
	return nil
}

func (dbdb *DB) DeleteCollection(_ context.Context, db string, name string) error {
	dbdb.mu.Lock()
	defer dbdb.mu.Unlock()

	colls, exists := dbdb.colls[db]
	if !exists {
		return database.ErrCollectionNotFound
	}

	if _, exists := colls[name]; !exists {
		return database.ErrCollectionNotFound
	}

	delete(colls, name)
	return nil
}
