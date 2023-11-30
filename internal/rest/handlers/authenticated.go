package handlers

import (
	"core/internal/authentication"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type AuthenticatedContext struct {
	echo.Context
	Log    *zap.SugaredLogger
	Claims authentication.AccessClaims
	JWT    string
}

type AuthenticatedHandlerFunc func(ctx AuthenticatedContext) error

func AuthenticatedHandlerFactory(logger *zap.SugaredLogger) func(handler AuthenticatedHandlerFunc) func(ctx echo.Context) error {
	return func(handler AuthenticatedHandlerFunc) func(ctx echo.Context) error {
		return func(ctx echo.Context) error {
			claims, ok := ctx.Get("jwt_claims").(*authentication.AccessClaims)
			if !ok {
				logger.Debug("missing jwt_claims")

				return echo.ErrUnauthorized
			}

			jwt, ok := ctx.Get("jwt").(string)
			if !ok {
				logger.Debug("missing jwt")

				return echo.ErrUnauthorized
			}

			logger = logger.With(
				"user_id", claims.Subject,
				"name", claims.Name,
			)

			return handler(AuthenticatedContext{
				Context: ctx,
				Claims:  *claims,
				Log:     logger,
				JWT:     jwt,
			})
		}
	}
}
