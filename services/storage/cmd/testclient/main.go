package main

import (
	"fmt"
	"log/slog"
	"time"

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
	defer c.Close()

	metadata, err := c.ObjGetMetadata("system", "objtest", "something.md")
	if err != nil {
		slog.Error("Failed to get metadata", sloki.WrapError(err))
		return
	}

	fmt.Printf("Metadata: %+v\n", metadata)

	time.Sleep(1 * time.Second)

	if err := c.ObjPut("system", "objtest", "something.md", []byte("Hello, World!!1!!!!!!!")); err != nil {
		slog.Error("Failed to put object", sloki.WrapError(err))
		return
	}
	slog.Info("Object put successfully")

	metadata, err = c.ObjGetMetadata("system", "objtest", "something.md")
	if err != nil {
		slog.Error("Failed to get metadata after put", sloki.WrapError(err))
		return
	}

	fmt.Printf("Metadata after put: %+v\n", metadata)

}
