package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config is the main app config
type Config struct {
	Server   Server
	Database Database
}

type Server struct {
	Name string `env:"APP_NAME"`
	Port int    `env:"APP_PORT"`
	Env  string `env:"APP_ENV"`
}

type Database struct {
	Host     string `env:"DB_HOST"`
	Name     string `env:"DB_NAME"`
	Port     int    `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	SSLMode  string `env:"DB_SSLMODE"`
	DSN      string
}

// Load loads .env file (if exists) + returns populated Config
func Load() (*Config, error) {
	// Optional: load .env file (ignored if not exists)
	_ = godotenv.Load()

	cfg := &Config{
		Server: Server{
			Name: getEnv("APP_NAME", "paybridge-transaction-service"),
			Port: getEnvInt("APP_PORT", 8083),
			Env:  getEnv("APP_ENV", "development"),
		},
		Database: Database{
			Host:     getEnv("DB_HOST", "localhost"),
			Name:     getEnv("DB_NAME", ""),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}

	cfg.Database.DSN = cfg.buildPostgresDSN()

	return cfg, nil
}

func (c *Config) buildPostgresDSN() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// Helper: get string from env or fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Helper: get int from env or fallback
func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}
