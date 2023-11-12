package grpc

import (
	"io"

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

func NewMockServer() (*FileServer, *MockFileManager) {
	logger, _ := zap.NewDevelopment()
	fm := &MockFileManager{}
	return &FileServer{logger: logger, fm: fm}, fm
}
