package routes

import (
	"lucasbonna/pulse/internal/api/handlers"

	"github.com/go-chi/chi/v5"
)

// JobsRoutes sets up all job-related routes
func JobsRoutes(jobsHandler *handlers.JobsHandler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", jobsHandler.GetAllJobs)
	r.Post("/", jobsHandler.CreateJob)
	// Add more job routes here as needed
	// r.Get("/{id}", jobsHandler.GetJobByID)
	// r.Put("/{id}", jobsHandler.UpdateJob)
	// r.Delete("/{id}", jobsHandler.DeleteJob)

	return r
}
