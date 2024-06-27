package cmd

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mcuadros/go-defaults"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const shutdownPeriod = 15 * time.Second

type Config struct {
	LogLevel           string        `mapstructure:"LOG_LEVEL" oneof:"debug info warn error" default:"info"`
	DatabaseDSN        string        `mapstructure:"DATABASE_DSN" validate:"required"`
	APIPort            int           `mapstructure:"API_PORT" default:"8080"`
	AuthNSecret        string        `mapstructure:"AUTHN_SECRET" default:"secret" `
	AuthNAccessExpiry  time.Duration `mapstructure:"AUTHN_ACCESS_EXPIRY" default:"3600"`
	AuthNRefreshExpiry time.Duration `mapstructure:"AUTHN_REFRESH_EXPIRY" default:"3600"`
}

var rootCmd = &cobra.Command{
	Use: "core",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		zap.L().Fatal("Failed to execute root command", zap.Error(err))
	}
}

func loadConfig() (*Config, error) {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		zap.L().Info("no .env file found, using strictly environment variables", zap.Error(err))
	}

	viper.AutomaticEnv()

	var config Config
	defaults.SetDefaults(&config)
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	if err := validator.New().Struct(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func getLogger(level string) *zap.SugaredLogger {
	l, err := zap.NewProduction()
	if err != nil {
		zap.L().Fatal("Failed to build logger", zap.Error(err))
	}

	l.Info("Logger initialized", zap.String("level", level))

	return l.Sugar()
}
