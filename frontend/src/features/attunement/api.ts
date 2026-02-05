import { request } from "@/shared/api/http";
import type { PlayerAttunements, AttunementWithBonuses, AggregatedBonuses } from './types'

const BASE = "/api";

export const attunementApi = {
  async getAttunements(): Promise<PlayerAttunements> {
    return request<PlayerAttunements>(`${BASE}/attunements`);
  },

  async getAttunement(element: string): Promise<AttunementWithBonuses> {
    return request<AttunementWithBonuses>(`${BASE}/attunement/${element}`);
  },

  async awardXP(element: string, source: string): Promise<{ success: boolean; level_up: boolean; new_level: number; xp_awarded: number }> {
    return request(`${BASE}/attunement/award-xp`, {
      method: "POST",
      body: { element, source }
    });
  },

  async getBonuses(): Promise<AggregatedBonuses> {
    return request<AggregatedBonuses>(`${BASE}/attunements/bonuses`);
  }
};

// Element display configuration
export const elementConfig: Record<string, { name: string; icon: string; color: string; description: string }> = {
  fire: { name: 'Fire', icon: '🔥', color: '#ef4444', description: 'Increases attack power' },
  water: { name: 'Water', icon: '💧', color: '#3b82f6', description: 'Improves defense' },
  earth: { name: 'Earth', icon: '🌍', color: '#22c55e', description: 'Maximizes defense' },
  air: { name: 'Air', icon: '💨', color: '#a8a29e', description: 'Boosts speed' },
  light: { name: 'Light', icon: '✨', color: '#fbbf24', description: 'Enhances precision' },
  dark: { name: 'Dark', icon: '🌑', color: '#7c3aed', description: 'Increases power' }
};

// Format XP with commas
export function formatXP(xp: number): string {
  return xp.toLocaleString();
}

// Calculate progress percentage
export function calculateProgress(current: number, required: number): number {
  if (required <= 0) return 100;
  return Math.min(100, Math.round((current / required) * 100));
}
