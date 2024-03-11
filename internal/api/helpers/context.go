package helpers

import (
	"core/internal/authentication"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type AuthContext struct {
	echo.Context
	Log    *zap.SugaredLogger
	Claims authentication.AccessClaims
	JWT    string
}

type AuthenticatedContextFunc func(ctx AuthContext) error

func AuthenticatedContextFactory(l *zap.SugaredLogger) func(handler AuthenticatedContextFunc) func(ctx echo.Context) error {
	return func(handler AuthenticatedContextFunc) func(ctx echo.Context) error {
		return func(ctx echo.Context) error {
			claims, ok := ctx.Get("jwt_claims").(*authentication.AccessClaims)
			if !ok {
				l.Debug("missing jwt_claims")

				return echo.ErrUnauthorized
			}

			jwt, ok := ctx.Get("jwt").(string)
			if !ok {
				l.Debug("missing jwt")

				return echo.ErrUnauthorized
			}

			l = l.With(
				"user_id", claims.Subject,
				"name", claims.Name,
			)

			return handler(AuthContext{
				Context: ctx,
				Claims:  *claims,
				Log:     l,
				JWT:     jwt,
			})
		}
	}
}
