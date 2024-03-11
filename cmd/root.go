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
	Log            LogConfig             `mapstructure:"log"`
	Database       database.Config       `mapstructure:"database"`
	API            api.Config            `mapstructure:"api"`
	Authentication authentication.Config `mapstructure:"authentication"`
}

type LogConfig struct {
	Level string `mapstructure:"level" validate:"required"`
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
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err == nil {
		zap.L().Info("Using config file", zap.String("file", viper.ConfigFileUsed()))
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	fmt.Println("init config done")
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
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("load config done")
	return &config, nil
}

func GetLogger(config LogConfig) *zap.SugaredLogger {
	l, err := zap.NewProduction()
	if err != nil {
		zap.L().Fatal("Failed to build logger", zap.Error(err))
	}

	l.Info("Logger initialized", zap.String("level", config.Level))

	return l.Sugar()
}

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
