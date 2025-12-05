package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/OliverSchlueter/goutils/containers"
	"github.com/OliverSchlueter/goutils/env"
	"github.com/OliverSchlueter/goutils/middleware"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/src/internal/app"
	"github.com/fancyinnovations/fancyspaces/src/internal/auth"
	"github.com/justinas/alice"
)

const (
	mongodbUrlEnv = "MONGODB_URL"

	clickhouseUser     = "CLICKHOUSE_USER"
	clickhousePassword = "CLICKHOUSE_PASSWORD"
)

func main() {
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

	// Setup databases
	mc := containers.ConnectToMongo(env.MustGetStr(mongodbUrlEnv), "fancyspaces")
	ch := containers.ConnectToClickhouse(
		"fancyanalytics.net:9000",
		"fancyspaces",
		env.MustGetStr(clickhouseUser),
		env.MustGetStr(clickhousePassword),
		"fancyspaces",
		"1.0.0",
	)

	// Setup HTTP server
	mux := http.NewServeMux()
	port := "8080"

	app.Start(app.Configuration{
		Mux:        mux,
		Mongo:      mc,
		ClickHouse: ch,
	})

	auth.ApiKey = env.MustGetStr("API_KEY")

	go func() {
		middleware.OnlyLogStatusAbove = 399 // log 4xx and 5xx status codes
		chain := alice.New(
			middleware.RequestLogging,
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

		slog.Info("Shutdown complete")
	}
}
