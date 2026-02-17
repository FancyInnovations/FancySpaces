package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/OliverSchlueter/goutils/broker"
	"github.com/OliverSchlueter/goutils/middleware"
	"github.com/OliverSchlueter/goutils/problems"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/core/internal/spaces"
	spaces2 "github.com/fancyinnovations/fancyspaces/integrations/spaces-go-sdk/spaces"
	"github.com/nats-io/nats.go"
)

type NatsHandler struct {
	broker broker.Broker
	store  *spaces.Store
}

type NatsConfiguration struct {
	Broker broker.Broker
	Store  *spaces.Store
}

func NewNatsHandler(cfg NatsConfiguration) *NatsHandler {
	return &NatsHandler{
		broker: cfg.Broker,
		store:  cfg.Store,
	}
}

func (h *NatsHandler) Register() error {
	if err := h.broker.SubscribeQueue("fancyspaces.core.spaces.get", "fancyspaces.core.spaces.get", middleware.NatsLogging(h.handleGet)); err != nil {
		return fmt.Errorf("could not subscribe to nats subject: %w", err)
	}

	return nil
}

func (h *NatsHandler) handleGet(msg *nats.Msg) {
	id := string(msg.Data)

	if id == "" {
		problems.ValidationError("space_id", "Space ID is required").WriteToBroker(h.broker, msg.Reply)
		return
	}

	space, err := h.store.GetByID(id)
	if err != nil {
		if errors.Is(err, spaces.ErrSpaceNotFound) {
			problems.NotFound("Space", id).WriteToBroker(h.broker, msg.Reply)
			return
		}
		slog.Error("Could not get space", sloki.WrapError(err))
		problems.InternalServerError("").WriteToBroker(h.broker, msg.Reply)
		return
	}

	response, err := json.Marshal(spaces2.InternalSpace{
		Space:             *space,
		AnalyticsWriteKey: space.AnalyticsSettings.WriteKey,
	})
	if err != nil {
		slog.Error("failed to marshal space", sloki.WrapError(err))
		problems.InternalServerError("Could not marshal space").WriteToBroker(h.broker, msg.Reply)
		return
	}

	if err := h.broker.Publish(msg.Reply, response); err != nil {
		slog.Error("failed to publish space response", sloki.WrapError(err))
	}
}
