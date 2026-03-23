package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL        string
	Port               string
	JWTSecret          string
	RefreshTokenSecret string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found (continuing)")
	}

	cfg := &Config{
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		Port:               os.Getenv("PORT"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		RefreshTokenSecret: os.Getenv("REFRESH_TOKEN_SECRET"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	if cfg.RefreshTokenSecret == "" {
		log.Fatal("REFRESH_TOKEN_SECRET is required")
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return cfg
}
