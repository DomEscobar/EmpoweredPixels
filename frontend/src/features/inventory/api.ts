import { endpoints } from "@/shared/api/endpoints";
import { request } from "@/shared/api/http";

export interface Equipment {
  id: string;
  userId: number;
  fighterId?: string;
  type: string; // Backend sends "type" which is the ItemID
  level: number;
  rarity: number;
  enhancement: number;
  isFavorite?: boolean;
}

export interface InventoryPageResponse {
  page: number;
  pageSize: number;
  totalCount: number;
  items: Equipment[];
}

export interface InventoryPageRequest {
  page: number;
  pageSize: number;
  filter?: {
    rarity?: number[];
    levelMin?: number;
    levelMax?: number;
    slot?: string;
  };
  sort?: {
    field: string;
    direction: 'asc' | 'desc';
  };
}

interface BalanceResponse {
  itemId: string;
  balance: number;
}

interface ItemDto {
  id: string;
  itemId: string;
  rarity: number;
}

export async function getParticleBalance(token: string): Promise<number> {
  const res = await request<BalanceResponse>(`${endpoints.inventory}/balance/particles`, { token });
  return res.balance;
}

export async function getTokenBalance(token: string, type: 'common' | 'rare' | 'fabled' | 'mythic'): Promise<number> {
  const res = await request<BalanceResponse>(`${endpoints.inventory}/balance/token/${type}`, { token });
  return res.balance;
}

export async function getInventoryPage(token: string, req: InventoryPageRequest): Promise<InventoryPageResponse> {
  return request<InventoryPageResponse>(`${endpoints.equipment}/inventory`, {
    method: "POST",
    token,
    body: req
  });
}

export async function getEquipment(token: string, id: string): Promise<Equipment> {
  return request<Equipment>(`${endpoints.equipment}/${id}`, { token });
}

export async function enhanceEquipment(token: string, equipmentId: string): Promise<Equipment> {
  return request<Equipment>(`${endpoints.equipment}/enhance`, {
    method: "POST",
    token,
    body: { equipmentId }
  });
}

export async function getEnhanceCost(token: string, equipmentId: string): Promise<{ particles: number; tokens: number; tokenType: string }> {
  return request<{ particles: number; tokens: number; tokenType: string }>(`${endpoints.equipment}/enhance/cost`, {
    method: "POST",
    token,
    body: { equipmentId }
  });
}

export async function salvageEquipment(token: string, equipmentId: string): Promise<{ particles: number }> {
  const items = await request<ItemDto[]>(`${endpoints.equipment}/salvage`, {
    method: "POST",
    token,
    body: { equipmentId }
  });
  
  // Count particles (assuming most salvaged items are particles for now, 
  // or we could filter by ID if we had the constant exported to frontend)
  return { particles: items.length };
}

export async function setFavorite(token: string, equipmentId: string, isFavorite: boolean): Promise<void> {
  return request<void>(`${endpoints.equipment}/${equipmentId}/favorite`, {
    method: isFavorite ? "POST" : "DELETE",
    token
  });
}
