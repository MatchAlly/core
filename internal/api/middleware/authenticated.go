package middleware

import (
	"core/internal/authentication"
	"core/internal/cache"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

func Authenticated(authenticationService authentication.Service, cacheService cache.Service) func(ctx huma.Context, next func(huma.Context)) {
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

		used, err := cacheService.GetTokenUsed(ctx.Context(), token)
		if err != nil {
			ctx.SetStatus(http.StatusInternalServerError)
			return
		}
		if used {
			ctx.SetStatus(http.StatusUnauthorized)
			return
		}

		ctx = huma.WithValue(ctx, "claims", claims)
		ctx = huma.WithValue(ctx, "token", token)

		next(ctx)
	}
}
