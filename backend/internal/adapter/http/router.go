package httpadapter

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"empoweredpixels/internal/adapter/http/handlers"
	mcphandlers "empoweredpixels/internal/adapter/http/handlers"
	inventoryhandlers "empoweredpixels/internal/adapter/http/handlers/inventory"
	leaguehandlers "empoweredpixels/internal/adapter/http/handlers/leagues"
	matchhandlers "empoweredpixels/internal/adapter/http/handlers/matches"
	rewardhandlers "empoweredpixels/internal/adapter/http/handlers/rewards"
	rosterhandlers "empoweredpixels/internal/adapter/http/handlers/roster"
	seasonhandlers "empoweredpixels/internal/adapter/http/handlers/seasons"
	shophandlers "empoweredpixels/internal/adapter/http/handlers/shop"
	attunementhandlers "empoweredpixels/internal/adapter/http/handlers/attunement"
	dailyhandlers "empoweredpixels/internal/adapter/http/handlers/daily"
	leaderboardhandlers "empoweredpixels/internal/adapter/http/handlers/leaderboard"
	eventhandlers "empoweredpixels/internal/adapter/http/handlers/events"
	guildhandlers "empoweredpixels/internal/adapter/http/handlers/guilds"
	weaponhandlers "empoweredpixels/internal/adapter/http/handlers/weapons"
	skillhandlers "empoweredpixels/internal/adapter/http/handlers/skills"
	"empoweredpixels/internal/adapter/http/middleware"
	"empoweredpixels/internal/adapter/http/responses"
	"empoweredpixels/internal/adapter/ws"
	"empoweredpixels/internal/config"
	"empoweredpixels/internal/infra/jobs"
	"empoweredpixels/internal/mcp"
	"empoweredpixels/internal/usecase/identity"
	inventoryusecase "empoweredpixels/internal/usecase/inventory"
	leaguesusecase "empoweredpixels/internal/usecase/leagues"
	matchesusecase "empoweredpixels/internal/usecase/matches"
	rewardsusecase "empoweredpixels/internal/usecase/rewards"
	rosterusecase "empoweredpixels/internal/usecase/roster"
	seasonsusecase "empoweredpixels/internal/usecase/seasons"
	shopusecase "empoweredpixels/internal/usecase/shop"
	attunementusecase "empoweredpixels/internal/usecase/attunement"
	weaponsusecase "empoweredpixels/internal/usecase/weapons"
	skillsusecase "empoweredpixels/internal/usecase/skills"
	dailyusecase "empoweredpixels/internal/usecase/daily"
	guildsusecase "empoweredpixels/internal/usecase/guilds"
	leaderboardusecase "empoweredpixels/internal/usecase/leaderboard"
	eventsusecase "empoweredpixels/internal/usecase/events"
)

type Dependencies struct {
	Config           config.Config
	IdentityService  *identity.Service
	RosterService    *rosterusecase.Service
	MatchService     *matchesusecase.Service
	InventoryService inventoryusecase.Service
	WeaponService    *weaponsusecase.Service
	SkillService     *skillsusecase.Service
	LeagueService    *leaguesusecase.Service
	LeagueJob        *jobs.LeagueJob
	RewardService    *rewardsusecase.Service
	SeasonService       *seasonsusecase.Service
	ShopService         *shopusecase.Service
	AttunementService   *attunementusecase.Service
	DailyService        *dailyusecase.Service
	LeaderboardService  *leaderboardusecase.Service
	EventService        *eventsusecase.Service
	GuildService        *guildsusecase.Service
	MatchHub            *ws.MatchHub
	MCPHandler       *mcp.MCPHandler
	MCPAuditLogger   *mcp.AuditLogger
	MCPFilter        *mcp.FairnessFilter
}

func NewRouter(deps Dependencies) http.Handler {
	r := mux.NewRouter().StrictSlash(true)
	
	// Create common middleware
	authMiddleware := func(next http.Handler) http.Handler { return next }
	if deps.Config.JWTSecret != "" {
		authMiddleware = func(next http.Handler) http.Handler {
			return middleware.WithUserID(next, []byte(deps.Config.JWTSecret))
		}
	}

	// Health check (no API prefix usually)
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		handlers.Health().ServeHTTP(w, r)
	}).Methods("GET")
	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		handlers.Health().ServeHTTP(w, r)
	}).Methods("GET")

	// API Routes
	api := r.PathPrefix("/api").Subrouter()
	api.Use(func(next http.Handler) http.Handler {
		return authMiddleware(next)
	})

	if deps.IdentityService != nil {
		authHandler := handlers.NewAuthHandler(deps.IdentityService)
		registerHandler := handlers.NewRegisterHandler(deps.IdentityService)
		api.HandleFunc("/authentication/token", authHandler.Token).Methods("POST")
		api.HandleFunc("/authentication/refresh", authHandler.Refresh).Methods("POST")
		api.HandleFunc("/register", registerHandler.Register).Methods("POST")
		api.HandleFunc("/register/verify", registerHandler.Verify).Methods("POST")
	}

	if deps.RosterService != nil {
		h := rosterhandlers.NewFighterHandler(deps.RosterService)
		api.HandleFunc("/fighter", h.List).Methods("GET")
		api.HandleFunc("/fighter", h.Create).Methods("PUT")
		api.HandleFunc("/fighter/{id}", func(w http.ResponseWriter, r *http.Request) {
			h.Get(w, r, mux.Vars(r)["id"])
		}).Methods("GET")
		api.HandleFunc("/fighter/{id}", func(w http.ResponseWriter, r *http.Request) {
			h.Delete(w, r, mux.Vars(r)["id"])
		}).Methods("DELETE")
		api.HandleFunc("/fighter/{id}/name", func(w http.ResponseWriter, r *http.Request) {
			h.GetName(w, r, mux.Vars(r)["id"])
		}).Methods("GET")
		api.HandleFunc("/fighter/{id}/experience", func(w http.ResponseWriter, r *http.Request) {
			h.GetExperience(w, r, mux.Vars(r)["id"])
		}).Methods("GET")
		api.HandleFunc("/fighter/{id}/configuration", func(w http.ResponseWriter, r *http.Request) {
			h.GetConfiguration(w, r, mux.Vars(r)["id"])
		}).Methods("GET")
		api.HandleFunc("/fighter/{id}/configuration", func(w http.ResponseWriter, r *http.Request) {
			h.UpdateConfiguration(w, r, mux.Vars(r)["id"])
		}).Methods("POST")

		// Squad routes
		squadHandler := rosterhandlers.NewSquadHandler(deps.RosterService.SquadService) // Assume SquadService is injected
		api.HandleFunc("/roster/squad/active", squadHandler.GetActive).Methods("GET")
		api.HandleFunc("/roster/squad/active", squadHandler.SetActive).Methods("POST")
	}

	if deps.MatchService != nil {
		h := matchhandlers.NewHandler(deps.MatchService)
		api.HandleFunc("/match/quick-join", h.QuickJoin).Methods("POST")
		api.HandleFunc("/match/online-players", h.GetOnlinePlayers).Methods("GET")
		api.HandleFunc("/match/{id}", func(w http.ResponseWriter, r *http.Request) {
			h.GetMatch(w, r, mux.Vars(r)["id"])
		}).Methods("GET")
		api.HandleFunc("/match/{id}/teams", func(w http.ResponseWriter, r *http.Request) {
			h.GetTeams(w, r, mux.Vars(r)["id"])
		}).Methods("GET")
		api.HandleFunc("/match/browse", h.Browse).Methods("POST")
		api.HandleFunc("/match/join", h.Join).Methods("POST")
		api.HandleFunc("/match/join/team", h.JoinTeam).Methods("POST")
		api.HandleFunc("/match/leave", h.Leave).Methods("POST")
		api.HandleFunc("/match/leave/team", h.LeaveTeam).Methods("POST")
		api.HandleFunc("/match/{id}/start", func(w http.ResponseWriter, r *http.Request) {
			h.Start(w, r, mux.Vars(r)["id"])
		}).Methods("POST")
		api.HandleFunc("/match/{id}/roundticks", func(w http.ResponseWriter, r *http.Request) {
			h.RoundTicks(w, r, mux.Vars(r)["id"])
		}).Methods("GET")
		api.HandleFunc("/match/{id}/fighterscores", func(w http.ResponseWriter, r *http.Request) {
			h.FighterScores(w, r, mux.Vars(r)["id"])
		}).Methods("GET")
		api.HandleFunc("/match/quick-join", h.QuickJoin).Methods("POST")
	}

	if deps.InventoryService != nil {
		h := inventoryhandlers.NewHandler(deps.InventoryService)
		api.HandleFunc("/inventory/balance/particles", h.BalanceParticles).Methods("GET")
		api.HandleFunc("/inventory/balance/token/common", h.BalanceTokenCommon).Methods("GET")
		api.HandleFunc("/inventory/balance/token/rare", h.BalanceTokenRare).Methods("GET")
		api.HandleFunc("/inventory/balance/token/fabled", h.BalanceTokenFabled).Methods("GET")
		api.HandleFunc("/inventory/balance/token/mythic", h.BalanceTokenMythic).Methods("GET")
		api.HandleFunc("/equipment/{id}", func(w http.ResponseWriter, r *http.Request) {
			h.GetEquipment(w, r, mux.Vars(r)["id"])
		}).Methods("GET")
		api.HandleFunc("/equipment/enhance/cost", h.EnhancementCost).Methods("POST")
		api.HandleFunc("/equipment/enhance", h.Enhance).Methods("POST")
		api.HandleFunc("/equipment/salvage", h.Salvage).Methods("POST")
		api.HandleFunc("/equipment/salvage/inventory", h.SalvageInventory).Methods("POST")
		api.HandleFunc("/equipment/inventory", h.InventoryPage).Methods("POST")
		api.HandleFunc("/equipment/fighter/{id}", func(w http.ResponseWriter, r *http.Request) {
			h.ListByFighter(w, r, mux.Vars(r)["id"])
		}).Methods("GET")
		api.HandleFunc("/equipment/{id}/favorite", func(w http.ResponseWriter, r *http.Request) {
			h.SetFavorite(w, r, mux.Vars(r)["id"], true)
		}).Methods("POST")
		api.HandleFunc("/equipment/{id}/favorite", func(w http.ResponseWriter, r *http.Request) {
			h.SetFavorite(w, r, mux.Vars(r)["id"], false)
		}).Methods("DELETE")
	}

	if deps.LeagueService != nil {
		h := leaguehandlers.NewHandler(deps.LeagueService)
		api.HandleFunc("/league", h.List).Methods("GET")
		api.HandleFunc("/league/{id}", func(w http.ResponseWriter, r *http.Request) {
			h.Get(w, r, mux.Vars(r)["id"])
		}).Methods("GET")
		api.HandleFunc("/league/{id}/winner", func(w http.ResponseWriter, r *http.Request) {
			h.LastWinner(w, r, mux.Vars(r)["id"])
		}).Methods("GET")
		api.HandleFunc("/league/subscribe", h.Subscribe).Methods("POST")
		api.HandleFunc("/league/unsubscribe", h.Unsubscribe).Methods("POST")
		api.HandleFunc("/league/{id}/subscriptions", func(w http.ResponseWriter, r *http.Request) {
			h.Subscriptions(w, r, mux.Vars(r)["id"])
		}).Methods("GET")
		api.HandleFunc("/league/{id}/subscriptions/user", func(w http.ResponseWriter, r *http.Request) {
			h.SubscriptionsForUser(w, r, mux.Vars(r)["id"])
		}).Methods("GET")
		api.HandleFunc("/league/{id}/matches", func(w http.ResponseWriter, r *http.Request) {
			h.Matches(w, r, mux.Vars(r)["id"])
		}).Methods("POST")
		api.HandleFunc("/league/{id}/highscores", func(w http.ResponseWriter, r *http.Request) {
			h.Highscores(w, r, mux.Vars(r)["id"])
		}).Methods("POST")
		if deps.LeagueJob != nil {
			api.HandleFunc("/league/{id}/run", func(w http.ResponseWriter, r *http.Request) {
				runLeagueJob(w, r, mux.Vars(r)["id"], deps.LeagueJob)
			}).Methods("POST")
		}
	}

	if deps.RewardService != nil {
		h := rewardhandlers.NewHandler(deps.RewardService)
		api.HandleFunc("/reward", h.List).Methods("GET")
		api.HandleFunc("/reward/claim", h.Claim).Methods("POST")
		api.HandleFunc("/reward/claim/all", h.ClaimAll).Methods("POST")
	}

	if deps.SeasonService != nil {
		h := seasonhandlers.NewHandler(deps.SeasonService)
		api.HandleFunc("/season/summary", h.Summary).Methods("POST")
	}

	if deps.WeaponService != nil {
		h := weaponhandlers.NewHandler(deps.WeaponService)
		api.HandleFunc("/weapons", h.List).Methods("GET")
		api.HandleFunc("/weapons/database", h.Database).Methods("GET")
		api.HandleFunc("/weapons/{id}", func(w http.ResponseWriter, r *http.Request) {
			h.Get(w, r, mux.Vars(r)["id"])
		}).Methods("GET")
		api.HandleFunc("/weapons/equip", h.Equip).Methods("POST")
		api.HandleFunc("/weapons/{id}/unequip", func(w http.ResponseWriter, r *http.Request) {
			h.Unequip(w, r, mux.Vars(r)["id"])
		}).Methods("POST")
		api.HandleFunc("/weapons/enhance", h.Enhance).Methods("POST")
		api.HandleFunc("/weapons/forge", h.ForgePreview).Methods("POST")
	}

	if deps.SkillService != nil {
		h := skillhandlers.NewHandler(deps.SkillService)
		api.HandleFunc("/skills/tree", h.GetSkillTree).Methods("GET")
		api.HandleFunc("/skills/fighter/{id}", func(w http.ResponseWriter, r *http.Request) {
			h.GetFighterSkills(w, r, mux.Vars(r)["id"])
		}).Methods("GET")
		api.HandleFunc("/skills/allocate", h.AllocateSkill).Methods("POST")
		api.HandleFunc("/skills/loadout/{id}", func(w http.ResponseWriter, r *http.Request) {
			h.GetLoadout(w, r, mux.Vars(r)["id"])
		}).Methods("GET")
		api.HandleFunc("/skills/loadout", h.SetLoadout).Methods("POST")
		api.HandleFunc("/skills/reset", h.ResetSkills).Methods("POST")
		api.HandleFunc("/skills/reset-cost/{id}", func(w http.ResponseWriter, r *http.Request) {
			h.GetResetCost(w, r, mux.Vars(r)["id"])
		}).Methods("GET")
	}

	if deps.ShopService != nil {
		h := shophandlers.NewHandler(deps.ShopService)
		api.HandleFunc("/shop/items", h.GetShopItems).Methods("GET")
		api.HandleFunc("/shop/gold", h.GetGoldPackages).Methods("GET")
		api.HandleFunc("/shop/bundles", h.GetBundles).Methods("GET")
		api.HandleFunc("/shop/item/{id}", h.GetShopItem).Methods("GET")
		api.HandleFunc("/shop/purchase", h.Purchase).Methods("POST")
		api.HandleFunc("/player/gold", h.GetPlayerGold).Methods("GET")
		api.HandleFunc("/player/transactions", h.GetTransactions).Methods("GET")
	}

	if deps.AttunementService != nil {
		h := attunementhandlers.NewHandler(deps.AttunementService)
		api.HandleFunc("/attunements", h.GetAttunements).Methods("GET")
		api.HandleFunc("/attunements/bonuses", h.GetBonuses).Methods("GET")
		api.HandleFunc("/attunement/{element}", h.GetAttunement).Methods("GET")
		api.HandleFunc("/attunement/award-xp", h.AwardXP).Methods("POST")
	}

	if deps.DailyService != nil {
		h := dailyhandlers.NewHandler(deps.DailyService)
		api.HandleFunc("/daily-reward", h.GetStatus).Methods("GET")
		api.HandleFunc("/daily-reward/claim", h.Claim).Methods("POST")
	}

	if deps.EventService != nil {
		h := eventhandlers.NewHandler(deps.EventService)
		api.HandleFunc("/events/current", h.GetCurrentEvents).Methods("GET")
		api.HandleFunc("/events/status", h.GetEventStatus).Methods("GET")
		api.HandleFunc("/events/next", h.GetNextEvent).Methods("GET")
	}

	if deps.LeaderboardService != nil {
		h := leaderboardhandlers.NewHandler(deps.LeaderboardService)
		api.HandleFunc("/leaderboard/{category}", h.GetLeaderboard).Methods("GET")
		api.HandleFunc("/leaderboard/{category}/nearby", h.GetNearbyRanks).Methods("GET")
		api.HandleFunc("/achievements", h.GetAchievements).Methods("GET")
		api.HandleFunc("/player/achievements", h.GetPlayerAchievements).Methods("GET")
		api.HandleFunc("/achievement/{id}/claim", h.ClaimAchievement).Methods("POST")
	}

	if deps.GuildService != nil {
		h := guildhandlers.NewHandler(deps.GuildService)
		api.HandleFunc("/guilds", h.List).Methods("GET")
		api.HandleFunc("/guilds", h.Create).Methods("POST")
		api.HandleFunc("/guilds/{id}", h.Get).Methods("GET")
		api.HandleFunc("/guilds/{id}/join", h.RequestJoin).Methods("POST")
	}

	if deps.MatchHub != nil {
		r.Handle("/ws/match", deps.MatchHub)
	}

	if deps.MCPHandler != nil {
		mcpH := mcphandlers.NewMCPHandler(deps.MCPHandler)
		api.HandleFunc("/mcp/tool", mcpH.Call).Methods("POST")

		if deps.MCPAuditLogger != nil && deps.MCPFilter != nil {
			restH := mcphandlers.NewMCPRESTHandler(deps.MCPHandler, deps.MCPAuditLogger, deps.MCPFilter)
			r.HandleFunc("/mcp/gameState", restH.GameState).Methods("GET")
			r.HandleFunc("/mcp/action", restH.SubmitAction).Methods("POST")
			r.HandleFunc("/mcp/player/stats", restH.PlayerStats).Methods("GET")
		}
	}

	return middleware.WithCORS(r)
}

func runLeagueJob(w http.ResponseWriter, r *http.Request, idStr string, job *jobs.LeagueJob) {
	leagueID, err := strconv.Atoi(idStr)
	if err != nil || leagueID <= 0 {
		responses.Error(w, http.StatusBadRequest, "invalid league id")
		return
	}
	if err := job.RunLeague(r.Context(), leagueID); err != nil {
		switch err {
		case jobs.ErrLeagueNotFound, jobs.ErrNoSubscriptions:
			responses.Error(w, http.StatusBadRequest, err.Error())
			return
		default:
			log.Printf("league run error: %v", err)
			responses.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
