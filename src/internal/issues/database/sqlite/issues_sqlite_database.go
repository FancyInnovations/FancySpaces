package sqlite

import (
	"database/sql"
	_ "embed"

	"github.com/fancyinnovations/fancyspaces/internal/issues"
)

//go:embed sqlStatements/createIssuesTable.sql
var createIssuesTable string

//go:embed sqlStatements/createCommentsTable.sql
var createCommentsTable string

type DB struct {
	conn *sql.DB
}

type Configuration struct {
	Conn *sql.DB
}

func New(cfg Configuration) *DB {
	return &DB{
		conn: cfg.Conn,
	}
}

func (db *DB) Setup() error {
	_, err := db.conn.Exec(createIssuesTable, nil)
	if err != nil {
		return err
	}

	_, err = db.conn.Exec(createCommentsTable, nil)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetIssues(space string) ([]issues.Issue, error) {
	query := `SELECT id, space, title, description, type, status, priority, assignee, reporter, created_at, updated_at, external_source FROM issues WHERE space = ?`
	rows, err := db.conn.Query(query, space)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issuesList []issues.Issue
	for rows.Next() {
		var i issues.Issue
		err := rows.Scan(&i.ID, &i.Space, &i.Title, &i.Description, &i.Type, &i.Status, &i.Priority, &i.Assignee, &i.Reporter, &i.CreatedAt, &i.UpdatedAt, &i.ExternalSource)
		if err != nil {
			return nil, err
		}
		issuesList = append(issuesList, i)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return issuesList, nil
}

func (db *DB) GetIssue(space, id string) (*issues.Issue, error) {
	query := `SELECT id, space, title, description, type, status, priority, assignee, reporter, created_at, updated_at, external_source FROM issues WHERE space = ? AND id = ?`
	row := db.conn.QueryRow(query, space, id)

	var i issues.Issue
	err := row.Scan(&i.ID, &i.Space, &i.Title, &i.Description, &i.Type, &i.Status, &i.Priority, &i.Assignee, &i.Reporter, &i.CreatedAt, &i.UpdatedAt, &i.ExternalSource)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, issues.ErrIssueNotFound
		}
		return nil, err
	}

	return &i, nil
}

func (db *DB) CreateIssue(issue *issues.Issue) error {
	query := `INSERT INTO issues (id, space, title, description, type, status, priority, assignee, reporter, created_at, updated_at, external_source) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := db.conn.Exec(query, issue.ID, issue.Space, issue.Title, issue.Description, issue.Type, issue.Status, issue.Priority, issue.Assignee, issue.Reporter, issue.CreatedAt, issue.UpdatedAt, issue.ExternalSource)

	if err != nil {
		if sqliteErr, ok := err.(interface{ ErrorCode() int }); ok && sqliteErr.ErrorCode() == 19 {
			return issues.ErrIssueAlreadyExists
		}
	}

	return err
}

func (db *DB) UpdateIssue(issue *issues.Issue) error {
	query := `UPDATE issues SET title = ?, description = ?, type = ?, status = ?, priority = ?, assignee = ?, reporter = ?, updated_at = ?, external_source = ? WHERE space = ? AND id = ?`
	_, err := db.conn.Exec(query, issue.Title, issue.Description, issue.Type, issue.Status, issue.Priority, issue.Assignee, issue.Reporter, issue.UpdatedAt, issue.ExternalSource, issue.Space, issue.ID)
	return err
}

func (db *DB) DeleteIssue(space, id string) error {
	query := `DELETE FROM issues WHERE space = ? AND id = ?`
	_, err := db.conn.Exec(query, space, id)
	return err
}

func (db *DB) GetComments(issue string) ([]issues.Comment, error) {
	query := `SELECT id, issue, author, content, created_at, updated_at FROM issue_comments WHERE issue = ?`
	rows, err := db.conn.Query(query, issue)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commentsList []issues.Comment
	for rows.Next() {
		var c issues.Comment
		err := rows.Scan(&c.ID, &c.Issue, &c.Author, &c.Content, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		commentsList = append(commentsList, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return commentsList, nil
}

func (db *DB) AddComment(comment *issues.Comment) error {
	query := `INSERT INTO issue_comments (id, issue, author, content, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.conn.Exec(query, comment.ID, comment.Issue, comment.Author, comment.Content, comment.CreatedAt, comment.UpdatedAt)
	return err
}

func (db *DB) UpdateComment(comment *issues.Comment) error {
	query := `UPDATE issue_comments SET content = ?, updated_at = ? WHERE issue = ? AND id = ?`
	_, err := db.conn.Exec(query, comment.Content, comment.UpdatedAt, comment.Issue, comment.ID)
	return err
}

func (db *DB) DeleteComment(issue, id string) error {
	query := `DELETE FROM issue_comments WHERE issue = ? AND id = ?`
	_, err := db.conn.Exec(query, issue, id)
	return err
}
