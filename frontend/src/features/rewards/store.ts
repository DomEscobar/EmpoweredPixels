import { defineStore } from "pinia";
import { Reward, RewardContent, getRewards, claimReward, claimAllRewards } from "./api";
import { useAuthStore } from "@/features/auth/store";

interface RewardsState {
  rewards: Reward[];
  lastClaimed: RewardContent | null;
  isLoading: boolean;
  error: string | null;
}

export const useRewardsStore = defineStore("rewards", {
  state: (): RewardsState => ({
    rewards: [],
    lastClaimed: null,
    isLoading: false,
    error: null,
  }),
  getters: {
    rewardCount: (state) => state.rewards.length,
  },
  actions: {
    async fetchRewards() {
      const auth = useAuthStore();
      if (!auth.token) return;

      this.isLoading = true;
      try {
        this.rewards = await getRewards(auth.token);
      } catch (e) {
        this.error = "Failed to load rewards";
      } finally {
        this.isLoading = false;
      }
    },
    async claim(reward: Reward) {
      const auth = useAuthStore();
      if (!auth.token) return;

      this.isLoading = true;
      try {
        const content = await claimReward(auth.token, reward.id, reward.poolId);
        this.lastClaimed = content;
        // Remove from list
        this.rewards = this.rewards.filter(r => r.id !== reward.id);
        return content;
      } catch (e) {
        this.error = "Failed to claim reward";
        throw e;
      } finally {
        this.isLoading = false;
      }
    },
    async claimAll() {
      const auth = useAuthStore();
      if (!auth.token) return;

      this.isLoading = true;
      try {
        const content = await claimAllRewards(auth.token);
        this.lastClaimed = content;
        this.rewards = [];
        return content;
      } catch (e) {
        this.error = "Failed to claim rewards";
        throw e;
      } finally {
        this.isLoading = false;
      }
    }
  }
});
