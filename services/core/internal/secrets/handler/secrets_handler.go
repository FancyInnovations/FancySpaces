package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/OliverSchlueter/goutils/problems"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/core/internal/secrets"
	spacesStore "github.com/fancyinnovations/fancyspaces/core/internal/spaces"
	"github.com/fancyinnovations/fancyspaces/integrations/idp-go-sdk/idp"
	"github.com/fancyinnovations/fancyspaces/integrations/spaces-go-sdk/spaces"
)

type Handler struct {
	store       *secrets.Store
	spaces      *spacesStore.Store
	userFromCtx func(ctx context.Context) *idp.User
}

type Configuration struct {
	Store       *secrets.Store
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
	mux.HandleFunc(prefix+"/spaces/{space_id}/secrets", h.handleSecrets)
	mux.HandleFunc(prefix+"/spaces/{space_id}/secrets/{key}", h.handleSecret)
	mux.HandleFunc(prefix+"/spaces/{space_id}/secrets/{key}/decrypted", h.handleGetDecryptedSecret)
}

func (h *Handler) handleSecrets(w http.ResponseWriter, r *http.Request) {
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

	if !space.SecretsSettings.Enabled {
		spacesStore.ProblemFeatureNotEnabled("secrets").WriteToHTTP(w)
		return
	}

	u := h.userFromCtx(r.Context())
	if !idp.IsUserValid(u) {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}

	if !space.IsOwner(u) {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetSecrets(w, r, space)
	case http.MethodPost:
		h.handleCreateSecret(w, r, space)
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodPost}).WriteToHTTP(w)
	}
}

func (h *Handler) handleGetSecrets(w http.ResponseWriter, _ *http.Request, space *spaces.Space) {
	all, err := h.store.GetSecrets(space.ID)
	if err != nil {
		slog.Error("Failed to get secrets for space", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(all)
}

func (h *Handler) handleCreateSecret(w http.ResponseWriter, r *http.Request, space *spaces.Space) {
	var req CreateSecretReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}

	if req.Key == "" {
		problems.ValidationError("key", "Secret key is required").WriteToHTTP(w)
		return
	}

	if req.Value == "" {
		problems.ValidationError("value", "Secret value is required").WriteToHTTP(w)
		return
	}

	if err := h.store.CreateSecret(space.ID, req.Key, req.Value, req.Description); err != nil {
		if errors.Is(err, secrets.ErrSecretAlreadyExists) {
			problems.AlreadyExists("Secret", req.Key).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to create secret", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) handleSecret(w http.ResponseWriter, r *http.Request) {
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

	if !space.SecretsSettings.Enabled {
		spacesStore.ProblemFeatureNotEnabled("secrets").WriteToHTTP(w)
		return
	}

	u := h.userFromCtx(r.Context())
	if !idp.IsUserValid(u) {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}

	if !space.IsOwner(u) {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	key := r.PathValue("key")
	if key == "" {
		problems.ValidationError("key", "Secret key is required").WriteToHTTP(w)
		return
	}

	secret, err := h.store.GetSecret(space.ID, key)
	if err != nil {
		if errors.Is(err, secrets.ErrSecretNotFound) {
			problems.NotFound("Secret", key).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get secret", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetSecret(w, r, space, secret)
	case http.MethodPut:
		h.handleUpdateSecret(w, r, space, secret)
	case http.MethodDelete:
		h.handleDeleteSecret(w, r, space, secret)
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodPut, http.MethodDelete}).WriteToHTTP(w)
	}
}

func (h *Handler) handleGetSecret(w http.ResponseWriter, _ *http.Request, space *spaces.Space, secret *secrets.Secret) {
	secret, err := h.store.GetSecret(space.ID, secret.Key)
	if err != nil {
		if errors.Is(err, secrets.ErrSecretNotFound) {
			problems.NotFound("Secret", secret.Key).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get secret", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(secret)
}

func (h *Handler) handleGetDecryptedSecret(w http.ResponseWriter, r *http.Request) {
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

	if !space.SecretsSettings.Enabled {
		spacesStore.ProblemFeatureNotEnabled("secrets").WriteToHTTP(w)
		return
	}

	u := h.userFromCtx(r.Context())
	if !idp.IsUserValid(u) {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}

	if !space.IsOwner(u) {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	key := r.PathValue("key")
	if key == "" {
		problems.ValidationError("key", "Secret key is required").WriteToHTTP(w)
		return
	}

	_, err = h.store.GetSecret(space.ID, key)
	if err != nil {
		if errors.Is(err, secrets.ErrSecretNotFound) {
			problems.NotFound("Secret", key).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get secret", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	decryptedValue, err := h.store.GetDecryptedSecret(space.ID, key)
	if err != nil {
		slog.Error("Failed to get decrypted secret", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(decryptedValue))
}

func (h *Handler) handleUpdateSecret(w http.ResponseWriter, r *http.Request, _ *spaces.Space, secret *secrets.Secret) {
	var req UpdateSecretReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}

	if req.Description != "" {
		secret.Description = req.Description
	}
	if req.Value != "" {
		secret.Value = []byte(req.Value)
	}

	if err := h.store.UpdateSecret(secret); err != nil {
		slog.Error("Failed to update secret", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleDeleteSecret(w http.ResponseWriter, _ *http.Request, space *spaces.Space, secret *secrets.Secret) {
	if err := h.store.DeleteSecret(space.ID, secret.Key); err != nil {
		slog.Error("Failed to delete secret", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
