package fake

import (
	"sync"

	"github.com/fancyinnovations/fancyspaces/internal/issues"
)

type DB struct {
	Issues   []issues.Issue
	Comments []issues.Comment
	Mu       *sync.Mutex
}

func New() *DB {
	return &DB{
		Issues:   []issues.Issue{},
		Comments: []issues.Comment{},
		Mu:       &sync.Mutex{},
	}
}

func (db *DB) GetIssues(space string) ([]issues.Issue, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	result := make([]issues.Issue, 0)
	for _, issue := range db.Issues {
		if issue.Space == space {
			result = append(result, issue)
		}
	}

	return result, nil
}

func (db *DB) GetIssue(space, id string) (*issues.Issue, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	for _, issue := range db.Issues {
		if issue.Space == space && issue.ID == id {
			return &issue, nil
		}
	}

	return nil, issues.ErrIssueNotFound
}

func (db *DB) CreateIssue(issue *issues.Issue) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	for _, item := range db.Issues {
		if item.ID == issue.ID {
			return issues.ErrIssueAlreadyExists
		}
	}

	db.Issues = append(db.Issues, *issue)
	return nil
}

func (db *DB) UpdateIssue(issue *issues.Issue) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	for i, item := range db.Issues {
		if item.ID == issue.ID {
			db.Issues[i] = *issue
			return nil
		}
	}

	return issues.ErrIssueNotFound
}

func (db *DB) DeleteIssue(space, id string) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	for i, item := range db.Issues {
		if item.Space == space && item.ID == id {
			db.Issues = append(db.Issues[:i], db.Issues[i+1:]...)
			return nil
		}
	}

	return issues.ErrIssueNotFound
}

func (db *DB) GetComments(issue string) ([]issues.Comment, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	result := make([]issues.Comment, 0)
	for _, comment := range db.Comments {
		if comment.Issue == issue {
			result = append(result, comment)
		}
	}

	return result, nil
}

func (db *DB) AddComment(comment *issues.Comment) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	for _, item := range db.Comments {
		if item.ID == comment.ID {
			return issues.ErrCommentAlreadyExists
		}
	}

	db.Comments = append(db.Comments, *comment)
	return nil
}

func (db *DB) UpdateComment(comment *issues.Comment) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	for i, item := range db.Comments {
		if item.ID == comment.ID && item.Issue == comment.Issue {
			db.Comments[i] = *comment
			return nil
		}
	}

	return issues.ErrCommentNotFound
}

func (db *DB) DeleteComment(issue, id string) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	for i, item := range db.Comments {
		if item.Issue == issue && item.ID == id {
			db.Comments = append(db.Comments[:i], db.Comments[i+1:]...)
			return nil
		}
	}

	return issues.ErrCommentNotFound
}
