package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/OliverSchlueter/goutils/broker"
	"github.com/OliverSchlueter/goutils/containers"
	"github.com/OliverSchlueter/goutils/middleware"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/cloud/internal/app"
	"github.com/fancyinnovations/fancyspaces/integrations/idp-go-sdk/idp"
	"github.com/fancyinnovations/fancyspaces/integrations/idp-go-sdk/keys"
	"github.com/justinas/alice"
)

func main() {
	nc := containers.ConnectToNatsE2E()
	b := broker.NewNatsBroker(&broker.NatsConfiguration{Nats: nc})

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

	// Setup HTTP server
	port := "8095"
	mux := http.NewServeMux()

	app.Start(app.Configuration{
		Mux: mux,
	})

	_, pubKey, err := keys.GetOrGenerateRSAKeys(keys.E2EKeyPath)
	if err != nil {
		slog.Error("Failed to get RSA keys", sloki.WrapError(err))
		os.Exit(1)
	}

	a := idp.NewService(idp.Configuration{
		PublicKey:      pubKey,
		Broker:         b,
		ExcludedRoutes: []string{},
	})

	chain := alice.New(
		middleware.RequestLogging,
		a.HTTPMiddleware,
		middleware.Recovery,
	).Then(mux)

	go func() {
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

		containers.DisconnectNats(nc)

		slog.Info("Shutdown complete")
	}
}
