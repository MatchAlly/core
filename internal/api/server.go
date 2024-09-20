package api

import (
	"context"
	"core/internal/api/handlers"
	"core/internal/api/helpers"
	"core/internal/authentication"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Config struct {
	Port int `mapstructure:"port" default:"8080"`
}

type Server struct {
	echo *echo.Echo
	port int
	l    *zap.SugaredLogger
}

func NewServer(
	port int,
	l *zap.SugaredLogger,
	handler *handlers.Handler,
	authService authentication.Service,
) (*Server, error) {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true
	e.Validator = helpers.NewValidator()

	e.Use(
		middleware.Recover(),
		middleware.Logger(),
		middleware.GzipWithConfig(middleware.GzipConfig{
			Skipper: middleware.DefaultGzipConfig.Skipper,
		}),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowCredentials: true,
			AllowOrigins: []string{
				"http://localhost:5173", // dev
				"https://matchally.me/", // prod
			},
		}),
	)

	Register(
		handler,
		e.Group(""),
		l.With("module", "api"),
		authService,
	)

	return &Server{
		echo: e,
		port: port,
		l:    l,
	}, nil
}

func (s *Server) Start() error {
	address := fmt.Sprintf("0.0.0.0:%d", s.port)
	if err := s.echo.Start(address); err != nil {
		return errors.Wrap(err, "Failed to start server")
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.echo.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "Failed to shutdown server")
	}

	return nil
}
