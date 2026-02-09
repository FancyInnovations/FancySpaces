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
		return
	}
	defer c.Close()

	//if err := c.ObjPut("system", "objtest", "something.md", []byte("Hello world!lolol")); err != nil {
	//	slog.Error("Failed to put object", sloki.WrapError(err))
	//	return
	//}
	//slog.Info("Object put successfully")

	data, err := c.ObjGet("system", "objtest", "something.md")
	if err != nil {
		slog.Error("Failed to get object", sloki.WrapError(err))
		return
	}
	slog.Info("Object retrieved successfully", slog.String("data", string(data)))
}
