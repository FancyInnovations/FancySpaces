package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OliverSchlueter/goutils/containers"
	"github.com/OliverSchlueter/goutils/middleware"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/internal/app"
	"github.com/fancyinnovations/fancyspaces/internal/auth"
	"github.com/fancyinnovations/fancyspaces/internal/fflags"
	"github.com/justinas/alice"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Feature flags
	fflags.DisableIssueSyncer.Enable()

	// Setup logging
	logService := sloki.NewService(sloki.Configuration{
		URL:          "http://localhost:3100/loki/api/v1/push",
		Service:      "fancyspaces",
		ConsoleLevel: slog.LevelDebug,
		LokiLevel:    slog.LevelInfo,
		EnableLoki:   false,
		Handlers:     []sloki.LogHandler{},
	})
	slog.SetDefault(slog.New(logService))

	// Connect to databases
	mc := containers.ConnectToMongoE2E("fancyspaces_e2e")
	ch := containers.ConnectToClickhouseE2E("fancyspaces_e2e")
	mio := containers.ConnectToMinIOE2E()

	// Setup default admin user
	auth.Users["oliver"] = &auth.User{
		ID:        "oliver",
		Provider:  auth.ProviderBasic,
		Name:      "Oliver",
		Email:     "oliver@fancyinnovations.com",
		Verified:  true,
		Password:  auth.Hash("hello"),
		Roles:     []string{"admin", "user"},
		CreatedAt: time.Date(2025, 12, 3, 19, 0, 0, 0, time.UTC),
		IsActive:  true,
		Metadata: map[string]string{
			"api_key": "hello",
		},
	}

	// Setup HTTP server
	mux := http.NewServeMux()
	port := "8080"

	app.Start(app.Configuration{
		Mux:        mux,
		Mongo:      mc,
		ClickHouse: ch,
		MinIO:      mio,
	})

	go func() {
		chain := alice.New(
			middleware.RequestLogging,
			auth.Middleware,
			middleware.Recovery,
		).Then(mux)

		err := http.ListenAndServe(":"+port, chain)
		if err != nil {
			slog.Error("Could not start server on port "+port, sloki.WrapError(err))
			os.Exit(1)
		}
	}()

	slog.Info(fmt.Sprintf("Started server on http://localhost:%s\n", port))

	// Wait for a signal to exit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	switch <-sig {
	case os.Interrupt:
		slog.Info("Received interrupt signal, shutting down...")

		containers.DisconnectMongo(mc)
		containers.DisconnectClickhouse(ch)
		containers.DisconnectMinIO(mio)

		slog.Info("Shutdown complete")
	}
}
