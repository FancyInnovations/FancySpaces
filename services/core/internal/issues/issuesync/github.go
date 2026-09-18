package issuesync

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/fancyinnovations/fancyspaces/core/internal/issues"
	"github.com/fancyinnovations/fancyspaces/integrations/spaces-go-sdk/spaces"
	"github.com/google/go-github/v82/github"
)

func (s *Service) syncGitHub(space spaces.Space) error {
	startTime := time.Now()

	listOpts := &github.IssueListByRepoOptions{
		State: "all",
	}
	if space.IssueSettings.GitHubSyncLabel != "" {
		listOpts.Labels = []string{space.IssueSettings.GitHubSyncLabel}
	}

	ghIssues, _, err := s.ghc.Issues.ListByRepo(context.Background(), space.IssueSettings.GitHubSyncOwner, space.IssueSettings.GitHubSyncRepo, listOpts)
	if err != nil {
		return fmt.Errorf("failed to fetch GitHub issues: %w", err)
	}

	fsIssues, err := s.issuesStore.GetIssues(space.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch FancySpaces issues: %w", err)
	}
	fsIssuesMap := issuesToMap(fsIssues)

	createdIssues := 0
	updatedIssues := 0
	seenIssueIDs := make(map[string]struct{}, len(ghIssues))

	for _, ghIssue := range ghIssues {
		if ghIssue.IsPullRequest() {
			continue // skip pull requests
		}

		var status issues.Status
		if ghIssue.GetState() == "open" {
			status = issues.StatusBacklog
		} else {
			status = issues.StatusClosed
		}

		var resolvedAt *time.Time
		if ghIssue.ClosedAt != nil {
			t := ghIssue.GetClosedAt().Time
			resolvedAt = &t
		}

		fsIssue := issues.Issue{
			ID:               "gh-" + strconv.Itoa(ghIssue.GetNumber()),
			Space:            space.ID,
			Title:            ghIssue.GetTitle(),
			Description:      ghIssue.GetBody(),
			Type:             issues.TypeBug,
			Status:           status,
			Priority:         issues.PriorityMedium,
			Assignee:         "",
			Reporter:         ghIssue.GetUser().GetLogin(),
			CreatedAt:        ghIssue.GetCreatedAt().Time,
			UpdatedAt:        ghIssue.GetUpdatedAt().Time,
			ExternalSource:   issues.ExternalSourceGitHub,
			ExternalID:       strconv.Itoa(ghIssue.GetNumber()),
			ExternalURL:      ghIssue.GetHTMLURL(),
			FixVersion:       "",
			AffectedVersions: nil,
			ResolvedAt:       resolvedAt,
			ParentIssue:      "",
			ExtraFields: map[string]interface{}{
				"github_url": ghIssue.GetHTMLURL(),
			},
		}
		seenIssueIDs[fsIssue.ID] = struct{}{}

		existingIssue, exists := fsIssuesMap[fsIssue.ID]
		if !exists {
			slog.Debug(
				"Creating FancySpaces issue from GitHub issue",
				slog.String("space_id", space.ID),
				slog.Int("github_issue_number", ghIssue.GetNumber()),
				slog.String("fancyspaces_issue_id", fsIssue.ID),
			)

			if err := s.issuesStore.ForceCreateIssue(&fsIssue); err != nil {
				return fmt.Errorf("failed to create FancySpaces issue for GitHub issue #%d: %w", ghIssue.GetNumber(), err)
			}

			createdIssues++
		} else {
			// GitHub owns the title, body and open/closed state. Preserve FancySpaces
			// planning metadata so a sync cannot erase local triage decisions.
			if existingIssue.ArchivedAt != nil {
				continue
			}
			fsIssue.CreatedAt = existingIssue.CreatedAt
			fsIssue.Type = existingIssue.Type
			fsIssue.Priority = existingIssue.Priority
			fsIssue.Assignee = existingIssue.Assignee
			fsIssue.FixVersion = existingIssue.FixVersion
			fsIssue.AffectedVersions = existingIssue.AffectedVersions
			fsIssue.ParentIssue = existingIssue.ParentIssue
			fsIssue.Labels = existingIssue.Labels
			fsIssue.Relationships = existingIssue.Relationships
			fsIssue.ExtraFields = existingIssue.ExtraFields
			if hasIssueChange(&existingIssue, &fsIssue) {
				slog.Debug(
					"Updating FancySpaces issue from GitHub issue",
					slog.String("space_id", space.ID),
					slog.Int("github_issue_number", ghIssue.GetNumber()),
					slog.String("fancyspaces_issue_id", fsIssue.ID),
				)

				if err := s.issuesStore.ForceUpdateIssue(&fsIssue); err != nil {
					return fmt.Errorf("failed to update FancySpaces issue for GitHub issue #%d: %w", ghIssue.GetNumber(), err)
				}
				updatedIssues++
			}
		}

	}

	// When no label filter is active, an external issue missing from GitHub's
	// complete listing was deleted. Archive the local projection instead of
	// leaving a permanently stale issue visible.
	if space.IssueSettings.GitHubSyncLabel == "" {
		for _, existing := range fsIssues {
			if existing.ExternalSource != issues.ExternalSourceGitHub || existing.ArchivedAt != nil {
				continue
			}
			if _, seen := seenIssueIDs[existing.ID]; seen {
				continue
			}
			now := time.Now()
			existing.ArchivedAt = &now
			existing.UpdatedAt = now
			if err := s.issuesStore.ForceUpdateIssue(&existing); err != nil {
				return fmt.Errorf("failed to archive deleted GitHub issue %s: %w", existing.ID, err)
			}
			updatedIssues++
		}
	}

	timeElapsed := time.Since(startTime)
	slog.Info("GitHub issue sync completed",
		slog.String("space_id", space.ID),
		slog.Int("github_issues_processed", len(ghIssues)),
		slog.Int("fancyspaces_issues_created", createdIssues),
		slog.Int("fancyspaces_issues_updated", updatedIssues),
		slog.Duration("time_elapsed", timeElapsed),
	)

	return nil
}

func hasIssueChange(one *issues.Issue, two *issues.Issue) bool {
	if one.Title != two.Title {
		return true
	}
	if one.Description != two.Description {
		return true
	}
	if one.Status != two.Status {
		return true
	}
	if (one.ResolvedAt == nil) != (two.ResolvedAt == nil) {
		return true
	}
	if one.ResolvedAt != nil && two.ResolvedAt != nil && !one.ResolvedAt.Equal(*two.ResolvedAt) {
		return true
	}
	if one.ExternalSource != two.ExternalSource {
		return true
	}
	return false
}

func issuesToMap(list []issues.Issue) map[string]issues.Issue {
	m := make(map[string]issues.Issue)
	for _, issue := range list {
		m[issue.ID] = issue
	}
	return m
}
