package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	ctx := context.Background()

	// Check if any league exists
	var count int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM leagues").Scan(&count)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	if count > 0 {
		fmt.Printf("Found %d existing leagues. Skipping initialization.\n", count)
		return
	}

	fmt.Println("Initializing test league...")

	// Create a test league
	name := "Founders Cup"
	options := `{"description": "The inaugural competition for the masters of the pixels.", "tier": "Legendary", "prizePool": "1000 Gold"}`
	isDeactivated := false

	var id int
	err = pool.QueryRow(ctx, "INSERT INTO leagues (name, options, is_deactivated) VALUES ($1, $2, $3) RETURNING id", name, options, isDeactivated).Scan(&id)
	if err != nil {
		log.Fatalf("Insert failed: %v", err)
	}

	fmt.Printf("Successfully created league '%s' with ID %d\n", name, id)
}
