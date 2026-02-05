package weapons

// WeaponDatabase contains all static weapon definitions
var WeaponDatabase = []Weapon{
	// Swords (4 weapons)
	{
		ID:          "wpn_sword_rusty_001",
		Name:        "Rusty Blade",
		Type:        Sword,
		Rarity:      Common,
		BaseDamage:  12,
		AttackSpeed: 1.0,
		CritChance:  5,
		Durability:  100,
		IconURL:     "https://vibemedia.space/wpn_sword_rusty_9a8f7b2c.png?prompt=rusty%20old%20sword%20with%20worn%20blade&style=pixel_game_asset&key=NOGON",
		Description: "An old blade that has seen better days.",
	},
	{
		ID:          "wpn_sword_iron_002",
		Name:        "Iron Sword",
		Type:        Sword,
		Rarity:      Common,
		BaseDamage:  18,
		AttackSpeed: 1.0,
		CritChance:  5,
		Durability:  100,
		IconURL:     "https://vibemedia.space/wpn_sword_iron_3k8m9n1p.png?prompt=iron%20sword%20with%20straight%20blade&style=pixel_game_asset&key=NOGON",
		Description: "A standard iron sword.",
	},
	{
		ID:          "wpn_sword_steel_003",
		Name:        "Steel Saber",
		Type:        Sword,
		Rarity:      Rare,
		BaseDamage:  28,
		AttackSpeed: 1.1,
		CritChance:  10,
		Durability:  120,
		IconURL:     "https://vibemedia.space/wpn_sword_steel_8d7e6f5a.png?prompt=steel%20saber%20with%20curved%20blade&style=pixel_game_asset&key=NOGON",
		Description: "A finely crafted steel saber.",
	},
	{
		ID:          "wpn_sword_knight_004",
		Name:        "Knight's Blade",
		Type:        Sword,
		Rarity:      Epic,
		BaseDamage:  45,
		AttackSpeed: 1.2,
		CritChance:  15,
		Durability:  150,
		IconURL:     "https://vibemedia.space/wpn_sword_knight_2c4b6a8e.png?prompt=knightly%20sword%20with%20ornate%20hilt%20and%20golden%20accents&style=pixel_game_asset&key=NOGON",
		Description: "A blade worthy of a knight.",
	},
	{
		ID:          "wpn_sword_excalibur_005",
		Name:        "Excalibur",
		Type:        Sword,
		Rarity:      Legendary,
		BaseDamage:  75,
		AttackSpeed: 1.3,
		CritChance:  22,
		Durability:  200,
		IconURL:     "https://vibemedia.space/wpn_sword_excalibur_1a3c5e7g.png?prompt=legendary%20excalibur%20sword%20glowing%20with%20holy%20light&style=pixel_game_asset&key=NOGON",
		Description: "The legendary sword of kings.",
	},

	// Bows (4 weapons)
	{
		ID:          "wpn_bow_short_001",
		Name:        "Short Bow",
		Type:        Bow,
		Rarity:      Common,
		BaseDamage:  10,
		AttackSpeed: 1.2,
		CritChance:  8,
		Durability:  80,
		IconURL:     "https://vibemedia.space/wpn_bow_short_4b5c6d7e.png?prompt=short%20bow%20made%20of%20wood&style=pixel_game_asset&key=NOGON",
		Description: "A simple wooden bow.",
	},
	{
		ID:          "wpn_bow_long_002",
		Name:        "Longbow",
		Type:        Bow,
		Rarity:      Rare,
		BaseDamage:  25,
		AttackSpeed: 1.3,
		CritChance:  12,
		Durability:  100,
		IconURL:     "https://vibemedia.space/wpn_bow_long_9f8e7d6c.png?prompt=longbow%20with%20elegant%20curves&style=pixel_game_asset&key=NOGON",
		Description: "A longbow with excellent range.",
	},
	{
		ID:          "wpn_bow_elven_003",
		Name:        "Elven Bow",
		Type:        Bow,
		Rarity:      Epic,
		BaseDamage:  40,
		AttackSpeed: 1.4,
		CritChance:  18,
		Durability:  120,
		IconURL:     "https://vibemedia.space/wpn_bow_elven_5a4b3c2d.png?prompt=elven%20bow%20with%20intricate%20carvings%20and%20silver%20accents&style=pixel_game_asset&key=NOGON",
		Description: "A bow crafted by elven masters.",
	},
	{
		ID:          "wpn_bow_dragon_004",
		Name:        "Dragon Bow",
		Type:        Bow,
		Rarity:      Legendary,
		BaseDamage:  65,
		AttackSpeed: 1.4,
		CritChance:  25,
		Durability:  180,
		IconURL:     "https://vibemedia.space/wpn_bow_dragon_7e6f5g4h.png?prompt=dragon%20bow%20with%20red%20scales%20and%20fiery%20glow&style=pixel_game_asset&key=NOGON",
		Description: "Forged from dragon bone.",
	},

	// Staves (4 weapons)
	{
		ID:          "wpn_staff_wooden_001",
		Name:        "Wooden Staff",
		Type:        Staff,
		Rarity:      Common,
		BaseDamage:  15,
		AttackSpeed: 0.9,
		CritChance:  5,
		Durability:  90,
		IconURL:     "https://vibemedia.space/wpn_staff_wooden_2b3c4d5e.png?prompt=wooden%20magic%20staff&style=pixel_game_asset&key=NOGON",
		Description: "A basic wooden staff.",
	},
	{
		ID:          "wpn_staff_mystic_002",
		Name:        "Mystic Staff",
		Type:        Staff,
		Rarity:      Rare,
		BaseDamage:  30,
		AttackSpeed: 1.0,
		CritChance:  10,
		Durability:  110,
		IconURL:     "https://vibemedia.space/wpn_staff_mystic_6f7g8h9i.png?prompt=mystic%20staff%20with%20glowing%20crystal%20orb&style=pixel_game_asset&key=NOGON",
		Description: "A staff infused with arcane energy.",
	},
	{
		ID:          "wpn_staff_archmage_003",
		Name:        "Archmage Rod",
		Type:        Staff,
		Rarity:      Epic,
		BaseDamage:  50,
		AttackSpeed: 1.1,
		CritChance:  15,
		Durability:  140,
		IconURL:     "https://vibemedia.space/wpn_staff_archmage_3j4k5l6m.png?prompt=archmage%20rod%20with%20multiple%20floating%20gems&style=pixel_game_asset&key=NOGON",
		Description: "Wielded by the greatest mages.",
	},
	{
		ID:          "wpn_staff_void_004",
		Name:        "Void Scepter",
		Type:        Staff,
		Rarity:      Legendary,
		BaseDamage:  70,
		AttackSpeed: 1.2,
		CritChance:  22,
		Durability:  170,
		IconURL:     "https://vibemedia.space/wpn_staff_void_8n9o0p1q.png?prompt=void%20scepter%20with%20dark%20purple%20energy%20and%20void%20portal&style=pixel_game_asset&key=NOGON",
		Description: "Channels the power of the void.",
	},

	// Daggers (3 weapons)
	{
		ID:          "wpn_dagger_steel_001",
		Name:        "Steel Dagger",
		Type:        Dagger,
		Rarity:      Common,
		BaseDamage:  8,
		AttackSpeed: 1.4,
		CritChance:  10,
		Durability:  70,
		IconURL:     "https://vibemedia.space/wpn_dagger_steel_4r5s6t7u.png?prompt=steel%20dagger%20with%20sharp%20blade&style=pixel_game_asset&key=NOGON",
		Description: "A sharp steel dagger.",
	},
	{
		ID:          "wpn_dagger_shadow_002",
		Name:        "Shadow Knife",
		Type:        Dagger,
		Rarity:      Rare,
		BaseDamage:  22,
		AttackSpeed: 1.5,
		CritChance:  15,
		Durability:  90,
		IconURL:     "https://vibemedia.space/wpn_dagger_shadow_2v3w4x5y.png?prompt=shadow%20knife%20with%20dark%20smoky%20aura&style=pixel_game_asset&key=NOGON",
		Description: "Blade forged in shadows.",
	},
	{
		ID:          "wpn_dagger_assassin_003",
		Name:        "Assassin's Blade",
		Type:        Dagger,
		Rarity:      Epic,
		BaseDamage:  38,
		AttackSpeed: 1.5,
		CritChance:  20,
		Durability:  110,
		IconURL:     "https://vibemedia.space/wpn_dagger_assassin_6z7a8b9c.png?prompt=assassin%20blade%20with%20black%20handle%20and%20poison%20green%20edge&style=pixel_game_asset&key=NOGON",
		Description: "Preferred by master assassins.",
	},

	// Axes (4 weapons - including the Mythic)
	{
		ID:          "wpn_axe_war_001",
		Name:        "War Axe",
		Type:        Axe,
		Rarity:      Common,
		BaseDamage:  20,
		AttackSpeed: 0.8,
		CritChance:  5,
		Durability:  100,
		IconURL:     "https://vibemedia.space/wpn_axe_war_1d2e3f4g.png?prompt=war%20axe%20with%20heavy%20blade&style=pixel_game_asset&key=NOGON",
		Description: "A heavy battle axe.",
	},
	{
		ID:          "wpn_axe_battle_002",
		Name:        "Battle Axe",
		Type:        Axe,
		Rarity:      Rare,
		BaseDamage:  35,
		AttackSpeed: 0.9,
		CritChance:  10,
		Durability:  120,
		IconURL:     "https://vibemedia.space/wpn_axe_battle_5h6i7j8k.png?prompt=battle%20axe%20with%20double%20blades&style=pixel_game_asset&key=NOGON",
		Description: "A balanced battle axe.",
	},
	{
		ID:          "wpn_axe_berserker_003",
		Name:        "Berserker Axe",
		Type:        Axe,
		Rarity:      Epic,
		BaseDamage:  55,
		AttackSpeed: 1.0,
		CritChance:  15,
		Durability:  150,
		IconURL:     "https://vibemedia.space/wpn_axe_berserker_9l0m1n2o.png?prompt=berserker%20axe%20with%20spiked%20design%20and%20blood%20grooves&style=pixel_game_asset&key=NOGON",
		Description: "Unleashes primal fury.",
	},
	{
		ID:          "wpn_axe_worldender_004",
		Name:        "World Ender",
		Type:        Axe,
		Rarity:      Mythic,
		BaseDamage:  95,
		AttackSpeed: 1.2,
		CritChance:  30,
		Durability:  250,
		IconURL:     "https://vibemedia.space/wpn_axe_worldender_3p4q5r6s.png?prompt=mythic%20world%20ender%20axe%20with%20cosmic%20energy%20and%20ancient%20runes&style=pixel_game_asset&key=NOGON",
		Description: "A weapon of apocalyptic power.",
	},
}

// GetWeaponByID retrieves a weapon definition by ID
func GetWeaponByID(id string) (*Weapon, bool) {
	for i := range WeaponDatabase {
		if WeaponDatabase[i].ID == id {
			return &WeaponDatabase[i], true
		}
	}
	return nil, false
}

// GetWeaponsByType returns all weapons of a specific type
func GetWeaponsByType(weaponType WeaponType) []Weapon {
	var result []Weapon
	for _, w := range WeaponDatabase {
		if w.Type == weaponType {
			result = append(result, w)
		}
	}
	return result
}

// GetWeaponsByRarity returns all weapons of a specific rarity
func GetWeaponsByRarity(rarity Rarity) []Weapon {
	var result []Weapon
	for _, w := range WeaponDatabase {
		if w.Rarity == rarity {
			result = append(result, w)
		}
	}
	return result
}