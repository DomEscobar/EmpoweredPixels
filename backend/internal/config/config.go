package config

import (
	"log"
	"os"
)

type Config struct {
	HTTPAddress string
	DatabaseURL string
	JWTSecret   string
	TokenDays   int
	EngineURL   string
}

func FromEnv() Config {
	address := os.Getenv("EP_HTTP_ADDRESS")
	if address == "" {
		address = ":54321"
	}

	databaseURL := os.Getenv("EP_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/empoweredpixels?sslmode=disable"
	}

	jwtSecret := os.Getenv("EP_JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("FATAL: EP_JWT_SECRET environment variable is required and must be at least 32 characters. Set a strong secret before starting.")
	}
	if len(jwtSecret) < 32 {
		log.Fatalf("FATAL: EP_JWT_SECRET must be at least 32 characters long (got %d). Use a cryptographically strong value.", len(jwtSecret))
	}

	tokenDays := 7

	engineURL := os.Getenv("EP_ENGINE_URL")

	return Config{
		HTTPAddress: address,
		DatabaseURL: databaseURL,
		JWTSecret:   jwtSecret,
		TokenDays:   tokenDays,
		EngineURL:   engineURL,
	}
}
