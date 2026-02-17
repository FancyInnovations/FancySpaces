package idp

import "time"

// User represents a user in the identity provider system.
// The ID and Email fields are unique identifiers.
type User struct {
	ID        string            `json:"id"`
	Provider  Provider          `json:"provider"`
	Name      string            `json:"name"`
	Email     string            `json:"email"`
	Verified  bool              `json:"verified"`
	Password  string            `json:"password"`
	Roles     []string          `json:"roles"` // e.g., ["admin", "user"]
	CreatedAt time.Time         `json:"created_at"`
	IsActive  bool              `json:"is_active"`
	Metadata  map[string]string `json:"metadata"` // Additional user metadata
}

func (u *User) IsAdmin() bool {
	for _, role := range u.Roles {
		if role == "admin" {
			return true
		}
	}

	return false
}

type Provider string

const (
	ProviderBasic   Provider = "basic"
	ProviderGoogle  Provider = "google"
	ProviderGithub  Provider = "github"
	ProviderDiscord Provider = "discord"
)

func IsUserValid(user *User) bool {
	if user == nil {
		return false
	}
	if !user.IsActive {
		return false
	}
	if !user.Verified {
		return false
	}
	return true
}
