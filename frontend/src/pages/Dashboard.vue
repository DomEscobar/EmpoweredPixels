<template>
  <div 
    class="min-h-screen pb-24 md:pb-8 px-4 md:p-6 font-mono"
    :style="{ 
      backgroundImage: `url('${PIXEL_ASSETS.BG_DUNGEON}')`,
      backgroundSize: '128px 128px',
      imageRendering: 'pixelated'
    }"
    data-testid="dashboard-page"
  >
    <!-- CRT Scanline Overlay -->
    <div class="pointer-events-none fixed inset-0 z-50 opacity-[0.03] bg-[repeating-linear-gradient(0deg,transparent,transparent_2px,rgba(0,0,0,0.3)_2px,rgba(0,0,0,0.3)_4px)]"></div>
    
    <!-- Vignette -->
    <div class="pointer-events-none fixed inset-0 z-40 bg-[radial-gradient(ellipse_at_center,transparent_0%,rgba(0,0,0,0.4)_100%)]"></div>

    <div class="relative z-10 max-w-4xl mx-auto space-y-4 md:space-y-6">
      
      <!-- Compact Header Banner -->
      <header class="pixel-box bg-slate-900/95 p-3 md:p-6" data-testid="dashboard-header">
        <div class="flex items-center justify-between gap-3">
          <div class="flex items-center gap-3 min-w-0">
            <img 
              :src="PIXEL_ASSETS.ICON_CASTLE" 
              alt="" 
              class="w-10 h-10 md:w-14 md:h-14 pixelated flex-shrink-0"
            />
            <div class="min-w-0">
              <h1 class="text-lg md:text-2xl font-bold text-amber-400 text-shadow-retro tracking-wide truncate">
                COMMAND CENTER
              </h1>
              <p class="text-slate-400 text-xs md:text-sm mt-0.5 truncate">
                Welcome back, Commander.
              </p>
            </div>
          </div>
          
          <!-- Status Indicator -->
          <div class="pixel-box-sm bg-emerald-900/30 border-emerald-500/50 px-2 md:px-4 py-1 md:py-2 flex items-center gap-1 md:gap-2 flex-shrink-0" data-testid="user-status">
            <div class="w-1.5 h-1.5 md:w-2 md:h-2 rounded-full bg-emerald-400 animate-pulse"></div>
             <span class="text-emerald-400 text-xs md:text-xs font-bold uppercase tracking-wider">Online</span>
          </div>
        </div>
      </header>

      <!-- Event Banner -->
      <EventBanner data-testid="event-banner" />

      <!-- League Deadlines (New Integration) -->
      <div v-if="imminentLeagues.length > 0" class="pixel-box bg-linear-to-r from-red-900/40 to-amber-900/20 border-red-500/30 p-3 md:p-4 flex flex-col md:flex-row items-center justify-between gap-3">
          <div class="flex items-center gap-2 md:gap-3 min-w-0">
              <span class="text-xl md:text-2xl animate-pulse">⏰</span>
              <div class="min-w-0">
                  <h4 class="text-red-400 font-black text-[10px] md:text-xs uppercase tracking-[0.1em] md:tracking-[0.2em]">League Deadline Approaching</h4>
                  <p class="text-slate-300 text-xs md:text-sm truncate">Competition entries closing soon for <span class="text-amber-400 font-bold">{{ imminentLeagues.map(l => l.name).join(', ') }}</span></p>
              </div>
          </div>
           <router-link to="/leagues" class="rpg-btn-small border-red-500/50 bg-red-950/50 hover:bg-red-900 text-red-200 text-sm md:text-xs px-3 md:px-3 py-2 w-full md:w-auto min-h-[44px]">
               SECURE SPOT
           </router-link>
      </div>

       <!-- KPI Grid - Mobile Optimized with Trend Indicators -->
        <div class="grid gap-2 md:gap-4 grid-cols-2 md:grid-cols-4">
          
          <!-- Active Roster -->
          <div class="pixel-box bg-slate-900/90 p-2 md:p-4 group hover:bg-slate-900/95 transition-all hover:scale-[1.02] min-h-[72px] md:min-h-0 touch-manipulation" data-testid="kpi-roster">
            <div class="flex items-center gap-2 md:gap-3">
              <div class="pixel-box-sm bg-indigo-900/50 border-indigo-500/50 flex items-center justify-center min-w-[44px] min-h-[44px] md:min-w-[44px] md:min-h-[44px]">
                <img :src="PIXEL_ASSETS.ICON_KNIGHT" alt="" class="w-6 h-6 md:w-8 md:h-8 pixelated" />
              </div>
              <div class="min-w-0 flex-1">
                <p class="text-[9px] md:text-xs text-slate-400 uppercase tracking-wider mb-0.5">Roster</p>
                <div class="flex items-baseline gap-1">
                  <p class="text-lg md:text-3xl font-bold text-white text-shadow-retro truncate">{{ rosterStore.fighters.length }}</p>
                  <span class="text-[9px] text-emerald-400 font-bold hidden sm:inline">↑</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Active Campaigns -->
          <div class="pixel-box bg-slate-900/90 p-2 md:p-4 group hover:bg-slate-900/95 transition-all hover:scale-[1.02] min-h-[72px] md:min-h-0 touch-manipulation" data-testid="kpi-campaigns">
            <div class="flex items-center gap-2 md:gap-3">
              <div class="pixel-box-sm bg-purple-900/50 border-purple-500/50 flex items-center justify-center min-w-[44px] min-h-[44px] md:min-w-[44px] md:min-h-[44px]">
                <img :src="PIXEL_ASSETS.ICON_FLAG" alt="" class="w-6 h-6 md:w-8 md:h-8 pixelated" />
              </div>
              <div class="min-w-0 flex-1">
                <p class="text-[9px] md:text-xs text-slate-400 uppercase tracking-wider mb-0.5">Campaigns</p>
                <div class="flex items-baseline gap-1">
                  <p class="text-lg md:text-3xl font-bold text-white text-shadow-retro truncate">{{ leaguesStore.activeLeagueCount }}</p>
                  <span class="text-[9px] text-amber-400 font-bold animate-pulse hidden sm:inline">●</span>
                </div>
              </div>
            </div>
            <div class="mt-1 flex items-center text-[9px]">
              <span class="flex items-center gap-1 text-emerald-400 font-bold">
                <span class="w-1.5 h-1.5 rounded-full bg-emerald-400"></span>
                LIVE
              </span>
            </div>
          </div>

          <!-- Combat Record -->
          <div class="pixel-box bg-slate-900/90 p-2 md:p-4 group hover:bg-slate-900/95 transition-all hover:scale-[1.02] min-h-[72px] md:min-h-0 touch-manipulation" data-testid="kpi-combat">
            <div class="flex items-center gap-2 md:gap-3">
              <div class="pixel-box-sm bg-pink-900/50 border-pink-500/50 flex items-center justify-center min-w-[44px] min-h-[44px] md:min-w-[44px] md:min-h-[44px]">
                <img :src="PIXEL_ASSETS.ICON_SWORDS" alt="" class="w-6 h-6 md:w-8 md:h-8 pixelated" />
              </div>
              <div class="min-w-0 flex-1">
                <p class="text-[9px] md:text-xs text-slate-400 uppercase tracking-wider mb-0.5">Battles</p>
                <div class="flex items-baseline gap-1">
                  <p class="text-lg md:text-3xl font-bold text-white text-shadow-retro truncate">{{ matchesStore.recentMatches.length }}</p>
                  <span v-if="matchesStore.recentMatches.length > 0" class="text-[9px] text-amber-400 font-bold hidden sm:inline">↑</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Pending Rewards -->
          <div 
            class="pixel-box bg-slate-900/90 p-2 md:p-4 group transition-all hover:scale-[1.02] min-h-[72px] md:min-h-0 touch-manipulation"
            :class="rewardsStore.rewardCount > 0 ? 'border-amber-500/50 animate-pulse' : ''"
            data-testid="kpi-rewards"
          >
            <div class="flex items-center gap-2 md:gap-3">
              <div 
                class="pixel-box-sm bg-amber-900/50 border-amber-500/50 flex items-center justify-center min-w-[44px] min-h-[44px] md:min-w-[44px] md:min-h-[44px]"
                :class="rewardsStore.rewardCount > 0 ? 'animate-pulse' : ''"
              >
                <img :src="PIXEL_ASSETS.ICON_CHEST" alt="" class="w-6 h-6 md:w-8 md:h-8 pixelated" />
              </div>
              <div class="min-w-0 flex-1">
                <p class="text-[9px] md:text-xs text-slate-400 uppercase tracking-wider mb-0.5">Rewards</p>
                <div class="flex items-baseline gap-1">
                  <p class="text-lg md:text-3xl font-bold text-white text-shadow-retro truncate">{{ rewardsStore.rewardCount }}</p>
                  <span v-if="rewardsStore.rewardCount > 0" class="text-[9px] text-emerald-400 font-bold hidden sm:inline">●</span>
                </div>
              </div>
            </div>
          </div>
        </div>

      <!-- Main Content - Single Column on Mobile -->
      <div class="space-y-4 md:space-y-6">
        
          <!-- Combat Log - Enhanced Mobile -->
          <div class="pixel-box bg-slate-900/90">
            <!-- Header -->
            <div class="p-2 md:p-4 border-b-4 border-slate-800 flex items-center justify-between">
              <div class="flex items-center gap-2 md:gap-3">
                <img :src="PIXEL_ASSETS.ICON_SCROLL" alt="" class="w-5 h-5 md:w-5 md:h-5 pixelated" />
                <h3 class="text-sm md:text-lg font-bold text-amber-300">BATTLE LOG</h3>
                <span v-if="matchesStore.recentMatches.length > 0" class="text-[9px] bg-slate-800 px-1.5 py-0.5 rounded hidden sm:inline">{{ matchesStore.recentMatches.length }}</span>
              </div>
               <router-link to="/matches" class="rpg-btn-small text-sm min-h-[44px]" data-testid="battle-log-view-all">
                  VIEW ALL
               </router-link>
            </div>
            
            <!-- Content -->
            <div class="p-2 md:p-4">
              <div v-if="matchesStore.isLoading" class="py-6 md:py-12 text-center">
                <div class="inline-block animate-spin text-2xl md:text-3xl mb-2">⚔️</div>
                <p class="text-amber-400 animate-pulse text-xs md:text-base">Loading battles...</p>
              </div>
              
              <div v-else-if="matchesStore.recentMatches.length === 0" class="py-6 md:py-12 text-center">
                <img :src="PIXEL_ASSETS.ICON_SCROLL_EMPTY" alt="" class="w-10 h-10 md:w-12 md:h-12 pixelated mx-auto mb-2 opacity-50" />
                <p class="text-slate-400 text-xs md:text-sm mb-3">No battles yet</p>
                <router-link to="/matches" class="rpg-btn-small inline-flex text-xs py-2 px-4 min-h-[44px]">
                  JOIN BATTLE
                </router-link>
              </div>
              
              <div v-else class="space-y-1.5 max-h-52 md:max-h-80 overflow-y-auto custom-scrollbar">
                <div 
                  v-for="match in matchesStore.recentMatches.slice(0, 5)" 
                  :key="match.id" 
                  class="pixel-box-sm bg-slate-800/60 p-2.5 flex items-center justify-between hover:bg-slate-800/80 hover:border-indigo-500/30 transition-all cursor-pointer min-h-[52px] active:scale-[0.98]"
                  @click="goToMatch(match.id)"
                >
                  <div class="flex items-center gap-2.5 min-w-0 flex-1">
                    <div 
                      class="w-2 h-2 rounded-full flex-shrink-0"
                      :class="match.ended ? 'bg-slate-500' : 'bg-emerald-400 animate-pulse'"
                    ></div>
                    <div class="min-w-0">
                      <p class="font-bold text-slate-200 text-xs truncate">Battle #{{ match.id.slice(0, 6) }}</p>
                      <p class="text-[9px] text-slate-500 truncate">{{ formatRelativeTime(match.created) }}</p>
                    </div>
                  </div>
                  <div class="flex items-center gap-2 flex-shrink-0">
                    <span class="text-sm font-bold" :class="getMatchResult(match) === 'win' ? 'text-emerald-400' : getMatchResult(match) === 'loss' ? 'text-red-400' : 'text-amber-400'">
                      {{ getMatchResultIcon(match) }}
                    </span>
                    <span 
                      class="pixel-badge px-2 py-0.5 text-[9px] font-bold uppercase flex-shrink-0"
                      :class="match.ended ? 'bg-slate-700 text-slate-300 border-slate-600' : 'bg-emerald-900/40 text-emerald-300 border-emerald-500/40'"
                    >
                      {{ match.ended ? 'ENDED' : 'LIVE' }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Champion Card - Enhanced -->
          <div class="pixel-box bg-slate-900/90" data-testid="champion-card">
            <div class="p-2 md:p-4 border-b-4 border-slate-800 flex items-center gap-2 md:gap-3">
              <img :src="PIXEL_ASSETS.ICON_CROWN" alt="" class="w-5 h-5 md:w-5 md:h-5 pixelated text-amber-400" />
              <h3 class="text-sm md:text-lg font-bold text-amber-300">CHAMPION</h3>
              <span v-if="topFighter" class="text-[9px] bg-indigo-900/60 border border-indigo-500/50 px-1.5 py-0.5 text-indigo-300 font-bold hidden sm:inline">ELITE</span>
            </div>
            
            <div class="p-3 md:p-4">
              <div v-if="topFighter" class="text-center">
                <div class="relative inline-block mb-3">
                  <div class="pixel-box-sm w-16 h-16 md:w-24 md:h-24 mx-auto flex items-center justify-center bg-gradient-to-br from-indigo-900/80 to-purple-900/80 border-indigo-500/60">
                    <img :src="PIXEL_ASSETS.ICON_WARRIOR" alt="" class="w-10 h-10 md:w-16 md:h-16 pixelated" />
                  </div>
                  <div class="absolute -bottom-2 -right-2 bg-gradient-to-br from-amber-600 to-amber-800 text-white text-[10px] px-1.5 py-1 font-bold border-2 border-slate-900 shadow-lg">
                    LVL {{ topFighter.level }}
                  </div>
                </div>
                <h4 class="text-sm md:text-lg font-bold text-white truncate max-w-full mb-2" data-testid="champion-name">{{ topFighter.name }}</h4>
                <div class="flex flex-wrap gap-1.5 justify-center mb-3">
                  <span class="pixel-badge bg-slate-800 text-slate-300 px-2 py-1 text-[10px] border-slate-600 inline-flex items-center gap-1">
                    <span class="w-2 h-2 rounded-full bg-pink-500"></span>
                    Power: {{ topFighter.power }}
                  </span>
                  <span class="pixel-badge bg-slate-800 text-slate-300 px-2 py-1 text-[10px] border-slate-600 inline-flex items-center gap-1">
                    <span class="w-2 h-2 rounded-full bg-emerald-500"></span>
                    Win Rate: 78%
                  </span>
                </div>
                <div class="w-full bg-slate-800/60 h-1.5 rounded-full overflow-hidden">
                  <div class="h-full bg-gradient-to-r from-amber-600 to-amber-400" :style="{ width: `${Math.min(100, topFighter.level * 5)}%` }"></div>
                </div>
              </div>
              
              <div v-else class="py-4 md:py-8 text-center">
                <p class="text-slate-500 text-xs md:text-sm mb-3">No champions yet</p>
                <router-link to="/roster" class="rpg-btn-small inline-flex text-xs py-2 px-4 min-h-[44px]">
                  RECRUIT
                </router-link>
              </div>
            </div>
          </div>

         <!-- Quick Actions - Compact -->
         <div class="pixel-box bg-slate-900/90">
            <div class="p-2 md:p-4 border-b-4 border-slate-800">
              <h3 class="text-xs font-bold text-slate-500 uppercase tracking-wider">Quick Actions</h3>
            </div>
           
           <div class="p-2 md:p-4 grid grid-cols-2 gap-2">
             <router-link 
               to="/matches" 
               class="pixel-box-sm bg-slate-800/60 p-2 flex flex-col items-center justify-center hover:bg-slate-800/80 hover:border-indigo-500/50 transition-all group min-h-[64px]"
               data-testid="quick-battle"
             >
               <img :src="PIXEL_ASSETS.ICON_SWORDS" alt="" class="w-5 h-5 pixelated mb-1 group-hover:scale-110 transition-transform" />
               <span class="text-[10px] font-bold text-slate-300">BATTLE</span>
             </router-link>
             
             <router-link 
               to="/roster" 
               class="pixel-box-sm bg-slate-800/60 p-2 flex flex-col items-center justify-center hover:bg-slate-800/80 hover:border-purple-500/50 transition-all group min-h-[64px]"
               data-testid="quick-roster"
             >
               <img :src="PIXEL_ASSETS.ICON_KNIGHT" alt="" class="w-5 h-5 pixelated mb-1 group-hover:scale-110 transition-transform" />
               <span class="text-[10px] font-bold text-slate-300">ROSTER</span>
             </router-link>
             
             <router-link 
               to="/inventory" 
               class="pixel-box-sm bg-slate-800/60 p-2 flex flex-col items-center justify-center hover:bg-slate-800/80 hover:border-emerald-500/50 transition-all group min-h-[64px]"
               data-testid="quick-vault"
             >
               <img :src="PIXEL_ASSETS.ICON_CHEST" alt="" class="w-5 h-5 pixelated mb-1 group-hover:scale-110 transition-transform" />
               <span class="text-[10px] font-bold text-slate-300">VAULT</span>
             </router-link>
             
             <router-link 
               to="/leagues" 
               class="pixel-box-sm bg-slate-800/60 p-2 flex flex-col items-center justify-center hover:bg-slate-800/80 hover:border-amber-500/50 transition-all group min-h-[64px]"
               data-testid="quick-leagues"
             >
               <img :src="PIXEL_ASSETS.ICON_TROPHY" alt="" class="w-5 h-5 pixelated mb-1 group-hover:scale-110 transition-transform" />
               <span class="text-[10px] font-bold text-slate-300">LEAGUES</span>
             </router-link>
           </div>
         </div>

      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, computed } from "vue";
import { useRouter } from "vue-router";
import { useRosterStore } from "@/features/roster/store";
import { useMatchesStore } from "@/features/matches/store";
import { useRewardsStore } from "@/features/rewards/store";
import { useLeaguesStore } from "@/features/leagues/store";
import EventBanner from "@/features/events/components/EventBanner.vue";

const router = useRouter();

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

const formatRelativeTime = (dateString: string) => {
  const now = new Date();
  const date = new Date(dateString);
  const diffMs = now.getTime() - date.getTime();
  const diffMins = Math.floor(diffMs / 60000);
  const diffHours = Math.floor(diffMs / 3600000);
  const diffDays = Math.floor(diffMs / 86400000);
  
  if (diffMins < 1) return 'Just now';
  if (diffMins < 60) return `${diffMins}m ago`;
  if (diffHours < 24) return `${diffHours}h ago`;
  return `${diffDays}d ago`;
};

const getMatchResult = (match: any): 'win' | 'loss' | 'unknown' => {
  if (!match.ended || !match.result) return 'unknown';
  return match.result === 'win' ? 'win' : 'loss';
};

const getMatchResultIcon = (match: any): string => {
  const result = getMatchResult(match);
  if (result === 'win') return '✓';
  if (result === 'loss') return '✗';
  return '●';
};

const goToMatch = (matchId: string) => {
  router.push(`/matches/${matchId}`);
};

const claimAllRewards = async () => {
   await rewardsStore.claimAll();
};

const imminentLeagues = computed(() => {
  return leaguesStore.leagues.slice(0, 2);
});

onMounted(() => {
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
