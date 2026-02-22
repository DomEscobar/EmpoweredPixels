export function getEquipmentImageUrl(type: string | null | undefined, itemId?: string): { url: string; testid: string } {
  // Defensive check for null/undefined
  if (!type || typeof type !== 'string') {
    const fallbackUrl = 'https://vibemedia.pace/item_generic_789.png?prompt=pixel%20mystery%20item%20sprite&style=pixel_game_asset&key=NOGON';
    return {
      url: fallbackUrl,
      testid: `equipment-image-${itemId || 'fallback'}`
    };
  }

  const t = type.toLowerCase();

  // Known specific item mappings for distinct visuals
   const specificImages: Record<string, string> = {
     // Armor pieces - Chest
     'armor_01': 'https://vibemedia.space/che_plate_armor01_v1.png?prompt=pixel%20ornate%20plate%20chest%20armor%20sprite&style=pixel_game_asset&key=NOGON',
     'chest_armor': 'https://vibemedia.space/che_generic_armor_v1.png?prompt=pixel%20chest%20armor%20sprite&style=pixel_game_asset&key=NOGON',
     'leather_armor': 'https://vibemedia.space/che_leather_armor_v1.png?prompt=pixel%20leather%20tunic%20chest%20armor%20sprite&style=pixel_game_asset&key=NOGON',
     'plate_chest': 'https://vibemedia.space/che_plate_armor02_v1.png?prompt=pixel%20heavy%20plate%20chest%20armor%20sprite&style=pixel_game_asset&key=NOGON',
     'mail_chest': 'https://vibemedia.space/che_mail_armor_v1.png?prompt=pixel%20chainmail%20chest%20armor%20sprite&style=pixel_game_asset&key=NOGON',
     'cloth_chest': 'https://vibemedia.space/che_cloth_armor_v1.png?prompt=pixel%20cloth%20robes%20chest%20armor%20sprite&style=pixel_game_asset&key=NOGON',

    // Armor pieces - Helmet
    'helmet_01': 'https://vibemedia.space/hel_plate_helmet01_v1.png?prompt=pixel%20plate%20helmet%20sprite&style=pixel_game_asset&key=NOGON',
    'leather_helmet': 'https://vibemedia.space/hel_leather_helmet_v1.png?prompt=pixel%20leather%20cap%20sprite&style=pixel_game_asset&key=NOGON',
    'iron_helmet': 'https://vibemedia.space/hel_iron_helmet_v1.png?prompt=pixel%20iron%20helmet%20sprite&style=pixel_game_asset&key=NOGON',
    'plate_helmet': 'https://vibemedia.space/hel_plate_helmet02_v1.png?prompt=pixel%20full%20plate%20helmet%20sprite&style=pixel_game_asset&key=NOGON',

     // Armor pieces - Gloves
     'gloves_01': 'https://vibemedia.space/glv_plate_gloves01_v1.png?prompt=pixel%20plate%20gauntlets%20sprite&style=pixel_game_asset&key=NOGON',
     'gauntlets': 'https://vibemedia.space/glv_generic_gauntlets_v1.png?prompt=pixel%20gauntlets%20sprite&style=pixel_game_asset&key=NOGON',
     'leather_gloves': 'https://vibemedia.space/glv_leather_gloves_v1.png?prompt=pixel%20leather%20gloves%20sprite&style=pixel_game_asset&key=NOGON',
     'mail_gloves': 'https://vibemedia.space/glv_mail_gloves_v1.png?prompt=pixel%20chainmail%20gloves%20sprite&style=pixel_game_asset&key=NOGON',

     // Armor pieces - Boots
     'boots_01': 'https://vibemedia.space/bot_plate_boots01_v1.png?prompt=pixel%20plate%20boots%20sprite&style=pixel_game_asset&key=NOGON',
     'greaves': 'https://vibemedia.space/bot_generic_greaves_v1.png?prompt=pixel%20greaves%20sprite&style=pixel_game_asset&key=NOGON',
     'leather_boots': 'https://vibemedia.space/bot_leather_boots_v1.png?prompt=pixel%20leather%20boots%20sprite&style=pixel_game_asset&key=NOGON',
     'mail_boots': 'https://vibemedia.space/bot_mail_boots_v1.png?prompt=pixel%20chainmail%20boots%20sprite&style=pixel_game_asset&key=NOGON',

    // Weapons
    'sword_01': 'https://vibemedia.space/wpn_steel_sword_123.png?prompt=pixel%20steel%20longsword%20sprite&style=pixel_game_asset&key=NOGON',
    'iron_sword': 'https://vibemedia.space/wpn_iron_sword_123.png?prompt=pixel%20iron%20broadsword%20sprite&style=pixel_game_asset&key=NOGON',
    'steel_sword': 'https://vibemedia.space/wpn_steel_sword_123.png?prompt=pixel%20steel%20longsword%20sprite&style=pixel_game_asset&key=NOGON',
    'greatsword': 'https://vibemedia.space/wpn_greatsword_123.png?prompt=pixel%20two-handed%20greatsword%20sprite&style=pixel_game_asset&key=NOGON',
    'axe_01': 'https://vibemedia.space/wpn_war_axe_123.png?prompt=pixel%20war%20axe%20sprite&style=pixel_game_asset&key=NOGON',
    'battle_axe': 'https://vibemedia.space/wpn_battle_axe_123.png?prompt=pixel%20battle%20axe%20sprite&style=pixel_game_asset&key=NOGON',
    'bow_01': 'https://vibemedia.space/wpn_longbow_123.png?prompt=pixel%20longbow%20sprite&style=pixel_game_asset&key=NOGON',
    'shortbow': 'https://vibemedia.space/wpn_shortbow_123.png?prompt=pixel%20shortbow%20sprite&style=pixel_game_asset&key=NOGON',
    'staff_01': 'https://vibemedia.space/wpn_staff_123.png?prompt=pixel%20wooden%20staff%20sprite&style=pixel_game_asset&key=NOGON',
    'magic_staff': 'https://vibemedia.space/wpn_magic_staff_123.png?prompt=pixel%20ornate%20magic%20staff%20sprite&style=pixel_game_asset&key=NOGON',
    'dagger': 'https://vibemedia.space/wpn_dagger_123.png?prompt=pixel%20dagger%20sprite&style=pixel_game_asset&key=NOGON',
    'spear': 'https://vibemedia.space/wpn_spear_123.png?prompt=pixel%20spear%20sprite&style=pixel_game_asset&key=NOGON',
    'hammer': 'https://vibemedia.space/wpn_warhammer_123.png?prompt=pixel%20warhammer%20sprite&style=pixel_game_asset&key=NOGON',

    // Shields
    'shield_01': 'https://vibemedia.space/shld_wooden_shield_123.png?prompt=pixel%20wooden%20shield%20sprite&style=pixel_game_asset&key=NOGON',
    'iron_shield': 'https://vibemedia.space/shld_iron_shield_123.png?prompt=pixel%20iron%20shield%20sprite&style=pixel_game_asset&key=NOGON',
    'tower_shield': 'https://vibemedia.space/shld_tower_shield_123.png?prompt=pixel%20tower%20shield%20sprite&style=pixel_game_asset&key=NOGON',
    'buckler': 'https://vibemedia.space/shld_buckler_123.png?prompt=pixel%20buckler%20shield%20sprite&style=pixel_game_asset&key=NOGON',

    // Trinkets
    'ring_01': 'https://vibemedia.space/ring_gold_456.png?prompt=pixel%20gold%20ring%20sprite&style=pixel_game_asset&key=NOGON',
    'silver_ring': 'https://vibemedia.space/ring_silver_456.png?prompt=pixel%20silver%20ring%20sprite&style=pixel_game_asset&key=NOGON',
    'necklace_01': 'https://vibemedia.space/neck_gold_456.png?prompt=pixel%20gold%20necklace%20sprite&style=pixel_game_asset&key=NOGON',
    'amulet_01': 'https://vibemedia.space/neck_amulet_456.png?prompt=pixel%20amulet%20sprite&style=pixel_game_asset&key=NOGON',
  };

  let url: string;

  if (specificImages[t]) {
    url = specificImages[t];
  } else {
    // Generate dynamic URLs for different equipment categories to ensure uniqueness
    const wpnPatterns = ['wpn', 'sword', 'axe', 'bow', 'staff', 'dagger', 'spear', 'hammer', 'blade'];
    const armorPatterns = ['helmet', 'chest', 'gloves', 'boots', 'armor', 'shield', 'plate', 'mail', 'leather', 'cloth'];
    const trinketPatterns = ['ring', 'necklace', 'amulet', 'trinket', 'charm', 'pendant'];

    const isWeapon = wpnPatterns.some(pattern => t.includes(pattern));
    const isArmor = armorPatterns.some(pattern => t.includes(pattern));
    const isTrinket = trinketPatterns.some(pattern => t.includes(pattern));

    if (isWeapon) {
      const encodedType = encodeURIComponent(t.replace(/_/g, ' '));
      url = `https://vibemedia.space/equipment/wpn_${t}_123.png?prompt=pixel%20${encodedType}%20sprite%20item%20${itemId || ''}&style=pixel_game_asset&key=NOGON`;
    } else if (isArmor) {
      const encodedType = encodeURIComponent(t.replace(/_/g, ' '));
      url = `https://vibemedia.space/equipment/arm_${t}_123.png?prompt=pixel%20${encodedType}%20sprite%20item%20${itemId || ''}&style=pixel_game_asset&key=NOGON`;
    } else if (isTrinket) {
      const encodedType = encodeURIComponent(t.replace(/_/g, ' '));
      url = `https://vibemedia.space/equipment/acc_${t}_123.png?prompt=pixel%20${encodedType}%20sprite%20item%20${itemId || ''}&style=pixel_game_asset&key=NOGON`;
    } else {
      // Final generic fallback - use type in the URL to differentiate
      const encodedType = encodeURIComponent(t.replace(/_/g, ' '));
      url = `https://vibemedia.space/equipment/unk_${t}_789.png?prompt=pixel%20${encodedType}%20sprite%20item%20${itemId || ''}&style=pixel_game_asset&key=NOGON`;
    }
  }

  // Append itemId seed parameter when provided to guarantee uniqueness
  if (itemId) {
    url = `${url}&seed=${itemId}`;
  }

  const testid = `equipment-image-${itemId || t}`;

  return { url, testid };
}

export function getEquipmentCategory(type: string | null | undefined): 'weapon' | 'armor' | 'trinket' | 'unknown' {
  if (!type || typeof type !== 'string') return 'unknown';
  
  const t = type.toLowerCase();
  
  const wpnPatterns = ['wpn', 'sword', 'axe', 'bow', 'staff', 'dagger', 'spear', 'hammer', 'blade'];
  const armorPatterns = ['helmet', 'chest', 'gloves', 'gauntlet', 'greave', 'boots', 'armor', 'shield', 'plate', 'mail', 'leather', 'cloth'];
  const trinketPatterns = ['ring', 'necklace', 'amulet', 'trinket', 'charm', 'pendant'];
  
  if (wpnPatterns.some(pattern => t.includes(pattern))) return 'weapon';
  if (armorPatterns.some(pattern => t.includes(pattern))) return 'armor';
  if (trinketPatterns.some(pattern => t.includes(pattern))) return 'trinket';
  
  return 'unknown';
}
