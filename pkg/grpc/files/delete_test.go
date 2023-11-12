package grpc

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestDeleteHandler(t *testing.T) {
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

	fm.On("Delete").Return(nil).Once()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf(err.Error())
	}
	client := NewFileServiceClient(conn)
	response, err := client.Delete(context.Background(), &DeleteRequest{Filename: "test.txt"})
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.Equal("ok", response.GetMessage(), "Wrong message returned")
}
