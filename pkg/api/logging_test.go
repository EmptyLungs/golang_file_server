package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogEntry struct {
	Message   string  `json:"message"`
	Proto     string  `json:"proto"`
	Uri       string  `json:"uri"`
	Method    string  `json:"method"`
	Remote    string  `json:"remote"`
	UserAgent string  `json:"user-agent"`
	Duration  float64 `json:"duration"`
	Source    string  `json:"source"`
}

type TestingLogSink struct {
	entries []string
}

func (s *TestingLogSink) Write(p []byte) (n int, err error) {
	s.entries = append(s.entries, string(p))
	return len(p), nil
}

func TestLoggingMiddleware(t *testing.T) {
	assert := assert.New(t)
	sink := &TestingLogSink{}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(sink),
		zap.InfoLevel,
	)
	logger := zap.New(core)
	srv, _ := NewServer(&Config{}, logger, new(MockFileManager))
	req, err := http.NewRequest("GET", "/echo", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	srv.handler.ServeHTTP(rr, req)
	assert.Equal(rr.Code, http.StatusOK, "Hanlder return wrong status code")
	var log LogEntry
	if err = json.Unmarshal([]byte(sink.entries[0]), &log); err != nil {
		t.Fatalf(err.Error())
	}
}
