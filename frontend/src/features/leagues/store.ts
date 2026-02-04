import { defineStore } from "pinia";
import {
  League,
  LeagueDetail,
  LeagueSubscription,
  LeagueMatch,
  LeagueHighscore,
  getLeagues,
  getLeagueDetail,
  getUserSubscriptions,
  getAllSubscriptions,
  subscribeToLeague,
  unsubscribeFromLeague,
  getLeagueMatches,
  getLeagueHighscores,
} from "./api";
import { useAuthStore } from "@/features/auth/store";

interface LeaguesState {
  leagues: League[];
  subscriptions: Record<number, LeagueSubscription[]>;
  allSubscriptions: Record<number, LeagueSubscription[]>;
  leagueMatches: Record<number, LeagueMatch[]>;
  leagueHighscores: Record<number, LeagueHighscore[]>;
  selectedLeague: LeagueDetail | null;
  isLoading: boolean;
  isLoadingDetail: boolean;
  isSubscribing: boolean;
  error: string | null;
}

export const useLeaguesStore = defineStore("leagues", {
  state: (): LeaguesState => ({
    leagues: [],
    subscriptions: {},
    allSubscriptions: {},
    leagueMatches: {},
    leagueHighscores: {},
    selectedLeague: null,
    isLoading: false,
    isLoadingDetail: false,
    isSubscribing: false,
    error: null,
  }),
  actions: {
    async fetchLeagues() {
      const auth = useAuthStore();
      if (!auth.token) return;

      this.isLoading = true;
      this.error = null;
      try {
        this.leagues = await getLeagues(auth.token);
        await Promise.all(this.leagues.map((l) => this.fetchSubscriptions(l.id)));
      } catch (e) {
        this.error = "Failed to load leagues";
      } finally {
        this.isLoading = false;
      }
    },

    async fetchSubscriptions(leagueId: number) {
      const auth = useAuthStore();
      if (!auth.token) return;

      try {
        const subs = await getUserSubscriptions(auth.token, leagueId);
        this.subscriptions[leagueId] = subs;
      } catch (e) {
        console.error(`Failed to load subscriptions for league ${leagueId}`, e);
      }
    },

    async fetchLeagueDetail(leagueId: number) {
      const auth = useAuthStore();
      if (!auth.token) return;

      this.isLoadingDetail = true;
      this.error = null;
      try {
        this.selectedLeague = await getLeagueDetail(auth.token, leagueId);
        this.allSubscriptions[leagueId] = this.selectedLeague.subscriptions;
      } catch (e) {
        this.error = "Failed to load league details";
        this.selectedLeague = null;
      } finally {
        this.isLoadingDetail = false;
      }
    },

    async fetchAllSubscriptions(leagueId: number) {
      const auth = useAuthStore();
      if (!auth.token) return;

      try {
        this.allSubscriptions[leagueId] = await getAllSubscriptions(auth.token, leagueId);
      } catch (e) {
        console.error(`Failed to load all subscriptions for league ${leagueId}`, e);
      }
    },

    async subscribe(leagueId: number, fighterId: string) {
      const auth = useAuthStore();
      if (!auth.token) return;

      this.isSubscribing = true;
      this.error = null;
      try {
        await subscribeToLeague(auth.token, leagueId, fighterId);
        if (!this.subscriptions[leagueId]) {
          this.subscriptions[leagueId] = [];
        }
        this.subscriptions[leagueId].push({ leagueId, fighterId });
        await this.fetchAllSubscriptions(leagueId);
      } catch (e) {
        this.error = "Failed to subscribe";
        throw e;
      } finally {
        this.isSubscribing = false;
      }
    },

    async unsubscribe(leagueId: number, fighterId: string) {
      const auth = useAuthStore();
      if (!auth.token) return;

      this.isSubscribing = true;
      this.error = null;
      try {
        await unsubscribeFromLeague(auth.token, leagueId, fighterId);
        if (this.subscriptions[leagueId]) {
          this.subscriptions[leagueId] = this.subscriptions[leagueId].filter(
            (s) => s.fighterId !== fighterId
          );
        }
        await this.fetchAllSubscriptions(leagueId);
      } catch (e) {
        this.error = "Failed to unsubscribe";
        throw e;
      } finally {
        this.isSubscribing = false;
      }
    },

    async fetchLeagueMatches(leagueId: number, page = 1) {
      const auth = useAuthStore();
      if (!auth.token) return;

      try {
        const result = await getLeagueMatches(auth.token, leagueId, page, 20);
        this.leagueMatches[leagueId] = result.items;
      } catch (e) {
        console.error(`Failed to load matches for league ${leagueId}`, e);
      }
    },

    async fetchLeagueHighscores(leagueId: number) {
      const auth = useAuthStore();
      if (!auth.token) return;

      try {
        this.leagueHighscores[leagueId] = await getLeagueHighscores(auth.token, leagueId);
      } catch (e) {
        console.error(`Failed to load highscores for league ${leagueId}`, e);
      }
    },

    clearSelectedLeague() {
      this.selectedLeague = null;
    },
  },
  getters: {
    activeLeagueCount: (state) => {
      let count = 0;
      for (const leagueId in state.subscriptions) {
        if (state.subscriptions[leagueId].length > 0) count++;
      }
      return count;
    },
    isSubscribedToLeague: (state) => (leagueId: number) => {
      return (state.subscriptions[leagueId]?.length ?? 0) > 0;
    },
    getSubscribedFighters: (state) => (leagueId: number) => {
      return state.subscriptions[leagueId] ?? [];
    },
    getParticipantCount: (state) => (leagueId: number) => {
      return state.allSubscriptions[leagueId]?.length ?? 0;
    },
  },
});
