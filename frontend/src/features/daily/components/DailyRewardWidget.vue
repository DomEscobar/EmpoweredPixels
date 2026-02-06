<template>
  <div 
    @click="openModal"
    class="pixel-box p-4 cursor-pointer transition-all hover:border-amber-500/50 group"
    :class="dailyStore.canClaim ? 'bg-amber-900/20 border-amber-500/50' : 'bg-slate-900/80'"
  >
    <div class="flex items-center gap-4">
      <!-- Icon -->
      <div 
        class="w-14 h-14 pixel-box-sm flex items-center justify-center text-3xl transition-all"
        :class="dailyStore.canClaim ? 'bg-amber-600 animate-pulse' : 'bg-slate-800'"
      >
        {{ dailyStore.nextReward?.icon || "🎁" }}
      </div>
      
      <!-- Content -->
      <div class="flex-1">
        <h3 class="font-bold text-white group-hover:text-amber-400 transition-colors">
          Daily Reward
        </h3>
        <p class="text-sm text-slate-400">
          <span v-if="dailyStore.canClaim" class="text-amber-400 font-bold">
            Ready to claim!
          </span>
          <span v-else-if="dailyStore.timeUntilReset">
            Next in {{ dailyStore.timeUntilReset }}
          </span>
          <span v-else>
            Streak: {{ dailyStore.currentStreak }} days 🔥
          </span>
        </p>
      </div>

      <!-- Arrow -->
      <div class="text-slate-600 group-hover:text-amber-400 transition-colors">
        →
      </div>
    </div>

    <!-- Progress Bar -->
    <div class="mt-3 h-2 bg-slate-800 rounded-full overflow-hidden">
      <div 
        class="h-full bg-gradient-to-r from-amber-600 to-amber-400 transition-all"
        :style="{ width: `${(dailyStore.currentStreak / 7) * 100}%` }"
      />
    </div>
    <p class="text-xs text-slate-500 mt-1 text-right">
      Day {{ dailyStore.currentStreak }} / 7
    </p>
  </div>

  <!-- Modal -->
  <DailyRewardModal :show="showModal" @close="showModal = false" />
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useDailyStore } from "../store";
import DailyRewardModal from "./DailyRewardModal.vue";

const dailyStore = useDailyStore();
const showModal = ref(false);

function openModal() {
  showModal.value = true;
}

onMounted(() => {
  dailyStore.fetchStatus();
});
</script>
