package cmd

import (
	"core/internal/api"
	"core/internal/authentication"
	"core/internal/database"
	"fmt"
	"reflect"
	"regexp"
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
	Log            LogConfig             `mapstructure:"database" validate:"dive"`
	Database       database.Config       `mapstructure:"database" validate:"dive"`
	API            api.Config            `mapstructure:"api" validate:"dive"`
	Authentication authentication.Config `mapstructure:"authentication" validate:"dive"`
}

type LogConfig struct {
	Environment string `mapstructure:"environment" validate:"required"`
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
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		zap.L().Fatal("Failed to read config", zap.Error(err))
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func loadConfig(configs ...string) (*Config, error) {
	var config Config
	defaults.SetDefaults(&config)
	bindEnvs(config)

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	match := regexp.MustCompile(`.*`)
	if len(configs) != 0 {
		match = regexp.MustCompile(strings.ToLower(fmt.Sprintf("^Config.(%s)", strings.Join(configs, "|"))))
	}

	err = validator.New().StructFiltered(config, func(ns []byte) bool {
		return !match.MatchString(strings.ToLower(string(ns)))
	})
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func GetLogger(config LogConfig) *zap.SugaredLogger {
	var logger *zap.Logger
	var err error
	switch config.Environment {
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

	logger.Info("Logger initialized", zap.String("environment", config.Environment))

	return logger.Sugar()
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
