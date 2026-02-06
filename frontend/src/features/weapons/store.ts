import { defineStore } from "pinia";
import { Weapon, getWeapons, equipWeapon, unequipWeapon } from "./api";
import { useAuthStore } from "@/features/auth/store";
import { useRosterStore } from "@/features/roster/store";

interface WeaponsState {
  weapons: Weapon[];
  isLoading: boolean;
  error: string | null;
}

export const useWeaponsStore = defineStore("weapons", {
  state: (): WeaponsState => ({
    weapons: [],
    isLoading: false,
    error: null,
  }),
  actions: {
    async fetchWeapons() {
      const auth = useAuthStore();
      if (!auth.token) return;

      this.isLoading = true;
      try {
        this.weapons = await getWeapons(auth.token);
      } catch (e) {
        this.error = "Failed to load weapons";
      } finally {
        this.isLoading = false;
      }
    },
    async equipWeapon(weaponId: string, fighterId: string) {
      const auth = useAuthStore();
      if (!auth.token) return;

      try {
        await equipWeapon(auth.token, weaponId, fighterId);
        // Refresh weapons list to get updated isEquipped/fighterId states
        await this.fetchWeapons();
        // Also refresh the specific fighter in roster if needed, although primary weapon is often just a field
        const roster = useRosterStore();
        const fighter = roster.fighters.find(f => f.id === fighterId);
        if (fighter) {
           fighter.weaponId = weaponId;
        }
      } catch (e) {
        this.error = "Failed to equip weapon";
        throw e;
      }
    },
    async unequipWeapon(weaponId: string) {
      const auth = useAuthStore();
      if (!auth.token) return;

      try {
        await unequipWeapon(auth.token, weaponId);
        
        const weapon = this.weapons.find(w => w.id === weaponId);
        if (weapon && weapon.fighterId) {
            const roster = useRosterStore();
            const fighter = roster.fighters.find(f => f.id === weapon.fighterId);
            if (fighter && fighter.weaponId === weaponId) {
                fighter.weaponId = undefined;
            }
        }
        
        await this.fetchWeapons();
      } catch (e) {
        this.error = "Failed to unequip weapon";
        throw e;
      }
    }
  }
});
