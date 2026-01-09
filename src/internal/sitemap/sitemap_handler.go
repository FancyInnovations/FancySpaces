package sitemap

import (
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/OliverSchlueter/goutils/problems"
	"github.com/OliverSchlueter/goutils/ratelimit"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/internal/spaces"
)

type Handler struct {
	spaces *spaces.Store
	rl     *ratelimit.Service
}

type Configuration struct {
	Spaces *spaces.Store
}

func NewHandler(cfg Configuration) *Handler {
	rl := ratelimit.NewService(ratelimit.Configuration{
		TokensPerSecond: 2,
		MaxTokens:       10,
	})

	return &Handler{
		spaces: cfg.Spaces,
		rl:     rl,
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/sitemap.xml", h.handle)
}

func (h *Handler) handle(w http.ResponseWriter, r *http.Request) {
	if err := h.rl.CheckRequest(r, "sitemap_xml"); err != nil {
		ratelimit.RateLimitExceededProblem().WriteToHTTP(w)
		return
	}

	allSpaces, err := h.spaces.GetAll()
	if err != nil {
		slog.Error("Failed to retrieve spaces", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	var filteredSpaces []spaces.Space
	for _, s := range allSpaces {
		if s.Status == spaces.StatusApproved || s.Status == spaces.StatusArchived {
			filteredSpaces = append(filteredSpaces, s)
			continue
		}
	}

	var urls []Url
	urls = append(urls,
		Url{
			Loc:        baseURL,
			ChangeFreq: "daily",
			Priority:   "1.0",
		},
		Url{
			Loc:        baseURL + "/explore",
			ChangeFreq: "daily",
			Priority:   "0.9",
		},
		Url{
			Loc:        baseURL + "/explore/minecraft-plugins",
			ChangeFreq: "daily",
			Priority:   "0.9",
		},
		Url{
			Loc:        baseURL + "/explore/hytale-plugins",
			ChangeFreq: "daily",
			Priority:   "0.9",
		},
		Url{
			Loc:        baseURL + "/explore/other-projects",
			ChangeFreq: "daily",
			Priority:   "0.9",
		},
	)

	for _, s := range filteredSpaces {
		urls = append(urls,
			Url{
				Loc:        fmt.Sprintf("%s/spaces/%s", baseURL, s.Slug),
				ChangeFreq: "weekly",
				Priority:   "0.75",
			},
			Url{
				Loc:        fmt.Sprintf("%s/spaces/%s/versions", baseURL, s.Slug),
				ChangeFreq: "daily",
				Priority:   "0.5",
			},
		)
	}

	sitemap := UrlSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Urls:  urls,
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Header().Set("Cache-Control", "public, max-age=10800") // 3 hours

	w.WriteHeader(http.StatusOK)

	enc := xml.NewEncoder(w)
	defer enc.Close()

	enc.Indent("", "  ")

	// Write XML header
	w.Write([]byte(xml.Header))

	if err := enc.Encode(sitemap); err != nil {
		slog.Error("Failed to encode sitemap XML", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}
}
