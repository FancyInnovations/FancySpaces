package spaces

import (
	"time"

	"github.com/fancyinnovations/fancyspaces/internal/auth"
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
	Members     []Member   `json:"members"`

	IssueSettings           IssueSettings           `json:"issue_settings"`
	ReleaseSettings         ReleaseSettings         `json:"release_settings"`
	MavenRepositorySettings MavenRepositorySettings `json:"maven_repository_settings"`
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
	RoleOwner  Role = "owner"
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

func (s *Space) IsMember(u *auth.User) bool {
	for _, m := range s.Members {
		if m.UserID == u.ID {
			return true
		}
	}

	return false
}

func (s *Space) IsOwner(u *auth.User) bool {
	for _, m := range s.Members {
		if m.UserID == u.ID {
			return m.Role == RoleOwner
		}
	}

	return false
}

func (s *Space) HasFullAccess(u *auth.User) bool {
	for _, m := range s.Members {
		if m.UserID == u.ID {
			return m.Role == RoleOwner || m.Role == RoleAdmin
		}
	}

	return false
}

func (s *Space) HasWriteAccess(u *auth.User) bool {
	for _, m := range s.Members {
		if m.UserID == u.ID {
			return m.Role == RoleOwner || m.Role == RoleAdmin || m.Role == RoleMember
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

type CreateOrUpdateSpaceReq struct {
	Slug        string     `json:"slug"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Categories  []Category `json:"categories"`
	IconURL     string     `json:"icon_url"`
}
