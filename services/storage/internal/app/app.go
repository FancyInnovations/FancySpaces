package app

import (
	"net/http"

	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
	fakeDatabaseDB "github.com/fancyinnovations/fancyspaces/storage/internal/database/databasedb/fake"
)

const apiPrefix = "/api/v1"

type Configuration struct {
	Mux *http.ServeMux
}

func Start(cfg Configuration) {

	databaseDB := fakeDatabaseDB.NewDatabaseDB()
	databaseStore := database.NewService(database.Configuration{
		DB: databaseDB,
	})

	if err := seedInternalDatabases(databaseStore); err != nil {
		panic(err)
	}
}
