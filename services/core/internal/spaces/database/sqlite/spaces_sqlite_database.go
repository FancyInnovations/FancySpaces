package sqlite

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/fancyinnovations/fancyspaces/internal/spaces"
)

//go:embed sqlStatements/createSpacesTable.sql
var createMembersTable string

//go:embed sqlStatements/createMembersTable.sql
var createSpacesTable string

//go:embed sqlStatements/createCategoriesTable.sql
var createCategoriesTable string

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
	_, err := db.conn.Exec(createSpacesTable, nil)
	if err != nil {
		return err
	}

	_, err = db.conn.Exec(createMembersTable, nil)
	if err != nil {
		return err
	}

	_, err = db.conn.Exec(createCategoriesTable, nil)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetByID(id string) (*spaces.Space, error) {
	query := `SELECT id, slug, title, description, icon_url, status, created_at FROM spaces WHERE id = ?`
	row := db.conn.QueryRow(query, id)

	var s spaces.Space
	err := row.Scan(&s.ID, &s.Slug, &s.Title, &s.Description, &s.IconURL, &s.Status, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, spaces.ErrSpaceNotFound
		}

		return nil, fmt.Errorf("failed to scan space: %w", err)
	}

	s.Members, err = db.fetchMembersForSpace(s.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch members for space: %w", err)
	}

	s.Categories, err = db.fetchCategoriesForSpace(s.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories for space: %w", err)
	}

	return &s, nil
}

func (db *DB) GetBySlug(slug string) (*spaces.Space, error) {
	query := `SELECT id, slug, title, description, icon_url, status, created_at FROM spaces WHERE slug = ?`
	row := db.conn.QueryRow(query, slug)

	var s spaces.Space
	err := row.Scan(&s.ID, &s.Slug, &s.Title, &s.Description, &s.IconURL, &s.Status, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, spaces.ErrSpaceNotFound
		}

		return nil, fmt.Errorf("failed to scan space: %w", err)
	}

	s.Members, err = db.fetchMembersForSpace(s.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch members for space: %w", err)
	}

	s.Categories, err = db.fetchCategoriesForSpace(s.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories for space: %w", err)
	}

	return &s, nil
}

func (db *DB) GetAll() ([]spaces.Space, error) {
	query := `SELECT id, slug, title, description, icon_url, status, created_at FROM spaces`
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query spaces: %w", err)
	}
	defer rows.Close()

	var spacesList = make([]spaces.Space, 0)
	for rows.Next() {
		var s spaces.Space
		err := rows.Scan(&s.ID, &s.Slug, &s.Title, &s.Description, &s.IconURL, &s.Status, &s.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan space: %w", err)
		}

		s.Members, err = db.fetchMembersForSpace(s.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch members for space: %w", err)
		}

		s.Categories, err = db.fetchCategoriesForSpace(s.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch categories for space: %w", err)
		}

		spacesList = append(spacesList, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return spacesList, nil
}

func (db *DB) Create(s *spaces.Space) error {
	stmt := `INSERT INTO spaces (id, slug, title, description, icon_url, status, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := db.conn.Exec(stmt, s.ID, s.Slug, s.Title, s.Description, s.IconURL, s.Status, s.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert space: %w", err)
	}

	for _, m := range s.Members {
		memberStmt := `INSERT INTO space_members (space_id, user_id, role) VALUES (?, ?, ?)`
		_, err := db.conn.Exec(memberStmt, s.ID, m.UserID, m.Role)
		if err != nil {
			return fmt.Errorf("failed to insert member: %w", err)
		}
	}

	for _, c := range s.Categories {
		categoryStmt := `INSERT INTO space_categories (space_id, category) VALUES (?, ?)`
		_, err := db.conn.Exec(categoryStmt, s.ID, c)
		if err != nil {
			return fmt.Errorf("failed to insert category: %w", err)
		}
	}

	return nil
}

func (db *DB) Update(id string, s *spaces.Space) error {
	stmt := `UPDATE spaces SET slug = ?, title = ?, description = ?, icon_url = ?, status = ? WHERE id = ?`
	_, err := db.conn.Exec(stmt, s.Slug, s.Title, s.Description, s.IconURL, s.Status, id)
	if err != nil {
		return fmt.Errorf("failed to update space: %w", err)
	}

	// For simplicity, we'll delete existing members and categories and re-insert them.
	delMembersStmt := `DELETE FROM space_members WHERE space_id = ?`
	_, err = db.conn.Exec(delMembersStmt, id)
	if err != nil {
		return fmt.Errorf("failed to delete existing members: %w", err)
	}

	for _, m := range s.Members {
		memberStmt := `INSERT INTO space_members (space_id, user_id, role) VALUES (?, ?, ?)`
		_, err := db.conn.Exec(memberStmt, id, m.UserID, m.Role)
		if err != nil {
			return fmt.Errorf("failed to insert member: %w", err)
		}
	}

	// Similarly for categories
	delCategoriesStmt := `DELETE FROM space_categories WHERE space_id = ?`
	_, err = db.conn.Exec(delCategoriesStmt, id)
	if err != nil {
		return fmt.Errorf("failed to delete existing categories: %w", err)
	}

	for _, c := range s.Categories {
		categoryStmt := `INSERT INTO categories (space_id, category) VALUES (?, ?)`
		_, err := db.conn.Exec(categoryStmt, id, c)
		if err != nil {
			return fmt.Errorf("failed to insert category: %w", err)
		}
	}

	return nil
}

func (db *DB) Delete(id string) error {
	stmt := `DELETE FROM spaces WHERE id = ?`
	_, err := db.conn.Exec(stmt, id)
	if err != nil {
		return fmt.Errorf("failed to delete space: %w", err)
	}

	// Also delete associated members and categories
	delMembersStmt := `DELETE FROM space_members WHERE space_id = ?`
	_, err = db.conn.Exec(delMembersStmt, id)
	if err != nil {
		return fmt.Errorf("failed to delete members: %w", err)
	}

	delCategoriesStmt := `DELETE FROM space_categories WHERE space_id = ?`
	_, err = db.conn.Exec(delCategoriesStmt, id)
	if err != nil {
		return fmt.Errorf("failed to delete categories: %w", err)
	}

	return nil
}

func (db *DB) fetchMembersForSpace(id string) ([]spaces.Member, error) {
	query := `SELECT user_id, role FROM space_members WHERE space_id = ?`
	rows, err := db.conn.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query members: %w", err)
	}
	defer rows.Close()

	var members []spaces.Member
	for rows.Next() {
		var m spaces.Member
		err := rows.Scan(&m.UserID, &m.Role)
		if err != nil {
			return nil, fmt.Errorf("failed to scan member: %w", err)
		}
		members = append(members, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return members, nil
}

func (db *DB) fetchCategoriesForSpace(id string) ([]spaces.Category, error) {
	query := `SELECT category FROM space_categories WHERE space_id = ?`
	rows, err := db.conn.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []spaces.Category
	for rows.Next() {
		var c spaces.Category
		err := rows.Scan(&c)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return categories, nil
}
