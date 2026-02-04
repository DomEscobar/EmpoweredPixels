package httpadapter

import (
	"log"
	"net/http"
	"strconv"

	"empoweredpixels/internal/adapter/http/handlers"
	inventoryhandlers "empoweredpixels/internal/adapter/http/handlers/inventory"
	leaguehandlers "empoweredpixels/internal/adapter/http/handlers/leagues"
	matchhandlers "empoweredpixels/internal/adapter/http/handlers/matches"
	rewardhandlers "empoweredpixels/internal/adapter/http/handlers/rewards"
	rosterhandlers "empoweredpixels/internal/adapter/http/handlers/roster"
	seasonhandlers "empoweredpixels/internal/adapter/http/handlers/seasons"
	"empoweredpixels/internal/adapter/http/middleware"
	"empoweredpixels/internal/adapter/http/responses"
	"empoweredpixels/internal/adapter/ws"
	"empoweredpixels/internal/config"
	"empoweredpixels/internal/infra/jobs"
	"empoweredpixels/internal/usecase/identity"
	inventoryusecase "empoweredpixels/internal/usecase/inventory"
	leaguesusecase "empoweredpixels/internal/usecase/leagues"
	matchesusecase "empoweredpixels/internal/usecase/matches"
	rewardsusecase "empoweredpixels/internal/usecase/rewards"
	rosterusecase "empoweredpixels/internal/usecase/roster"
	seasonsusecase "empoweredpixels/internal/usecase/seasons"
)

type Dependencies struct {
	Config           config.Config
	IdentityService  *identity.Service
	RosterService    *rosterusecase.Service
	MatchService     *matchesusecase.Service
	InventoryService *inventoryusecase.Service
	LeagueService    *leaguesusecase.Service
	LeagueJob        *jobs.LeagueJob
	RewardService    *rewardsusecase.Service
	SeasonService    *seasonsusecase.Service
	MatchHub         *ws.MatchHub
}

func NewRouter(deps Dependencies) http.Handler {
	mux := http.NewServeMux()
	authMiddleware := func(next http.Handler) http.Handler { return next }
	if deps.Config.JWTSecret != "" {
		authMiddleware = func(next http.Handler) http.Handler {
			return middleware.WithUserID(next, []byte(deps.Config.JWTSecret))
		}
	}

	mux.Handle("GET /health", handlers.Health())

	if deps.IdentityService != nil {
		authHandler := handlers.NewAuthHandler(deps.IdentityService)
		registerHandler := handlers.NewRegisterHandler(deps.IdentityService)

		mux.Handle("POST /api/authentication/token", authMiddleware(http.HandlerFunc(authHandler.Token)))
		mux.Handle("POST /api/authentication/refresh", authMiddleware(http.HandlerFunc(authHandler.Refresh)))
		mux.Handle("POST /api/register", http.HandlerFunc(registerHandler.Register))
		mux.Handle("POST /api/register/verify", http.HandlerFunc(registerHandler.Verify))
	}

	if deps.RosterService != nil {
		fighterHandler := rosterhandlers.NewFighterHandler(deps.RosterService)

		mux.Handle("GET /api/fighter", authMiddleware(http.HandlerFunc(fighterHandler.List)))
		mux.Handle("PUT /api/fighter", authMiddleware(http.HandlerFunc(fighterHandler.Create)))
		mux.Handle("GET /api/fighter/{id}", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fighterHandler.Get(w, r, r.PathValue("id"))
		})))
		mux.Handle("DELETE /api/fighter/{id}", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fighterHandler.Delete(w, r, r.PathValue("id"))
		})))
		mux.Handle("GET /api/fighter/{id}/name", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fighterHandler.GetName(w, r, r.PathValue("id"))
		})))
		mux.Handle("GET /api/fighter/{id}/experience", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fighterHandler.GetExperience(w, r, r.PathValue("id"))
		})))
		mux.Handle("GET /api/fighter/{id}/configuration", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fighterHandler.GetConfiguration(w, r, r.PathValue("id"))
		})))
		mux.Handle("POST /api/fighter/{id}/configuration", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fighterHandler.UpdateConfiguration(w, r, r.PathValue("id"))
		})))
	}

	if deps.MatchService != nil {
		matchHandler := matchhandlers.NewHandler(deps.MatchService)

		mux.Handle("GET /api/match/options/default", authMiddleware(http.HandlerFunc(matchHandler.GetDefaultOptions)))
		mux.Handle("GET /api/match/options/sizes", authMiddleware(http.HandlerFunc(matchHandler.GetBattlefieldSizes)))
		mux.Handle("GET /api/match/current", authMiddleware(http.HandlerFunc(matchHandler.GetCurrentMatch)))
		mux.Handle("PUT /api/match/create", authMiddleware(http.HandlerFunc(matchHandler.CreateMatch)))
		mux.Handle("PUT /api/match/create/team", authMiddleware(http.HandlerFunc(matchHandler.CreateTeam)))
		mux.Handle("GET /api/match/{id}", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			matchHandler.GetMatch(w, r, r.PathValue("id"))
		})))
		mux.Handle("GET /api/match/{id}/teams", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			matchHandler.GetTeams(w, r, r.PathValue("id"))
		})))
		mux.Handle("POST /api/match/browse", authMiddleware(http.HandlerFunc(matchHandler.Browse)))
		mux.Handle("POST /api/match/join", authMiddleware(http.HandlerFunc(matchHandler.Join)))
		mux.Handle("POST /api/match/join/team", authMiddleware(http.HandlerFunc(matchHandler.JoinTeam)))
		mux.Handle("POST /api/match/leave", authMiddleware(http.HandlerFunc(matchHandler.Leave)))
		mux.Handle("POST /api/match/leave/team", authMiddleware(http.HandlerFunc(matchHandler.LeaveTeam)))
		mux.Handle("POST /api/match/{id}/start", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			matchHandler.Start(w, r, r.PathValue("id"))
		})))
		mux.Handle("GET /api/match/{id}/roundticks", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			matchHandler.RoundTicks(w, r, r.PathValue("id"))
		})))
		mux.Handle("GET /api/match/{id}/fighterscores", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			matchHandler.FighterScores(w, r, r.PathValue("id"))
		})))
	}

	if deps.InventoryService != nil {
		inventoryHandler := inventoryhandlers.NewHandler(deps.InventoryService)

		mux.Handle("GET /api/inventory/balance/particles", authMiddleware(http.HandlerFunc(inventoryHandler.BalanceParticles)))
		mux.Handle("GET /api/inventory/balance/token/common", authMiddleware(http.HandlerFunc(inventoryHandler.BalanceTokenCommon)))
		mux.Handle("GET /api/inventory/balance/token/rare", authMiddleware(http.HandlerFunc(inventoryHandler.BalanceTokenRare)))
		mux.Handle("GET /api/inventory/balance/token/fabled", authMiddleware(http.HandlerFunc(inventoryHandler.BalanceTokenFabled)))
		mux.Handle("GET /api/inventory/balance/token/mythic", authMiddleware(http.HandlerFunc(inventoryHandler.BalanceTokenMythic)))

		mux.Handle("GET /api/equipment/{id}", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			inventoryHandler.GetEquipment(w, r, r.PathValue("id"))
		})))
		mux.Handle("POST /api/equipment/enhance/cost", authMiddleware(http.HandlerFunc(inventoryHandler.EnhancementCost)))
		mux.Handle("POST /api/equipment/enhance", authMiddleware(http.HandlerFunc(inventoryHandler.Enhance)))
		mux.Handle("POST /api/equipment/salvage", authMiddleware(http.HandlerFunc(inventoryHandler.Salvage)))
		mux.Handle("POST /api/equipment/salvage/inventory", authMiddleware(http.HandlerFunc(inventoryHandler.SalvageInventory)))
		mux.Handle("POST /api/equipment/inventory", authMiddleware(http.HandlerFunc(inventoryHandler.InventoryPage)))
		mux.Handle("GET /api/equipment/fighter/{id}", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			inventoryHandler.ListByFighter(w, r, r.PathValue("id"))
		})))
		mux.Handle("POST /api/equipment/{id}/favorite", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			inventoryHandler.SetFavorite(w, r, r.PathValue("id"), true)
		})))
		mux.Handle("DELETE /api/equipment/{id}/favorite", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			inventoryHandler.SetFavorite(w, r, r.PathValue("id"), false)
		})))
	}

	if deps.LeagueService != nil {
		leagueHandler := leaguehandlers.NewHandler(deps.LeagueService)

		mux.Handle("GET /api/league", authMiddleware(http.HandlerFunc(leagueHandler.List)))
		mux.Handle("GET /api/league/{id}", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			leagueHandler.Get(w, r, r.PathValue("id"))
		})))
		mux.Handle("GET /api/league/{id}/winner", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			leagueHandler.LastWinner(w, r, r.PathValue("id"))
		})))
		mux.Handle("POST /api/league/subscribe", authMiddleware(http.HandlerFunc(leagueHandler.Subscribe)))
		mux.Handle("POST /api/league/unsubscribe", authMiddleware(http.HandlerFunc(leagueHandler.Unsubscribe)))
		mux.Handle("GET /api/league/{id}/subscriptions", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			leagueHandler.Subscriptions(w, r, r.PathValue("id"))
		})))
		mux.Handle("GET /api/league/{id}/subscriptions/user", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			leagueHandler.SubscriptionsForUser(w, r, r.PathValue("id"))
		})))
		mux.Handle("POST /api/league/{id}/matches", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			leagueHandler.Matches(w, r, r.PathValue("id"))
		})))
		mux.Handle("POST /api/league/{id}/highscores", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			leagueHandler.Highscores(w, r, r.PathValue("id"))
		})))
		if deps.LeagueJob != nil {
			mux.Handle("POST /api/league/{id}/run", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				runLeagueJob(w, r, r.PathValue("id"), deps.LeagueJob)
			})))
		}
	}

	if deps.RewardService != nil {
		rewardHandler := rewardhandlers.NewHandler(deps.RewardService)

		mux.Handle("GET /api/reward", authMiddleware(http.HandlerFunc(rewardHandler.List)))
		mux.Handle("POST /api/reward/claim", authMiddleware(http.HandlerFunc(rewardHandler.Claim)))
		mux.Handle("POST /api/reward/claim/all", authMiddleware(http.HandlerFunc(rewardHandler.ClaimAll)))
	}

	if deps.SeasonService != nil {
		seasonHandler := seasonhandlers.NewHandler(deps.SeasonService)
		mux.Handle("POST /api/season/summary", authMiddleware(http.HandlerFunc(seasonHandler.Summary)))
	}

	if deps.MatchHub != nil {
		mux.Handle("GET /ws/match", deps.MatchHub)
	}

	return middleware.WithCORS(mux)
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
