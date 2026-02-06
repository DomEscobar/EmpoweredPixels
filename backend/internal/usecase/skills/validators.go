package skills

import (
	"context"
	"empoweredpixels/internal/domain/skills"
)

type ValidationResult struct {
	IsValid bool
	Error   error
}

type Validator interface {
	Validate(ctx context.Context, s *Service, fighterID string, data interface{}) ValidationResult
}

type AllocationValidator struct{}

func (v *AllocationValidator) Validate(ctx context.Context, s *Service, fighterID string, data interface{}) ValidationResult {
	skillID, ok := data.(string)
	if !ok {
		return ValidationResult{false, ErrInvalidData}
	}

	// Verify skill exists
	_, found := skills.GetSkillByID(skillID)
	if !found {
		return ValidationResult{false, ErrSkillNotFound}
	}

	// Get fighter level
	level, err := s.skillRepo.GetFighterLevel(ctx, fighterID)
	if err != nil {
		return ValidationResult{false, err}
	}

	// Get current skills
	fighterSkills, err := s.skillRepo.GetFighterSkills(ctx, fighterID)
	if err != nil {
		return ValidationResult{false, err}
	}

	// Logic check from skills.CanAllocate
	if !skills.CanAllocate(fighterSkills.AllocatedPoints, skillID, level) {
		currentRank := fighterSkills.AllocatedPoints[skillID]
		if currentRank >= 3 {
			return ValidationResult{false, ErrSkillMaxRank}
		}
		maxPoints := level * skills.PointsPerLevel
		totalAllocated := 0
		for _, p := range fighterSkills.AllocatedPoints {
			totalAllocated += p
		}
		if totalAllocated >= maxPoints {
			return ValidationResult{false, ErrNoSkillPoints}
		}
		return ValidationResult{false, ErrPrerequisitesNotMet}
	}

	return ValidationResult{true, nil}
}

type LoadoutValidator struct{}

func (v *LoadoutValidator) Validate(ctx context.Context, s *Service, fighterID string, data interface{}) ValidationResult {
	loadout, ok := data.([]string)
	if !ok {
		return ValidationResult{false, ErrInvalidData}
	}

	if len(loadout) > skills.MaxActiveSkills {
		return ValidationResult{false, ErrLoadoutTooLarge}
	}

	// Get current skills
	fighterSkills, err := s.skillRepo.GetFighterSkills(ctx, fighterID)
	if err != nil {
		return ValidationResult{false, err}
	}

	// Validate loadout
	if !skills.CanSetLoadout(loadout, fighterSkills.AllocatedPoints) {
		for _, skillID := range loadout {
			if fighterSkills.AllocatedPoints[skillID] == 0 {
				return ValidationResult{false, ErrSkillNotAllocated}
			}
			if skill, found := skills.GetSkillByID(skillID); found && skill.Type != skills.Active {
				return ValidationResult{false, ErrNotActiveSkill}
			}
		}
		return ValidationResult{false, ErrInvalidLoadout}
	}

	return ValidationResult{true, nil}
}
