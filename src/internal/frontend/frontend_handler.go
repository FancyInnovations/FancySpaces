package frontend

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/OliverSchlueter/goutils/sloki"
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

type FS interface {
	ReadFile(name string) ([]byte, error)
}

type Handler struct {
	files FS
}

type Configuration struct {
	Files FS
}

func NewHandler(cfg Configuration) *Handler {
	return &Handler{
		files: cfg.Files,
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	pages := []string{
		"explore",
		"explore/minecraft-plugins",
		"explore/hytale-plugins",
		"explore/other-projects",
		"explore/by-other-creators",
		"tools/markdown-editor",
		"auth/login",
		"spaces/{space_id}",
		"spaces/{space_id}/versions",
		"spaces/{space_id}/versions/{version_id}",
		"spaces/{space_id}/issues",
		"spaces/{space_id}/issues/new",
		"spaces/{space_id}/issues/{issue_id}",
		"spaces/{space_id}/issues/{issue_id}/edit",
	}

	mux.HandleFunc("/{$}", h.handleIndex)
	for _, p := range pages {
		mux.HandleFunc(fmt.Sprintf("/%s", p), h.handleIndex)
		mux.HandleFunc(fmt.Sprintf("/%s/", p), h.handleIndex)
	}

	mux.HandleFunc("/", h.handleAssets)
}

func (h *Handler) handleIndex(w http.ResponseWriter, r *http.Request) {
	file, err := h.files.ReadFile("assets/index.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		slog.Error("Could not read index.html", sloki.WrapError(err), sloki.WrapRequest(r))
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write(file)
}

func (h *Handler) handleAssets(w http.ResponseWriter, r *http.Request) {
	path := "assets" + r.URL.Path

	file, err := h.files.ReadFile(path)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	for ext, ct := range contentTypes {
		if strings.HasSuffix(path, ext) {
			w.Header().Set("Content-Type", ct)
			break
		}
	}

	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write(file)
}
