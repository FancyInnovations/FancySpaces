package issuesync

import (
	"log/slog"
	"time"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/core/internal/issues"
	"github.com/fancyinnovations/fancyspaces/core/internal/spaces"
	"github.com/gofri/go-github-pagination/githubpagination"
	"github.com/google/go-github/v82/github"
)

type Service struct {
	spacesStore *spaces.Store
	issuesStore *issues.Store

	ghc *github.Client
}

type Configuration struct {
	SpacesStore  *spaces.Store
	IssuesStore  *issues.Store
	GitHubClient *github.Client
}

func NewService(cfg *Configuration) *Service {
	if cfg.GitHubClient == nil {
		paginator := githubpagination.NewClient(nil,
			githubpagination.WithPerPage(100),
		)
		cfg.GitHubClient = github.NewClient(paginator)
	}

	return &Service{
		spacesStore: cfg.SpacesStore,
		issuesStore: cfg.IssuesStore,
		ghc:         cfg.GitHubClient,
	}
}

func (s *Service) StartScheduler() {
	ticker := time.NewTicker(4 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				s.SyncIssuesForAllSpaces()
			}
		}
	}()
}

func (s *Service) SyncIssuesForAllSpaces() {
	allSpaces, err := s.spacesStore.GetAll()
	if err != nil {
		slog.Error("failed to fetch all spaces", sloki.WrapError(err))
		return
	}

	slog.Info("Starting issue sync for all spaces")
	for _, space := range allSpaces {
		if err := s.SyncIssuesForSpace(&space); err != nil {
			slog.Error("failed to sync issues for space", slog.String("space_id", space.ID), sloki.WrapError(err))
		}
	}
}

func (s *Service) SyncIssuesForSpace(space *spaces.Space) error {
	if !space.IssueSettings.Enabled {
		return nil
	}

	if space.IssueSettings.GitHubSync {
		if err := s.syncGitHub(*space); err != nil {
			return err
		}
	}

	return nil
}
