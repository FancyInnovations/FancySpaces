package sqlite

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fancyinnovations/fancyspaces/src/internal/versions"
)

//go:embed sqlStatements/createVersionsTable.sql
var createVersionsTable string

//go:embed sqlStatements/createVersionFilesTable.sql
var createVersionFilesTable string

//go:embed sqlStatements/createVersionDownloadsTable.sql
var createVersionDownloadsTable string

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
	_, err := db.conn.Exec(createVersionsTable, nil)
	if err != nil {
		return err
	}

	_, err = db.conn.Exec(createVersionFilesTable, nil)
	if err != nil {
		return err
	}

	_, err = db.conn.Exec(createVersionDownloadsTable, nil)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetAll(spaceID string) ([]versions.Version, error) {
	q := `SELECT id, name, channel, published_at, changelog, supported_platform_versions, downloads
	      FROM versions
	      WHERE space_id = ?
	      ORDER BY published_at DESC`
	rows, err := db.conn.Query(q, spaceID)
	if err != nil {
		return nil, fmt.Errorf("query versions: %w", err)
	}
	defer rows.Close()

	var result []versions.Version
	for rows.Next() {
		var (
			id                     string
			name                   string
			channel                string
			publishedAtStr         string
			changelog              sql.NullString
			supportedPlatformsJSON sql.NullString
			downloads              int64
		)
		if err := rows.Scan(&id, &name, &channel, &publishedAtStr, &changelog, &supportedPlatformsJSON, &downloads); err != nil {
			return nil, fmt.Errorf("scan version row: %w", err)
		}

		publishedAt, err := time.Parse(time.RFC3339, publishedAtStr)
		if err != nil {
			return nil, fmt.Errorf("parse published_at: %w", err)
		}

		var supported []string
		if supportedPlatformsJSON.Valid && supportedPlatformsJSON.String != "" {
			if err := json.Unmarshal([]byte(supportedPlatformsJSON.String), &supported); err != nil {
				return nil, fmt.Errorf("unmarshal supported_platform_versions: %w", err)
			}
		}

		files, err := db.getFiles(spaceID, id)
		if err != nil {
			return nil, fmt.Errorf("load files for version %s: %w", id, err)
		}

		v := versions.Version{
			SpaceID:                   spaceID,
			ID:                        id,
			Name:                      name,
			Channel:                   channel,
			PublishedAt:               publishedAt,
			Changelog:                 changelog.String,
			SupportedPlatformVersions: supported,
			Files:                     files,
			Downloads:                 downloads,
		}
		result = append(result, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate versions rows: %w", err)
	}

	return result, nil
}

func (db *DB) GetByID(spaceID, versionID string) (*versions.Version, error) {
	q := `SELECT id, name, channel, published_at, changelog, supported_platform_versions, downloads
	      FROM versions
	      WHERE space_id = ? AND id = ?`
	row := db.conn.QueryRow(q, spaceID, versionID)

	var (
		id                     string
		name                   string
		channel                string
		publishedAtStr         string
		changelog              sql.NullString
		supportedPlatformsJSON sql.NullString
		downloads              int64
	)

	if err := row.Scan(&id, &name, &channel, &publishedAtStr, &changelog, &supportedPlatformsJSON, &downloads); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan version: %w", err)
	}

	publishedAt, err := time.Parse(time.RFC3339, publishedAtStr)
	if err != nil {
		return nil, fmt.Errorf("parse published_at: %w", err)
	}

	var supported []string
	if supportedPlatformsJSON.Valid && supportedPlatformsJSON.String != "" {
		if err := json.Unmarshal([]byte(supportedPlatformsJSON.String), &supported); err != nil {
			return nil, fmt.Errorf("unmarshal supported_platform_versions: %w", err)
		}
	}

	files, err := db.getFiles(spaceID, id)
	if err != nil {
		return nil, fmt.Errorf("load files: %w", err)
	}

	v := &versions.Version{
		SpaceID:                   spaceID,
		ID:                        id,
		Name:                      name,
		Channel:                   channel,
		PublishedAt:               publishedAt,
		Changelog:                 changelog.String,
		SupportedPlatformVersions: supported,
		Files:                     files,
		Downloads:                 downloads,
	}
	return v, nil
}

func (db *DB) GetByName(spaceID, versionNumber string) (*versions.Version, error) {
	q := `SELECT id FROM versions WHERE space_id = ? AND name = ? LIMIT 1`
	row := db.conn.QueryRow(q, spaceID, versionNumber)

	var id string
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find version id by name: %w", err)
	}

	return db.GetByID(spaceID, id)
}

func (db *DB) Create(v *versions.Version) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	supportedJSON, err := json.Marshal(v.SupportedPlatformVersions)
	if err != nil {
		return fmt.Errorf("marshal supported_platform_versions: %w", err)
	}

	publishedAt := v.PublishedAt.Format(time.RFC3339)
	q := `INSERT INTO versions
	      (space_id, id, name, channel, published_at, changelog, supported_platform_versions, downloads)
	      VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	if _, err := tx.Exec(q, v.SpaceID, v.ID, v.Name, v.Channel, publishedAt, v.Changelog, string(supportedJSON), v.Downloads); err != nil {
		return fmt.Errorf("insert version: %w", err)
	}

	if err := db.insertFilesTx(tx, v.SpaceID, v.ID, v.Files); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (db *DB) Update(spaceID, versionID string, v *versions.Version) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	supportedJSON, err := json.Marshal(v.SupportedPlatformVersions)
	if err != nil {
		return fmt.Errorf("marshal supported_platform_versions: %w", err)
	}

	publishedAt := v.PublishedAt.Format(time.RFC3339)
	q := `UPDATE versions SET name = ?, channel = ?, published_at = ?, changelog = ?, supported_platform_versions = ?, downloads = ?
	      WHERE space_id = ? AND id = ?`
	if _, err := tx.Exec(q, v.Name, v.Channel, publishedAt, v.Changelog, string(supportedJSON), v.Downloads, spaceID, versionID); err != nil {
		return fmt.Errorf("update version: %w", err)
	}

	// Replace files: delete existing and insert new
	if _, err := tx.Exec(`DELETE FROM version_files WHERE space_id = ? AND version_id = ?`, spaceID, versionID); err != nil {
		return fmt.Errorf("delete old version_files: %w", err)
	}

	if err := db.insertFilesTx(tx, spaceID, versionID, v.Files); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (db *DB) Delete(spaceID, versionID string) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM version_files WHERE space_id = ? AND version_id = ?`, spaceID, versionID); err != nil {
		return fmt.Errorf("delete version_files: %w", err)
	}

	if _, err := tx.Exec(`DELETE FROM versions WHERE space_id = ? AND id = ?`, spaceID, versionID); err != nil {
		return fmt.Errorf("delete version: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (db *DB) LogDownload(spaceID, versionID string) error {
	q := `INSERT INTO version_downloads (timestamp, space_id, version_id) VALUES (?, ?, ?)`
	downloadedAt := time.Now().Format(time.RFC3339)
	if _, err := db.conn.Exec(q, downloadedAt, spaceID, versionID); err != nil {
		return fmt.Errorf("insert version_download: %w", err)
	}

	return nil
}

func (db *DB) getFiles(spaceID, versionID string) ([]versions.VersionFile, error) {
	q := `SELECT name, url, size FROM version_files WHERE space_id = ? AND version_id = ? ORDER BY name`
	rows, err := db.conn.Query(q, spaceID, versionID)
	if err != nil {
		return nil, fmt.Errorf("query version_files: %w", err)
	}
	defer rows.Close()

	var files []versions.VersionFile
	for rows.Next() {
		var f versions.VersionFile
		if err := rows.Scan(&f.Name, &f.URL, &f.Size); err != nil {
			return nil, fmt.Errorf("scan version_file: %w", err)
		}
		files = append(files, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate version_files: %w", err)
	}
	return files, nil
}

func (db *DB) insertFilesTx(tx *sql.Tx, spaceID, versionID string, files []versions.VersionFile) error {
	if len(files) == 0 {
		return nil
	}
	stmt, err := tx.Prepare(`INSERT INTO version_files (space_id, version_id, name, url, size) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("prepare insert version_files: %w", err)
	}
	defer stmt.Close()

	for _, f := range files {
		if _, err := stmt.Exec(spaceID, versionID, f.Name, f.URL, f.Size); err != nil {
			return fmt.Errorf("insert version_file %s: %w", f.Name, err)
		}
	}
	return nil
}
