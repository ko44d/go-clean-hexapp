package config

import (
	"fmt"
	"os"
	"strconv"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type HTTPConfig struct {
	Port int
}

type Config struct {
	DB   DBConfig
	HTTP HTTPConfig
}

func Load() (*Config, error) {
	cfg := &Config{}
	var err error

	cfg.DB.Host, err = lookupRequiredEnv("POSTGRES_HOST")
	if err != nil {
		return nil, err
	}
	cfg.DB.Port, err = lookupRequiredEnvInt("POSTGRES_PORT")
	if err != nil {
		return nil, err
	}
	cfg.DB.User, err = lookupRequiredEnv("POSTGRES_USER")
	if err != nil {
		return nil, err
	}
	cfg.DB.Password, err = lookupRequiredEnv("POSTGRES_PASSWORD")
	if err != nil {
		return nil, err
	}
	cfg.DB.Name, err = lookupRequiredEnv("POSTGRES_DB")
	if err != nil {
		return nil, err
	}
	cfg.DB.SSLMode, err = lookupRequiredEnv("POSTGRES_SSLMODE")
	if err != nil {
		return nil, err
	}

	cfg.HTTP.Port = lookupEnvInt("PORT", 8080)

	return cfg, nil
}

func (c Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.Name, c.DB.SSLMode,
	)
}

func lookupEnv(key string, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func lookupRequiredEnv(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return "", fmt.Errorf("required environment variable %s is unset or empty", key)
	}

	return value, nil
}

func lookupEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func lookupRequiredEnvInt(key string) (int, error) {
	value, err := lookupRequiredEnv(key)
	if err != nil {
		return 0, err
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("required environment variable %s must be a valid integer: %w", key, err)
	}

	return parsedValue, nil
}
