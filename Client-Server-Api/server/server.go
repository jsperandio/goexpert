package server

import (
	"log/slog"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	Options *ServerOptions
	mux     *http.ServeMux
}

func NewServer(opt *ServerOptions) *Server {
	if opt == nil {
		opt = NewDefaultServerOptions()
	}

	return &Server{
		mux:     http.NewServeMux(),
		Options: opt,
	}
}

func (s *Server) RegisterRoute(path string, handler http.Handler) {
	slog.Info("registering route", "path", path)
	s.mux.Handle(path, handler)
}

func (s *Server) Start() error {
	slog.Info("starting server", "port", s.Options.Port)
	return http.ListenAndServe(":"+s.Options.Port, s.mux)
}
