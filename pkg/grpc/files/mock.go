package grpc

import (
	"context"
	"io"
	"net"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
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

func NewMockServer() (FileServiceClient, *MockFileManager, func()) {
	logger, _ := zap.NewDevelopment()
	fm := &MockFileManager{}
	server := &FileServer{logger: logger, fm: fm}
	lis := bufconn.Listen(1024 * 1024)

	srv := grpc.NewServer()
	RegisterFileServiceServer(srv, server)

	go func() {
		if err := srv.Serve(lis); err != nil {
			logger.Fatal(err.Error())
		}
	}()

	closer := func() {
		if err := lis.Close(); err != nil {
			logger.Fatal(err.Error())
		}
		srv.Stop()
	}
	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal(err.Error())
	}
	client := NewFileServiceClient(conn)
	return client, fm, closer
}
