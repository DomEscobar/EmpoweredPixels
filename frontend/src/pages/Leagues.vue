<template>
  <div class="retro-rpg min-h-screen p-4 md:p-8 font-mono text-amber-100" :style="{ backgroundImage: `url('${PIXEL_ASSETS.BG_DUNGEON}')` }">
    <div class="fixed inset-0 bg-slate-950/80 pointer-events-none z-0"></div>
    
    <div class="relative z-10 max-w-7xl mx-auto space-y-8">
      <!-- Title Section -->
      <header class="flex flex-col md:flex-row md:items-end justify-between gap-6 pb-6 border-b-4 border-amber-900/50">
        <div class="flex items-center gap-4">
          <div class="w-16 h-16 bg-slate-900 border-4 border-amber-600 flex items-center justify-center shadow-lg">
            <img :src="PIXEL_ASSETS.ICON_CROWN" class="w-10 h-10 pixelated" />
          </div>
          <div>
            <h1 class="text-4xl md:text-6xl font-black tracking-tight text-amber-500 uppercase text-shadow-retro leading-none">
              WAR LEAGUES
            </h1>
            <p class="text-amber-200/60 mt-2 font-bold text-xs uppercase tracking-widest">
              <span class="text-red-500">>></span> Eternal glory awaits the champions
            </p>
          </div>
        </div>
        
        <div class="flex items-center gap-4">
          <div class="bg-slate-900 border-2 border-slate-700 px-4 py-2">
            <div class="text-[10px] text-slate-500 uppercase tracking-wider">Active Campaigns</div>
            <div class="text-2xl font-black text-emerald-400">{{ leaguesStore.activeLeagueCount }}</div>
          </div>
        </div>
      </header>

      <!-- Loading State -->
      <div v-if="leaguesStore.isLoading" class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <div v-for="i in 3" :key="i" class="h-80 bg-slate-900/50 border-4 border-slate-800 animate-pulse flex items-center justify-center">
          <span class="text-slate-700 uppercase font-bold text-xs">Loading Campaigns...</span>
        </div>
      </div>

      <!-- Error State -->
      <div v-else-if="leaguesStore.error" class="py-20 text-center border-4 border-dashed border-red-900/50 bg-red-950/20">
        <img :src="PIXEL_ASSETS.ICON_SKULL" class="w-16 h-16 mx-auto opacity-50 pixelated mb-4" />
        <h3 class="text-red-400 uppercase font-bold tracking-widest">Connection Lost</h3>
        <p class="text-slate-600 text-xs mt-2">{{ leaguesStore.error }}</p>
        <button @click="leaguesStore.fetchLeagues()" class="mt-4 rpg-btn-small bg-red-900 border-red-700 text-red-200 hover:bg-red-800">
          Retry
        </button>
      </div>

      <!-- Empty State -->
      <div v-else-if="!leaguesStore.leagues.length" class="py-20 text-center border-4 border-dashed border-slate-800 bg-slate-900/20">
        <img :src="PIXEL_ASSETS.ICON_CROWN" class="w-16 h-16 mx-auto opacity-20 pixelated grayscale mb-4" />
        <h3 class="text-slate-500 uppercase font-bold tracking-widest">No Active Campaigns</h3>
        <p class="text-slate-600 text-xs mt-2">The war council has not declared any leagues</p>
      </div>

      <!-- League Grid -->
      <div v-else class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <div 
          v-for="league in leaguesStore.leagues" 
          :key="league.id"
          class="group rpg-card bg-slate-900 border-4 transition-all duration-200 flex flex-col relative overflow-hidden cursor-pointer"
          :class="isSubscribed(league.id) 
            ? 'border-emerald-600 hover:border-emerald-500 ring-2 ring-emerald-500/20' 
            : 'border-slate-800 hover:border-amber-600/50'"
          @click="openLeagueDetail(league)"
        >
          <!-- Active Subscription Badge -->
          <div v-if="isSubscribed(league.id)" class="absolute top-0 right-0 z-20">
            <div class="bg-emerald-600 text-emerald-100 text-[10px] font-bold uppercase tracking-widest px-3 py-1 flex items-center gap-1">
              <span class="w-2 h-2 bg-emerald-300 rounded-full animate-pulse"></span>
              ACTIVE
            </div>
          </div>

          <!-- Card Top Decoration -->
          <div class="h-2 transition-colors" :class="isSubscribed(league.id) ? 'bg-emerald-600' : 'bg-slate-800 group-hover:bg-amber-600/50'"></div>
          
          <!-- League Icon Banner -->
          <div class="relative h-24 bg-gradient-to-b from-slate-800 to-slate-900 flex items-center justify-center overflow-hidden">
            <div class="absolute inset-0 opacity-10" :style="{ backgroundImage: `url('${PIXEL_ASSETS.BG_DUNGEON}')`, backgroundSize: '64px' }"></div>
            <div class="w-16 h-16 bg-slate-950 border-4 border-amber-700 flex items-center justify-center relative z-10 transform group-hover:scale-110 transition-transform">
              <img :src="getLeagueIcon(league)" class="w-10 h-10 pixelated" />
            </div>
          </div>

          <div class="p-5 flex flex-col flex-1 relative z-10">
            <div class="mb-4">
              <h3 class="text-xl font-black text-amber-100 uppercase leading-tight group-hover:text-amber-400 transition-colors">
                {{ league.name }}
              </h3>
              <p class="text-xs text-slate-500 mt-2 leading-relaxed line-clamp-2">
                {{ getLeagueDescription(league) }}
              </p>
            </div>

            <!-- Stats Grid -->
            <div class="grid grid-cols-2 gap-3 mb-5">
              <div class="bg-slate-950/50 p-3 border border-slate-800">
                <div class="text-[10px] text-slate-600 uppercase tracking-wider">Combatants</div>
                <div class="text-lg font-black text-slate-200">{{ leaguesStore.getParticipantCount(league.id) || '—' }}</div>
              </div>
              <div class="bg-slate-950/50 p-3 border border-slate-800">
                <div class="text-[10px] text-slate-600 uppercase tracking-wider">Tier</div>
                <div class="text-lg font-black" :class="getTierColor(league)">{{ getLeagueTier(league) }}</div>
              </div>
            </div>

            <!-- Prize Pool -->
            <div class="bg-amber-900/20 border border-amber-800/30 p-3 mb-5">
              <div class="flex items-center justify-between">
                <span class="text-[10px] text-amber-700 uppercase tracking-wider">Spoils of War</span>
                <span class="text-amber-400 font-black">{{ getLeaguePrize(league) }}</span>
              </div>
            </div>

            <!-- Subscribed Fighters -->
            <div v-if="getMyFighters(league.id).length" class="mb-4">
              <div class="text-[10px] text-emerald-600 uppercase tracking-wider mb-2">Your Champions</div>
              <div class="flex flex-wrap gap-2">
                <div 
                  v-for="sub in getMyFighters(league.id)" 
                  :key="sub.fighterId"
                  class="bg-emerald-900/30 border border-emerald-700/50 px-2 py-1 text-xs text-emerald-300 font-bold flex items-center gap-1"
                >
                  <img :src="PIXEL_ASSETS.ICON_FIGHTER" class="w-4 h-4 pixelated" />
                  {{ getFighterName(sub.fighterId) }}
                </div>
              </div>
            </div>

            <!-- Actions -->
            <div class="mt-auto space-y-2">
              <button 
                @click.stop="openLeagueDetail(league)"
                class="w-full rpg-btn-small bg-slate-800 border-slate-600 text-slate-300 hover:bg-slate-700 font-bold uppercase"
              >
                View Campaign
              </button>
              <button 
                v-if="!isSubscribed(league.id)"
                @click.stop="openSubscribeModal(league)"
                class="w-full rpg-btn-small bg-amber-600 border-amber-800 text-slate-900 hover:bg-amber-500 font-black uppercase"
              >
                Enlist Fighter
              </button>
              <button 
                v-else
                @click.stop="openManageModal(league)"
                class="w-full rpg-btn-small bg-emerald-700 border-emerald-900 text-emerald-100 hover:bg-emerald-600 font-bold uppercase"
              >
                Manage Squad
              </button>
            </div>
          </div>

          <!-- BG Texture -->
          <div class="absolute inset-0 opacity-5 pointer-events-none" :style="{ backgroundImage: `url('${PIXEL_ASSETS.BG_DUNGEON}')`, backgroundSize: '64px' }"></div>
        </div>
      </div>
    </div>

    <!-- Subscribe Modal -->
    <BaseModal :show="showSubscribeModal" @close="showSubscribeModal = false">
      <template #title>
        <div class="flex items-center gap-3 text-amber-500">
          <img :src="PIXEL_ASSETS.ICON_SCROLL" class="w-6 h-6 pixelated" />
          <span class="uppercase font-black text-xl tracking-wide">Enlist Champion</span>
        </div>
      </template>
      <div class="space-y-4 font-mono text-slate-200">
        <div v-if="selectedLeagueForAction" class="bg-slate-950 border-2 border-amber-900/50 p-4 mb-4">
          <div class="text-[10px] text-amber-700 uppercase tracking-wider mb-1">Campaign</div>
          <div class="text-lg font-black text-amber-400">{{ selectedLeagueForAction.name }}</div>
        </div>

        <div v-if="!rosterStore.fighters.length" class="text-center py-8">
          <img :src="PIXEL_ASSETS.ICON_SKULL" class="w-12 h-12 mx-auto opacity-30 pixelated mb-4" />
          <p class="text-slate-500 uppercase text-xs font-bold">No fighters available</p>
          <p class="text-slate-600 text-xs mt-1">Recruit a champion in the Armory first</p>
        </div>

        <div v-else class="space-y-2 max-h-[400px] overflow-y-auto custom-scrollbar pr-2">
          <div
            v-for="fighter in availableFighters"
            :key="fighter.id"
            class="relative flex items-center gap-4 p-3 border-2 cursor-pointer transition-all duration-200 group bg-slate-900"
            :class="selectedFighterId === fighter.id ? 'border-amber-500 bg-amber-900/20' : 'border-slate-800 hover:border-slate-600 hover:bg-slate-800'"
            @click="selectedFighterId = fighter.id"
          >
            <div class="w-12 h-12 bg-black border-2 border-slate-700 flex items-center justify-center">
              <img :src="PIXEL_ASSETS.ICON_FIGHTER" class="w-10 h-10 pixelated object-cover" />
            </div>
            <div class="flex-1">
              <h4 class="font-bold text-white uppercase tracking-wider" :class="selectedFighterId === fighter.id ? 'text-amber-300' : ''">
                {{ fighter.name }}
              </h4>
              <div class="text-[10px] text-slate-500 uppercase mt-1">Lvl {{ fighter.level ?? 1 }}</div>
            </div>
            <div v-if="selectedFighterId === fighter.id" class="text-amber-400 font-bold text-lg">
              <<
            </div>
          </div>
        </div>

        <div class="flex justify-between gap-3 pt-6 border-t-2 border-slate-800 border-dashed">
          <button type="button" @click="showSubscribeModal = false" class="px-4 py-2 text-xs uppercase font-bold text-slate-500 hover:text-slate-300">
            Cancel
          </button>
          <button 
            :disabled="!selectedFighterId || leaguesStore.isSubscribing" 
            @click="confirmSubscribe"
            class="rpg-btn bg-amber-600 border-amber-800 text-slate-900 hover:bg-amber-500 font-black uppercase tracking-wider px-6 py-2 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ leaguesStore.isSubscribing ? 'Enlisting...' : 'Enlist' }}
          </button>
        </div>
      </div>
    </BaseModal>

    <!-- Manage Squad Modal -->
    <BaseModal :show="showManageModal" @close="showManageModal = false">
      <template #title>
        <div class="flex items-center gap-3 text-emerald-400">
          <img :src="PIXEL_ASSETS.ICON_SHIELD" class="w-6 h-6 pixelated" />
          <span class="uppercase font-black text-xl tracking-wide">Manage Squad</span>
        </div>
      </template>
      <div class="space-y-4 font-mono text-slate-200">
        <div v-if="selectedLeagueForAction" class="bg-slate-950 border-2 border-emerald-900/50 p-4 mb-4">
          <div class="text-[10px] text-emerald-700 uppercase tracking-wider mb-1">Campaign</div>
          <div class="text-lg font-black text-emerald-400">{{ selectedLeagueForAction.name }}</div>
        </div>

        <!-- Current Subscriptions -->
        <div class="space-y-3">
          <div class="text-[10px] text-slate-500 uppercase tracking-wider">Enlisted Champions</div>
          <div 
            v-for="sub in currentSubscriptions" 
            :key="sub.fighterId"
            class="flex items-center gap-4 p-3 bg-slate-900 border-2 border-slate-800"
          >
            <div class="w-12 h-12 bg-black border-2 border-emerald-700 flex items-center justify-center">
              <img :src="PIXEL_ASSETS.ICON_FIGHTER" class="w-10 h-10 pixelated object-cover" />
            </div>
            <div class="flex-1">
              <h4 class="font-bold text-emerald-300 uppercase tracking-wider">{{ getFighterName(sub.fighterId) }}</h4>
            </div>
            <button 
              @click="handleUnsubscribe(sub.fighterId)"
              :disabled="leaguesStore.isSubscribing"
              class="rpg-btn-small bg-red-900/50 border-red-800 text-red-300 hover:bg-red-900 hover:text-white"
            >
              Withdraw
            </button>
          </div>
        </div>

        <!-- Add More -->
        <div v-if="availableFighters.length" class="pt-4 border-t border-slate-800">
          <div class="text-[10px] text-slate-500 uppercase tracking-wider mb-3">Enlist Additional Champion</div>
          <div class="flex gap-3">
            <select 
              v-model="selectedFighterId" 
              class="flex-1 bg-slate-950 border-2 border-slate-700 p-2 text-amber-100 uppercase text-xs font-bold focus:outline-none focus:border-amber-500"
            >
              <option value="">Select Fighter</option>
              <option v-for="f in availableFighters" :key="f.id" :value="f.id">{{ f.name }}</option>
            </select>
            <button 
              @click="confirmSubscribe"
              :disabled="!selectedFighterId || leaguesStore.isSubscribing"
              class="rpg-btn-small bg-amber-600 border-amber-800 text-slate-900 hover:bg-amber-500 font-black disabled:opacity-50"
            >
              Enlist
            </button>
          </div>
        </div>

        <div class="flex justify-end gap-3 pt-6 border-t-2 border-slate-800 border-dashed">
          <button type="button" @click="showManageModal = false" class="rpg-btn-small bg-slate-800 border-slate-600 text-slate-300 hover:bg-slate-700 font-bold uppercase">
            Done
          </button>
        </div>
      </div>
    </BaseModal>

    <!-- League Detail Modal -->
    <BaseModal :show="showDetailModal" @close="closeDetailModal">
      <template #title>
        <div class="flex items-center gap-3 text-amber-500">
          <img :src="PIXEL_ASSETS.ICON_CROWN" class="w-6 h-6 pixelated" />
          <span class="uppercase font-black text-xl tracking-wide">{{ leaguesStore.selectedLeague?.name ?? 'Campaign Intel' }}</span>
        </div>
      </template>
      <div class="font-mono text-slate-200">
        <!-- Tabs -->
        <div class="flex border-b-2 border-slate-800 mb-6">
          <button 
            v-for="tab in detailTabs" 
            :key="tab.id"
            @click="activeDetailTab = tab.id"
            class="px-4 py-3 text-xs uppercase font-bold tracking-wider transition-all border-b-2 -mb-[2px]"
            :class="activeDetailTab === tab.id 
              ? 'text-amber-400 border-amber-500 bg-amber-900/10' 
              : 'text-slate-500 border-transparent hover:text-slate-300'"
          >
            {{ tab.label }}
          </button>
        </div>

        <!-- Loading -->
        <div v-if="leaguesStore.isLoadingDetail" class="py-12 text-center">
          <div class="w-8 h-8 border-4 border-amber-600 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
          <p class="text-slate-500 uppercase text-xs">Gathering intelligence...</p>
        </div>

        <div v-else-if="leaguesStore.selectedLeague">
          <!-- Participants Tab -->
          <div v-if="activeDetailTab === 'participants'" class="space-y-4">
            <div class="flex items-center justify-between mb-4">
              <span class="text-[10px] text-slate-500 uppercase tracking-wider">Registered Warriors</span>
              <span class="text-amber-400 font-bold">{{ leaguesStore.selectedLeague.subscriptions?.length ?? 0 }}</span>
            </div>
            
            <div v-if="!leaguesStore.selectedLeague.subscriptions?.length" class="text-center py-8 bg-slate-950/50 border border-slate-800">
              <p class="text-slate-600 uppercase text-xs">No warriors have enlisted yet</p>
            </div>
            
            <div v-else class="space-y-2 max-h-[300px] overflow-y-auto custom-scrollbar">
              <div 
                v-for="(sub, idx) in leaguesStore.selectedLeague.subscriptions" 
                :key="sub.fighterId"
                class="flex items-center gap-4 p-3 bg-slate-900 border border-slate-800"
              >
                <div class="w-8 h-8 bg-slate-950 border border-slate-700 flex items-center justify-center text-xs font-bold text-slate-500">
                  {{ idx + 1 }}
                </div>
                <div class="w-10 h-10 bg-black border-2 border-slate-700 flex items-center justify-center">
                  <img :src="PIXEL_ASSETS.ICON_FIGHTER" class="w-8 h-8 pixelated" />
                </div>
                <div class="flex-1">
                  <h4 class="font-bold text-slate-200 uppercase text-sm">{{ getFighterName(sub.fighterId) || sub.fighterId.substring(0, 8) }}</h4>
                </div>
                <div v-if="isMyFighter(sub.fighterId)" class="text-emerald-400 text-[10px] uppercase font-bold">
                  YOURS
                </div>
              </div>
            </div>
          </div>

          <!-- Matches Tab -->
          <div v-if="activeDetailTab === 'matches'" class="space-y-4">
            <div v-if="!currentLeagueMatches.length" class="text-center py-8 bg-slate-950/50 border border-slate-800">
              <img :src="PIXEL_ASSETS.ICON_SWORDS" class="w-12 h-12 mx-auto opacity-20 pixelated mb-4" />
              <p class="text-slate-600 uppercase text-xs">No battles recorded yet</p>
            </div>
            
            <div v-else class="space-y-2 max-h-[300px] overflow-y-auto custom-scrollbar">
              <router-link 
                v-for="match in currentLeagueMatches" 
                :key="match.matchId"
                :to="'/matches/' + match.matchId"
                class="flex items-center gap-4 p-3 bg-slate-900 border border-slate-800 hover:border-amber-600/50 transition-colors group"
              >
                <div class="w-10 h-10 bg-slate-950 border-2 border-slate-700 flex items-center justify-center group-hover:border-amber-600">
                  <img :src="PIXEL_ASSETS.ICON_SWORDS" class="w-6 h-6 pixelated" />
                </div>
                <div class="flex-1">
                  <h4 class="font-bold text-slate-300 uppercase text-sm group-hover:text-amber-400">
                    Battle #{{ match.matchId.substring(0, 8) }}
                  </h4>
                </div>
                <div class="text-slate-600 text-xs uppercase">View >></div>
              </router-link>
            </div>
          </div>

          <!-- Leaderboard Tab -->
          <div v-if="activeDetailTab === 'leaderboard'" class="space-y-4">
            <div v-if="!currentLeagueHighscores.length" class="text-center py-8 bg-slate-950/50 border border-slate-800">
              <img :src="PIXEL_ASSETS.ICON_TROPHY" class="w-12 h-12 mx-auto opacity-20 pixelated mb-4" />
              <p class="text-slate-600 uppercase text-xs">Rankings will appear after battles</p>
            </div>
            
            <div v-else class="space-y-2 max-h-[300px] overflow-y-auto custom-scrollbar">
              <div 
                v-for="(score, idx) in currentLeagueHighscores" 
                :key="score.fighterId"
                class="flex items-center gap-4 p-3 border"
                :class="idx === 0 ? 'bg-amber-900/20 border-amber-700' : idx === 1 ? 'bg-slate-400/10 border-slate-600' : idx === 2 ? 'bg-orange-900/20 border-orange-800' : 'bg-slate-900 border-slate-800'"
              >
                <div 
                  class="w-8 h-8 flex items-center justify-center text-sm font-black"
                  :class="idx === 0 ? 'text-amber-400' : idx === 1 ? 'text-slate-300' : idx === 2 ? 'text-orange-400' : 'text-slate-500'"
                >
                  {{ idx + 1 }}
                </div>
                <div class="w-10 h-10 bg-black border-2 flex items-center justify-center" :class="idx === 0 ? 'border-amber-600' : 'border-slate-700'">
                  <img :src="PIXEL_ASSETS.ICON_FIGHTER" class="w-8 h-8 pixelated" />
                </div>
                <div class="flex-1">
                  <h4 class="font-bold uppercase text-sm" :class="idx === 0 ? 'text-amber-300' : 'text-slate-200'">
                    {{ score.fighterName || score.fighterId.substring(0, 8) }}
                  </h4>
                  <div class="text-[10px] text-slate-600">{{ score.username }}</div>
                </div>
                <div class="text-right">
                  <div class="text-lg font-black" :class="idx === 0 ? 'text-amber-400' : 'text-slate-300'">{{ score.score }}</div>
                  <div class="text-[10px] text-slate-600 uppercase">Points</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </BaseModal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import BaseModal from '@/shared/ui/BaseModal.vue';
import { useLeaguesStore } from '@/features/leagues/store';
import { useRosterStore } from '@/features/roster/store';
import type { League } from '@/features/leagues/api';

const PIXEL_ASSETS = {
  BG_DUNGEON: 'https://vibemedia.space/bg_dungeon_v2_99283.png?prompt=dark%20dungeon%20floor%20tile%20texture%20seamless&style=pixel_game_asset&key=NOGON',
  ICON_CROWN: 'https://vibemedia.space/icon_crown_v2_88912.png?prompt=golden%20royal%20crown%20pixel%20art&style=pixel_game_asset&key=NOGON',
  ICON_SWORDS: 'https://vibemedia.space/icon_swords_v2_11234.png?prompt=crossed%20steel%20swords%20pixel%20art&style=pixel_game_asset&key=NOGON',
  ICON_TROPHY: 'https://vibemedia.space/trophy_icon_4d5e6f_v1.png?prompt=golden%20trophy%20pixel%20art&style=pixel_game_asset&key=NOGON',
  ICON_SCROLL: 'https://vibemedia.space/icon_scroll_v2_66789.png?prompt=ancient%20magic%20scroll%20pixel%20art&style=pixel_game_asset&key=NOGON',
  ICON_SHIELD: 'https://vibemedia.space/icon_shield_v2_77432.png?prompt=knight%20shield%20pixel%20art&style=pixel_game_asset&key=NOGON',
  ICON_SKULL: 'https://vibemedia.space/icon_skull_v2_55432.png?prompt=skull%20icon%20pixel%20art&style=pixel_game_asset&key=NOGON',
  ICON_FIGHTER: 'https://vibemedia.space/fighter_hooded_8p9q0r_v1.png?prompt=mystery%20hooded%20figure%20pixel%20art&style=pixel_game_asset&key=NOGON',
};

const leaguesStore = useLeaguesStore();
const rosterStore = useRosterStore();

const showSubscribeModal = ref(false);
const showManageModal = ref(false);
const showDetailModal = ref(false);
const selectedLeagueForAction = ref<League | null>(null);
const selectedFighterId = ref<string>('');
const activeDetailTab = ref('participants');

const detailTabs = [
  { id: 'participants', label: 'Warriors' },
  { id: 'matches', label: 'Battles' },
  { id: 'leaderboard', label: 'Rankings' },
];

const isSubscribed = (leagueId: number) => leaguesStore.isSubscribedToLeague(leagueId);

const getMyFighters = (leagueId: number) => leaguesStore.getSubscribedFighters(leagueId);

const isMyFighter = (fighterId: string) => {
  return rosterStore.fighters.some(f => f.id === fighterId);
};

const getFighterName = (fighterId: string) => {
  const fighter = rosterStore.fighters.find(f => f.id === fighterId);
  return fighter?.name ?? fighterId.substring(0, 8);
};

const currentSubscriptions = computed(() => {
  if (!selectedLeagueForAction.value) return [];
  return leaguesStore.subscriptions[selectedLeagueForAction.value.id] ?? [];
});

const availableFighters = computed(() => {
  const subscribedIds = new Set(currentSubscriptions.value.map(s => s.fighterId));
  return rosterStore.fighters.filter(f => !subscribedIds.has(f.id));
});

const currentLeagueMatches = computed(() => {
  if (!leaguesStore.selectedLeague) return [];
  return leaguesStore.leagueMatches[leaguesStore.selectedLeague.id] ?? [];
});

const currentLeagueHighscores = computed(() => {
  if (!leaguesStore.selectedLeague) return [];
  return leaguesStore.leagueHighscores[leaguesStore.selectedLeague.id] ?? [];
});

function getLeagueIcon(league: League) {
  const tier = league.options?.tier?.toLowerCase() ?? '';
  if (tier === 'mythic') return PIXEL_ASSETS.ICON_CROWN;
  if (tier === 'legendary') return PIXEL_ASSETS.ICON_TROPHY;
  return PIXEL_ASSETS.ICON_SWORDS;
}

function getLeagueDescription(league: League) {
  return league.options?.description ?? 'A fierce competition awaits those brave enough to enter.';
}

function getLeagueTier(league: League) {
  return league.options?.tier ?? 'Standard';
}

function getTierColor(league: League) {
  const tier = league.options?.tier?.toLowerCase() ?? '';
  if (tier === 'mythic') return 'text-purple-400';
  if (tier === 'legendary') return 'text-amber-400';
  if (tier === 'epic') return 'text-indigo-400';
  if (tier === 'rare') return 'text-blue-400';
  return 'text-slate-400';
}

function getLeaguePrize(league: League) {
  return league.options?.prizePool ?? 'Glory & Honor';
}

function openSubscribeModal(league: League) {
  selectedLeagueForAction.value = league;
  selectedFighterId.value = rosterStore.fighters[0]?.id ?? '';
  showSubscribeModal.value = true;
}

function openManageModal(league: League) {
  selectedLeagueForAction.value = league;
  selectedFighterId.value = '';
  showManageModal.value = true;
}

async function openLeagueDetail(league: League) {
  selectedLeagueForAction.value = league;
  activeDetailTab.value = 'participants';
  showDetailModal.value = true;
  await Promise.all([
    leaguesStore.fetchLeagueDetail(league.id),
    leaguesStore.fetchLeagueMatches(league.id),
    leaguesStore.fetchLeagueHighscores(league.id),
  ]);
}

function closeDetailModal() {
  showDetailModal.value = false;
  leaguesStore.clearSelectedLeague();
}

async function confirmSubscribe() {
  if (!selectedLeagueForAction.value || !selectedFighterId.value) return;
  
  try {
    await leaguesStore.subscribe(selectedLeagueForAction.value.id, selectedFighterId.value);
    selectedFighterId.value = '';
    
    if (!showManageModal.value) {
      showSubscribeModal.value = false;
    }
  } catch (e) {
    console.error('Failed to subscribe', e);
  }
}

async function handleUnsubscribe(fighterId: string) {
  if (!selectedLeagueForAction.value) return;
  
  try {
    await leaguesStore.unsubscribe(selectedLeagueForAction.value.id, fighterId);
    
    if (!currentSubscriptions.value.length) {
      showManageModal.value = false;
    }
  } catch (e) {
    console.error('Failed to unsubscribe', e);
  }
}

watch(activeDetailTab, async (tab) => {
  if (!leaguesStore.selectedLeague) return;
  const leagueId = leaguesStore.selectedLeague.id;
  
  if (tab === 'matches' && !leaguesStore.leagueMatches[leagueId]?.length) {
    await leaguesStore.fetchLeagueMatches(leagueId);
  }
  if (tab === 'leaderboard' && !leaguesStore.leagueHighscores[leagueId]?.length) {
    await leaguesStore.fetchLeagueHighscores(leagueId);
  }
});

onMounted(async () => {
  await rosterStore.fetchFighters();
  await leaguesStore.fetchLeagues();
  
  for (const league of leaguesStore.leagues) {
    await leaguesStore.fetchAllSubscriptions(league.id);
  }
});
</script>

<style scoped>
.retro-rpg {
  image-rendering: pixelated;
  background-repeat: repeat;
  background-size: 128px;
}

.pixelated {
  image-rendering: pixelated;
}

.text-shadow-retro {
  text-shadow: 3px 3px 0 #000, -1px -1px 0 #000;
}

.rpg-btn {
  border-bottom-width: 4px;
  border-right-width: 4px;
  transition: all 0.1s;
}

.rpg-btn:active {
  border-bottom-width: 0px;
  border-right-width: 0px;
  transform: translate(4px, 4px);
}

.rpg-btn-small {
  padding: 0.5rem 1rem;
  font-size: 0.75rem;
  border-width: 2px;
  transition: all 0.1s;
}

.rpg-btn-small:hover {
  filter: brightness(1.1);
}

.rpg-btn-small:active {
  transform: translateY(1px);
}

.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}

.custom-scrollbar::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.2);
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(251, 191, 36, 0.3);
  border-radius: 3px;
}

.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: rgba(251, 191, 36, 0.5);
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
