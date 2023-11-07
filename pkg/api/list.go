package api

import (
	"net/http"
)

func (s *Server) listFileHandler(w http.ResponseWriter, r *http.Request) {
	fileNames, err := s.fileManager.List()
	if err != nil {
		s.JsonError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	s.JsonResponse(w, r, http.StatusOK, fileNames)
}
