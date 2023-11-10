package api

import (
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"

	"go.uber.org/zap"
)

type Payload struct {
	Filename string `json:"filename"`
}

func (s *Server) deleteFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		s.JsonError(w, r, http.StatusBadRequest, "Empty request body")
		return
	}
	var payload Payload
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		s.logger.Debug("error", zap.Error(err))
		s.JsonError(w, r, http.StatusBadRequest, "Failed to parse request body")
		return
	}
	if err := s.fileManager.Delete(payload.Filename); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			s.JsonResponse(w, r, http.StatusNotFound, nil)
			return
		}
		s.JsonError(w, r, http.StatusInternalServerError, "Failed to delete file")
		return
	}
	if payload.Filename == "" {
		s.JsonError(w, r, http.StatusBadRequest, "Missing filename in request body")
		return
	}
	s.logger.Info("asdsadadssad", zap.String("path", payload.Filename))
	if err := s.fileManager.Delete(payload.Filename); err != nil {
		s.JsonError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	s.JsonResponse(w, r, http.StatusNoContent, nil)
}
