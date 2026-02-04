export const endpoints = {
  health: "/health",
  authToken: "/api/authentication/token",
  authRefresh: "/api/authentication/refresh",
  register: "/api/register",
  fighter: "/api/fighter",
  match: "/api/match",
  inventory: "/api/inventory",
  equipment: "/api/equipment",
  reward: "/api/reward",
  league: "/api/league",
  season: "/api/season",
} as const;
