package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) listFileHandler(w http.ResponseWriter, r *http.Request) {
	fileNames, err := s.fileManager.ListFiles()
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't list files\n%s", err.Error()), http.StatusBadRequest)
		return
	}
	response, err := json.Marshal(fileNames)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
