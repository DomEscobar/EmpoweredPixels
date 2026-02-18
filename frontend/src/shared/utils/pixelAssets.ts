export const PIXEL_ASSETS = {
  BG_DUNGEON: 'https://vibemedia.space/bg_dungeon_v2_99283.png?prompt=dark%20dungeon%20floor%20tile%20texture%20seamless&style=pixel_game_asset&key=NOGON',
  BG_DUNGEON_DASH: 'https://vibemedia.space/bg_dungeon_dash_5y6t7u_v1.png?prompt=dark%20dungeon%20stone%20floor%20tile%20texture%20seamless&style=pixel_game_asset&key=NOGON',
  ICON_SWORDS: 'https://vibemedia.space/icon_swords_v2_11234.png?prompt=crossed%20steel%20swords%20pixel%20art&style=pixel_game_asset&key=NOGON',
  ICON_TROPHY: 'https://vibemedia.space/trophy_icon_4d5e6f_v1.png?prompt=golden%20trophy%20pixel%20art&style=pixel_game_asset&key=NOGON',
  ICON_SCROLL: 'https://vibemedia.space/icon_scroll_v2_66789.png?prompt=ancient%20magic%20scroll%20pixel%20art&style=pixel_game_asset&key=NOGON',
  ICON_SCROLL_EMPTY: 'https://vibemedia.space/icon_scroll_blank_6z7a8s_v1.png?prompt=blank%20parchment%20scroll%20pixel%20art&style=pixel_game_asset&key=NOGON',
  ICON_CHEST: 'https://vibemedia.space/icon_chest_v2_55667.png?prompt=ancient%20wooden%20treasure%20chest%20with%20golden%20brass%20corners%20and%20mystical%20glow&style=pixel_game_asset&key=NOGON',
  ICON_POTION: 'https://vibemedia.space/icon_potion_v2_77890.png?prompt=red%20health%20potion%20pixel%20art&style=pixel_game_asset&key=NOGON',
  ICON_BOT: 'https://vibemedia.space/icon_skull_v2_55432.png?prompt=skull%20icon%20pixel%20art&style=pixel_game_asset&key=NOGON',
  ICON_FIGHTER: 'https://vibemedia.space/fighter_hooded_8p9q0r_v1.png?prompt=mystery%20hooded%20figure%20pixel%20art&style=pixel_game_asset&key=NOGON',
  ICON_CASTLE: 'https://vibemedia.space/icon_castle_cmd_8i9o0p_v1.png?prompt=medieval%20castle%20tower%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  ICON_KNIGHT: 'https://vibemedia.space/icon_knight_1a2s3d_v1.png?prompt=armored%20knight%20helmet%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  ICON_FLAG: 'https://vibemedia.space/icon_flag_war_4f5g6h_v1.png?prompt=war%20banner%20flag%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  ICON_CROWN: 'https://vibemedia.space/icon_crown_gold_9d0f1g_v1.png?prompt=golden%20royal%20crown%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  ICON_WARRIOR: 'https://vibemedia.space/icon_warrior_elite_2h3j4k_v1.png?prompt=elite%20warrior%20knight%20pixel%20art%20character&style=pixel_game_asset&key=NOGON',
} as const;

export type PixelAssetKey = keyof typeof PIXEL_ASSETS;

export function getMatchStatusIcon(status: string): string {
  if (status === 'running') return PIXEL_ASSETS.ICON_SWORDS;
  if (status === 'completed') return PIXEL_ASSETS.ICON_TROPHY;
  return PIXEL_ASSETS.ICON_SCROLL;
}

export function getMatchHeroImage(status: string): string {
  if (status === 'running') return PIXEL_ASSETS.ICON_SWORDS;
  if (status === 'completed') return PIXEL_ASSETS.ICON_TROPHY;
  return PIXEL_ASSETS.ICON_CHEST;
}
