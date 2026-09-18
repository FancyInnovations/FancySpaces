package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/OliverSchlueter/goutils/problems"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/core/internal/issues"
	spacesStore "github.com/fancyinnovations/fancyspaces/core/internal/spaces"
	"github.com/fancyinnovations/fancyspaces/integrations/idp-go-sdk/idp"
	"github.com/fancyinnovations/fancyspaces/integrations/spaces-go-sdk/spaces"
)

type Handler struct {
	store       *issues.Store
	spaces      *spacesStore.Store
	userFromCtx func(ctx context.Context) *idp.User
}

type Configuration struct {
	Store       *issues.Store
	Spaces      *spacesStore.Store
	UserFromCtx func(ctx context.Context) *idp.User
}

type issueWriteRequest struct {
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	Type             issues.Type            `json:"type"`
	Status           issues.Status          `json:"status"`
	Priority         issues.Priority        `json:"priority"`
	Assignee         string                 `json:"assignee"`
	FixVersion       string                 `json:"fix_version"`
	AffectedVersions []string               `json:"affected_versions"`
	ParentIssue      string                 `json:"parent_issue"`
	ExtraFields      map[string]interface{} `json:"extra_fields"`
}

type issuePatchRequest struct {
	Title            *string                 `json:"title"`
	Description      *string                 `json:"description"`
	Type             *issues.Type            `json:"type"`
	Status           *issues.Status          `json:"status"`
	Priority         *issues.Priority        `json:"priority"`
	Assignee         *string                 `json:"assignee"`
	FixVersion       *string                 `json:"fix_version"`
	AffectedVersions *[]string               `json:"affected_versions"`
	ParentIssue      *string                 `json:"parent_issue"`
	ExtraFields      *map[string]interface{} `json:"extra_fields"`
}

type issueListResponse struct {
	Items  []issues.Issue `json:"items"`
	Total  int            `json:"total"`
	Offset int            `json:"offset"`
	Limit  int            `json:"limit"`
}

func New(cfg Configuration) *Handler {
	return &Handler{store: cfg.Store, spaces: cfg.Spaces, userFromCtx: cfg.UserFromCtx}
}

func (h *Handler) Register(prefix string, mux *http.ServeMux) {
	mux.HandleFunc(prefix+"/spaces/{space_id}/issues", h.handleIssues)
	mux.HandleFunc(prefix+"/spaces/{space_id}/issues/{issue_id}", h.handleIssue)
	mux.HandleFunc(prefix+"/spaces/{space_id}/issues/{issue_id}/comments", h.handleComments)
	mux.HandleFunc(prefix+"/spaces/{space_id}/issues/{issue_id}/comments/{comment_id}", h.handleComment)
}

func (h *Handler) loadSpace(w http.ResponseWriter, r *http.Request) (*spaces.Space, bool) {
	sid := r.PathValue("space_id")
	if sid == "" {
		problems.ValidationError("space_id", "Space ID is required").WriteToHTTP(w)
		return nil, false
	}
	space, err := h.spaces.Get(sid)
	if err != nil {
		if errors.Is(err, spaces.ErrSpaceNotFound) {
			problems.NotFound("Space", sid).WriteToHTTP(w)
			return nil, false
		}
		slog.Error("Failed to get space by id", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return nil, false
	}
	if space.Status != spaces.StatusApproved && space.Status != spaces.StatusArchived {
		u := h.userFromCtx(r.Context())
		if u == nil || !u.Verified || !u.IsActive || !space.IsMember(u) {
			problems.NotFound("Space", space.ID).WriteToHTTP(w)
			return nil, false
		}
	}
	if !space.IssueSettings.Enabled {
		spacesStore.ProblemFeatureNotEnabled("issues").WriteToHTTP(w)
		return nil, false
	}
	return space, true
}

func (h *Handler) handleIssues(w http.ResponseWriter, r *http.Request) {
	space, ok := h.loadSpace(w, r)
	if !ok {
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.handleListIssues(w, r, space)
	case http.MethodPost:
		h.handleCreateIssue(w, r, space)
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodPost}).WriteToHTTP(w)
	}
}

func (h *Handler) handleIssue(w http.ResponseWriter, r *http.Request) {
	space, ok := h.loadSpace(w, r)
	if !ok {
		return
	}
	iid := r.PathValue("issue_id")
	if iid == "" {
		problems.ValidationError("issue_id", "Issue ID is required").WriteToHTTP(w)
		return
	}
	issue, err := h.store.GetIssue(space.ID, iid)
	if err != nil {
		h.writeStoreError(w, "Issue", iid, err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.writeJSON(w, http.StatusOK, issue)
	case http.MethodPut:
		h.handleUpdateIssue(w, r, space, issue, false)
	case http.MethodPatch:
		h.handleUpdateIssue(w, r, space, issue, true)
	case http.MethodDelete:
		h.handleDeleteIssue(w, r, space, issue)
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodDelete}).WriteToHTTP(w)
	}
}

func (h *Handler) handleListIssues(w http.ResponseWriter, r *http.Request, space *spaces.Space) {
	q := r.URL.Query()
	limit, offset, err := parsePagination(q.Get("limit"), q.Get("offset"))
	if err != nil {
		problems.ValidationError("pagination", err.Error()).WriteToHTTP(w)
		return
	}
	opts := issues.ListOptions{
		Query: q.Get("q"), Type: issues.Type(q.Get("type")), Status: issues.Status(q.Get("status")),
		Priority: issues.Priority(q.Get("priority")), Assignee: q.Get("assignee"),
		ExternalSource: issues.ExternalSource(q.Get("external_source")), Offset: offset, Limit: limit,
	}
	list, total, err := h.store.ListIssues(space.ID, opts)
	if err != nil {
		slog.Error("Failed to list issues", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}
	h.writeJSON(w, http.StatusOK, issueListResponse{Items: list, Total: total, Offset: offset, Limit: limit})
}

func parsePagination(limitValue, offsetValue string) (int, int, error) {
	limit, offset := 100, 0
	var err error
	if limitValue != "" {
		limit, err = strconv.Atoi(limitValue)
		if err != nil || limit < 1 || limit > 200 {
			return 0, 0, errors.New("limit must be between 1 and 200")
		}
	}
	if offsetValue != "" {
		offset, err = strconv.Atoi(offsetValue)
		if err != nil || offset < 0 {
			return 0, 0, errors.New("offset must be zero or greater")
		}
	}
	return limit, offset, nil
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

	var req issueWriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}
	issue := issues.Issue{Space: space.ID, Reporter: u.ID, Title: req.Title, Description: req.Description,
		Type: req.Type, Status: issues.StatusBacklog, Priority: req.Priority, Assignee: req.Assignee,
		FixVersion: req.FixVersion, AffectedVersions: req.AffectedVersions, ParentIssue: req.ParentIssue, ExtraFields: req.ExtraFields}
	if err := h.validateIssueWrite(space, &issue); err != nil {
		h.writeIssueValidationError(w, err)
		return
	}
	if err := h.store.CreateIssue(&issue); err != nil {
		h.writeIssueValidationError(w, err)
		return
	}
	h.writeJSON(w, http.StatusCreated, issue)
}

func (h *Handler) handleUpdateIssue(w http.ResponseWriter, r *http.Request, space *spaces.Space, issue *issues.Issue, partial bool) {
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}
	if !space.HasWriteAccess(u) {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	updated := *issue
	if partial {
		var req issuePatchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
			return
		}
		applyPatch(&updated, req)
	} else {
		var req issueWriteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
			return
		}
		updated.Title, updated.Description, updated.Type, updated.Status, updated.Priority = req.Title, req.Description, req.Type, req.Status, req.Priority
		updated.Assignee, updated.FixVersion, updated.AffectedVersions, updated.ParentIssue, updated.ExtraFields = req.Assignee, req.FixVersion, req.AffectedVersions, req.ParentIssue, req.ExtraFields
	}
	if updated.Status != issue.Status && (updated.Status == issues.StatusDone || updated.Status == issues.StatusClosed) {
		now := time.Now()
		updated.ResolvedAt = &now
	} else if updated.Status != issues.StatusDone && updated.Status != issues.StatusClosed {
		updated.ResolvedAt = nil
	}
	if err := h.validateIssueWrite(space, &updated); err != nil {
		h.writeIssueValidationError(w, err)
		return
	}
	if err := h.store.UpdateIssue(&updated); err != nil {
		h.writeIssueValidationError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, updated)
}

func applyPatch(issue *issues.Issue, req issuePatchRequest) {
	if req.Title != nil {
		issue.Title = *req.Title
	}
	if req.Description != nil {
		issue.Description = *req.Description
	}
	if req.Type != nil {
		issue.Type = *req.Type
	}
	if req.Status != nil {
		issue.Status = *req.Status
	}
	if req.Priority != nil {
		issue.Priority = *req.Priority
	}
	if req.Assignee != nil {
		issue.Assignee = *req.Assignee
	}
	if req.FixVersion != nil {
		issue.FixVersion = *req.FixVersion
	}
	if req.AffectedVersions != nil {
		issue.AffectedVersions = *req.AffectedVersions
	}
	if req.ParentIssue != nil {
		issue.ParentIssue = *req.ParentIssue
	}
	if req.ExtraFields != nil {
		issue.ExtraFields = *req.ExtraFields
	}
}

func (h *Handler) validateIssueWrite(space *spaces.Space, issue *issues.Issue) error {
	if err := issue.Validate(); err != nil {
		return err
	}
	if issue.ParentIssue != "" {
		if issue.ParentIssue == issue.ID {
			return issues.ErrInvalidParent
		}
		parent, err := h.store.GetIssue(space.ID, issue.ParentIssue)
		if err != nil || parent.ID == issue.ID {
			return issues.ErrInvalidParent
		}
	}
	return nil
}

func (h *Handler) handleDeleteIssue(w http.ResponseWriter, r *http.Request, space *spaces.Space, issue *issues.Issue) {
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}
	if !space.HasWriteAccess(u) || issue.ExternalSource != "" {
		problems.Forbidden().WriteToHTTP(w)
		return
	}
	if err := h.store.DeleteIssue(space.ID, issue.ID); err != nil {
		h.writeStoreError(w, "Issue", issue.ID, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleComments(w http.ResponseWriter, r *http.Request) {
	space, ok := h.loadSpace(w, r)
	if !ok {
		return
	}
	issueID := r.PathValue("issue_id")
	issue, err := h.store.GetIssue(space.ID, issueID)
	if err != nil {
		h.writeStoreError(w, "Issue", issueID, err)
		return
	}
	if r.Method == http.MethodGet {
		comments, err := h.store.GetComments(space.ID, issue.ID)
		if err != nil {
			problems.InternalServerError("").WriteToHTTP(w)
			return
		}
		h.writeJSON(w, http.StatusOK, comments)
		return
	}
	if r.Method != http.MethodPost {
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodPost}).WriteToHTTP(w)
		return
	}
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}
	if !space.HasWriteAccess(u) {
		problems.Forbidden().WriteToHTTP(w)
		return
	}
	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}
	comment := &issues.Comment{Space: space.ID, Issue: issue.ID, Author: u.ID, Content: strings.TrimSpace(req.Content)}
	if err := h.store.AddComment(comment); err != nil {
		h.writeIssueValidationError(w, err)
		return
	}
	h.writeJSON(w, http.StatusCreated, comment)
}

func (h *Handler) handleComment(w http.ResponseWriter, r *http.Request) {
	space, ok := h.loadSpace(w, r)
	if !ok {
		return
	}
	issueID, commentID := r.PathValue("issue_id"), r.PathValue("comment_id")
	if r.Method != http.MethodPatch && r.Method != http.MethodDelete {
		problems.MethodNotAllowed(r.Method, []string{http.MethodPatch, http.MethodDelete}).WriteToHTTP(w)
		return
	}
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}
	if !space.HasWriteAccess(u) {
		problems.Forbidden().WriteToHTTP(w)
		return
	}
	comments, err := h.store.GetComments(space.ID, issueID)
	if err != nil {
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}
	var current *issues.Comment
	for i := range comments {
		if comments[i].ID == commentID {
			current = &comments[i]
			break
		}
	}
	if current == nil {
		problems.NotFound("Comment", commentID).WriteToHTTP(w)
		return
	}
	if current.Author != u.ID && !space.HasFullAccess(u) {
		problems.Forbidden().WriteToHTTP(w)
		return
	}
	if r.Method == http.MethodDelete {
		if err := h.store.DeleteComment(space.ID, issueID, commentID); err != nil {
			h.writeStoreError(w, "Comment", commentID, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}
	current.Content = strings.TrimSpace(req.Content)
	if current.Content == "" || len(current.Content) > 10_000 {
		problems.ValidationError("content", "Comment must be between 1 and 10000 characters").WriteToHTTP(w)
		return
	}
	if err := h.store.UpdateComment(current); err != nil {
		h.writeStoreError(w, "Comment", commentID, err)
		return
	}
	h.writeJSON(w, http.StatusOK, current)
}

func (h *Handler) writeIssueValidationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, issues.ErrIssueNotFound), errors.Is(err, issues.ErrIssueArchived):
		problems.NotFound("Issue", "").WriteToHTTP(w)
	case errors.Is(err, issues.ErrInvalidParent), errors.Is(err, issues.ErrTitleTooShort), errors.Is(err, issues.ErrTitleTooLong), errors.Is(err, issues.ErrDescriptionTooLong), errors.Is(err, issues.ErrInvalidType), errors.Is(err, issues.ErrInvalidStatus), errors.Is(err, issues.ErrInvalidPriority), errors.Is(err, issues.ErrCommentTooLong):
		problems.ValidationError("body", err.Error()).WriteToHTTP(w)
	default:
		slog.Error("Failed to write issue", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
	}
}

func (h *Handler) writeStoreError(w http.ResponseWriter, resource, id string, err error) {
	if errors.Is(err, issues.ErrIssueNotFound) || errors.Is(err, issues.ErrCommentNotFound) {
		problems.NotFound(resource, id).WriteToHTTP(w)
		return
	}
	slog.Error("Failed to access issue resource", sloki.WrapError(err))
	problems.InternalServerError("").WriteToHTTP(w)
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
