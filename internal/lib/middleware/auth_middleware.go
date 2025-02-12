package middleware

import (
	"avito_shop/internal/usecases"
	"context"
	"errors"
	"net/http"
)

var ErrContextParsing = errors.New("can't parse from context")

const AuthContextKey = "user_id"

func WithTokenAuth(authService usecases.Auth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Unauthorized: token is empty", http.StatusUnauthorized)
				return
			}

			userID, err := authService.ParseToken(token)
			if err != nil {
				http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), AuthContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
