package clickhouse

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/fancyinnovations/fancyspaces/internal/analytics"
)

//go:embed create_database.sql
var createDatabaseSQL string

//go:embed create_version_downloads_table.sql
var createVersionDownloadsTableSQL string

//go:embed create_maven_artifact_downloads_table.sql
var createMavenArtifactDownloadsTableSQL string

type DB struct {
	ch driver.Conn
}

type Configuration struct {
	CH driver.Conn
}

func NewDB(cfg *Configuration) *DB {
	return &DB{
		ch: cfg.CH,
	}
}

func (db *DB) Setup(ctx context.Context) error {
	if err := db.ch.Exec(ctx, createDatabaseSQL); err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	if err := db.ch.Exec(ctx, createVersionDownloadsTableSQL); err != nil {
		return fmt.Errorf("failed to create version_downloads table: %w", err)
	}

	if err := db.ch.Exec(ctx, createMavenArtifactDownloadsTableSQL); err != nil {
		return fmt.Errorf("failed to create maven_artifact_downloads table: %w", err)
	}

	return nil
}

func (db *DB) GetDownloadCountForSpace(ctx context.Context, spaceID string) (uint64, error) {
	var count uint64
	query := `
		SELECT COUNT(*) 
		FROM fancyspaces.version_downloads 
		WHERE space_id = ?`
	if err := db.ch.QueryRow(ctx, query, spaceID).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (db *DB) GetDownloadCountForVersion(ctx context.Context, spaceID, versionID string) (uint64, error) {
	var count uint64
	query := `
		SELECT COUNT(*) 
		FROM fancyspaces.version_downloads 
		WHERE space_id = ? 
	    AND version_id = ?`
	if err := db.ch.QueryRow(ctx, query, spaceID, versionID).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (db *DB) GetDownloadCountForVersions(ctx context.Context, spaceID string) (map[string]uint64, error) {
	var res []struct {
		VersionID string `ch:"version_id"`
		Downloads uint64 `ch:"downloads"`
	}

	query := `
		SELECT
			version_id,
			count(*) AS downloads
		FROM fancyspaces.version_downloads
		WHERE space_id = ?
		GROUP BY space_id, version_id`
	if err := db.ch.Select(ctx, &res, query, spaceID); err != nil {
		return nil, err
	}

	asMap := make(map[string]uint64)
	for _, r := range res {
		asMap[r.VersionID] = r.Downloads
	}

	return asMap, nil
}

func (db *DB) StoreVersionDownloads(ctx context.Context, records []analytics.VersionDownload) error {
	batch, err := db.ch.PrepareBatch(ctx, "INSERT INTO fancyspaces.version_downloads")
	if err != nil {
		return err
	}

	for _, r := range records {
		if err := batch.AppendStruct(&r); err != nil {
			return err
		}
	}

	if err := batch.Send(); err != nil {
		return err
	}

	return nil
}

func (db *DB) StoreMavenArtifactDownloads(ctx context.Context, records []analytics.MavenArtifactDownload) error {
	batch, err := db.ch.PrepareBatch(ctx, "INSERT INTO fancyspaces.maven_artifact_downloads")
	if err != nil {
		return err
	}

	for _, r := range records {
		if err := batch.AppendStruct(&r); err != nil {
			return err
		}
	}

	if err := batch.Send(); err != nil {
		return err
	}

	return nil
}
