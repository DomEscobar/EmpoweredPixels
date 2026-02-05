package skills

// SkillDatabase contains all 15 skill definitions (5 per branch, tiers 1-3 MVP)
var SkillDatabase = []Skill{
	// ===== OFFENSE BRANCH =====
	{
		ID:          "skl_power_strike",
		Name:        "Power Strike",
		Branch:      Offense,
		Tier:        1,
		Type:        Active,
		Description: "A powerful strike dealing 40-70% bonus damage based on rank.",
		ManaCost:    20,
		Cooldown:    8,
		EffectValue: 40,
		Duration:    0,
		IconURL:     "https://vibemedia.space/skl_power_strike_001.png?prompt=power%20strike%20skill%20icon%20with%20glowing%20sword&style=pixel_game_asset&key=NOGON",
	},
	{
		ID:          "skl_bleed",
		Name:        "Bleed",
		Branch:      Offense,
		Tier:        1,
		Type:        Passive,
		Description: "Attacks have a chance to cause bleeding, dealing damage over time.",
		ManaCost:    0,
		Cooldown:    0,
		EffectValue: 5, // damage per second
		Duration:    4,
		IconURL:     "https://vibemedia.space/skl_bleed_001.png?prompt=bleed%20skill%20icon%20with%20blood%20drops&style=pixel_game_asset&key=NOGON",
	},
	{
		ID:          "skl_whirlwind",
		Name:        "Whirlwind",
		Branch:      Offense,
		Tier:        2,
		Type:        Active,
		Description: "Spin attack hitting all nearby enemies for 80-120% weapon damage.",
		ManaCost:    40,
		Cooldown:    15,
		EffectValue: 80,
		Duration:    0,
		IconURL:     "https://vibemedia.space/skl_whirlwind_001.png?prompt=whirlwind%20attack%20skill%20icon%20with%20spinning%20blades&style=pixel_game_asset&key=NOGON",
	},
	{
		ID:          "skl_berserk",
		Name:        "Berserk",
		Branch:      Offense,
		Tier:        2,
		Type:        Active,
		Description: "Trade armor for damage. +20-40% damage, -15% to -5% armor.",
		ManaCost:    50,
		Cooldown:    30,
		EffectValue: 20, // % damage increase
		Duration:    10,
		IconURL:     "https://vibemedia.space/skl_berserk_001.png?prompt=berserk%20rage%20skill%20icon%20with%20red%20aura&style=pixel_game_asset&key=NOGON",
	},
	{
		ID:          "skl_execute",
		Name:        "Execute",
		Branch:      Offense,
		Tier:        3,
		Type:        Active,
		Description: "Instantly kill enemies below 15-25% health.",
		ManaCost:    70,
		Cooldown:    45,
		EffectValue: 15, // % health threshold
		Duration:    0,
		IconURL:     "https://vibemedia.space/skl_execute_001.png?prompt=execute%20skill%20icon%20with%20skull%20and%20scythe&style=pixel_game_asset&key=NOGON",
	},

	// ===== DEFENSE BRANCH =====
	{
		ID:          "skl_block",
		Name:        "Block",
		Branch:      Defense,
		Tier:        1,
		Type:        Active,
		Description: "Block 1-3 incoming attacks completely.",
		ManaCost:    15,
		Cooldown:    8,
		EffectValue: 1, // number of attacks
		Duration:    5,
		IconURL:     "https://vibemedia.space/skl_block_001.png?prompt=shield%20block%20skill%20icon%20with%20blue%20barrier&style=pixel_game_asset&key=NOGON",
	},
	{
		ID:          "skl_heal",
		Name:        "Heal",
		Branch:      Defense,
		Tier:        1,
		Type:        Active,
		Description: "Restore 15-35% of maximum health.",
		ManaCost:    30,
		Cooldown:    18,
		EffectValue: 15, // % of max HP
		Duration:    0,
		IconURL:     "https://vibemedia.space/skl_heal_001.png?prompt=heal%20skill%20icon%20with%20green%20cross%20and%20light&style=pixel_game_asset&key=NOGON",
	},
	{
		ID:          "skl_shield_wall",
		Name:        "Shield Wall",
		Branch:      Defense,
		Tier:        2,
		Type:        Active,
		Description: "Increase armor by 40-60% for 6-10 seconds.",
		ManaCost:    45,
		Cooldown:    25,
		EffectValue: 40, // % armor increase
		Duration:    6,
		IconURL:     "https://vibemedia.space/skl_shield_wall_001.png?prompt=shield%20wall%20skill%20icon%20with%20golden%20barrier&style=pixel_game_asset&key=NOGON",
	},
	{
		ID:          "skl_reflect",
		Name:        "Reflect",
		Branch:      Defense,
		Tier:        2,
		Type:        Passive,
		Description: "Return 20-40% of damage taken back to attacker.",
		ManaCost:    0,
		Cooldown:    0,
		EffectValue: 20, // % damage reflected
		Duration:    0,
		IconURL:     "https://vibemedia.space/skl_reflect_001.png?prompt=reflect%20skill%20icon%20with%20mirrored%20shield&style=pixel_game_asset&key=NOGON",
	},
	{
		ID:          "skl_immortal",
		Name:        "Immortal",
		Branch:      Defense,
		Tier:        3,
		Type:        Active,
		Description: "Become invulnerable for 3-5 seconds.",
		ManaCost:    90,
		Cooldown:    75,
		EffectValue: 3, // seconds
		Duration:    3,
		IconURL:     "https://vibemedia.space/skl_immortal_001.png?prompt=immortal%20skill%20icon%20with%20divine%20light%20and%20wings&style=pixel_game_asset&key=NOGON",
	},

	// ===== UTILITY BRANCH =====
	{
		ID:          "skl_haste",
		Name:        "Haste",
		Branch:      Utility,
		Tier:        1,
		Type:        Active,
		Description: "Increase attack and movement speed by 15-25% for 8 seconds.",
		ManaCost:    25,
		Cooldown:    20,
		EffectValue: 15, // % speed increase
		Duration:    8,
		IconURL:     "https://vibemedia.space/skl_haste_001.png?prompt=haste%20skill%20icon%20with%20wind%20swirls%20and%20speed%20lines&style=pixel_game_asset&key=NOGON",
	},
	{
		ID:          "skl_mana_regen",
		Name:        "Mana Regeneration",
		Branch:      Utility,
		Tier:        1,
		Type:        Passive,
		Description: "Regenerate 5-15 mana per second passively.",
		ManaCost:    0,
		Cooldown:    0,
		EffectValue: 5, // mana per second
		Duration:    0,
		IconURL:     "https://vibemedia.space/skl_mana_regen_001.png?prompt=mana%20regeneration%20skill%20icon%20with%20blue%20energy%20orb&style=pixel_game_asset&key=NOGON",
	},
	{
		ID:          "skl_teleport",
		Name:        "Teleport",
		Branch:      Utility,
		Tier:        2,
		Type:        Active,
		Description: "Instantly teleport to target location within range.",
		ManaCost:    35,
		Cooldown:    12,
		EffectValue: 10, // range in tiles
		Duration:    0,
		IconURL:     "https://vibemedia.space/skl_teleport_001.png?prompt=teleport%20skill%20icon%20with%20purple%20portal%20swirl&style=pixel_game_asset&key=NOGON",
	},
	{
		ID:          "skl_stun",
		Name:        "Stun Strike",
		Branch:      Utility,
		Tier:        2,
		Type:        Active,
		Description: "Stun target enemy for 1.5-3 seconds.",
		ManaCost:    30,
		Cooldown:    16,
		EffectValue: 1500, // milliseconds
		Duration:    0,
		IconURL:     "https://vibemedia.space/skl_stun_001.png?prompt=stun%20skill%20icon%20with%20stars%20and%20impact%20burst&style=pixel_game_asset&key=NOGON",
	},
	{
		ID:          "skl_life_steal",
		Name:        "Life Steal Aura",
		Branch:      Utility,
		Tier:        3,
		Type:        Passive,
		Description: "Gain 8-20% lifesteal on all attacks.",
		ManaCost:    0,
		Cooldown:    0,
		EffectValue: 8, // % lifesteal
		Duration:    0,
		IconURL:     "https://vibemedia.space/skl_lifesteal_001.png?prompt=life%20steal%20skill%20icon%20with%20red%20heart%20and%20vampire%20fangs&style=pixel_game_asset&key=NOGON",
	},
}

// UltimateSkills contains the ultimate abilities (unlocked at level 50)
var UltimateSkills = []Skill{
	{
		ID:          "ult_meteor_strike",
		Name:        "Meteor Strike",
		Branch:      Offense,
		Tier:        5,
		Type:        Ultimate,
		Description: "Call down a meteor dealing 250 AOE damage to all enemies. Charges from dealing damage.",
		ManaCost:    100,
		Cooldown:    0, // Uses charge system
		EffectValue: 250,
		Duration:    0,
		IconURL:     "https://vibemedia.space/ult_meteor_001.png?prompt=meteor%20strike%20ultimate%20icon%20with%20falling%20fire%20rock&style=pixel_game_asset&key=NOGON",
	},
	{
		ID:          "ult_divine_protection",
		Name:        "Divine Protection",
		Branch:      Defense,
		Tier:        5,
		Type:        Ultimate,
		Description: "Full heal + 4 seconds invincibility. Charges from taking damage.",
		ManaCost:    100,
		Cooldown:    0,
		EffectValue: 100, // % heal
		Duration:    4,
		IconURL:     "https://vibemedia.space/ult_divine_001.png?prompt=divine%20protection%20ultimate%20icon%20with%20golden%20dome%20shield&style=pixel_game_asset&key=NOGON",
	},
	{
		ID:          "ult_time_warp",
		Name:        "Time Warp",
		Branch:      Utility,
		Tier:        5,
		Type:        Ultimate,
		Description: "Reset all cooldowns and gain +40% speed for 10 seconds. Charges from dealing damage.",
		ManaCost:    100,
		Cooldown:    0,
		EffectValue: 40, // % speed
		Duration:    10,
		IconURL:     "https://vibemedia.space/ult_time_001.png?prompt=time%20warp%20ultimate%20icon%20with%20clock%20and%20time%20spiral&style=pixel_game_asset&key=NOGON",
	},
}

// GetSkillByID retrieves a skill definition by ID
func GetSkillByID(id string) (*Skill, bool) {
	for i := range SkillDatabase {
		if SkillDatabase[i].ID == id {
			return &SkillDatabase[i], true
		}
	}
	for i := range UltimateSkills {
		if UltimateSkills[i].ID == id {
			return &UltimateSkills[i], true
		}
	}
	return nil, false
}

// GetSkillsByBranch returns all skills in a branch
func GetSkillsByBranch(branch SkillBranch) []Skill {
	var result []Skill
	for _, s := range SkillDatabase {
		if s.Branch == branch {
			result = append(result, s)
		}
	}
	return result
}

// GetSkillsByTier returns all skills at a specific tier
func GetSkillsByTier(tier int) []Skill {
	var result []Skill
	for _, s := range SkillDatabase {
		if s.Tier == tier {
			result = append(result, s)
		}
	}
	return result
}

// GetUltimateByBranch returns the ultimate skill for a branch
func GetUltimateByBranch(branch SkillBranch) (*Skill, bool) {
	for i := range UltimateSkills {
		if UltimateSkills[i].Branch == branch {
			return &UltimateSkills[i], true
		}
	}
	return nil, false
}