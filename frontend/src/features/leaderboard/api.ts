// Leaderboard API
const API_URL = import.meta.env.VITE_API_URL || "";

export interface LeaderboardEntry {
  id: string;
  category: string;
  user_id: number;
  username: string;
  avatar?: string;
  rank: number;
  score: number;
  previous_rank: number;
  trend: "up" | "down" | "same";
  updated_at: string;
}

export interface LeaderboardResponse {
  category: string;
  total_count: number;
  user_rank: number;
  user_entry?: LeaderboardEntry;
  entries: LeaderboardEntry[];
}

export interface Achievement {
  id: string;
  key: string;
  name: string;
  description: string;
  icon: string;
  category: string;
  requirement_type: string;
  requirement_value: number;
  reward_gold: number;
  reward_title?: string;
  hidden: boolean;
}

export interface PlayerAchievement {
  id: string;
  user_id: number;
  achievement_id: string;
  achievement?: Achievement;
  progress: number;
  completed: boolean;
  completed_at?: string;
  claimed: boolean;
  claimed_at?: string;
}

export type LeaderboardCategory = "power" | "wealth" | "combat" | "achievements" | "streak";

export const CATEGORY_LABELS: Record<LeaderboardCategory, { name: string; icon: string; description: string }> = {
  power: { name: "Power Ranking", icon: "⚔️", description: "Total fighter power" },
  wealth: { name: "Wealth Ranking", icon: "💰", description: "Gold accumulated" },
  combat: { name: "Combat Ranking", icon: "🏆", description: "Matches won" },
  achievements: { name: "Achievement Points", icon: "⭐", description: "Total achievement score" },
  streak: { name: "Win Streak", icon: "🔥", description: "Current consecutive wins" },
};

export async function getLeaderboard(
  token: string,
  category: LeaderboardCategory,
  limit = 10,
  offset = 0
): Promise<LeaderboardResponse> {
  const response = await fetch(
    `${API_URL}/api/leaderboard/${category}?limit=${limit}&offset=${offset}`,
    { headers: { Authorization: `Bearer ${token}` } }
  );
  if (!response.ok) throw new Error("Failed to fetch leaderboard");
  return response.json();
}

export async function getNearbyRanks(
  token: string,
  category: LeaderboardCategory,
  range_size = 5
): Promise<LeaderboardResponse> {
  const response = await fetch(
    `${API_URL}/api/leaderboard/${category}/nearby?range=${range_size}`,
    { headers: { Authorization: `Bearer ${token}` } }
  );
  if (!response.ok) throw new Error("Failed to fetch nearby ranks");
  return response.json();
}

export async function getAchievements(token: string): Promise<Achievement[]> {
  const response = await fetch(`${API_URL}/api/achievements`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!response.ok) throw new Error("Failed to fetch achievements");
  return response.json();
}

export async function getPlayerAchievements(token: string): Promise<PlayerAchievement[]> {
  const response = await fetch(`${API_URL}/api/player/achievements`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!response.ok) throw new Error("Failed to fetch player achievements");
  return response.json();
}

export async function claimAchievement(token: string, achievementId: string): Promise<{ success: boolean }> {
  const response = await fetch(`${API_URL}/api/achievement/${achievementId}/claim`, {
    method: "POST",
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!response.ok) throw new Error("Failed to claim achievement");
  return response.json();
}
