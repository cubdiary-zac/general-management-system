package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	AppEnv            string
	ServerPort        string
	DatabaseURL       string
	JWTSecret         string
	JWTTTLHours       int
	SeedAdminName     string
	SeedAdminEmail    string
	SeedAdminPassword string
}

func Load() Config {
	cfg := Config{
		AppEnv:            getEnv("APP_ENV", "development"),
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		DatabaseURL:       getDatabaseURL(),
		JWTSecret:         getEnv("JWT_SECRET", "dev-only-secret"),
		JWTTTLHours:       getEnvInt("JWT_TTL_HOURS", 24),
		SeedAdminName:     getEnv("SEED_ADMIN_NAME", "System Owner"),
		SeedAdminEmail:    getEnv("SEED_ADMIN_EMAIL", "admin@gms.local"),
		SeedAdminPassword: getEnv("SEED_ADMIN_PASSWORD", "admin123"),
	}

	return cfg
}

func getDatabaseURL() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}

	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	user := getEnv("POSTGRES_USER", "gms")
	password := getEnv("POSTGRES_PASSWORD", "gms")
	dbName := getEnv("POSTGRES_DB", "gms")
	sslMode := getEnv("POSTGRES_SSLMODE", "disable")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbName, sslMode)
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.Atoi(value)
		if err == nil {
			return parsed
		}
	}
	return fallback
}
