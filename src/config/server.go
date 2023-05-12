package config

import (
	"database/sql"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Server struct {
    Router *chi.Mux
	Databse *sql.DB
}


func CreateNewServer() *Server {
    s := &Server{}
	s.Databse = ConnectDB()
    s.Router = chi.NewRouter()
    return s
}

func (s *Server) MountMiddlewares() {
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
}

func (s *Server) MountHandlers(handlers func(chi.Router)) {
	s.Router.Route("/api/v1", handlers)
}
