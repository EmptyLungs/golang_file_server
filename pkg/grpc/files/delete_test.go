package grpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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
