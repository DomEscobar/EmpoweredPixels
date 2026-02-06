<template>
  <div 
    class="min-h-screen p-4 md:p-8 font-mono"
    :style="{ 
      backgroundImage: `url('${PIXEL_ASSETS.BG_DUNGEON}')`,
      backgroundSize: '128px 128px',
      imageRendering: 'pixelated'
    }"
  >
    <!-- CRT Scanline Overlay -->
    <div class="pointer-events-none fixed inset-0 z-50 opacity-[0.03] bg-[repeating-linear-gradient(0deg,transparent,transparent_2px,rgba(0,0,0,0.3)_2px,rgba(0,0,0,0.3)_4px)]"></div>
    
    <!-- Vignette -->
    <div class="pointer-events-none fixed inset-0 z-40 bg-[radial-gradient(ellipse_at_center,transparent_0%,rgba(0,0,0,0.4)_100%)]"></div>

    <div class="relative z-10 max-w-7xl mx-auto space-y-6">
      
      <!-- Header Banner -->
      <header class="pixel-box bg-slate-900/95 p-6">
        <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <div class="flex items-center gap-4">
            <img 
              :src="PIXEL_ASSETS.ICON_CASTLE" 
              alt="" 
              class="w-14 h-14 pixelated"
            />
            <div>
              <h1 class="text-2xl md:text-3xl font-bold text-amber-400 text-shadow-retro tracking-wide">
                COMMAND CENTER
              </h1>
              <p class="text-slate-400 text-sm mt-1">
                Welcome back, Commander. All systems operational.
              </p>
            </div>
          </div>
          
          <!-- Status Indicator -->
          <div class="pixel-box-sm bg-emerald-900/30 border-emerald-500/50 px-4 py-2 flex items-center gap-2">
            <div class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></div>
            <span class="text-emerald-400 text-xs font-bold uppercase tracking-wider">Online</span>
          </div>
        </div>
      </header>

      <!-- Event Banner -->
      <EventBanner />

      <!-- Daily Reward Modal -->
      <DailyRewardModal 
        :show="showDailyReward" 
        @close="showDailyReward = false" 
      />

      <!-- KPI Grid -->
      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        
        <!-- Active Roster -->
        <div class="pixel-box bg-slate-900/90 p-4 group hover:bg-slate-900/95 transition-colors">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs text-slate-500 uppercase tracking-wider">Active Roster</p>
              <p class="mt-2 text-3xl font-bold text-white text-shadow-retro">{{ rosterStore.fighters.length }}</p>
            </div>
            <div class="p-3 pixel-box-sm bg-indigo-900/30 border-indigo-500/30 group-hover:border-indigo-400/50 transition-colors">
              <img :src="PIXEL_ASSETS.ICON_KNIGHT" alt="" class="w-8 h-8 pixelated" />
            </div>
          </div>
          <div class="mt-4 h-2 pixel-box-sm bg-slate-800/80 overflow-hidden">
            <div class="h-full bg-linear-to-r from-indigo-600 to-indigo-400" style="width: 75%"></div>
          </div>
        </div>

        <!-- Active Campaigns -->
        <div class="pixel-box bg-slate-900/90 p-4 group hover:bg-slate-900/95 transition-colors">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs text-slate-500 uppercase tracking-wider">Active Campaigns</p>
              <p class="mt-2 text-3xl font-bold text-white text-shadow-retro">{{ leaguesStore.activeLeagueCount }}</p>
            </div>
            <div class="p-3 pixel-box-sm bg-purple-900/30 border-purple-500/30 group-hover:border-purple-400/50 transition-colors">
              <img :src="PIXEL_ASSETS.ICON_FLAG" alt="" class="w-8 h-8 pixelated" />
            </div>
          </div>
          <div class="mt-4 flex items-center text-xs">
            <span class="flex items-center gap-1 text-emerald-400 font-bold">
              <span class="w-1.5 h-1.5 rounded-full bg-emerald-400"></span>
              LIVE
            </span>
            <span class="text-slate-500 ml-2">Season 1</span>
          </div>
        </div>

        <!-- Combat Record -->
        <div class="pixel-box bg-slate-900/90 p-4 group hover:bg-slate-900/95 transition-colors">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs text-slate-500 uppercase tracking-wider">Combat Record</p>
              <p class="mt-2 text-3xl font-bold text-white text-shadow-retro">{{ matchesStore.recentMatches.length }}</p>
            </div>
            <div class="p-3 pixel-box-sm bg-pink-900/30 border-pink-500/30 group-hover:border-pink-400/50 transition-colors">
              <img :src="PIXEL_ASSETS.ICON_SWORDS" alt="" class="w-8 h-8 pixelated" />
            </div>
          </div>
          <div class="mt-4 text-xs text-slate-500">
            Recent engagements
          </div>
        </div>

        <!-- Pending Rewards -->
        <div 
          class="pixel-box bg-slate-900/90 p-4 group transition-all"
          :class="rewardsStore.rewardCount > 0 ? 'border-amber-500/50 shadow-lg shadow-amber-500/10' : ''"
        >
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs text-slate-500 uppercase tracking-wider">Pending Rewards</p>
              <p class="mt-2 text-3xl font-bold text-white text-shadow-retro">{{ rewardsStore.rewardCount }}</p>
            </div>
            <div 
              class="p-3 pixel-box-sm bg-amber-900/30 border-amber-500/30 transition-colors"
              :class="rewardsStore.rewardCount > 0 ? 'animate-pulse border-amber-400/60' : ''"
            >
              <img :src="PIXEL_ASSETS.ICON_CHEST" alt="" class="w-8 h-8 pixelated" />
            </div>
          </div>
          <div class="mt-4">
            <button 
              v-if="rewardsStore.rewardCount > 0" 
              class="rpg-btn w-full text-sm py-2"
              @click="claimAllRewards"
            >
              CLAIM ALL
            </button>
            <p v-else class="text-xs text-slate-600">No pending rewards</p>
          </div>
        </div>
      </div>

      <!-- Main Content Split -->
      <div class="grid gap-6 lg:grid-cols-3">
        
        <!-- Left Column: Combat Log -->
        <div class="lg:col-span-2">
          <div class="pixel-box bg-slate-900/90 h-full">
            <!-- Header -->
            <div class="p-4 border-b-4 border-slate-800 flex items-center justify-between">
              <div class="flex items-center gap-3">
                <img :src="PIXEL_ASSETS.ICON_SCROLL" alt="" class="w-5 h-5 pixelated" />
                <h3 class="text-lg font-bold text-amber-300">BATTLE LOG</h3>
              </div>
              <router-link to="/matches" class="rpg-btn-small">
                VIEW ALL
              </router-link>
            </div>
            
            <!-- Content -->
            <div class="p-4">
              <div v-if="matchesStore.isLoading" class="py-12 text-center">
                <div class="inline-block animate-spin text-3xl mb-3">⚔️</div>
                <p class="text-amber-400 animate-pulse">Loading battle data...</p>
              </div>
              
              <div v-else-if="matchesStore.recentMatches.length === 0" class="py-12 text-center">
                <img :src="PIXEL_ASSETS.ICON_SCROLL_EMPTY" alt="" class="w-12 h-12 pixelated mx-auto mb-3 opacity-50" />
                <p class="text-slate-400">No recent combat data</p>
                <router-link to="/matches" class="rpg-btn-small mt-4 inline-flex">
                  JOIN BATTLE
                </router-link>
              </div>
              
              <div v-else class="space-y-2 max-h-80 overflow-y-auto custom-scrollbar">
                <div 
                  v-for="match in matchesStore.recentMatches" 
                  :key="match.id" 
                  class="pixel-box-sm bg-slate-800/60 p-3 flex items-center justify-between hover:bg-slate-800/80 transition-colors cursor-pointer"
                >
                  <div class="flex items-center gap-3">
                    <div 
                      class="w-2 h-2 rounded-full"
                      :class="match.ended ? 'bg-slate-500' : 'bg-emerald-400 animate-pulse'"
                    ></div>
                    <div>
                      <p class="font-bold text-slate-200 text-sm">Battle #{{ match.id.slice(0, 8) }}</p>
                      <p class="text-xs text-slate-500">{{ formatDate(match.created) }}</p>
                    </div>
                  </div>
                  <span 
                    class="pixel-badge px-2 py-0.5 text-xs font-bold uppercase"
                    :class="match.ended ? 'bg-slate-700 text-slate-400' : 'bg-emerald-900/50 text-emerald-400 border-emerald-500/50'"
                  >
                    {{ match.ended ? 'ENDED' : 'LIVE' }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Right Column -->
        <div class="space-y-6">
          
          <!-- Elite Operative -->
          <div class="pixel-box bg-slate-900/90">
            <div class="p-4 border-b-4 border-slate-800 flex items-center gap-3">
              <img :src="PIXEL_ASSETS.ICON_CROWN" alt="" class="w-5 h-5 pixelated" />
              <h3 class="text-lg font-bold text-amber-300">CHAMPION</h3>
            </div>
            
            <div class="p-4">
              <div v-if="topFighter" class="text-center">
                <div class="relative inline-block">
                  <div class="pixel-box-sm w-20 h-20 mx-auto flex items-center justify-center bg-slate-800/80 border-indigo-500/50 mb-3">
                    <img :src="PIXEL_ASSETS.ICON_WARRIOR" alt="" class="w-14 h-14 pixelated" />
                  </div>
                  <div class="absolute -bottom-1 -right-1 bg-indigo-600 text-white text-xs px-2 py-0.5 font-bold border-2 border-slate-900">
                    LVL {{ topFighter.level }}
                  </div>
                </div>
                <h4 class="text-lg font-bold text-white mt-2">{{ topFighter.name }}</h4>
                <div class="mt-2">
                  <span class="pixel-badge bg-slate-800 text-slate-300 px-2 py-1 text-xs border-slate-600">
                    Power: {{ topFighter.power }}
                  </span>
                </div>
              </div>
              
              <div v-else class="py-6 text-center">
                <p class="text-slate-500 text-sm">No champions found</p>
                <router-link to="/roster" class="rpg-btn-small mt-3 inline-flex">
                  RECRUIT
                </router-link>
              </div>
            </div>
          </div>

          <!-- Quick Actions -->
          <div class="pixel-box bg-slate-900/90">
            <div class="p-4 border-b-4 border-slate-800">
              <h3 class="text-xs font-bold text-slate-500 uppercase tracking-wider">Quick Actions</h3>
            </div>
            
            <div class="p-4 grid grid-cols-2 gap-3">
              <router-link 
                to="/matches" 
                class="pixel-box-sm bg-slate-800/60 p-4 flex flex-col items-center justify-center hover:bg-slate-800/80 hover:border-indigo-500/50 transition-all group"
              >
                <img :src="PIXEL_ASSETS.ICON_SWORDS" alt="" class="w-8 h-8 pixelated mb-2 group-hover:scale-110 transition-transform" />
                <span class="text-xs font-bold text-slate-300">BATTLE</span>
              </router-link>
              
              <router-link 
                to="/roster" 
                class="pixel-box-sm bg-slate-800/60 p-4 flex flex-col items-center justify-center hover:bg-slate-800/80 hover:border-purple-500/50 transition-all group"
              >
                <img :src="PIXEL_ASSETS.ICON_KNIGHT" alt="" class="w-8 h-8 pixelated mb-2 group-hover:scale-110 transition-transform" />
                <span class="text-xs font-bold text-slate-300">ROSTER</span>
              </router-link>
              
              <router-link 
                to="/inventory" 
                class="pixel-box-sm bg-slate-800/60 p-4 flex flex-col items-center justify-center hover:bg-slate-800/80 hover:border-emerald-500/50 transition-all group"
              >
                <img :src="PIXEL_ASSETS.ICON_CHEST" alt="" class="w-8 h-8 pixelated mb-2 group-hover:scale-110 transition-transform" />
                <span class="text-xs font-bold text-slate-300">VAULT</span>
              </router-link>
              
              <router-link 
                to="/leagues" 
                class="pixel-box-sm bg-slate-800/60 p-4 flex flex-col items-center justify-center hover:bg-slate-800/80 hover:border-amber-500/50 transition-all group"
              >
                <img :src="PIXEL_ASSETS.ICON_TROPHY" alt="" class="w-8 h-8 pixelated mb-2 group-hover:scale-110 transition-transform" />
                <span class="text-xs font-bold text-slate-300">LEAGUES</span>
              </router-link>
            </div>
          </div>

        </div>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, computed, ref } from "vue";
import { useRosterStore } from "@/features/roster/store";
import { useMatchesStore } from "@/features/matches/store";
import { useRewardsStore } from "@/features/rewards/store";
import { useLeaguesStore } from "@/features/leagues/store";
import { useDailyStore } from "@/features/daily/store";
import EventBanner from "@/features/events/components/EventBanner.vue";
import DailyRewardModal from "@/features/daily/components/DailyRewardModal.vue";

const PIXEL_ASSETS = {
  BG_DUNGEON: 'https://vibemedia.space/bg_dungeon_dash_5y6t7u_v1.png?prompt=dark%20dungeon%20stone%20floor%20tile%20texture%20seamless&style=pixel_game_asset&key=NOGON',
  ICON_CASTLE: 'https://vibemedia.space/icon_castle_cmd_8i9o0p_v1.png?prompt=medieval%20castle%20tower%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  ICON_KNIGHT: 'https://vibemedia.space/icon_knight_1a2s3d_v1.png?prompt=armored%20knight%20helmet%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  ICON_FLAG: 'https://vibemedia.space/icon_flag_war_4f5g6h_v1.png?prompt=war%20banner%20flag%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  ICON_SWORDS: 'https://vibemedia.space/icon_crossed_swords_7j8k9l_v1.png?prompt=crossed%20steel%20swords%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  ICON_CHEST: 'https://vibemedia.space/icon_chest_gold_0m1n2b_v1.png?prompt=golden%20treasure%20chest%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  ICON_SCROLL: 'https://vibemedia.space/icon_scroll_log_3v4c5x_v1.png?prompt=ancient%20scroll%20with%20writing%20pixel%20art&style=pixel_game_asset&key=NOGON',
  ICON_SCROLL_EMPTY: 'https://vibemedia.space/icon_scroll_blank_6z7a8s_v1.png?prompt=blank%20parchment%20scroll%20pixel%20art&style=pixel_game_asset&key=NOGON',
  ICON_CROWN: 'https://vibemedia.space/icon_crown_gold_9d0f1g_v1.png?prompt=golden%20royal%20crown%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  ICON_WARRIOR: 'https://vibemedia.space/icon_warrior_elite_2h3j4k_v1.png?prompt=elite%20warrior%20knight%20pixel%20art%20character&style=pixel_game_asset&key=NOGON',
  ICON_TROPHY: 'https://vibemedia.space/icon_trophy_league_5l6p7q_v1.png?prompt=golden%20trophy%20cup%20pixel%20art&style=pixel_game_asset&key=NOGON'
};

const rosterStore = useRosterStore();
const matchesStore = useMatchesStore();
const rewardsStore = useRewardsStore();
const leaguesStore = useLeaguesStore();
const dailyStore = useDailyStore();

const showDailyReward = ref(false);

const topFighter = computed(() => {
  if (!rosterStore.fighters.length) return null;
  return [...rosterStore.fighters].sort((a, b) => {
     if (b.level !== a.level) return b.level - a.level;
     return b.power - a.power;
  })[0];
});

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString(undefined, {
    month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit'
  });
};

const claimAllRewards = async () => {
   await rewardsStore.claimAll();
};

onMounted(async () => {
  await dailyStore.fetchStatus();
  if (dailyStore.canClaim) {
    showDailyReward.value = true;
  }

  Promise.all([
    rosterStore.fetchFighters(),
    matchesStore.fetchRecentMatches(),
    rewardsStore.fetchRewards(),
    leaguesStore.fetchLeagues()
  ]);
});
</script>

<style scoped>
.pixelated {
  image-rendering: pixelated;
}

.text-shadow-retro {
  text-shadow: 2px 2px 0 #000, 4px 4px 0 rgba(0, 0, 0, 0.3);
}

.pixel-box {
  border: 4px solid #1e293b;
  box-shadow: 
    inset 0 0 0 2px #334155,
    4px 4px 0 #0f172a;
  image-rendering: pixelated;
}

.pixel-box-sm {
  border: 2px solid #334155;
  box-shadow: 2px 2px 0 #0f172a;
}

.pixel-badge {
  border: 2px solid #475569;
}

.rpg-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.625rem 1.25rem;
  font-weight: 700;
  font-size: 0.875rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #fef3c7;
  background: linear-gradient(to bottom, #d97706, #b45309);
  border: 3px solid #92400e;
  box-shadow: 
    inset 0 1px 0 rgba(255, 255, 255, 0.2),
    0 4px 0 #78350f,
    0 6px 4px rgba(0, 0, 0, 0.4);
  cursor: pointer;
  transition: all 0.1s;
}

.rpg-btn:hover:not(:disabled) {
  background: linear-gradient(to bottom, #f59e0b, #d97706);
  transform: translateY(-1px);
  box-shadow: 
    inset 0 1px 0 rgba(255, 255, 255, 0.3),
    0 5px 0 #78350f,
    0 7px 6px rgba(0, 0, 0, 0.4);
}

.rpg-btn:active:not(:disabled) {
  transform: translateY(2px);
  box-shadow: 
    inset 0 1px 0 rgba(255, 255, 255, 0.1),
    0 2px 0 #78350f,
    0 3px 2px rgba(0, 0, 0, 0.3);
}

.rpg-btn-small {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.375rem 0.75rem;
  font-weight: 600;
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #e2e8f0;
  background: linear-gradient(to bottom, #475569, #334155);
  border: 2px solid #1e293b;
  box-shadow: 
    inset 0 1px 0 rgba(255, 255, 255, 0.1),
    0 2px 0 #0f172a;
  cursor: pointer;
  transition: all 0.1s;
}

.rpg-btn-small:hover:not(:disabled) {
  background: linear-gradient(to bottom, #64748b, #475569);
}

.rpg-btn-small:active:not(:disabled) {
  transform: translateY(1px);
  box-shadow: 0 1px 0 #0f172a;
}

.custom-scrollbar::-webkit-scrollbar {
  width: 8px;
}

.custom-scrollbar::-webkit-scrollbar-track {
  background: #1e293b;
  border: 1px solid #334155;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #475569;
  border: 1px solid #64748b;
}

.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #64748b;
}
</style>
