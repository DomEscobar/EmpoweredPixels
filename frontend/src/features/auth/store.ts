import { defineStore } from "pinia";
import { login, register, TokenResponse } from "./api";

interface AuthState {
  token: string | null;
  refresh: string | null;
  userId: number | null;
  isLoading: boolean;
  error: string | null;
}

export const useAuthStore = defineStore("auth", {
  state: (): AuthState => ({
    token: localStorage.getItem("ep_token"),
    refresh: localStorage.getItem("ep_refresh"),
    userId: localStorage.getItem("ep_user_id") ? Number(localStorage.getItem("ep_user_id")) : null,
    isLoading: false,
    error: null,
  }),
  actions: {
    async login(user: string, password: string) {
      this.isLoading = true;
      this.error = null;
      try {
        const data = await login(user, password);
        this.applyToken(data);
      } catch (error) {
        this.error = error instanceof Error ? error.message : "Login failed";
        throw error;
      } finally {
        this.isLoading = false;
      }
    },
    async register(username: string, email: string, password: string) {
      this.isLoading = true;
      this.error = null;
      try {
        await register(username, email, password);
      } catch (error) {
        this.error = error instanceof Error ? error.message : "Register failed";
        throw error;
      } finally {
        this.isLoading = false;
      }
    },
    logout() {
      this.token = null;
      this.refresh = null;
      this.userId = null;
      localStorage.removeItem("ep_token");
      localStorage.removeItem("ep_refresh");
      localStorage.removeItem("ep_user_id");
    },
    applyToken(data: TokenResponse) {
      this.token = data.token;
      this.refresh = data.refresh;
      this.userId = data.userId;
      localStorage.setItem("ep_token", data.token);
      localStorage.setItem("ep_refresh", data.refresh);
      localStorage.setItem("ep_user_id", String(data.userId));
    },
  },
});
