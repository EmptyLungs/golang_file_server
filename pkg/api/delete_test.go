package api

import (
	"bytes"
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteHandler(t *testing.T) {
	assert, mfs, srv := Setup(t)
	mfs.On("Delete").Return(nil).Once()
	payload := Payload{Filename: "test.txt"}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf(err.Error())
	}
	req, err := http.NewRequest("POST", "/delete", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	srv.handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusNoContent, rr.Code, "Hanlder return wrong status code")
}

func TestDeleteHandlerNotExists(t *testing.T) {
	assert, mfs, srv := Setup(t)
	mfs.On("Delete").Return(fs.ErrNotExist).Once()
	payload := Payload{Filename: "test.txt"}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf(err.Error())
	}
	req, err := http.NewRequest("POST", "/delete", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	srv.handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusNotFound, rr.Code, "Hanlder return wrong status code")
}

func TestDeleteHandlerEmptyBody(t *testing.T) {
	assert, _, srv := Setup(t)
	req, err := http.NewRequest("POST", "/delete", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	// srv.echoHandler(rr, req)
	srv.handler.ServeHTTP(rr, req)
	assert.Equal(rr.Code, http.StatusBadRequest, "Hanlder return wrong status code")
	var response ErrorResponse
	body, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if err = json.Unmarshal(body, &response); err != nil {
		t.Fatalf(err.Error())
	}
	assert.Equal("Empty request body", response.Error, "Wrong error message")
}

func TestDeleteHandlerWrongJsonPayload(t *testing.T) {
	assert, _, srv := Setup(t)
	var badPayload struct {
		Test string `json:"test"`
	}
	jsonData, err := json.Marshal(badPayload)
	if err != nil {
		t.Fatalf(err.Error())
	}

	req, err := http.NewRequest("POST", "/delete", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	// srv.echoHandler(rr, req)
	srv.handler.ServeHTTP(rr, req)
	assert.Equal(rr.Code, http.StatusBadRequest, "Hanlder returned wrong status code")
	var response ErrorResponse
	body, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if err = json.Unmarshal(body, &response); err != nil {
		t.Fatalf(err.Error())
	}
	assert.Equal("Missing filename in request body", response.Error, "Wrong error message")
}
func TestDeleteHandlerBadJsonPayload(t *testing.T) {
	assert, _, srv := Setup(t)
	badJsonData := []byte("\n\n\nasdasdasddas")
	req, err := http.NewRequest("POST", "/delete", bytes.NewBuffer(badJsonData))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	// srv.echoHandler(rr, req)
	srv.handler.ServeHTTP(rr, req)
	assert.Equal(rr.Code, http.StatusBadRequest, "Hanlder returned wrong status code")
	var response ErrorResponse
	body, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if err = json.Unmarshal(body, &response); err != nil {
		t.Fatalf(err.Error())
	}
	assert.Equal("Failed to parse request body", response.Error, "Wrong error message")
}
