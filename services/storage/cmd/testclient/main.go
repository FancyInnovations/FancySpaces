package main

import (
	"fmt"
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

	if err := c.KVSet("system", "collections", "mykey", 42.2); err != nil {
		slog.Error("Failed to set key", sloki.WrapError(err))
	}

	val, err := c.KVGet("system", "collections", "mykey")
	if err != nil {
		slog.Error("Failed to get key", sloki.WrapError(err))
	} else {
		fmt.Printf("Got value: %v\n", val.AsMap())
	}
}
