package spaces

import (
	"time"

	"github.com/fancyinnovations/fancyspaces/src/internal/auth"
)

type Space struct {
	ID          string     `json:"id"`
	Slug        string     `json:"slug"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Categories  []Category `json:"categories"`
	IconURL     string     `json:"icon_url"`
	Status      Status     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	Members     []Member   `json:"members"`
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
	StatusArchived Status = "archived"
	StatusRejected Status = "rejected"
	StatusBanned   Status = "banned"
)

type Category string

const (
	CategoryMinecraftPlugin Category = "minecraft_plugin"
	CategoryMinecraftServer Category = "minecraft_server"
	CategoryMinecraftMod    Category = "minecraft_mod"
	CategoryWebApp          Category = "web_app"
	CategoryMobileApp       Category = "mobile_app"
	CategoryOther           Category = "other"
)

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
