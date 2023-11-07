package api

import (
	"os"
	"path"
	"time"

	"go.uber.org/zap"
)

func NewMockServer() *Server {
	dir, _ := os.Getwd()
	testDataDir := path.Join(dir, "test-data")
	os.RemoveAll(testDataDir)
	config := &Config{
		Host:                  "",
		Port:                  "8080",
		HttpServerTimeout:     5 * time.Second,
		UploaderDir:           testDataDir,
		UploaderMaxFileSizeMB: 5,
	}
	logger, _ := zap.NewDevelopment()
	srv, _ := NewServer(config, logger)
	return srv
}
