package issues

import (
	"strings"
	"time"
)

type Issue struct {
	ID               string         `json:"id" bson:"id"`
	Space            string         `json:"space" bson:"space"`
	Title            string         `json:"title" bson:"title"`
	Description      string         `json:"description" bson:"description"`
	Type             Type           `json:"type" bson:"type"`
	Status           Status         `json:"status" bson:"status"`
	Priority         Priority       `json:"priority" bson:"priority"`
	Assignee         string         `json:"assignee,omitempty" bson:"assignee"`
	Reporter         string         `json:"reporter" bson:"reporter"`
	CreatedAt        time.Time      `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" bson:"updated_at"`
	ExternalSource   ExternalSource `json:"external_source,omitempty" bson:"external_source"`
	ExternalID       string         `json:"external_id,omitempty" bson:"external_id"`
	ExternalURL      string         `json:"external_url,omitempty" bson:"external_url"`
	FixVersion       string         `json:"fix_version,omitempty" bson:"fix_version"`
	AffectedVersions []string       `json:"affected_versions,omitempty" bson:"affected_versions"`
	ResolvedAt       *time.Time     `json:"resolved_at,omitempty" bson:"resolved_at"`
	ParentIssue      string         `json:"parent_issue,omitempty" bson:"parent_issue"`
	Labels           []string       `json:"labels,omitempty" bson:"labels"`
	Relationships    []Relationship `json:"relationships,omitempty" bson:"relationships"`
	// ExtraFields is retained for imported integration metadata. HTTP write DTOs do
	// not expose it, so clients cannot use it as an unvalidated escape hatch.
	ExtraFields map[string]interface{} `json:"extra_fields,omitempty" bson:"extra_fields"`
	ArchivedAt  *time.Time             `json:"archived_at,omitempty" bson:"archived_at"`
}

// Relationship records a directed connection to another issue in the same space.
// The inverse is intentionally not written automatically: callers can render it
// from either side without risking two records that disagree.
type Relationship struct {
	Issue string           `json:"issue" bson:"issue"`
	Type  RelationshipType `json:"type" bson:"type"`
}

type RelationshipType string

const (
	RelationshipBlocks     RelationshipType = "blocks"
	RelationshipDuplicates RelationshipType = "duplicates"
	RelationshipRelated    RelationshipType = "related"
)

type Activity struct {
	ID        string    `json:"id" bson:"id"`
	Space     string    `json:"space" bson:"space"`
	Issue     string    `json:"issue" bson:"issue"`
	Actor     string    `json:"actor" bson:"actor"`
	Kind      string    `json:"kind" bson:"kind"`
	Field     string    `json:"field,omitempty" bson:"field,omitempty"`
	OldValue  string    `json:"old_value,omitempty" bson:"old_value,omitempty"`
	NewValue  string    `json:"new_value,omitempty" bson:"new_value,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
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
	StatusPlanned    Status = "planned"
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
	i.Title = strings.TrimSpace(i.Title)
	if len(i.Title) == 0 {
		return ErrTitleTooShort
	}
	if len(i.Title) > 100 {
		return ErrTitleTooLong
	}

	if len(i.Description) > 10_000 {
		return ErrDescriptionTooLong
	}
	if !validType(i.Type) {
		return ErrInvalidType
	}
	if !validStatus(i.Status) {
		return ErrInvalidStatus
	}
	if !validPriority(i.Priority) {
		return ErrInvalidPriority
	}
	seenLabels := make(map[string]struct{}, len(i.Labels))
	labels := make([]string, 0, len(i.Labels))
	for _, label := range i.Labels {
		label = strings.TrimSpace(label)
		if label == "" || len(label) > 50 {
			return ErrInvalidLabel
		}
		key := strings.ToLower(label)
		if _, duplicate := seenLabels[key]; !duplicate {
			seenLabels[key] = struct{}{}
			labels = append(labels, label)
		}
	}
	i.Labels = labels
	seenRelationships := make(map[string]struct{}, len(i.Relationships))
	for _, relationship := range i.Relationships {
		if relationship.Issue == "" || relationship.Issue == i.ID || !validRelationshipType(relationship.Type) {
			return ErrInvalidRelationship
		}
		key := string(relationship.Type) + ":" + relationship.Issue
		if _, duplicate := seenRelationships[key]; duplicate {
			return ErrInvalidRelationship
		}
		seenRelationships[key] = struct{}{}
	}

	return nil
}

func validRelationshipType(value RelationshipType) bool {
	return value == RelationshipBlocks || value == RelationshipDuplicates || value == RelationshipRelated
}

func ValidType(value Type) bool         { return validType(value) }
func ValidStatus(value Status) bool     { return validStatus(value) }
func ValidPriority(value Priority) bool { return validPriority(value) }

func validType(value Type) bool {
	switch value {
	case TypeEpic, TypeBug, TypeTask, TypeStory, TypeIdea:
		return true
	default:
		return false
	}
}

func validStatus(value Status) bool {
	switch value {
	case StatusBacklog, StatusPlanned, StatusInProgress, StatusDone, StatusClosed:
		return true
	default:
		return false
	}
}

func validPriority(value Priority) bool {
	switch value {
	case PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical:
		return true
	default:
		return false
	}
}

type Comment struct {
	ID        string    `json:"id" bson:"id"`
	Space     string    `json:"space" bson:"space"`
	Issue     string    `json:"issue" bson:"issue"`
	Author    string    `json:"author" bson:"author"`
	Content   string    `json:"content" bson:"content"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
