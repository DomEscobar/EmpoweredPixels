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

	// Hydrate fighter data for each member
	for i := range squad.Members {
		fighter, err := s.fighterRepo.GetByID(ctx, squad.Members[i].FighterID)
		if err == nil && fighter != nil {
			squad.Members[i].Fighter = fighter
		}
	}

	return squad, nil
}
