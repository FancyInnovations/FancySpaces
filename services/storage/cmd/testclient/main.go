package main

import (
	"log/slog"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/client"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/protocol"
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

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.CommandKVGet,
		DatabaseName:   "system",
		CollectionName: "collections",
		Payload:        make([]byte, 0),
	})
	if err != nil {
		slog.Error("Command failed", sloki.WrapError(err))
	}

	slog.Info("Command response", slog.Int("code", int(resp.Code)), slog.Int("payload_length", len(resp.Payload)))
}
