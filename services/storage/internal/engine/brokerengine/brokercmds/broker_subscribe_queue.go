package brokercmds

import (
	"encoding/binary"
	"errors"
	"log/slog"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/storage/internal/auth"
	"github.com/fancyinnovations/fancyspaces/storage/internal/command"
	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
	"github.com/fancyinnovations/fancyspaces/storage/internal/engine/brokerengine"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/commonresponses"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/protocol"
)

// handleSubscribeQueue handles the protocol.ServerCommandBrokerSubscribeQueue command, which subscribes the client to a given subject on the broker engine with a queue group.
// Payload format: | subject Length (2 bytes) | subject (variable) | queue group Length (2 bytes) | queue group (variable) |
func (c *Commands) handleSubscribeQueue(ctx *command.ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
	u := auth.UserFromContext(ctx.Ctx)
	if u == nil || !u.Verified || !u.IsActive {
		return commonresponses.Unauthorized, nil
	}

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
	if len(data) < 4 {
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

	if len(data) < 2+subjectLen+2 {
		return &protocol.Response{
			Code:    protocol.StatusBadRequest,
			Payload: []byte("invalid payload length for queue group length"),
		}, nil
	}

	queueGroupLen := int(binary.BigEndian.Uint16(data[2+subjectLen : 2+subjectLen+2]))
	if len(data) < 2+subjectLen+2+queueGroupLen {
		return &protocol.Response{
			Code:    protocol.StatusBadRequest,
			Payload: []byte("invalid payload length for queue group"),
		}, nil
	}

	queueGroup := string(data[2+subjectLen+2 : 2+subjectLen+2+queueGroupLen])

	subject := string(data[2 : 2+subjectLen])

	be.Subscribe(subject, &brokerengine.Subscriber{
		ID:    ctx.ID,
		Queue: queueGroup,
	})

	return commonresponses.OK, nil
}
