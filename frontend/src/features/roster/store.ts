import { defineStore } from "pinia";
import { Fighter, getFighters, createFighter, deleteFighter, getFighterEquipment, Equipment, updateFighterConfiguration } from "./api";
import { useAuthStore } from "@/features/auth/store";

interface RosterState {
  fighters: Fighter[];
  equipment: Record<string, Equipment[]>;
  isLoading: boolean;
  error: string | null;
}

export const useRosterStore = defineStore("roster", {
  state: (): RosterState => ({
    fighters: [],
    equipment: {},
    isLoading: false,
    error: null,
  }),
  actions: {
    async fetchFighters() {
      const auth = useAuthStore();
      if (!auth.token) return;

      this.isLoading = true;
      try {
        this.fighters = await getFighters(auth.token);
      } catch (e) {
        this.error = "Failed to load roster";
      } finally {
        this.isLoading = false;
      }
    },
    async fetchFighterEquipment(fighterId: string) {
      const auth = useAuthStore();
      if (!auth.token) return;

      try {
        const items = await getFighterEquipment(auth.token, fighterId);
        this.equipment[fighterId] = items;
      } catch (e) {
        console.error("Failed to load equipment", e);
      }
    },
    async addFighter(name: string): Promise<Fighter | undefined> {
      const auth = useAuthStore();
      if (!auth.token) return;

      this.isLoading = true;
      try {
        const newFighter = await createFighter(auth.token, name);
        this.fighters.push(newFighter);
        return newFighter;
      } catch (e) {
        this.error = "Failed to create fighter";
        throw e;
      } finally {
        this.isLoading = false;
      }
    },
    async removeFighter(id: string) {
      const auth = useAuthStore();
      if (!auth.token) return;

      try {
        await deleteFighter(auth.token, id);
        this.fighters = this.fighters.filter(f => f.id !== id);
        delete this.equipment[id];
      } catch (e) {
        this.error = "Failed to delete fighter";
      }
    },
    async updateAttunement(fighterId: string, attunementId: string | null) {
      const auth = useAuthStore();
      if (!auth.token) return;

      try {
        await updateFighterConfiguration(auth.token, fighterId, attunementId);
        const fighter = this.fighters.find(f => f.id === fighterId);
        if (fighter) {
          fighter.attunementId = attunementId ?? undefined;
        }
      } catch (e) {
        this.error = "Failed to update attunement";
        console.error("Failed to update attunement", e);
      }
    }
  }
});
