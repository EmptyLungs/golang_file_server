package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUploadHandler(t *testing.T) {
	assert := assert.New(t)
	client, mfs, closer := NewMockServer()
	t.Cleanup(closer)
	mfs.On("Create").Return(nil)
	stream, err := client.Upload(context.Background())
	if err != nil {
		t.Fatalf(err.Error())
	}

	firstChunk := []byte("first")
	secondChunk := []byte("second")
	expectedSize := uint32(len(firstChunk) + len(secondChunk))
	fileName := "test.txt"
	req1 := &UploadRequest{
		Chunk:    firstChunk,
		Filename: fileName,
	}
	req2 := &UploadRequest{
		Chunk:    secondChunk,
		Filename: "",
	}
	if err := stream.Send(req1); err != nil {
		t.Fatalf(err.Error())
	}
	if err := stream.Send(req2); err != nil {
		t.Fatalf(err.Error())
	}
	response, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.Equal("ok", response.GetMessage(), "wrong message")
	assert.Equal(expectedSize, response.GetFilesize(), "worng file size")
	assert.Equal(fileName, response.GetFilename(), "wrong file name")
}
func TestUploadHandlerFail_FMError(t *testing.T) {
	assert := assert.New(t)
	client, mfs, closer := NewMockServer()
	t.Cleanup(closer)
	mfs.On("Create").Return(errors.New("failed to create"))
	fileData := []byte("test file content")
	fileName := "test.txt"
	req := &UploadRequest{
		Chunk:    fileData,
		Filename: fileName,
	}

	stream, err := client.Upload(context.Background())
	if err != nil {
		t.Fatalf(err.Error())
	}
	if err := stream.Send(req); err != nil {
		t.Fatalf(err.Error())
	}
	_, err = stream.CloseAndRecv()
	expected := status.Error(codes.Internal, "failed to create")
	assert.Equal(expected, err, "failed to create")
}
func TestUploadHandlerFail_MissingFilename(t *testing.T) {
	assert := assert.New(t)
	client, mfs, closer := NewMockServer()
	t.Cleanup(closer)
	mfs.On("Create").Return(nil)
	fileData := []byte("test file content")
	fileName := ""
	req := &UploadRequest{
		Chunk:    fileData,
		Filename: fileName,
	}

	stream, err := client.Upload(context.Background())
	if err != nil {
		t.Fatalf(err.Error())
	}
	if err := stream.Send(req); err != nil {
		t.Fatalf(err.Error())
	}
	_, err = stream.CloseAndRecv()
	expected := status.Error(codes.InvalidArgument, "missing file name")
	assert.Equal(expected, err, "missing file name")
}
