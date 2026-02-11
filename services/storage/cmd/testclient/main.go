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
	defer c.Close()

	//coll := collection.NewObjectCollection(c, "system", "obj_test")

	db, err := c.DBDatabaseGet("system")
	if err != nil {
		slog.Error("Failed to get database", sloki.WrapError(err))
		return
	}

	fmt.Printf("Database: %#v\n", db)

	coll, err := c.DBCollectionGet("system", "obj_test")
	if err != nil {
		slog.Error("Failed to get collection", sloki.WrapError(err))
		return
	}

	fmt.Printf("Collection: %#v\n", coll)
}
