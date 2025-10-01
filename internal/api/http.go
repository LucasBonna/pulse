package api

import (
	"log"
	"lucasbonna/pulse/db"
	internal_middleware "lucasbonna/pulse/internal/api/middleware"
	"lucasbonna/pulse/internal/api/routes"
	"lucasbonna/pulse/internal/config"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	db     *db.Queries
	config *config.Env
}

func NewServer(database *db.Queries, config *config.Env) *Server {
	return &Server{
		db:     database,
		config: config,
	}
}

func (s *Server) startRoutes() http.Handler {
	r := chi.NewRouter()

	jobResource := routes.NewJobResource(s.db)

	r.Mount("/jobs", jobResource.Routes())

	return r
}

func (s *Server) Start() error {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.CleanPath)
	r.Use(middleware.RedirectSlashes)

	r.Use(internal_middleware.AuthenticationMiddleware(s.config.Token))

	r.Mount("/api", s.startRoutes())

	log.Println("starting http server on port ", s.config.Port)
	if err := http.ListenAndServe(":"+s.config.Port, r); err != nil {
		return err
	}

	return nil
}
