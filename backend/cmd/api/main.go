package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	httpadapter "empoweredpixels/internal/adapter/http"
	"empoweredpixels/internal/adapter/ws"
	"empoweredpixels/internal/config"
	"empoweredpixels/internal/infra/db"
	"empoweredpixels/internal/infra/db/repositories"
	"empoweredpixels/internal/infra/engine"
	"empoweredpixels/internal/infra/jobs"
	identityusecase "empoweredpixels/internal/usecase/identity"
	inventoryusecase "empoweredpixels/internal/usecase/inventory"
	leaguesusecase "empoweredpixels/internal/usecase/leagues"
	matchesusecase "empoweredpixels/internal/usecase/matches"
	rewardsusecase "empoweredpixels/internal/usecase/rewards"
	rosterusecase "empoweredpixels/internal/usecase/roster"
	seasonsusecase "empoweredpixels/internal/usecase/seasons"
	shopusecase "empoweredpixels/internal/usecase/shop"
	attunementusecase "empoweredpixels/internal/usecase/attunement"
	weaponsusecase "empoweredpixels/internal/usecase/weapons"
	dailyusecase "empoweredpixels/internal/usecase/daily"
	"empoweredpixels/internal/mcp"
)

func main() {
	cfg := config.FromEnv()

	database, err := db.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Printf("database error: %v", err)
		os.Exit(1)
	}
	defer database.Pool.Close()

	migrationsDir := filepath.Join("internal", "infra", "db", "migrations")
	if err := db.ApplyMigrations(context.Background(), database.Pool, migrationsDir); err != nil {
		log.Printf("migration error: %v", err)
		os.Exit(1)
	}

	userRepo := repositories.NewUserRepository(database.Pool)
	tokenRepo := repositories.NewTokenRepository(database.Pool)
	verificationRepo := repositories.NewVerificationRepository(database.Pool)
	identityService := identityusecase.NewService(
		userRepo,
		tokenRepo,
		verificationRepo,
		cfg.JWTSecret,
		cfg.TokenDays,
		time.Now,
	)

	fighterRepo := repositories.NewFighterRepository(database.Pool)
	experienceRepo := repositories.NewExperienceRepository(database.Pool)
	configurationRepo := repositories.NewConfigurationRepository(database.Pool)
	rosterService := rosterusecase.NewService(
		fighterRepo,
		experienceRepo,
		configurationRepo,
		time.Now,
	)

	itemRepo := repositories.NewItemRepository(database.Pool)
	equipmentRepo := repositories.NewEquipmentRepository(database.Pool)
	equipmentOptionRepo := repositories.NewEquipmentOptionRepository(database.Pool)
	inventoryService := inventoryusecase.NewService(itemRepo, equipmentRepo, equipmentOptionRepo, time.Now)

	rewardRepo := repositories.NewRewardRepository(database.Pool)
	rewardService := rewardsusecase.NewService(rewardRepo, itemRepo, equipmentRepo, time.Now)

	matchRepo := repositories.NewMatchRepository(database.Pool)
	matchTeamRepo := repositories.NewMatchTeamRepository(database.Pool)
	matchRegistrationRepo := repositories.NewMatchRegistrationRepository(database.Pool)
	matchResultRepo := repositories.NewMatchResultRepository(database.Pool)
	matchScoreRepo := repositories.NewMatchScoreRepository(database.Pool)
	engineClient := engine.NewClient(engine.Config{BaseURL: cfg.EngineURL})
	matchHub := ws.NewMatchHub()
	matchService := matchesusecase.NewService(
		matchRepo,
		matchTeamRepo,
		matchRegistrationRepo,
		matchResultRepo,
		matchScoreRepo,
		fighterRepo,
		inventoryService,
		rewardService,
		rosterService,
		engineClient,
		matchHub,
		time.Now,
	)

	leagueRepo := repositories.NewLeagueRepository(database.Pool)
	leagueSubRepo := repositories.NewLeagueSubscriptionRepository(database.Pool)
	leagueMatchRepo := repositories.NewLeagueMatchRepository(database.Pool)
	leagueService := leaguesusecase.NewService(leagueRepo, leagueSubRepo, leagueMatchRepo, fighterRepo, time.Now)
	leagueJob := jobs.NewLeagueJob(matchService, leagueRepo, leagueSubRepo, leagueMatchRepo, fighterRepo, 4*time.Hour)
	leagueJob.Start()

	lobbyCleanupJob := jobs.NewLobbyCleanupJob(matchService, 60, 5*time.Minute)
	lobbyCleanupJob.Start()

	seasonSummaryRepo := repositories.NewSeasonSummaryRepository(database.Pool)
	seasonService := seasonsusecase.NewService(seasonSummaryRepo)

	weaponRepo := repositories.NewWeaponRepository(database.Pool)
	weaponService := weaponsusecase.NewService(weaponRepo)

	// Shop service initialization
	shopRepo := repositories.NewShopRepository(database.Pool)
	goldRepo := repositories.NewPlayerGoldRepository(database.Pool)
	txRepo := repositories.NewTransactionRepository(database.Pool)
	shopService := shopusecase.NewService(shopRepo, goldRepo, txRepo)

	// Attunement service initialization
	attunementRepo := repositories.NewAttunementRepository(database.Pool)
	attunementService := attunementusecase.NewService(attunementRepo)

	// Daily reward service initialization
	dailyRepo := repositories.NewDailyRewardRepository(database.Pool)
	dailyService := dailyusecase.NewService(dailyRepo, goldRepo)

	mcpFilter := mcp.NewFairnessFilter(100, 1*time.Minute)
	mcpHandler := mcp.NewMCPHandler(mcpFilter, identityService, rosterService, inventoryService, leagueService, matchService, rewardService)
	mcpAuditLogger, _ := mcp.NewAuditLogger("")

	server := &http.Server{
		Addr: cfg.HTTPAddress,
		Handler: httpadapter.NewRouter(httpadapter.Dependencies{
			Config:           cfg,
			IdentityService:  identityService,
			RosterService:    rosterService,
			MatchService:     matchService,
			InventoryService: inventoryService,
			WeaponService:    weaponService,
			SkillService:        nil,
			ShopService:         shopService,
			AttunementService:   attunementService,
			DailyService:        dailyService,
			LeagueService:       leagueService,
			LeagueJob:        leagueJob,
			RewardService:    rewardService,
			SeasonService:    seasonService,
			MatchHub:         matchHub,
			MCPHandler:       mcpHandler,
			MCPAuditLogger:   mcpAuditLogger,
			MCPFilter:        mcpFilter,
		}),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("api listening on %s", cfg.HTTPAddress)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("server error: %v", err)
		os.Exit(1)
	}
}
