package api

import (
	"time"

	"go.uber.org/zap"
)

func NewMockServer() *Server {
	config := &Config{
		Host:                  "",
		Port:                  "8080",
		HttpServerTimeout:     5 * time.Second,
		UploaderDir:           "./test-data",
		UploaderMaxFileSizeMB: 5,
	}
	logger, _ := zap.NewDevelopment()
	srv, _ := NewServer(config, logger)
	return srv
}
