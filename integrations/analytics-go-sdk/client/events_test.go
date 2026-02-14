package client

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_SendEvent(t *testing.T) {
	for _, tc := range []struct {
		name        string
		statusCode  int
		serverError bool
		expErr      error
	}{
		{
			name:       "successful event send",
			statusCode: http.StatusOK,
			expErr:     nil,
		},
		{
			name:       "unexpected status code",
			statusCode: http.StatusBadRequest,
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
				server = httptest.NewServer(http.NotFoundHandler())
				server.Close()
			} else {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tc.statusCode)
				}))
			}

			c, err := New(Configuration{
				BaseURL:   server.URL,
				AuthToken: "token",
				ProjectID: "project",
			})
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			e := &Event{
				Name: "test_event",
				Properties: map[string]string{
					"key": "value",
				},
			}

			err = c.SendEvent(e)
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
