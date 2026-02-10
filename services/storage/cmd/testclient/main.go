package main

import (
	"log/slog"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/client"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/client/collection"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
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

	coll := collection.NewMessageBrokerCollection(c, "system", "kv_test")
	_ = coll
}
