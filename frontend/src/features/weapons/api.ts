import { endpoints } from "@/shared/api/endpoints";
import { request } from "@/shared/api/http";

export interface Weapon {
  id: string; // The unique user weapon instance ID
  weaponId: string; // The weapon definition ID
  name: string;
  type: string;
  rarity: string;
  enhancement: number;
  durability: number;
  isEquipped: boolean;
  fighterId?: string;
  damage: number;
  attackSpeed: number;
  critChance: number;
  iconUrl: string;
  description: string;
}

export async function getWeapons(token: string) {
  return request<Weapon[]>(endpoints.weapons, { token });
}

export async function equipWeapon(token: string, weaponId: string, fighterId: string) {
  return request<void>(`${endpoints.weapons}/equip`, {
    method: "POST",
    token,
    body: { weaponId, fighterId },
  });
}

export async function unequipWeapon(token: string, weaponId: string) {
  return request<void>(`${endpoints.weapons}/${weaponId}/unequip`, {
    method: "POST",
    token,
  });
}
