package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"empoweredpixels/internal/usecase/leagues"
	matchesusecase "empoweredpixels/internal/usecase/matches"
)

var (
	ErrLeagueNotFound  = errors.New("league not found")
	ErrNoSubscriptions = errors.New("league has no subscriptions")
)

type LeagueJob struct {
	matchService    *matchesusecase.Service
	leagueRepo      leagues.LeagueRepository
	subRepo         leagues.SubscriptionRepository
	leagueMatchRepo leagues.LeagueMatchRepository
	fighterRepo     leagues.FighterRepository
	interval        time.Duration
	mu              sync.Mutex // Prevents concurrent league runs
}

func NewLeagueJob(
	matchService *matchesusecase.Service,
	leagueRepo leagues.LeagueRepository,
	subRepo leagues.SubscriptionRepository,
	leagueMatchRepo leagues.LeagueMatchRepository,
	fighterRepo leagues.FighterRepository,
	interval time.Duration,
) *LeagueJob {
	return &LeagueJob{
		matchService:    matchService,
		leagueRepo:      leagueRepo,
		subRepo:         subRepo,
		leagueMatchRepo: leagueMatchRepo,
		fighterRepo:     fighterRepo,
		interval:        interval,
	}
}

func (j *LeagueJob) Start() {
	if j.interval <= 0 {
		return
	}

	go func() {
		ticker := time.NewTicker(j.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				j.runAllActiveLeagues()
			}
		}
	}()
}

func (j *LeagueJob) runAllActiveLeagues() {
	// Prevent concurrent executions if previous run is still in progress
	j.mu.Lock()
	defer j.mu.Unlock()

	ctx := context.Background()
	leagues, err := j.leagueRepo.List(ctx)
	if err != nil {
		return
	}

	for _, league := range leagues {
		_ = j.RunLeague(ctx, league.ID)
	}
}

func (j *LeagueJob) RunMatch(ctx context.Context, matchID string) error {
	if j.matchService == nil {
		return nil
	}
	return j.matchService.ExecuteMatch(ctx, matchID)
}

func (j *LeagueJob) RunLeague(ctx context.Context, leagueID int) error {
	if j.leagueRepo == nil || j.subRepo == nil || j.leagueMatchRepo == nil || j.fighterRepo == nil || j.matchService == nil {
		return nil
	}

	league, err := j.leagueRepo.GetByID(ctx, leagueID)
	if err != nil || league == nil {
		return ErrLeagueNotFound
	}

	subs, err := j.subRepo.ListByLeague(ctx, leagueID)
	if err != nil {
		return err
	}
	if len(subs) == 0 {
		return ErrNoSubscriptions
	}

	firstFighter, err := j.fighterRepo.GetByID(ctx, subs[0].FighterID)
	if err != nil || firstFighter == nil {
		return err
	}
	creatorUserID := firstFighter.UserID

	options := j.matchService.DefaultOptions()
	if len(league.Options) > 0 {
		var leagueOpts struct {
			BotCount      *int `json:"botCount"`
			BotPowerlevel *int `json:"botPowerlevel"`
		}
		_ = json.Unmarshal(league.Options, &leagueOpts)
		if leagueOpts.BotCount != nil {
			// Validate botCount (1-10)
			if *leagueOpts.BotCount < 1 || *leagueOpts.BotCount > 10 {
				defaultVal := 5
				options.BotCount = &defaultVal
			} else {
				options.BotCount = leagueOpts.BotCount
			}
		}
		if leagueOpts.BotPowerlevel != nil {
			// Validate botPowerlevel (1-100)
			if *leagueOpts.BotPowerlevel < 1 || *leagueOpts.BotPowerlevel > 100 {
				defaultVal := 50
				options.BotPowerlevel = &defaultVal
			} else {
				options.BotPowerlevel = leagueOpts.BotPowerlevel
			}
		}
	}

	match, err := j.matchService.CreateMatch(ctx, creatorUserID, options)
	if err != nil {
		return fmt.Errorf("failed to create match: %w", err)
	}

	// Track successful joins for potential rollback
	var successfulJoins []string
	defer func() {
		// If we fail after creating the match, try to clean up by disassociating the match from the league
		if err != nil && match != nil {
			j.leagueMatchRepo.Delete(ctx, leagueID, match.ID)
			// Note: We don't delete the match itself as it might be used elsewhere
		}
	}()

	for _, sub := range subs {
		fighter, err := j.fighterRepo.GetByID(ctx, sub.FighterID)
		if err != nil || fighter == nil {
			// Continue with other fighters even if one fails
			continue
		}
		if err := j.matchService.Join(ctx, fighter.UserID, match.ID, sub.FighterID); err != nil {
			// Log but continue - a partial league match is still valid
			log.Printf("Failed to add fighter %s to league match: %v", sub.FighterID, err)
			continue
		}
		successfulJoins = append(successfulJoins, sub.FighterID)
	}

	if len(successfulJoins) == 0 {
		return errors.New("no fighters successfully joined the match")
	}

	if err := j.leagueMatchRepo.Create(ctx, leagueID, match.ID); err != nil {
		return fmt.Errorf("failed to record league match: %w", err)
	}

	if err := j.matchService.ExecuteMatch(ctx, match.ID); err != nil {
		return fmt.Errorf("failed to execute match: %w", err)
	}

	updated, _ := j.matchService.GetMatch(ctx, match.ID)
	if updated != nil && updated.Started != nil {
		_ = j.leagueMatchRepo.UpdateStarted(ctx, leagueID, match.ID, updated.Started)
	}

	return nil
}
