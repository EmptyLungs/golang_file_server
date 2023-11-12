package grpc

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestUploadHandler(t *testing.T) {
	assert := assert.New(t)
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})
	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})
	service, fm := NewMockServer()

	RegisterFileServiceServer(srv, service)
	// todo: fix this somehow
	go func() {
		if err := srv.Serve(lis); err != nil {
			t.Fatalf(err.Error())
		}
	}()

	fm.On("Create").Return(nil).Once()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf(err.Error())
	}
	client := NewFileServiceClient(conn)
	stream, err := client.Upload(context.Background())
	if err != nil {
		t.Fatalf(err.Error())
	}

	fileData := []byte("test file content")
	fileName := "test.txt"
	req := &UploadRequest{
		Chunk:    fileData,
		Filename: fileName,
	}

	if err := stream.Send(req); err != nil {
		t.Fatalf(err.Error())
	}
	response, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.Equal("ok", response.GetMessage(), "Wrong message returned")
}
