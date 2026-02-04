import { endpoints } from "@/shared/api/endpoints";
import { request } from "@/shared/api/http";

export interface Reward {
  id: string;
  poolId: string;
}

export interface ItemDto {
  id: string;
  itemId: string;
  rarity: number;
}

export interface EquipmentDto {
  id: string;
  type: string;
  userId: number;
  fighterId?: string;
  isFavorite: boolean;
  level: number;
  rarity: number;
  enhancement: number;
}

export interface RewardContent {
  items: ItemDto[];
  equipment: EquipmentDto[];
}

export async function getRewards(token: string) {
  return request<Reward[]>(endpoints.reward, { token });
}

export async function claimReward(token: string, id: string, poolId: string) {
  return request<RewardContent>(`${endpoints.reward}/claim`, {
    method: "POST",
    token,
    body: { id, poolId }
  });
}

export async function claimAllRewards(token: string) {
  return request<RewardContent>(`${endpoints.reward}/claim/all`, {
    method: "POST",
    token,
    body: {}
  });
}
