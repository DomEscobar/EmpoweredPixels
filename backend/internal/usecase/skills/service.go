package skills

import (
	"context"
	"errors"

	"empoweredpixels/internal/domain/skills"
)

var (
	ErrSkillNotFound       = errors.New("skill not found")
	ErrNoSkillPoints       = errors.New("no skill points available")
	ErrSkillMaxRank        = errors.New("skill already at maximum rank")
	ErrPrerequisitesNotMet = errors.New("prerequisites not met for this skill")
	ErrInvalidLoadout      = errors.New("invalid loadout configuration")
	ErrLoadoutTooLarge     = errors.New("loadout exceeds maximum active skills")
	ErrSkillNotAllocated   = errors.New("skill not allocated")
	ErrNotActiveSkill      = errors.New("only active skills can be in loadout")
	ErrInvalidData         = errors.New("invalid data provided for validation")
)

type SkillRepository interface {
	GetFighterSkills(ctx context.Context, fighterID string) (*skills.FighterSkills, error)
	AllocateSkillPoint(ctx context.Context, fighterID string, skillID string, points int) error
	SetLoadout(ctx context.Context, fighterID string, loadout []string) error
	ResetSkills(ctx context.Context, fighterID string) error
	UpdateUltimateCharge(ctx context.Context, fighterID string, charge int) error
}

type FighterRepository interface {
	GetFighterLevel(ctx context.Context, fighterID string) (int, error)
}

type Service struct {
	skillRepo           SkillRepository
	fighterRepo         FighterRepository
	allocationValidator Validator
	loadoutValidator    Validator
}

func NewService(skillRepo SkillRepository, fighterRepo FighterRepository) *Service {
	return &Service{
		skillRepo:           skillRepo,
		fighterRepo:         fighterRepo,
		allocationValidator: &AllocationValidator{},
		loadoutValidator:    &LoadoutValidator{},
	}
}

// GetSkillTree returns all skills organized by branch
func (s *Service) GetSkillTree(ctx context.Context) *SkillTreeResponse {
	return &SkillTreeResponse{
		Offense: skills.GetSkillsByBranch(skills.Offense),
		Defense: skills.GetSkillsByBranch(skills.Defense),
		Utility: skills.GetSkillsByBranch(skills.Utility),
		Ultimates: skills.UltimateSkills,
	}
}

// GetFighterSkillState returns the complete skill state for a fighter
func (s *Service) GetFighterSkillState(ctx context.Context, fighterID string) (*FighterSkillStateResponse, error) {
	fighterSkills, err := s.skillRepo.GetFighterSkills(ctx, fighterID)
	if err != nil {
		return nil, err
	}

	level, err := s.skillRepo.GetFighterLevel(ctx, fighterID)
	if err != nil {
		return nil, err
	}

	// Calculate available points
	maxPoints := level * skills.PointsPerLevel
	allocatedPoints := 0
	for _, p := range fighterSkills.AllocatedPoints {
		allocatedPoints += p
	}
	availablePoints := maxPoints - allocatedPoints

	// Build allocated skills list
	var allocated []AllocatedSkillResponse
	for skillID, points := range fighterSkills.AllocatedPoints {
		if skill, found := skills.GetSkillByID(skillID); found {
			allocated = append(allocated, AllocatedSkillResponse{
				SkillID:     skillID,
				Name:        skill.Name,
				Branch:      skill.Branch.String(),
				Rank:        points,
				EffectValue: skills.CalculateSkillEffect(skill, points),
			})
		}
	}

	// Build loadout
	var loadout []LoadoutSkillResponse
	for _, skillID := range fighterSkills.Loadout {
		if skill, found := skills.GetSkillByID(skillID); found {
			rank := fighterSkills.AllocatedPoints[skillID]
			loadout = append(loadout, LoadoutSkillResponse{
				SkillID:     skillID,
				Name:        skill.Name,
				ManaCost:    skill.ManaCost,
				Cooldown:    skill.Cooldown,
				EffectValue: skills.CalculateSkillEffect(skill, rank),
				IconURL:     skill.IconURL,
			})
		}
	}

	// Get progress for each branch
	progress := skills.GetSkillProgress(fighterSkills.AllocatedPoints)
	var branchProgress []BranchProgressResponse
	for _, p := range progress {
		branchProgress = append(branchProgress, BranchProgressResponse{
			Branch:          p.Branch.String(),
			PointsAllocated: p.PointsAllocated,
			SkillsUnlocked:  p.SkillsUnlocked,
			MaxTierUnlocked: p.MaxTierUnlocked,
		})
	}

	return &FighterSkillStateResponse{
		FighterID:       fighterID,
		Level:           level,
		AvailablePoints: availablePoints,
		AllocatedPoints: allocatedPoints,
		Skills:          allocated,
		Loadout:         loadout,
		BranchProgress:  branchProgress,
		UltimateCharge:  fighterSkills.UltimateCharge,
	}, nil
}

// AllocateSkillPoint allocates a skill point to a skill
func (s *Service) AllocateSkillPoint(ctx context.Context, fighterID string, skillID string) error {
	res := s.allocationValidator.Validate(ctx, s, fighterID, skillID)
	if !res.IsValid {
		return res.Error
	}

	// Get current skills
	fighterSkills, err := s.skillRepo.GetFighterSkills(ctx, fighterID)
	if err != nil {
		return err
	}

	// Add point
	newPoints := fighterSkills.AllocatedPoints[skillID] + 1
	return s.skillRepo.AllocateSkillPoint(ctx, fighterID, skillID, newPoints)
}

// SetLoadout sets the active skills loadout
func (s *Service) SetLoadout(ctx context.Context, fighterID string, loadout []string) error {
	res := s.loadoutValidator.Validate(ctx, s, fighterID, loadout)
	if !res.IsValid {
		return res.Error
	}

	return s.skillRepo.SetLoadout(ctx, fighterID, loadout)
}

// ResetSkills resets all skill allocations (costs gold)
func (s *Service) ResetSkills(ctx context.Context, fighterID string) error {
	return s.skillRepo.ResetSkills(ctx, fighterID)
}

// GetResetCost returns the gold cost to reset skills
func (s *Service) GetResetCost(ctx context.Context, fighterID string) (int, error) {
	level, err := s.skillRepo.GetFighterLevel(ctx, fighterID)
	if err != nil {
		return 0, err
	}
	return skills.ResetCost(level), nil
}

// AddUltimateCharge adds charge to ultimate from combat
func (s *Service) AddUltimateCharge(ctx context.Context, fighterID string, damageDealt, damageTaken int) error {
	fighterSkills, err := s.skillRepo.GetFighterSkills(ctx, fighterID)
	if err != nil {
		return err
	}

	newCharge := skills.AddUltimateCharge(fighterSkills.UltimateCharge, damageDealt, damageTaken)
	return s.skillRepo.UpdateUltimateCharge(ctx, fighterID, newCharge)
}

// CanUseUltimate checks if ultimate is ready
func (s *Service) CanUseUltimate(ctx context.Context, fighterID string) (bool, error) {
	fighterSkills, err := s.skillRepo.GetFighterSkills(ctx, fighterID)
	if err != nil {
		return false, err
	}
	return fighterSkills.UltimateCharge >= skills.UltimateThreshold, nil
}

// UseUltimate resets ultimate charge after use
func (s *Service) UseUltimate(ctx context.Context, fighterID string) error {
	return s.skillRepo.UpdateUltimateCharge(ctx, fighterID, 0)
}

// Response types
type SkillTreeResponse struct {
	Offense   []skills.Skill `json:"offense"`
	Defense   []skills.Skill `json:"defense"`
	Utility   []skills.Skill `json:"utility"`
	Ultimates []skills.Skill `json:"ultimates"`
}

type AllocatedSkillResponse struct {
	SkillID     string `json:"skillId"`
	Name        string `json:"name"`
	Branch      string `json:"branch"`
	Rank        int    `json:"rank"`
	EffectValue int    `json:"effectValue"`
}

type LoadoutSkillResponse struct {
	SkillID     string `json:"skillId"`
	Name        string `json:"name"`
	ManaCost    int    `json:"manaCost"`
	Cooldown    int    `json:"cooldown"`
	EffectValue int    `json:"effectValue"`
	IconURL     string `json:"iconUrl"`
}

type BranchProgressResponse struct {
	Branch          string `json:"branch"`
	PointsAllocated int    `json:"pointsAllocated"`
	SkillsUnlocked  int    `json:"skillsUnlocked"`
	MaxTierUnlocked int    `json:"maxTierUnlocked"`
}

type FighterSkillStateResponse struct {
	FighterID       string                   `json:"fighterId"`
	Level           int                      `json:"level"`
	AvailablePoints int                      `json:"availablePoints"`
	AllocatedPoints int                      `json:"allocatedPoints"`
	Skills          []AllocatedSkillResponse `json:"skills"`
	Loadout         []LoadoutSkillResponse   `json:"loadout"`
	BranchProgress  []BranchProgressResponse `json:"branchProgress"`
	UltimateCharge  int                      `json:"ultimateCharge"`
}