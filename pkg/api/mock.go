package api

import (
	"io"
	"testing"
	"time"

	"github.com/EmptyLungs/golang_file_server/pkg/files"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockFileManager struct {
	mock.Mock
}

func (m *MockFileManager) Create(file io.Reader, filename string) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockFileManager) Delete(filename string) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockFileManager) List() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func NewMockServer(fileManager files.IFileManager) *Server {
	config := &Config{
		Host:                  "",
		Port:                  "8080",
		HttpServerTimeout:     5 * time.Second,
		UploaderMaxFileSizeMB: 5,
	}
	logger, _ := zap.NewDevelopment()
	srv, _ := NewServer(config, logger, fileManager)
	return srv
}

func Setup(t *testing.T) (*assert.Assertions, *MockFileManager, *Server) {
	mockfs := new(MockFileManager)
	return assert.New(t), mockfs, NewMockServer(mockfs)
}
