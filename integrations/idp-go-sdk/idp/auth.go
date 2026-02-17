package idp

import (
	"crypto/rsa"
	"encoding/json"

	"github.com/OliverSchlueter/goutils/broker"
)

var ServiceBaseURL = "https://fancyanalytics.net/idp/api/v1"

// Service provides methods to interact with the IDP service.
type Service struct {
	broker         broker.Broker
	excludedRoutes []string
	publicKey      *rsa.PublicKey
	usersCache     *usersCache
}

// Configuration holds the necessary configuration for initializing the IDP service.
type Configuration struct {
	// Broker is the message broker used for communication with the IDP service.
	Broker broker.Broker

	// PublicKey is the RSA public key used for validating JWT tokens issued by the IDP service.
	PublicKey *rsa.PublicKey

	// ExcludedRoutes is a list of (HTTP) routes that should be excluded from authentication checks. (optional)
	ExcludedRoutes []string
}

// NewService initializes and returns a new instance of the IDP service with the provided configuration.
func NewService(cfg Configuration) *Service {
	if cfg.ExcludedRoutes == nil {
		cfg.ExcludedRoutes = []string{}
	}

	return &Service{
		broker:         cfg.Broker,
		excludedRoutes: cfg.ExcludedRoutes,
		publicKey:      cfg.PublicKey,
		usersCache:     newUsersCache(),
	}
}

// GetUser retrieves a user by their ID or email.
func (s *Service) GetUser(id string) (*User, error) {
	userFromCache, err := s.usersCache.GetByID(id)
	if err == nil {
		return userFromCache, nil
	}

	resp, err := s.broker.Request("idp.user.get", []byte(id))
	if err != nil {
		return nil, err
	}

	var u User
	if err := json.Unmarshal(resp.Data, &u); err != nil {
		return nil, err
	}

	s.usersCache.UpsertUser(&u)

	return &u, nil
}

// ValidateUser validates a user's credentials and returns the user if valid.
func (s *Service) ValidateUser(userID, password string) (*User, error) {
	userFromCache, err := s.usersCache.GetByID(userID)
	if err == nil {
		// user found in cache, validate password
		if userFromCache.Password != PasswordHash(password) {
			return nil, ErrInvalidBasicCredentials
		}
		return userFromCache, nil
	}

	resp, err := s.broker.Request("idp.user.validate", []byte(`{"username":"`+userID+`", "password":"`+password+`"}`))
	if err != nil {
		return nil, err
	}

	var u User
	if err := json.Unmarshal(resp.Data, &u); err != nil {
		return nil, err
	}

	s.usersCache.UpsertUser(&u)

	return &u, nil
}
