import { endpoints } from "@/shared/api/endpoints";
import { request } from "@/shared/api/http";

export interface Fighter {
  id: string;
  name: string;
  level: number;
  currentExp: number;
  levelExp: number;
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
  attunementId?: string;
  class?: string;
  created: string;
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

export interface FighterConfiguration {
  fighterId: string;
  attunementId: string | null;
}

export async function updateFighterConfiguration(token: string, fighterId: string, attunementId: string | null) {
  return request<FighterConfiguration>(`${endpoints.fighter}/${fighterId}/configuration`, {
    method: "POST",
    token,
    body: { fighterId, attunementId },
  });
}

export interface Equipment {
  id: string;
  type: string;
  level: number;
  rarity: number;
  enhancement: number;
  fighterId?: string;
}

export async function getFighterEquipment(token: string, fighterId: string) {
  return request<Equipment[]>(`${endpoints.equipment}/fighter/${fighterId}`, { token });
}
