package grpc

import (
	context "context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestEchoGRPC(t *testing.T) {
	assert := assert.New(t)
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})
	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})
	logger, _ := zap.NewDevelopment()
	service := NewEchoService(logger)

	RegisterEchoServiceServer(srv, service)
	// todo: fix this somehow
	go func() {
		if err := srv.Serve(lis); err != nil {
			t.Fatalf(err.Error())
		}
	}()
	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	// todo: fix deprecated grpc.WithInsecure
	conn, err := grpc.DialContext(context.Background(), "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf(err.Error())
	}
	stream := NewEchoServiceClient(conn)
	req := &EchoRequest{Message: "wassup"}
	resp, err := stream.Echo(context.Background(), req)
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.Equal("wassup", resp.Message, "Wrong message returned")
}
