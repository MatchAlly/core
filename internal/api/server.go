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

	"go.uber.org/zap"
)

type Config struct {
	Port    int
	Version string
}

type Server struct {
	config Config
	e      *echo.Echo
	api    huma.API
	l      *zap.SugaredLogger
}

func NewServer(config Config, version string, l *zap.SugaredLogger, handler *handlers.Handler, authService authentication.Service, cacheService cache.Service) *Server {
	var api huma.API

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(echoMiddleware.Recover())

	humaConfig := huma.DefaultConfig("MatchAlly", config.Version)
	humaConfig.OpenAPI.Servers = []*huma.Server{{URL: "https://matchally.me/api"}}
	humaConfig.OpenAPI.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearerAuth": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}

	baseAPI := humaecho.NewWithGroup(e, e.Group("/api"), humaConfig)

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
		config: config,
		e:      e,
		api:    api,
		l:      l,
	}
}

func (s *Server) Start() error {
	address := fmt.Sprintf("0.0.0.0:%d", s.config.Port)
	if err := s.e.Start(address); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.e.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}
