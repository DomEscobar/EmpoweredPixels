package jobs

import (
	"context"
	"log"
	"time"

	matchesusecase "empoweredpixels/internal/usecase/matches"
)

type LobbyCleanupJob struct {
	matchService     *matchesusecase.Service
	olderThanMinutes int
	interval         time.Duration
	stop             chan struct{}
}

func NewLobbyCleanupJob(matchService *matchesusecase.Service, olderThanMinutes int, interval time.Duration) *LobbyCleanupJob {
	return &LobbyCleanupJob{
		matchService:     matchService,
		olderThanMinutes: olderThanMinutes,
		interval:         interval,
		stop:             make(chan struct{}),
	}
}

func (j *LobbyCleanupJob) Start() {
	go func() {
		ticker := time.NewTicker(j.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				j.Run()
			case <-j.stop:
				return
			}
		}
	}()
}

func (j *LobbyCleanupJob) Stop() {
	close(j.stop)
}

func (j *LobbyCleanupJob) Run() {
	ctx := context.Background()
	cancelled, err := j.matchService.CleanupStaleLobbies(ctx, j.olderThanMinutes)
	if err != nil {
		log.Printf("lobby cleanup error: %v", err)
		return
	}
	if cancelled > 0 {
		log.Printf("lobby cleanup: cancelled %d stale lobbies", cancelled)
	}
}
