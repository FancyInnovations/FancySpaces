package handler

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/OliverSchlueter/goutils/problems"
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
	sid := r.PathValue("space_id")
	// TODO check if space exists

	repoName := r.PathValue("repository_name")
	// TODO check if repository exists

	// TODO check if repo is private and if user has access

	group, err := GroupFromURL(r.URL.String())
	if err != nil {
		slog.Error("Failed to parse group ID from URL", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	artifactID, err := ArtifactFromURL(r.URL.String())
	if err != nil {
		slog.Error("Failed to parse artifact ID from URL", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// TODO get artifact / create if not exists

	version, err := VersionFromURL(r.URL.String())
	if err != nil {
		slog.Error("Failed to parse version from URL", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// TODO get version / create if not exists

	fileName, err := FilenameFromURL(r.URL.String())
	if err != nil {
		slog.Error("Failed to parse filename from URL", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// TODO add file to version / replace if exists

	slog.Info("----------------------")
	slog.Info(
		"Received Maven request",
		slog.String("method", r.Method),
		slog.String("url", r.URL.String()),
		slog.String("space_id", sid),
		slog.String("repository_name", repoName),
		slog.String("group", group),
		slog.String("artifact", artifactID),
		slog.String("version", version),
		slog.String("filename", fileName),
	)

	// write body to file system if method is PUT
	if r.Method == http.MethodPut {
		u := h.userFromCtx(r.Context())
		if u == nil || !u.Verified || !u.IsActive {
			problems.Unauthorized().WriteToHTTP(w)
			return
		}

		body, _ := io.ReadAll(r.Body)

		filePath := "data/maven/" + strings.TrimPrefix(r.URL.String(), "/maven/"+sid+"/"+repoName+"/")

		// make sure directory exists
		dir := strings.TrimSuffix(filePath, "/"+strings.Split(filePath, "/")[len(strings.Split(filePath, "/"))-1])
		if err := os.MkdirAll(dir, 0755); err != nil {
			slog.Error("Failed to create directory for Maven artifact", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// write file
		if err := os.WriteFile(filePath, body, 0644); err != nil {
			slog.Error("Failed to store Maven artifact", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// TODO store file in S3 instead

		slog.Info("Maven artifact stored successfully", slog.String("path", filePath))
	}

	slog.Info("----------------------")

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleStoreFile(w http.ResponseWriter, r *http.Request) {

}
