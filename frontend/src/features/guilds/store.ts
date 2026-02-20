import { defineStore } from "pinia";
import { Guild, getGuilds, createGuild, joinGuild } from "./api";
import { useAuthStore } from "@/features/auth/store";

interface GuildsState {
  guilds: Guild[];
  isLoading: boolean;
  error: string | null;
}

export const useGuildsStore = defineStore("guilds", {
  state: (): GuildsState => ({
    guilds: [],
    isLoading: false,
    error: null,
  }),
  actions: {
    async fetchGuilds() {
      const auth = useAuthStore();
      if (!auth.token) return;

      this.isLoading = true;
      try {
        this.guilds = await getGuilds(auth.token);
      } catch (e: any) {
        this.error = e.message || "Failed to load guilds";
      } finally {
        this.isLoading = false;
      }
    },
    async createNewGuild(name: string, description: string) {
      const auth = useAuthStore();
      if (!auth.token) return;

      this.isLoading = true;
      try {
        const guild = await createGuild(auth.token, name, description);
        this.guilds.push(guild);
        return guild;
      } catch (e: any) {
        this.error = e.message || "Failed to create guild";
        throw e;
      } finally {
        this.isLoading = false;
      }
    },
    async requestJoin(guildId: string) {
      const auth = useAuthStore();
      if (!auth.token) return;

      try {
        return await joinGuild(auth.token, guildId);
      } catch (e: any) {
        this.error = e.message || "Failed to join guild";
        throw e;
      }
    }
  }
});
