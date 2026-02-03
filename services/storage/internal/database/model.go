package database

import "time"

type Database struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Collection struct {
	Database  string    `json:"database"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Engine    Engine    `json:"engine"`
}

type Engine string

const (
	EngineKeyValue Engine = "kv"
)
