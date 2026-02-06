import { defineStore } from "pinia";
import { getEventStatus, type EventStatus } from "./api";
import { useAuthStore } from "@/features/auth/store";

interface EventsState {
  eventStatus: EventStatus | null;
  isLoading: boolean;
  error: string | null;
}

export const useEventsStore = defineStore("events", {
  state: (): EventsState => ({
    eventStatus: null,
    isLoading: false,
    error: null,
  }),

  actions: {
    async fetchStatus() {
      const auth = useAuthStore();
      if (!auth.token) return;

      this.isLoading = true;
      try {
        this.eventStatus = await getEventStatus(auth.token);
      } catch (e: any) {
        this.error = e.message;
      } finally {
        this.isLoading = false;
      }
    },
  },
});
