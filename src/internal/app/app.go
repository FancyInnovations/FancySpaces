package app

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/fancyinnovations/fancyspaces/src/internal/auth"
	"github.com/fancyinnovations/fancyspaces/src/internal/spaces"
	"github.com/fancyinnovations/fancyspaces/src/internal/spaces/database/sqlite"
	"github.com/fancyinnovations/fancyspaces/src/internal/spaces/handler"
)

const apiPrefix = "/api/v1"

type Configuration struct {
	Mux    *http.ServeMux
	SQLite *sql.DB
}

func Start(cfg Configuration) {
	spacesDB := sqlite.New(sqlite.Configuration{
		Conn: cfg.SQLite,
	})
	if err := spacesDB.Setup(); err != nil {
		panic(fmt.Errorf("failed to setup spaces database: %w", err))
	}
	spacesStore := spaces.New(spaces.Configuration{
		DB: spacesDB,
	})
	spacesHandler := handler.New(handler.Configuration{
		Store:       spacesStore,
		UserFromCtx: auth.UserFromContext,
	})
	spacesHandler.Register(apiPrefix, cfg.Mux)

}
