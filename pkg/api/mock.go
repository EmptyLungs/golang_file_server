package api

import (
	"mime/multipart"
	"os"
	"path"
	"time"

	"github.com/EmptyLungs/golang_file_server/pkg/files"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockFileManager struct {
	mock.Mock
}

func (m *MockFileManager) Create(file multipart.File, handler *multipart.FileHeader) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockFileManager) Delete() error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockFileManager) List() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func NewMockServer(fileManager files.IFileManager) *Server {
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
	srv, _ := NewServer(config, logger, fileManager)
	return srv
}
