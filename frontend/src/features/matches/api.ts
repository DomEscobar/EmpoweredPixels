import { endpoints } from "@/shared/api/endpoints";
import { request } from "@/shared/api/http";

export interface MatchRegistration {
  matchId: string;
  fighterId: string;
  teamId?: string;
}

export interface MatchOptions {
  isPrivate: boolean;
  maxPowerlevel?: number;
  actionsPerRound: number;
  maxFightersPerUser?: number;
  botCount?: number;
  botPowerlevel?: number;
  features: string[];
  battlefield: string;
  bounds: string;
  positionGenerator: string;
  moveOrder: string;
  winCondition: string;
  staleCondition: string;
}

export interface Match {
  id: string;
  creatorUserId?: number;
  created: string;
  started?: string;
  completedAt?: string;
  cancelledAt?: string;
  status: string;
  ended: boolean;
  registrations: MatchRegistration[];
  options: MatchOptions;
}

export interface PagedResponse<T> {
  page: number;
  pageSize: number;
  totalCount: number;
  items: T[];
}

export async function getMatches(token: string, page: number = 1, pageSize: number = 20, status?: string) {
  return request<PagedResponse<Match>>(`${endpoints.match}/browse`, {
    method: "POST",
    token,
    body: { page, pageSize, status }
  });
}

export async function getCurrentMatch(token: string) {
  try {
    return await request<Match>(`${endpoints.match}/current`, { token });
  } catch (error) {
    // 204 No Content might throw or return null depending on http client
    return null; 
  }
}
