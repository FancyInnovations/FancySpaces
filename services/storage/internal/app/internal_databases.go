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

	users := map[string]database.PermissionLevel{
		"oliver": database.PermissionLevelAdmin,
	}
	if err := dbStore.UpdateDatabaseUsers(ctx, systemDB, users); err != nil {
		return fmt.Errorf("failed to create internal database 'system': %w", err)
	}

	// key-value test collection
	if err := dbStore.CreateCollectionIfNotExists(ctx, systemDB, "kv_test", database.EngineKeyValue); err != nil {
		return fmt.Errorf("failed to create internal collection 'kv_test' in database 'system': %w", err)
	}

	// object test collection
	if err := dbStore.CreateCollectionIfNotExists(ctx, systemDB, "obj_test", database.EngineObject); err != nil {
		return fmt.Errorf("failed to create internal collection 'obj_test' in database 'system': %w", err)
	}

	// broker test collection
	if err := dbStore.CreateCollectionIfNotExists(ctx, systemDB, "broker_test", database.EngineBroker); err != nil {
		return fmt.Errorf("failed to create internal collection 'broker_test' in database 'system': %w", err)
	}

	return nil
}
