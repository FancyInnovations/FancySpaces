package database

import "time"

type Database struct {
	Name      string                     `json:"name"`
	CreatedAt time.Time                  `json:"created_at"`
	Users     map[string]PermissionLevel `json:"users"`
}

type PermissionLevel string

const (
	PermissionLevelReadOnly  PermissionLevel = "read_only"
	PermissionLevelReadWrite PermissionLevel = "read_write"
	PermissionLevelAdmin     PermissionLevel = "admin"
)

func (db *Database) HasPermission(userID string, requiredLevel PermissionLevel) bool {
	level, exists := db.Users[userID]
	if !exists {
		return false
	}

	switch requiredLevel {
	case PermissionLevelReadOnly:
		return level == PermissionLevelReadOnly || level == PermissionLevelReadWrite || level == PermissionLevelAdmin
	case PermissionLevelReadWrite:
		return level == PermissionLevelReadWrite || level == PermissionLevelAdmin
	case PermissionLevelAdmin:
		return level == PermissionLevelAdmin
	default:
		return false
	}
}

type Collection struct {
	Database   string      `json:"database"`
	Name       string      `json:"name"`
	CreatedAt  time.Time   `json:"created_at"`
	Engine     Engine      `json:"engine"`
	KVSettings *KVSettings `json:"kv_settings,omitempty"`
}

type KVSettings struct {
	DisableTTL bool `json:"disable_ttl"`
}

type Engine string

const (
	EngineKeyValue Engine = "kv"
	EngineObject   Engine = "object"
	EngineBroker   Engine = "broker"
)
