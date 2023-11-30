package rest

import (
	"context"
	"core/internal/authentication"
	"core/internal/club"
	"core/internal/leaderboard"
	"core/internal/match"
	"core/internal/rating"
	"core/internal/rest/controllers"
	"core/internal/rest/helpers"
	"core/internal/statistic"
	"core/internal/user"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Config struct {
	Port int `mapstructure:"port" default:"8000"`
}

type Server struct {
	echo    *echo.Echo
	address string
	port    int
}

func NewServer(
	address string,
	port int,
	logger *zap.SugaredLogger,
	authService authentication.Service,
	userService user.Service,
	clubService club.Service,
	matchService match.Service,
	ratingService rating.Service,
	statisticService statistic.Service,
	leaderboardService leaderboard.Service,
) (*Server, error) {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true
	e.Validator = helpers.NewValidator()

	e.Use(
		middleware.Recover(),
		middleware.Logger(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
		}),
		middleware.GzipWithConfig(middleware.GzipConfig{
			Skipper: middleware.DefaultGzipConfig.Skipper,
		}),
	)

	controllers.Register(
		e.Group("/api/v1"),
		logger.With("module", "rest"),
		authService,
		userService,
		clubService,
		matchService,
		ratingService,
		statisticService,
		leaderboardService,
	)

	return &Server{
		echo:    e,
		address: address,
		port:    port,
	}, nil
}

func (s *Server) Start() error {
	err := s.echo.Start(fmt.Sprintf("%s:%d", s.address, s.port))
	if err != nil {
		return errors.Wrap(err, "Failed to start server")
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	err := s.echo.Shutdown(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to shutdown server")
	}

	return nil
}
