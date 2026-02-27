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
	"github.com/fancyinnovations/fancyspaces/core/internal/blogs"
	spacesStore "github.com/fancyinnovations/fancyspaces/core/internal/spaces"
	"github.com/fancyinnovations/fancyspaces/integrations/idp-go-sdk/idp"
	"github.com/fancyinnovations/fancyspaces/integrations/spaces-go-sdk/spaces"
)

type Handler struct {
	store       *blogs.Store
	spaces      *spacesStore.Store
	userFromCtx func(ctx context.Context) *idp.User
}

type Configuration struct {
	Store       *blogs.Store
	Spaces      *spacesStore.Store
	UserFromCtx func(ctx context.Context) *idp.User
}

func New(cfg Configuration) *Handler {
	return &Handler{
		store:       cfg.Store,
		spaces:      cfg.Spaces,
		userFromCtx: cfg.UserFromCtx,
	}
}

func (h *Handler) Register(prefix string, mux *http.ServeMux) {
	mux.HandleFunc(prefix+"/blog-articles", h.handleArticles)

	mux.HandleFunc(prefix+"/blog-articles/{article_id}", h.handleArticle)
	mux.HandleFunc(prefix+"/blog-articles/{article_id}/content", h.handleGetArticleContent)
	mux.HandleFunc(prefix+"/spaces/{space_id}/blog-articles", h.handleArticlesForSpace)
	mux.HandleFunc(prefix+"/users/{user_id}/blog-articles", h.handleArticlesForUser)
}

func (h *Handler) handleArticles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		problems.MethodNotAllowed(r.Method, []string{http.MethodPost}).WriteToHTTP(w)
		return
	}

	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}

	var req CreateOrUpdateArticleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}

	if req.Title == "" {
		problems.ValidationError("title", "Title is required").WriteToHTTP(w)
		return
	}
	if len(req.Title) > blogs.MaxTitleSize {
		problems.ValidationError("title", "Title is too long").WriteToHTTP(w)
		return
	}

	if req.Summary == "" {
		problems.ValidationError("summary", "Summary is required").WriteToHTTP(w)
		return
	}
	if len(req.Summary) > blogs.MaxSummarySize {
		problems.ValidationError("summary", "Summary is too long").WriteToHTTP(w)
		return
	}

	if req.Content == "" {
		problems.ValidationError("content", "Content is required").WriteToHTTP(w)
		return
	}
	if len(req.Content) > blogs.MaxContentSize {
		problems.ValidationError("content", "Content is too long").WriteToHTTP(w)
		return
	}

	if req.SpaceID != "" {
		space, err := h.spaces.Get(req.SpaceID)
		if err != nil {
			if errors.Is(err, spaces.ErrSpaceNotFound) {
				problems.NotFound("Space", req.SpaceID).WriteToHTTP(w)
				return
			}

			slog.Error("Failed to get space by id", sloki.WrapError(err))
			problems.InternalServerError("").WriteToHTTP(w)
			return
		}
		if space.Status != spaces.StatusApproved && space.Status != spaces.StatusArchived {
			if !space.IsMember(u) {
				problems.NotFound("Space", space.ID).WriteToHTTP(w)
				return
			}
		}
	}

	a := &blogs.Article{
		ID:          "", // ID will be generated
		SpaceID:     req.SpaceID,
		Author:      u.ID,
		Title:       req.Title,
		Summary:     req.Summary,
		Content:     req.Content,
		PublishedAt: time.Now(), // PublishedAt will be generated
	}

	if err := h.store.CreateArticle(a, req.Content); err != nil {
		slog.Error("Failed to create article", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(a)
}

func (h *Handler) handleArticle(w http.ResponseWriter, r *http.Request) {
	articleID := r.PathValue("article_id")
	if articleID == "" {
		problems.ValidationError("article_id", "Article ID is required").WriteToHTTP(w)
		return
	}

	a, err := h.store.GetArticleByID(articleID)
	if err != nil {
		if err == blogs.ErrArticleNotFound {
			problems.NotFound("Article", articleID).WriteToHTTP(w)
			return
		}

		problems.InternalServerError("Failed to get article by id").WriteToHTTP(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetArticle(w, r, a)
	case http.MethodPut:
		h.handleUpdateArticle(w, r, a)
	case http.MethodDelete:
		h.handleDeleteArticle(w, r, a)
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodPut, http.MethodDelete}).WriteToHTTP(w)
	}
}

func (h *Handler) handleArticlesForSpace(w http.ResponseWriter, r *http.Request) {
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

	if !space.BlogSettings.Enabled {
		spacesStore.ProblemFeatureNotEnabled("blog").WriteToHTTP(w)
		return
	}

	articles, err := h.store.GetArticlesForSpace(sid)
	if err != nil {
		slog.Error("Failed to get articles for space", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	json.NewEncoder(w).Encode(articles)
}

func (h *Handler) handleArticlesForUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	if userID == "" {
		problems.ValidationError("user_id", "User ID is required").WriteToHTTP(w)
		return
	}

	articles, err := h.store.GetArticlesForUser(userID)
	if err != nil {
		slog.Error("Failed to get articles for user", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	json.NewEncoder(w).Encode(articles)
}

func (h *Handler) handleGetArticle(w http.ResponseWriter, r *http.Request, a *blogs.Article) {
	// if the article is space-owned, we need to check if the user has access to the space
	if a.SpaceID != "" {
		space, err := h.spaces.Get(a.SpaceID)
		if err != nil {
			if errors.Is(err, spaces.ErrSpaceNotFound) {
				problems.NotFound("Space", a.SpaceID).WriteToHTTP(w)
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
	}

	json.NewEncoder(w).Encode(a)
}

func (h *Handler) handleGetArticleContent(w http.ResponseWriter, r *http.Request) {
	articleID := r.PathValue("article_id")
	if articleID == "" {
		problems.ValidationError("article_id", "Article ID is required").WriteToHTTP(w)
		return
	}

	a, err := h.store.GetArticleByID(articleID)
	if err != nil {
		if errors.Is(err, blogs.ErrArticleNotFound) {
			problems.NotFound("Article", articleID).WriteToHTTP(w)
			return
		}

		problems.InternalServerError("Failed to get article by id").WriteToHTTP(w)
		return
	}

	// if the article is space-owned, we need to check if the user has access to the space
	if a.SpaceID != "" {
		space, err := h.spaces.Get(a.SpaceID)
		if err != nil {
			if errors.Is(err, spaces.ErrSpaceNotFound) {
				problems.NotFound("Space", a.SpaceID).WriteToHTTP(w)
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
	}

	content, err := h.store.GetArticleContentByID(a.ID)
	if err != nil {
		slog.Error("Failed to get article content", sloki.WrapError(err))
		problems.InternalServerError("Failed to get article content").WriteToHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(content))
}

func (h *Handler) handleUpdateArticle(w http.ResponseWriter, r *http.Request, a *blogs.Article) {
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}

	if a.SpaceID != "" {
		space, err := h.spaces.Get(a.SpaceID)
		if err != nil {
			if errors.Is(err, spaces.ErrSpaceNotFound) {
				problems.NotFound("Space", a.SpaceID).WriteToHTTP(w)
				return
			}

			slog.Error("Failed to get space by id", sloki.WrapError(err))
			problems.InternalServerError("").WriteToHTTP(w)
			return
		}

		if !space.IsOwner(u) {
			problems.Forbidden().WriteToHTTP(w)
			return
		}
	} else {
		if u.ID != a.Author {
			problems.Forbidden().WriteToHTTP(w)
			return
		}
	}

	var req CreateOrUpdateArticleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}

	if req.Title != "" {
		if len(req.Title) > blogs.MaxTitleSize {
			problems.ValidationError("title", "Title is too long").WriteToHTTP(w)
			return
		}
		a.Title = req.Title
	}

	if req.Summary != "" {
		if len(req.Summary) > blogs.MaxSummarySize {
			problems.ValidationError("summary", "Summary is too long").WriteToHTTP(w)
			return
		}
		a.Summary = req.Summary
	}

	if req.Content != "" {
		if len(req.Content) > blogs.MaxContentSize {
			problems.ValidationError("content", "Content is too long").WriteToHTTP(w)
			return
		}
		a.Content = req.Content
	}

	if err := h.store.UpdateArticle(a); err != nil {
		slog.Error("Failed to update article", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	json.NewEncoder(w).Encode(a)
}

func (h *Handler) handleDeleteArticle(w http.ResponseWriter, r *http.Request, a *blogs.Article) {
	if a.SpaceID != "" {
		space, err := h.spaces.Get(a.SpaceID)
		if err != nil {
			if errors.Is(err, spaces.ErrSpaceNotFound) {
				problems.NotFound("Space", a.SpaceID).WriteToHTTP(w)
				return
			}

			slog.Error("Failed to get space by id", sloki.WrapError(err))
			problems.InternalServerError("").WriteToHTTP(w)
			return
		}

		u := h.userFromCtx(r.Context())
		if u == nil || !u.Verified || !u.IsActive || !space.IsOwner(u) {
			problems.Forbidden().WriteToHTTP(w)
			return
		}
	} else {
		u := h.userFromCtx(r.Context())
		if u == nil || !u.Verified || !u.IsActive || u.ID != a.Author {
			problems.Forbidden().WriteToHTTP(w)
			return
		}
	}

	err := h.store.DeleteArticle(a.ID)
	if err != nil {
		slog.Error("Failed to delete article", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
