package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/OliverSchlueter/goutils/problems"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/src/internal/auth"
	"github.com/fancyinnovations/fancyspaces/src/internal/spaces"
	"github.com/fancyinnovations/fancyspaces/src/internal/versions"
)

type Handler struct {
	store       *versions.Store
	spaces      *spaces.Store
	userFromCtx func(ctx context.Context) *auth.User
}

type Configuration struct {
	Store       *versions.Store
	Spaces      *spaces.Store
	UserFromCtx func(ctx context.Context) *auth.User
}

func New(cfg Configuration) *Handler {
	return &Handler{
		store:       cfg.Store,
		spaces:      cfg.Spaces,
		userFromCtx: cfg.UserFromCtx,
	}
}

func (h *Handler) Register(prefix string, mux *http.ServeMux) {
	mux.HandleFunc(prefix+"/spaces/{space_id}/versions", h.handleVersions)
	mux.HandleFunc(prefix+"/spaces/{space_id}/versions/{version_id}", h.handleVersion)
	mux.HandleFunc(prefix+"/spaces/{space_id}/versions/{version_id}/files/{file_name}", h.handleVersionFile)
}

func (h *Handler) handleVersions(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("space_id")
	if sid == "" {
		problems.ValidationError("space_id", "Space ID is required").WriteToHTTP(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetVersions(w, r, sid)
	case http.MethodPost:
		h.handleCreateVersion(w, r, sid)
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

	switch r.Method {
	case http.MethodGet:
		h.handleGetVersion(w, r, sid, vid)
	case http.MethodDelete:
		h.handleDeleteVersion(w, r, sid, vid)
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodDelete}).WriteToHTTP(w)
	}
}

func (h *Handler) handleVersionFile(w http.ResponseWriter, r *http.Request) {
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

	fileName := r.PathValue("file_name")
	if fileName == "" {
		problems.ValidationError("file_name", "File name is required").WriteToHTTP(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleDownloadVersionFile(w, r, sid, vid, fileName)
	case http.MethodPost:
		h.handleUploadVersionFile(w, r, sid, vid, fileName)
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodPost}).WriteToHTTP(w)
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
	w.Header().Set("Cache-Control", "public, 3600")
	json.NewEncoder(w).Encode(all)
}

// no auth required
func (h *Handler) handleGetVersion(w http.ResponseWriter, r *http.Request, spaceID, versionID string) {
	if versionID == "latest" {
		channel := r.URL.Query().Get("channel")
		if channel == "" {
			channel = "release"
		}

		ver, err := h.store.GetLatest(r.Context(), spaceID, channel)
		if err != nil {
			slog.Error("Failed to get latest version", sloki.WrapError(err))
			problems.InternalServerError("").WriteToHTTP(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, 60") // 1 minute
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
	w.Header().Set("Cache-Control", "public, 3600") // 1 hour
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
		ID:                        req.Name,
		Name:                      req.Name,
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

func (h *Handler) handleUploadVersionFile(w http.ResponseWriter, r *http.Request, spaceID, versionID, fileName string) {
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

	ver, err := h.store.Get(r.Context(), spaceID, versionID)
	if err != nil {
		slog.Error("Failed to get version", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	data := make([]byte, r.ContentLength)
	_, err = r.Body.Read(data)
	if err != nil {
		slog.Error("Failed to read file data", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	if err := h.store.UploadVersionFile(r.Context(), ver, fileName, data); err != nil {
		slog.Error("Failed to upload version file", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// no auth required
func (h *Handler) handleDownloadVersionFile(w http.ResponseWriter, r *http.Request, spaceID, versionID, fileName string) {
	if versionID == "latest" {
		channel := r.URL.Query().Get("channel")
		if channel == "" {
			channel = "release"
		}

		ver, err := h.store.GetLatest(r.Context(), spaceID, channel)
		if err != nil {
			slog.Error("Failed to get latest version", sloki.WrapError(err))
			problems.InternalServerError("").WriteToHTTP(w)
			return
		}
		versionID = ver.ID
	}

	data, err := h.store.DownloadVersionFile(r.Context(), r, spaceID, versionID, fileName)
	if err != nil {
		slog.Error("Failed to download version file", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	w.Header().Set("Cache-Control", "public, 86400") // 24h
	w.Write(data)
}
