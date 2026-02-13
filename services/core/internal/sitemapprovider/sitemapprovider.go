package sitemapprovider

import (
	"fmt"
	"log/slog"

	"github.com/OliverSchlueter/goutils/sitemapgen"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/core/internal/spaces"
)

const baseURL = "https://fancyspaces.net"

type Service struct {
	spaces *spaces.Store
}

type Configuration struct {
	Spaces *spaces.Store
}

func NewService(cfg *Configuration) *Service {
	return &Service{
		spaces: cfg.Spaces,
	}
}

func (s *Service) GenerateUrls() []sitemapgen.Url {
	var urls []sitemapgen.Url
	urls = append(urls,
		sitemapgen.Url{
			Loc:        baseURL,
			ChangeFreq: "daily",
			Priority:   "1.0",
		},
		sitemapgen.Url{
			Loc:        baseURL + "/explore",
			ChangeFreq: "daily",
			Priority:   "0.9",
		},
		sitemapgen.Url{
			Loc:        baseURL + "/explore/minecraft-plugins",
			ChangeFreq: "daily",
			Priority:   "0.9",
		},
		sitemapgen.Url{
			Loc:        baseURL + "/explore/hytale-plugins",
			ChangeFreq: "daily",
			Priority:   "0.9",
		},
		sitemapgen.Url{
			Loc:        baseURL + "/explore/other-projects",
			ChangeFreq: "daily",
			Priority:   "0.9",
		},
		sitemapgen.Url{
			Loc:        baseURL + "/explore/by-other-creators",
			ChangeFreq: "daily",
			Priority:   "0.9",
		},
	)

	allSpaces, err := s.spaces.GetAll()
	if err != nil {
		slog.Error("Failed to retrieve spaces", sloki.WrapError(err))
		return urls
	}

	var filteredSpaces []spaces.Space
	for _, space := range allSpaces {
		if space.Status == spaces.StatusApproved || space.Status == spaces.StatusArchived {
			filteredSpaces = append(filteredSpaces, space)
			continue
		}
	}

	for _, space := range filteredSpaces {
		urls = append(urls,
			sitemapgen.Url{
				Loc:        fmt.Sprintf("%s/spaces/%s", baseURL, space.Slug),
				ChangeFreq: "weekly",
				Priority:   "0.75",
			},
		)

		if space.ReleaseSettings.Enabled {
			urls = append(urls,
				sitemapgen.Url{
					Loc:        fmt.Sprintf("%s/spaces/%s/versions", baseURL, space.Slug),
					ChangeFreq: "daily",
					Priority:   "0.5",
				},
			)
		}

		if space.IssueSettings.Enabled {
			urls = append(urls,
				sitemapgen.Url{
					Loc:        fmt.Sprintf("%s/spaces/%s/issues", baseURL, space.Slug),
					ChangeFreq: "daily",
					Priority:   "0.5",
				},
			)
		}
	}

	return urls
}
