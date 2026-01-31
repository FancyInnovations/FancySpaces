package issues

import (
	"time"
)

type Issue struct {
	ID               string                 `json:"id" bson:"id"`
	Space            string                 `json:"space" bson:"space"`
	Title            string                 `json:"title" bson:"title"`
	Description      string                 `json:"description" bson:"description"`
	Type             Type                   `json:"type" bson:"type"`
	Status           Status                 `json:"status" bson:"status"`
	Priority         Priority               `json:"priority" bson:"priority"`
	Assignee         string                 `json:"assignee,omitempty" bson:"assignee"`
	Reporter         string                 `json:"reporter" bson:"reporter"`
	CreatedAt        time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at" bson:"updated_at"`
	ExternalSource   string                 `json:"external_source,omitempty" bson:"external_source"`
	FixVersion       string                 `json:"fix_version,omitempty" json:"fix_version"`
	AffectedVersions []string               `json:"affected_versions,omitempty" bson:"affected_versions"`
	ResolvedAt       *time.Time             `json:"resolved_at,omitempty" bson:"resolved_at"`
	ParentIssue      string                 `json:"parent_issue,omitempty" bson:"parent_issue"`
	ExtraFields      map[string]interface{} `json:"extra_fields,omitempty" bson:"extra_fields"`
}

type Type string

const (
	TypeEpic  Type = "epic"
	TypeBug   Type = "bug"
	TypeTask  Type = "task"
	TypeStory Type = "story"
	TypeIdea  Type = "idea"
)

type Status string

const (
	StatusBacklog    Status = "backlog"
	StatusToDo       Status = "todo"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
	StatusClosed     Status = "closed"
)

type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityMedium   Priority = "medium"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

type ExternalSource string

const (
	ExternalSourceGitHub           ExternalSource = "github"
	ExternalSourceDiscordForumPost ExternalSource = "discord_forum_post"
	ExternalSourceDiscordTicketBot ExternalSource = "discord_ticket_bot"
)

func (i *Issue) Validate() error {
	if len(i.Title) == 0 {
		return ErrTitleTooShort
	}
	if len(i.Title) > 100 {
		return ErrTitleTooLong
	}

	if len(i.Description) > 1000 {
		return ErrDescriptionTooLong
	}

	return nil
}

type Comment struct {
	ID        string    `json:"id" bson:"id"`
	Issue     string    `json:"issue" bson:"issue"`
	Author    string    `json:"author"`
	Content   string    `json:"content" bson:"content"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
