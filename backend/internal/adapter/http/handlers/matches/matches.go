package matches

import (
	"encoding/json"
	"log"
	"net/http"

	"empoweredpixels/internal/adapter/http/middleware"
	"empoweredpixels/internal/adapter/http/responses"
	"empoweredpixels/internal/domain/matches"
	matchesusecase "empoweredpixels/internal/usecase/matches"
)

type Handler struct {
	service *matchesusecase.Service
}

func NewHandler(service *matchesusecase.Service) *Handler {
	return &Handler{service: service}
}

type matchOptionsDto struct {
	IsPrivate          bool     `json:"isPrivate"`
	MaxPowerlevel      *int     `json:"maxPowerlevel"`
	ActionsPerRound    int      `json:"actionsPerRound"`
	MaxFightersPerUser *int     `json:"maxFightersPerUser"`
	BotCount           *int     `json:"botCount"`
	BotPowerlevel      *int     `json:"botPowerlevel"`
	AutoStart          bool     `json:"autoStart"`
	Features           []string `json:"features"`
	Battlefield        string   `json:"battlefield"`
	Bounds             string   `json:"bounds"`
	PositionGenerator  string   `json:"positionGenerator"`
	MoveOrder          string   `json:"moveOrder"`
	WinCondition       string   `json:"winCondition"`
	StaleCondition     string   `json:"staleCondition"`
}

type matchDto struct {
	ID            string                 `json:"id"`
	CreatorUserID *int64                 `json:"creatorUserId"`
	Created       string                 `json:"created"`
	Started       *string                `json:"started"`
	CompletedAt   *string                `json:"completedAt"`
	CancelledAt   *string                `json:"cancelledAt"`
	Status        string                 `json:"status"`
	Ended         bool                   `json:"ended"`
	Registrations []matchRegistrationDto `json:"registrations"`
	Options       matchOptionsDto        `json:"options"`
}

type matchTeamDto struct {
	ID          string `json:"id"`
	MatchID     string `json:"matchId"`
	HasPassword bool   `json:"hasPassword"`
}

type matchTeamOperationDto struct {
	ID        string  `json:"id"`
	MatchID   string  `json:"matchId"`
	FighterID *string `json:"fighterId"`
	Password  *string `json:"password"`
}

type matchRegistrationDto struct {
	MatchID   string  `json:"matchId"`
	FighterID string  `json:"fighterId"`
	TeamID    *string `json:"teamId"`
}

type pagingOptions struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Status   string `json:"status"`
}

type pageDto[T any] struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalCount int `json:"totalCount"`
	Items      []T `json:"items"`
}

type matchScoreFighterDto struct {
	MatchID      string `json:"matchId"`
	FighterID    string `json:"fighterId"`
	TotalKills   int    `json:"totalKills"`
	TotalDeaths  int    `json:"totalDeaths"`
	TotalAssists int    `json:"totalAssists"`
}

func (h *Handler) GetDefaultOptions(w http.ResponseWriter, r *http.Request) {
	options := h.service.DefaultOptions()
	responses.JSON(w, http.StatusOK, matchOptionsDto{
		IsPrivate:          options.IsPrivate,
		MaxPowerlevel:      options.MaxPowerlevel,
		MaxFightersPerUser: options.MaxFightersPerUser,
	})
}

func (h *Handler) GetBattlefieldSizes(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, h.service.BattleFieldSizes())
}

func (h *Handler) CreateMatch(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var payload matchOptionsDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	match, err := h.service.CreateMatch(r.Context(), userID, matchesusecase.MatchOptions{
		IsPrivate:          payload.IsPrivate,
		MaxFightersPerUser: payload.MaxFightersPerUser,
		MaxPowerlevel:      payload.MaxPowerlevel,
		BotCount:           payload.BotCount,
		BotPowerlevel:      payload.BotPowerlevel,
		AutoStart:          payload.AutoStart,
	})
	if err != nil {
		log.Printf("match create error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	registrations, _ := h.service.GetRegistrations(r.Context(), match.ID)
	responses.JSON(w, http.StatusOK, mapMatch(match, registrations))
}

func (h *Handler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var payload matchTeamOperationDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	team, err := h.service.CreateTeam(r.Context(), payload.MatchID, payload.Password)
	if err != nil {
		if err == matchesusecase.ErrInvalidMatch {
			responses.Error(w, http.StatusBadRequest, "invalid match")
			return
		}
		log.Printf("match create team error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	responses.JSON(w, http.StatusOK, matchTeamDto{
		ID:          team.ID,
		MatchID:     team.MatchID,
		HasPassword: team.Password != nil,
	})
}

func (h *Handler) GetMatch(w http.ResponseWriter, r *http.Request, id string) {
	match, err := h.service.GetMatch(r.Context(), id)
	if err != nil {
		log.Printf("match get error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}
	if match == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	registrations, _ := h.service.GetRegistrations(r.Context(), match.ID)
	responses.JSON(w, http.StatusOK, mapMatch(match, registrations))
}

func (h *Handler) GetCurrentMatch(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	match, err := h.service.GetCurrentMatch(r.Context(), userID)
	if err != nil {
		log.Printf("match get current error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}
	if match == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	registrations, _ := h.service.GetRegistrations(r.Context(), match.ID)
	responses.JSON(w, http.StatusOK, mapMatch(match, registrations))
}

func (h *Handler) GetTeams(w http.ResponseWriter, r *http.Request, id string) {
	teams, err := h.service.GetTeams(r.Context(), id)
	if err != nil {
		log.Printf("match get teams error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	result := make([]matchTeamDto, 0, len(teams))
	for _, team := range teams {
		result = append(result, matchTeamDto{
			ID:          team.ID,
			MatchID:     team.MatchID,
			HasPassword: team.Password != nil,
		})
	}

	responses.JSON(w, http.StatusOK, result)
}

func (h *Handler) Browse(w http.ResponseWriter, r *http.Request) {
	var payload pagingOptions
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		payload = pagingOptions{Page: 1, PageSize: 20}
	}

	var matchesList []matches.Match
	var err error
	if payload.Status != "" {
		matchesList, err = h.service.BrowseByStatus(r.Context(), payload.Status, payload.Page, payload.PageSize)
	} else {
		matchesList, err = h.service.BrowseMatches(r.Context(), payload.Page, payload.PageSize)
	}
	if err != nil {
		log.Printf("match browse error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	items := make([]matchDto, 0, len(matchesList))
	for _, match := range matchesList {
		registrations, _ := h.service.GetRegistrations(r.Context(), match.ID)
		item := mapMatch(&match, registrations)
		items = append(items, item)
	}

	responses.JSON(w, http.StatusOK, pageDto[matchDto]{
		Page:       payload.Page,
		PageSize:   payload.PageSize,
		TotalCount: len(items),
		Items:      items,
	})
}

func (h *Handler) Join(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var payload matchRegistrationDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	err := h.service.Join(r.Context(), userID, payload.MatchID, payload.FighterID)
	if err != nil {
		switch err {
		case matchesusecase.ErrInvalidMatch, matchesusecase.ErrInvalidFighter, matchesusecase.ErrMatchLimit:
			responses.Error(w, http.StatusBadRequest, err.Error())
			return
		default:
			log.Printf("match join error: %v", err)
			responses.Error(w, http.StatusInternalServerError, "server error")
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) JoinTeam(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var payload matchTeamOperationDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.FighterID == nil {
		responses.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	err := h.service.JoinTeam(r.Context(), userID, payload.MatchID, payload.ID, *payload.FighterID, payload.Password)
	if err != nil {
		switch err {
		case matchesusecase.ErrInvalidMatch, matchesusecase.ErrInvalidFighter, matchesusecase.ErrInvalidTeam, matchesusecase.ErrInvalidTeamPass:
			responses.Error(w, http.StatusBadRequest, err.Error())
			return
		default:
			log.Printf("match join team error: %v", err)
			responses.Error(w, http.StatusInternalServerError, "server error")
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Leave(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var payload matchRegistrationDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	err := h.service.Leave(r.Context(), userID, payload.MatchID, payload.FighterID)
	if err != nil {
		switch err {
		case matchesusecase.ErrInvalidMatch, matchesusecase.ErrInvalidFighter:
			responses.Error(w, http.StatusBadRequest, err.Error())
			return
		default:
			log.Printf("match leave error: %v", err)
			responses.Error(w, http.StatusInternalServerError, "server error")
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) LeaveTeam(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var payload matchTeamOperationDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.FighterID == nil {
		responses.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	err := h.service.LeaveTeam(r.Context(), userID, payload.MatchID, *payload.FighterID, payload.ID)
	if err != nil {
		switch err {
		case matchesusecase.ErrInvalidFighter, matchesusecase.ErrRegistration:
			responses.Error(w, http.StatusBadRequest, err.Error())
			return
		default:
			log.Printf("match leave team error: %v", err)
			responses.Error(w, http.StatusInternalServerError, "server error")
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Start(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.ExecuteMatch(r.Context(), id); err != nil {
		switch err {
		case matchesusecase.ErrInvalidMatch, matchesusecase.ErrMatchNotLobby, matchesusecase.ErrNotEnoughFighters:
			log.Printf("ExecuteMatch error: %v", err)
			responses.Error(w, http.StatusBadRequest, err.Error())
			return
		default:
			log.Printf("ExecuteMatch error: %v", err)
			responses.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) RoundTicks(w http.ResponseWriter, r *http.Request, id string) {
	data, err := h.service.RoundTicks(r.Context(), id)
	if err != nil {
		if err == matchesusecase.ErrInvalidMatch {
			responses.Error(w, http.StatusBadRequest, "invalid match result")
			return
		}
		log.Printf("match round ticks error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (h *Handler) FighterScores(w http.ResponseWriter, r *http.Request, id string) {
	scores, err := h.service.FighterScores(r.Context(), id)
	if err != nil {
		log.Printf("match fighter scores error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	result := make([]matchScoreFighterDto, 0, len(scores))
	for _, score := range scores {
		result = append(result, matchScoreFighterDto{
			MatchID:      score.MatchID,
			FighterID:    score.FighterID,
			TotalKills:   score.TotalKills,
			TotalDeaths:  score.TotalDeaths,
			TotalAssists: score.TotalAssists,
		})
	}
	responses.JSON(w, http.StatusOK, result)
}

func mapMatch(match *matches.Match, registrations []matches.MatchRegistration) matchDto {
	var options matchOptionsDto
	_ = json.Unmarshal(match.Options, &options)

	var started, completedAt, cancelledAt *string
	if match.Started != nil {
		v := match.Started.Format(timeLayout)
		started = &v
	}
	if match.CompletedAt != nil {
		v := match.CompletedAt.Format(timeLayout)
		completedAt = &v
	}
	if match.CancelledAt != nil {
		v := match.CancelledAt.Format(timeLayout)
		cancelledAt = &v
	}
	status := match.Status
	if status == "" {
		status = matches.MatchStatusLobby
	}

	regs := make([]matchRegistrationDto, 0, len(registrations))
	for _, r := range registrations {
		regs = append(regs, matchRegistrationDto{
			MatchID:   r.MatchID,
			FighterID: r.FighterID,
			TeamID:    r.TeamID,
		})
	}

	return matchDto{
		ID:            match.ID,
		CreatorUserID: match.CreatorUserID,
		Created:       match.Created.Format(timeLayout),
		Started:       started,
		CompletedAt:   completedAt,
		CancelledAt:   cancelledAt,
		Status:        status,
		Ended:         status == matches.MatchStatusCompleted,
		Registrations: regs,
		Options:       options,
	}
}

const timeLayout = "2006-01-02T15:04:05Z07:00"

// QuickJoin handles quick joining an open lobby
func (h *Handler) QuickJoin(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		FighterID string `json:"fighterId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	match, err := h.service.QuickJoin(r.Context(), int64(userID), req.FighterID)
	if err != nil {
		responses.Error(w, http.StatusNotFound, err.Error())
		return
	}

	registrations, _ := h.service.GetRegistrations(r.Context(), match.ID)
	responses.JSON(w, http.StatusOK, mapMatch(match, registrations))
}

// GetOnlinePlayers returns the count of online players
func (h *Handler) GetOnlinePlayers(w http.ResponseWriter, r *http.Request) {
	count, err := h.service.GetOnlinePlayersCount(r.Context())
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	responses.JSON(w, http.StatusOK, map[string]any{
		"onlinePlayers": count,
	})
}
