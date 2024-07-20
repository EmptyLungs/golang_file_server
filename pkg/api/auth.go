package api

import (
	"net/http"
)

type AuthMiddleware struct {
	token string
}

func NewAuthMiddleware(token string) *AuthMiddleware {
	return &AuthMiddleware{
		token: token,
	}
}

func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token == "" {
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		if token != "Bearer "+m.token {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
