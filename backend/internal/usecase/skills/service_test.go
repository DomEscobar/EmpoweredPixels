package skills_test

import (
	"context"
	"testing"

	"empoweredpixels/internal/domain/skills"
	skillsusecase "empoweredpixels/internal/usecase/skills"
)

// Mock repositories
type mockSkillRepo struct {
	fighterSkills map[string]*skills.FighterSkills
}

func newMockSkillRepo() *mockSkillRepo {
	return &mockSkillRepo{
		fighterSkills: make(map[string]*skills.FighterSkills),
	}
}

func (m *mockSkillRepo) GetFighterSkills(ctx context.Context, fighterID string) (*skills.FighterSkills, error) {
	if fs, ok := m.fighterSkills[fighterID]; ok {
		return fs, nil
	}
	return &skills.FighterSkills{
		FighterID:       fighterID,
		AllocatedPoints: make(map[string]int),
		Loadout:         []string{},
		UltimateCharge:  0,
	}, nil
}

func (m *mockSkillRepo) AllocateSkillPoint(ctx context.Context, fighterID string, skillID string, points int) error {
	if _, ok := m.fighterSkills[fighterID]; !ok {
		m.fighterSkills[fighterID] = &skills.FighterSkills{
			FighterID:       fighterID,
			AllocatedPoints: make(map[string]int),
			Loadout:         []string{},
		}
	}
	m.fighterSkills[fighterID].AllocatedPoints[skillID] = points
	return nil
}

func (m *mockSkillRepo) SetLoadout(ctx context.Context, fighterID string, loadout []string) error {
	if _, ok := m.fighterSkills[fighterID]; !ok {
		m.fighterSkills[fighterID] = &skills.FighterSkills{
			FighterID:       fighterID,
			AllocatedPoints: make(map[string]int),
		}
	}
	m.fighterSkills[fighterID].Loadout = loadout
	return nil
}

func (m *mockSkillRepo) ResetSkills(ctx context.Context, fighterID string) error {
	m.fighterSkills[fighterID] = &skills.FighterSkills{
		FighterID:       fighterID,
		AllocatedPoints: make(map[string]int),
		Loadout:         []string{},
		UltimateCharge:  0,
	}
	return nil
}

func (m *mockSkillRepo) UpdateUltimateCharge(ctx context.Context, fighterID string, charge int) error {
	if _, ok := m.fighterSkills[fighterID]; !ok {
		m.fighterSkills[fighterID] = &skills.FighterSkills{
			FighterID:       fighterID,
			AllocatedPoints: make(map[string]int),
		}
	}
	m.fighterSkills[fighterID].UltimateCharge = charge
	return nil
}

type mockFighterRepo struct {
	levels map[string]int
}

func newMockFighterRepo() *mockFighterRepo {
	return &mockFighterRepo{
		levels: make(map[string]int),
	}
}

func (m *mockFighterRepo) GetFighterLevel(ctx context.Context, fighterID string) (int, error) {
	if level, ok := m.levels[fighterID]; ok {
		return level, nil
	}
	return 1, nil
}

func TestService_GetSkillTree(t *testing.T) {
	skillRepo := newMockSkillRepo()
	fighterRepo := newMockFighterRepo()
	service := skillsusecase.NewService(skillRepo, fighterRepo)

	tree := service.GetSkillTree(context.Background())

	if len(tree.Offense) == 0 {
		t.Error("Expected offense skills in tree")
	}
	if len(tree.Defense) == 0 {
		t.Error("Expected defense skills in tree")
	}
	if len(tree.Utility) == 0 {
		t.Error("Expected utility skills in tree")
	}
}

func TestService_AllocateSkillPoint(t *testing.T) {
	skillRepo := newMockSkillRepo()
	fighterRepo := newMockFighterRepo()
	fighterRepo.levels["fighter1"] = 10
	service := skillsusecase.NewService(skillRepo, fighterRepo)

	ctx := context.Background()

	// Test allocating tier 1 skill
	err := service.AllocateSkillPoint(ctx, "fighter1", "skl_power_strike")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test allocating without prerequisites
	err = service.AllocateSkillPoint(ctx, "fighter1", "skl_whirlwind")
	if err != skillsusecase.ErrPrerequisitesNotMet {
		t.Errorf("Expected ErrPrerequisitesNotMet, got: %v", err)
	}

	// Test allocating non-existent skill
	err = service.AllocateSkillPoint(ctx, "fighter1", "nonexistent")
	if err != skillsusecase.ErrSkillNotFound {
		t.Errorf("Expected ErrSkillNotFound, got: %v", err)
	}
}

func TestService_SetLoadout(t *testing.T) {
	skillRepo := newMockSkillRepo()
	fighterRepo := newMockFighterRepo()
	fighterRepo.levels["fighter1"] = 10
	service := skillsusecase.NewService(skillRepo, fighterRepo)

	ctx := context.Background()

	// Allocate some skills first
	service.AllocateSkillPoint(ctx, "fighter1", "skl_power_strike")
	service.AllocateSkillPoint(ctx, "fighter1", "skl_block")

	// Test valid loadout
	err := service.SetLoadout(ctx, "fighter1", []string{"skl_power_strike", "skl_block"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test loadout too large
	err = service.SetLoadout(ctx, "fighter1", []string{"skl_power_strike", "skl_block", "skl_heal", "skl_haste"})
	if err != skillsusecase.ErrLoadoutTooLarge {
		t.Errorf("Expected ErrLoadoutTooLarge, got: %v", err)
	}

	// Test loadout with unallocated skill
	err = service.SetLoadout(ctx, "fighter1", []string{"skl_execute"})
	if err != skillsusecase.ErrSkillNotAllocated {
		t.Errorf("Expected ErrSkillNotAllocated, got: %v", err)
	}
}

func TestService_ResetSkills(t *testing.T) {
	skillRepo := newMockSkillRepo()
	fighterRepo := newMockFighterRepo()
	fighterRepo.levels["fighter1"] = 10
	service := skillsusecase.NewService(skillRepo, fighterRepo)

	ctx := context.Background()

	// Allocate skills and set loadout
	service.AllocateSkillPoint(ctx, "fighter1", "skl_power_strike")
	service.SetLoadout(ctx, "fighter1", []string{"skl_power_strike"})

	// Reset skills
	err := service.ResetSkills(ctx, "fighter1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify reset
	state, err := service.GetFighterSkillState(ctx, "fighter1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(state.Skills) != 0 {
		t.Error("Expected skills to be reset")
	}
	if len(state.Loadout) != 0 {
		t.Error("Expected loadout to be reset")
	}
}

func TestService_GetResetCost(t *testing.T) {
	skillRepo := newMockSkillRepo()
	fighterRepo := newMockFighterRepo()
	fighterRepo.levels["fighter1"] = 50
	service := skillsusecase.NewService(skillRepo, fighterRepo)

	ctx := context.Background()

	cost, err := service.GetResetCost(ctx, "fighter1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedCost := 2500 // 50 * 50
	if cost != expectedCost {
		t.Errorf("Expected cost %d, got %d", expectedCost, cost)
	}
}

func TestService_AddUltimateCharge(t *testing.T) {
	skillRepo := newMockSkillRepo()
	fighterRepo := newMockFighterRepo()
	service := skillsusecase.NewService(skillRepo, fighterRepo)

	ctx := context.Background()

	// Add charge from combat
	err := service.AddUltimateCharge(ctx, "fighter1", 100, 50)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check if ultimate is ready
	ready, err := service.CanUseUltimate(ctx, "fighter1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if ready {
		t.Error("Ultimate should not be ready yet")
	}
}

func TestService_GetFighterSkillState(t *testing.T) {
	skillRepo := newMockSkillRepo()
	fighterRepo := newMockFighterRepo()
	fighterRepo.levels["fighter1"] = 20
	service := skillsusecase.NewService(skillRepo, fighterRepo)

	ctx := context.Background()

	// Allocate some skills
	service.AllocateSkillPoint(ctx, "fighter1", "skl_power_strike")
	service.AllocateSkillPoint(ctx, "fighter1", "skl_power_strike") // Rank 2
	service.SetLoadout(ctx, "fighter1", []string{"skl_power_strike"})

	state, err := service.GetFighterSkillState(ctx, "fighter1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if state.FighterID != "fighter1" {
		t.Errorf("Expected fighter1, got %s", state.FighterID)
	}
	if state.Level != 20 {
		t.Errorf("Expected level 20, got %d", state.Level)
	}
	if len(state.Skills) != 1 {
		t.Errorf("Expected 1 skill, got %d", len(state.Skills))
	}
	if len(state.Loadout) != 1 {
		t.Errorf("Expected 1 loadout slot, got %d", len(state.Loadout))
	}
}