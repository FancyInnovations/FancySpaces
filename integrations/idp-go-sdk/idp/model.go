package idp

import "time"

// User represents a user in the identity provider system.
// The ID and Email fields are unique identifiers.
type User struct {
	ID        string            `json:"id" bson:"id"`
	Provider  Provider          `json:"provider" bson:"provider"`
	Name      string            `json:"name" bson:"name"`
	Email     string            `json:"email" bson:"email"`
	Verified  bool              `json:"verified" bson:"verified"`
	Password  string            `json:"password" bson:"password"`
	Roles     []string          `json:"roles" bson:"roles"` // e.g., ["admin", "user"]
	CreatedAt time.Time         `json:"created_at" bson:"created_at"`
	IsActive  bool              `json:"is_active" bson:"is_active"`
	Metadata  map[string]string `json:"metadata" bson:"metadata"` // Additional user metadata
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

// IsUserValid checks if the user is valid for authentication or authorization purposes.
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
