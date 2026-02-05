package weapons

import (
	"testing"
)

func TestWeaponDatabase_Completeness(t *testing.T) {
	// Verify we have exactly 24 weapons (added Divine and Unique tiers)
	if len(WeaponDatabase) != 24 {
		t.Errorf("Expected 24 weapons, got %d", len(WeaponDatabase))
	}

	// Count by type
	typeCounts := make(map[WeaponType]int)
	rarityCounts := make(map[Rarity]int)

	for _, w := range WeaponDatabase {
		typeCounts[w.Type]++
		rarityCounts[w.Rarity]++
	}

	// Verify type distribution
	expectedTypes := map[WeaponType]int{
		Sword:  8,
		Bow:    4,
		Staff:  5,
		Dagger: 3,
		Axe:    4,
	}

	for weaponType, expected := range expectedTypes {
		if typeCounts[weaponType] != expected {
			t.Errorf("Expected %d %s weapons, got %d", expected, weaponType.String(), typeCounts[weaponType])
		}
	}

	// Verify we have at least 1 of each rarity (including new ones)
	expectedRarities := []Rarity{Broken, Common, Uncommon, Rare, Epic, Legendary, Mythic, Divine, Unique}
	for _, r := range expectedRarities {
		if rarityCounts[r] == 0 {
			t.Errorf("Expected at least 1 %s weapon", r.String())
		}
	}
}

func TestWeaponDatabase_UniqueIDs(t *testing.T) {
	ids := make(map[string]bool)
	for _, w := range WeaponDatabase {
		if ids[w.ID] {
			t.Errorf("Duplicate weapon ID: %s", w.ID)
		}
		ids[w.ID] = true
	}
}

func TestWeaponDatabase_ValidStats(t *testing.T) {
	for _, w := range WeaponDatabase {
		if w.BaseDamage <= 0 {
			t.Errorf("Weapon %s has invalid damage: %d", w.Name, w.BaseDamage)
		}
		if w.AttackSpeed <= 0 {
			t.Errorf("Weapon %s has invalid attack speed: %f", w.Name, w.AttackSpeed)
		}
		if w.CritChance < 0 || w.CritChance > 100 {
			t.Errorf("Weapon %s has invalid crit chance: %d", w.Name, w.CritChance)
		}
		if w.Durability <= 0 {
			t.Errorf("Weapon %s has invalid durability: %d", w.Name, w.Durability)
		}
	}
}

func TestWeaponDatabase_WorldEnder(t *testing.T) {
	// Special test for the Mythic weapon
	worldEnder, found := GetWeaponByID("wpn_axe_worldender_004")
	if !found {
		t.Fatal("World Ender not found")
	}

	if worldEnder.Rarity != Mythic {
		t.Errorf("World Ender should be Mythic rarity")
	}
	if worldEnder.BaseDamage != 95 {
		t.Errorf("World Ender should have 95 damage, got %d", worldEnder.BaseDamage)
	}
	if worldEnder.CritChance != 30 {
		t.Errorf("World Ender should have 30%% crit chance, got %d", worldEnder.CritChance)
	}
}
