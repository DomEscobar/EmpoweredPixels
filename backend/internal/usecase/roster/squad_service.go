package roster

import (
	"context"

	"empoweredpixels/internal/domain/roster"
	"github.com/google/uuid"
)

type SquadService struct {
	repo        SquadRepository
	fighterRepo FighterRepository
}

func NewSquadService(repo SquadRepository, fighterRepo FighterRepository) *SquadService {
	return &SquadService{
		repo:        repo,
		fighterRepo: fighterRepo,
	}
}

func (s *SquadService) SetActiveSquad(ctx context.Context, userID int64, name string, fighterIDs []string) (*roster.Squad, error) {
	if len(fighterIDs) > 3 {
		fighterIDs = fighterIDs[:3]
	}

	err := s.repo.DeactivateAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	squad := &roster.Squad{
		ID:       uuid.NewString(),
		UserID:   userID,
		Name:     name,
		IsActive: true,
	}

	for i, id := range fighterIDs {
		squad.Members = append(squad.Members, roster.Member{
			FighterID: id,
			SlotIndex: i,
		})
	}

	err = s.repo.Create(ctx, squad)
	if err != nil {
		return nil, err
	}

	return s.GetActiveSquad(ctx, userID)
}

func (s *SquadService) GetActiveSquad(ctx context.Context, userID int64) (*roster.Squad, error) {
	squad, err := s.repo.GetActiveByUserID(ctx, userID)
	if err != nil || squad == nil {
		return squad, err
	}

	// Hydrate fighter data for each member in parallel to reduce latency
	type result struct {
		index   int
		fighter *roster.Fighter
		err     error
	}

	resChan := make(chan result, len(squad.Members))

	for i := range squad.Members {
		go func(idx int, fighterID string) {
			fighter, err := s.fighterRepo.GetByID(ctx, fighterID)
			resChan <- result{index: idx, fighter: fighter, err: err}
		}(i, squad.Members[i].FighterID)
	}

	for i := 0; i < len(squad.Members); i++ {
		res := <-resChan
		if res.err == nil && res.fighter != nil {
			squad.Members[res.index].Fighter = res.fighter
		}
	}

	return squad, nil
}
