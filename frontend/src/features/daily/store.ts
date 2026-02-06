import { defineStore } from "pinia";
import { getDailyRewardStatus, claimDailyReward, DailyRewardStatus, ClaimResult } from "./api";
import { useAuthStore } from "@/features/auth/store";

interface DailyState {
  status: DailyRewardStatus | null;
  isLoading: boolean;
  isClaiming: boolean;
  error: string | null;
}

export const useDailyStore = defineStore("daily", {
  state: (): DailyState => ({
    status: null,
    isLoading: false,
    isClaiming: false,
    error: null,
  }),
  
  getters: {
    canClaim: (state): boolean => state.status?.can_claim ?? false,
    currentStreak: (state): number => state.status?.streak ?? 0,
    nextReward: (state) => state.status?.next_reward ?? null,
    timeUntilReset: (state): string => state.status?.time_until_reset ?? "",
    isStreakBroken: (state): boolean => {
      if (!state.status?.last_claimed) return false;
      // Simple check: if can't claim and streak > 0, we already claimed today
      return !state.status.can_claim && state.status.streak > 0;
    }
  },

  actions: {
    async fetchStatus() {
      const auth = useAuthStore();
      if (!auth.token) return;

      this.isLoading = true;
      this.error = null;
      try {
        this.status = await getDailyRewardStatus(auth.token);
      } catch (e: any) {
        this.error = e.message || "Failed to load daily reward";
      } finally {
        this.isLoading = false;
      }
    },

    async claimReward(): Promise<ClaimResult | null> {
      const auth = useAuthStore();
      if (!auth.token) return null;

      this.isClaiming = true;
      this.error = null;
      try {
        const result = await claimDailyReward(auth.token);
        // Refresh status after claim
        await this.fetchStatus();
        return result;
      } catch (e: any) {
        this.error = e.message || "Failed to claim reward";
        throw e;
      } finally {
        this.isClaiming = false;
      }
    },

    clearError() {
      this.error = null;
    }
  }
});
