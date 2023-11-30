package cmd

import (
	"time"

	"github.com/caarlos0/env"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const shutdownPeriod = 15 * time.Second

type Config struct {
	LogEnv        string        `env:"LOG_ENV"`
	DBDSN         string        `env:"DB_DSN"`
	Port          int           `env:"PORT" envDefault:"8000"`
	JWTSecret     string        `env:"JWT_SECRET"`
	JWTExpiration time.Duration `env:"JWT_EXPIRATION"`
}

var rootCmd = &cobra.Command{
	Use: "core",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		zap.L().Fatal("Failed to execute root command", zap.Error(err))
	}
}

func init() { //nolint:gochecknoinits
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	err := godotenv.Load()
	if err != nil {
		zap.L().Warn("Error loading .env file")
	}
}

func loadConfig() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		zap.L().Warn("Error loading config from environment", zap.Error(err))
	}

	return &cfg
}

func GetLogger(logEnv string) *zap.SugaredLogger {
	var logger *zap.Logger
	var err error
	switch logEnv {
	case "dev":
		logger, err = zap.NewDevelopment()
	case "prod":
		logger, err = zap.NewProduction()
	default:
		logger, err = zap.NewProduction()
	}
	if err != nil {
		zap.L().Fatal("Failed to build logger", zap.Error(err))
	}

	logger.Info("Logger initialized",
		zap.String("logEnv", logEnv),
	)

	return logger.Sugar()
}
