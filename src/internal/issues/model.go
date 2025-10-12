package issues

type Issue struct {
	ID             string   `json:"id"`
	Space          string   `json:"space"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Type           Type     `json:"type"`
	Status         Status   `json:"status"`
	Priority       Priority `json:"priority"`
	Assignee       string   `json:"assignee"`
	Reporter       string   `json:"reporter"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
	ExternalSource string   `json:"external_source,omitempty"`
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
	ExternalSourceGitHub  ExternalSource = "github"
	ExternalSourceDiscord ExternalSource = "discord"
)

type Comment struct {
	ID        string `json:"id"`
	Issue     string `json:"issue"`
	Author    string `json:"author"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
