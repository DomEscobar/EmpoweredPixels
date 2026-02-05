# TOOLS.md - Local Notes

Skills define _how_ tools work. This file is for _your_ specifics — the stuff that's unique to your setup.

## What Goes Here

Things like:

- Camera names and locations
- SSH hosts and aliases
- Preferred voices for TTS
- Speaker/room names
- Device nicknames
- Anything environment-specific

## Examples

```markdown
### Cameras

- living-room → Main area, 180° wide angle
- front-door → Entrance, motion-triggered

### SSH

- home-server → 192.168.1.100, user: admin

### TTS

- Preferred voice: "Nova" (warm, slightly British)
- Default speaker: Kitchen HomePod
```

## Why Separate?

Skills are shared. Your setup is yours. Keeping them apart means you can update skills without losing your notes, and share skills without leaking your infrastructure.

---

## Asset Generation - Vibemedia.space API

For Sprites, Icons, and Images, use the vibemedia.space API which creates pixel art on the fly.

### API Format
```
https://vibemedia.space/[FIXED_ID].png?prompt=[URL_ENCODED_PROMPT]&style=pixel_game_asset&key=NOGON
```

**Important:** Use a FIXED complex UUID-like string for [FIXED_ID], not a randomizer function.

### Examples

**Weapons:**
```
vibemedia.space/wpn_fire_sword_9a8f7b2c.png?prompt=legendary%20fire%20sword%20with%20flaming%20blade%20and%20ornate%20hilt&style=pixel_game_asset&key=NOGON
vibemedia.space/wpn_ice_staff_3k8m9n1p.png?prompt=crystal%20ice%20staff%20with%20glowing%20blue%20tip%20and%20frosty%20aura&style=pixel_game_asset&key=NOGON
```

**Character Spritesheets (4x4 grid, 192x192 total, 48x48 per frame):**
```
vibemedia.space/char_knight_walk_012jasd.png?prompt=fantasy%20knight%20character%20spritesheet%20walking%20with%20silver%20armor%20and%20blue%20cape&style=pixel_game_asset&key=NOGON
```
- Row 1: DOWN (walking cycle)
- Row 2: RIGHT (walking cycle)
- Row 3: UP (walking cycle)
- Row 4: LEFT (walking cycle)

**Effects & VFX (4x4 spritesheet):**
```
vibemedia.space/fx_fire_burst_8d7e6f5a.png?prompt=magical%20fire%20effect%20spritesheet%20with%20orange%20and%20red%20flames&style=pixel_game_asset&key=NOGON
```

**Item Icons:**
```
vibemedia.space/icon_health_pot_4b5c6d7e.png?prompt=health%20potion%20sprite%20with%20red%20liquid%20in%20glass%20bottle&style=pixel_game_asset&key=NOGON
vibemedia.space/icon_gold_coin_1a2b3c4d.png?prompt=golden%20coin%20sprite%20with%20detailed%20engravings&style=pixel_game_asset&key=NOGON
```

**Textures (64x64 default):**
```
vibemedia.space/tex_grass_5e6f7g8h.png?prompt=grass%20texture%20with%20vibrant%20green%20blades&style=pixel_game_asset&key=NOGON
```

### Tips for Pixel Game Assets

**General Sprites:**
- Use "sprite" in prompt for single objects
- Describe materials, colors, and key details clearly

**Character Spritesheets:**
- MUST include "spritesheet" in prompt
- Add "walking animation" for 4-directional movement
- Add "one direction" for single pose

**Effects:**
- MUST include "effect" for VFX
- Use terms like "fire", "explosion", "magic sparkles"

**Textures:**
- MUST include "texture" in prompt
- Describe surface details: "weathered", "polished", "rough"

---

Add whatever helps you do your job. This is your cheat sheet.
