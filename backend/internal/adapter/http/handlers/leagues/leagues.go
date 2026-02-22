package leagues

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"empoweredpixels/internal/adapter/http/middleware"
	"empoweredpixels/internal/adapter/http/responses"
	"empoweredpixels/internal/domain/leagues"
	leaguesusecase "empoweredpixels/internal/usecase/leagues"
)

// Service defines the contract for league business logic needed by the HTTP handler.
type Service interface {
	List(ctx context.Context) ([]leagues.League, error)
	Get(ctx context.Context, id int) (*leagues.League, error)
	Subscribe(ctx context.Context, userID int64, leagueID int, fighterID string) error
	Unsubscribe(ctx context.Context, userID int64, leagueID int, fighterID string) error
	Subscriptions(ctx context.Context, leagueID int) ([]leagues.LeagueSubscription, error)
	SubscriptionsForUser(ctx context.Context, leagueID int, userID int64) ([]leagues.LeagueSubscription, error)
	Matches(ctx context.Context, leagueID int, page int, pageSize int) ([]leagues.LeagueMatch, int, error)
	GetLastWinner(ctx context.Context, leagueID int) (*leagues.LeagueWinner, error)
	GetHighScores(ctx context.Context, leagueID int, lastMatches int) ([]leagues.LeagueHighscore, error)
	CreateLeague(ctx context.Context, name string, options []byte, isDeactivated bool) (*leagues.League, error)
	UpdateLeague(ctx context.Context, id int, name string, options []byte, isDeactivated bool) (*leagues.League, error)
	DeleteLeague(ctx context.Context, id int) error
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type createLeagueDto struct {
	Name     string          `json:"name"`
	Options  json.RawMessage `json:"options"`
	IsActive bool            `json:"isActive"`
}

type updateLeagueDto struct {
	Name     *string          `json:"name"`
	Options  *json.RawMessage `json:"options"`
	IsActive *bool            `json:"isActive"`
}

type leagueDto struct {
	ID       int             `json:"id"`
	Name     string          `json:"name"`
	Options  json.RawMessage `json:"options"`
	IsActive bool            `json:"isActive"`
}

type leagueDetailDto struct {
	leagueDto
	Subscriptions []leagueSubscriptionDto `json:"subscriptions"`
}

type leagueLastWinnerDto struct {
	LeagueID int `json:"leagueId"`
}

type leagueSubscriptionDto struct {
	LeagueID  int    `json:"leagueId"`
	FighterID string `json:"fighterId"`
}

type leagueMatchDto struct {
	LeagueID int        `json:"leagueId"`
	MatchID  string     `json:"matchId"`
	Started  *time.Time `json:"started"`
}

type leagueHighscoreDto struct {
	FighterID   string `json:"fighterId"`
	FighterName string `json:"fighterName"`
	Username    string `json:"username"`
	Score       int    `json:"score"`
}

type pagingOptions struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type pageDto[T any] struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalCount int `json:"totalCount"`
	Items      []T `json:"items"`
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	leaguesList, err := h.service.List(r.Context())
	if err != nil {
		log.Printf("league list error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	result := make([]leagueDto, 0, len(leaguesList))
	for _, league := range leaguesList {
		result = append(result, mapLeague(league))
	}

	responses.JSON(w, http.StatusOK, result)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request, id string) {
	leagueID, err := strconv.Atoi(id)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid league")
		return
	}

	league, err := h.service.Get(r.Context(), leagueID)
	if err != nil {
		log.Printf("league get error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}
	if league == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	subs, err := h.service.Subscriptions(r.Context(), leagueID)
	if err != nil {
		log.Printf("league subscriptions error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	detail := leagueDetailDto{
		leagueDto:     mapLeague(*league),
		Subscriptions: mapSubscriptions(subs),
	}

	responses.JSON(w, http.StatusOK, detail)
}

func (h *Handler) LastWinner(w http.ResponseWriter, r *http.Request, id string) {
	leagueID, err := strconv.Atoi(id)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid league")
		return
	}

	winner, err := h.service.GetLastWinner(r.Context(), leagueID)
	if err != nil {
		log.Printf("league last winner error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}
	if winner == nil {
		// No completed matches found
		w.WriteHeader(http.StatusNoContent)
		return
	}

	responses.JSON(w, http.StatusOK, winner)
}

func (h *Handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var payload leagueSubscriptionDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	if err := h.service.Subscribe(r.Context(), userID, payload.LeagueID, payload.FighterID); err != nil {
		switch err {
		case leaguesusecase.ErrInvalidLeague, leaguesusecase.ErrInvalidFighter:
			responses.Error(w, http.StatusBadRequest, err.Error())
			return
		default:
			log.Printf("league subscribe error: %v", err)
			responses.Error(w, http.StatusInternalServerError, "server error")
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var payload leagueSubscriptionDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	if err := h.service.Unsubscribe(r.Context(), userID, payload.LeagueID, payload.FighterID); err != nil {
		switch err {
		case leaguesusecase.ErrInvalidLeague, leaguesusecase.ErrInvalidFighter, leaguesusecase.ErrInvalidSubscription:
			responses.Error(w, http.StatusBadRequest, err.Error())
			return
		default:
			log.Printf("league unsubscribe error: %v", err)
			responses.Error(w, http.StatusInternalServerError, "server error")
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Subscriptions(w http.ResponseWriter, r *http.Request, id string) {
	leagueID, err := strconv.Atoi(id)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid league")
		return
	}

	subs, err := h.service.Subscriptions(r.Context(), leagueID)
	if err != nil {
		log.Printf("league subscriptions by id error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	responses.JSON(w, http.StatusOK, mapSubscriptions(subs))
}

func (h *Handler) SubscriptionsForUser(w http.ResponseWriter, r *http.Request, id string) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	leagueID, err := strconv.Atoi(id)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid league")
		return
	}

	subs, err := h.service.SubscriptionsForUser(r.Context(), leagueID, userID)
	if err != nil {
		log.Printf("league subscriptions for user error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	responses.JSON(w, http.StatusOK, mapSubscriptions(subs))
}

func (h *Handler) Matches(w http.ResponseWriter, r *http.Request, id string) {
	leagueID, err := strconv.Atoi(id)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid league")
		return
	}

	var payload pagingOptions
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		payload = pagingOptions{Page: 1, PageSize: 20}
	}

	matchesList, totalCount, err := h.service.Matches(r.Context(), leagueID, payload.Page, payload.PageSize)
	if err != nil {
		log.Printf("league matches error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	items := make([]leagueMatchDto, 0, len(matchesList))
	for _, match := range matchesList {
		items = append(items, leagueMatchDto{
			LeagueID: match.LeagueID,
			MatchID:  match.MatchID,
			Started:  match.Started,
		})
	}

	responses.JSON(w, http.StatusOK, pageDto[leagueMatchDto]{
		Page:       payload.Page,
		PageSize:   payload.PageSize,
		TotalCount: totalCount,
		Items:      items,
	})
}

func (h *Handler) Highscores(w http.ResponseWriter, r *http.Request, id string) {
	leagueID, err := strconv.Atoi(id)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid league")
		return
	}

	// Parse lastMatches query parameter (default 50)
	lastMatchesStr := r.URL.Query().Get("lastMatches")
	lastMatches := 50
	if lastMatchesStr != "" {
		if lm, err := strconv.Atoi(lastMatchesStr); err == nil && lm > 0 {
			lastMatches = lm
		}
	}

	highscores, err := h.service.GetHighScores(r.Context(), leagueID, lastMatches)
	if err != nil {
		log.Printf("league highscores error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	result := make([]leagueHighscoreDto, 0, len(highscores))
	for _, hs := range highscores {
		result = append(result, leagueHighscoreDto{
			FighterID:   hs.FighterID,
			FighterName: hs.FighterName,
			Username:    hs.Username,
			Score:       hs.Score,
		})
	}

	responses.JSON(w, http.StatusOK, result)
}

func mapLeague(league leagues.League) leagueDto {
	return leagueDto{
		ID:       league.ID,
		Name:     league.Name,
		Options:  league.Options,
		IsActive: !league.IsDeactivated,
	}
}

func (h *Handler) CreateLeague(w http.ResponseWriter, r *http.Request) {
	var payload createLeagueDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	// Convert isActive to isDeactivated (inverse)
	isDeactivated := !payload.IsActive

	// Use options as provided (json.RawMessage is already []byte)
	optionsBytes := payload.Options
	if optionsBytes == nil {
		optionsBytes = []byte("{}")
	}

	league, err := h.service.CreateLeague(r.Context(), payload.Name, optionsBytes, isDeactivated)
	if err != nil {
		log.Printf("league create error: %v", err)
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusCreated, mapLeague(*league))
}

func (h *Handler) UpdateLeague(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid league id")
		return
	}

	var payload updateLeagueDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	// Get current league to merge fields
	current, err := h.service.Get(r.Context(), id)
	if err != nil {
		log.Printf("league get for update error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}
	if current == nil {
		responses.Error(w, http.StatusNotFound, "league not found")
		return
	}

	// Apply updates
	name := current.Name
	if payload.Name != nil {
		if *payload.Name == "" {
			responses.Error(w, http.StatusBadRequest, "name cannot be empty")
			return
		}
		name = *payload.Name
	}

	options := current.Options
	if payload.Options != nil {
		// payload.Options is *json.RawMessage, dereference and handle nil
		opt := *payload.Options
		if opt == nil {
			options = []byte("{}")
		} else {
			options = opt
		}
	}

	isActive := !current.IsDeactivated
	if payload.IsActive != nil {
		isActive = *payload.IsActive
	}
	isDeactivated := !isActive

	league, err := h.service.UpdateLeague(r.Context(), id, name, options, isDeactivated)
	if err != nil {
		log.Printf("league update error: %v", err)
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, mapLeague(*league))
}

func (h *Handler) DeleteLeague(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid league id")
		return
	}

	if err := h.service.DeleteLeague(r.Context(), id); err != nil {
		log.Printf("league delete error: %v", err)
		switch err {
		case leaguesusecase.ErrInvalidLeague:
			responses.Error(w, http.StatusNotFound, err.Error())
		default:
			responses.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func mapSubscriptions(subs []leagues.LeagueSubscription) []leagueSubscriptionDto {
	result := make([]leagueSubscriptionDto, 0, len(subs))
	for _, sub := range subs {
		result = append(result, leagueSubscriptionDto{
			LeagueID:  sub.LeagueID,
			FighterID: sub.FighterID,
		})
	}
	return result
}
