import { defineStore } from "pinia";
import {
  getCurrentEvents,
  getEventStatus,
  getNextEvent,
  type ActiveEvent,
  type EventStatus,
  type NextEventInfo,
} from "./api";
import { useAuthStore } from "@/features/auth/store";

interface EventsState {
  currentEvents: ActiveEvent[];
  eventStatus: EventStatus | null;
  nextEvent: NextEventInfo | null;
  isLoading: boolean;
  error: string | null;
}

export const useEventsStore = defineStore("events", {
  state: (): EventsState => ({
    currentEvents: [],
    eventStatus: null,
    nextEvent: null,
    isLoading: false,
    error: null,
  }),

  getters: {
    hasActiveEvent: (state) => state.eventStatus?.has_active_event ?? false,
    activeMultiplier: (state) => state.eventStatus?.multiplier ?? 1.0,
    activeEventType: (state) => state.eventStatus?.type,
    timeRemaining: (state) => state.eventStatus?.time_remaining,
  },

  actions: {
    async fetchCurrentEvents() {
      const auth = useAuthStore();
      if (!auth.token) return;

      this.isLoading = true;
      try {
        this.currentEvents = await getCurrentEvents(auth.token);
      } catch (e: any) {
        this.error = e.message;
      } finally {
        this.isLoading = false;
      }
    },

    async fetchEventStatus() {
      const auth = useAuthStore();
      if (!auth.token) return;

      try {
        this.eventStatus = await getEventStatus(auth.token);
      } catch (e: any) {
        console.error("Failed to fetch event status", e);
      }
    },

    async fetchNextEvent() {
      const auth = useAuthStore();
      if (!auth.token) return;

      try {
        this.nextEvent = await getNextEvent(auth.token);
      } catch (e: any) {
        console.error("Failed to fetch next event", e);
      }
    },

    async refreshAll() {
      await Promise.all([
        this.fetchCurrentEvents(),
        this.fetchEventStatus(),
        this.fetchNextEvent(),
      ]);
    },
  },
});
