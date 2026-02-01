package handler

import (
	"context"
	"encoding/xml"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/OliverSchlueter/goutils/problems"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/internal/analytics"
	"github.com/fancyinnovations/fancyspaces/internal/auth"
	"github.com/fancyinnovations/fancyspaces/internal/maven"
	"github.com/fancyinnovations/fancyspaces/internal/spaces"
)

type Handler struct {
	store       *maven.Store
	spaces      *spaces.Store
	analytics   *analytics.Store
	userFromCtx func(ctx context.Context) *auth.User
}

type Configuration struct {
	Store       *maven.Store
	Spaces      *spaces.Store
	Analytics   *analytics.Store
	UserFromCtx func(ctx context.Context) *auth.User
}

func New(cfg Configuration) *Handler {
	return &Handler{
		store:       cfg.Store,
		spaces:      cfg.Spaces,
		analytics:   cfg.Analytics,
		userFromCtx: cfg.UserFromCtx,
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	// https://fancyspaces.net/maven/{space_id}/{repository_name}/{group_id}/{artifact_id}/{version}/{filename}
	mux.HandleFunc("/maven/{space_id}/{repository_name}/", h.handleMavenRequest)
}

func (h *Handler) handleMavenRequest(w http.ResponseWriter, r *http.Request) {
	// get space
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
		spaces.ProblemFeatureNotEnabled("maven-repository").WriteToHTTP(w)
		return
	}

	// get repository
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
		h.handleFetchFile(w, r, space, repo)
		return
	case http.MethodPut:
		h.handleStoreFile(w, r, space, repo)
		return
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodPut}).WriteToHTTP(w)
		return
	}
}

func (h *Handler) handleStoreFile(w http.ResponseWriter, r *http.Request, space *spaces.Space, repo *maven.Repository) {
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive || !space.HasWriteAccess(u) {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	if IsMetadataURL(r.URL.String()) {
		w.WriteHeader(http.StatusOK)
		return
	}

	group, err := GroupFromURL(r.URL.String())
	if err != nil {
		problems.NotFound("Maven Artifact", "<url>").WriteToHTTP(w)
		return
	}
	groupPath := strings.ReplaceAll(group, ".", "/")

	artifactID, err := ArtifactFromURL(r.URL.String())
	if err != nil {
		problems.NotFound("Maven Artifact", "<url>").WriteToHTTP(w)
		return
	}

	artifact, err := h.store.GetArtifact(r.Context(), space.ID, repo.Name, group, artifactID)
	if err != nil {
		if errors.Is(err, maven.ErrArtifactNotFound) {
			// create artifact
			artifact = &maven.Artifact{
				SpaceID:    space.ID,
				Repository: repo.Name,
				Group:      group,
				ID:         artifactID,
				Versions:   nil,
			}
			if err := h.store.CreateArtifact(r.Context(), space.ID, repo.Name, *artifact); err != nil {
				slog.Error("Failed to create artifact", sloki.WrapError(err))
				problems.InternalServerError("").WriteToHTTP(w)
				return
			}
		} else {
			slog.Error("Failed to get artifact", sloki.WrapError(err))
			problems.InternalServerError("").WriteToHTTP(w)
			return
		}
	}

	version, err := VersionFromURL(r.URL.String())
	if err != nil {
		slog.Error("Failed to parse version from URL", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	artifactVersion := artifact.GetVersion(version)
	if artifactVersion == nil {
		// create version
		artifactVersion = &maven.ArtifactVersion{
			Version:     version,
			PublishedAt: time.Now(),
			Files:       []*maven.ArtifactVersionFile{},
		}
		artifact.Versions = append(artifact.Versions, artifactVersion)
	} else {
		// update published at
		artifactVersion.PublishedAt = time.Now()
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed to read request body", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	fileName, err := FilenameFromURL(r.URL.String())
	if err != nil {
		problems.NotFound("Maven Artifact File", "<url>").WriteToHTTP(w)
		return
	}

	artifactVersionFile := artifactVersion.GetFile(fileName)
	if artifactVersionFile == nil {
		// create file
		artifactVersionFile = &maven.ArtifactVersionFile{
			Name: fileName,
			Size: int64(len(body)),
			URL:  "https://fancyspaces.net/maven/" + space.ID + "/" + repo.Name + "/" + groupPath + "/" + artifactID + "/" + version + "/" + fileName,
		}
		artifactVersion.Files = append(artifactVersion.Files, artifactVersionFile)
	} else {
		// update file size
		artifactVersionFile.Size = int64(len(body))
	}

	if err := h.store.UpdateArtifact(r.Context(), space.ID, repo.Name, *artifact); err != nil {
		slog.Error("Failed to update artifact version file", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	if err := h.store.UploadArtifactFile(r.Context(), space.ID, repo.Name, group, artifactID, version, fileName, body); err != nil {
		slog.Error("Failed to upload artifact file", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) handleFetchFile(w http.ResponseWriter, r *http.Request, space *spaces.Space, repo *maven.Repository) {
	group, err := GroupFromURL(r.URL.String())
	if err != nil {
		problems.NotFound("Maven Artifact", "<url>").WriteToHTTP(w)
		return
	}

	artifactID, err := ArtifactFromURL(r.URL.String())
	if err != nil {
		problems.NotFound("Maven Artifact", "<url>").WriteToHTTP(w)
		return
	}
	artifact, err := h.store.GetArtifact(r.Context(), space.ID, repo.Name, group, artifactID)
	if err != nil {
		if errors.Is(err, maven.ErrArtifactNotFound) {
			problems.NotFound("Maven Artifact", artifactID).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get artifact", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	version, err := VersionFromURL(r.URL.String())
	if err != nil {
		slog.Error("Failed to parse version from URL", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if IsMetadataURL(r.URL.String()) {
		metadata := artifact.ToMetadataXML()

		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)

		encoder := xml.NewEncoder(w)
		defer encoder.Close()

		encoder.Indent("", "  ")
		if err := encoder.Encode(metadata); err != nil {
			slog.Error("Failed to encode metadata XML", sloki.WrapError(err))
			problems.InternalServerError("").WriteToHTTP(w)
			return
		}
		return
	}

	artifactVersion := artifact.GetVersion(version)
	if artifactVersion == nil {
		problems.NotFound("Maven Artifact Version", version).WriteToHTTP(w)
		return
	}

	fileName, err := FilenameFromURL(r.URL.String())
	if err != nil {
		problems.NotFound("Maven Artifact File", "<url>").WriteToHTTP(w)
		return
	}
	artifactVersionFile := artifactVersion.GetFile(fileName)
	if artifactVersionFile == nil {
		problems.NotFound("Maven Artifact File", fileName).WriteToHTTP(w)
		return
	}

	data, err := h.store.DownloadArtifactFile(r.Context(), space.ID, repo.Name, group, artifactID, version, fileName)
	if err != nil {
		slog.Error("Failed to get artifact file", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
