package issues

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
