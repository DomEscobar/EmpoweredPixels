package config

import "os"

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
		jwtSecret = "dev-secret"
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
