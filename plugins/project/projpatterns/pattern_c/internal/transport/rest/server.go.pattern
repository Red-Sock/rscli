package rest

import (
	"context"
	"encoding/json"
	"net/http"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka/api"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"

	"proj_name/internal/config"
)

type Server struct {
	HttpServer *http.Server

	version string
}

func NewServer(cfg config.Config, server *api.Rest) *Server {
	r := mux.NewRouter()

	s := &Server{
		HttpServer: &http.Server{
			Addr:    "0.0.0.0:" + server.GetPortStr(),
			Handler: setUpCors().Handler(r),
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

func setUpCors() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodGet,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})
}
