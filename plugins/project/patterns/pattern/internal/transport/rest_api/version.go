package rest_api

import (
	"net/http"
)

func (s *Server) Version(resp http.ResponseWriter, _ *http.Request) {
	vr := VersionResponse{
		Version: s.version,
	}
	b, err := s.formResponse(vr)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = resp.Write(b)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
	}
}

type VersionResponse struct {
	Version string
}
