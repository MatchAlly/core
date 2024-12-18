package api

import (
	"context"
	"core/internal/api/handlers"
	"core/internal/api/middleware"
	"core/internal/authentication"
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

func NewServer(port int, version string, l *zap.SugaredLogger, handler *handlers.Handler, authService authentication.Service) *Server {
	var api huma.API

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(echoMiddleware.Recover())

	// Public API setup
	publicConfig := huma.DefaultConfig("MatchAlly", version)
	publicConfig.Servers = []*huma.Server{{URL: "https://matchally.me/public"}}
	publicGroup := e.Group("/public")
	publicAPI := humaecho.NewWithGroup(e, publicGroup, publicConfig)
	publicAPI.UseMiddleware(middleware.CanonicalLogger(l))
	addPublicRoutes(publicAPI, handler)

	// Authenticated API setup
	authenticatedConfig := huma.DefaultConfig("MatchAlly", version)
	authenticatedConfig.Servers = []*huma.Server{{URL: "https://matchally.me/api"}}
	authenticatedGroup := e.Group("/api")
	authenticatedAPI := humaecho.NewWithGroup(e, authenticatedGroup, authenticatedConfig)
	authenticatedAPI.UseMiddleware(
		middleware.CanonicalLogger(l),
		middleware.Authenticated(authService),
	)
	addAuthenticatedRoutes(authenticatedAPI, handler)

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
