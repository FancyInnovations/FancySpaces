package spaces

import (
	"errors"
	"fmt"
	"time"

	"github.com/OliverSchlueter/goutils/idgen"
	"github.com/fancyinnovations/fancyspaces/integrations/idp-go-sdk/idp"
	"github.com/fancyinnovations/fancyspaces/integrations/spaces-go-sdk/spaces"
)

type DB interface {
	GetByID(id string) (*spaces.Space, error)
	GetBySlug(slug string) (*spaces.Space, error)
	GetForUser(userID string) ([]spaces.Space, error)
	GetAll() ([]spaces.Space, error)
	Create(s *spaces.Space) error
	Update(id string, s *spaces.Space) error
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
func (s *Store) Get(id string) (*spaces.Space, error) {
	sp, err := s.db.GetBySlug(id)
	if err != nil {
		if errors.Is(err, spaces.ErrSpaceNotFound) {
			sp, err = s.db.GetByID(id)
			if err == nil {
				return sp, nil
			}
		}

		return nil, err
	}

	return sp, nil
}

func (s *Store) GetByID(id string) (*spaces.Space, error) {
	return s.db.GetByID(id)
}

func (s *Store) GetBySlug(slug string) (*spaces.Space, error) {
	return s.db.GetBySlug(slug)
}

func (s *Store) GetForUser(userID string) ([]spaces.Space, error) {
	return s.db.GetForUser(userID)
}

func (s *Store) GetAll() ([]spaces.Space, error) {
	return s.db.GetAll()
}

func (s *Store) Create(creator *idp.User, req *CreateOrUpdateSpaceReq) (*spaces.Space, error) {
	if !creator.IsActive {
		return nil, ErrUserNotActive
	}
	if !creator.Verified {
		return nil, ErrUserNotVerified
	}

	space := &spaces.Space{
		ID:          idgen.GenerateID(8),
		Slug:        req.Slug,
		Title:       req.Title,
		Description: req.Description,
		Categories:  req.Categories,
		IconURL:     req.IconURL,
		Status:      spaces.StatusDraft,
		CreatedAt:   time.Now(),
		Creator:     creator.ID,
		Members:     []spaces.Member{},
	}

	// check if slug is already taken by another space
	if _, err := s.db.GetBySlug(space.Slug); err == nil {
		return nil, ErrSpaceAlreadyExists
	}

	if err := space.Validate(); err != nil {
		return nil, fmt.Errorf("invalid space: %w", err)
	}

	if err := s.db.Create(space); err != nil {
		return nil, err
	}

	return space, nil
}

func (s *Store) CreateFull(space *spaces.Space) error {
	// check if slug is already taken by another space
	if _, err := s.db.GetBySlug(space.Slug); err == nil {
		return ErrSpaceAlreadyExists
	}

	if err := space.Validate(); err != nil {
		return fmt.Errorf("invalid space: %w", err)
	}

	return s.db.Create(space)
}

func (s *Store) Update(id string, req *CreateOrUpdateSpaceReq) error {
	space, err := s.db.GetByID(id)
	if err != nil {
		return err
	}

	// check if slug is already taken by another space
	if space.Slug != req.Slug {
		if _, err := s.db.GetBySlug(req.Slug); !errors.Is(err, spaces.ErrSpaceNotFound) {
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

func (s *Store) UpdateFull(space *spaces.Space) error {
	if err := space.Validate(); err != nil {
		return fmt.Errorf("invalid space: %w", err)
	}

	return s.db.Update(space.ID, space)
}

func (s *Store) Delete(id string) error {
	return s.db.Delete(id)
}

func (s *Store) ChangeStatus(space *spaces.Space, to spaces.Status) error {
	if to == space.Status {
		return nil // no change
	}

	space.Status = to

	if err := space.Validate(); err != nil {
		return fmt.Errorf("invalid space: %w", err)
	}

	return s.db.Update(space.ID, space)
}
