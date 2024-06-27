package cmd

import (
	"reflect"
	"strings"
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
	v := viper.New()

	v.SetConfigFile(".env")
	if err := v.ReadInConfig(); err != nil {
		zap.L().Info("no .env file found, using strictly environment variables", zap.Error(err))
	}

	v.AutomaticEnv()

	var config Config
	defaults.SetDefaults(&config)
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}
	bindEnvs(config)

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

// Adapted from https://github.com/spf13/viper/issues/188#issuecomment-401431526
func bindEnvs(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		fieldv := ifv.Field(i)
		t := ift.Field(i)
		name := strings.ToLower(t.Name)
		tag, ok := t.Tag.Lookup("mapstructure")
		if ok {
			name = tag
		}
		parts := append(parts, name)
		switch fieldv.Kind() { //nolint:exhaustive
		case reflect.Struct:
			bindEnvs(fieldv.Interface(), parts...)
		default:
			viper.BindEnv(strings.Join(parts, ".")) //nolint:errcheck
		}
	}
}
