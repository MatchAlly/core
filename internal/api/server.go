package api

import (
	"context"
	"core/internal/api/handlers"
	"core/internal/api/middleware"
	"core/internal/authentication"
	"core/internal/cache"
	"fmt"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type Config struct {
	Port    int
	Version string
}

type Server struct {
	config Config
	e      *echo.Echo
	api    huma.API
	l      *slog.Logger
}

func NewServer(config Config, version string, l *slog.Logger, handler *handlers.Handler, authService authentication.Service, cacheService cache.Service) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(echoMiddleware.Recover())

	humaConfig := huma.DefaultConfig("MatchAlly", config.Version)
	humaConfig.OpenAPI.Servers = []*huma.Server{{URL: "http://localhost:8080"}, {URL: "https://matchally.me"}}
	humaConfig.OpenAPI.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearerAuth": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}

	api := humaecho.New(e, humaConfig)
	api.UseMiddleware(middleware.CanonicalLogger(l))

	authGroup := huma.NewGroup(api, "/api/v1")
	authGroup.UseMiddleware(middleware.Authenticated(authService, cacheService))
	addAuthRoutes(authGroup, handler)

	baseGroup := huma.NewGroup(api, "/api")
	addPublicRoutes(baseGroup, handler)

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
