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
		zap.L().Fatal("failed to execute root command", zap.Error(err))
	}
}

func loadConfig() (*Config, error) {
	viper.SetOptions(viper.ExperimentalBindStruct()) // This is required to bind env vars untill it releases fully in v1.20

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

func getLogger() *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.Encoding = "console"
	l, err := config.Build()
	if err != nil {
		zap.L().Fatal("failed to build logger", zap.Error(err))
	}

	return l.Sugar()
}
