package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/OliverSchlueter/goutils/problems"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/core/internal/analytics"
	"github.com/fancyinnovations/fancyspaces/core/internal/spaces"
	"github.com/fancyinnovations/fancyspaces/integrations/idp-go-sdk/idp"
)

type Handler struct {
	store       *spaces.Store
	analytics   *analytics.Store
	userFromCtx func(ctx context.Context) *idp.User
}

type Configuration struct {
	Store       *spaces.Store
	Analytics   *analytics.Store
	UserFromCtx func(ctx context.Context) *idp.User
}

func New(cfg Configuration) *Handler {
	return &Handler{
		store:       cfg.Store,
		analytics:   cfg.Analytics,
		userFromCtx: cfg.UserFromCtx,
	}
}

func (h *Handler) Register(prefix string, mux *http.ServeMux) {
	mux.HandleFunc(prefix+"/spaces", h.handleSpaces)
	mux.HandleFunc(prefix+"/spaces/{space_id}", h.handleSpace)
	mux.HandleFunc(prefix+"/spaces/{space_id}/status", h.handleChangeStatus)
	mux.HandleFunc(prefix+"/spaces/{space_id}/downloads", h.handleDownloads)

	mux.HandleFunc(prefix+"/spaces/{space_id}/members", h.handleMembers)
	mux.HandleFunc(prefix+"/spaces/{space_id}/members/{user_id}", h.handleMember)
}

func (h *Handler) handleSpaces(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetSpaces(w, r)
	case http.MethodPost:
		h.handleCreateSpace(w, r)
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodPost}).WriteToHTTP(w)
	}
}

func (h *Handler) handleSpace(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("space_id")
	if sid == "" {
		problems.ValidationError("space_id", "Space ID is required").WriteToHTTP(w)
		return
	}

	s, err := h.store.Get(sid)
	if err != nil {
		if errors.Is(err, spaces.ErrSpaceNotFound) {
			problems.NotFound("Space", sid).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get space by id", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetSpace(w, r, s)
	case http.MethodPut:
		h.handleUpdateSpace(w, r, s)
	case http.MethodDelete:
		h.handleDeleteSpace(w, r, s)
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodPut, http.MethodDelete}).WriteToHTTP(w)
	}
}

// no auth required
func (h *Handler) handleGetSpaces(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Accept") != "application/json" {
		problems.WrongAcceptType("application/json", r.Header.Get("Accept")).WriteToHTTP(w)
		return
	}

	all, err := h.store.GetAll()
	if err != nil {
		slog.Error("Failed to get spaces", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	u := h.userFromCtx(r.Context())

	var res []spaces.Space
	for _, s := range all {
		if s.Status == spaces.StatusApproved || s.Status == spaces.StatusArchived {
			res = append(res, s)
			continue
		}

		if u == nil || !u.Verified || !u.IsActive {
			continue
		}
		if s.IsMember(u) {
			res = append(res, s)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600") // 1 hour
	json.NewEncoder(w).Encode(res)
}

// no auth required
func (h *Handler) handleGetSpace(w http.ResponseWriter, r *http.Request, s *spaces.Space) {
	if r.Header.Get("Accept") != "application/json" {
		problems.WrongAcceptType("application/json", r.Header.Get("Accept")).WriteToHTTP(w)
		return
	}

	s, err := h.store.Get(s.ID)
	if err != nil {
		if errors.Is(err, spaces.ErrSpaceNotFound) {
			problems.NotFound("Space", s.ID).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get space by id", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	if s.Status != spaces.StatusApproved && s.Status != spaces.StatusArchived {
		u := h.userFromCtx(r.Context())
		if u == nil || !u.Verified || !u.IsActive || !s.IsMember(u) {
			problems.NotFound("Space", s.ID).WriteToHTTP(w)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	json.NewEncoder(w).Encode(s)
}

func (h *Handler) handleCreateSpace(w http.ResponseWriter, r *http.Request) {
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}

	if ct := r.Header.Get("Content-Type"); ct != "application/json" {
		problems.WrongContentType("application/json", ct).WriteToHTTP(w)
		return
	}

	var req spaces.CreateOrUpdateSpaceReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}

	s, err := h.store.Create(u, &req)
	if err != nil {
		if errors.Is(err, spaces.ErrSpaceAlreadyExists) {
			problems.AlreadyExists("Space", req.Slug).WriteToHTTP(w)
			return
		}
		if errors.Is(err, spaces.ErrUserNotActive) || errors.Is(err, spaces.ErrUserNotVerified) {
			problems.Unauthorized().WriteToHTTP(w)
			return
		}
		if errors.Is(err, spaces.ErrSlugTooLong) || errors.Is(err, spaces.ErrSlugTooShort) ||
			errors.Is(err, spaces.ErrTitleTooLong) || errors.Is(err, spaces.ErrTitleTooShort) ||
			errors.Is(err, spaces.ErrDescriptionTooLong) {
			problems.ValidationError("body", err.Error()).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to create space", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) handleUpdateSpace(w http.ResponseWriter, r *http.Request, s *spaces.Space) {
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}

	if !s.HasFullAccess(u) {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	if ct := r.Header.Get("Content-Type"); ct != "application/json" {
		problems.WrongContentType("application/json", ct).WriteToHTTP(w)
		return
	}

	var req spaces.CreateOrUpdateSpaceReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}

	if err := h.store.Update(s.ID, &req); err != nil {
		if errors.Is(err, spaces.ErrSpaceNotFound) {
			problems.NotFound("Space", s.ID).WriteToHTTP(w)
			return
		}
		if errors.Is(err, spaces.ErrSpaceAlreadyExists) {
			problems.AlreadyExists("Space", req.Slug).WriteToHTTP(w)
			return
		}
		if errors.Is(err, spaces.ErrUserNotActive) || errors.Is(err, spaces.ErrUserNotVerified) {
			problems.Unauthorized().WriteToHTTP(w)
			return
		}
		if errors.Is(err, spaces.ErrSlugTooLong) || errors.Is(err, spaces.ErrSlugTooShort) ||
			errors.Is(err, spaces.ErrTitleTooLong) || errors.Is(err, spaces.ErrTitleTooShort) ||
			errors.Is(err, spaces.ErrDescriptionTooLong) {
			problems.ValidationError("body", err.Error()).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to update space", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleDeleteSpace(w http.ResponseWriter, r *http.Request, s *spaces.Space) {
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}

	if !s.IsOwner(u) {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	if err := h.store.Delete(s.ID); err != nil {
		slog.Error("Failed to delete space", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleChangeStatus(w http.ResponseWriter, r *http.Request) {
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}

	sid := r.PathValue("space_id")
	if sid == "" {
		problems.ValidationError("space_id", "Space ID is required").WriteToHTTP(w)
		return
	}

	s, err := h.store.Get(sid)
	if err != nil {
		if errors.Is(err, spaces.ErrSpaceNotFound) {
			problems.NotFound("Space", sid).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get space by id", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	if r.Method != http.MethodPut {
		problems.MethodNotAllowed(r.Method, []string{http.MethodPut}).WriteToHTTP(w)
		return
	}

	var req ChangeStatusReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}

	if req.To == spaces.StatusApproved || req.To == spaces.StatusRejected || req.To == spaces.StatusBanned {
		if !u.IsAdmin() {
			problems.Forbidden().WriteToHTTP(w)
			return
		}
	} else if req.To == spaces.StatusPrivate || req.To == spaces.StatusArchived {
		if !s.HasFullAccess(u) {
			problems.Forbidden().WriteToHTTP(w)
			return
		}

		if req.To == spaces.StatusArchived && s.Status != spaces.StatusApproved {
			problems.ValidationError("to", "Space must be approved before it can be archived").WriteToHTTP(w)
			return
		}
		if req.To == spaces.StatusPrivate && s.Status != spaces.StatusApproved {
			problems.ValidationError("to", "Space must be approved before it can be made private").WriteToHTTP(w)
			return
		}
	} else {
		problems.ValidationError("to", "Invalid status").WriteToHTTP(w)
		return
	}

	if err := h.store.ChangeStatus(s, req.To); err != nil {
		if errors.Is(err, spaces.ErrSpaceNotFound) {
			problems.NotFound("Space", s.ID).WriteToHTTP(w)
			return
		}
		slog.Error("Failed to change space status", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// no auth required
func (h *Handler) handleDownloads(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet}).WriteToHTTP(w)
		return
	}

	sid := r.PathValue("space_id")
	if sid == "" {
		problems.ValidationError("space_id", "Space ID is required").WriteToHTTP(w)
		return
	}

	s, err := h.store.Get(sid)
	if err != nil {
		if errors.Is(err, spaces.ErrSpaceNotFound) {
			problems.NotFound("Space", sid).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get space by id", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	count, err := h.analytics.GetDownloadCountForSpace(r.Context(), s.ID)
	if err != nil {
		slog.Error("Failed to get download count for space", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	versions, err := h.analytics.GetDownloadCountForVersions(r.Context(), s.ID)
	if err != nil {
		slog.Error("Failed to get download count for space versions", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	resp := SpaceDownloadsResp{
		Downloads: count,
		Versions:  versions,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60") // 1 minute
	json.NewEncoder(w).Encode(resp)
}
