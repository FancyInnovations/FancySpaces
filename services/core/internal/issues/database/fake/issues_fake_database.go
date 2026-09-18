package fake

import (
	"sort"
	"strings"
	"sync"

	"github.com/fancyinnovations/fancyspaces/core/internal/issues"
)

type DB struct {
	Issues     []issues.Issue
	Comments   []issues.Comment
	Activities []issues.Activity
	Mu         *sync.Mutex
}

func New() *DB {
	return &DB{
		Issues:     []issues.Issue{},
		Comments:   []issues.Comment{},
		Activities: []issues.Activity{},
		Mu:         &sync.Mutex{},
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

func (db *DB) ListIssues(space string, opts issues.ListOptions) ([]issues.Issue, int, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	result := make([]issues.Issue, 0)
	query := strings.ToLower(strings.TrimSpace(opts.Query))
	for _, issue := range db.Issues {
		if issue.Space != space || issue.ArchivedAt != nil {
			continue
		}
		if query != "" && !strings.Contains(strings.ToLower(issue.Title), query) && !strings.Contains(strings.ToLower(issue.Description), query) && !strings.Contains(strings.ToLower(issue.ID), query) {
			continue
		}
		if opts.Type != "" && issue.Type != opts.Type || opts.Status != "" && issue.Status != opts.Status || opts.Priority != "" && issue.Priority != opts.Priority || opts.Assignee != "" && issue.Assignee != opts.Assignee || opts.ExternalSource != "" && issue.ExternalSource != opts.ExternalSource {
			continue
		}
		if opts.Label != "" {
			matched := false
			for _, label := range issue.Labels {
				if strings.EqualFold(label, opts.Label) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}
		result = append(result, issue)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].UpdatedAt.After(result[j].UpdatedAt) })
	total := len(result)
	start := opts.Offset
	if start > total {
		start = total
	}
	end := start + opts.Limit
	if end > total {
		end = total
	}
	return result[start:end], total, nil
}

func (db *DB) GetIssue(space, id string) (*issues.Issue, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	for _, issue := range db.Issues {
		if issue.Space == space && issue.ID == id && issue.ArchivedAt == nil {
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
		if item.Space == issue.Space && item.ID == issue.ID {
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

func (db *DB) GetComments(space, issue string) ([]issues.Comment, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	result := make([]issues.Comment, 0)
	for _, comment := range db.Comments {
		if comment.Space == space && comment.Issue == issue {
			result = append(result, comment)
		}
	}

	return result, nil
}

func (db *DB) AddComment(comment *issues.Comment) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	for _, item := range db.Comments {
		if item.Space == comment.Space && item.Issue == comment.Issue && item.ID == comment.ID {
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

func (db *DB) DeleteComment(space, issue, id string) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	for i, item := range db.Comments {
		if item.Space == space && item.Issue == issue && item.ID == id {
			db.Comments = append(db.Comments[:i], db.Comments[i+1:]...)
			return nil
		}
	}

	return issues.ErrCommentNotFound
}

func (db *DB) GetActivities(space, issue string) ([]issues.Activity, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()
	result := make([]issues.Activity, 0)
	for _, activity := range db.Activities {
		if activity.Space == space && activity.Issue == issue {
			result = append(result, activity)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].CreatedAt.Before(result[j].CreatedAt) })
	return result, nil
}

func (db *DB) AddActivity(activity *issues.Activity) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()
	db.Activities = append(db.Activities, *activity)
	return nil
}
