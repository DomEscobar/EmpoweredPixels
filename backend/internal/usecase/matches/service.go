package matches

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"empoweredpixels/internal/domain/combat"
	"empoweredpixels/internal/domain/inventory"
	"empoweredpixels/internal/domain/matches"
	"empoweredpixels/internal/infra/engine"
	inventoryusecase "empoweredpixels/internal/usecase/inventory"
	"empoweredpixels/internal/usecase/rewards"
	rosterusecase "empoweredpixels/internal/usecase/roster"

	"github.com/google/uuid"
)

var (
	ErrInvalidMatch      = errors.New("invalid match")
	ErrInvalidFighter    = errors.New("invalid fighter")
	ErrInvalidTeam       = errors.New("invalid team")
	ErrInvalidTeamPass   = errors.New("invalid team password")
	ErrRegistration      = errors.New("invalid match registration")
	ErrMatchLimit        = errors.New("match fighter limit exceeded")
	ErrMatchNotLobby     = errors.New("match is not in lobby state")
	ErrNotEnoughFighters = errors.New("not enough fighters")
)

type Hub interface {
	Broadcast(matchID string, payload any)
}

type Service struct {
	matches       MatchRepository
	teams         TeamRepository
	registrations RegistrationRepository
	results       ResultRepository
	scores        ScoreRepository
	fighters      FighterRepository
	inventory     inventoryusecase.Service
	rewards       *rewards.Service
	roster        *rosterusecase.Service
	engine        *engine.Client
	hub           Hub
	now           func() time.Time
}

func NewService(
	matches MatchRepository,
	teams TeamRepository,
	registrations RegistrationRepository,
	results ResultRepository,
	scores ScoreRepository,
	fighters FighterRepository,
	inventory inventoryusecase.Service,
	rewards *rewards.Service,
	roster *rosterusecase.Service,
	engineClient *engine.Client,
	hub Hub,
	now func() time.Time,
) *Service {
	if now == nil {
		now = time.Now
	}

	return &Service{
		matches:       matches,
		teams:         teams,
		registrations: registrations,
		results:       results,
		scores:        scores,
		fighters:      fighters,
		inventory:     inventory,
		rewards:       rewards,
		roster:        roster,
		engine:        engineClient,
		hub:           hub,
		now:           now,
	}
}

type MatchOptions struct {
	IsPrivate          bool `json:"isPrivate"`
	MaxFightersPerUser *int `json:"maxFightersPerUser"`
	MaxPowerlevel      *int `json:"maxPowerlevel"`
	BotCount           *int `json:"botCount"`
	BotPowerlevel      *int `json:"botPowerlevel"`
	AutoStart          bool `json:"autoStart"`
}

func (s *Service) DefaultOptions() MatchOptions {
	return MatchOptions{
		IsPrivate:          false,
		MaxFightersPerUser: nil,
		MaxPowerlevel:      nil,
	}
}

func (s *Service) BattleFieldSizes() []string {
	return []string{}
}

func (s *Service) CreateMatch(ctx context.Context, userID int64, options MatchOptions) (*matches.Match, error) {
	data, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	match := &matches.Match{
		ID:            uuid.NewString(),
		CreatorUserID: &userID,
		Created:       s.now(),
		Status:        matches.MatchStatusLobby,
		Options:       data,
	}

	if err := s.matches.Create(ctx, match); err != nil {
		return nil, err
	}

	// Auto-join creator's first fighter if available
	fighters, err := s.fighters.ListByUser(ctx, userID)
	if err == nil && len(fighters) > 0 {
		_ = s.Join(ctx, userID, match.ID, fighters[0].ID)
	}

	return match, nil
}

func (s *Service) GetCurrentMatch(ctx context.Context, userID int64) (*matches.Match, error) {
	return s.matches.GetCurrentMatch(ctx, userID)
}

func (s *Service) CreateTeam(ctx context.Context, matchID string, password *string) (*matches.MatchTeam, error) {
	match, err := s.matches.GetByID(ctx, matchID)
	if err != nil {
		return nil, err
	}
	if match == nil || match.Status != matches.MatchStatusLobby {
		return nil, ErrInvalidMatch
	}

	team := &matches.MatchTeam{
		ID:       uuid.NewString(),
		MatchID:  matchID,
		Password: password,
	}
	if err := s.teams.Create(ctx, team); err != nil {
		return nil, err
	}

	return team, nil
}

func (s *Service) GetMatch(ctx context.Context, id string) (*matches.Match, error) {
	return s.matches.GetByID(ctx, id)
}

func (s *Service) GetTeams(ctx context.Context, matchID string) ([]matches.MatchTeam, error) {
	return s.teams.ListByMatch(ctx, matchID)
}

func (s *Service) GetRegistrations(ctx context.Context, matchID string) ([]matches.MatchRegistration, error) {
	return s.registrations.ListByMatch(ctx, matchID)
}

func (s *Service) BrowseMatches(ctx context.Context, page int, pageSize int) ([]matches.Match, error) {
	return s.BrowseByStatus(ctx, matches.MatchStatusLobby, page, pageSize)
}

func (s *Service) BrowseByStatus(ctx context.Context, status string, page int, pageSize int) ([]matches.Match, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	if status == "" {
		status = matches.MatchStatusLobby
	}
	return s.matches.ListByStatus(ctx, status, pageSize, offset)
}

func (s *Service) Join(ctx context.Context, userID int64, matchID string, fighterID string) error {
	match, err := s.matches.GetByID(ctx, matchID)
	if err != nil {
		return err
	}
	if match == nil || match.Status != matches.MatchStatusLobby {
		return ErrInvalidMatch
	}

	fighter, err := s.fighters.GetByUserAndID(ctx, userID, fighterID)
	if err != nil {
		return err
	}
	if fighter == nil {
		return ErrInvalidFighter
	}

	var options MatchOptions
	_ = json.Unmarshal(match.Options, &options)
	if options.MaxFightersPerUser != nil {
		count, err := s.registrations.CountByMatchAndUser(ctx, matchID, userID)
		if err != nil {
			return err
		}
		if count >= *options.MaxFightersPerUser {
			return ErrMatchLimit
		}
	}

	registration := &matches.MatchRegistration{
		MatchID:   matchID,
		FighterID: fighterID,
		Date:      s.now(),
	}
	if err := s.registrations.Upsert(ctx, registration); err != nil {
		return err
	}
	if s.hub != nil {
		s.hub.Broadcast(matchID, map[string]any{"type": "lobbyUpdate", "matchId": matchID})
	}

	if options.AutoStart {
		s.tryAutoStart(matchID, options)
	}

	return nil
}

func (s *Service) JoinTeam(ctx context.Context, userID int64, matchID string, teamID string, fighterID string, password *string) error {
	match, err := s.matches.GetByID(ctx, matchID)
	if err != nil {
		return err
	}
	if match == nil || match.Status != matches.MatchStatusLobby {
		return ErrInvalidMatch
	}

	team, err := s.teams.GetByID(ctx, teamID)
	if err != nil {
		return err
	}
	if team == nil {
		return ErrInvalidTeam
	}
	if team.Password != nil {
		if password == nil || *password != *team.Password {
			return ErrInvalidTeamPass
		}
	}

	fighter, err := s.fighters.GetByUserAndID(ctx, userID, fighterID)
	if err != nil {
		return err
	}
	if fighter == nil {
		return ErrInvalidFighter
	}

	registration, err := s.registrations.GetByMatchAndFighter(ctx, matchID, fighterID)
	if err != nil {
		return err
	}
	if registration == nil {
		registration = &matches.MatchRegistration{
			MatchID:   matchID,
			FighterID: fighterID,
			Date:      s.now(),
			TeamID:    &teamID,
		}
		return s.registrations.Upsert(ctx, registration)
	}

	registration.TeamID = &teamID
	return s.registrations.Upsert(ctx, registration)
}

func (s *Service) Leave(ctx context.Context, userID int64, matchID string, fighterID string) error {
	match, err := s.matches.GetByID(ctx, matchID)
	if err != nil {
		return err
	}
	if match == nil || match.Status != matches.MatchStatusLobby {
		return ErrInvalidMatch
	}
	fighter, err := s.fighters.GetByUserAndID(ctx, userID, fighterID)
	if err != nil {
		return err
	}
	if fighter == nil {
		return ErrInvalidFighter
	}
	if err := s.registrations.Delete(ctx, matchID, fighterID); err != nil {
		return err
	}
	if s.hub != nil {
		s.hub.Broadcast(matchID, map[string]any{"type": "lobbyUpdate", "matchId": matchID})
	}
	return nil
}

func (s *Service) LeaveTeam(ctx context.Context, userID int64, matchID string, fighterID string, teamID string) error {
	match, err := s.matches.GetByID(ctx, matchID)
	if err != nil {
		return err
	}
	if match == nil || match.Status != matches.MatchStatusLobby {
		return ErrInvalidMatch
	}
	fighter, err := s.fighters.GetByUserAndID(ctx, userID, fighterID)
	if err != nil {
		return err
	}
	if fighter == nil {
		return ErrInvalidFighter
	}

	registration, err := s.registrations.GetByMatchAndFighter(ctx, matchID, fighterID)
	if err != nil {
		return err
	}
	if registration == nil {
		return ErrRegistration
	}
	if registration.TeamID == nil || *registration.TeamID != teamID {
		return nil
	}
	registration.TeamID = nil
	return s.registrations.Upsert(ctx, registration)
}

func (s *Service) RoundTicks(ctx context.Context, matchID string) ([]byte, error) {
	result, err := s.results.GetByMatch(ctx, matchID)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, ErrInvalidMatch
	}
	return result.RoundTicks, nil
}

func (s *Service) FighterScores(ctx context.Context, matchID string) ([]matches.MatchScoreFighter, error) {
	return s.scores.ListByMatch(ctx, matchID)
}

func (s *Service) ExecuteMatch(ctx context.Context, matchID string) error {
	match, err := s.matches.GetByID(ctx, matchID)
	if err != nil {
		return err
	}
	if match == nil {
		return ErrInvalidMatch
	}
	if match.Status != matches.MatchStatusLobby {
		return ErrMatchNotLobby
	}

	fighters, err := s.fighters.ListByMatch(ctx, matchID)
	if err != nil {
		return err
	}

	var options MatchOptions
	_ = json.Unmarshal(match.Options, &options)

	totalFighters := len(fighters)
	if options.BotCount != nil {
		totalFighters += *options.BotCount
	}

	if totalFighters < 2 {
		return ErrNotEnoughFighters
	}

	now := s.now()
	match.Status = matches.MatchStatusRunning
	match.Started = &now
	if err := s.matches.Update(ctx, match); err != nil {
		return err
	}

	if s.hub != nil {
		s.hub.Broadcast(matchID, map[string]any{"type": "matchStatus", "status": matches.MatchStatusRunning, "matchId": matchID})
	}

	// Load equipment for all participants
	fighterEquipment := make(map[string][]inventory.Equipment)
	if s.inventory != nil {
		for _, f := range fighters {
			items, err := s.inventory.ListByFighter(ctx, f.UserID, f.ID)
			if err == nil {
				fighterEquipment[f.ID] = items
			}
		}
	}

	simulator := NewSimulator()
	result, err := simulator.Run(matchID, fighters, fighterEquipment, options)
	if err != nil {
		match.Status = matches.MatchStatusLobby
		match.Started = nil
		_ = s.matches.Update(ctx, match)
		return err
	}

	roundTicksJson, _ := json.Marshal(result.RoundTicks)
	matchResult := &matches.MatchResult{
		ID:         uuid.NewString(),
		MatchID:    matchID,
		RoundTicks: roundTicksJson,
	}
	if err := s.results.Upsert(ctx, matchResult); err != nil {
		return err
	}

	scoresMapping := make(map[string]combat.FighterScore)
	for _, score := range result.Scores {
		scoresMapping[score.FighterID] = score
	}

	scores := make([]matches.MatchScoreFighter, 0, len(result.Scores))
	for _, score := range result.Scores {
		scores = append(scores, matches.MatchScoreFighter{
			MatchID:      matchID,
			FighterID:    score.FighterID,
			TotalKills:   score.Kills,
			TotalDeaths:  score.Deaths,
			TotalAssists: score.Assists,
		})
	}
	if len(scores) > 0 {
		if err := s.scores.Upsert(ctx, scores); err != nil {
			return err
		}
	}

	completedAt := s.now()
	match.Status = matches.MatchStatusCompleted
	match.CompletedAt = &completedAt
	if err := s.matches.Update(ctx, match); err != nil {
		return err
	}

	// Award rewards and experience to all participants
	if s.rewards != nil {
		rewardedUsers := make(map[int64]bool)

		var winnerID string
		maxKills := -1
		for _, score := range result.Scores {
			if score.Deaths == 0 && score.Kills > maxKills {
				maxKills = score.Kills
				winnerID = score.FighterID
			}
		}

		// Calculate Bot difficulty bonus
		botBonusExp := 0
		if options.BotCount != nil && options.BotPowerlevel != nil {
			// +1 EXP for every 5 powerlevels of bots, scaled by bot count
			botBonusExp = (*options.BotPowerlevel / 5) * (*options.BotCount / 2)
		}

		for _, f := range fighters {
			// Award Loot (per User)
			if !rewardedUsers[f.UserID] {
				pool := "match_participation"
				if f.ID == winnerID {
					pool = "match_win"
				}

				if _, err := s.rewards.IssueReward(ctx, f.UserID, pool); err != nil {
					// log error but don't fail the match execution
				}
				rewardedUsers[f.UserID] = true
			}

			// Award Experience (per Fighter)
			if s.roster != nil {
				score, ok := scoresMapping[f.ID]
				expAmount := 10 + botBonusExp // Base EXP + difficulty bonus
				if ok {
					expAmount += score.Kills * 5
					if f.ID == winnerID {
						expAmount += 20 // Winner bonus EXP
					}
				}

				currentExp, err := s.roster.GetExperience(ctx, f.ID)
				if err == nil {
					currentExp.Experience += expAmount
					_ = s.roster.UpdateExperience(ctx, currentExp)
				}
			}
		}
	}

	if s.hub != nil {
		s.hub.Broadcast(matchID, map[string]any{"type": "matchEnded", "matchId": matchID, "status": matches.MatchStatusCompleted})
	}

	return nil
}


func (s *Service) tryAutoStart(matchID string, options MatchOptions) {
	go func() {
		ctx := context.Background()

		fighters, err := s.fighters.ListByMatch(ctx, matchID)
		if err != nil {
			return
		}

		totalFighters := len(fighters)
		if options.BotCount != nil {
			totalFighters += *options.BotCount
		}

		if totalFighters >= 2 {
			_ = s.ExecuteMatch(ctx, matchID)
		}
	}()
}

func (s *Service) CancelMatch(ctx context.Context, matchID string) error {
	match, err := s.matches.GetByID(ctx, matchID)
	if err != nil {
		return err
	}
	if match == nil {
		return ErrInvalidMatch
	}
	if match.Status != matches.MatchStatusLobby {
		return ErrMatchNotLobby
	}

	now := s.now()
	match.Status = matches.MatchStatusCancelled
	match.CancelledAt = &now
	if err := s.matches.Update(ctx, match); err != nil {
		return err
	}

	if s.hub != nil {
		s.hub.Broadcast(matchID, map[string]any{"type": "matchCancelled", "matchId": matchID})
	}

	return nil
}

func (s *Service) CleanupStaleLobbies(ctx context.Context, olderThanMinutes int) (int, error) {
	stale, err := s.matches.ListStaleLobbies(ctx, olderThanMinutes)
	if err != nil {
		return 0, err
	}

	cancelled := 0
	for _, match := range stale {
		if err := s.CancelMatch(ctx, match.ID); err == nil {
			cancelled++
		}
	}

	return cancelled, nil
}
