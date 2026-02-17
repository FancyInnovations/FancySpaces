package handler

import (
	"encoding/json"
	"fmt"

	"github.com/OliverSchlueter/goutils/broker"
	"github.com/OliverSchlueter/goutils/middleware"
	"github.com/fancyinnovations/fancyspaces/core/internal/secrets"
	"github.com/fancyinnovations/fancyspaces/integrations/spaces-go-sdk/spaces"
	"github.com/nats-io/nats.go"
)

type NatsHandler struct {
	broker broker.Broker
	store  *secrets.Store
}

type NatsConfiguration struct {
	Broker broker.Broker
	Store  *secrets.Store
}

func NewNatsHandler(cfg NatsConfiguration) *NatsHandler {
	return &NatsHandler{
		broker: cfg.Broker,
		store:  cfg.Store,
	}
}

func (h *NatsHandler) Register() error {
	if err := h.broker.SubscribeQueue("fancyspaces.core.secrets.get", "fancyspaces.core.secrets.get", middleware.NatsLogging(h.handleGet)); err != nil {
		return fmt.Errorf("could not subscribe to nats subject: %w", err)
	}

	return nil
}

func (h *NatsHandler) handleGet(msg *nats.Msg) {
	var req spaces.GetSecretReqDTO
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		h.broker.Publish(msg.Reply, []byte(fmt.Sprintf("failed to unmarshal request: %v", err)))
		return
	}

	secret, err := h.store.GetDecryptedSecret(req.SpaceID, req.Key)
	if err != nil {
		h.broker.Publish(msg.Reply, []byte(fmt.Sprintf("failed to get secret: %v", err)))
		return
	}

	resp := spaces.GetSecretRespDTO{
		Value: secret,
	}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		h.broker.Publish(msg.Reply, []byte(fmt.Sprintf("failed to marshal response: %v", err)))
		return
	}

	h.broker.Publish(msg.Reply, respBytes)
}
