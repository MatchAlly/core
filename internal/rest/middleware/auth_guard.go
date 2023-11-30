package middleware

import (
	"core/internal/authentication"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthGuard(authenticationService authentication.Service) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if header == "" {
				return echo.ErrUnauthorized
			}

			typ, token, ok := strings.Cut(header, " ")
			if !ok || typ != "Bearer" {
				return echo.ErrUnauthorized
			}

			valid, claims, err := authenticationService.VerifyAccessToken(c.Request().Context(), token)
			if !valid || err != nil {
				return echo.ErrUnauthorized
			}

			c.Set("jwt", token)
			c.Set("jwt_claims", claims)

			return next(c)
		}
	}

}
