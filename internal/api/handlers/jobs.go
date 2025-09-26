package handlers

import (
	"encoding/json"
	"net/http"

	"lucasbonna/pulse/db"
)

// JobsHandler holds the database dependency
type JobsHandler struct {
	db *db.Queries
}

// NewJobsHandler creates a new JobsHandler with database dependency
func NewJobsHandler(database *db.Queries) *JobsHandler {
	return &JobsHandler{
		db: database,
	}
}

// GetAllJobs handles GET /jobs
func (h *JobsHandler) GetAllJobs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	jobs, err := h.db.GetAllJobs(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch jobs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

// CreateJob handles POST /jobs
func (h *JobsHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement job creation
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{"message": "Not implemented yet"})
}
