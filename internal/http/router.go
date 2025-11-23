package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"pr-reviewer/internal/handlers"
)

func NewRouter(
	userHandler *handlers.UsersHandler,
	teamHandler *handlers.TeamsHandler,
	prHandler *handlers.PRHandler,
) http.Handler {

	r := chi.NewRouter()

	// Users
	r.Post("/users", userHandler.CreateUser)
	r.Get("/users", userHandler.ListUsers)
	r.Get("/users/{id}", userHandler.GetUser)
	r.Put("/users/{id}", userHandler.UpdateUser)
	r.Delete("/users/{id}", userHandler.DeleteUser)

	// Teams
	r.Post("/teams", teamHandler.CreateTeam)
	r.Get("/teams", teamHandler.ListTeams)
	r.Get("/teams/{name}", teamHandler.GetTeam)
	r.Delete("/teams/{name}", teamHandler.DeleteTeam)

	// Pull Requests
	r.Post("/pullRequest/create", prHandler.CreatePR)
	r.Get("/pullRequest/{id}", prHandler.GetPR)
	r.Post("/pullRequest/reassign", prHandler.ReassignReviewer)
	r.Post("/pullRequest/merge", prHandler.MergePR)

	return r
}
