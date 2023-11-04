package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEchohandler(t *testing.T) {
	cases := []struct {
		url    string
		method string
		data   string
		status int
	}{
		{url: "/echo", method: "GET", status: http.StatusOK},
	}

	srv := NewMockServer()
	handler := http.HandlerFunc(srv.echoHandler)

	for _, c := range cases {
		req, err := http.NewRequest(c.method, c.url, strings.NewReader(c.data))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != c.status {
			t.Errorf("hanlder return wrong status code: got %v want %v", status, c.status)
		}
	}
}
