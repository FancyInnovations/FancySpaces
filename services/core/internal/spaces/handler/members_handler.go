package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/OliverSchlueter/goutils/problems"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/integrations/spaces-go-sdk/spaces"
)

func (h *Handler) handleMembers(w http.ResponseWriter, r *http.Request) {
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
		h.handleGetMembers(w, r, s)
	case http.MethodPost:
		h.handleAddMember(w, r, s)
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodPost}).WriteToHTTP(w)
	}
}

func (h *Handler) handleMember(w http.ResponseWriter, r *http.Request) {
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

	uid := r.PathValue("user_id")
	if uid == "" {
		problems.ValidationError("user_id", "User ID is required").WriteToHTTP(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetMember(w, r, s, uid)
	case http.MethodPut:
		h.handleUpdateMember(w, r, s, uid)
	case http.MethodDelete:
		h.handleDeleteMember(w, r, s, uid)
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodPut, http.MethodDelete}).WriteToHTTP(w)
	}
}

func (h *Handler) handleGetMembers(w http.ResponseWriter, r *http.Request, s *spaces.Space) {
	if r.Header.Get("Accept") != "application/json" {
		problems.WrongAcceptType("application/json", r.Header.Get("Accept")).WriteToHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600") // 1 hour
	json.NewEncoder(w).Encode(s.Members)
}

func (h *Handler) handleGetMember(w http.ResponseWriter, r *http.Request, s *spaces.Space, uid string) {
	if r.Header.Get("Accept") != "application/json" {
		problems.WrongAcceptType("application/json", r.Header.Get("Accept")).WriteToHTTP(w)
		return
	}

	for _, m := range s.Members {
		if m.UserID == uid {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Cache-Control", "public, max-age=3600") // 1 hour
			json.NewEncoder(w).Encode(m)
			return
		}
	}

	problems.NotFound("Member", uid).WriteToHTTP(w)
}

func (h *Handler) handleAddMember(w http.ResponseWriter, r *http.Request, s *spaces.Space) {
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}

	if !s.IsOwner(u) {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	if ct := r.Header.Get("Content-Type"); ct != "application/json" {
		problems.WrongContentType("application/json", ct).WriteToHTTP(w)
		return
	}

	var req spaces.Member
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}

	if req.UserID == "" {
		problems.ValidationError("user_id", "User ID is required").WriteToHTTP(w)
		return
	}

	if req.Role != spaces.RoleMember && req.Role != spaces.RoleAdmin {
		problems.ValidationError("role", "Role must be 'member' or 'admin'").WriteToHTTP(w)
		return
	}

	// check if user is already a member
	for _, m := range s.Members {
		if m.UserID == req.UserID {
			problems.AlreadyExists("User", req.UserID).WriteToHTTP(w)
			return
		}
	}

	s.Members = append(s.Members, req)

	if err := h.store.UpdateFull(s); err != nil {
		slog.Error("Failed to update space", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) handleUpdateMember(w http.ResponseWriter, r *http.Request, s *spaces.Space, uid string) {
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}

	if !s.IsOwner(u) {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	if ct := r.Header.Get("Content-Type"); ct != "application/json" {
		problems.WrongContentType("application/json", ct).WriteToHTTP(w)
		return
	}

	var req struct {
		Role spaces.Role `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}

	if req.Role != spaces.RoleMember && req.Role != spaces.RoleAdmin {
		problems.ValidationError("role", "Role must be 'member' or 'admin'").WriteToHTTP(w)
		return
	}

	found := false
	for i, m := range s.Members {
		if m.UserID == uid {
			s.Members[i].Role = req.Role
			found = true
			break
		}
	}

	if !found {
		problems.NotFound("Member", uid).WriteToHTTP(w)
		return
	}

	if err := h.store.UpdateFull(s); err != nil {
		slog.Error("Failed to update space", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleDeleteMember(w http.ResponseWriter, r *http.Request, s *spaces.Space, uid string) {
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}

	// only the user themselves or an owner can remove a member
	if u.ID != uid && !s.IsOwner(u) {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	found := false
	for i, m := range s.Members {
		if m.UserID == uid {
			s.Members = append(s.Members[:i], s.Members[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		problems.NotFound("Member", uid).WriteToHTTP(w)
		return
	}

	if err := h.store.UpdateFull(s); err != nil {
		slog.Error("Failed to update space", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
