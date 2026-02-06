<template>
  <div v-if="eventsStore.eventStatus?.has_active_event" class="pixel-box bg-gradient-to-r from-purple-900/50 to-amber-900/50 border-amber-500/50 p-4 mb-6">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-4">
        <div class="text-4xl animate-pulse">
          {{ EVENT_TYPE_LABELS[eventsStore.eventStatus.type || '']?.icon || "🎉" }}
        </div>
        <div>
          <h3 class="text-lg font-black text-amber-400 uppercase tracking-wider">
            {{ EVENT_TYPE_LABELS[eventsStore.eventStatus.type || '']?.name || "Special Event" }} Active!
          </h3>
          <p class="text-sm text-slate-300">
            {{ eventsStore.eventStatus.multiplier }}x rewards for a limited time!
          </p>
        </div>
      </div>
      <div class="text-right">
        <div class="text-xs text-slate-400 uppercase">Ends in</div>
        <div class="text-2xl font-mono font-bold text-amber-400">
          {{ eventsStore.eventStatus.time_remaining }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useEventsStore } from '@/features/events/store';
import { EVENT_TYPE_LABELS } from '@/features/events/api';

const eventsStore = useEventsStore();

onMounted(() => {
  eventsStore.fetchStatus();
});
</script>
