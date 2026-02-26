package spaces

import (
	"regexp"
	"time"

	"github.com/fancyinnovations/fancyspaces/integrations/idp-go-sdk/idp"
)

type Space struct {
	ID          string     `json:"id" bson:"space_id"`
	Slug        string     `json:"slug" bson:"slug"`
	Title       string     `json:"title" bson:"title"`
	Summary     string     `json:"summary" bson:"summary"`
	Description string     `json:"description" bson:"description"`
	Categories  []Category `json:"categories" bson:"categories"`
	Links       []Link     `json:"links" bson:"links"`
	IconURL     string     `json:"icon_url" bson:"icon_url"`
	Status      Status     `json:"status" bson:"status"`
	CreatedAt   time.Time  `json:"created_at" bson:"created_at"`
	Creator     string     `json:"creator" bson:"creator"` // UserID of the creator
	Members     []Member   `json:"members" bson:"members"`

	IssueSettings           IssueSettings           `json:"issue_settings" bson:"issue_settings"`
	ReleaseSettings         ReleaseSettings         `json:"release_settings" bson:"release_settings"`
	MavenRepositorySettings MavenRepositorySettings `json:"maven_repository_settings" bson:"maven_repository_settings"`
	StorageSettings         StorageSettings         `json:"storage_settings" bson:"storage_settings"`
	AnalyticsSettings       AnalyticsSettings       `json:"analytics_settings" bson:"analytics_settings"`
	SecretsSettings         SecretsSettings         `json:"secrets_settings" bson:"secrets_settings"`
	BlogSettings            BlogSettings            `json:"blog_settings" bson:"blog_settings"`
}

// InternalSpace is the internal representation of a Space, containing additional fields that are not exposed to clients.
type InternalSpace struct {
	Space

	AnalyticsWriteKey string `json:"analytics_write_key"`
}

type IssueSettings struct {
	Enabled bool `json:"enabled" bson:"enabled"`

	GitHubSync      bool   `json:"github_sync" bson:"github_sync"`
	GitHubSyncOwner string `json:"github_sync_owner" bson:"github_sync_owner"` // GitHub username or organization (required if GitHubSync is true)
	GitHubSyncRepo  string `json:"github_sync_repo" bson:"github_sync_repo"`   // GitHub repository name (required if GitHubSync is true)
	GitHubSyncLabel string `json:"github_sync_label" bson:"github_sync_label"` // only sync issues with this label (optional)
}

type ReleaseSettings struct {
	Enabled bool `json:"enabled" bson:"enabled"`

	// TODO: Implement these settings in the future
	//DiscordNotifications          bool   `json:"discord_notifications"`
	//DiscordNotificationWebhookURL string `json:"discord_notification_webhook_url"`
}

type MavenRepositorySettings struct {
	Enabled bool `json:"enabled" bson:"enabled"`
}

type StorageSettings struct {
	Enabled bool `json:"enabled" bson:"enabled"`
}

type AnalyticsSettings struct {
	Enabled         bool   `json:"enabled" bson:"enabled"`
	RequireWriteKey bool   `json:"require_write_key" bson:"require_write_key"`
	WriteKey        string `json:"-" bson:"write_key"`
}

type SecretsSettings struct {
	Enabled bool `json:"enabled" bson:"enabled"`
}

type BlogSettings struct {
	Enabled bool `json:"enabled" bson:"enabled"`
}

type Link struct {
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

type Member struct {
	UserID string `json:"user_id" bson:"user_id"`
	Role   Role   `json:"role" bson:"role"`
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

	slugPattern := `^[a-z0-9]+(?:-[a-z0-9]+)*$`
	matched, err := regexp.MatchString(slugPattern, s.Slug)
	if err != nil {
		return err
	}
	if !matched {
		return ErrSlugInvalidFormat
	}

	if len(s.Title) > 100 {
		return ErrTitleTooLong
	}
	if len(s.Title) < 3 {
		return ErrTitleTooShort
	}

	if len(s.Summary) > 300 {
		return ErrSummaryTooLong
	}

	if len(s.Description) > 10_000 {
		return ErrDescriptionTooLong
	}

	return nil
}
