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

	//fn := func(msg []byte) {
	//	fmt.Printf("Received message: %s\n", string(msg))
	//}
	//
	//if err := c.BrokerSubscribe("system", "brokertest", "user.created", fn); err != nil {
	//	slog.Error("Failed to subscribe", sloki.WrapError(err))
	//	return
	//}
	//
	//time.Sleep(5 * time.Second)
	if err := c.BrokerPublish("system", "brokertest", "user.created", []byte("Hello, World!")); err != nil {
		slog.Error("Failed to publish", sloki.WrapError(err))
		return
	}

	select {}
}
