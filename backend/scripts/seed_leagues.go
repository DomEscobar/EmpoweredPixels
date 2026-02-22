package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"empoweredpixels/internal/domain/leagues"
	"empoweredpixels/internal/infra/db"
	"empoweredpixels/internal/infra/db/repositories"
)

func main() {
	databaseURL := os.Getenv("EP_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/empoweredpixels?sslmode=disable"
	}

	ctx := context.Background()
	database, err := db.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
	defer database.Pool.Close()

	leagueRepo := repositories.NewLeagueRepository(database.Pool)

	// Check if "Emerald League" exists
	existing, err := leagueRepo.List(ctx)
	if err != nil {
		log.Fatalf("failed to list leagues: %v", err)
	}

	for _, l := range existing {
		if l.Name == "Emerald League" {
			fmt.Println("Emerald League already exists")
			return
		}
	}

	options := []byte(`{"tier": "standard", "description": "The proving grounds for rising champions.", "prizePool": "1000 Gold", "botCount": 5, "botPowerlevel": 30}`)
	league := &leagues.League{
		Name:          "Emerald League",
		Options:       options,
		IsDeactivated: false,
	}

	if err := leagueRepo.Create(ctx, league); err != nil {
		log.Fatalf("failed to create league: %v", err)
	}

	fmt.Printf("Created Emerald League with ID: %d\n", league.ID)
}
