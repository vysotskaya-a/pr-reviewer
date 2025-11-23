package handlers

import (
	"encoding/json"
	"net/http"

	"pr-reviewer/internal/service"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type UsersHandler struct {
	users service.UserService
	log   *zap.Logger
}

func NewUsersHandler(us service.UserService, log *zap.Logger) *UsersHandler {
	return &UsersHandler{users: us, log: log}
}

// CreateUser POST /users
func (h *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Username    string  `json:"username"`
		DisplayName *string `json:"display_name"`
		TeamName    *string `json:"team_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		h.log.Error("CreateUser: decode", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if in.Username == "" {
		http.Error(w, "username required", http.StatusBadRequest)
		return
	}

	u, err := h.users.CreateUser(r.Context(), in.Username, in.DisplayName, in.TeamName)
	if err != nil {
		// сервис user возвращает строковую ошибку "team not found" при отсутствии команды
		if err.Error() == "team not found" {
			h.log.Info("CreateUser: team not found", zap.String("team", stringPtrToValue(in.TeamName)))
			http.Error(w, "team not found", http.StatusBadRequest)
			return
		}
		h.log.Error("CreateUser: service error", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(u)
}

// ListUsers GET /users
func (h *UsersHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	list, err := h.users.ListUsers(r.Context())
	if err != nil {
		h.log.Error("ListUsers: service error", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}

// GetUser GET /users/{id}
func (h *UsersHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	u, err := h.users.GetUser(r.Context(), id)
	if err != nil {
		// repository/service returns non-nil error if not found — отображаем 404
		h.log.Info("GetUser: not found", zap.String("id", id), zap.Error(err))
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(u)
}

// UpdateUser PUT /users/{id}
func (h *UsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var in struct {
		DisplayName *string `json:"display_name"`
		IsActive    *bool   `json:"is_active"`
		TeamName    *string `json:"team_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		h.log.Error("UpdateUser: decode", zap.Error(err))
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err := h.users.UpdateUser(r.Context(), id, in.DisplayName, in.IsActive, in.TeamName)
	if err != nil {
		// as with CreateUser, service returns "team not found" string when team absent
		if err.Error() == "team not found" {
			h.log.Info("UpdateUser: team not found", zap.String("team", stringPtrToValue(in.TeamName)))
			http.Error(w, "team not found", http.StatusBadRequest)
			return
		}
		// if update failed because user not found or other DB error -- respond 500 (could refine)
		h.log.Error("UpdateUser: service error", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteUser DELETE /users/{id}
func (h *UsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}
	if err := h.users.DeleteUser(r.Context(), id); err != nil {
		h.log.Error("DeleteUser: service error", zap.String("id", id), zap.Error(err))
		// Could return 404 if not found, but repository Delete returns nil even if row missing.
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// helper: safe deref for logging
func stringPtrToValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
