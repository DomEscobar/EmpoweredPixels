<template>
  <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <!-- Backdrop -->
    <div class="absolute inset-0 bg-black/80 backdrop-blur-sm" @click="close"></div>
    
    <!-- Modal -->
    <div class="relative w-full max-w-lg pixel-box bg-slate-900 overflow-hidden">
      <!-- Header -->
      <div class="bg-gradient-to-r from-amber-600 to-amber-500 p-4 text-center border-b-4 border-amber-800">
        <h2 class="text-2xl font-black text-white text-shadow-retro uppercase tracking-wider">
          Daily Reward
        </h2>
        <p class="text-amber-100 text-sm mt-1">Day {{ dailyStore.currentStreak + 1 }} of 7</p>
      </div>

      <!-- Streak Tracker -->
      <div class="p-6 space-y-6">
        <!-- Streak Fire -->
        <div class="flex items-center justify-center gap-2">
          <span v-for="day in 7" :key="day" 
            class="text-2xl transition-all"
            :class="day <= dailyStore.currentStreak ? 'grayscale-0 scale-110' : 'grayscale opacity-30'">
            🔥
          </span>
        </div>

        <!-- Reward Calendar -->
        <div class="grid grid-cols-7 gap-2">
          <div v-for="reward in REWARD_SCHEDULE" :key="reward.day"
            class="pixel-box-sm p-2 text-center transition-all"
            :class="getRewardClass(reward.day)">
            <div class="text-2xl mb-1">{{ reward.icon }}</div>
            <div class="text-xs font-bold truncate">{{ reward.name }}</div>
            <div class="text-[10px] text-slate-400">{{ reward.description }}</div>
          </div>
        </div>

        <!-- Current Reward Highlight -->
        <div v-if="dailyStore.nextReward" class="pixel-box bg-gradient-to-br from-amber-900/50 to-slate-900 p-4 text-center">
          <p class="text-slate-400 text-xs uppercase tracking-wider mb-2">Today's Reward</p>
          <div class="text-5xl mb-2 animate-bounce">{{ dailyStore.nextReward.icon }}</div>
          <h3 class="text-xl font-black text-amber-400">{{ dailyStore.nextReward.name }}</h3>
          <p class="text-slate-300 text-sm">{{ dailyStore.nextReward.description }}</p>
        </div>

        <!-- Time Until Reset -->
        <div v-if="dailyStore.timeUntilReset && !dailyStore.canClaim" class="text-center">
          <p class="text-slate-400 text-sm">Next reward in</p>
          <p class="text-2xl font-mono text-amber-400 font-bold">{{ dailyStore.timeUntilReset }}</p>
        </div>

        <!-- Warning -->
        <div v-if="dailyStore.currentStreak > 0" class="pixel-box-sm bg-red-900/30 border-red-500/30 p-3 text-center">
          <p class="text-red-400 text-xs">
            ⚠️ Miss a day and your streak resets to 0!
          </p>
        </div>
      </div>

      <!-- Actions -->
      <div class="p-4 border-t border-slate-800 flex gap-3">
        <button @click="close" class="flex-1 rpg-btn-secondary">
          Close
        </button>
        <button 
          v-if="dailyStore.canClaim"
          @click="handleClaim"
          :disabled="dailyStore.isClaiming"
          class="flex-1 rpg-btn bg-amber-600 border-amber-800 hover:bg-amber-500 text-white font-bold">
          <span v-if="dailyStore.isClaiming">Claiming...</span>
          <span v-else>CLAIM REWARD 🎁</span>
        </button>
        <button 
          v-else
          disabled
          class="flex-1 pixel-box-sm bg-slate-800 text-slate-500 cursor-not-allowed">
          Already Claimed ✓
        </button>
      </div>

      <!-- Close X -->
      <button @click="close" class="absolute top-2 right-2 text-white/50 hover:text-white">
        ✕
      </button>
    </div>
  </div>

  <!-- Claim Success Animation -->
  <Teleport to="body">
    <div v-if="showSuccess" class="fixed inset-0 z-[60] pointer-events-none flex items-center justify-center">
      <div class="text-center animate-bounce">
        <div class="text-8xl mb-4">{{ claimedReward?.icon }}</div>
        <h3 class="text-4xl font-black text-amber-400 text-shadow-retro">{{ claimedReward?.name }}</h3>
        <p class="text-2xl text-white mt-2">Claimed!</p>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { useDailyStore } from "../store";
import { REWARD_SCHEDULE, type DailyReward } from "../api";

const props = defineProps<{
  show: boolean;
}>();

const emit = defineEmits<{ close: [] }>();

const dailyStore = useDailyStore();
const showSuccess = ref(false);
const claimedReward = ref<DailyReward | null>(null);

function close() {
  emit("close");
}

function getRewardClass(day: number) {
  const currentDay = dailyStore.currentStreak + 1;
  if (day < currentDay) {
    return "bg-green-900/30 border-green-500/30 opacity-70"; // Claimed
  } else if (day === currentDay && dailyStore.canClaim) {
    return "bg-amber-900/50 border-amber-500/50 ring-2 ring-amber-400 animate-pulse"; // Ready
  } else if (day === currentDay) {
    return "bg-slate-800 border-slate-600"; // Current but claimed
  } else {
    return "bg-slate-900/50 border-slate-800 opacity-50"; // Future
  }
}

async function handleClaim() {
  try {
    const result = await dailyStore.claimReward();
    if (result) {
      claimedReward.value = result.reward;
      showSuccess.value = true;
      setTimeout(() => {
        showSuccess.value = false;
      }, 2000);
    }
  } catch (e) {
    // Error handled by store
  }
}

// Auto-fetch when modal opens
watch(() => props.show, (isOpen) => {
  if (isOpen) {
    dailyStore.fetchStatus();
  }
});
</script>
