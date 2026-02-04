package app

import (
	"log/slog"
	"net/http"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/storage/internal/command"
	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
	fakeDatabaseDB "github.com/fancyinnovations/fancyspaces/storage/internal/database/databasedb/fake"
	"github.com/fancyinnovations/fancyspaces/storage/internal/engine"
	"github.com/fancyinnovations/fancyspaces/storage/internal/engine/kv/kvcmds"
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
		slog.Error("Could not seed internal databases", sloki.WrapError(err))
		panic(err)
	}

	engineService := engine.NewService(engine.Configuration{
		DatabaseStore: databaseStore,
	})
	if err := engineService.LoadEngines(); err != nil {
		slog.Error("Could not load engines", sloki.WrapError(err))
		panic(err)
	}

	// tcp server
	cmdService := command.NewService()
	cmdService.RegisterHandlers(command.SystemCommands())
	cmdService.RegisterHandlers(kvcmds.Commands())

	// TODO: register more command handlers here, e.g. commands for database operations or engine-specific commands

	srv := server.New(server.Configuration{
		Addr:       ":" + cfg.ServerPort,
		CmdService: cmdService,
	})

	return srv
}
