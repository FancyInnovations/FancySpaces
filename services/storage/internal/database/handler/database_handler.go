package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/OliverSchlueter/goutils/problems"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/storage/internal/auth"
	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
)

type Handler struct {
	store       *database.Store
	userFromCtx func(ctx context.Context) *auth.User
}

type Configuration struct {
	Store       *database.Store
	UserFromCtx func(ctx context.Context) *auth.User
}

func New(cfg Configuration) *Handler {
	return &Handler{
		store:       cfg.Store,
		userFromCtx: cfg.UserFromCtx,
	}
}

func (h *Handler) Register(prefix string, mux *http.ServeMux) {
	mux.HandleFunc(prefix+"/databases", h.handleDatabases)
	mux.HandleFunc(prefix+"/databases/{db_id}", h.handleDatabase)
	mux.HandleFunc(prefix+"/databases/{db_id}/collections", h.handleCollections)
	mux.HandleFunc(prefix+"/databases/{db_id}/collections/{coll_name}", h.handleCollection)
}

func (h *Handler) handleDatabases(w http.ResponseWriter, r *http.Request) {
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetDatabases(w, r)
	case http.MethodPost:
		h.handleCreateDatabase(w, r)
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodPost}).WriteToHTTP(w)
	}
}

func (h *Handler) handleGetDatabases(w http.ResponseWriter, r *http.Request) {
	u := h.userFromCtx(r.Context())

	var dbs []*database.Database
	var err error
	if u.IsAdmin() {
		// Admins can see all databases
		dbs, err = h.store.GetAllDatabases(r.Context())
		if err != nil {
			slog.Error("Failed to get all databases", sloki.WrapError(err))
			problems.InternalServerError("").WriteToHTTP(w)
			return
		}
	} else {
		// Regular users can only see databases they are members of
		dbs, err = h.store.GetDatabasesForUser(r.Context(), u.ID)
		if err != nil {
			slog.Error("Failed to get databases for user", sloki.WrapError(err))
			problems.InternalServerError("").WriteToHTTP(w)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60") // 1 minute
	json.NewEncoder(w).Encode(dbs)
}

func (h *Handler) handleCreateDatabase(w http.ResponseWriter, r *http.Request) {
	u := h.userFromCtx(r.Context())

	var req CreateDatabaseReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode create database request", sloki.WrapError(err))
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}

	if req.Name == "" {
		problems.ValidationError("name", "Name is required").WriteToHTTP(w)
		return
	}

	if err := h.store.CreateDatabase(r.Context(), req.Name, u); err != nil {
		if errors.Is(err, database.ErrDatabaseAlreadyExists) {
			problems.AlreadyExists("Database", req.Name).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to create database", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) handleDatabase(w http.ResponseWriter, r *http.Request) {
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}

	dbid := r.PathValue("db_id")
	if dbid == "" {
		problems.ValidationError("db_id", "Database ID is required").WriteToHTTP(w)
		return
	}

	db, err := h.store.GetDatabase(r.Context(), dbid)
	if err != nil {
		if errors.Is(err, database.ErrDatabaseNotFound) {
			problems.NotFound("Database", dbid).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get database by id", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}
	if !db.IsMember(u.ID) && !u.IsAdmin() {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetDatabase(w, r, db)
	case http.MethodPut:
		h.handleUpdateDatabase(w, r, db)
	case http.MethodDelete:
		h.handleDeleteDatabase(w, r, db)
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodPut, http.MethodDelete}).WriteToHTTP(w)
	}
}

func (h *Handler) handleGetDatabase(w http.ResponseWriter, _ *http.Request, db *database.Database) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600") // 1 hour
	json.NewEncoder(w).Encode(db)
}

func (h *Handler) handleUpdateDatabase(w http.ResponseWriter, r *http.Request, db *database.Database) {
	u := h.userFromCtx(r.Context())
	if !db.HasPermission(u.ID, database.PermissionLevelAdmin) && !u.IsAdmin() {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	var req UpdateDatabaseReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode update database request", sloki.WrapError(err))
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}

	if err := h.store.UpdateDatabaseUsers(r.Context(), db, req.Users); err != nil {
		slog.Error("Failed to update database users", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleDeleteDatabase(w http.ResponseWriter, r *http.Request, db *database.Database) {
	u := h.userFromCtx(r.Context())
	if !db.HasPermission(u.ID, database.PermissionLevelAdmin) && !u.IsAdmin() {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	if err := h.store.DeleteDatabase(r.Context(), db.Name); err != nil {
		slog.Error("Failed to delete database", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleCollections(w http.ResponseWriter, r *http.Request) {
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}

	dbid := r.PathValue("db_id")
	if dbid == "" {
		problems.ValidationError("db_id", "Database ID is required").WriteToHTTP(w)
		return
	}

	db, err := h.store.GetDatabase(r.Context(), dbid)
	if err != nil {
		if errors.Is(err, database.ErrDatabaseNotFound) {
			problems.NotFound("Database", dbid).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get database by id", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}
	if !db.IsMember(u.ID) && !u.IsAdmin() {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetCollections(w, r, db)
	case http.MethodPost:
		h.handleCreateCollection(w, r, db)
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodPost}).WriteToHTTP(w)
	}
}

func (h *Handler) handleGetCollections(w http.ResponseWriter, r *http.Request, db *database.Database) {
	colls, err := h.store.GetAllCollections(r.Context(), db)
	if err != nil {
		slog.Error("Failed to get collections for database", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600") // 1 hour
	json.NewEncoder(w).Encode(colls)
}

func (h *Handler) handleCreateCollection(w http.ResponseWriter, r *http.Request, db *database.Database) {
	u := h.userFromCtx(r.Context())
	if !db.HasPermission(u.ID, database.PermissionLevelAdmin) && !u.IsAdmin() {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	var req CreateCollectionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode create collection request", sloki.WrapError(err))
		problems.ValidationError("body", "Invalid JSON").WriteToHTTP(w)
		return
	}

	if req.Name == "" {
		problems.ValidationError("name", "Name is required").WriteToHTTP(w)
		return
	}
	if req.Engine == "" {
		problems.ValidationError("engine", "Engine is required").WriteToHTTP(w)
		return
	}

	if err := h.store.CreateCollection(r.Context(), db, req.Name, req.Engine); err != nil {
		if errors.Is(err, database.ErrCollectionAlreadyExists) {
			problems.AlreadyExists("Collection", req.Name).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to create collection", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) handleCollection(w http.ResponseWriter, r *http.Request) {
	u := h.userFromCtx(r.Context())
	if u == nil || !u.Verified || !u.IsActive {
		problems.Unauthorized().WriteToHTTP(w)
		return
	}

	dbid := r.PathValue("db_id")
	if dbid == "" {
		problems.ValidationError("db_id", "Database ID is required").WriteToHTTP(w)
		return
	}

	db, err := h.store.GetDatabase(r.Context(), dbid)
	if err != nil {
		if errors.Is(err, database.ErrDatabaseNotFound) {
			problems.NotFound("Database", dbid).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get database by id", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}
	if !db.IsMember(u.ID) && !u.IsAdmin() {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	collname := r.PathValue("coll_name")
	if collname == "" {
		problems.ValidationError("coll_name", "Collection name is required").WriteToHTTP(w)
		return
	}

	coll, err := h.store.GetCollection(r.Context(), db, collname)
	if err != nil {
		if errors.Is(err, database.ErrCollectionNotFound) {
			problems.NotFound("Collection", collname).WriteToHTTP(w)
			return
		}

		slog.Error("Failed to get collection by name", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetCollection(w, r, db, coll)
	case http.MethodDelete:
		h.handleDeleteCollection(w, r, db, coll)
	default:
		problems.MethodNotAllowed(r.Method, []string{http.MethodGet, http.MethodDelete}).WriteToHTTP(w)
	}
}

func (h *Handler) handleGetCollection(w http.ResponseWriter, _ *http.Request, _ *database.Database, coll *database.Collection) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600") // 1 hour
	json.NewEncoder(w).Encode(coll)
}

func (h *Handler) handleDeleteCollection(w http.ResponseWriter, r *http.Request, db *database.Database, coll *database.Collection) {
	u := h.userFromCtx(r.Context())
	if !db.HasPermission(u.ID, database.PermissionLevelAdmin) && !u.IsAdmin() {
		problems.Forbidden().WriteToHTTP(w)
		return
	}

	if err := h.store.DeleteCollection(r.Context(), db, coll.Name); err != nil {
		slog.Error("Failed to delete collection", sloki.WrapError(err))
		problems.InternalServerError("").WriteToHTTP(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
