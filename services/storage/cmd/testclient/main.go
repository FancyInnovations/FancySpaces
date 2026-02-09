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

	//data, err := c.ObjGet("system", "objtest", "something.md")
	//if err != nil {
	//	slog.Error("Failed to get object", sloki.WrapError(err))
	//	return
	//}
	//slog.Info("Object retrieved successfully", slog.String("data", string(data)))

	if err := c.ObjDelete("system", "objtest", "something.md"); err != nil {
		slog.Error("Failed to delete object", sloki.WrapError(err))
		return
	}
	slog.Info("Object deleted successfully")

	md, err := c.ObjGetMetadata("system", "objtest", "something.md")
	if err != nil {
		slog.Error("Failed to get object metadata", sloki.WrapError(err))
		return
	}
	slog.Info("Object metadata retrieved successfully", slog.Uint64("size", uint64(md.Size)), slog.Uint64("crc32", uint64(md.Checksum)))
}
