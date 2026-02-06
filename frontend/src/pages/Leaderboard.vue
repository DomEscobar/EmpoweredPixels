<template>
  <div class="min-h-screen p-4 md:p-8 font-mono" :style="{ backgroundImage: `url('${PIXEL_ASSETS.BG_DUNGEON}')` }">
    <div class="fixed inset-0 bg-slate-950/80 pointer-events-none z-0"></div>
    
    <div class="relative z-10 max-w-6xl mx-auto space-y-6">
      <!-- Header -->
      <header class="pixel-box bg-slate-900/95 p-6">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div class="flex items-center gap-4">
            <div class="w-14 h-14 bg-amber-600 pixel-box-sm flex items-center justify-center">
              <span class="text-3xl">🏆</span>
            </div>
            <div>
              <h1 class="text-2xl md:text-3xl font-bold text-amber-400 text-shadow-retro">HALL OF FAME</h1>
              <p class="text-slate-400 text-sm">Prove your worth. Claim your glory.</p>
            </div>
          </div>
          
          <!-- User Rank Card -->
          <div v-if="leaderboardStore.userEntry" class="pixel-box-sm bg-amber-900/30 border-amber-500/50 px-4 py-2">
            <div class="text-xs text-amber-400 uppercase">Your Rank</div>
            <div class="flex items-center gap-2">
              <span class="text-2xl font-black text-white">#{{ leaderboardStore.userRank }}</span>
              <span v-if="leaderboardStore.userEntry.trend === 'up'" class="text-green-400">↑</span>
              <span v-else-if="leaderboardStore.userEntry.trend === 'down'" class="text-red-400">↓</span>
            </div>
            <div class="text-xs text-slate-400">{{ formatScore(leaderboardStore.userEntry.score) }} pts</div>
          </div>
        </div>
      </header>

      <!-- Category Tabs -->
      <div class="flex flex-wrap gap-2">
        <button
          v-for="(label, key) in CATEGORY_LABELS"
          :key="key"
          @click="setCategory(key as LeaderboardCategory)"
          class="pixel-box-sm px-4 py-2 text-sm font-bold transition-all"
          :class="currentCategory === key ? 'bg-amber-600 text-white border-amber-400' : 'bg-slate-800 text-slate-400 hover:bg-slate-700'"
        >
          <span class="mr-1">{{ label.icon }}</span>
          {{ label.name }}
        </button>
      </div>

      <!-- Main Content Grid -->
      <div class="grid lg:grid-cols-3 gap-6">
        <!-- Leaderboard -->
        <div class="lg:col-span-2 space-y-4">
          <div class="pixel-box bg-slate-900/90 p-4">
            <h2 class="text-lg font-bold text-amber-400 mb-4 flex items-center gap-2">
              <span>{{ CATEGORY_LABELS[currentCategory].icon }}</span>
              {{ CATEGORY_LABELS[currentCategory].name }}
              <span class="text-xs text-slate-500 font-normal">- {{ CATEGORY_LABELS[currentCategory].description }}</span>
            </h2>

            <!-- Loading State -->
            <div v-if="leaderboardStore.isLoading" class="py-12 text-center">
              <div class="animate-spin text-3xl mb-2">⚔️</div>
              <p class="text-amber-400 animate-pulse">Loading rankings...</p>
            </div>

            <!-- Entries List -->
            <div v-else class="space-y-2">
              <div
                v-for="(entry, index) in leaderboardStore.topEntries"
                :key="entry.id"
                class="pixel-box-sm p-3 flex items-center gap-4 transition-all"
                :class="getEntryClass(entry, index)"
              >
                <!-- Rank -->
                <div class="w-10 text-center">
                  <span v-if="entry.rank <= 3" class="text-2xl">{{ getRankEmoji(entry.rank) }}</span>
                  <span v-else class="text-lg font-bold text-slate-400">#{{ entry.rank }}</span>
                </div>

                <!-- Avatar placeholder -->
                <div class="w-10 h-10 pixel-box-sm bg-slate-800 flex items-center justify-center">
                  <span class="text-xl">🎮</span>
                </div>

                <!-- Username -->
                <div class="flex-1">
                  <div class="font-bold text-white">{{ entry.username }}</div>
                  <div class="text-xs text-slate-500">
                    <span v-if="entry.trend === 'up'" class="text-green-400">↑ Rising</span>
                    <span v-else-if="entry.trend === 'down'" class="text-red-400">↓ Falling</span>
                    <span v-else>→ Stable</span>
                  </div>
                </div>

                <!-- Score -->
                <div class="text-right">
                  <div class="font-black text-amber-400">{{ formatScore(entry.score) }}</div>
                  <div class="text-xs text-slate-500">points</div>
                </div>
              </div>

              <!-- Empty State -->
              <div v-if="leaderboardStore.topEntries.length === 0" class="py-12 text-center text-slate-500">
                <p>No rankings available yet.</p>
                <p class="text-sm">Be the first to climb the ranks!</p>
              </div>
            </div>
          </div>

          <!-- Nearby Ranks -->
          <div v-if="leaderboardStore.nearbyRanks?.entries?.length" class="pixel-box bg-slate-900/90 p-4">
            <h3 class="text-sm font-bold text-slate-400 mb-3 uppercase tracking-wider">Nearby Competition</h3>
            <div class="space-y-1">
              <div
                v-for="entry in leaderboardStore.nearbyRanks.entries"
                :key="entry.id"
                class="pixel-box-sm p-2 flex items-center gap-3 text-sm"
                :class="entry.user_id === leaderboardStore.leaderboard?.user_entry?.user_id ? 'bg-amber-900/20 border-amber-500/30' : 'bg-slate-800/50'"
              >
                <span class="w-8 text-center font-bold" :class="entry.rank === leaderboardStore.userRank ? 'text-amber-400' : 'text-slate-500'">
                  #{{ entry.rank }}
                </span>
                <span class="flex-1 truncate">{{ entry.username }}</span>
                <span class="font-bold text-amber-400">{{ formatScore(entry.score) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Achievements Sidebar -->
        <div class="space-y-4">
          <div class="pixel-box bg-slate-900/90 p-4">
            <h2 class="text-lg font-bold text-amber-400 mb-4 flex items-center gap-2">
              <span>⭐</span> Achievements
            </h2>

            <!-- Progress Summary -->
            <div class="mb-4 space-y-2">
              <div v-for="(progress, cat) in leaderboardStore.progressByCategory" :key="cat" class="pixel-box-sm bg-slate-800/50 p-2">
                <div class="flex justify-between text-xs mb-1">
                  <span class="text-slate-400 capitalize">{{ cat }}</span>
                  <span class="text-amber-400">{{ progress.completed }}/{{ progress.total }}</span>
                </div>
                <div class="h-1.5 bg-slate-700 rounded-full overflow-hidden">
                  <div class="h-full bg-amber-500 transition-all" :style="{ width: `${(progress.completed / progress.total) * 100}%` }"></div>
                </div>
              </div>
            </div>

            <!-- Achievement List -->
            <div class="space-y-2 max-h-[400px] overflow-y-auto custom-scrollbar">
              <div
                v-for="pa in unclaimedAchievements.slice(0, 5)"
                :key="pa.id"
                class="pixel-box-sm p-3 bg-amber-900/20 border-amber-500/30 cursor-pointer hover:bg-amber-900/30 transition-colors"
                @click="claimReward(pa.achievement_id)"
              >
                <div class="flex items-start gap-3">
                  <span class="text-2xl">{{ pa.achievement?.icon }}</span>
                  <div class="flex-1">
                    <div class="font-bold text-amber-400 text-sm">{{ pa.achievement?.name }}</div>
                    <div class="text-xs text-slate-400">{{ pa.achievement?.description }}</div>
                    <div class="mt-1 text-xs text-green-400 font-bold">Click to claim {{ pa.achievement?.reward_gold }} gold!</div>
                  </div>
                </div>
              </div>

              <div
                v-for="pa in completedAchievements.filter(pa => pa.claimed).slice(0, 5)"
                :key="pa.id"
                class="pixel-box-sm p-3 bg-slate-800/50 opacity-60"
              >
                <div class="flex items-start gap-3">
                  <span class="text-2xl">{{ pa.achievement?.icon }}</span>
                  <div class="flex-1">
                    <div class="font-bold text-slate-300 text-sm">{{ pa.achievement?.name }}</div>
                    <div class="text-xs text-slate-500">{{ pa.achievement?.description }}</div>
                  </div>
                  <span class="text-green-500 text-xs">✓ Claimed</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useLeaderboardStore } from '@/features/leaderboard/store';
import { CATEGORY_LABELS, type LeaderboardCategory } from '@/features/leaderboard/api';

const PIXEL_ASSETS = {
  BG_DUNGEON: 'https://vibemedia.space/bg_dungeon_v2_99283.png?prompt=dark%20dungeon%20floor%20tile%20texture%20seamless&style=pixel_game_asset&key=NOGON',
};

const leaderboardStore = useLeaderboardStore();
const currentCategory = ref<LeaderboardCategory>('power');

const unclaimedAchievements = computed(() => leaderboardStore.unclaimedAchievements);
const completedAchievements = computed(() => leaderboardStore.completedAchievements);

function setCategory(cat: LeaderboardCategory) {
  currentCategory.value = cat;
  leaderboardStore.fetchLeaderboard(cat);
  leaderboardStore.fetchNearbyRanks(cat);
}

function getEntryClass(entry: any, index: number) {
  if (entry.user_id === leaderboardStore.leaderboard?.user_entry?.user_id) {
    return 'bg-amber-900/30 border-amber-500/50';
  }
  if (index === 0) return 'bg-yellow-900/20 border-yellow-500/30';
  if (index === 1) return 'bg-slate-700/30 border-slate-400/30';
  if (index === 2) return 'bg-orange-900/20 border-orange-500/30';
  return 'bg-slate-800/50';
}

function getRankEmoji(rank: number) {
  if (rank === 1) return '🥇';
  if (rank === 2) return '🥈';
  if (rank === 3) return '🥉';
  return `#${rank}`;
}

function formatScore(score: number) {
  if (score >= 1000000) return `${(score / 1000000).toFixed(1)}M`;
  if (score >= 1000) return `${(score / 1000).toFixed(1)}K`;
  return score.toLocaleString();
}

async function claimReward(achievementId: string) {
  try {
    await leaderboardStore.claimAchievementReward(achievementId);
  } catch (e) {
    // Error handled by store
  }
}

onMounted(() => {
  leaderboardStore.fetchLeaderboard('power');
  leaderboardStore.fetchNearbyRanks('power');
  leaderboardStore.fetchAchievements();
  leaderboardStore.fetchPlayerAchievements();
});
</script>
