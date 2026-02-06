<template>
  <div v-if="eventsStore.hasActiveEvent" class="pixel-box bg-gradient-to-r from-purple-900/50 to-amber-900/50 border-amber-500/50 p-4 mb-6">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-4">
        <div class="text-4xl animate-pulse">
          {{ EVENT_TYPE_LABELS[eventsStore.activeEventType || '']?.icon || "🎉" }}
        </div>
        <div>
          <h3 class="text-lg font-black text-amber-400 uppercase tracking-wider">
            {{ EVENT_TYPE_LABELS[eventsStore.activeEventType || '']?.name || "Special Event" }} Active!
          </h3>
          <p class="text-sm text-slate-300">
            {{ eventsStore.activeMultiplier }}x rewards for a limited time!
          </p>
        </div>
      </div>
      <div class="text-right">
        <div class="text-xs text-slate-400 uppercase">Ends in</div>
        <div class="text-2xl font-mono font-bold text-amber-400">
          {{ eventsStore.timeRemaining }}
        </div>
      </div>
    </div>
  </div>

  <!-- Next Event Preview -->
  <div v-else-if="eventsStore.nextEvent" class="pixel-box bg-slate-900/50 border-slate-700/50 p-3 mb-6">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-3">
        <span class="text-2xl">📅</span>
        <div>
          <span class="text-sm text-slate-400">Next Event:</span>
          <span class="text-sm font-bold text-amber-400 ml-2">
            {{ eventsStore.nextEvent.event.name }}
          </span>
        </div>
      </div>
      <div class="text-xs text-slate-500">
        Starts in {{ formatTime(eventsStore.nextEvent.wait_ms) }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useEventsStore } from '@/features/events/store';
import { EVENT_TYPE_LABELS } from '@/features/events/api';

const eventsStore = useEventsStore();

function formatTime(ms: number) {
  const hours = Math.floor(ms / 3600000);
  const minutes = Math.floor((ms % 3600000) / 60000);
  if (hours > 0) return `${hours}h ${minutes}m`;
  return `${minutes}m`;
}

onMounted(() => {
  eventsStore.refreshAll();
  // Refresh every minute
  setInterval(() => eventsStore.refreshAll(), 60000);
});
</script>
