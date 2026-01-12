package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/OliverSchlueter/goutils/idgen"
	"github.com/OliverSchlueter/goutils/problems"
	"github.com/OliverSchlueter/goutils/ratelimit"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/internal/analytics"
	"github.com/fancyinnovations/fancyspaces/internal/auth"
	"github.com/fancyinnovations/fancyspaces/internal/spaces"
	"github.com/fancyinnovations/fancyspaces/internal/versions"
)

type Handler struct {
	store             *versions.Store
	spaces            *spaces.Store
	analytics         *analytics.Store
	userFromCtx       func(ctx context.Context) *auth.User
	downloadRatelimit *ratelimit.Service
}

type Configuration struct {
	Store       *versions.Store
	Spaces      *spaces.Store
	Analytics   *analytics.Store
	UserFromCtx func(ctx context.Context) *auth.User
}

func New(cfg Configuration) *Handler {
	downloadRatelimit := ratelimit.NewService(ratelimit.Configuration{
		TokensPerSecond: 1,
		MaxTokens:       5,
	})

	return &Handler{
		store:             cfg.Store,
		spaces:            cfg.Spaces,
		analytics:         cfg.Analytics,
		userFromCtx:       cfg.UserFromCtx,
		downloadRatelimit: downloadRatelimit,
	}
}

func (h *Handler) Register(prefix string, mux *http.ServeMux) {
	mux.HandleFunc(prefix+"/spaces/{space_id}/versions", h.handleVersions)
	mux.HandleFunc(prefix+"/spaces/{space_id}/versions/{version_id}", h.handleVersion)
	mux.HandleFunc(prefix+"/spaces/{space_id}/versions/{version_id}/downloads", h.handleVersionDownloads)
	mux.HandleFunc(prefix+"/spaces/{space_id}/versions/{version_id}/files/{file_name}", h.handleVersionFile)
}

func (h *Handler) handleVersions(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("space_id")
	if sid == "" {
		problems.ValidationError("space_id", "Space ID is required").WriteToHTTP(w)
		return
	}

	space, err := h.spaces.Get(sid)
	if err != nil {
		if errors.Is(err, spaces.ErrSpaceNotFound) {
			problems.NotFound("Space", sid).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get space by id", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}
	if space.Status != spaces.StatusApproved && space.Status != spaces.StatusArchived {
		u := h.userFromCtx(r.Context())
		if u == nil || !u.Verified || !u.IsActive || !space.IsMember(u) {
			problems.NotFound("Space", space.ID).WriteToHTTP(w)
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetVersions(w, r, space.ID)
	case http.MethodPost:
		h.handleCreateVersion(w, r, space.ID)
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodPost}).WriteToHTTP(w)
	}
}

func (h *Handler) handleVersion(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("space_id")
	if sid == "" {
		problems.ValidationError("space_id", "Space ID is required").WriteToHTTP(w)
		return
	}

	vid := r.PathValue("version_id")
	if vid == "" {
		problems.ValidationError("version_id", "Version ID or name is required").WriteToHTTP(w)
		return
	}

	space, err := h.spaces.Get(sid)
	if err != nil {
		if errors.Is(err, spaces.ErrSpaceNotFound) {
			problems.NotFound("Space", sid).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get space by id", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}
	if space.Status != spaces.StatusApproved && space.Status != spaces.StatusArchived {
		u := h.userFromCtx(r.Context())
		if u == nil || !u.Verified || !u.IsActive || !space.IsMember(u) {
			problems.NotFound("Space", space.ID).WriteToHTTP(w)
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetVersion(w, r, space.ID, vid)
	case http.MethodDelete:
		h.handleDeleteVersion(w, r, space.ID, vid)
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodDelete}).WriteToHTTP(w)
	}
}

// no auth required
func (h *Handler) handleGetVersions(w http.ResponseWriter, r *http.Request, spaceID string) {
	all, err := h.store.GetAll(r.Context(), spaceID)
	if err != nil {
		slog.Error("Failed to get versions", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60") // 1 minute
	json.NewEncoder(w).Encode(all)
}

// no auth required
func (h *Handler) handleGetVersion(w http.ResponseWriter, r *http.Request, spaceID, versionID string) {
	if versionID == "latest" {
		channel := r.URL.Query().Get("channel")
		platform := r.URL.Query().Get("platform")

		ver, err := h.store.GetLatest(r.Context(), spaceID, channel, platform)
		if err != nil {
			if errors.Is(err, versions.ErrVersionNotFound) {
				problems.NotFound("Version", "latest").WriteToHTTP(w)
				return
			}

			slog.Error("Failed to get latest version", sloki.WrapError(err))
			problems.InternalServerError("").WriteToHTTP(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=60") // 1 minute
		json.NewEncoder(w).Encode(ver)
		return
	}

	ver, err := h.store.Get(r.Context(), spaceID, versionID)
	if err != nil {
		slog.Error("Failed to get version", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600") // 1 hour
	json.NewEncoder(w).Encode(ver)
}

func (h *Handler) handleCreateVersion(w http.ResponseWriter, r *http.Request, spaceID string) {
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}

	space, err := h.spaces.Get(spaceID)
	if err != nil {
		if errors.Is(err, spaces.ErrSpaceNotFound) {
			problems.NotFound("Space", spaceID).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get space by id", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	if !space.HasWriteAccess(u) {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	var req CreateVersionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode create version request", sloki.WrapError(err))
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}

	ver := versions.Version{
		SpaceID:                   spaceID,
		ID:                        idgen.GenerateID(8),
		Name:                      req.Name,
		Platform:                  req.Platform,
		Channel:                   req.Channel,
		PublishedAt:               time.Now(),
		Changelog:                 req.Changelog,
		SupportedPlatformVersions: req.SupportedPlatformVersions,
		Files:                     []versions.VersionFile{},
	}

	if err := h.store.Create(r.Context(), &ver); err != nil {
		slog.Error("Failed to create version", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

func (h *Handler) handleDeleteVersion(w http.ResponseWriter, r *http.Request, spaceID, versionID string) {
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}

	space, err := h.spaces.Get(spaceID)
	if err != nil {
		if errors.Is(err, spaces.ErrSpaceNotFound) {
			problems.NotFound("Space", spaceID).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get space by id", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	if !space.HasWriteAccess(u) {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	if err := h.store.Delete(r.Context(), spaceID, versionID); err != nil {
		slog.Error("Failed to delete version", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// no auth required
func (h *Handler) handleVersionDownloads(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet}).WriteToHTTP(w)
		return
	}

	sid := r.PathValue("space_id")
	if sid == "" {
		problems.ValidationError("space_id", "Space ID is required").WriteToHTTP(w)
		return
	}
	space, err := h.spaces.Get(sid)
	if err != nil {
		if errors.Is(err, spaces.ErrSpaceNotFound) {
			problems.NotFound("Space", sid).WriteToHTTP(w)
			return
		}
		slog.Error("Failed to get space by id", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	vid := r.PathValue("version_id")
	if vid == "" {
		problems.ValidationError("version_id", "Version ID is required").WriteToHTTP(w)
		return
	}
	ver, err := h.store.Get(r.Context(), space.ID, vid)
	if err != nil {
		if errors.Is(err, versions.ErrVersionNotFound) {
			problems.NotFound("Version", vid).WriteToHTTP(w)
			return
		}
		slog.Error("Failed to get version by id", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	count, err := h.analytics.GetDownloadCountForVersion(r.Context(), space.ID, ver.ID)
	if err != nil {
		slog.Error("Failed to get download count for version", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	resp := VersionDownloadsResp{
		Downloads: count,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60") // 1 minute
	json.NewEncoder(w).Encode(resp)
}
