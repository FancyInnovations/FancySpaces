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
		return
	}

	_ = c

	if err := c.KVSet("system", "collections", "mykey", 42); err != nil {
		slog.Error("Failed to set key", sloki.WrapError(err))
	}

	if err := c.KVSet("system", "collections", "mykey2", 69); err != nil {
		slog.Error("Failed to set key", sloki.WrapError(err))
	}

	if err := c.KVSet("system", "collections", "mykey3", 1233); err != nil {
		slog.Error("Failed to set key", sloki.WrapError(err))
	}

	count, err := c.KVCount("system", "collections")
	if err != nil {
		slog.Error("Failed to count keys", sloki.WrapError(err))
	} else {
		fmt.Printf("Key count: %d\n", count)
	}

	if err := c.KVDeleteAll("system", "collections"); err != nil {
		slog.Error("Failed to delete all keys", sloki.WrapError(err))
	} else {
		fmt.Println("All keys deleted successfully")
	}

	count, err = c.KVCount("system", "collections")
	if err != nil {
		slog.Error("Failed to count keys after deletion", sloki.WrapError(err))
	} else {
		fmt.Printf("Key count after deletion: %d\n", count)
	}
}
