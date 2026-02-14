package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OliverSchlueter/goutils/middleware"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/storage/internal/app"
	"github.com/fancyinnovations/fancyspaces/storage/internal/auth"
	"github.com/justinas/alice"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Setup logging
	logService := sloki.NewService(sloki.Configuration{
		URL:          "http://localhost:3100/loki/api/v1/push",
		Service:      "fancyspaces-storage",
		ConsoleLevel: slog.LevelDebug,
		LokiLevel:    slog.LevelInfo,
		EnableLoki:   false,
		Handlers:     []sloki.LogHandler{},
	})
	slog.SetDefault(slog.New(logService))

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
	httpPort := "8090"
	serverPort := "8091"
	mux := http.NewServeMux()

	srv := app.Start(app.Configuration{
		Mux:        mux,
		ServerPort: serverPort,
	})

	chain := alice.New(
		middleware.RequestLogging,
		auth.Middleware,
		middleware.Recovery,
	).Then(mux)

	go func() {
		if err := http.ListenAndServe(":"+httpPort, chain); err != nil {
			slog.Error("Could not start http server on port "+httpPort, sloki.WrapError(err))
			os.Exit(1)
		}
	}()

	go func() {
		if err := srv.Run(); err != nil {
			slog.Error("Could not start tcp server on port "+serverPort, sloki.WrapError(err))
			os.Exit(1)
		}
	}()

	slog.Info(fmt.Sprintf("Started http server on http://localhost:%s", httpPort))
	slog.Info(fmt.Sprintf("Started TCP server on localhost:%s", serverPort))

	// Wait for a signal to exit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	switch <-sig {
	case os.Interrupt:
		slog.Info("Received interrupt signal, shutting down...")

		slog.Info("Shutdown complete")
	}
}
