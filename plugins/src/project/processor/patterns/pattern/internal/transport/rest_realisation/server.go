package rest_realisation

import (
	"context"
	"encoding/json"
	"financial-microservice/internal/config"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	HttpServer *http.Server

	version string
}

func NewServer(cfg *config.Config) *Server {
	r := mux.NewRouter()
	s := &Server{
		HttpServer: &http.Server{
			Addr:    "0.0.0.0:" + cfg.GetString(config.ServerRestAPIPort),
			Handler: r,
		},

		version: cfg.GetString(config.AppInfoVersion),
	}

	r.HandleFunc("/version", s.Version)

	return s
}

func (s *Server) Start(ctx context.Context) error {
	return s.HttpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.HttpServer.Shutdown(ctx)
}

func (s *Server) formResponse(r interface{}) ([]byte, error) {
	return json.Marshal(r)
}
