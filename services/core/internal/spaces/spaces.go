package spaces

import (
	"errors"
	"fmt"
	"time"

	"github.com/OliverSchlueter/goutils/idgen"
	"github.com/fancyinnovations/fancyspaces/core/internal/auth"
)

type DB interface {
	GetByID(id string) (*Space, error)
	GetBySlug(slug string) (*Space, error)
	GetAll() ([]Space, error)
	Create(s *Space) error
	Update(id string, s *Space) error
	Delete(id string) error
}

type Store struct {
	db DB
}

type Configuration struct {
	DB DB
}

func New(cfg Configuration) *Store {
	return &Store{
		db: cfg.DB,
	}
}

// Get tries to get a space by ID first, then by slug if not found by ID.
func (s *Store) Get(id string) (*Space, error) {
	sp, err := s.db.GetBySlug(id)
	if err != nil {
		if errors.Is(err, ErrSpaceNotFound) {
			sp, err = s.db.GetByID(id)
			if err == nil {
				return sp, nil
			}
		}

		return nil, err
	}

	return sp, nil
}

func (s *Store) GetByID(id string) (*Space, error) {
	return s.db.GetByID(id)
}

func (s *Store) GetBySlug(slug string) (*Space, error) {
	return s.db.GetBySlug(slug)
}

func (s *Store) GetAll() ([]Space, error) {
	return s.db.GetAll()
}

func (s *Store) Create(creator *auth.User, req *CreateOrUpdateSpaceReq) (*Space, error) {
	if !creator.IsActive {
		return nil, ErrUserNotActive
	}
	if !creator.Verified {
		return nil, ErrUserNotVerified
	}

	space := &Space{
		ID:          idgen.GenerateID(8),
		Slug:        req.Slug,
		Title:       req.Title,
		Description: req.Description,
		Categories:  req.Categories,
		IconURL:     req.IconURL,
		Status:      StatusDraft,
		CreatedAt:   time.Now(),
		Members: []Member{
			{
				UserID: creator.ID,
				Role:   RoleOwner,
			},
		},
	}

	// check if slug is already taken by another space
	if space.Slug != req.Slug {
		if _, err := s.db.GetBySlug(req.Slug); err != nil {
			return nil, ErrSpaceAlreadyExists
		}
	}

	if err := space.Validate(); err != nil {
		return nil, fmt.Errorf("invalid space: %w", err)
	}

	if err := s.db.Create(space); err != nil {
		return nil, err
	}

	return space, nil
}

func (s *Store) Update(id string, req *CreateOrUpdateSpaceReq) error {
	space, err := s.db.GetByID(id)
	if err != nil {
		return err
	}

	// check if slug is already taken by another space
	if space.Slug != req.Slug {
		if _, err := s.db.GetBySlug(req.Slug); !errors.Is(err, ErrSpaceNotFound) {
			return ErrSpaceAlreadyExists
		}
	}

	space.Slug = req.Slug
	space.Title = req.Title
	space.Description = req.Description
	space.Categories = req.Categories
	space.IconURL = req.IconURL

	if err := space.Validate(); err != nil {
		return fmt.Errorf("invalid space: %w", err)
	}

	return s.db.Update(id, space)
}

func (s *Store) UpdateFull(space *Space) error {
	if err := space.Validate(); err != nil {
		return fmt.Errorf("invalid space: %w", err)
	}

	return s.db.Update(space.ID, space)
}

func (s *Store) Delete(id string) error {
	return s.db.Delete(id)
}

func (s *Store) ChangeStatus(space *Space, to Status) error {
	if to == space.Status {
		return nil // no change
	}

	space.Status = to

	if err := space.Validate(); err != nil {
		return fmt.Errorf("invalid space: %w", err)
	}

	return s.db.Update(space.ID, space)
}
