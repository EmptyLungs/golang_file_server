package api

import (
	"fmt"
	"net/http"
)

func (s *Server) uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse form data\n%s", err.Error()), http.StatusBadRequest)
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = s.fileManager.CreateFile(file, handler)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to upload file\n%s", err.Error()), http.StatusBadRequest)
		return
	}
}
