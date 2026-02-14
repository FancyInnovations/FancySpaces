package brokercmds

import (
	"encoding/binary"
	"errors"
	"log/slog"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/commonresponses"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/protocol"
	"github.com/fancyinnovations/fancyspaces/storage/internal/command"
	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
	"github.com/fancyinnovations/fancyspaces/storage/internal/engine/brokerengine"
)

// handleSubscribe handles the protocol.ServerCommandBrokerSubscribe command, which subscribes the client to a given subject on the broker engine.
// Payload format: | subject Length (2 bytes) | subject (variable) |
func (c *Commands) handleSubscribe(ctx *command.ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
	e, err := c.engineService.GetEngine(cmd.DatabaseName, cmd.CollectionName)
	if err != nil {
		if errors.Is(err, database.ErrCollectionNotFound) {
			return commonresponses.CollectionNotFound, nil
		}

		slog.Error("Failed to get engine",
			slog.String("database", cmd.DatabaseName),
			slog.String("collection", cmd.CollectionName),
			sloki.WrapError(err),
		)
		return commonresponses.InternalServerError, nil
	}

	if e.Type != database.EngineBroker {
		return commonresponses.CommandNotAllowed, nil
	}
	be := e.AsBrokerEngine()

	data := cmd.Payload
	if len(data) < 2 {
		return &protocol.Response{
			Code:    protocol.StatusBadRequest,
			Payload: []byte("invalid payload length"),
		}, nil
	}

	subjectLen := int(binary.BigEndian.Uint16(data[0:2]))
	if len(data) < 2+subjectLen {
		return &protocol.Response{
			Code:    protocol.StatusBadRequest,
			Payload: []byte("invalid payload length for subject"),
		}, nil
	}

	subject := string(data[2 : 2+subjectLen])

	be.Subscribe(subject, &brokerengine.Subscriber{
		ID:    ctx.ID,
		Queue: "",
	})

	return commonresponses.OK, nil
}
