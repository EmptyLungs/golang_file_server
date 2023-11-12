package grpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadHandlerw(t *testing.T) {
	assert := assert.New(t)
	client, mfs, closer := NewMockServer()
	t.Cleanup(closer)
	mfs.On("Create").Return(nil)
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
