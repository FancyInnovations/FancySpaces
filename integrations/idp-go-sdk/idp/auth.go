package idp

import (
	"crypto/rsa"
	"encoding/json"
	"strings"

	"github.com/OliverSchlueter/goutils/broker"
)

var ServiceBaseURL = "https://fancyanalytics.net/idp/api/v1"

// Service provides methods to interact with the IDP service.
type Service struct {
	broker     broker.Broker
	publicKey  *rsa.PublicKey
	usersCache *usersCache
}

// Configuration holds the necessary configuration for initializing the IDP service.
type Configuration struct {
	// Broker is the message broker used for communication with the IDP service.
	Broker broker.Broker

	// PublicKey is the RSA public key used for validating JWT tokens issued by the IDP service.
	PublicKey *rsa.PublicKey
}

// NewService initializes and returns a new instance of the IDP service with the provided configuration.
func NewService(cfg Configuration) *Service {
	return &Service{
		broker:     cfg.Broker,
		publicKey:  cfg.PublicKey,
		usersCache: newUsersCache(),
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

	if strings.Contains(string(resp.Data), "NotFound") {
		return nil, ErrUserNotFound
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
		if isValid, err := CheckPassword(password, userFromCache.Password); err != nil || !isValid {
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
