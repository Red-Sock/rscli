package transport

import (
	"net"
	"net/http"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/rs/cors"
)

type httpServer struct {
	server *http.Server

	listener net.Listener
	serveMux *http.ServeMux
}

func newHttpServer(listener net.Listener, httpMux *http.ServeMux) httpServer {
	return httpServer{
		server: &http.Server{
			Handler: setUpCors().Handler(httpMux),
		},

		listener: listener,
		serveMux: httpMux,
	}
}

func (s *httpServer) AddHttpHandler(path string, handler http.Handler) {
	s.serveMux.Handle(path, handler)
}

func (s *httpServer) start() error {
	err := s.server.Serve(s.listener)
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return errors.Wrap(err, "error listening http server")
		}
	}

	return nil
}

func (s *httpServer) stop() error {
	return nil
}

func setUpCors() *cors.Cors {
	return cors.New(
		cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{
				http.MethodPost,
				http.MethodGet,
			},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: false,
		})
}
