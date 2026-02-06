import { defineStore } from "pinia";
import {
  getLeaderboard,
  getNearbyRanks,
  getAchievements,
  getPlayerAchievements,
  claimAchievement,
  type LeaderboardResponse,
  type PlayerAchievement,
  type Achievement,
  type LeaderboardCategory,
} from "./api";
import { useAuthStore } from "@/features/auth/store";

interface LeaderboardState {
  currentCategory: LeaderboardCategory;
  leaderboard: LeaderboardResponse | null;
  nearbyRanks: LeaderboardResponse | null;
  achievements: Achievement[];
  playerAchievements: PlayerAchievement[];
  isLoading: boolean;
  isLoadingAchievements: boolean;
  error: string | null;
}

export const useLeaderboardStore = defineStore("leaderboard", {
  state: (): LeaderboardState => ({
    currentCategory: "power",
    leaderboard: null,
    nearbyRanks: null,
    achievements: [],
    playerAchievements: [],
    isLoading: false,
    isLoadingAchievements: false,
    error: null,
  }),

  getters: {
    topEntries: (state) => state.leaderboard?.entries.slice(0, 10) || [],
    userRank: (state) => state.leaderboard?.user_rank || 0,
    userEntry: (state) => state.leaderboard?.user_entry,
    completedAchievements: (state) =>
      state.playerAchievements.filter((pa) => pa.completed),
    unclaimedAchievements: (state) =>
      state.playerAchievements.filter((pa) => pa.completed && !pa.claimed),
    progressByCategory: (state) => {
      const categories: Record<string, { total: number; completed: number }> = {};
      state.playerAchievements.forEach((pa) => {
        const cat = pa.achievement?.category || "other";
        if (!categories[cat]) categories[cat] = { total: 0, completed: 0 };
        categories[cat].total++;
        if (pa.completed) categories[cat].completed++;
      });
      return categories;
    },
  },

  actions: {
    async fetchLeaderboard(category?: LeaderboardCategory, limit = 10, offset = 0) {
      const auth = useAuthStore();
      if (!auth.token) return;

      const cat = category || this.currentCategory;
      this.currentCategory = cat;
      this.isLoading = true;
      this.error = null;

      try {
        this.leaderboard = await getLeaderboard(auth.token, cat, limit, offset);
      } catch (e: any) {
        this.error = e.message || "Failed to load leaderboard";
      } finally {
        this.isLoading = false;
      }
    },

    async fetchNearbyRanks(category?: LeaderboardCategory, rangeSize = 5) {
      const auth = useAuthStore();
      if (!auth.token) return;

      const cat = category || this.currentCategory;
      this.isLoading = true;
      this.error = null;

      try {
        this.nearbyRanks = await getNearbyRanks(auth.token, cat, rangeSize);
      } catch (e: any) {
        this.error = e.message || "Failed to load nearby ranks";
      } finally {
        this.isLoading = false;
      }
    },

    async fetchAchievements() {
      const auth = useAuthStore();
      if (!auth.token) return;

      this.isLoadingAchievements = true;
      try {
        this.achievements = await getAchievements(auth.token);
      } catch (e: any) {
        console.error("Failed to fetch achievements", e);
      } finally {
        this.isLoadingAchievements = false;
      }
    },

    async fetchPlayerAchievements() {
      const auth = useAuthStore();
      if (!auth.token) return;

      this.isLoadingAchievements = true;
      try {
        this.playerAchievements = await getPlayerAchievements(auth.token);
      } catch (e: any) {
        console.error("Failed to fetch player achievements", e);
      } finally {
        this.isLoadingAchievements = false;
      }
    },

    async claimAchievementReward(achievementId: string) {
      const auth = useAuthStore();
      if (!auth.token) return;

      try {
        await claimAchievement(auth.token, achievementId);
        // Refresh player achievements
        await this.fetchPlayerAchievements();
      } catch (e: any) {
        this.error = e.message || "Failed to claim reward";
        throw e;
      }
    },

    clearError() {
      this.error = null;
    },
  },
});
