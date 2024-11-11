package api

import (
	"context"
	"core/internal/api/handlers"
	"core/internal/api/middleware"
	"core/internal/authentication"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Server struct {
	port   int
	router *chi.Mux
	api    huma.API
	l      *zap.SugaredLogger
}

func NewServer(port int, version string, l *zap.SugaredLogger, handler *handlers.Handler, authService authentication.Service) *Server {
	var api huma.API

	router := chi.NewMux()
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.Recoverer)

	router.Route("/public", func(r chi.Router) {
		config := huma.DefaultConfig("MatchAlly", version)
		config.Servers = []*huma.Server{
			{URL: "https://matchally.me/public"},
		}
		api = humachi.New(r, config)
		api.UseMiddleware(middleware.CanonicalLogger(l))

		addPublicRoutes(api, handler)
	})

	router.Route("/api", func(r chi.Router) {
		config := huma.DefaultConfig("MatchAlly", version)
		config.Servers = []*huma.Server{
			{URL: "https://matchally.me/api"},
		}
		api = humachi.New(r, config)
		api.UseMiddleware(
			middleware.Authenticated(authService),
			middleware.CanonicalLogger(l),
		)

		addAuthenticatedRoutes(api, handler)
	})

	return &Server{
		port:   port,
		router: router,
		api:    api,
		l:      l,
	}
}

func (s *Server) Start() error {
	address := fmt.Sprintf("0.0.0.0:%d", s.port)
	if err := http.ListenAndServe(address, s.router); err != nil {
		return errors.Wrap(err, "failed to start api server")
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	// TODO implement gracefull shutdown
	return nil
}
