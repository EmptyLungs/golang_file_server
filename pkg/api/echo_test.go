package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEchohandler(t *testing.T) {
	assert, _, srv := Setup(t)
	cases := []struct {
		url    string
		method string
		status int
	}{
		{url: "/echo", method: "GET", status: http.StatusOK},
	}
	for _, c := range cases {
		req, err := http.NewRequest(c.method, c.url, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		// srv.echoHandler(rr, req)
		srv.handler.ServeHTTP(rr, req)
		assert.Equal(rr.Code, c.status, "Hanlder return wrong status code")
	}
}
