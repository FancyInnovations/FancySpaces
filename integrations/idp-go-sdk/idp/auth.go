package idp

import (
	"crypto/rsa"

	"github.com/OliverSchlueter/goutils/broker"
)

type Service struct {
	broker         broker.Broker
	excludedRoutes []string
	publicKey      *rsa.PublicKey
	usersCache     *usersCache
}

type Configuration struct {
	PublicKey      *rsa.PublicKey
	Broker         broker.Broker
	ExcludedRoutes []string
}

func NewService(cfg Configuration) *Service {
	return &Service{
		broker:         cfg.Broker,
		excludedRoutes: cfg.ExcludedRoutes,
		publicKey:      cfg.PublicKey,
		usersCache:     newUsersCache(),
	}
}
