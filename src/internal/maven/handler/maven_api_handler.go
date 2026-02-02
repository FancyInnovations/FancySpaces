package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/OliverSchlueter/goutils/problems"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/internal/maven"
	"github.com/fancyinnovations/fancyspaces/internal/spaces"
)

func (h *Handler) handleRepositories(w http.ResponseWriter, r *http.Request) {
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

	if !space.MavenRepositorySettings.Enabled {
		spaces.ProblemFeatureNotEnabled("releases").WriteToHTTP(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetRepositories(w, r, space)
	case http.MethodPost:
		h.handleCreateRepository(w, r, space)
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodPost}).WriteToHTTP(w)
	}
}

// no auth required
func (h *Handler) handleGetRepositories(w http.ResponseWriter, r *http.Request, space *spaces.Space) {
	all, err := h.store.GetRepositories(r.Context(), space.ID)
	if err != nil {
		slog.Error("Failed to get maven repositories", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	u := h.userFromCtx(r.Context())

	filtered := []maven.Repository{}
	for _, repo := range all {
		if repo.Public {
			filtered = append(filtered, repo)
		} else {
			if u != nil && u.Verified && u.IsActive && space.IsMember(u) {
				filtered = append(filtered, repo)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=600") // 10 minutes
	json.NewEncoder(w).Encode(filtered)
}

func (h *Handler) handleCreateRepository(w http.ResponseWriter, r *http.Request, space *spaces.Space) {
	// TODO: implement
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) handleRepository(w http.ResponseWriter, r *http.Request) {
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

	if !space.MavenRepositorySettings.Enabled {
		spaces.ProblemFeatureNotEnabled("releases").WriteToHTTP(w)
		return
	}

	repoName := r.PathValue("repository_name")
	repo, err := h.store.GetRepository(r.Context(), space.ID, repoName)
	if err != nil {
		if errors.Is(err, maven.ErrRepositoryNotFound) {
			problems.NotFound("Maven Repository", repoName).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get maven repository", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	if !repo.Public {
		u := h.userFromCtx(r.Context())
		if u == nil || !u.Verified || !u.IsActive || !space.IsMember(u) {
			problems.NotFound("Maven Repository", repo.Name).WriteToHTTP(w)
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetRepository(w, r, space, repo)
	case http.MethodPut:
		h.handleUpdateRepository(w, r, space, repo)
	case http.MethodDelete:
		h.handleDeleteRepository(w, r, space, repo)
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodGet, http.MethodDelete}).WriteToHTTP(w)
	}
}

// no auth required
func (h *Handler) handleGetRepository(w http.ResponseWriter, r *http.Request, space *spaces.Space, repo *maven.Repository) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=600") // 10 minutes
	json.NewEncoder(w).Encode(repo)
}

func (h *Handler) handleUpdateRepository(w http.ResponseWriter, r *http.Request, space *spaces.Space, repo *maven.Repository) {
	// TODO: implement
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) handleDeleteRepository(w http.ResponseWriter, r *http.Request, space *spaces.Space, repo *maven.Repository) {
	// TODO: implement
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) handleArtifacts(w http.ResponseWriter, r *http.Request) {
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

	if !space.MavenRepositorySettings.Enabled {
		spaces.ProblemFeatureNotEnabled("releases").WriteToHTTP(w)
		return
	}

	repoName := r.PathValue("repository_name")
	repo, err := h.store.GetRepository(r.Context(), space.ID, repoName)
	if err != nil {
		if errors.Is(err, maven.ErrRepositoryNotFound) {
			problems.NotFound("Maven Repository", repoName).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get maven repository", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	if !repo.Public {
		u := h.userFromCtx(r.Context())
		if u == nil || !u.Verified || !u.IsActive || !space.IsMember(u) {
			problems.NotFound("Maven Repository", repo.Name).WriteToHTTP(w)
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetArtifacts(w, r, space, repo)
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet}).WriteToHTTP(w)
	}
}

func (h *Handler) handleGetArtifacts(w http.ResponseWriter, r *http.Request, space *spaces.Space, repo *maven.Repository) {
	all, err := h.store.GetArtifacts(r.Context(), space.ID, repo.Name)
	if err != nil {
		slog.Error("Failed to get maven artifacts", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60") // 1 minute
	json.NewEncoder(w).Encode(all)
}

func (h *Handler) handleArtifact(w http.ResponseWriter, r *http.Request) {
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

	if !space.MavenRepositorySettings.Enabled {
		spaces.ProblemFeatureNotEnabled("releases").WriteToHTTP(w)
		return
	}

	repoName := r.PathValue("repository_name")
	repo, err := h.store.GetRepository(r.Context(), space.ID, repoName)
	if err != nil {
		if errors.Is(err, maven.ErrRepositoryNotFound) {
			problems.NotFound("Maven Repository", repoName).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get maven repository", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	if !repo.Public {
		u := h.userFromCtx(r.Context())
		if u == nil || !u.Verified || !u.IsActive || !space.IsMember(u) {
			problems.NotFound("Maven Repository", repo.Name).WriteToHTTP(w)
			return
		}
	}

	groupArtifactID := r.PathValue("group_artifact_id") // {groupId}:{artifactId}
	groupID := strings.SplitN(groupArtifactID, ":", 2)[0]
	artifactID := strings.SplitN(groupArtifactID, ":", 2)[1]
	artifact, err := h.store.GetArtifact(r.Context(), space.ID, repo.Name, groupID, artifactID)
	if err != nil {
		if errors.Is(err, maven.ErrArtifactNotFound) {
			problems.NotFound("Maven Artifact", groupArtifactID).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get maven artifact", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetArtifact(w, r, space, repo, artifact)
	case http.MethodDelete:
		h.handleDeleteArtifact(w, r, space, repo, artifact)
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodDelete}).WriteToHTTP(w)
	}
}

// no auth required
func (h *Handler) handleGetArtifact(w http.ResponseWriter, r *http.Request, space *spaces.Space, repo *maven.Repository, artifact *maven.Artifact) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60") // 1 minute
	json.NewEncoder(w).Encode(artifact)
}

func (h *Handler) handleDeleteArtifact(w http.ResponseWriter, r *http.Request, space *spaces.Space, repo *maven.Repository, artifact *maven.Artifact) {
	// TODO: implement
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) handleJavadoc(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	w.WriteHeader(http.StatusNotImplemented)
}
