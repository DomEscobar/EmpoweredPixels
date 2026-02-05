package weapons

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math/big"
	"time"

	"empoweredpixels/internal/domain/weapons"
)

var (
	ErrWeaponNotFound      = errors.New("weapon not found")
	ErrInventoryFull       = errors.New("weapon inventory is full")
	ErrWeaponAlreadyEquipped = errors.New("weapon already equipped")
	ErrFighterHasWeapon    = errors.New("fighter already has a weapon equipped")
	ErrMaxEnhancement      = errors.New("weapon is at maximum enhancement")
	ErrEnhancementFailed   = errors.New("enhancement failed")
)

type WeaponRepository interface {
	GetByID(ctx context.Context, userID int64, id string) (*weapons.UserWeapon, error)
	ListByUser(ctx context.Context, userID int64, limit int, offset int) ([]weapons.UserWeapon, error)
	ListByUserAll(ctx context.Context, userID int64) ([]weapons.UserWeapon, error)
	CountByUser(ctx context.Context, userID int64) (int, error)
	Create(ctx context.Context, userWeapon *weapons.UserWeapon) error
	UpdateEnhancement(ctx context.Context, id string, enhancement int) error
	UpdateFighter(ctx context.Context, id string, fighterID *string) error
	Delete(ctx context.Context, id string) error
	GetEquippedByFighter(ctx context.Context, fighterID string) (*weapons.UserWeapon, error)
}

type Service struct {
	repo WeaponRepository
}

func NewService(repo WeaponRepository) *Service {
	return &Service{repo: repo}
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// ListUserWeapons returns all weapons owned by a user
func (s *Service) ListUserWeapons(ctx context.Context, userID int64) ([]weapons.UserWeapon, error) {
	return s.repo.ListByUserAll(ctx, userID)
}

// GetWeaponDetails returns full weapon details with stats
func (s *Service) GetWeaponDetails(ctx context.Context, userID int64, weaponID string) (*weapons.UserWeapon, *weapons.Weapon, *weapons.WeaponStats, error) {
	uw, err := s.repo.GetByID(ctx, userID, weaponID)
	if err != nil {
		return nil, nil, nil, err
	}
	if uw == nil {
		return nil, nil, nil, ErrWeaponNotFound
	}

	weaponDef, found := weapons.GetWeaponByID(uw.WeaponID)
	if !found {
		return nil, nil, nil, ErrWeaponNotFound
	}

	stats := weapons.CalculateStats(weaponDef, uw.Enhancement)
	return uw, weaponDef, &stats, nil
}

// EquipWeapon equips a weapon to a fighter
func (s *Service) EquipWeapon(ctx context.Context, userID int64, weaponID string, fighterID string) error {
	// Check if weapon exists and belongs to user
	uw, err := s.repo.GetByID(ctx, userID, weaponID)
	if err != nil {
		return err
	}
	if uw == nil {
		return ErrWeaponNotFound
	}
	if uw.IsEquipped {
		return ErrWeaponAlreadyEquipped
	}

	// Check if fighter already has a weapon
	equipped, err := s.repo.GetEquippedByFighter(ctx, fighterID)
	if err != nil {
		return err
	}
	if equipped != nil {
		return ErrFighterHasWeapon
	}

	// Equip the weapon
	return s.repo.UpdateFighter(ctx, weaponID, &fighterID)
}

// UnequipWeapon removes a weapon from a fighter
func (s *Service) UnequipWeapon(ctx context.Context, userID int64, weaponID string) error {
	uw, err := s.repo.GetByID(ctx, userID, weaponID)
	if err != nil {
		return err
	}
	if uw == nil {
		return ErrWeaponNotFound
	}
	if !uw.IsEquipped {
		return nil // Already unequipped
	}

	return s.repo.UpdateFighter(ctx, weaponID, nil)
}

// EnhanceWeapon attempts to enhance a weapon
func (s *Service) EnhanceWeapon(ctx context.Context, userID int64, weaponID string) (*weapons.EnhancementResult, error) {
	uw, err := s.repo.GetByID(ctx, userID, weaponID)
	if err != nil {
		return nil, err
	}
	if uw == nil {
		return nil, ErrWeaponNotFound
	}

	if !weapons.CanEnhance(uw.Enhancement) {
		return nil, ErrMaxEnhancement
	}

	// Get weapon definition for rarity
	weaponDef, found := weapons.GetWeaponByID(uw.WeaponID)
	if !found {
		return nil, ErrWeaponNotFound
	}

	// Calculate success chance
	failureChance := weapons.EnhancementFailureChance(uw.Enhancement + 1)
	successChance := 1.0 - failureChance

	// Roll for success
	success := rollSuccess(successChance)

	// Apply enhancement
	result := weapons.ApplyEnhancement(uw, success)

	// Update database
	if err := s.repo.UpdateEnhancement(ctx, weaponID, uw.Enhancement); err != nil {
		return nil, err
	}

	// Log the cost (would integrate with economy service in full implementation)
	_ = weapons.EnhancementCost(weaponDef.Rarity, uw.Enhancement)

	return &result, nil
}

// PreviewEnhancement returns enhancement odds without modifying anything
func (s *Service) PreviewEnhancement(ctx context.Context, userID int64, weaponID string) (*weapons.Weapon, int, float64, int, error) {
	uw, err := s.repo.GetByID(ctx, userID, weaponID)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	if uw == nil {
		return nil, 0, 0, 0, ErrWeaponNotFound
	}

	weaponDef, found := weapons.GetWeaponByID(uw.WeaponID)
	if !found {
		return nil, 0, 0, 0, ErrWeaponNotFound
	}

	nextLevel := uw.Enhancement + 1
	if nextLevel > weapons.MaxEnhancement {
		nextLevel = weapons.MaxEnhancement
	}

	successChance := 1.0 - weapons.EnhancementFailureChance(nextLevel)
	cost := weapons.EnhancementCost(weaponDef.Rarity, uw.Enhancement)

	return weaponDef, nextLevel, successChance, cost, nil
}

// AddWeaponToInventory adds a new weapon to user's inventory
func (s *Service) AddWeaponToInventory(ctx context.Context, userID int64, weaponDefID string) (*weapons.UserWeapon, error) {
	// Check if weapon definition exists
	_, found := weapons.GetWeaponByID(weaponDefID)
	if !found {
		return nil, ErrWeaponNotFound
	}

	// Check inventory space
	count, err := s.repo.CountByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if count >= weapons.MaxInventorySlots {
		return nil, ErrInventoryFull
	}

	// Create the user weapon
	uw := &weapons.UserWeapon{
		ID:          generateID(),
		UserID:      userID,
		WeaponID:    weaponDefID,
		Enhancement: 0,
		Durability:  100,
		IsEquipped:  false,
		Created:     time.Now(),
	}

	if err := s.repo.Create(ctx, uw); err != nil {
		return nil, err
	}

	return uw, nil
}

// GetFighterWeapon gets the weapon equipped by a fighter
func (s *Service) GetFighterWeapon(ctx context.Context, fighterID string) (*weapons.UserWeapon, *weapons.Weapon, error) {
	uw, err := s.repo.GetEquippedByFighter(ctx, fighterID)
	if err != nil {
		return nil, nil, err
	}
	if uw == nil {
		return nil, nil, nil
	}

	weaponDef, found := weapons.GetWeaponByID(uw.WeaponID)
	if !found {
		return nil, nil, ErrWeaponNotFound
	}

	return uw, weaponDef, nil
}

// rollSuccess returns true if the roll succeeds based on probability (0.0-1.0)
func rollSuccess(probability float64) bool {
	n, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		// Fallback to time-based seed if crypto rand fails
		return probability > 0.5
	}
	roll := float64(n.Int64()) / 10000.0
	return roll < probability
}