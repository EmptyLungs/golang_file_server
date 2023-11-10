package api

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type statusCodeInterceptor struct {
	http.ResponseWriter
	statusCode int
}

func (sci *statusCodeInterceptor) WriteHeader(code int) {
	sci.statusCode = code
	sci.ResponseWriter.WriteHeader(code)
}

type LoggingMiddleware struct {
	logger *zap.Logger
}

func NewLoggingMiddleware(logger *zap.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

func (m *LoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		statusCodeWrapper := &statusCodeInterceptor{w, http.StatusOK}
		next.ServeHTTP(statusCodeWrapper, r)

		m.logger.Info(
			"request done",
			zap.String("proto", r.Proto),
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Int("code", statusCodeWrapper.statusCode),
			zap.String("remote", r.RemoteAddr),
			zap.String("user-agent", r.UserAgent()),
			zap.Duration("duration", time.Since(start)),
		)
	})
}
