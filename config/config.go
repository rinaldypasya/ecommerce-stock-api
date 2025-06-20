package config

import (
	"log"
	"os"
)

type Config struct {
	DBDSN     string
	JWTSecret string
}

var AppConfig *Config

func Load() {
	AppConfig = &Config{
		DBDSN:     getEnv("DB_DSN", "postgres://postgres:postgres@localhost:5432/ecommerce?sslmode=disable"),
		JWTSecret: getEnv("JWT_SECRET", "dev-secret"),
	}
	log.Println("✅ Configuration loaded.")
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
