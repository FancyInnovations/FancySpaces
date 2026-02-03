package badges

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/OliverSchlueter/goutils/badgegen"
	"github.com/OliverSchlueter/goutils/problems"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/internal/analytics"
	"github.com/fancyinnovations/fancyspaces/internal/spaces"
	"github.com/fancyinnovations/fancyspaces/internal/versions"
)

type Handler struct {
	spaces    *spaces.Store
	versions  *versions.Store
	analytics *analytics.Store
}

type Configuration struct {
	Spaces    *spaces.Store
	Versions  *versions.Store
	Analytics *analytics.Store
}

func NewHandler(config Configuration) *Handler {
	return &Handler{
		spaces:    config.Spaces,
		versions:  config.Versions,
		analytics: config.Analytics,
	}
}

func (h *Handler) Register(prefix string, mux *http.ServeMux) {
	mux.HandleFunc(prefix+"/badges/downloads", h.handleSpaceDownloadsBadge)
	mux.HandleFunc(prefix+"/badges/latest-version", h.handleSpaceLatestVersionBadge)
}

func (h *Handler) handleSpaceDownloadsBadge(w http.ResponseWriter, r *http.Request) {
	sid := r.URL.Query().Get("space_id")
	if sid == "" {
		problems.ValidationError("space_id", "Space ID is required").WriteToHTTP(w)
		return
	}

	s, err := h.spaces.Get(sid)
	if err != nil {
		if errors.Is(err, spaces.ErrSpaceNotFound) {
			problems.NotFound("Space", sid).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get space by id", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	if !s.ReleaseSettings.Enabled {
		spaces.ProblemFeatureNotEnabled("releases").WriteToHTTP(w)
		return
	}

	downloads, err := h.analytics.GetDownloadCountForSpace(r.Context(), s.ID)
	if err != nil {
		slog.Error("Failed to get download count for space", sloki.WrapError(err), slog.String("space_id", s.ID))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	badge := badgegen.Generate("Downloads", fmt.Sprintf("%d", downloads), "#541ea6")

	w.Header().Set("Content-Type", "image/svg+xml;charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=300") // 5 minutes
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(badge))
}

func (h *Handler) handleSpaceLatestVersionBadge(w http.ResponseWriter, r *http.Request) {
	sid := r.URL.Query().Get("space_id")
	if sid == "" {
		problems.ValidationError("space_id", "Space ID is required").WriteToHTTP(w)
		return
	}

	s, err := h.spaces.Get(sid)
	if err != nil {
		if errors.Is(err, spaces.ErrSpaceNotFound) {
			problems.NotFound("Space", sid).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get space by id", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	if !s.ReleaseSettings.Enabled {
		spaces.ProblemFeatureNotEnabled("releases").WriteToHTTP(w)
		return
	}

	latestVer, err := h.versions.GetLatest(r.Context(), s.ID, "", "")
	if err != nil {
		if errors.Is(err, versions.ErrVersionNotFound) {
			latestVer = &versions.Version{
				Name: "N/A",
			}
		} else {
			slog.Error("Failed to get latest version for space", sloki.WrapError(err), slog.String("space_id", s.ID))
			problems.InternalServerError("").WriteToHTTP(w)
			return
		}
	}

	badge := badgegen.Generate("Latest version", latestVer.Name, "green")

	w.Header().Set("Content-Type", "image/svg+xml;charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=300") // 5 minutes
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(badge))
}
