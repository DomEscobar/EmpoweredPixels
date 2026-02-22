import { defineStore } from "pinia";
import { Fighter, getFighters, createFighter, deleteFighter, getFighterEquipment, Equipment, equipItem, unequipItem } from "./api";
import { useAuthStore } from "@/features/auth/store";
import { useInventoryStore } from "@/features/inventory/store";

interface RosterState {
  fighters: Fighter[];
  equipment: Record<string, Equipment[]>;
  isLoading: boolean;
  error: string | null;
}

export const useFighterStore = defineStore("roster", {
  state: (): RosterState => ({
    fighters: [],
    equipment: {},
    isLoading: false,
    error: null,
  }),
  actions: {
    async fetchFighters() {
      const auth = useAuthStore();
      if (!auth.token) {
        this.error = "Not authenticated";
        return;
      }

      this.isLoading = true;
      try {
        this.fighters = await getFighters(auth.token);
      } catch (e: any) {
        this.error = e.message || "Failed to load roster";
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
    async addFighter(name: string): Promise<Fighter> {
      const auth = useAuthStore();
      if (!auth.token) {
        throw new Error("Not authenticated");
      }

      this.isLoading = true;
      try {
        const newFighter = await createFighter(auth.token, name);
        this.fighters.push(newFighter);
        return newFighter;
      } catch (e: any) {
        this.error = e.message || "Failed to create fighter";
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
    async equipItemToFighter(fighterId: string, equipmentId: string) {
      const auth = useAuthStore();
      const inventory = useInventoryStore();
      if (!auth.token) {
        throw new Error("Not authenticated");
      }

      // Validate equipment exists in inventory before equip
      const item = inventory.equipment.find(i => i.id === equipmentId);
      if (!item) {
        throw new Error("Equipment not found in inventory");
      }
      if (item.fighterId) {
        throw new Error("Equipment is already bound to another fighter");
      }

      this.isLoading = true;
      try {
        await equipItem(auth.token, fighterId, equipmentId);
        await this.fetchFighterEquipment(fighterId);
        await this.fetchFighters(); // Power might change
        await inventory.fetchInventory(); // Item is now bound
      } catch (e: any) {
        const errorMsg = e.message || "Failed to equip item";
        this.error = errorMsg;
        throw new Error(errorMsg);
      } finally {
        this.isLoading = false;
      }
    },
    async unequipItemFromFighter(fighterId: string, equipmentId: string) {
      const auth = useAuthStore();
      const inventory = useInventoryStore();
      if (!auth.token) return;

      try {
        await unequipItem(auth.token, equipmentId);
        await this.fetchFighterEquipment(fighterId);
        await this.fetchFighters(); // Power might change
        await inventory.fetchInventory(); // Item is now free
      } catch (e: any) {
        this.error = e.message || "Failed to unequip item";
        throw e;
      }
    }
  }
});

export const useRosterStore = useFighterStore;
