package handlers

import (
	"encoding/json"
	"net/http"
	"pr-reviewer/internal/models"
	"pr-reviewer/internal/service"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type PRHandler struct {
	pr  service.PRService
	log *zap.Logger
}

func NewPRHandler(pr service.PRService, log *zap.Logger) *PRHandler {
	return &PRHandler{pr: pr, log: log}
}

func (h *PRHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name     string `json:"name"`
		AuthorID string `json:"author_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	pr, reviewers, err := h.pr.CreatePR(r.Context(), in.Name, in.AuthorID)
	if err != nil {
		h.log.Error("CreatePR", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := struct {
		PR        *models.PullRequest `json:"pr"`
		Reviewers []models.User       `json:"reviewers"`
	}{pr, reviewers}

	json.NewEncoder(w).Encode(resp)
}

func (h *PRHandler) GetPR(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	pr, revs, err := h.pr.GetPR(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	resp := struct {
		PR        *models.PullRequest `json:"pr"`
		Reviewers []models.User       `json:"reviewers"`
	}{pr, revs}

	json.NewEncoder(w).Encode(resp)
}

func (h *PRHandler) MergePR(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	pr, err := h.pr.MergePR(r.Context(), id)
	if err != nil {
		h.log.Error("MergePR", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(pr)
}

func (h *PRHandler) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var in struct {
		OldReviewer string `json:"old_reviewer"`
	}
	json.NewDecoder(r.Body).Decode(&in)

	newR, err := h.pr.ReassignReviewer(r.Context(), id, in.OldReviewer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	json.NewEncoder(w).Encode(newR)
}
