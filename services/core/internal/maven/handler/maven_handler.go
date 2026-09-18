package handler

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/OliverSchlueter/goutils/problems"
	"github.com/OliverSchlueter/goutils/ratelimit"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/core/internal/analytics"
	"github.com/fancyinnovations/fancyspaces/core/internal/maven"
	spacesStore "github.com/fancyinnovations/fancyspaces/core/internal/spaces"
	"github.com/fancyinnovations/fancyspaces/integrations/idp-go-sdk/idp"
	"github.com/fancyinnovations/fancyspaces/integrations/spaces-go-sdk/spaces"
)

type Handler struct {
	store             *maven.Store
	spaces            *spacesStore.Store
	analytics         *analytics.Store
	userFromCtx       func(ctx context.Context) *idp.User
	downloadRatelimit *ratelimit.Service
}

type Configuration struct {
	Store       *maven.Store
	Spaces      *spacesStore.Store
	Analytics   *analytics.Store
	UserFromCtx func(ctx context.Context) *idp.User
}

func New(cfg Configuration) *Handler {
	downloadRatelimit := ratelimit.NewService(ratelimit.Configuration{
		TokensPerSecond: 1,
		MaxTokens:       5,
	})

	return &Handler{
		store:             cfg.Store,
		spaces:            cfg.Spaces,
		analytics:         cfg.Analytics,
		userFromCtx:       cfg.UserFromCtx,
		downloadRatelimit: downloadRatelimit,
	}
}

func (h *Handler) RegisterAPIEndpoints(prefix string, mux *http.ServeMux) {
	mux.HandleFunc(prefix+"/spaces/{space_id}/maven-repositories", h.handleRepositories)
	mux.HandleFunc(prefix+"/spaces/{space_id}/maven-repositories/{repository_name}", h.handleRepository)
	mux.HandleFunc(prefix+"/spaces/{space_id}/maven-repositories/{repository_name}/artifacts", h.handleArtifacts)
	mux.HandleFunc(prefix+"/spaces/{space_id}/maven-repositories/{repository_name}/artifacts/{group_artifact_id}", h.handleArtifact)
	mux.HandleFunc(prefix+"/spaces/{space_id}/maven-repositories/{repository_name}/artifacts/{group_artifact_id}/versions/{version}", h.handleArtifactVersion)

	mux.HandleFunc("/javadoc/{space_id}/{repository_name}/{group_artifact_id}/{version}/{file_path...}", h.handleJavadoc)
}

// RegisterMavenEndpoints registers the endpoints for Maven client requests
// These endpoints will be registered on the maven.fancyspaces.net domain
func (h *Handler) RegisterMavenEndpoints(mux *http.ServeMux) {
	// https://maven.fancyspaces.net/{space_id}/{repository_name}/{group_path}/{artifact_id}/{version}/{filename}
	mux.HandleFunc("/{space_id}/{repository_name}/", h.handleMavenRequest)
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
		spacesStore.ProblemFeatureNotEnabled("maven-repository").WriteToHTTP(w)
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
	case http.MethodGet, http.MethodHead:
		h.handleFetchFile(w, r, space, repo)
		return
	case http.MethodPut:
		h.handleStoreFile(w, r, space, repo)
		return
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodHead, http.MethodPut}).WriteToHTTP(w)
		return
	}
}

func (h *Handler) handleStoreFile(w http.ResponseWriter, r *http.Request, space *spaces.Space, repo *maven.Repository) {
	if repo.InternalMirror != nil {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

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

	fileName, err := FilenameFromURL(r.URL.String())
	if err != nil {
		problems.NotFound("Maven Artifact File", "<url>").WriteToHTTP(w)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed to read request body", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	artifactVersion := artifact.GetVersion(version)
	if artifactVersion == nil {
		artifactVersion = &maven.ArtifactVersion{
			Version:     version,
			PublishedAt: time.Now(),
			Files:       []*maven.ArtifactVersionFile{},
		}
		artifact.Versions = append(artifact.Versions, artifactVersion)
	} else {
		artifactVersion.PublishedAt = time.Now()
	}

	artifactVersionFile := artifactVersion.GetFile(fileName)
	if artifactVersionFile == nil {
		// create file
		artifactVersionFile = &maven.ArtifactVersionFile{
			Name: fileName,
			Size: int64(len(body)),
			URL:  "https://maven.fancyspaces.net/" + space.ID + "/" + repo.Name + "/" + groupPath + "/" + artifactID + "/" + version + "/" + fileName,
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

	if err := h.downloadRatelimit.CheckRequest(r, group+":"+artifactID); err != nil {
		ratelimit.RateLimitExceededProblem().WriteToHTTP(w)
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

	if IsMetadataURL(r.URL.String()) {
		var metadata maven.MetadataXML
		if version, err := VersionFromURL(r.URL.String()); err == nil && version != "" {
			metadata = artifact.ToVersionMetadataXML(version)
		} else {
			metadata = artifact.ToMetadataXML()
		}

		xmlData, err := xml.MarshalIndent(metadata, "", "  ")
		if err != nil {
			slog.Error("Failed to encode metadata XML", sloki.WrapError(err))
			problems.InternalServerError("").WriteToHTTP(w)
			return
		}
		xmlPayload := append([]byte(xml.Header), xmlData...)
		xmlPayload = append(xmlPayload, '\n')

		reqFilename, _ := FilenameFromURL(r.URL.String())
		switch {
		case reqFilename == "maven-metadata.xml":
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusOK)
			if r.Method == http.MethodHead {
				return
			}
			w.Write(xmlPayload)
		case reqFilename == "maven-metadata.xml.sha1" || strings.HasSuffix(reqFilename, ".xml.sha1"):
			hash := fmt.Sprintf("%x", sha1.Sum(xmlPayload))
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			if r.Method == http.MethodHead {
				return
			}
			w.Write([]byte(hash))
		case reqFilename == "maven-metadata.xml.md5" || strings.HasSuffix(reqFilename, ".xml.md5"):
			hash := fmt.Sprintf("%x", md5.Sum(xmlPayload))
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			if r.Method == http.MethodHead {
				return
			}
			w.Write([]byte(hash))
		case reqFilename == "maven-metadata.xml.sha256" || strings.HasSuffix(reqFilename, ".xml.sha256"):
			hash := fmt.Sprintf("%x", sha256.Sum256(xmlPayload))
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			if r.Method == http.MethodHead {
				return
			}
			w.Write([]byte(hash))
		case reqFilename == "maven-metadata.xml.sha512" || strings.HasSuffix(reqFilename, ".xml.sha512"):
			hash := fmt.Sprintf("%x", sha512.Sum512(xmlPayload))
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			if r.Method == http.MethodHead {
				return
			}
			w.Write([]byte(hash))
		default:
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusOK)
			if r.Method == http.MethodHead {
				return
			}
			w.Write(xmlPayload)
		}
		return
	}

	version, err := VersionFromURL(r.URL.String())
	if err != nil {
		slog.Error("Failed to parse version from URL", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
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

	if (artifactVersionFile.Name == artifactID+"-"+version+".jar" || (strings.HasSuffix(artifactVersionFile.Name, ".jar") && !strings.Contains(artifactVersionFile.Name, "-sources") && !strings.Contains(artifactVersionFile.Name, "-javadoc"))) && h.analytics != nil {
		if err := h.analytics.LogMavenArtifactDownload(r.Context(), space.ID, repo.Name, group, artifactID, version, r); err != nil {
			slog.Error("Failed to log maven artifact download", sloki.WrapError(err))
		}
	}

	data, err := h.store.DownloadArtifactFile(r.Context(), space.ID, repo.Name, group, artifactID, version, artifactVersionFile.Name)
	if err != nil {
		slog.Error("Failed to get artifact file", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)

	// HEAD request, no body
	if r.Method == http.MethodHead {
		return
	}

	w.Write(data)
}
