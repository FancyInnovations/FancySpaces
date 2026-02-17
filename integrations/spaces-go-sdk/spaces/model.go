package spaces

import (
	"time"

	"github.com/fancyinnovations/fancyspaces/integrations/idp-go-sdk/idp"
)

type Space struct {
	ID          string     `json:"id"`
	Slug        string     `json:"slug"`
	Title       string     `json:"title"`
	Summary     string     `json:"summary"`
	Description string     `json:"description"`
	Categories  []Category `json:"categories"`
	Links       []Link     `json:"links"`
	IconURL     string     `json:"icon_url"`
	Status      Status     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	Creator     string     `json:"creator"` // UserID of the creator
	Members     []Member   `json:"members"`

	IssueSettings           IssueSettings           `json:"issue_settings"`
	ReleaseSettings         ReleaseSettings         `json:"release_settings"`
	MavenRepositorySettings MavenRepositorySettings `json:"maven_repository_settings"`
	StorageSettings         StorageSettings         `json:"storage_settings"`
	AnalyticsSettings       AnalyticsSettings       `json:"analytics_settings"`
}

// InternalSpace is the internal representation of a Space, containing additional fields that are not exposed to clients.
type InternalSpace struct {
	Space

	AnalyticsWriteKey string `json:"analytics_write_key"`
}

type IssueSettings struct {
	Enabled bool `json:"enabled"`

	GitHubSync      bool   `json:"github_sync"`
	GitHubSyncOwner string `json:"github_sync_owner"` // GitHub username or organization (required if GitHubSync is true)
	GitHubSyncRepo  string `json:"github_sync_repo"`  // GitHub repository name (required if GitHubSync is true)
	GitHubSyncLabel string `json:"github_sync_label"` // only sync issues with this label (optional)
}

type ReleaseSettings struct {
	Enabled bool `json:"enabled"`

	// TODO: Implement these settings in the future
	//DiscordNotifications          bool   `json:"discord_notifications"`
	//DiscordNotificationWebhookURL string `json:"discord_notification_webhook_url"`
}

type MavenRepositorySettings struct {
	Enabled bool `json:"enabled"`
}

type StorageSettings struct {
	Enabled bool `json:"enabled"`
}

type AnalyticsSettings struct {
	Enabled         bool   `json:"enabled"`
	RequireWriteKey bool   `json:"require_write_key"`
	WriteKey        string `json:"-"`
}

type Link struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Member struct {
	UserID string `json:"user_id"`
	Role   Role   `json:"role"`
}

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

type Status string

const (
	StatusDraft    Status = "draft"
	StatusReview   Status = "review"
	StatusApproved Status = "approved"
	StatusPrivate  Status = "private"
	StatusArchived Status = "archived"
	StatusRejected Status = "rejected"
	StatusBanned   Status = "banned"
)

type Category string

const (
	CategoryMinecraftPlugin Category = "minecraft_plugin"
	CategoryMinecraftServer Category = "minecraft_server"
	CategoryMinecraftMod    Category = "minecraft_mod"
	CategoryHytalePlugin    Category = "hytale_plugin"
	CategoryWebApp          Category = "web_app"
	CategoryMobileApp       Category = "mobile_app"
	CategoryOther           Category = "other"
)

func (s *Space) IsMember(u *idp.User) bool {
	if !idp.IsUserValid(u) {
		return false
	}

	if s.Creator == u.ID {
		return true
	}

	for _, m := range s.Members {
		if m.UserID == u.ID {
			return true
		}
	}

	return false
}

func (s *Space) IsOwner(u *idp.User) bool {
	if !idp.IsUserValid(u) {
		return false
	}

	return s.Creator == u.ID
}

func (s *Space) HasFullAccess(u *idp.User) bool {
	if !idp.IsUserValid(u) {
		return false
	}

	if s.Creator == u.ID {
		return true
	}

	for _, m := range s.Members {
		if m.UserID == u.ID {
			return m.Role == RoleAdmin
		}
	}

	return false
}

func (s *Space) HasWriteAccess(u *idp.User) bool {
	if !idp.IsUserValid(u) {
		return false
	}

	if s.Creator == u.ID {
		return true
	}

	for _, m := range s.Members {
		if m.UserID == u.ID {
			return m.Role == RoleAdmin || m.Role == RoleMember
		}
	}

	return false
}

func (s *Space) Validate() error {
	if len(s.Slug) < 3 {
		return ErrSlugTooShort
	}
	if len(s.Slug) > 20 {
		return ErrSlugTooLong
	}

	if len(s.Title) > 100 {
		return ErrTitleTooLong
	}
	if len(s.Title) < 3 {
		return ErrTitleTooShort
	}

	if len(s.Description) > 500 {
		return ErrDescriptionTooLong
	}

	return nil
}
