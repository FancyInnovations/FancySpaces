package fake

import (
	"sync"

	spacesStore "github.com/fancyinnovations/fancyspaces/core/internal/spaces"
	"github.com/fancyinnovations/fancyspaces/integrations/spaces-go-sdk/spaces"
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

func (db *DB) GetForCreator(userID string) ([]spaces.Space, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	var result []spaces.Space
	for _, s := range db.Items {
		if s.Creator == userID {
			result = append(result, s)
			continue
		}
	}

	return result, nil
}

func (db *DB) GetForCategory(category string) ([]spaces.Space, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	var result []spaces.Space
	for _, s := range db.Items {
		for _, c := range s.Categories {
			if string(c) == category {
				result = append(result, s)
				break
			}
		}
	}

	return result, nil
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
			return spacesStore.ErrSpaceAlreadyExists
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
