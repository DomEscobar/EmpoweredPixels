import { endpoints } from "@/shared/api/endpoints";
import { request } from "@/shared/api/http";

export interface LeagueOptions {
  description?: string;
  matchInterval?: number;
  prizePool?: string;
  tier?: string;
}

export interface League {
  id: number;
  name: string;
  options: LeagueOptions | null;
}

export interface LeagueDetail extends League {
  subscriptions: LeagueSubscription[];
}

export interface LeagueSubscription {
  leagueId: number;
  fighterId: string;
}

export interface LeagueMatch {
  leagueId: number;
  matchId: string;
}

export interface LeagueHighscore {
  fighterId: string;
  fighterName: string;
  username: string;
  score: number;
}

export interface PagedResult<T> {
  page: number;
  pageSize: number;
  totalCount: number;
  items: T[];
}

export async function getLeagues(token: string) {
  return request<League[]>(endpoints.league, { token });
}

export async function getLeagueDetail(token: string, leagueId: number) {
  return request<LeagueDetail>(`${endpoints.league}/${leagueId}`, { token });
}

export async function getUserSubscriptions(token: string, leagueId: number) {
  return request<LeagueSubscription[]>(`${endpoints.league}/${leagueId}/subscriptions/user`, { token });
}

export async function getAllSubscriptions(token: string, leagueId: number) {
  return request<LeagueSubscription[]>(`${endpoints.league}/${leagueId}/subscriptions`, { token });
}

export async function subscribeToLeague(token: string, leagueId: number, fighterId: string) {
  return request<void>(`${endpoints.league}/subscribe`, {
    method: "POST",
    token,
    body: { leagueId, fighterId },
  });
}

export async function unsubscribeFromLeague(token: string, leagueId: number, fighterId: string) {
  return request<void>(`${endpoints.league}/unsubscribe`, {
    method: "POST",
    token,
    body: { leagueId, fighterId },
  });
}

export async function getLeagueMatches(token: string, leagueId: number, page = 1, pageSize = 20) {
  return request<PagedResult<LeagueMatch>>(`${endpoints.league}/${leagueId}/matches`, {
    method: "POST",
    token,
    body: { page, pageSize },
  });
}

export async function getLeagueHighscores(token: string, leagueId: number, lastMatches = 50) {
  return request<LeagueHighscore[]>(`${endpoints.league}/${leagueId}/highscores`, {
    method: "POST",
    token,
    body: { lastMatches },
  });
}
