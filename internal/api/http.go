package api

import (
	"log"
	"lucasbonna/pulse/db"
	"lucasbonna/pulse/internal/api/routes"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	db *db.Queries
}

func NewServer(database *db.Queries) *Server {
	return &Server{
		db: database,
	}
}

func (s *Server) startRoutes() http.Handler {
	r := chi.NewRouter()

	// Middlewares

	jobResource := routes.NewJobResource(s.db)

	r.Mount("/jobs", jobResource.Routes())

	return r
}

func (s *Server) Start(port string) error {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.CleanPath)
	r.Use(middleware.RedirectSlashes)

	r.Mount("/api", s.startRoutes())

	log.Println("starting http server on port ", port)
	if err := http.ListenAndServe(port, r); err != nil {
		return err
	}

	return nil
}
