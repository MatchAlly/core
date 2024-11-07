package cmd

import (
	"context"
	"core/internal/api"
	"core/internal/api/handlers"
	"core/internal/authentication"
	"core/internal/club"
	"core/internal/database"
	"core/internal/match"
	"core/internal/member"
	"core/internal/rating"
	"core/internal/user"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var serveCmd = &cobra.Command{
	Use:  "serve",
	Long: "Start the service",
	Run:  serve,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) {
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
	}
	authenticationService := authentication.NewService(authenticationConfig, userService)

	matchRepository := match.NewRepository(db)
	matchService := match.NewService(matchRepository)

	ratingRepository := rating.NewRepository(db)
	ratingService := rating.NewService(ratingRepository)

	// Initialize API server
	handler := handlers.NewHandler(l, authenticationService, userService, clubService, memberService, matchService, ratingService)
	apiServer := api.NewServer(config.APIPort, config.APIVersion, l, handler, authenticationService)
	if err != nil {
		l.Fatal("failed to create api server", zap.Error(err))
	}

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Start the API server
	l.Info("api server starting", zap.String("port", fmt.Sprint(config.APIPort)))
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
