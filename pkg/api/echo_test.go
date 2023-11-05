package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEchohandler(t *testing.T) {
	assert := assert.New(t)
	cases := []struct {
		url    string
		method string
		data   string
		status int
	}{
		{url: "/echo", method: "GET", status: http.StatusOK},
	}

	srv := NewMockServer()
	for _, c := range cases {
		req, err := http.NewRequest(c.method, c.url, strings.NewReader(c.data))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		srv.echoHandler(rr, req)
		assert.Equal(rr.Code, c.status, "hanlder return wrong status code:\ngot %v want %v", rr.Code, c.status)
	}
}
