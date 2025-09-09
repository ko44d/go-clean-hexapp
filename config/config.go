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

	cfg.DB.Host = lookupEnv("POSTGRES_HOST", "localhost")
	cfg.DB.Port = lookupEnvInt("POSTGRES_PORT", 5432)
	cfg.DB.User = lookupEnv("POSTGRES_USER", "clean-hexuser")
	cfg.DB.Password = lookupEnv("POSTGRES_PASSWORD", "clean-hexpass")
	cfg.DB.Name = lookupEnv("POSTGRES_DB", "clean-hexapp")
	cfg.DB.SSLMode = lookupEnv("POSTGRES_SSLMODE", "disable")

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

func lookupEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
