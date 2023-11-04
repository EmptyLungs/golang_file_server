package api

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockRandomBytesFile(size int64, filename string) (*bytes.Buffer, string) {
	randomData := make([]byte, size<<20)
	_, _ = rand.Read(randomData)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	multipart, _ := writer.CreateFormFile("file", filename)
	multipart.Write(randomData)
	writer.Close()
	return body, writer.FormDataContentType()
}

func TestUploadHandler(t *testing.T) {
	assert := assert.New(t)
	cases := []struct {
		url      string
		method   string
		fileSize int64
		fileName string
		status   int
	}{
		{url: "/upload", method: "POST", fileSize: 2, fileName: "small.txt", status: http.StatusOK},
		{url: "/upload", method: "POST", fileSize: 155, fileName: "big.txt", status: http.StatusBadRequest},
		{url: "/upload", method: "POST", fileSize: 0, fileName: "empty.txt", status: http.StatusBadRequest},
	}

	srv := NewMockServer()
	handler := http.HandlerFunc(srv.uploadFileHandler)

	for _, c := range cases {
		file, contentType := mockRandomBytesFile(c.fileSize, c.fileName)
		t.Logf(fmt.Sprintf("Uploading file %s with size of %d bytes", c.fileName, len(file.Bytes())))
		req, err := http.NewRequest(c.method, c.url, file)
		req.Header.Add("Content-Type", contentType)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		assert.Equal(rr.Code, c.status, "Hanlder return wrong status code:\nGot %v want %v\nResponse: %s", rr.Code, c.status, rr.Body.String())
	}
}
