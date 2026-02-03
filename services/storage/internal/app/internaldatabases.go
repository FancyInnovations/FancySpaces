package app

import (
	"context"
	"fmt"

	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
)

func seedInternalDatabases(dbStore *database.Store) error {
	ctx := context.Background()

	// system database
	if err := dbStore.CreateDatabaseIfNotExists(ctx, "system"); err != nil {
		return fmt.Errorf("failed to create internal database 'system': %w", err)
	}
	systemDB, err := dbStore.GetDatabase(ctx, "system")
	if err != nil {
		return fmt.Errorf("failed to get internal database 'system': %w", err)
	}

	// databases collection
	if err := dbStore.CreateCollectionIfNotExists(ctx, systemDB, "databases", database.EngineKeyValue); err != nil {
		return fmt.Errorf("failed to create internal collection 'databases' in database 'system': %w", err)
	}

	// collections collection
	if err := dbStore.CreateCollectionIfNotExists(ctx, systemDB, "collections", database.EngineKeyValue); err != nil {
		return fmt.Errorf("failed to create internal collection 'collections' in database 'system': %w", err)
	}

	return nil
}
