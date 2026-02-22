import { endpoints } from "@/shared/api/endpoints";
import { request } from "@/shared/api/http";

export interface Fighter {
  id: string;
  name: string;
  level: number;
  currentExp: number;
  levelExp: number;
  xp?: number; // Kept for compatibility if used elsewhere
  xpToNextLevel?: number; // Kept for compatibility
  power: number;
  conditionPower: number;
  precision: number;
  ferocity: number;
  accuracy: number;
  agility: number;
  armor: number;
  vitality: number;
  parryChance: number;
  healingPower: number;
  speed: number;
  vision: number;
  weaponId?: string;
  class?: string;
  created: string;
  // Match Statistics
  matchesWon: number;
  matchesLost: number;
  totalMatches: number;
  totalDamageDealt: number;
  totalDamageTaken: number;
}

export async function getFighters(token: string) {
  return request<Fighter[]>(endpoints.fighter, { token });
}

export async function createFighter(token: string, name: string) {
  return request<Fighter>(endpoints.fighter, {
    method: "PUT",
    token,
    body: { name },
  });
}

export async function deleteFighter(token: string, id: string) {
  return request<void>(`${endpoints.fighter}/${id}`, {
    method: "DELETE",
    token,
  });
}

export interface Equipment {
  id: string;
  type: string;
  category: string;
  level: number;
  rarity: number;
  enhancement: number;
  fighterId?: string;
}

export async function getFighterEquipment(token: string, fighterId: string) {
  return request<Equipment[]>(`${endpoints.equipment}/fighter/${fighterId}`, { token });
}

export async function equipItem(token: string, fighterId: string, equipmentId: string) {
  return request<void>(`${endpoints.equipment}/${equipmentId}/equip/${fighterId}`, {
    method: "POST",
    token
  });
}

export async function unequipItem(token: string, equipmentId: string) {
  return request<void>(`${endpoints.equipment}/${equipmentId}/unequip`, {
    method: "POST",
    token
  });
}
