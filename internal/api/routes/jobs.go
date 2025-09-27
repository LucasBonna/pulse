package routes

import (
	"context"
	"lucasbonna/pulse/db"
	"lucasbonna/pulse/internal/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type JobsResource struct {
	db *db.Queries
}

func NewJobResource(database *db.Queries) *JobsResource {
	return &JobsResource{
		db: database,
	}
}

func (js JobsResource) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", js.GetAllJobs)
	r.Post("/", js.CreateJob)

	return r
}

func (js JobsResource) GetAllJobs(w http.ResponseWriter, r *http.Request) {
	allJobs, err := js.db.GetAllJobs(context.Background())
	if err != nil {
		utils.WriteJsonError(w, http.StatusInternalServerError, "failed to fetch all jobs")
		return
	}

	utils.WriteJsonResponse(w, http.StatusOK, allJobs)
}

func (js JobsResource) CreateJob(w http.ResponseWriter, r *http.Request) {
}
