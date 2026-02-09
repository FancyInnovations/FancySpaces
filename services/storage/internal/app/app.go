package app

import (
	"log/slog"
	"net/http"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/storage/internal/command"
	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
	fakeDatabaseDB "github.com/fancyinnovations/fancyspaces/storage/internal/database/databasedb/fake"
	"github.com/fancyinnovations/fancyspaces/storage/internal/engine"
	"github.com/fancyinnovations/fancyspaces/storage/internal/engine/broker/brokercmds"
	"github.com/fancyinnovations/fancyspaces/storage/internal/engine/kvengine/kvcmds"
	"github.com/fancyinnovations/fancyspaces/storage/internal/engine/objectengine/objectcmds"
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

	srv := server.New(server.Configuration{
		Addr: ":" + cfg.ServerPort,
	})

	engineService := engine.NewService(engine.Configuration{
		DatabaseStore:       databaseStore,
		SendBrokerMessage:   srv.SendBrokerMessage,
		IsConnectionHealthy: srv.IsConnectionHealthy,
	})
	if err := engineService.LoadEngines(); err != nil {
		slog.Error("Could not load engines", sloki.WrapError(err))
		panic(err)
	}

	// tcp server
	cmdService := command.NewService()
	cmdService.RegisterHandlers(command.SystemCommands())

	kvCommands := kvcmds.New(kvcmds.Configuration{
		DatabaseStore: databaseStore,
		EngineService: engineService,
	})
	cmdService.RegisterHandlers(kvCommands.Get())

	objectCommands := objectcmds.New(objectcmds.Configuration{
		DatabaseStore: databaseStore,
		EngineService: engineService,
	})
	cmdService.RegisterHandlers(objectCommands.Get())

	brokerCommands := brokercmds.New(brokercmds.Configuration{
		DatabaseStore: databaseStore,
		EngineService: engineService,
	})
	cmdService.RegisterHandlers(brokerCommands.Get())

	srv.SetCommandService(cmdService)

	return srv
}
