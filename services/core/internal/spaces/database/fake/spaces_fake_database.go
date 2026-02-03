package fake

import (
	"sync"

	"github.com/fancyinnovations/fancyspaces/internal/spaces"
)

type DB struct {
	Items []spaces.Space
	Mu    *sync.Mutex
}

func New() *DB {
	return &DB{
		Items: []spaces.Space{},
		Mu:    &sync.Mutex{},
	}
}

func (db *DB) GetByID(id string) (*spaces.Space, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	for _, s := range db.Items {
		if s.ID == id {
			return &s, nil
		}
	}

	return nil, spaces.ErrSpaceNotFound
}

func (db *DB) GetBySlug(slug string) (*spaces.Space, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	for _, s := range db.Items {
		if s.Slug == slug {
			return &s, nil
		}
	}

	return nil, spaces.ErrSpaceNotFound
}

func (db *DB) GetAll() ([]spaces.Space, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	return db.Items, nil
}

func (db *DB) Create(s *spaces.Space) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	for _, item := range db.Items {
		if item.ID == s.ID || item.Slug == s.Slug {
			return spaces.ErrSpaceAlreadyExists
		}
	}

	db.Items = append(db.Items, *s)
	return nil
}

func (db *DB) Update(id string, s *spaces.Space) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	for i, item := range db.Items {
		if item.ID == id {
			db.Items[i] = *s
			return nil
		}
	}

	return spaces.ErrSpaceNotFound
}

func (db *DB) Delete(id string) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	for i, item := range db.Items {
		if item.ID == id {
			db.Items = append(db.Items[:i], db.Items[i+1:]...)
			return nil
		}
	}

	return nil
}
