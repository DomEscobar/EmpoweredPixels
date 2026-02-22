<template>
  <div class="leagues-list" data-testid="leagues-list">
    <!-- Loading State -->
    <div v-if="isLoading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <div v-for="i in 3" :key="i" class="h-80 bg-slate-900/50 border-4 border-slate-800 animate-pulse flex items-center justify-center">
        <span class="text-slate-700 uppercase font-bold text-xs">Loading Campaigns...</span>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="py-20 text-center border-4 border-dashed border-red-900/50 bg-red-950/20">
      <img :src="PIXEL_ASSETS.ICON_SKULL" class="w-16 h-16 mx-auto opacity-50 pixelated mb-4" />
      <h3 class="text-red-400 uppercase font-bold tracking-widest">Connection Lost</h3>
      <p class="text-slate-600 text-xs mt-2">{{ error }}</p>
      <button @click="$emit('retry')" class="mt-4 rpg-btn-small bg-red-900 border-red-700 text-red-200 hover:bg-red-800">
        Retry
      </button>
    </div>

    <!-- Empty State -->
    <div v-else-if="!leagues.length" class="py-20 text-center border-4 border-dashed border-slate-800 bg-slate-900/20">
      <img :src="PIXEL_ASSETS.ICON_CROWN" class="w-16 h-16 mx-auto opacity-20 pixelated grayscale mb-4" />
      <h3 class="text-slate-500 uppercase font-bold tracking-widest">No Active Campaigns</h3>
      <p class="text-slate-600 text-xs mt-2">The war council has not declared any leagues</p>
    </div>

    <!-- League Grid -->
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <div 
         v-for="league in leagues" 
         :key="league.id"
         class="group rpg-card bg-slate-900 border-4 transition-all duration-200 flex flex-col relative overflow-hidden cursor-pointer"
         :class="isSubscribed(league.id) 
           ? 'border-emerald-600 hover:border-emerald-500 ring-2 ring-emerald-500/20' 
           : 'border-slate-800 hover:border-amber-600/50'"
         :data-testid="'league-card-' + league.id"
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

        <div class="flex-1 flex flex-col relative z-10">
          <div class="mb-4">
            <h3 class="text-lg sm:text-xl font-black text-amber-100 uppercase leading-tight group-hover:text-amber-400 transition-colors">
              {{ league.name }}
            </h3>
             <p class="text-[10px] sm:text-xs text-slate-500 mt-2 leading-relaxed line-clamp-2">
              {{ getLeagueDescription(league) }}
            </p>
          </div>

          <!-- Stats Grid -->
          <div class="grid grid-cols-2 gap-3 mb-5">
            <div class="bg-slate-950/50 p-3 border border-slate-800">
              <div class="text-[10px] text-slate-600 uppercase tracking-wider">Combatants</div>
               <div class="text-base sm:text-lg font-black text-slate-200">{{ participantCount(league.id) || '—' }}</div>
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

          <!-- Last Winner -->
          <div v-if="getLastWinner(league.id)" class="bg-slate-950/50 border border-slate-800 p-3 mb-5">
            <div class="text-[10px] text-slate-600 uppercase tracking-wider mb-1">Last Victor</div>
            <div class="flex items-center gap-2">
              <img :src="PIXEL_ASSETS.ICON_TROPHY" class="w-6 h-6 pixelated text-amber-400" />
              <div class="flex-1 min-w-0">
                <div class="text-sm font-bold text-amber-300 truncate">{{ getLastWinner(league.id)?.fighterName }}</div>
                <div class="text-[10px] text-slate-500">@{{ getLastWinner(league.id)?.username }}</div>
              </div>
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
              @click.stop="$emit('view-detail', league)"
              class="w-full rpg-btn-small bg-slate-800 border-slate-600 text-slate-300 hover:bg-slate-700 font-bold uppercase"
              :data-testid="'view-detail-btn-' + league.id"
            >
              View Campaign
            </button>
            <button 
              v-if="!isSubscribed(league.id)"
              @click.stop="$emit('subscribe', league)"
              class="w-full rpg-btn-small bg-amber-600 border-amber-800 text-slate-900 hover:bg-amber-500 font-black uppercase"
              :data-testid="'subscribe-btn-' + league.id"
            >
              Enlist Fighter
            </button>
            <button 
              v-else
              @click.stop="$emit('manage', league)"
              class="w-full rpg-btn-small bg-emerald-700 border-emerald-900 text-emerald-100 hover:bg-emerald-600 font-bold uppercase"
              :data-testid="'manage-btn-' + league.id"
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
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useLeaguesStore } from './store';
import { useRosterStore } from '../roster/store';
import type { League } from './api';

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

const props = defineProps<{
  leagues: League[];
  isLoading?: boolean;
  error?: string | null;
  participantCount: (leagueId: number) => number | null;
  isSubscribed: (leagueId: number) => boolean;
  getMyFighters: (leagueId: number) => { fighterId: string }[];
  getLastWinner: (leagueId: number) => { fighterName: string; username: string } | null;
}>();

defineEmits<{
  (e: 'subscribe', league: League): void;
  (e: 'manage', league: League): void;
  (e: 'view-detail', league: League): void;
  (e: 'retry'): void;
}>();

const leaguesStore = useLeaguesStore();
const rosterStore = useRosterStore();

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

function getFighterName(fighterId: string) {
  const fighter = rosterStore.fighters.find(f => f.id === fighterId);
  return fighter?.name ?? fighterId.substring(0, 8);
}
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

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
