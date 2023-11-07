package api

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
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
	var expectedFiles []string
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
	for _, c := range cases {
		file, contentType := mockRandomBytesFile(c.fileSize, c.fileName)
		t.Logf(fmt.Sprintf("Uploading file %s with size of %d bytes", c.fileName, len(file.Bytes())))
		req, err := http.NewRequest(c.method, c.url, file)
		req.Header.Add("Content-Type", contentType)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		srv.uploadFileHandler(rr, req)
		assert.Equal(rr.Code, c.status, "Hanlder return wrong status code")
		t.Logf(rr.Body.String())
		if rr.Code == http.StatusOK {
			expectedFiles = append(expectedFiles, c.fileName)
		}
	}

	t.Logf("Requesting list of files after upload")

	req, err := http.NewRequest("GET", "/list", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	srv.listFileHandler(rr, req)

	assert.Equal(200, rr.Code, "List handler returned non-200 status code")
	var responseFiles []string
	t.Logf(rr.Body.String())
	err = json.Unmarshal(rr.Body.Bytes(), &responseFiles)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(expectedFiles, responseFiles, "Got unexpected list files response")
}
