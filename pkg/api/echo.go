package api

import (
	"net/http"
)

func (s *Server) echoHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
