package api

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func (s *Server) uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		s.JsonError(w, r, http.StatusBadRequest, "Failed to parse form data")
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		s.JsonError(w, r, http.StatusBadRequest, "Failed to read file")
		return
	}
	defer file.Close()

	fileSize := handler.Size

	if fileSize > s.config.UploaderMaxFileSizeMB*1024*1024 {
		s.JsonError(w, r, http.StatusBadRequest, fmt.Sprintf("File is too big %d %d", handler.Size, s.config.UploaderMaxFileSizeMB*1024*1024))
		return
	} else if fileSize == 0 {
		s.JsonResponse(w, r, http.StatusBadRequest, "Empty file")
		return
	}

	err = s.fileManager.CreateFile(file, handler)
	if err != nil {
		s.JsonError(w, r, http.StatusInternalServerError, "Failed to upload file")
		s.logger.Error("Failed to upload file", zap.Error(err))
		return
	}
}
