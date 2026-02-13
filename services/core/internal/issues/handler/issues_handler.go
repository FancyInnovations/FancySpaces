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
	"github.com/fancyinnovations/fancyspaces/core/internal/auth"
	"github.com/fancyinnovations/fancyspaces/core/internal/issues"
	"github.com/fancyinnovations/fancyspaces/core/internal/spaces"
)

type Handler struct {
	store       *issues.Store
	spaces      *spaces.Store
	userFromCtx func(ctx context.Context) *auth.User
}

type Configuration struct {
	Store       *issues.Store
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
	mux.HandleFunc(prefix+"/spaces/{space_id}/issues", h.handleIssues)
	mux.HandleFunc(prefix+"/spaces/{space_id}/issues/{issue_id}", h.handleIssue)
}

func (h *Handler) handleIssues(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("space_id")
	if sid == "" {
		http.Error(w, "missing space_id", http.StatusBadRequest)
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

	if !space.ReleaseSettings.Enabled {
		spaces.ProblemFeatureNotEnabled("releases").WriteToHTTP(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleListIssues(w, r, space)
	case http.MethodPost:
		h.handleCreateIssue(w, r, space)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleIssue(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("space_id")
	if sid == "" {
		http.Error(w, "missing space_id", http.StatusBadRequest)
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

	if !space.ReleaseSettings.Enabled {
		spaces.ProblemFeatureNotEnabled("releases").WriteToHTTP(w)
		return
	}

	iid := r.PathValue("issue_id")
	if iid == "" {
		http.Error(w, "missing issue_id", http.StatusBadRequest)
		return
	}
	issue, err := h.store.GetIssue(space.ID, iid)
	if err != nil {
		if errors.Is(err, issues.ErrIssueNotFound) {
			problems.NotFound("Issue", iid).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get issue by id", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetIssue(w, r, space, issue)
	case http.MethodPut:
		h.handleUpdateIssue(w, r, space, issue)
	case http.MethodDelete:
		h.handleDeleteIssue(w, r, space, issue)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// no auth required
func (h *Handler) handleListIssues(w http.ResponseWriter, r *http.Request, space *spaces.Space) {
	all, err := h.store.GetIssues(space.ID)
	if err != nil {
		slog.Error("Failed to list issues", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=5") // 5 seconds
	json.NewEncoder(w).Encode(all)
}

func (h *Handler) handleCreateIssue(w http.ResponseWriter, r *http.Request, space *spaces.Space) {
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}
	if !space.HasWriteAccess(u) {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	var issue issues.Issue
	if err := json.NewDecoder(r.Body).Decode(&issue); err != nil {
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}
	issue.Space = space.ID
	issue.Reporter = u.ID
	issue.Status = issues.StatusBacklog
	issue.ExternalSource = ""

	if err := h.store.CreateIssue(&issue); err != nil {
		if errors.Is(err, issues.ErrIssueAlreadyExists) {
			problems.AlreadyExists("Issue", issue.ID).WriteToHTTP(w)
			return
		}
		slog.Error("Failed to create issue", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(issue)
}

// no auth required
func (h *Handler) handleGetIssue(w http.ResponseWriter, r *http.Request, space *spaces.Space, issue *issues.Issue) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=5") // 5 seconds
	json.NewEncoder(w).Encode(issue)
}

func (h *Handler) handleUpdateIssue(w http.ResponseWriter, r *http.Request, space *spaces.Space, issue *issues.Issue) {
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}
	if !space.HasWriteAccess(u) {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	var req issues.Issue
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}

	// If status is being changed to done or closed, set resolved at time
	if (req.Status == issues.StatusDone || req.Status == issues.StatusClosed) && (issue.Status != issues.StatusDone && issue.Status != issues.StatusClosed) {
		now := time.Now()
		issue.ResolvedAt = &now
	}
	if (req.Status != issues.StatusDone && req.Status != issues.StatusClosed) && (issue.Status == issues.StatusDone || issue.Status == issues.StatusClosed) {
		issue.ResolvedAt = nil
	}

	issue.Title = req.Title
	issue.Description = req.Description
	issue.Type = req.Type
	issue.Status = req.Status
	issue.Priority = req.Priority
	issue.Assignee = req.Assignee
	issue.FixVersion = req.FixVersion
	issue.AffectedVersions = req.AffectedVersions
	issue.ParentIssue = req.ParentIssue
	issue.ExtraFields = req.ExtraFields

	if err := h.store.UpdateIssue(issue); err != nil {
		slog.Error("Failed to update issue", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(issue)
}

func (h *Handler) handleDeleteIssue(w http.ResponseWriter, r *http.Request, space *spaces.Space, issue *issues.Issue) {
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}
	if !space.HasWriteAccess(u) {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	if err := h.store.DeleteIssue(space.ID, issue.ID); err != nil {
		slog.Error("Failed to delete issue", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
