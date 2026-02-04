package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"empoweredpixels/internal/infra/db/repositories"
	matchesusecase "empoweredpixels/internal/usecase/matches"
)

var (
	ErrLeagueNotFound   = errors.New("league not found")
	ErrNoSubscriptions  = errors.New("league has no subscriptions")
)

type LeagueJob struct {
	matchService    *matchesusecase.Service
	leagueRepo      *repositories.LeagueRepository
	subRepo         *repositories.LeagueSubscriptionRepository
	leagueMatchRepo *repositories.LeagueMatchRepository
	fighterRepo     *repositories.FighterRepository
	interval        time.Duration
}

func NewLeagueJob(
	matchService *matchesusecase.Service,
	leagueRepo *repositories.LeagueRepository,
	subRepo *repositories.LeagueSubscriptionRepository,
	leagueMatchRepo *repositories.LeagueMatchRepository,
	fighterRepo *repositories.FighterRepository,
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
			options.BotCount = leagueOpts.BotCount
		}
		if leagueOpts.BotPowerlevel != nil {
			options.BotPowerlevel = leagueOpts.BotPowerlevel
		}
	}

	match, err := j.matchService.CreateMatch(ctx, creatorUserID, options)
	if err != nil {
		return err
	}

	for _, sub := range subs {
		fighter, err := j.fighterRepo.GetByID(ctx, sub.FighterID)
		if err != nil || fighter == nil {
			continue
		}
		_ = j.matchService.Join(ctx, fighter.UserID, match.ID, sub.FighterID)
	}

	if err := j.leagueMatchRepo.Create(ctx, leagueID, match.ID); err != nil {
		return err
	}

	if err := j.matchService.ExecuteMatch(ctx, match.ID); err != nil {
		return err
	}

	updated, _ := j.matchService.GetMatch(ctx, match.ID)
	if updated != nil && updated.Started != nil {
		_ = j.leagueMatchRepo.UpdateStarted(ctx, leagueID, match.ID, updated.Started)
	}

	return nil
}
