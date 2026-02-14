package client

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	for _, tc := range []struct {
		name   string
		cfg    Configuration
		exp    *Client
		expErr error
	}{
		{
			name: "default values",
			cfg: Configuration{
				AuthToken: "my-token",
				ProjectID: "my-project",
			},
			expErr: nil,
		},
		{
			name: "all values set",
			cfg: Configuration{
				BaseURL:               "https://custom-url.com/",
				AuthToken:             "my-token",
				ProjectID:             "my-project",
				SenderID:              "custom-sender-id",
				MaxIdleConnections:    50,
				IdleConnectionTimeout: 60 * time.Second,
				RequestTimeout:        15 * time.Second,
			},
			expErr: nil,
		},
		{
			name: "missing project ID",
			cfg: Configuration{
				AuthToken: "my-token",
			},
			expErr: ErrNoProjectID,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := New(tc.cfg)
			if !errors.Is(err, tc.expErr) {
				t.Errorf("expected error %v, got %v", tc.expErr, err)
			}

			// If an error was expected, no need to check further.
			if tc.expErr != nil {
				return
			}

			if tc.cfg.BaseURL == "" {
				if got.baseURL != DefaultBaseURL {
					t.Errorf("expected baseURL %s, got %s", DefaultBaseURL, got.baseURL)
				}
			} else if got.baseURL != tc.cfg.BaseURL {
				t.Errorf("expected baseURL %s, got %s", tc.cfg.BaseURL, got.baseURL)
			}

			if got.authToken != tc.cfg.AuthToken {
				t.Errorf("expected authToken %s, got %s", tc.cfg.AuthToken, got.authToken)
			}

			if got.projectID != tc.cfg.ProjectID {
				t.Errorf("expected projectID %s, got %s", tc.cfg.ProjectID, got.projectID)
			}

			if tc.cfg.SenderID == "" {
				if got.senderID == "" {
					t.Error("expected non-empty senderID, got empty")
				}
			} else if got.senderID != tc.cfg.SenderID {
				t.Errorf("expected senderID %s, got %s", tc.cfg.SenderID, got.senderID)
			}

			if tc.cfg.MaxIdleConnections == 0 {
				if got.client.Transport.(*http.Transport).MaxIdleConns != DefaultMaxIdleConnections {
					t.Errorf("expected MaxIdleConnections %d, got %d", DefaultMaxIdleConnections, got.client.Transport.(*http.Transport).MaxIdleConns)
				}
			} else if got.client.Transport.(*http.Transport).MaxIdleConns != tc.cfg.MaxIdleConnections {
				t.Errorf("expected MaxIdleConnections %d, got %d", tc.cfg.MaxIdleConnections, got.client.Transport.(*http.Transport).MaxIdleConns)
			}
		})
	}
}

func TestClient_Ping(t *testing.T) {
	for _, tc := range []struct {
		name        string
		statusCode  int
		serverError bool
		expErr      error
	}{
		{
			name:       "successful ping",
			statusCode: http.StatusOK,
			expErr:     nil,
		},
		{
			name:       "unexpected status code",
			statusCode: http.StatusInternalServerError,
			expErr:     ErrUnexpectedStatusCode,
		},
		{
			name:        "network error",
			serverError: true,
			expErr:      errors.New("network error"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var server *httptest.Server
			if tc.serverError {
				// Simulate network error by using a closed server
				server = httptest.NewServer(http.NotFoundHandler())
				server.Close()
			} else {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tc.statusCode)
				}))
			}

			client, err := New(Configuration{
				BaseURL:   server.URL,
				AuthToken: "test-token",
				ProjectID: "test-project",
			})
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			err = client.Ping()
			if tc.expErr == nil {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			} else if tc.serverError {
				if err == nil {
					t.Errorf("expected network error, got nil")
				}
			} else if !errors.Is(err, tc.expErr) {
				t.Errorf("expected error %v, got %v", tc.expErr, err)
			}
		})
	}
}
