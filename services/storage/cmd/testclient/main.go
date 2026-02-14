package main

import (
	"log/slog"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/client"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/client/collection"
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

	coll, err := collection.NewKeyValueCollection(c, "system", "kv_test")
	if err != nil {
		slog.Error("Failed to create collection", sloki.WrapError(err))
		return
	}

	coll.Set("foo", []byte("bar"))
	coll.Set("foo2", []byte("bar"))
	coll.Set("foo3", []byte("bar"))

}
