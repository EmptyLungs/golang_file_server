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

func TestUploadHandler(t *testing.T) {
	assert := assert.New(t)
	cases := []struct {
		url       string
		method    string
		fileSize  int64
		fileName  string
		formField string
		status    int
	}{
		{fileSize: 2, fileName: "small.txt", status: http.StatusOK},
		{fileSize: 155, fileName: "big.txt", status: http.StatusBadRequest},
		{fileSize: 0, fileName: "empty.txt", status: http.StatusBadRequest},
		{fileSize: 1, formField: "test", fileName: "empty.txt", status: http.StatusBadRequest},
	}

	mockfs := new(MockFileManager)
	mockfs.On("Create").Return(nil)
	srv := NewMockServer(mockfs)
	for _, c := range cases {
		file, contentType := mockRandomBytesFile(c.fileSize, c.fileName, c.formField)

		t.Logf(fmt.Sprintf("Uploading file %s with size of %d bytes", c.fileName, len(file.Bytes())))
		req, err := http.NewRequest("POST", "/upload", file)
		req.Header.Add("Content-Type", contentType)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		srv.uploadFileHandler(rr, req)
		assert.Equal(rr.Code, c.status, "Hanlder return wrong status code")
	}
}
