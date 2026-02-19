import { endpoints } from "@/shared/api/endpoints";
import { request } from "@/shared/api/http";

export type MatchStatus = 'lobby' | 'running' | 'completed' | 'cancelled';

export interface MatchRegistration {
  matchId: string;
  fighterId: string;
  fighterName?: string;
  fighterLevel?: number;
  fighterPower?: number;
  teamId?: string;
  wins?: number;
  losses?: number;
}

export interface MatchOptions {
  isPrivate: boolean;
  maxPowerlevel?: number;
  maxFightersPerUser?: number;
  botCount?: number;
  botPowerlevel?: number;
  autoStart?: boolean;
}

export interface MatchRewards {
  gold: number;
  experience: number;
  items?: string[];
}

export interface Match {
  id: string;
  creatorUserId?: number;
  created: string;
  started?: string;
  completedAt?: string;
  cancelledAt?: string;
  status: MatchStatus;
  ended: boolean;
  registrations: MatchRegistration[];
  options: MatchOptions;
  rewards?: MatchRewards;
  totalRewards?: number;
  estimatedDuration?: number;
  roundTicks?: any[];
}

export interface MatchWithDetails extends Match {
  myFighterId?: string;
  isHot?: boolean;
  timeAgo?: string;
}

export interface PagedResponse<T> {
  page: number;
  pageSize: number;
  totalCount: number;
  items: T[];
}

export interface MatchCounts {
  lobby: number;
  running: number;
  completed: number;
}

export interface OnlinePlayersResponse {
  onlinePlayers: number;
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
    return null; 
  }
}

export async function getOnlinePlayers(token: string): Promise<OnlinePlayersResponse> {
  return request<OnlinePlayersResponse>(`${endpoints.match}/online-players`, { token });
}

export async function quickJoinMatch(token: string, fighterId: string): Promise<Match> {
  return request<Match>(`${endpoints.match}/quick-join`, {
    method: 'POST',
    token,
    body: { fighterId }
  });
}

export async function joinMatch(token: string, matchId: string, fighterId: string): Promise<void> {
  return request<void>(`${endpoints.match}/join`, {
    method: 'POST',
    token,
    body: { matchId, fighterId }
  });
}

export async function createMatch(token: string, options: Partial<MatchOptions>): Promise<Match> {
  return request<Match>(`${endpoints.match}/create`, {
    method: 'PUT',
    token,
    body: options
  });
}

export async function leaveMatch(token: string, matchId: string, fighterId: string): Promise<void> {
  return request<void>(`${endpoints.match}/leave`, {
    method: 'POST',
    token,
    body: { matchId, fighterId }
  });
}

export async function startMatch(token: string, matchId: string): Promise<void> {
  return request<void>(`${endpoints.match}/${matchId}/start`, {
    method: 'POST',
    token
  });
}

export async function getMatch(token: string, matchId: string): Promise<Match> {
  return request<Match>(`${endpoints.match}/${matchId}`, { token });
}

export async function getMatchRoundTicks(token: string, matchId: string): Promise<any[]> {
  const data = await request<any[]>(`${endpoints.match}/${matchId}/roundticks`, { token });
  return data;
}
