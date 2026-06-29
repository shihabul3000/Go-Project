package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv            string
	Port              string
	DatabaseURL       string
	JWTSecret         string
	JWTExpiresIn      time.Duration
	BcryptCost        int
	AllowedOrigins    []string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		AppEnv:            getString("APP_ENV", "development"),
		Port:              getString("PORT", "8080"),
		DatabaseURL:       getString("DATABASE_URL", ""),
		JWTSecret:         getString("JWT_SECRET", "change-this-secret-before-deploying"),
		JWTExpiresIn:      getDuration("JWT_EXPIRES_IN", 24*time.Hour),
		BcryptCost:        getInt("BCRYPT_COST", 12),
		AllowedOrigins:    getCSV("ALLOWED_ORIGINS", []string{"*"}),
		DBMaxOpenConns:    getInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:    getInt("DB_MAX_IDLE_CONNS", 10),
		DBConnMaxLifetime: getDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute),
	}
}

func getString(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getDuration(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getCSV(key string, fallback []string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			items = append(items, item)
		}
	}

	if len(items) == 0 {
		return fallback
	}
	return items
}
