package api

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/m4tthewde/fdmxyz/internal/config"
	"github.com/m4tthewde/fdmxyz/internal/db"
)

type Server struct {
	server       *http.Server
	router       *chi.Mux
	routeHandler *RouteHandler
}

func NewServer(config *config.Config) *Server {
	s := Server{
		server: &http.Server{Addr: ":" + config.Port},
		router: chi.NewRouter(),
		routeHandler: &RouteHandler{
			locked: false,
			config: config,
			mongoHandler: &db.MongoHandler{
				Config: config,
			},
		},
	}

	return &s
}

func (s *Server) Run() {
	s.router.Use(middleware.Logger)
	s.registerRoutes()
	s.server.Handler = s.router

	err := s.server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func (s *Server) registerRoutes() {
	s.router.Get("/webhook", s.routeHandler.get())
	s.router.Post("/webhook", s.routeHandler.register())
	s.router.Delete("/webhook", s.routeHandler.delete())
	s.router.Post("/twitch", s.routeHandler.twitch())
}
