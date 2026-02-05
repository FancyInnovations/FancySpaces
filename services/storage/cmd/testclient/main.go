package main

import (
	"log/slog"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/client"
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

	c, err := client.NewClient(client.Configuration{
		Host:   "localhost",
		Port:   "8091",
		ApiKey: "hello",
		//Username: "oliver",
		//Password: "hello",
	})
	if err != nil {
		slog.Error("Failed to connect", sloki.WrapError(err))
	}

	_ = c

	if err := c.KVSet("system", "collections", "mykey", 42); err != nil {
		slog.Error("Failed to set key", sloki.WrapError(err))
	}

	if err := c.KVSet("system", "collections", "mykey2", 69); err != nil {
		slog.Error("Failed to set key", sloki.WrapError(err))
	}

	if err := c.KVSet("system", "collections", "mykey3", 123); err != nil {
		slog.Error("Failed to set key", sloki.WrapError(err))
	}

	keys, err := c.KVKeys("system", "collections")
	if err != nil {
		slog.Error("Failed to get keys", sloki.WrapError(err))
	} else {
		slog.Info("Keys in collection", slog.Any("keys", keys))
	}
}
