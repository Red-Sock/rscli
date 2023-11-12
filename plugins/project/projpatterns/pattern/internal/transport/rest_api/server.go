package rest_api

import (
	"context"
	"encoding/json"
	"net/http"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"financial-microservice/internal/config"
)

type Server struct {
	HttpServer *http.Server

	version string
}

func NewServer(cfg config.Config) *Server {
	r := mux.NewRouter()

	s := &Server{
		HttpServer: &http.Server{
			Addr:    "0.0.0.0:" + cfg.Api().Get(config.ApiRestAPI).GetPortStr(),
			Handler: r,
		},

		version: cfg.AppInfo().Version,
	}

	r.HandleFunc("/version", s.Version)

	return s
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		err := s.HttpServer.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			logrus.Fatal(err)
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.HttpServer.Shutdown(ctx)
}

func (s *Server) formResponse(r interface{}) ([]byte, error) {
	return json.Marshal(r)
}
