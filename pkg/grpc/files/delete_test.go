package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDeleteHandler(t *testing.T) {
	assert := assert.New(t)
	client, mfs, closer := NewMockServer()
	t.Cleanup(closer)
	mfs.On("Delete").Return(nil)
	response, err := client.Delete(context.Background(), &DeleteRequest{Filename: "test.txt"})
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.Equal("ok", response.GetMessage(), "Wrong message returned")
}

func TestDeleteHandlerFail(t *testing.T) {
	assert := assert.New(t)
	client, mfs, closer := NewMockServer()
	t.Cleanup(closer)
	mfs.On("Delete").Return(errors.New("failed to delete"))
	_, err := client.Delete(context.Background(), &DeleteRequest{Filename: "test.txt"})
	expected := status.Error(codes.Internal, "failed to delete")
	assert.Equal(expected, err, "Wrong error message")
}
