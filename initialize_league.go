package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbURL := "postgres://ep_user:test_session_password_1234567890_abcdefgh@127.0.0.1:5432/empoweredpixels?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	name := "Emerald Challenger League"
	options := `{"tier": "rare", "description": "A league for aspiring heroes.", "prizePool": "500 Gold"}`
	isDeactivated := false

	var id int
	err = pool.QueryRow(context.Background(), "INSERT INTO leagues (name, options, is_deactivated) VALUES ($1, $2, $3) RETURNING id", name, options, isDeactivated).Scan(&id)
	if err != nil {
		log.Fatalf("QueryRow failed: %v\n", err)
	}

	fmt.Printf("Created league with ID: %d\n", id)
}
