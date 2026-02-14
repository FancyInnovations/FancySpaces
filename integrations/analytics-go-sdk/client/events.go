package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (c *Client) SendEvent(evt *Event) error {
	dto := createEventDTO{
		ProjectID:  c.projectID,
		Name:       evt.Name,
		Timestamp:  time.Now(),
		Properties: evt.Properties,
		WriteKey:   c.writeKey,
	}

	body, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("could not marshal event: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/collector/api/v1/events", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.authToken)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrUnexpectedStatusCode
	}

	return nil
}
