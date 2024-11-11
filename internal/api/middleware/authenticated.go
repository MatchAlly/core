package middleware

import (
	"core/internal/authentication"
	"net/http"
	"strconv"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

type contextKey string

const (
	userIDKey contextKey = "user_id"
	nameKey   contextKey = "name"
)

func Authenticated(authenticationService authentication.Service) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		header := ctx.Header("Authorization")
		if len(header) == 0 {
			ctx.SetStatus(http.StatusUnauthorized)
			return
		}

		typ, token, ok := strings.Cut(header, " ")
		if !ok || typ != "Bearer" {
			ctx.SetStatus(http.StatusUnauthorized)
			return
		}

		valid, claims, err := authenticationService.VerifyAccessToken(ctx.Context(), token)
		if !valid || err != nil {
			ctx.SetStatus(http.StatusUnauthorized)
			return
		}

		userID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			ctx.SetStatus(http.StatusUnauthorized)
			return
		}

		ctx = huma.WithValue(ctx, userIDKey, userID)
		ctx = huma.WithValue(ctx, nameKey, claims.Name)

		next(ctx)
	}
}
