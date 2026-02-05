export interface Attunement {
  element: string;
  level: number;
  current_xp: number;
  total_xp: number;
}

export interface PlayerAttunements {
  user_id: number;
  attunements: Attunement[];
  total_level: number;
}

export interface AttunementBonus {
  level: number;
  power: number;
  defense: number;
  speed: number;
  precision: number;
}

export interface AttunementWithBonuses {
  element: string;
  level: number;
  current_xp: number;
  total_xp: number;
  bonuses: AttunementBonus;
  xp_required: number;
  progress: number;
}

export interface AggregatedBonuses {
  by_element: Record<string, AttunementBonus>;
  total_power: number;
  total_defense: number;
  total_speed: number;
  total_precision: number;
}
