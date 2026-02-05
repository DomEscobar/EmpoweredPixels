package skillhandlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	skillhandlers "empoweredpixels/internal/adapter/http/handlers/skills"
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

// Mock context with userID
func withUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, "userID", userID)
}

func TestE2E_SkillTree(t *testing.T) {
	skillRepo := newMockSkillRepo()
	fighterRepo := newMockFighterRepo()
	service := skillsusecase.NewService(skillRepo, fighterRepo)
	handler := skillhandlers.NewHandler(service)

	req := httptest.NewRequest("GET", "/api/skills/tree", nil)
	req = req.WithContext(withUserID(req.Context(), 1))
	w := httptest.NewRecorder()

	handler.GetSkillTree(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if _, ok := response["offense"]; !ok {
		t.Error("Expected offense in response")
	}
	if _, ok := response["defense"]; !ok {
		t.Error("Expected defense in response")
	}
	if _, ok := response["utility"]; !ok {
		t.Error("Expected utility in response")
	}
}

func TestE2E_AllocateSkill(t *testing.T) {
	skillRepo := newMockSkillRepo()
	fighterRepo := newMockFighterRepo()
	fighterRepo.levels["fighter1"] = 10
	service := skillsusecase.NewService(skillRepo, fighterRepo)
	handler := skillhandlers.NewHandler(service)

	body := map[string]string{
		"fighterId": "fighter1",
		"skillId":   "skl_power_strike",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/skills/allocate", bytes.NewReader(jsonBody))
	req = req.WithContext(withUserID(req.Context(), 1))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.AllocateSkill(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestE2E_SetLoadout(t *testing.T) {
	skillRepo := newMockSkillRepo()
	fighterRepo := newMockFighterRepo()
	fighterRepo.levels["fighter1"] = 10
	service := skillsusecase.NewService(skillRepo, fighterRepo)
	handler := skillhandlers.NewHandler(service)

	ctx := withUserID(context.Background(), 1)

	// Allocate skill first
	service.AllocateSkillPoint(ctx, "fighter1", "skl_power_strike")

	// Set loadout
	body := map[string]interface{}{
		"fighterId": "fighter1",
		"loadout":   []string{"skl_power_strike"},
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/skills/loadout", bytes.NewReader(jsonBody))
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.SetLoadout(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestE2E_ResetSkills(t *testing.T) {
	skillRepo := newMockSkillRepo()
	fighterRepo := newMockFighterRepo()
	fighterRepo.levels["fighter1"] = 10
	service := skillsusecase.NewService(skillRepo, fighterRepo)
	handler := skillhandlers.NewHandler(service)

	ctx := withUserID(context.Background(), 1)

	// Allocate skill first
	service.AllocateSkillPoint(ctx, "fighter1", "skl_power_strike")

	body := map[string]string{
		"fighterId": "fighter1",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/skills/reset", bytes.NewReader(jsonBody))
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ResetSkills(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || !success {
		t.Error("Expected success to be true")
	}
}

func TestE2E_FullSkillFlow(t *testing.T) {
	// Complete E2E flow: Get tree -> Allocate skills -> Set loadout -> Reset
	skillRepo := newMockSkillRepo()
	fighterRepo := newMockFighterRepo()
	fighterRepo.levels["fighter1"] = 20
	service := skillsusecase.NewService(skillRepo, fighterRepo)
	handler := skillhandlers.NewHandler(service)
	ctx := withUserID(context.Background(), 1)

	// Step 1: Get skill tree
	req := httptest.NewRequest("GET", "/api/skills/tree", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	handler.GetSkillTree(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Step 1 failed: Expected status 200, got %d", w.Code)
	}

	// Step 2: Allocate tier 1 skills
	skills := []string{"skl_power_strike", "skl_bleed", "skl_block", "skl_heal"}
	for _, skillID := range skills {
		body := map[string]string{
			"fighterId": "fighter1",
			"skillId":   skillID,
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/skills/allocate", bytes.NewReader(jsonBody))
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.AllocateSkill(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Step 2 failed for %s: Expected status 200, got %d: %s", skillID, w.Code, w.Body.String())
		}
	}

	// Step 3: Set loadout
	loadoutBody := map[string]interface{}{
		"fighterId": "fighter1",
		"loadout":   []string{"skl_power_strike", "skl_block"},
	}
	jsonBody, _ := json.Marshal(loadoutBody)
	req = httptest.NewRequest("POST", "/api/skills/loadout", bytes.NewReader(jsonBody))
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handler.SetLoadout(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Step 3 failed: Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	// Step 4: Get fighter skills state
	req = httptest.NewRequest("GET", "/api/skills/fighter/fighter1", nil)
	req = req.WithContext(ctx)
	w = httptest.NewRecorder()
	handler.GetFighterSkills(w, req, "fighter1")

	if w.Code != http.StatusOK {
		t.Fatalf("Step 4 failed: Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var state map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &state); err != nil {
		t.Fatalf("Failed to parse state: %v", err)
	}

	if allocated, ok := state["allocatedPoints"].(float64); !ok || allocated != 4 {
		t.Errorf("Expected 4 allocated points, got %v", state["allocatedPoints"])
	}

	loadout, ok := state["loadout"].([]interface{})
	if !ok || len(loadout) != 2 {
		t.Errorf("Expected 2 loadout skills, got %v", state["loadout"])
	}

	// Step 5: Reset skills
	resetBody := map[string]string{"fighterId": "fighter1"}
	jsonBody, _ = json.Marshal(resetBody)
	req = httptest.NewRequest("POST", "/api/skills/reset", bytes.NewReader(jsonBody))
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handler.ResetSkills(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Step 5 failed: Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}