package issues

import (
	"time"

	"github.com/OliverSchlueter/goutils/idgen"
)

type DB interface {
	GetIssues(space string) ([]Issue, error)
	GetIssue(space, id string) (*Issue, error)
	CreateIssue(issue *Issue) error
	UpdateIssue(issue *Issue) error
	DeleteIssue(space, id string) error

	GetComments(issue string) ([]Comment, error)
	AddComment(comment *Comment) error
	UpdateComment(comment *Comment) error
	DeleteComment(issue, id string) error
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

func (s *Store) GetIssues(space string) ([]Issue, error) {
	return s.db.GetIssues(space)
}

func (s *Store) GetIssue(space, id string) (*Issue, error) {
	return s.db.GetIssue(space, id)
}

func (s *Store) CreateIssue(issue *Issue) error {
	if _, err := s.db.GetIssue(issue.Space, issue.ID); err == nil {
		return ErrIssueAlreadyExists
	}

	issue.ID = idgen.GenerateID(8)
	issue.CreatedAt = time.Now()
	issue.UpdatedAt = time.Now()

	if err := issue.Validate(); err != nil {
		return err
	}

	return s.db.CreateIssue(issue)
}

func (s *Store) UpdateIssue(issue *Issue) error {
	issue.UpdatedAt = time.Now()

	if err := issue.Validate(); err != nil {
		return err
	}

	return s.db.UpdateIssue(issue)
}

func (s *Store) DeleteIssue(space, id string) error {
	return s.db.DeleteIssue(space, id)
}

func (s *Store) GetComments(issue string) ([]Comment, error) {
	return s.db.GetComments(issue)
}

func (s *Store) AddComment(comment *Comment) error {
	if _, err := s.db.GetIssue(comment.Issue, comment.ID); err != nil {
		return ErrCommentAlreadyExists
	}

	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	return s.db.AddComment(comment)
}

func (s *Store) UpdateComment(comment *Comment) error {
	comment.UpdatedAt = time.Now()

	return s.db.UpdateComment(comment)
}

func (s *Store) DeleteComment(issue, id string) error {
	return s.db.DeleteComment(issue, id)
}
