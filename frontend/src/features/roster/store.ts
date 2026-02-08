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
// ... existing code
});

export const useFighterStore = useRosterStore;
