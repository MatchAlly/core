package cmd

import (
	"context"
	"core/internal/api"
	"core/internal/api/handlers"
	"core/internal/authentication"
	"core/internal/authorization"
	"core/internal/cache"
	"core/internal/club"
	"core/internal/database"
	"core/internal/game"
	"core/internal/match"
	"core/internal/member"
	"core/internal/rating"
	"core/internal/subscription"
	"core/internal/user"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

const shutdownPeriod = 15 * time.Second

// StartAPIserver initializes and starts the API server
func StartAPIserver(l *slog.Logger) {
	ctx := context.Background()

	config, err := loadConfig()
	if err != nil {
		slog.Error("Failed to read config", "error", err)
		os.Exit(1)
	}

	// Initialize connections to dependencies
	db, err := database.NewClient(ctx, config.DatabaseDSN)
	if err != nil {
		l.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	client := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("redis:%d", config.RedisPort)})

	cacheService := cache.NewService(client, config.DenylistExpiry)

	// Initialize services
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository, config.Pepper)

	memberRepository := member.NewRepository(db)
	memberService := member.NewService(memberRepository)

	subscriptionRepository := subscription.NewRepository(db)
	subscriptionService := subscription.NewService(subscriptionRepository)

	clubRepository := club.NewRepository(db)
	clubService := club.NewService(clubRepository, memberService, subscriptionService)

	authenticationConfig := authentication.Config{
		Secret:        config.AuthNSecret,
		AccessExpiry:  config.AuthNAccessExpiry,
		RefreshExpiry: config.AuthNRefreshExpiry,
		Pepper:        config.Pepper,
	}
	authenticationService := authentication.NewService(authenticationConfig, userService, subscriptionService, cacheService)

	authorizationService := authorization.NewService(memberService)

	gameRepository := game.NewRepository(db)
	gameService := game.NewService(gameRepository)

	matchRepository := match.NewRepository(db)
	matchService := match.NewService(matchRepository, gameService)

	ratingRepository := rating.NewRepository(db)
	ratingService := rating.NewService(ratingRepository)

	// Initialize API server
	handlerConfig := handlers.Config{}

	apiConfig := api.Config{
		Port:    config.APIPort,
		Version: config.APIVersion,
	}

	handler := handlers.NewHandler(l, handlerConfig, authenticationService, authorizationService, userService, clubService, memberService, matchService, ratingService, gameService, subscriptionService)
	apiServer := api.NewServer(apiConfig, config.APIVersion, l, handler, authenticationService, cacheService)
	if err != nil {
		l.Error("Failed to create api server", "error", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Start the API server
	l.Info("API server starting", "port", config.APIPort, "version", config.APIVersion)
	go func() {
		if err := apiServer.Start(); err != nil {
			l.Error("Failed to start api server", "error", err)
			cancel()
		}
	}()

	l.Info("Ready")

	<-ctx.Done()

	l.Info("Shutting down")

	shutdownctx, stop := context.WithTimeout(context.Background(), shutdownPeriod)
	defer stop()

	if err := apiServer.Shutdown(shutdownctx); err != nil {
		l.Error("Failed to shutdown api server", "error", err)
	}
}
