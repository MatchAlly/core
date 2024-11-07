package middleware

import (
	"context"
	"core/internal/authentication"
	"net/http"
	"strconv"
	"strings"
)

type contextKey string

const (
	userIDKey contextKey = "user_id"
	nameKey   contextKey = "name"
)

func Authorization(authenticationService authentication.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if len(header) == 0 {
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			typ, token, ok := strings.Cut(header, " ")
			if !ok || typ != "Bearer" {
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			valid, claims, err := authenticationService.VerifyAccessToken(r.Context(), token)
			if !valid || err != nil {
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			userID, err := strconv.Atoi(claims.Subject)
			if err != nil {
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			ctx = context.WithValue(ctx, nameKey, claims.Name)

			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
