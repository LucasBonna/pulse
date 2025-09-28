package routes

import (
	"context"
	"database/sql"
	"log"
	"lucasbonna/pulse/db"
	"lucasbonna/pulse/internal/api/dto"
	"lucasbonna/pulse/internal/api/middleware"
	"lucasbonna/pulse/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type JobsResource struct {
	db         *db.Queries
	validation *middleware.ValidationMiddleware
}

func NewJobResource(database *db.Queries) *JobsResource {
	return &JobsResource{
		db:         database,
		validation: middleware.NewValidationMiddleware(),
	}
}

func (js JobsResource) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", js.GetAllJobs)
	r.With(middleware.ValidateBody(js.validation, dto.CreateJobRequest{})).Post("/", js.CreateJob)
	r.Delete("/{id}", js.DeleteJob)
	return r
}

func (js JobsResource) GetAllJobs(w http.ResponseWriter, r *http.Request) {
	allJobs, err := js.db.GetAllJobs(context.Background())
	if err != nil {
		utils.WriteJsonError(w, http.StatusInternalServerError, "failed to fetch all jobs")
		return
	}

	var jobResponses []dto.CreateJobResponse
	for _, job := range allJobs {
		jobResponses = append(jobResponses, fromDBJob(job))
	}

	utils.WriteJsonResponse(w, http.StatusOK, jobResponses)
}

func (js JobsResource) DeleteJob(w http.ResponseWriter, r *http.Request) {
	jobIDStr := chi.URLParam(r, "id")
	jobID, err := strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		utils.WriteJsonError(w, http.StatusBadRequest, "invalid job ID")
		return
	}

	err = js.db.DeleteJob(context.Background(), jobID)
	if err != nil {
		utils.WriteJsonError(w, http.StatusInternalServerError, "failed to delete job")
		return
	}

	utils.WriteJsonResponse(w, http.StatusOK, "job deleted")
}

func (js JobsResource) CreateJob(w http.ResponseWriter, r *http.Request) {
	data := middleware.GetValidatedData[dto.CreateJobRequest](r)

	createdJob, err := js.db.CreateJob(context.Background(), db.CreateJobParams{
		Name:            data.Name,
		Url:             data.URL,
		Method:          data.Method,
		Headers:         sql.NullString{String: data.Headers, Valid: data.Method != ""},
		IntervalSeconds: data.IntervalSeconds,
		NextRunAt:       sql.NullTime{Time: time.Now(), Valid: true},
		Active:          sql.NullBool{Bool: data.Active, Valid: true},
	})
	if err != nil {
		log.Println("error creating job", err)
		utils.WriteJsonResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJsonResponse(w, http.StatusOK, fromDBJob(createdJob))
}

func fromDBJob(dbJob db.Job) dto.CreateJobResponse {
	response := dto.CreateJobResponse{
		Id:              dbJob.ID,
		Name:            dbJob.Name,
		Method:          dbJob.Method.(string),
		Url:             dbJob.Url,
		IntervalSeconds: dbJob.IntervalSeconds,
	}

	if dbJob.Headers.Valid && dbJob.Headers.String != "" {
		response.Headers = &dbJob.Headers.String
	}

	if dbJob.NextRunAt.Valid {
		response.NextRunAt = &dbJob.NextRunAt.Time
	}

	if dbJob.Active.Valid {
		response.Active = &dbJob.Active.Bool
	}

	return response
}
