package api

import (
	"acmilanbot/cfg"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Server struct {
	config cfg.Configuration
	router *chi.Mux
}

func (s *Server) Start() {
	port := s.config.Port
	if port == 0 {
		port = 80
	}
	http.ListenAndServe(fmt.Sprintf(":%d", port), s.router)
}

func NewServer(c cfg.Configuration) *Server {
	s := &Server{
		config: c,
		router: chi.NewRouter(),
	}
	s.configureRouter()
	s.addRoutes()
	return s
}

func (s *Server) configureRouter() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))
	s.router.Use(middleware.Heartbeat("/ping"))
	s.router.Use(render.SetContentType(render.ContentTypeJSON))
}

func (s *Server) addRoutes() {
	s.router.Get("/fixtures", s.handleJobsScheduling)
	s.router.Route("/threads", func(r chi.Router) {
		s.router.Post("/pre-match", s.handlePreMatchThreadCreation)
		s.router.Post("/match", s.handleMatchThreadCreation)
		s.router.Post("/post-match", s.handlePostMatchThreadCreation)
	})
}
