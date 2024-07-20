package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	cases := []struct {
		name           string
		token          string
		requestToken   string
		expectedStatus int
	}{
		{
			name:           "valid token",
			token:          "123",
			requestToken:   "123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "empty token both server and request",
			token:          "",
			requestToken:   "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "empty server token non-empty request",
			token:          "",
			requestToken:   "123123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid token",
			token:          "123",
			requestToken:   "321",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			middlware := NewAuthMiddleware(tt.token)

			req, err := http.NewRequest("GET", "/test", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Authorization", "Bearer "+tt.requestToken)
			rr := httptest.NewRecorder()
			handler := middlware.Handler(mockHandler)
			handler.ServeHTTP(rr, req)
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
