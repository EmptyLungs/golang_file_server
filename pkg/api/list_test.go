package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func TestListHandler_OK(t *testing.T) {
	assert := assert.New(t)
	mockfs := new(MockFileManager)
	srv := NewMockServer(mockfs)

	cases := []struct {
		files      []string
		err        error
		statusCode int
	}{
		{
			files:      []string{"test.txt", "asdasd.json"},
			err:        nil,
			statusCode: http.StatusOK,
		},
		{
			files:      nil,
			err:        errors.New("Failed to upload file"),
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, c := range cases {
		mockfs.On("List").Return(c.files, c.err).Once()
		req, err := http.NewRequest("GET", "/list", nil)
		req.Header.Add("Authorization", "Bearer test")
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		srv.handler.ServeHTTP(rr, req)
		assert.Equal(rr.Code, c.statusCode, "Handler returned wrong status")
		if c.err != nil {
			var response ErrorResponse
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Fatalf(err.Error())
			}
			assert.Equal(response.Error, c.err.Error(), "Wrong error message in response")
		} else {
			var responseFiles []string
			if err := json.Unmarshal(rr.Body.Bytes(), &responseFiles); err != nil {
				t.Fatalf(err.Error())
			}
			assert.Equal(responseFiles, c.files, "Wrong list returned in response")
		}
	}
}
