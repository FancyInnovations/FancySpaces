package main

import (
	"log/slog"
	"strconv"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/client"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/client/collection"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/codex"
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

	coll := collection.NewObjectCollection(c, "system", "obj_test")

	for i := 0; i < 100; i++ {
		p := Person{
			Name: "Person " + strconv.Itoa(200+i),
			Age:  i,
		}
		data, err := codex.Marshal(&p)
		if err != nil {
			slog.Error("Failed to marshal person", sloki.WrapError(err))
			return
		}

		err = coll.Put("person"+strconv.Itoa(200+i), data)
		if err != nil {
			slog.Error("Failed to set value", sloki.WrapError(err))
			return
		}
	}

	count, err := coll.Count()
	if err != nil {
		slog.Error("Failed to get count", sloki.WrapError(err))
		return
	}
	slog.Info("Count retrieved successfully", slog.Uint64("count", uint64(count)))

	size, err := coll.Size()
	if err != nil {
		slog.Error("Failed to get size", sloki.WrapError(err))
		return
	}
	slog.Info("Size retrieved successfully", slog.Uint64("size", size))
}
