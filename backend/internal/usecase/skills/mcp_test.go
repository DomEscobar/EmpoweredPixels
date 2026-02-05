package skills_test

import (
	"context"
	"testing"

	"empoweredpixels/internal/domain/skills"
	skillsusecase "empoweredpixels/internal/usecase/skills"
)

// MCP Test validates that external AI agents can interact with the skill system
func TestMCP_SkillQuery(t *testing.T) {
	skillRepo := newMockSkillRepo()
	fighterRepo := newMockFighterRepo()
	fighterRepo.levels["ai_fighter_001"] = 50
	service := skillsusecase.NewService(skillRepo, fighterRepo)

	ctx := context.Background()

	// Simulate AI agent querying skill tree
	tree := service.GetSkillTree(ctx)

	// AI agent should be able to see all branches
	if len(tree.Offense) == 0 {
		t.Error("AI agent cannot see offense skills")
	}
	if len(tree.Defense) == 0 {
		t.Error("AI agent cannot see defense skills")
	}
	if len(tree.Utility) == 0 {
		t.Error("AI agent cannot see utility skills")
	}
	if len(tree.Ultimates) == 0 {
		t.Error("AI agent cannot see ultimate skills")
	}

	// AI agent should be able to identify skill types
	for _, skill := range tree.Offense {
		if skill.Branch != skills.Offense {
			t.Errorf("Skill %s has wrong branch", skill.ID)
		}
	}
}

func TestMCP_SkillAllocation(t *testing.T) {
	skillRepo := newMockSkillRepo()
	fighterRepo := newMockFighterRepo()
	fighterRepo.levels["ai_fighter_001"] = 30
	service := skillsusecase.NewService(skillRepo, fighterRepo)

	ctx := context.Background()

	// AI agent allocates skills strategically
	allocationPlan := []string{
		"skl_power_strike", // Tier 1 Offense
		"skl_power_strike", // Rank 2
		"skl_bleed",        // Tier 1 Offense (passive)
		"skl_block",        // Tier 1 Defense
		"skl_heal",         // Tier 1 Defense
	}

	for _, skillID := range allocationPlan {
		err := service.AllocateSkillPoint(ctx, "ai_fighter_001", skillID)
		if err != nil {
			t.Errorf("AI agent failed to allocate %s: %v", skillID, err)
		}
	}

	// Verify AI agent's allocations
	state, err := service.GetFighterSkillState(ctx, "ai_fighter_001")
	if err != nil {
		t.Fatalf("Failed to get skill state: %v", err)
	}

	if state.AllocatedPoints != 5 {
		t.Errorf("AI agent expected 5 allocated points, got %d", state.AllocatedPoints)
	}
}

func TestMCP_LoadoutManagement(t *testing.T) {
	skillRepo := newMockSkillRepo()
	fighterRepo := newMockFighterRepo()
	fighterRepo.levels["ai_fighter_001"] = 20
	service := skillsusecase.NewService(skillRepo, fighterRepo)

	ctx := context.Background()

	// AI agent allocates active skills
	service.AllocateSkillPoint(ctx, "ai_fighter_001", "skl_power_strike")
	service.AllocateSkillPoint(ctx, "ai_fighter_001", "skl_block")

	// AI agent sets loadout
	err := service.SetLoadout(ctx, "ai_fighter_001", []string{"skl_power_strike", "skl_block"})
	if err != nil {
		t.Errorf("AI agent failed to set loadout: %v", err)
	}

	// AI agent queries current loadout
	state, err := service.GetFighterSkillState(ctx, "ai_fighter_001")
	if err != nil {
		t.Fatalf("Failed to get skill state: %v", err)
	}

	if len(state.Loadout) != 2 {
		t.Errorf("AI agent expected 2 loadout skills, got %d", len(state.Loadout))
	}
}

func TestMCP_UltimateCharge(t *testing.T) {
	skillRepo := newMockSkillRepo()
	fighterRepo := newMockFighterRepo()
	service := skillsusecase.NewService(skillRepo, fighterRepo)

	ctx := context.Background()

	// Simulate combat adding ultimate charge
	// AI deals 5000 damage and takes 2500 damage
	err := service.AddUltimateCharge(ctx, "ai_fighter_001", 5000, 2500)
	if err != nil {
		t.Errorf("Failed to add ultimate charge: %v", err)
	}

	// Check if ultimate is ready (5000/50 + 2500/25 = 100 + 100 = 200, capped at 100)
	ready, err := service.CanUseUltimate(ctx, "ai_fighter_001")
	if err != nil {
		t.Fatalf("Failed to check ultimate status: %v", err)
	}

	if !ready {
		t.Error("AI agent's ultimate should be ready after combat")
	}

	// AI agent uses ultimate
	err = service.UseUltimate(ctx, "ai_fighter_001")
	if err != nil {
		t.Errorf("AI agent failed to use ultimate: %v", err)
	}

	// Verify ultimate charge reset
	ready, _ = service.CanUseUltimate(ctx, "ai_fighter_001")
	if ready {
		t.Error("Ultimate should not be ready after use")
	}
}

func TestMCP_SkillTreeTraversal(t *testing.T) {
	// Test that AI agents can traverse the skill tree programmatically
	branches := []skills.SkillBranch{skills.Offense, skills.Defense, skills.Utility}

	for _, branch := range branches {
		branchSkills := skills.GetSkillsByBranch(branch)

		// Group by tier
		skillsByTier := make(map[int][]skills.Skill)
		for _, skill := range branchSkills {
			skillsByTier[skill.Tier] = append(skillsByTier[skill.Tier], skill)
		}

		// Verify tier progression (tier 1 should have most skills in MVP)
		if len(skillsByTier[1]) == 0 {
			t.Errorf("No tier 1 skills in %s branch", branch.String())
		}

		// AI agent can identify active vs passive skills
		for _, skill := range branchSkills {
			if skill.Type != skills.Active && skill.Type != skills.Passive && skill.Type != skills.Ultimate {
				t.Errorf("Skill %s has invalid type", skill.ID)
			}
		}
	}
}

func TestMCP_ErrorHandling(t *testing.T) {
	skillRepo := newMockSkillRepo()
	fighterRepo := newMockFighterRepo()
	fighterRepo.levels["ai_fighter_001"] = 1 // Low level
	service := skillsusecase.NewService(skillRepo, fighterRepo)

	ctx := context.Background()

	// AI agent tries to allocate without enough points
	err := service.AllocateSkillPoint(ctx, "ai_fighter_001", "skl_power_strike")
	if err != nil {
		t.Errorf("First allocation should succeed: %v", err)
	}

	// Second allocation should fail (no points left at level 1)
	err = service.AllocateSkillPoint(ctx, "ai_fighter_001", "skl_bleed")
	if err == nil {
		t.Error("AI agent should get error when allocating without points")
	}

	// AI agent tries to allocate non-existent skill
	err = service.AllocateSkillPoint(ctx, "ai_fighter_001", "invalid_skill")
	if err != skillsusecase.ErrSkillNotFound {
		t.Errorf("Expected ErrSkillNotFound, got: %v", err)
	}
}