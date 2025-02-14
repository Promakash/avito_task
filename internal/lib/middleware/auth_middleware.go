package middleware

import (
	"avito_shop/internal/domain"
	"avito_shop/internal/usecases"
	"avito_shop/pkg/http/handlers"
	resp "avito_shop/pkg/http/responses"
	"context"
	"errors"
	"net/http"
)

var ErrContextParsing = errors.New("can't parse from context")

type UserIDCtxKey string

const AuthContextKey UserIDCtxKey = "user_id"

func WithTokenAuth(authService usecases.Auth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				handlers.WriteResponse(w, r, resp.Unauthorized(errors.New("token is empty")))
				return
			}

			userID, err := authService.ParseToken(token)
			if err != nil {
				handlers.WriteResponse(w, r, resp.Unauthorized(errors.New("invalid token")))
				return
			}

			ctx := context.WithValue(r.Context(), AuthContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(r *http.Request) (domain.UserID, error) {
	id, ok := r.Context().Value(AuthContextKey).(domain.UserID)
	if !ok {
		return 0, ErrContextParsing
	}

	return id, nil
}
