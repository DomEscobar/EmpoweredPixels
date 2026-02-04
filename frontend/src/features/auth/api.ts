import { endpoints } from "@/shared/api/endpoints";
import { request } from "@/shared/api/http";

export interface TokenResponse {
  userId: number;
  token: string;
  refresh: string;
}

export async function login(user: string, password: string) {
  return request<TokenResponse>(endpoints.authToken, {
    method: "POST",
    body: { user, password },
  });
}

export async function refreshToken(userId: number, refresh: string) {
  return request<TokenResponse>(endpoints.authRefresh, {
    method: "POST",
    body: { userId, refresh },
  });
}

export async function register(username: string, email: string, password: string) {
  return request<void>(endpoints.register, {
    method: "POST",
    body: { username, email, password },
  });
}
