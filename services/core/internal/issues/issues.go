package issues

import (
	"time"

	"github.com/OliverSchlueter/goutils/idgen"
)

type DB interface {
	GetIssues(space string) ([]Issue, error)
	ListIssues(space string, opts ListOptions) ([]Issue, int, error)
	GetIssue(space, id string) (*Issue, error)
	CreateIssue(issue *Issue) error
	UpdateIssue(issue *Issue) error
	DeleteIssue(space, id string) error

	GetComments(space, issue string) ([]Comment, error)
	AddComment(comment *Comment) error
	UpdateComment(comment *Comment) error
	DeleteComment(space, issue, id string) error
}

type ListOptions struct {
	Query          string
	Type           Type
	Status         Status
	Priority       Priority
	Assignee       string
	ExternalSource ExternalSource
	Offset         int
	Limit          int
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

func (s *Store) ListIssues(space string, opts ListOptions) ([]Issue, int, error) {
	if opts.Limit <= 0 || opts.Limit > 200 {
		opts.Limit = 100
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}
	return s.db.ListIssues(space, opts)
}

func (s *Store) GetIssue(space, id string) (*Issue, error) {
	return s.db.GetIssue(space, id)
}

func (s *Store) CreateIssue(issue *Issue) error {
	issue.ID = idgen.GenerateID(8)
	issue.CreatedAt = time.Now()
	issue.UpdatedAt = time.Now()

	if err := issue.Validate(); err != nil {
		return err
	}

	return s.db.CreateIssue(issue)
}

func (s *Store) ForceCreateIssue(issue *Issue) error {
	return s.db.CreateIssue(issue)
}

func (s *Store) UpdateIssue(issue *Issue) error {
	if issue.ArchivedAt != nil {
		return ErrIssueArchived
	}
	issue.UpdatedAt = time.Now()

	if err := issue.Validate(); err != nil {
		return err
	}

	return s.db.UpdateIssue(issue)
}

func (s *Store) ForceUpdateIssue(issue *Issue) error {
	return s.db.UpdateIssue(issue)
}

func (s *Store) DeleteIssue(space, id string) error {
	issue, err := s.db.GetIssue(space, id)
	if err != nil {
		return err
	}
	now := time.Now()
	issue.ArchivedAt = &now
	issue.UpdatedAt = now
	return s.db.UpdateIssue(issue)
}

func (s *Store) GetComments(space, issue string) ([]Comment, error) {
	return s.db.GetComments(space, issue)
}

func (s *Store) AddComment(comment *Comment) error {
	if _, err := s.db.GetIssue(comment.Space, comment.Issue); err != nil {
		return err
	}
	if len(comment.Content) == 0 || len(comment.Content) > 10_000 {
		return ErrCommentTooLong
	}

	comment.ID = idgen.GenerateID(8)
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	return s.db.AddComment(comment)
}

func (s *Store) UpdateComment(comment *Comment) error {
	comment.UpdatedAt = time.Now()

	return s.db.UpdateComment(comment)
}

func (s *Store) DeleteComment(space, issue, id string) error {
	return s.db.DeleteComment(space, issue, id)
}
