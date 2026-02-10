package main

import (
	"log/slog"
	"time"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/client"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/codex"
)

type Sample struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Something int       `json:"something"`
}

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

	s := Sample{
		ID:        "sample1",
		Name:      "Test Sample",
		CreatedAt: time.Now(),
		Something: 42,
	}

	data, err := codex.Marshal(s)
	if err != nil {
		slog.Error("Failed to marshal sample", sloki.WrapError(err))
		return
	}

	if err := c.KVSet("system", "collections", s.ID, data); err != nil {
		slog.Error("Failed to set key", sloki.WrapError(err))
		return
	}
	slog.Info("Successfully set key")

	retrievedData, err := c.KVGet("system", "collections", s.ID)
	if err != nil {
		slog.Error("Failed to get key", sloki.WrapError(err))
		return
	}
	out := &Sample{}
	if err := codex.Unmarshal(retrievedData.AsBinary(), out); err != nil {
		slog.Error("Failed to unmarshal data", sloki.WrapError(err))
		return
	}
	slog.Info("Successfully retrieved and unmarshaled key", slog.Any("data", out))

}
