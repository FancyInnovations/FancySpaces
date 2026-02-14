// Package client provides a client for the FancyAnalytics API.
package client

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	DefaultBaseURL            = "http://fancyanalytics.net"
	DefaultMaxIdleConnections = 100
	DefaultIdleTimeout        = 90 * time.Second
	DefaultTimeout            = 30 * time.Second
)

var DefaultClient *Client

// Client is a client for the FancyAnalytics API.
type Client struct {
	baseURL   string
	authToken string
	senderID  string
	projectID string
	writeKey  string
	client    *http.Client
}

// Configuration holds the configuration for the FancyAnalytics client.
type Configuration struct {
	// BaseURL is the base URL of the FancyAnalytics API.
	// If empty, DefaultBaseURL is used.
	BaseURL string

	// SenderID is the sender ID to use for requests.
	// If empty a random UUID will be used.
	SenderID string

	// AuthToken is the authentication token to use for requests. Optional.
	AuthToken string

	// ProjectID is the project ID to use for requests.
	ProjectID string

	// WriteKey is the write key for the project. Optional.
	WriteKey string

	// MaxIdleConnections sets the maximum number of idle (keep-alive) connections across all hosts.
	// If zero, DefaultMaxIdleConnections is used.
	MaxIdleConnections int

	// IdleConnectionTimeout is the maximum amount of time an idle (keep-alive) connection will remain idle before closing itself.
	// If zero, DefaultIdleTimeout is used.
	IdleConnectionTimeout time.Duration

	// RequestTimeout is the maximum amount of time a request can take before timing out.
	// If zero, DefaultTimeout is used.
	RequestTimeout time.Duration
}

// New creates a new FancyAnalytics client with the given configuration.
func New(cfg Configuration) (*Client, error) {
	if cfg.ProjectID == "" {
		return nil, ErrNoProjectID
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultBaseURL
	}

	if cfg.SenderID == "" {
		cfg.SenderID = uuid.New().String()
	}

	if cfg.MaxIdleConnections == 0 {
		cfg.MaxIdleConnections = DefaultMaxIdleConnections
	}

	if cfg.IdleConnectionTimeout == 0 {
		cfg.IdleConnectionTimeout = DefaultIdleTimeout
	}

	if cfg.RequestTimeout == 0 {
		cfg.RequestTimeout = DefaultTimeout
	}

	transport := http.Transport{
		MaxIdleConns:       cfg.MaxIdleConnections,
		IdleConnTimeout:    cfg.IdleConnectionTimeout,
		DisableKeepAlives:  false,
		DisableCompression: false,
	}

	c := &http.Client{
		Transport: &transport,
		Timeout:   cfg.RequestTimeout,
	}

	return &Client{
		baseURL:   cfg.BaseURL,
		authToken: cfg.AuthToken,
		projectID: cfg.ProjectID,
		writeKey:  cfg.WriteKey,
		senderID:  cfg.SenderID,
		client:    c,
	}, nil
}

// Ping checks if the FancyAnalytics API is reachable.
func (c *Client) Ping() error {
	resp, err := c.client.Get(c.baseURL + "/collector/api/v1/ping")
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrUnexpectedStatusCode
	}

	return nil
}
