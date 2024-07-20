package api

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockRandomBytesFile(size int64, filename string, formField string) (*bytes.Buffer, string) {
	randomData := make([]byte, size<<20)
	_, _ = rand.Read(randomData)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if formField == "" {
		formField = "file"
	}
	multipart, _ := writer.CreateFormFile(formField, filename)
	multipart.Write(randomData)
	writer.Close()
	return body, writer.FormDataContentType()
}

func TestUploadHandler_OK(t *testing.T) {
	cases := []struct {
		fileSize  int64
		fileName  string
		formField string
		status    int
	}{
		{fileSize: 2, fileName: "small.txt", status: http.StatusOK},
		{fileSize: 2, fileName: "/root/small.txt", status: http.StatusOK},
		{fileSize: 155, fileName: "big.txt", status: http.StatusBadRequest},
		{fileSize: 0, fileName: "empty.txt", status: http.StatusBadRequest},
		{fileSize: 1, formField: "test", fileName: "empty.txt", status: http.StatusBadRequest},
	}
	assert, mockfs, srv := Setup(t)
	mockfs.On("Create").Return(nil)
	for _, c := range cases {
		file, contentType := mockRandomBytesFile(c.fileSize, c.fileName, c.formField)

		t.Logf(fmt.Sprintf("Uploading file %s with size of %d bytes", c.fileName, len(file.Bytes())))
		req, err := http.NewRequest("POST", "/upload", file)
		req.Header.Add("Content-Type", contentType)
		req.Header.Add("Authorization", "Bearer test")
		if err != nil {
			t.Fatalf(err.Error())
		}
		rr := httptest.NewRecorder()
		srv.handler.ServeHTTP(rr, req)
		assert.Equal(rr.Code, c.status, "Hanlder return wrong status code")
	}
}

func TestUploadHandler_FailFS(t *testing.T) {
	assert, mockfs, srv := Setup(t)
	mockfs.On("Create").Return(errors.New("Test"))

	file, ct := mockRandomBytesFile(srv.config.UploaderMaxFileSizeMB-1, "test.json", "file")
	req, err := http.NewRequest("POST", "/upload", file)
	if err != nil {
		t.Fatalf(err.Error())
	}
	req.Header.Add("Content-Type", ct)
	req.Header.Add("Authorization", "Bearer test")
	rr := httptest.NewRecorder()
	srv.handler.ServeHTTP(rr, req)
	assert.Equal(rr.Code, http.StatusInternalServerError, "Handler returned wrong status")
	var response ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf(err.Error())
	}
	assert.Equal("Failed to upload file", response.Error, "Wrong error message returned")
}

func TestUploadHandler_FailForm(t *testing.T) {
	assert, _, srv := Setup(t)
	req, err := http.NewRequest("POST", "/upload", nil)
	if err != nil {
		t.Fatalf(err.Error())
	}
	req.Header.Add("Authorization", "Bearer test")
	rr := httptest.NewRecorder()
	srv.handler.ServeHTTP(rr, req)
	assert.Equal(rr.Code, http.StatusBadRequest, "Handler returned wrong status")
	var response ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf(err.Error())
	}
	assert.Equal("Failed to parse form data", response.Error, "Wrong error message returned")
}
