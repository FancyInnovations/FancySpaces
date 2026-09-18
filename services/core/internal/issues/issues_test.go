package issues_test

import (
	"errors"
	"testing"

	"github.com/fancyinnovations/fancyspaces/core/internal/issues"
	issuesfake "github.com/fancyinnovations/fancyspaces/core/internal/issues/database/fake"
)

func newStore() *issues.Store {
	return issues.New(issues.Configuration{DB: issuesfake.New()})
}

func TestIssueValidationRejectsUnknownEnums(t *testing.T) {
	issue := issues.Issue{
		Title:    "A valid title",
		Type:     issues.Type("unknown"),
		Status:   issues.StatusBacklog,
		Priority: issues.PriorityMedium,
	}
	if !errors.Is(issue.Validate(), issues.ErrInvalidType) {
		t.Fatalf("expected invalid type error, got %v", issue.Validate())
	}
}

func TestStoreCreatesAndArchivesIssue(t *testing.T) {
	store := newStore()
	issue := &issues.Issue{Space: "space-1", Title: "Archive me", Type: issues.TypeTask, Status: issues.StatusBacklog, Priority: issues.PriorityMedium}
	if err := store.CreateIssue(issue); err != nil {
		t.Fatalf("create issue: %v", err)
	}
	if issue.ID == "" || issue.CreatedAt.IsZero() || issue.UpdatedAt.IsZero() {
		t.Fatal("store did not populate issue identity and timestamps")
	}
	if err := store.DeleteIssue(issue.Space, issue.ID); err != nil {
		t.Fatalf("archive issue: %v", err)
	}
	if _, err := store.GetIssue(issue.Space, issue.ID); !errors.Is(err, issues.ErrIssueNotFound) {
		t.Fatalf("expected archived issue to be hidden, got %v", err)
	}
}

func TestStoreCreatesScopedComments(t *testing.T) {
	store := newStore()
	issue := &issues.Issue{Space: "space-1", Title: "Comment me", Type: issues.TypeTask, Status: issues.StatusBacklog, Priority: issues.PriorityMedium}
	if err := store.CreateIssue(issue); err != nil {
		t.Fatalf("create issue: %v", err)
	}
	comment := &issues.Comment{Space: issue.Space, Issue: issue.ID, Author: "user-1", Content: "hello"}
	if err := store.AddComment(comment); err != nil {
		t.Fatalf("add comment: %v", err)
	}
	comments, err := store.GetComments(issue.Space, issue.ID)
	if err != nil || len(comments) != 1 || comments[0].ID == "" {
		t.Fatalf("unexpected comments result: %#v, %v", comments, err)
	}
}

func TestListIssuesFiltersAndPaginates(t *testing.T) {
	store := newStore()
	for _, title := range []string{"Bug in login", "Feature request"} {
		issue := &issues.Issue{Space: "space-1", Title: title, Type: issues.TypeTask, Status: issues.StatusBacklog, Priority: issues.PriorityMedium}
		if err := store.CreateIssue(issue); err != nil {
			t.Fatalf("create issue: %v", err)
		}
	}
	items, total, err := store.ListIssues("space-1", issues.ListOptions{Query: "login", Limit: 1})
	if err != nil || total != 1 || len(items) != 1 {
		t.Fatalf("unexpected list result: items=%d total=%d err=%v", len(items), total, err)
	}
}
