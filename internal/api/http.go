package api

import (
	"net/http"

	"lucasbonna/pulse/db"
	"lucasbonna/pulse/internal/api/handlers"
	"lucasbonna/pulse/internal/api/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server holds all dependencies for the HTTP server
type Server struct {
	db *db.Queries
}

// NewServer creates a new HTTP server with dependencies
func NewServer(database *db.Queries) *Server {
	return &Server{
		db: database,
	}
}

// SetupRoutes configures all routes with their dependencies
func (s *Server) SetupRoutes() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.CleanPath)
	r.Use(middleware.RedirectSlashes)

	// Health check
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	// Initialize handlers with dependencies
	jobsHandler := handlers.NewJobsHandler(s.db)

	// Mount routes
	r.Mount("/jobs", routes.JobsRoutes(jobsHandler))

	return r
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	handler := s.SetupRoutes()
	return http.ListenAndServe(addr, handler)
}
