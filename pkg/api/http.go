package api

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func (s *Server) JsonResponse(w http.ResponseWriter, r *http.Request, status int, result interface{}) {
	body, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("JSON Marshal failed", zap.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	w.Write(body)
}

func (s *Server) JsonError(w http.ResponseWriter, r *http.Request, status int, err string) {
	errorMsg := struct {
		Error string `json:"error"`
	}{Error: err}
	s.JsonResponse(w, r, status, errorMsg)
}
