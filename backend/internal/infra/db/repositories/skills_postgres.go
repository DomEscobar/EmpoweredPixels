package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"empoweredpixels/internal/domain/skills"
)

var ErrFighterSkillsNotFound = errors.New("fighter skills not found")

type SkillsPostgres struct {
	db *sql.DB
}

func NewSkillsPostgres(db *sql.DB) *SkillsPostgres {
	return &SkillsPostgres{db: db}
}

func (s *SkillsPostgres) GetFighterSkills(ctx context.Context, fighterID string) (*skills.FighterSkills, error) {
	query := `
		SELECT allocated_points, loadout, ultimate_charge
		FROM fighter_skills
		WHERE fighter_id = $1
	`
	var allocatedJSON []byte
	var loadoutJSON []byte
	var ultimateCharge int

	err := s.db.QueryRowContext(ctx, query, fighterID).Scan(&allocatedJSON, &loadoutJSON, &ultimateCharge)
	if err == sql.ErrNoRows {
		// Return empty skills for new fighters
		return &skills.FighterSkills{
			FighterID:       fighterID,
			AllocatedPoints: make(map[string]int),
			Loadout:         []string{},
			UltimateCharge:  0,
		}, nil
	}
	if err != nil {
		return nil, err
	}

	allocated := make(map[string]int)
	if len(allocatedJSON) > 0 {
		if err := json.Unmarshal(allocatedJSON, &allocated); err != nil {
			return nil, err
		}
	}

	var loadout []string
	if len(loadoutJSON) > 0 {
		if err := json.Unmarshal(loadoutJSON, &loadout); err != nil {
			return nil, err
		}
	}

	return &skills.FighterSkills{
		FighterID:       fighterID,
		AllocatedPoints: allocated,
		Loadout:         loadout,
		UltimateCharge:  ultimateCharge,
	}, nil
}

func (s *SkillsPostgres) AllocateSkillPoint(ctx context.Context, fighterID string, skillID string, points int) error {
	// First try to update existing record
	allocatedJSON, _ := json.Marshal(map[string]int{skillID: points})

	query := `
		INSERT INTO fighter_skills (fighter_id, allocated_points, loadout, ultimate_charge, updated_at)
		VALUES ($1, $2, '[]', 0, NOW())
		ON CONFLICT (fighter_id)
		DO UPDATE SET
			allocated_points = fighter_skills.allocated_points || $2::jsonb,
			updated_at = NOW()
	`
	_, err := s.db.ExecContext(ctx, query, fighterID, allocatedJSON)
	return err
}

func (s *SkillsPostgres) SetLoadout(ctx context.Context, fighterID string, loadout []string) error {
	loadoutJSON, _ := json.Marshal(loadout)

	query := `
		INSERT INTO fighter_skills (fighter_id, allocated_points, loadout, ultimate_charge, updated_at)
		VALUES ($1, '{}', $2, 0, NOW())
		ON CONFLICT (fighter_id)
		DO UPDATE SET
			loadout = $2,
			updated_at = NOW()
	`
	_, err := s.db.ExecContext(ctx, query, fighterID, loadoutJSON)
	return err
}

func (s *SkillsPostgres) ResetSkills(ctx context.Context, fighterID string) error {
	query := `
		UPDATE fighter_skills
		SET allocated_points = '{}',
		    loadout = '[]',
		    ultimate_charge = 0,
		    updated_at = NOW()
		WHERE fighter_id = $1
	`
	_, err := s.db.ExecContext(ctx, query, fighterID)
	return err
}

func (s *SkillsPostgres) UpdateUltimateCharge(ctx context.Context, fighterID string, charge int) error {
	query := `
		INSERT INTO fighter_skills (fighter_id, allocated_points, loadout, ultimate_charge, updated_at)
		VALUES ($1, '{}', '[]', $2, NOW())
		ON CONFLICT (fighter_id)
		DO UPDATE SET
			ultimate_charge = $2,
			updated_at = NOW()
	`
	_, err := s.db.ExecContext(ctx, query, fighterID, charge)
	return err
}