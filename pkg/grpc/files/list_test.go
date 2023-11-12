package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestListHandler(t *testing.T) {
	assert := assert.New(t)
	client, mfs, closer := NewMockServer()
	t.Cleanup(closer)
	mfs.On("List").Return([]string{"test.txt", "123.json"}, nil)

	resp, err := client.List(context.Background(), &emptypb.Empty{})
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.ElementsMatch([]string{"test.txt", "123.json"}, resp.Files, "Wrong response")
}

func TestListHandlerError(t *testing.T) {
	assert := assert.New(t)
	client, mfs, closer := NewMockServer()
	t.Cleanup(closer)
	mfs.On("List").Return([]string{}, errors.New("test"))
	_, err := client.List(context.Background(), &emptypb.Empty{})
	grpcErr := status.Error(codes.Internal, "test")
	assert.Equal(grpcErr, err)
}
