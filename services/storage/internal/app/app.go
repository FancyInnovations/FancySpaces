package app

import (
	"net/http"

	"github.com/fancyinnovations/fancyspaces/storage/internal/command"
	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
	fakeDatabaseDB "github.com/fancyinnovations/fancyspaces/storage/internal/database/databasedb/fake"
	"github.com/fancyinnovations/fancyspaces/storage/internal/server"
)

const apiPrefix = "/api/v1"

type Configuration struct {
	Mux        *http.ServeMux
	ServerPort string
}

func Start(cfg Configuration) *server.Server {
	// database
	databaseDB := fakeDatabaseDB.NewDatabaseDB()
	databaseStore := database.NewService(database.Configuration{
		DB: databaseDB,
	})

	if err := seedInternalDatabases(databaseStore); err != nil {
		panic(err)
	}

	// tcp server
	cmdService := command.NewService()

	// TODO: register more command handlers here, e.g. commands for database operations or engine-specific commands

	srv := server.New(server.Configuration{
		Addr:       ":" + cfg.ServerPort,
		CmdService: cmdService,
	})

	return srv
}
