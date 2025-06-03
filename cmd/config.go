package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	DatabaseDSN        string        `mapstructure:"DATABASE_DSN"`
	RedisPort          int           `mapstructure:"REDIS_PORT"`
	DenylistExpiry     time.Duration `mapstructure:"DENYLIST_EXPIRY"`
	APIPort            int           `mapstructure:"API_PORT"`
	APIVersion         string        `mapstructure:"API_VERSION"`
	AuthNSecret        string        `mapstructure:"AUTHN_SECRET"`
	AuthNAccessExpiry  time.Duration `mapstructure:"AUTHN_ACCESS_EXPIRY"`
	AuthNRefreshExpiry time.Duration `mapstructure:"AUTHN_REFRESH_EXPIRY"`
	Pepper             string        `mapstructure:"PEPPER"`
}

func loadConfig() (*Config, error) {
	var cfg Config
	err := load(".env", &cfg)
	if err != nil {
		if !os.IsNotExist(err) {
			slog.Info("No .env file found, using strictly environment variables")
		} else {
			return nil, err
		}
	}

	return &cfg, nil
}

// Load reads configuration from .env file and environment variables
// into the provided config struct. Environment variables take precedence
// over .env file values. The struct should have `mapstructure` tags.
func load(configPath string, cfg any) error {
	// Load .env file if it exists
	envMap := make(map[string]string)
	if err := loadEnvFile(configPath, envMap); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	// Populate the config struct
	return populateConfig(envMap, cfg)
}

// loadEnvFile reads a .env file and returns a map of key/value pairs
func loadEnvFile(filePath string, envMap map[string]string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if len(value) > 1 && (value[0] == '"' && value[len(value)-1] == '"' ||
			value[0] == '\'' && value[len(value)-1] == '\'') {
			value = value[1 : len(value)-1]
		}

		envMap[key] = value
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .env file: %w", err)
	}

	return nil
}

// populateConfig fills the config struct with values from envMap and environment variables
func populateConfig(envMap map[string]string, cfg any) error {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("config must be a non-nil pointer")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		// Get the environment variable name from the mapstructure tag
		tag := field.Tag.Get("mapstructure")
		if tag == "" {
			continue
		}

		// Get value first from actual environment, then fall back to .env file
		value := os.Getenv(tag)
		if value == "" {
			value = envMap[tag]
			if value == "" {
				continue // Skip if no value is found
			}
		}

		// Set the field value based on its type
		if err := setFieldValue(fieldValue, value); err != nil {
			return fmt.Errorf("failed to set field %s: %w", field.Name, err)
		}
	}

	return nil
}

// setFieldValue converts and sets the field value based on the field type
func setFieldValue(fieldValue reflect.Value, value string) error {
	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(value)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if fieldValue.Type() == reflect.TypeOf(time.Duration(0)) {
			duration, err := time.ParseDuration(value)
			if err != nil {
				return fmt.Errorf("invalid duration value %q: %w", value, err)
			}
			fieldValue.Set(reflect.ValueOf(duration))
		} else {
			intValue, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid integer value %q: %w", value, err)
			}
			fieldValue.SetInt(intValue)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid unsigned integer value %q: %w", value, err)
		}
		fieldValue.SetUint(uintValue)

	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid float value %q: %w", value, err)
		}
		fieldValue.SetFloat(floatValue)

	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid boolean value %q: %w", value, err)
		}
		fieldValue.SetBool(boolValue)

	default:
		return fmt.Errorf("unsupported field type: %s", fieldValue.Kind())
	}

	return nil
}
