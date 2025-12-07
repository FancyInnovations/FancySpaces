package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/OliverSchlueter/goutils/problems"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/src/internal/spaces"
	"github.com/fancyinnovations/fancyspaces/src/internal/versions"
)

func (h *Handler) handleVersionFile(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("space_id")
	if sid == "" {
		problems.ValidationError("space_id", "Space ID is required").WriteToHTTP(w)
		return
	}

	vid := r.PathValue("version_id")
	if vid == "" {
		problems.ValidationError("version_id", "Version ID or name is required").WriteToHTTP(w)
		return
	}

	fileName := r.PathValue("file_name")
	if fileName == "" {
		problems.ValidationError("file_name", "File name is required").WriteToHTTP(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleDownloadVersionFile(w, r, sid, vid, fileName)
	case http.MethodPost:
		h.handleUploadVersionFile(w, r, sid, vid, fileName)
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodPost}).WriteToHTTP(w)
	}
}

func (h *Handler) handleUploadVersionFile(w http.ResponseWriter, r *http.Request, spaceID, versionID, fileName string) {
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}

	space, err := h.spaces.Get(spaceID)
	if err != nil {
		if errors.Is(err, spaces.ErrSpaceNotFound) {
			problems.NotFound("Space", spaceID).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get space by id", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	if !space.HasWriteAccess(u) {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	ver, err := h.store.Get(r.Context(), spaceID, versionID)
	if err != nil {
		if errors.Is(err, versions.ErrVersionNotFound) {
			problems.NotFound("Version", "latest").WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get version", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	data := make([]byte, r.ContentLength)
	_, err = r.Body.Read(data)
	if err != nil {
		slog.Error("Failed to read file data", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	if err := h.store.UploadVersionFile(r.Context(), ver, fileName, data); err != nil {
		slog.Error("Failed to upload version file", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// no auth required
func (h *Handler) handleDownloadVersionFile(w http.ResponseWriter, r *http.Request, spaceID, versionID, fileName string) {
	if versionID == "latest" {
		channel := r.URL.Query().Get("channel")
		platform := r.URL.Query().Get("platform")

		ver, err := h.store.GetLatest(r.Context(), spaceID, channel, platform)
		if err != nil {
			if errors.Is(err, versions.ErrVersionNotFound) {
				problems.NotFound("Version", "latest").WriteToHTTP(w)
				return
			}

			slog.Error("Failed to get latest version", sloki.WrapError(err))
			problems.InternalServerError("").WriteToHTTP(w)
			return
		}
		versionID = ver.ID
	}

	data, err := h.store.DownloadVersionFile(r.Context(), r, spaceID, versionID, fileName)
	if err != nil {
		if errors.Is(err, versions.ErrVersionNotFound) {
			problems.NotFound("Version", "latest").WriteToHTTP(w)
			return
		}

		slog.Error("Failed to download version file", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	w.Header().Set("Cache-Control", "public, 86400") // 24h
	w.Write(data)
}
