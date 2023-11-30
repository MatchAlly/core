package cmd

import (
	"context"
	"core/internal/authentication"
	"core/internal/club"
	"core/internal/leaderboard"
	"core/internal/match"
	"core/internal/rating"
	"core/internal/rest"
	"core/internal/statistic"
	"core/internal/user"
	"core/pkg/database"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command.
var serveCmd = &cobra.Command{
	Use:  "serve",
	Long: "Start the service",
	Run:  serve,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	config := loadConfig()

	l := GetLogger(config.LogEnv)

	// Initialize database connection
	db, err := database.NewClient(ctx, config.DBDSN)
	if err != nil {
		l.Fatal("Failed to connect to database",
			"error", err)
	}

	// Initialize User service
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	// Initialize Club service
	clubRepository := club.NewRepository(db)
	clubService := club.NewService(clubRepository)

	// Initialize Authentication service
	authenticationService := authentication.NewService(config.JWTSecret, userService)

	// Initialize Match service
	matchRepository := match.NewRepository(db)
	matchService := match.NewService(matchRepository)

	// Initialize Rating service
	ratingRepository := rating.NewRepository(db)
	ratingService := rating.NewService(ratingRepository)

	// Initialize Statistics service
	statisticRepository := statistic.NewRepository(db)
	statisticService := statistic.NewService(statisticRepository)

	// Initialize Leaderboard service
	leaderboardService := leaderboard.NewService(clubService, userService, ratingService, statisticService)

	// Initialize REST server
	restServer, err := rest.NewServer(
		config.Address,
		config.Port,
		l,
		authenticationService,
		userService,
		clubService,
		matchService,
		ratingService,
		statisticService,
		leaderboardService,
	)
	if err != nil {
		l.Fatal("Failed to create rest server",
			"error", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Start the REST server
	fmt.Println("starting server...")
	l.Infow("REST server starting",
		"port", config.Port)
	go func() {
		if err := restServer.Start(); err != nil {
			l.Fatal("Failed to start rest server",
				"error", err)
			cancel()
		}
	}()

	l.Info("Ready")

	fmt.Println("ready")
	// Wait for shutdown signal
	<-ctx.Done()

	// Stop the servers
	l.Info("Shutting down")

	shutdownctx, stop := context.WithTimeout(context.Background(), shutdownPeriod)
	defer stop()

	if err := restServer.Shutdown(shutdownctx); err != nil {
		l.Error("Failed to shutdown rest server",
			"error", err)
	}
}
