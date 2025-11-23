package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"pr-reviewer/internal/service"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// TeamsHandler отвечает за команды (teams).
type TeamsHandler struct {
	teams service.TeamService
	log   *zap.Logger
}

func NewTeamsHandler(ts service.TeamService, log *zap.Logger) *TeamsHandler {
	return &TeamsHandler{teams: ts, log: log}
}

// CreateTeam POST /teams
func (h *TeamsHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var in struct {
		TeamName    string  `json:"team_name"`
		Description *string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		h.log.Error("CreateTeam: decode", zap.Error(err))
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if in.TeamName == "" {
		http.Error(w, "team_name required", http.StatusBadRequest)
		return
	}
	t, err := h.teams.CreateTeam(r.Context(), in.TeamName, in.Description)
	if err != nil {
		h.log.Error("CreateTeam: service error", zap.Error(err))
		// If duplicate key or validation error — could return 409. For now return 500.
		http.Error(w, "could not create team", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(t)
}

// GetTeam GET /teams/{team_name}
func (h *TeamsHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "team_name")
	if name == "" {
		http.Error(w, "team_name required", http.StatusBadRequest)
		return
	}
	t, err := h.teams.GetTeam(r.Context(), name)
	if err != nil {
		h.log.Info("GetTeam: not found", zap.String("team", name), zap.Error(err))
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(t)
}

// ListTeams GET /teams
func (h *TeamsHandler) ListTeams(w http.ResponseWriter, r *http.Request) {
	list, err := h.teams.ListTeams(r.Context())
	if err != nil {
		h.log.Error("ListTeams: service error", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}

// DeleteTeam DELETE /teams/{team_name}
func (h *TeamsHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "team_name")
	if name == "" {
		http.Error(w, "team_name required", http.StatusBadRequest)
		return
	}
	err := h.teams.DeleteTeam(r.Context(), name)
	if err != nil {
		// service returns ErrTeamHasMembers when FK violation occurs
		if errors.Is(err, service.ErrTeamHasMembers) {
			h.log.Info("DeleteTeam: has members", zap.String("team", name))
			http.Error(w, "team has members, cannot delete", http.StatusConflict)
			return
		}
		h.log.Error("DeleteTeam: service error", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
