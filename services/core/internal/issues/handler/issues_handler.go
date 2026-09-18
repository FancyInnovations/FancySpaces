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
	Title            string                `json:"title"`
	Description      string                `json:"description"`
	Type             issues.Type           `json:"type"`
	Status           issues.Status         `json:"status"`
	Priority         issues.Priority       `json:"priority"`
	Assignee         string                `json:"assignee"`
	FixVersion       string                `json:"fix_version"`
	AffectedVersions []string              `json:"affected_versions"`
	ParentIssue      string                `json:"parent_issue"`
	Labels           []string              `json:"labels"`
	Relationships    []issues.Relationship `json:"relationships"`
}

type issuePatchRequest struct {
	Title            *string                `json:"title"`
	Description      *string                `json:"description"`
	Type             *issues.Type           `json:"type"`
	Status           *issues.Status         `json:"status"`
	Priority         *issues.Priority       `json:"priority"`
	Assignee         *string                `json:"assignee"`
	FixVersion       *string                `json:"fix_version"`
	AffectedVersions *[]string              `json:"affected_versions"`
	ParentIssue      *string                `json:"parent_issue"`
	Labels           *[]string              `json:"labels"`
	Relationships    *[]issues.Relationship `json:"relationships"`
}

type issueListResponse struct {
	Items       []issues.Issue   `json:"items"`
	Total       int              `json:"total"`
	Offset      int              `json:"offset"`
	Limit       int              `json:"limit"`
	Permissions issuePermissions `json:"permissions"`
}

type issuePermissions struct {
	CanWrite   bool `json:"can_write"`
	CanArchive bool `json:"can_archive"`
}

func New(cfg Configuration) *Handler {
	return &Handler{store: cfg.Store, spaces: cfg.Spaces, userFromCtx: cfg.UserFromCtx}
}

func (h *Handler) Register(prefix string, mux *http.ServeMux) {
	mux.HandleFunc(prefix+"/spaces/{space_id}/issues", h.handleIssues)
	mux.HandleFunc(prefix+"/spaces/{space_id}/issues/{issue_id}", h.handleIssue)
	mux.HandleFunc(prefix+"/spaces/{space_id}/issues/{issue_id}/comments", h.handleComments)
	mux.HandleFunc(prefix+"/spaces/{space_id}/issues/{issue_id}/comments/{comment_id}", h.handleComment)
	mux.HandleFunc(prefix+"/spaces/{space_id}/issues/{issue_id}/activity", h.handleActivity)
	mux.HandleFunc(prefix+"/spaces/{space_id}/issues/bulk", h.handleBulkUpdate)
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
		ExternalSource: issues.ExternalSource(q.Get("external_source")), Label: q.Get("label"), Offset: offset, Limit: limit,
	}
	if opts.Type != "" && !issues.ValidType(opts.Type) || opts.Status != "" && !issues.ValidStatus(opts.Status) || opts.Priority != "" && !issues.ValidPriority(opts.Priority) {
		problems.ValidationError("filter", "Invalid issue filter").WriteToHTTP(w)
		return
	}
	list, total, err := h.store.ListIssues(space.ID, opts)
	if err != nil {
		slog.Error("Failed to list issues", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}
	u := h.userFromCtx(r.Context())
	canWrite := u != nil && u.Verified && u.IsActive && space.HasWriteAccess(u)
	h.writeJSON(w, http.StatusOK, issueListResponse{Items: list, Total: total, Offset: offset, Limit: limit, Permissions: issuePermissions{CanWrite: canWrite, CanArchive: canWrite}})
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
		FixVersion: req.FixVersion, AffectedVersions: req.AffectedVersions, ParentIssue: req.ParentIssue, Labels: req.Labels, Relationships: req.Relationships}
	if err := h.validateIssueWrite(space, &issue); err != nil {
		h.writeIssueValidationError(w, err)
		return
	}
	if err := h.store.CreateIssue(&issue); err != nil {
		h.writeIssueValidationError(w, err)
		return
	}
	h.recordActivity(&issue, u.ID, "created", "", "", "")
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
		updated.Assignee, updated.FixVersion, updated.AffectedVersions, updated.ParentIssue, updated.Labels, updated.Relationships = req.Assignee, req.FixVersion, req.AffectedVersions, req.ParentIssue, req.Labels, req.Relationships
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
	h.recordChanges(issue, &updated, u.ID)
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
	if req.Labels != nil {
		issue.Labels = *req.Labels
	}
	if req.Relationships != nil {
		issue.Relationships = *req.Relationships
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
		parentID := issue.ParentIssue
		seen := map[string]struct{}{issue.ID: {}}
		for depth := 0; depth < 50 && parentID != ""; depth++ {
			if _, cycle := seen[parentID]; cycle {
				return issues.ErrInvalidParent
			}
			seen[parentID] = struct{}{}
			parent, err := h.store.GetIssue(space.ID, parentID)
			if err != nil {
				return issues.ErrInvalidParent
			}
			parentID = parent.ParentIssue
		}
		if parentID != "" {
			return issues.ErrInvalidParent
		}
	}
	for _, relationship := range issue.Relationships {
		if _, err := h.store.GetIssue(space.ID, relationship.Issue); err != nil {
			return issues.ErrInvalidRelationship
		}
	}
	return nil
}

func (h *Handler) recordActivity(issue *issues.Issue, actor, kind, field, oldValue, newValue string) {
	if err := h.store.AddActivity(&issues.Activity{Space: issue.Space, Issue: issue.ID, Actor: actor, Kind: kind, Field: field, OldValue: oldValue, NewValue: newValue}); err != nil {
		slog.Warn("failed to record issue activity", sloki.WrapError(err))
	}
}

func (h *Handler) recordChanges(before, after *issues.Issue, actor string) {
	changes := []struct{ field, old, new string }{
		{"title", before.Title, after.Title}, {"status", string(before.Status), string(after.Status)},
		{"priority", string(before.Priority), string(after.Priority)}, {"assignee", before.Assignee, after.Assignee},
		{"parent_issue", before.ParentIssue, after.ParentIssue},
		{"labels", strings.Join(before.Labels, ", "), strings.Join(after.Labels, ", ")},
	}
	for _, change := range changes {
		if change.old != change.new {
			h.recordActivity(after, actor, "changed", change.field, change.old, change.new)
		}
	}
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
	h.recordActivity(issue, u.ID, "commented", "", "", "")
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
		h.recordActivity(&issues.Issue{ID: issueID, Space: space.ID}, u.ID, "deleted_comment", "", "", "")
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
	h.recordActivity(&issues.Issue{ID: issueID, Space: space.ID}, u.ID, "edited_comment", "", "", "")
	h.writeJSON(w, http.StatusOK, current)
}

func (h *Handler) handleActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet}).WriteToHTTP(w)
		return
	}
	space, ok := h.loadSpace(w, r)
	if !ok {
		return
	}
	issueID := r.PathValue("issue_id")
	if _, err := h.store.GetIssue(space.ID, issueID); err != nil {
		h.writeStoreError(w, "Issue", issueID, err)
		return
	}
	activities, err := h.store.GetActivities(space.ID, issueID)
	if err != nil {
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}
	h.writeJSON(w, http.StatusOK, activities)
}

// handleBulkUpdate intentionally only accepts the same controlled fields as PATCH.
// This keeps bulk triage from becoming a way to overwrite reporter or external-sync data.
func (h *Handler) handleBulkUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		problems.MethodNotAllowed(r.Method, []string{http.MethodPatch}).WriteToHTTP(w)
		return
	}
	space, ok := h.loadSpace(w, r)
	if !ok {
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
		IDs     []string          `json:"ids"`
		Changes issuePatchRequest `json:"changes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}
	if len(req.IDs) == 0 || len(req.IDs) > 100 {
		problems.ValidationError("ids", "Select between 1 and 100 issues").WriteToHTTP(w)
		return
	}
	type pendingUpdate struct{ before, after *issues.Issue }
	pending := make([]pendingUpdate, 0, len(req.IDs))
	for _, id := range req.IDs {
		issue, err := h.store.GetIssue(space.ID, id)
		if err != nil {
			h.writeStoreError(w, "Issue", id, err)
			return
		}
		updated := *issue
		applyPatch(&updated, req.Changes)
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
		pending = append(pending, pendingUpdate{before: issue, after: &updated})
	}
	updatedIssues := make([]issues.Issue, 0, len(pending))
	for _, update := range pending {
		if err := h.store.UpdateIssue(update.after); err != nil {
			h.writeIssueValidationError(w, err)
			return
		}
		h.recordChanges(update.before, update.after, u.ID)
		updatedIssues = append(updatedIssues, *update.after)
	}
	h.writeJSON(w, http.StatusOK, updatedIssues)
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
