package cmd

import (
	"context"
	"core/internal/club"
	"core/internal/database"
	"core/internal/game"
	"core/internal/match"
	"core/internal/rating"
	"core/internal/user"
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var seedCmd = &cobra.Command{
	Use:  "seed",
	Long: "Seed the database",
	Run:  seed,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(seedCmd)
}

func seed(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	config, err := loadConfig()
	if err != nil {
		zap.L().Fatal("failed to read config", zap.Error(err))
	}

	l := getLogger(config.LogLevel)

	db, err := database.NewClient(ctx, config.DatabaseDSN)
	if err != nil {
		l.Fatal("failed to connect to database", zap.Error(err))
	}

	// Seed the database
	if err := seedUsers(ctx, db, l); err != nil {
		l.Fatal("failed to seed users", zap.Error(err))
	}
	if err := seedClubs(ctx, db, l); err != nil {
		l.Fatal("failed to seed clubs", zap.Error(err))
	}
	if err := seedGames(ctx, db, l); err != nil {
		l.Fatal("failed to seed games", zap.Error(err))
	}
	if err := seedRatings(ctx, db, l); err != nil {
		l.Fatal("failed to seed ratings", zap.Error(err))
	}
	if err := seedMatches(ctx, db, l); err != nil {
		l.Fatal("failed to seed matches", zap.Error(err))
	}

	l.Info("Finsihed seeding database")
}

func seedUsers(ctx context.Context, db *gorm.DB, l *zap.SugaredLogger) error {
	userRepo := user.NewRepository(db)

	f, err := os.Open("internal/database/seeds/users.json")
	if err != nil {
		l.Error("failed to open seed file", zap.Error(err))
	}
	defer f.Close()

	var users []user.User
	if err = json.NewDecoder(f).Decode(&users); err != nil {
		l.Error("failed to decode users file", zap.Error(err))
	}

	for _, u := range users {
		if err := userRepo.CreateUser(ctx, &u); err != nil {
			l.Error("failed to create user", zap.Error(err))
		}
	}

	return nil
}

func seedClubs(ctx context.Context, db *gorm.DB, l *zap.SugaredLogger) error {
	repo := club.NewRepository(db)

	f, err := os.Open("internal/database/seeds/clubs.json")
	if err != nil {
		l.Error("failed to open seed file", zap.Error(err))
	}
	defer f.Close()

	var clubs []club.Club
	if err = json.NewDecoder(f).Decode(&clubs); err != nil {
		l.Error("failed to decode clubs file", zap.Error(err))
	}

	for _, c := range clubs {
		if _, err := repo.CreateClub(ctx, &c); err != nil {
			l.Error("failed to create club", zap.Error(err))
		}
	}

	return nil
}

func seedGames(ctx context.Context, db *gorm.DB, l *zap.SugaredLogger) error {
	repo := game.NewRepository(db)

	f, err := os.Open("internal/database/seeds/games.json")
	if err != nil {
		l.Error("failed to open seed file", zap.Error(err))
	}
	defer f.Close()

	var games []game.Game
	if err = json.NewDecoder(f).Decode(&games); err != nil {
		l.Error("failed to decode games file", zap.Error(err))
	}

	for _, g := range games {
		if err := repo.CreateGame(ctx, &g); err != nil {
			l.Error("failed to create game", zap.Error(err))
		}
	}

	return nil
}

func seedRatings(ctx context.Context, db *gorm.DB, l *zap.SugaredLogger) error {
	repo := rating.NewRepository(db)

	f, err := os.Open("internal/database/seeds/ratings.json")
	if err != nil {
		l.Error("failed to open seed file", zap.Error(err))
	}
	defer f.Close()

	var ratings []rating.Rating
	if err = json.NewDecoder(f).Decode(&ratings); err != nil {
		l.Error("failed to decode ratings file", zap.Error(err))
	}

	for _, r := range ratings {
		if err := repo.CreateRating(ctx, &r); err != nil {
			l.Error("failed to create rating", zap.Error(err))
		}
	}

	return nil
}

func seedMatches(ctx context.Context, db *gorm.DB, l *zap.SugaredLogger) error {
	repo := match.NewRepository(db)

	f, err := os.Open("internal/database/seeds/matches.json")
	if err != nil {
		l.Error("failed to open seed file", zap.Error(err))
	}
	defer f.Close()

	var matches []match.Match
	if err = json.NewDecoder(f).Decode(&matches); err != nil {
		l.Error("failed to decode matches file", zap.Error(err))
	}

	for _, m := range matches {
		if err := repo.CreateMatch(ctx, &m); err != nil {
			l.Error("failed to create match", zap.Error(err))
		}
	}

	return nil
}
