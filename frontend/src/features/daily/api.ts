// Daily Rewards API
const API_URL = import.meta.env.VITE_API_URL || "";

export interface DailyReward {
  day: number;
  name: string;
  description: string;
  icon: string;
  type: string;
  value?: number;
  rarity?: number;
}

export interface DailyRewardStatus {
  user_id: number;
  streak: number;
  last_claimed: string | null;
  total_claimed: number;
  can_claim: boolean;
  next_reward: DailyReward;
  time_until_reset?: string;
}

export interface ClaimResult {
  success: boolean;
  reward: DailyReward;
  reward_value: number;
  new_streak: number;
  day: number;
  next_reward: DailyReward;
}

export async function getDailyRewardStatus(token: string): Promise<DailyRewardStatus> {
  const response = await fetch(`${API_URL}/api/daily-reward`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!response.ok) throw new Error("Failed to fetch daily reward status");
  return response.json();
}

export async function claimDailyReward(token: string): Promise<ClaimResult> {
  const response = await fetch(`${API_URL}/api/daily-reward/claim`, {
    method: "POST",
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || "Failed to claim reward");
  }
  return response.json();
}

// Reward schedule for display
export const REWARD_SCHEDULE: DailyReward[] = [
  { day: 1, name: "Small Pouch", description: "100 Gold", icon: "🪙", type: "gold", value: 100 },
  { day: 2, name: "Common Chest", description: "250 Gold", icon: "📦", type: "gold", value: 250 },
  { day: 3, name: "Rare Cache", description: "500 Gold", icon: "💎", type: "gold", value: 500 },
  { day: 4, name: "Energy Boost", description: "2x XP for 1h", icon: "⚡", type: "boost", value: 2 },
  { day: 5, name: "Mystery Box", description: "Random Reward", icon: "🎁", type: "mystery" },
  { day: 6, name: "Fabled Vault", description: "1000 Gold", icon: "🏆", type: "gold", value: 1000 },
  { day: 7, name: "Legendary Crate", description: "2000 Gold + Legendary", icon: "👑", type: "gold", value: 2000, rarity: 5 },
];
