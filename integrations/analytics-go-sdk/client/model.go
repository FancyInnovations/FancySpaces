package client

import "time"

type Event struct {
	Name       string            `json:"name"`
	Properties map[string]string `json:"properties"`
}

func (e *Event) WithProperty(key, value string) {
	if e.Properties == nil {
		e.Properties = make(map[string]string)
	}

	e.Properties[key] = value
}

type createEventDTO struct {
	ProjectID  string            `json:"project_id"`
	Name       string            `json:"name"`
	Timestamp  time.Time         `json:"timestamp,omitempty"`
	Properties map[string]string `json:"properties"`
	WriteKey   string            `json:"write_key,omitempty"`
}

type RecordData struct {
	Metric string  `json:"metric"`
	Label  string  `json:"label,omitempty"`
	Value  float64 `json:"value"`
}

type createMetricRecordDto struct {
	SenderID  string       `json:"sender_id"`
	ProjectID string       `json:"project_id"`
	Timestamp int64        `json:"timestamp,omitempty"`
	WriteKey  string       `json:"write_key,omitempty"`
	Data      []RecordData `json:"data"`
}
