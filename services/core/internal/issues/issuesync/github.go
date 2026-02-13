package issuesync

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/fancyinnovations/fancyspaces/core/internal/issues"
	"github.com/fancyinnovations/fancyspaces/core/internal/spaces"
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

	for _, ghIssue := range ghIssues {
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
			ExternalSource:   "github",
			FixVersion:       "",
			AffectedVersions: nil,
			ResolvedAt:       resolvedAt,
			ParentIssue:      "",
			ExtraFields: map[string]interface{}{
				"github_url": ghIssue.GetHTMLURL(),
			},
		}

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
