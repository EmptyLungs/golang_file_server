package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEchohandler(t *testing.T) {
	assert := assert.New(t)
	cases := []struct {
		url    string
		method string
		status int
	}{
		{url: "/echo", method: "GET", status: http.StatusOK},
	}
	mockfs := new(MockFileManager)
	srv := NewMockServer(mockfs)
	for _, c := range cases {
		req, err := http.NewRequest(c.method, c.url, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		srv.echoHandler(rr, req)
		assert.Equal(rr.Code, c.status, "Hanlder return wrong status code")
	}
}
