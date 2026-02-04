package leagues

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"empoweredpixels/internal/adapter/http/middleware"
	"empoweredpixels/internal/adapter/http/responses"
	"empoweredpixels/internal/domain/leagues"
	leaguesusecase "empoweredpixels/internal/usecase/leagues"
)

type Handler struct {
	service *leaguesusecase.Service
}

func NewHandler(service *leaguesusecase.Service) *Handler {
	return &Handler{service: service}
}

type leagueDto struct {
	ID      int             `json:"id"`
	Name    string          `json:"name"`
	Options json.RawMessage `json:"options"`
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
	LeagueID int    `json:"leagueId"`
	MatchID  string `json:"matchId"`
}

type leagueHighscoreDto struct {
	FighterID   string `json:"fighterId"`
	FighterName string `json:"fighterName"`
	Username    string `json:"username"`
	Score       int    `json:"score"`
}

type leagueHighscoreOptionsDto struct {
	LastMatches int `json:"lastMatches"`
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

	responses.JSON(w, http.StatusOK, leagueLastWinnerDto{LeagueID: leagueID})
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

	matchesList, err := h.service.Matches(r.Context(), leagueID, payload.Page, payload.PageSize)
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
		})
	}

	responses.JSON(w, http.StatusOK, pageDto[leagueMatchDto]{
		Page:       payload.Page,
		PageSize:   payload.PageSize,
		TotalCount: len(items),
		Items:      items,
	})
}

func (h *Handler) Highscores(w http.ResponseWriter, r *http.Request, id string) {
	responses.JSON(w, http.StatusOK, []leagueHighscoreDto{})
}

func mapLeague(league leagues.League) leagueDto {
	return leagueDto{
		ID:      league.ID,
		Name:    league.Name,
		Options: league.Options,
	}
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
