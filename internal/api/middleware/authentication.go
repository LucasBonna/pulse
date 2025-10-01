package middleware

import (
	"lucasbonna/pulse/internal/utils"
	"net/http"
	"strings"
)

func AuthenticationMiddleware(validToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				utils.WriteJsonError(w, http.StatusUnauthorized, "Authorization header is required")
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				utils.WriteJsonError(w, http.StatusUnauthorized, "Authorization header must start with 'Bearer '")
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			if strings.TrimSpace(token) == "" {
				utils.WriteJsonError(w, http.StatusUnauthorized, "Bearer token is required")
				return
			}

			if token != validToken {
				utils.WriteJsonError(w, http.StatusUnauthorized, "Invalid Token")
			}

			next.ServeHTTP(w, r)
		})
	}
}
