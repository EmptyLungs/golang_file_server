package api

import (
	"net/http"
)

func (s *Server) echoHandler(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Message string
	}{Message: "Hello!"}
	s.JsonResponse(w, r, http.StatusOK, resp)
}
