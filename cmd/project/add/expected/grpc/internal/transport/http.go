package transport

import (
	"bytes"
	"net"
	"net/http"
	"text/template"

	"github.com/rs/cors"
	"go.redsock.ru/rerrors"
)

type httpServer struct {
	server *http.Server

	listener net.Listener
	serveMux *http.ServeMux

	registeredPaths map[string]struct{}
}

func newHttpServer(listener net.Listener, httpMux *http.ServeMux) httpServer {
	return httpServer{
		server: &http.Server{
			Handler: setUpCors().Handler(httpMux),
		},
		registeredPaths: make(map[string]struct{}),
		listener:        listener,
		serveMux:        httpMux,
	}
}

func (s *httpServer) AddHttpHandler(path string, handler http.Handler) {
	s.registeredPaths[path] = struct{}{}
	s.serveMux.Handle(path, handler)
}

func (s *httpServer) start() error {
	homePageHandler := s.buildHomePageHandler()

	_, ok := s.registeredPaths["/"]
	if ok {
		s.AddHttpHandler("/about/", homePageHandler)
	} else {
		s.AddHttpHandler("/", homePageHandler)
	}

	err := s.server.Serve(s.listener)
	if err != nil {
		if !rerrors.Is(err, http.ErrServerClosed) {
			return rerrors.Wrap(err, "error listening http server")
		}
	}

	return nil
}

func (s *httpServer) stop() error {
	return nil
}

func (s *httpServer) buildHomePageHandler() http.Handler {
	var err error
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Error during server creation " + err.Error()))
	})

	aboutHtml := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>HomePage</title>
</head>
<body>
<ul>
    {{ range .Routes }} <li>
        <a href="{{ . }}"> {{ . }}</a>
    </li>{{ end}}
</ul>
</body>
</html>`
	tmpl, err := template.New("about").Parse(aboutHtml)
	if err != nil {
		return dummyHandler
	}

	type AboutPage struct {
		Routes []string
	}
	ap := AboutPage{
		Routes: make([]string, 0, len(s.registeredPaths)),
	}
	for p := range s.registeredPaths {
		ap.Routes = append(ap.Routes, p)
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, ap)
	if err != nil {
		return dummyHandler
	}

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write(buf.Bytes())
	})
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
