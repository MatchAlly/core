package api

import (
	"context"
	"core/internal/api/handlers"
	"core/internal/api/middleware"
	"core/internal/authentication"
	"core/internal/cache"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Server struct {
	port int
	e    *echo.Echo
	api  huma.API
	l    *zap.SugaredLogger
}

func NewServer(port int, version string, l *zap.SugaredLogger, handler *handlers.Handler, authService authentication.Service, cacheService cache.Service) *Server {
	var api huma.API

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(echoMiddleware.Recover())

	config := huma.Config{
		OpenAPI: &huma.OpenAPI{
			Info: &huma.Info{
				Title:   "MatchAlly",
				Version: "1.0.0",
			},
			Servers: []*huma.Server{{URL: "https://matchally.me/api"}},
			Components: &huma.Components{
				SecuritySchemes: map[string]*huma.SecurityScheme{
					"bearerAuth": {
						Type:         "http",
						Scheme:       "bearer",
						BearerFormat: "JWT",
					},
				},
			},
		},
	}

	baseAPI := humaecho.NewWithGroup(e, e.Group("/api"), config)

	publicAPI := baseAPI
	publicAPI.UseMiddleware(middleware.CanonicalLogger(l))
	addPublicRoutes(publicAPI, handler)

	authAPI := baseAPI
	openapi := authAPI.OpenAPI()
	openapi.Security = append(openapi.Security, map[string][]string{"bearerAuth": {}})
	authAPI.UseMiddleware(
		middleware.CanonicalLogger(l),
		middleware.Authenticated(authService, cacheService),
	)
	addAuthRoutes(authAPI, handler)

	return &Server{
		port: port,
		e:    e,
		api:  api,
		l:    l,
	}
}

func (s *Server) Start() error {
	address := fmt.Sprintf("0.0.0.0:%d", s.port)
	if err := s.e.Start(address); err != nil {
		return errors.Wrap(err, "failed to start server")
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.e.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "failed to shutdown server")
	}

	return nil
}
