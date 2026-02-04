import { defineStore } from "pinia";
import { Match, getCurrentMatch, getMatches } from "./api";
import { useAuthStore } from "@/features/auth/store";

interface MatchesState {
  currentMatch: Match | null;
  recentMatches: Match[];
  isLoading: boolean;
  error: string | null;
}

export const useMatchesStore = defineStore("matches", {
  state: (): MatchesState => ({
    currentMatch: null,
    recentMatches: [],
    isLoading: false,
    error: null,
  }),
  actions: {
    async fetchCurrentMatch() {
      const auth = useAuthStore();
      if (!auth.token) return;

      this.isLoading = true;
      try {
        const match = await getCurrentMatch(auth.token);
        this.currentMatch = match ?? null;
      } catch (e) {
        this.currentMatch = null;
      } finally {
        this.isLoading = false;
      }
    },
    async fetchRecentMatches() {
      const auth = useAuthStore();
      if (!auth.token) return;

      this.isLoading = true;
      try {
        // Fetch completed matches for history
        const response = await getMatches(auth.token, 1, 5, "completed");
        this.recentMatches = response.items || [];
      } catch (e) {
        this.error = "Failed to load match history";
      } finally {
        this.isLoading = false;
      }
    }
  }
});
