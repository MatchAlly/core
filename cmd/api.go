package cmd

import (
	"context"
	xapi "core/internal/api"
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
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var apiCmd = &cobra.Command{
	Use:  "api",
	Long: "Start the api server",
	Run:  api,
}

func init() {
	rootCmd.AddCommand(apiCmd)
}

func api(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()

	config, err := loadConfig()
	if err != nil {
		zap.L().Fatal("failed to read config", zap.Error(err))
	}

	l := getLogger()

	// Initialize database connection
	db, err := database.NewClient(ctx, config.DatabaseDSN)
	if err != nil {
		l.Fatal("failed to connect to database", zap.Error(err))
	}

	client := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("redis:%d", config.RedisPort)})

	cacheService := cache.NewService(client, config.DenylistExpiry)

	// Initialize Services
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	clubRepository := club.NewRepository(db)
	clubService := club.NewService(clubRepository)

	memberRepository := member.NewRepository(db)
	memberService := member.NewService(memberRepository)

	authenticationConfig := authentication.Config{
		Secret:        config.AuthNSecret,
		AccessExpiry:  config.AuthNAccessExpiry,
		RefreshExpiry: config.AuthNRefreshExpiry,
		Pepper:        config.Pepper,
	}
	authenticationService := authentication.NewService(authenticationConfig, userService, cacheService)

	authorizationService := authorization.NewService(memberService)

	matchRepository := match.NewRepository(db)
	matchService := match.NewService(matchRepository)

	ratingRepository := rating.NewRepository(db)
	ratingService := rating.NewService(ratingRepository)

	gameRepository := game.NewRepository(db)
	gameService := game.NewService(gameRepository)

	subscriptionRepository := subscription.NewRepository(db)
	subscriptionService := subscription.NewService(subscriptionRepository)

	// Initialize API server
	handlerConfig := handlers.Config{}

	apiConfig := xapi.Config{
		Port:    config.APIPort,
		Version: config.APIVersion,
	}

	handler := handlers.NewHandler(l, handlerConfig, authenticationService, authorizationService, userService, clubService, memberService, matchService, ratingService, gameService, subscriptionService)
	apiServer := xapi.NewServer(apiConfig, config.APIVersion, l, handler, authenticationService, cacheService)
	if err != nil {
		l.Fatal("failed to create api server", zap.Error(err))
	}

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Start the API server
	l.Info("api server starting", zap.Int("port", config.APIPort), zap.String("version", config.APIVersion))
	go func() {
		if err := apiServer.Start(); err != nil {
			l.Fatal("failed to start api server", zap.Error(err))
			cancel()
		}
	}()

	l.Info("ready")

	// Wait for shutdown signal
	<-ctx.Done()

	// Stop the servers
	l.Info("shutting down")

	shutdownctx, stop := context.WithTimeout(context.Background(), shutdownPeriod)
	defer stop()

	if err := apiServer.Shutdown(shutdownctx); err != nil {
		l.Error("failed to shutdown api server", zap.Error(err))
	}
}
