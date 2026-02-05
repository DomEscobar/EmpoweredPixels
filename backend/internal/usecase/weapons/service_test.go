package weapons

import (
	"context"
	"errors"
	"testing"
	"time"

	"empoweredpixels/internal/domain/weapons"
)

type mockWeaponRepo struct {
	weapons      map[string]*weapons.UserWeapon
	fighterEquip map[string]*weapons.UserWeapon
	count        int
}

func newMockRepo() *mockWeaponRepo {
	return &mockWeaponRepo{
		weapons:      make(map[string]*weapons.UserWeapon),
		fighterEquip: make(map[string]*weapons.UserWeapon),
	}
}

func (m *mockWeaponRepo) GetByID(ctx context.Context, userID int64, id string) (*weapons.UserWeapon, error) {
	if w, ok := m.weapons[id]; ok && w.UserID == userID {
		return w, nil
	}
	return nil, nil
}

func (m *mockWeaponRepo) ListByUser(ctx context.Context, userID int64, limit int, offset int) ([]weapons.UserWeapon, error) {
	var result []weapons.UserWeapon
	for _, w := range m.weapons {
		if w.UserID == userID {
			result = append(result, *w)
		}
	}
	return result, nil
}

func (m *mockWeaponRepo) ListByUserAll(ctx context.Context, userID int64) ([]weapons.UserWeapon, error) {
	return m.ListByUser(ctx, userID, 100, 0)
}

func (m *mockWeaponRepo) CountByUser(ctx context.Context, userID int64) (int, error) {
	return m.count, nil
}

func (m *mockWeaponRepo) Create(ctx context.Context, userWeapon *weapons.UserWeapon) error {
	m.weapons[userWeapon.ID] = userWeapon
	m.count++
	return nil
}

func (m *mockWeaponRepo) UpdateEnhancement(ctx context.Context, id string, enhancement int) error {
	if w, ok := m.weapons[id]; ok {
		w.Enhancement = enhancement
		return nil
	}
	return errors.New("weapon not found")
}

func (m *mockWeaponRepo) UpdateFighter(ctx context.Context, id string, fighterID *string) error {
	if w, ok := m.weapons[id]; ok {
		// Remove from old fighter if equipped
		if w.FighterID != nil {
			delete(m.fighterEquip, *w.FighterID)
		}
		w.FighterID = fighterID
		w.IsEquipped = fighterID != nil
		if fighterID != nil {
			m.fighterEquip[*fighterID] = w
		}
		return nil
	}
	return errors.New("weapon not found")
}

func (m *mockWeaponRepo) Delete(ctx context.Context, id string) error {
	delete(m.weapons, id)
	return nil
}

func (m *mockWeaponRepo) GetEquippedByFighter(ctx context.Context, fighterID string) (*weapons.UserWeapon, error) {
	if w, ok := m.fighterEquip[fighterID]; ok {
		return w, nil
	}
	return nil, nil
}

func TestService_ListUserWeapons(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	// Add test weapons
	repo.Create(context.Background(), &weapons.UserWeapon{
		ID:       "w1",
		UserID:   1,
		WeaponID: "wpn_sword_iron_002",
		Created:  time.Now(),
	})
	repo.Create(context.Background(), &weapons.UserWeapon{
		ID:       "w2",
		UserID:   1,
		WeaponID: "wpn_bow_short_001",
		Created:  time.Now(),
	})

	weapons, err := svc.ListUserWeapons(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(weapons) != 2 {
		t.Errorf("expected 2 weapons, got %d", len(weapons))
	}
}

func TestService_GetWeaponDetails(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	repo.Create(context.Background(), &weapons.UserWeapon{
		ID:       "w1",
		UserID:   1,
		WeaponID: "wpn_sword_excalibur_006",
		Enhancement: 5,
		Created:  time.Now(),
	})

	uw, def, stats, err := svc.GetWeaponDetails(context.Background(), 1, "w1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if uw == nil {
		t.Fatal("expected user weapon")
	}
	if def == nil {
		t.Fatal("expected weapon definition")
	}
	if def.Name != "Excalibur" {
		t.Errorf("expected Excalibur, got %s", def.Name)
	}
	if stats.Damage <= 0 {
		t.Error("expected positive damage")
	}
}

func TestService_GetWeaponDetails_NotFound(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	_, _, _, err := svc.GetWeaponDetails(context.Background(), 1, "nonexistent")
	if err != ErrWeaponNotFound {
		t.Errorf("expected ErrWeaponNotFound, got %v", err)
	}
}

func TestService_EquipWeapon(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	repo.Create(context.Background(), &weapons.UserWeapon{
		ID:       "w1",
		UserID:   1,
		WeaponID: "wpn_sword_iron_002",
		Created:  time.Now(),
	})

	err := svc.EquipWeapon(context.Background(), 1, "w1", "fighter1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify equipped
	uw, _ := repo.GetByID(context.Background(), 1, "w1")
	if !uw.IsEquipped {
		t.Error("expected weapon to be equipped")
	}
	if uw.FighterID == nil || *uw.FighterID != "fighter1" {
		t.Error("expected fighter ID to be set")
	}
}

func TestService_EquipWeapon_AlreadyEquipped(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	fid := "fighter1"
	repo.Create(context.Background(), &weapons.UserWeapon{
		ID:        "w1",
		UserID:    1,
		WeaponID:  "wpn_sword_iron_002",
		FighterID: &fid,
		IsEquipped: true,
		Created:   time.Now(),
	})

	err := svc.EquipWeapon(context.Background(), 1, "w1", "fighter2")
	if err != ErrWeaponAlreadyEquipped {
		t.Errorf("expected ErrWeaponAlreadyEquipped, got %v", err)
	}
}

func TestService_UnequipWeapon(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	fid := "fighter1"
	repo.Create(context.Background(), &weapons.UserWeapon{
		ID:        "w1",
		UserID:    1,
		WeaponID:  "wpn_sword_iron_002",
		FighterID: &fid,
		IsEquipped: true,
		Created:   time.Now(),
	})

	err := svc.UnequipWeapon(context.Background(), 1, "w1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	uw, _ := repo.GetByID(context.Background(), 1, "w1")
	if uw.IsEquipped {
		t.Error("expected weapon to be unequipped")
	}
}

func TestService_EnhanceWeapon(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	repo.Create(context.Background(), &weapons.UserWeapon{
		ID:          "w1",
		UserID:      1,
		WeaponID:    "wpn_sword_iron_002",
		Enhancement: 0,
		Created:     time.Now(),
	})

	result, err := svc.EnhanceWeapon(context.Background(), 1, "w1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Level 0-3 should always succeed (0% failure chance)
	if !result.Success {
		t.Error("expected enhancement to succeed at low levels")
	}
	if result.NewLevel != 1 {
		t.Errorf("expected level 1, got %d", result.NewLevel)
	}
}

func TestService_EnhanceWeapon_MaxLevel(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	repo.Create(context.Background(), &weapons.UserWeapon{
		ID:          "w1",
		UserID:      1,
		WeaponID:    "wpn_sword_iron_002",
		Enhancement: 10,
		Created:     time.Now(),
	})

	_, err := svc.EnhanceWeapon(context.Background(), 1, "w1")
	if err != ErrMaxEnhancement {
		t.Errorf("expected ErrMaxEnhancement, got %v", err)
	}
}

func TestService_AddWeaponToInventory(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	uw, err := svc.AddWeaponToInventory(context.Background(), 1, "wpn_sword_excalibur_006")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if uw.WeaponID != "wpn_sword_excalibur_006" {
		t.Errorf("expected Excalibur ID, got %s", uw.WeaponID)
	}
	if uw.Enhancement != 0 {
		t.Error("expected enhancement level 0")
	}
}

func TestService_AddWeaponToInventory_InvalidWeapon(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	_, err := svc.AddWeaponToInventory(context.Background(), 1, "invalid_weapon_id")
	if err != ErrWeaponNotFound {
		t.Errorf("expected ErrWeaponNotFound, got %v", err)
	}
}

func TestService_AddWeaponToInventory_Full(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)
	repo.count = 50 // Simulate full inventory

	_, err := svc.AddWeaponToInventory(context.Background(), 1, "wpn_sword_iron_002")
	if err != ErrInventoryFull {
		t.Errorf("expected ErrInventoryFull, got %v", err)
	}
}

func TestService_PreviewEnhancement(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	repo.Create(context.Background(), &weapons.UserWeapon{
		ID:          "w1",
		UserID:      1,
		WeaponID:    "wpn_sword_iron_002",
		Enhancement: 5,
		Created:     time.Now(),
	})

	weaponDef, nextLevel, successChance, cost, err := svc.PreviewEnhancement(context.Background(), 1, "w1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if weaponDef == nil {
		t.Fatal("expected weapon definition")
	}
	if nextLevel != 6 {
		t.Errorf("expected next level 6, got %d", nextLevel)
	}
	if successChance != 0.85 { // 1 - 0.15 failure chance
		t.Errorf("expected 85%% success chance, got %f", successChance)
	}
	if cost <= 0 {
		t.Error("expected positive cost")
	}
}