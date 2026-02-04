package database

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type DB interface {
	GetDatabase(ctx context.Context, name string) (*Database, error)
	GetAllDatabases(ctx context.Context) ([]*Database, error)
	CreateDatabase(ctx context.Context, db Database) error
	UpdateDatabase(ctx context.Context, db Database) error
	DeleteDatabase(ctx context.Context, name string) error

	GetCollection(ctx context.Context, db string, name string) (*Collection, error)
	GetAllCollections(ctx context.Context, db string) ([]*Collection, error)
	CreateCollection(ctx context.Context, collection Collection) error
	UpdateCollection(ctx context.Context, collection Collection) error
	DeleteCollection(ctx context.Context, db string, name string) error
}

type Store struct {
	db DB
}

type Configuration struct {
	DB DB
}

func NewService(cfg Configuration) *Store {
	return &Store{
		db: cfg.DB,
	}
}

func (s *Store) GetDatabase(ctx context.Context, name string) (*Database, error) {
	return s.db.GetDatabase(ctx, name)
}

func (s *Store) GetAllDatabases(ctx context.Context) ([]*Database, error) {
	return s.db.GetAllDatabases(ctx)
}

func (s *Store) CreateDatabase(ctx context.Context, name string) error {
	_, err := s.db.GetDatabase(ctx, name)
	if err == nil {
		return ErrDatabaseAlreadyExists
	}

	db := Database{
		Name:      name,
		CreatedAt: time.Now(),
	}

	return s.db.CreateDatabase(ctx, db)
}

func (s *Store) CreateDatabaseIfNotExists(ctx context.Context, name string) error {
	if err := s.CreateDatabase(ctx, name); err != nil && !errors.Is(err, ErrDatabaseAlreadyExists) {
		return err
	}

	return nil
}

func (s *Store) UpdateDatabaseUsers(ctx context.Context, db *Database, updatedUsers map[string]PermissionLevel) error {
	db.Users = updatedUsers
	return s.db.UpdateDatabase(ctx, *db)
}

func (s *Store) DeleteDatabase(ctx context.Context, name string) error {
	colls, err := s.db.GetAllCollections(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to get collections for database %s: %w", name, err)
	}

	for _, coll := range colls {
		if err := s.db.DeleteCollection(ctx, name, coll.Name); err != nil {
			return fmt.Errorf("failed to delete collection %s in database %s: %w", coll.Name, name, err)
		}
	}

	return s.db.DeleteDatabase(ctx, name)
}

func (s *Store) GetCollection(ctx context.Context, db *Database, name string) (*Collection, error) {
	return s.db.GetCollection(ctx, db.Name, name)
}

func (s *Store) GetAllCollections(ctx context.Context, db *Database) ([]*Collection, error) {
	return s.db.GetAllCollections(ctx, db.Name)
}

func (s *Store) CreateCollection(ctx context.Context, db *Database, name string, engine Engine) error {
	_, err := s.db.GetCollection(ctx, db.Name, name)
	if err == nil {
		return ErrCollectionAlreadyExists
	}

	coll := Collection{
		Database:  db.Name,
		Name:      name,
		CreatedAt: time.Now(),
		Engine:    engine,
	}

	return s.db.CreateCollection(ctx, coll)
}

func (s *Store) CreateCollectionIfNotExists(ctx context.Context, db *Database, name string, engine Engine) error {
	if err := s.CreateCollection(ctx, db, name, engine); err != nil && !errors.Is(err, ErrCollectionAlreadyExists) {
		return err
	}

	return nil
}

func (s *Store) ChangeCollectionName(ctx context.Context, coll *Collection, newName string) error {
	coll.Name = newName

	return s.db.UpdateCollection(ctx, *coll)
}

func (s *Store) DeleteCollection(ctx context.Context, db *Database, name string) error {
	_, err := s.db.GetCollection(ctx, db.Name, name)
	if err != nil {
		return ErrCollectionNotFound
	}

	return s.db.DeleteCollection(ctx, db.Name, name)
}
