package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/OliverSchlueter/goutils/problems"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/core/internal/maven"
	"github.com/fancyinnovations/fancyspaces/core/internal/maven/javadoccache"
	spacesStore "github.com/fancyinnovations/fancyspaces/core/internal/spaces"
	"github.com/fancyinnovations/fancyspaces/integrations/spaces-go-sdk/spaces"
)

var contentTypes = map[string]string{
	".css":   "text/css",
	".js":    "application/javascript",
	".html":  "text/html",
	".json":  "application/json",
	".png":   "image/png",
	".jpg":   "image/jpeg",
	".gif":   "image/gif",
	".svg":   "image/svg+xml",
	".woff":  "font/woff",
	".woff2": "font/woff2",
	".ttf":   "font/ttf",
	".ico":   "image/x-icon",
	".webp":  "image/webp",
	".mp4":   "video/mp4",
	".mp3":   "audio/mpeg",
	".ogg":   "audio/ogg",
	".wav":   "audio/wav",
	".pdf":   "application/pdf",
	".xml":   "application/xml",
	".zip":   "application/zip",
	".tar":   "application/x-tar",
	".gz":    "application/gzip",
	".xz":    "application/x-xz",
	".rar":   "application/x-rar-compressed",
	".csv":   "text/csv",
	".txt":   "text/plain",
}

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
		spacesStore.ProblemFeatureNotEnabled("releases").WriteToHTTP(w)
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
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive || !space.HasWriteAccess(u) {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	var repo maven.Repository
	if err := json.NewDecoder(r.Body).Decode(&repo); err != nil {
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}

	if repo.Name == "" {
		problems.ValidationError("name", "Repository name is required").WriteToHTTP(w)
		return
	}

	if err := h.store.CreateRepository(r.Context(), space.ID, repo); err != nil {
		if errors.Is(err, maven.ErrRepositoryAlreadyExists) {
			problems.AlreadyExists("Maven Repository", repo.Name).WriteToHTTP(w)
			return
		}
		slog.Error("Failed to create maven repository", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(repo)
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
		spacesStore.ProblemFeatureNotEnabled("releases").WriteToHTTP(w)
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
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive || !space.HasWriteAccess(u) {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	var body maven.Repository
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}

	// Only certain fields are updatable: Public
	repo.Public = body.Public

	if err := h.store.UpdateRepository(r.Context(), space.ID, *repo); err != nil {
		slog.Error("Failed to update maven repository", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repo)
}

func (h *Handler) handleDeleteRepository(w http.ResponseWriter, r *http.Request, space *spaces.Space, repo *maven.Repository) {
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive || !space.HasWriteAccess(u) {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	if err := h.store.DeleteRepository(r.Context(), space.ID, repo.Name); err != nil {
		if errors.Is(err, maven.ErrRepositoryNotFound) {
			problems.NotFound("Maven Repository", repo.Name).WriteToHTTP(w)
			return
		}
		slog.Error("Failed to delete maven repository", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
		spacesStore.ProblemFeatureNotEnabled("releases").WriteToHTTP(w)
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
		spacesStore.ProblemFeatureNotEnabled("releases").WriteToHTTP(w)
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

func (h *Handler) handleArtifactVersion(w http.ResponseWriter, r *http.Request) {
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
		spacesStore.ProblemFeatureNotEnabled("releases").WriteToHTTP(w)
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

	groupArtifactID := r.PathValue("group_artifact_id")
	if groupArtifactID == "" {
		problems.ValidationError("group_artifact_id", "Maven Artifact is required").WriteToHTTP(w)
		return
	}
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

	versionName := r.PathValue("version")
	if versionName == "" {
		problems.ValidationError("version", "Maven Artifact Version is required").WriteToHTTP(w)
		return
	}
	version := artifact.GetVersion(versionName)
	if version == nil {
		problems.NotFound("Maven Artifact Version", versionName).WriteToHTTP(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetArtifactVersion(w, r, space, repo, artifact, *version)
	case http.MethodDelete:
		h.handleDeleteArtifactVersion(w, r, space, repo, artifact, *version)
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodDelete}).WriteToHTTP(w)
	}
}

// no auth required
func (h *Handler) handleGetArtifactVersion(w http.ResponseWriter, r *http.Request, space *spaces.Space, repo *maven.Repository, artifact *maven.Artifact, version maven.ArtifactVersion) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60")
	json.NewEncoder(w).Encode(version)
}

func (h *Handler) handleDeleteArtifactVersion(w http.ResponseWriter, r *http.Request, space *spaces.Space, repo *maven.Repository, artifact *maven.Artifact, version maven.ArtifactVersion) {
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive || !space.HasWriteAccess(u) {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	// If repository is a mirror, forbid
	repoCheck, err := h.store.GetRepository(r.Context(), space.ID, repo.Name)
	if err != nil {
		slog.Error("Failed to get maven repository", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}
	if repoCheck.InternalMirror != nil {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	// Remove the version from artifact and update
	newVersions := []*maven.ArtifactVersion{}
	found := false
	for _, v := range artifact.Versions {
		if v.Version == version.Version {
			found = true
			continue
		}
		newVersions = append(newVersions, v)
	}
	if !found {
		problems.NotFound("Maven Artifact Version", version.Version).WriteToHTTP(w)
		return
	}

	artifact.Versions = newVersions
	if err := h.store.UpdateArtifact(r.Context(), space.ID, repo.Name, *artifact); err != nil {
		slog.Error("Failed to delete maven artifact version", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// no auth required
func (h *Handler) handleGetArtifact(w http.ResponseWriter, r *http.Request, space *spaces.Space, repo *maven.Repository, artifact *maven.Artifact) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60") // 1 minute
	json.NewEncoder(w).Encode(artifact)
}

func (h *Handler) handleDeleteArtifact(w http.ResponseWriter, r *http.Request, space *spaces.Space, repo *maven.Repository, artifact *maven.Artifact) {
	// Only allow write access
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive || !space.HasWriteAccess(u) {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	if err := h.store.DeleteArtifact(r.Context(), space.ID, repo.Name, artifact.Group, artifact.ID); err != nil {
		if errors.Is(err, maven.ErrArtifactNotFound) {
			problems.NotFound("Maven Artifact", artifact.Group+":"+artifact.ID).WriteToHTTP(w)
			return
		}
		slog.Error("Failed to delete maven artifact", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleJavadoc(w http.ResponseWriter, r *http.Request) {
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
		spacesStore.ProblemFeatureNotEnabled("releases").WriteToHTTP(w)
		return
	}

	repoName := r.PathValue("repository_name")
	if repoName == "" {
		problems.ValidationError("repository_name", "Maven Repository is required").WriteToHTTP(w)
		return
	}
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
	if groupArtifactID == "" {
		problems.ValidationError("group_artifact_id", "Maven Artifact is required").WriteToHTTP(w)
		return
	}
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

	versionName := r.PathValue("version")
	if versionName == "" {
		problems.ValidationError("version", "Maven Artifact Version is required").WriteToHTTP(w)
		return
	}
	version := artifact.GetVersion(versionName)
	if version == nil {
		problems.NotFound("Maven Artifact Version", versionName).WriteToHTTP(w)
		return
	}

	filePath := r.PathValue("file_path")
	if filePath == "" {
		problems.ValidationError("file_path", "Javadoc file path is required").WriteToHTTP(w)
		return
	}

	data, err := h.store.GetJavadocFile(r.Context(), space, repo, artifact, version.Version, filePath)
	if err != nil {
		if errors.Is(err, javadoccache.ErrJavadocNotFound) {
			problems.NotFound("Javadoc File", filePath).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get javadoc file", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	for ext, ct := range contentTypes {
		if strings.HasSuffix(filePath, ext) {
			w.Header().Set("Content-Type", ct)
			break
		}
	}

	w.Header().Set("Cache-Control", "public, max-age=86400") // 1 day
	w.Write(data)
}
