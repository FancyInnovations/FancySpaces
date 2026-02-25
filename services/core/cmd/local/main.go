package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/OliverSchlueter/goutils/broker"
	"github.com/OliverSchlueter/goutils/containers"
	"github.com/OliverSchlueter/goutils/middleware"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/core/internal/app"
	"github.com/fancyinnovations/fancyspaces/core/internal/auth"
	"github.com/fancyinnovations/fancyspaces/core/internal/fflags"
	"github.com/fancyinnovations/fancyspaces/integrations/idp-go-sdk/idp"
	"github.com/fancyinnovations/fancyspaces/integrations/idp-go-sdk/keys"
	"github.com/justinas/alice"
)

func main() {
	nc := containers.ConnectToNatsE2E()
	b := broker.NewNatsBroker(&broker.NatsConfiguration{Nats: nc})

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
	auth.Users["oliver"] = &idp.User{
		ID:        "oliver",
		Provider:  idp.ProviderBasic,
		Name:      "Oliver",
		Email:     "oliver@fancyinnovations.com",
		Verified:  true,
		Password:  idp.PasswordHash("hello"),
		Roles:     []string{"admin", "user"},
		CreatedAt: time.Date(2025, 12, 3, 19, 0, 0, 0, time.UTC),
		IsActive:  true,
		Metadata: map[string]string{
			"api_key": "hello",
		},
	}

	// Setup HTTP server
	port := "8080"
	mux := http.NewServeMux()
	mavenMux := http.NewServeMux()

	app.Start(app.Configuration{
		Mux:              mux,
		MavenMux:         mavenMux,
		Broker:           b,
		Mongo:            mc,
		ClickHouse:       ch,
		MinIO:            mio,
		SecretsMasterKey: []byte("fooooooooooooooobarrrrrrrrrrrrrr"), // 32 bytes for AES-256
	})

	hostRouter := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		if strings.HasPrefix(host, "maven.") {
			mavenMux.ServeHTTP(w, r)
		} else {
			mux.ServeHTTP(w, r)
		}
	})

	pubKey, err := keys.LoadPublicKey(keys.E2EKeyPath)
	if err != nil {
		slog.Error("Failed to load public key", sloki.WrapError(err))
		os.Exit(1)
	}
	idp.ServiceBaseURL = "http://localhost:8083"
	a := idp.NewService(idp.Configuration{
		Broker:         b,
		PublicKey:      pubKey,
		ExcludedRoutes: []string{},
	})

	chain := alice.New(
		middleware.RequestLogging,
		a.HTTPMiddleware,
		middleware.Recovery,
	).Then(hostRouter)

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
		containers.DisconnectMongo(mc)
		containers.DisconnectClickhouse(ch)
		containers.DisconnectMinIO(mio)

		slog.Info("Shutdown complete")
	}
}
