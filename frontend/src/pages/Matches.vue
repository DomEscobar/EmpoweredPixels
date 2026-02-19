<template>
  <div class="retro-rpg min-h-screen p-4 md:p-8 font-mono text-amber-100" :style="{ backgroundImage: `url('${PIXEL_ASSETS.BG_DUNGEON}')` }" data-testid="matches-page">
    <div class="fixed inset-0 bg-slate-950/80 pointer-events-none z-0"></div>
    
    <div class="relative z-10 max-w-7xl mx-auto space-y-6">
      <!-- Header Section -->
      <header class="flex flex-col lg:flex-row lg:items-end justify-between gap-4 pb-6 border-b-4 border-amber-900/50" data-testid="matches-header">
        <div class="flex items-center gap-4">
          <div class="w-14 h-14 md:w-16 md:h-16 bg-slate-900 border-4 border-amber-600 flex items-center justify-center shadow-lg">
            <img :src="PIXEL_ASSETS.ICON_CHEST" class="w-8 h-8 md:w-10 md:h-10 pixelated" />
          </div>
          <div>
            <h1 class="text-3xl md:text-5xl font-black tracking-tight text-amber-500 uppercase text-shadow-retro leading-none">
              QUEST BOARD
            </h1>
            <p class="text-amber-200/60 mt-1 md:mt-2 font-bold text-xs uppercase tracking-widest hidden md:block">
              <span class="text-red-500">>></span> Find your glory in the arena
            </p>
          </div>
        </div>
        
        <div class="flex flex-col gap-3">
          <div class="pixel-box-sm bg-slate-800/80 px-3 py-1.5 flex items-center gap-2 self-end">
            <span class="w-2 h-2 rounded-full bg-green-500 animate-pulse"></span>
            <span class="text-xs text-slate-400">Online:</span>
            <span class="text-sm font-bold text-green-400">{{ onlinePlayers }}</span>
          </div>
          <div class="flex gap-2 md:gap-3">
            <button 
              @click="handleQuickJoin"
              :disabled="isQuickJoining || roster.fighters.length === 0"
              class="flex-1 md:flex-none rpg-btn bg-blue-600 border-blue-800 text-white hover:bg-blue-500 disabled:opacity-50 disabled:cursor-not-allowed group relative px-4 md:px-6 py-3 font-bold uppercase tracking-wider flex items-center justify-center gap-2"
              data-testid="quick-join-button"
            >
              <span v-if="isQuickJoining" class="animate-spin">⚡</span>
              <span v-else>⚡ Quick Join</span>
            </button>
            <button 
              @click="showCreate = true"
              class="rpg-btn bg-emerald-700 border-emerald-900 text-white hover:bg-emerald-600 group relative px-4 md:px-6 py-3 font-bold uppercase tracking-wider flex items-center gap-2"
              data-testid="create-match-button"
            >
              <span class="absolute inset-0 border-2 border-white/20 pointer-events-none"></span>
              <span class="hidden md:inline"><img :src="PIXEL_ASSETS.ICON_SCROLL" class="w-5 h-5 pixelated" /></span>
              <span>New</span>
            </button>
          </div>
        </div>
      </header>

      <!-- Active Quest Banner -->
      <transition
        enter-active-class="transition ease-out duration-300"
        enter-from-class="transform opacity-0 scale-95"
        enter-to-class="transform opacity-100 scale-100"
        leave-active-class="transition ease-in duration-200"
        leave-from-class="transform opacity-100 scale-100"
        leave-to-class="transform opacity-0 scale-95"
      >
        <div v-if="currentMatchId" class="rpg-panel border-4 border-amber-600 bg-slate-900 shadow-xl p-1 relative overflow-hidden" data-testid="active-match-banner">
          <div class="absolute top-0 left-0 w-2 h-2 bg-amber-500"></div>
          <div class="absolute top-0 right-0 w-2 h-2 bg-amber-500"></div>
          <div class="absolute bottom-0 left-0 w-2 h-2 bg-amber-500"></div>
          <div class="absolute bottom-0 right-0 w-2 h-2 bg-amber-500"></div>

          <div class="bg-slate-900/90 p-4 md:p-6 flex flex-col md:flex-row items-center justify-between gap-4 md:gap-6 border border-amber-900/30">
            <div class="flex items-center gap-4 md:gap-6">
              <div class="relative">
                <div class="w-16 h-16 md:w-20 md:h-20 bg-slate-950 border-4 border-slate-700 flex items-center justify-center overflow-hidden">
                  <img :src="statusIconImg" class="w-12 h-12 md:w-16 md:h-16 pixelated object-contain animate-bounce-slow" />
                </div>
                <div class="absolute -bottom-2 -right-2 bg-slate-900 border-2 border-amber-600 text-amber-500 text-xs font-bold px-2 py-0.5">
                  #{{ currentMatchId.substring(0, 4) }}
                </div>
              </div>
              
              <div>
                <h3 class="font-black text-xl md:text-2xl text-amber-100 uppercase tracking-wide flex flex-col">
                  <span class="text-xs text-amber-500 font-bold mb-1">Current Objective</span>
                  {{ currentMatchStatusLabel }}
                </h3>
                <div class="mt-2 flex items-center gap-2 text-sm font-bold" :class="statusTextClass">
                  <span class="w-2 h-2 bg-current animate-pulse"></span>
                  {{ statusMessage || 'Awaiting input...' }}
                </div>
              </div>
            </div>

            <div class="flex flex-col gap-2 w-full md:w-auto min-w-[180px] md:min-w-[200px]">
              <template v-if="currentMatchStatus === 'lobby'">
                <div class="flex justify-between items-center bg-black/40 px-3 py-1 border border-slate-700 mb-1">
                  <span class="text-[10px] text-slate-400 uppercase">Party Size</span>
                  <span class="text-amber-400 font-bold">{{ currentMatch?.registrations?.length ?? 1 }}/{{ options.maxPlayers || 2 }}</span>
                </div>
                <div class="grid grid-cols-2 gap-2">
                  <button @click="handleLeave" class="rpg-btn-small bg-red-900/50 border-red-800 text-red-300 hover:bg-red-900 hover:text-white" data-testid="leave-match-button">
                    Flee
                  </button>
                  <button @click="handleStart" :disabled="isStarting" class="rpg-btn-small bg-amber-600 border-amber-800 text-slate-900 hover:bg-amber-500 font-black" data-testid="begin-match-button">
                    BEGIN
                  </button>
                </div>
              </template>

              <template v-else-if="currentMatchStatus === 'running'">
                <div class="text-center mb-2">
                  <span class="text-red-500 font-bold animate-pulse text-xs uppercase tracking-widest">>> COMBAT ENGAGED <<</span>
                </div>
                <router-link :to="'/matches/' + currentMatchId" class="block">
                  <button class="w-full rpg-btn bg-red-700 border-red-900 text-white hover:bg-red-600 font-bold uppercase tracking-wider py-3">
                    Spectate
                  </button>
                </router-link>
              </template>

              <template v-else-if="currentMatchStatus === 'completed'">
                <div class="grid grid-cols-2 gap-2">
                  <button @click="clearCurrentMatch" class="rpg-btn-small bg-slate-800 border-slate-600 text-slate-400 hover:bg-slate-700">Dismiss</button>
                  <router-link :to="'/matches/' + currentMatchId" class="block">
                    <button class="w-full rpg-btn-small bg-emerald-700 border-emerald-900 text-emerald-100 hover:bg-emerald-600">Rewards</button>
                  </router-link>
                </div>
              </template>
            </div>
          </div>
        </div>
      </transition>

      <!-- Main Quest Board -->
      <div class="space-y-4 md:space-y-6">
        <!-- Compact Filter Pills -->
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div class="flex items-center gap-2 overflow-x-auto pb-2 md:pb-0 custom-scrollbar">
            <button
              v-for="stat in filterOptions"
              :key="stat.value"
              @click="browseStatus = stat.value"
              class="px-3 md:px-4 py-2 text-xs font-bold uppercase tracking-wider rounded-full transition-all whitespace-nowrap flex items-center gap-2"
              :class="browseStatus === stat.value 
                ? 'bg-amber-600 text-slate-900 shadow-lg scale-105 z-10' 
                : 'bg-slate-900 border-2 border-slate-700 text-slate-500 hover:text-amber-200 hover:border-amber-700'"
            >
              {{ stat.label }}
              <span class="text-[10px] opacity-70">({{ getCountForStatus(stat.value) }})</span>
            </button>
          </div>
          
          <div class="relative w-full md:w-56 shrink-0">
            <input
              v-model="search"
              placeholder="Search..."
              class="w-full bg-slate-900 border-2 border-slate-700 p-2 pl-3 pr-8 text-amber-100 placeholder-slate-600 focus:outline-none focus:border-amber-500 uppercase text-xs font-bold rounded"
            />
            <div class="absolute right-2 top-1/2 -translate-y-1/2 pointer-events-none opacity-50">
               🔍
            </div>
          </div>
        </div>

        <!-- Loading Skeletons -->
        <div v-if="isLoading" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          <div v-for="i in 8" :key="i" class="match-card-skeleton bg-slate-900/80 border-2 border-slate-800 rounded-lg overflow-hidden">
            <Skeleton height="80px" />
            <div class="p-4 space-y-3">
              <Skeleton width="60%" height="20px" />
              <Skeleton width="80%" height="14px" />
              <Skeleton width="40%" height="14px" />
              <div class="pt-2">
                <Skeleton height="36px" />
              </div>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <EmptyState
          v-else-if="!filteredMatches.length"
          icon="📜"
          :title="emptyStateTitle"
          :message="emptyStateMessage"
        >
          <button 
            v-if="browseStatus === 'lobby' && roster.fighters.length === 0"
            @click="$router.push('/roster')" 
            class="rpg-btn-small mt-4 bg-indigo-600 border-indigo-800 text-white hover:bg-indigo-500"
          >
            RECRUIT HERO
          </button>
          <button 
            v-else-if="browseStatus === 'lobby'"
            @click="showCreate = true" 
            class="rpg-btn-small mt-4 bg-emerald-700 border-emerald-900 text-emerald-100 hover:bg-emerald-600"
          >
            CREATE MATCH
          </button>
        </EmptyState>

        <!-- Match Cards Grid -->
        <div v-else class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4" data-testid="matches-grid">
          <div 
            v-for="match in filteredMatches" 
            :key="match.id"
            class="match-card group bg-slate-900 border-2 border-slate-800 hover:border-amber-600/50 transition-all duration-200 flex flex-col relative overflow-hidden"
            :class="{ 'border-amber-500 ring-4 ring-amber-500/20 z-10': currentMatchId === match.id }"
            :data-testid="`match-card-${match.id}`"
            @mouseenter="hoveredMatchId = match.id"
            @mouseleave="hoveredMatchId = null"
          >
            <!-- Card Hero Header -->
            <div class="h-20 relative overflow-hidden" :class="getMatchCardGradient(match.status)">
              <div class="absolute inset-0 flex items-center justify-center">
                <img :src="getMatchHeroImage(match.status)" class="w-12 h-12 pixelated opacity-60" />
              </div>
              <div class="absolute top-2 left-2">
                <span class="status-badge" :class="getStatusBadgeClass(match.status)">
                  {{ match.status }}
                </span>
              </div>
              <div class="absolute top-2 right-2">
                <span class="text-[10px] text-slate-400 uppercase font-bold">#{{ match.id.substring(0, 4) }}</span>
              </div>
              <div v-if="match.status === 'lobby' && isMatchHot(match)" class="absolute bottom-2 right-2">
                <span class="text-[10px] text-red-400 animate-pulse font-bold">🔥 HOT</span>
              </div>
            </div>

            <!-- Card Body -->
            <div class="p-4 flex flex-col flex-1">
              <h3 class="text-lg font-bold uppercase leading-tight group-hover:text-amber-400 transition-colors mb-3">
                {{ getMatchTitle(match) }}
              </h3>

              <!-- Key Metrics -->
              <div class="space-y-2 mb-4 flex-1">
                <div class="flex items-center justify-between text-xs bg-slate-950/50 p-2 rounded border border-slate-800">
                  <span class="text-slate-500 uppercase flex items-center gap-1">👥 Fighters</span>
                  <span class="text-amber-200 font-bold">{{ match.registrations?.length || 0 }}/{{ match.options?.maxFightersPerUser || 2 }}</span>
                </div>
                <div class="flex items-center justify-between text-xs bg-slate-950/50 p-2 rounded border border-slate-800">
                  <span class="text-slate-500 uppercase flex items-center gap-1">🤖 Bots</span>
                  <span class="text-amber-200 font-bold">{{ match.options?.botCount || 0 }}</span>
                </div>
                <div v-if="match.status === 'completed'" class="flex items-center justify-between text-xs bg-slate-950/50 p-2 rounded border border-slate-800">
                  <span class="text-slate-500 uppercase flex items-center gap-1">💎 Rewards</span>
                  <span class="text-emerald-400 font-bold">{{ match.totalRewards || 0 }}</span>
                </div>
              </div>

              <!-- CTA -->
              <div class="mt-auto space-y-2">
                <template v-if="(match.status ?? 'lobby') === 'lobby'">
                  <button 
                    v-if="currentMatchId !== match.id"
                    @click="openJoinModal(match.id)"
                    class="w-full rpg-btn-small bg-indigo-900 border-indigo-700 text-indigo-200 hover:bg-indigo-800 hover:text-white font-bold uppercase"
                    :data-testid="`join-match-${match.id}`"
                  >
                   Join Party
                  </button>
                  <template v-else>
                    <button 
                      @click="handleStart"
                      :disabled="isStarting"
                      class="w-full rpg-btn-small bg-emerald-700 border-emerald-900 text-emerald-100 hover:bg-emerald-600 font-bold uppercase"
                    >
                      {{ isStarting ? 'Starting...' : 'Start Battle' }}
                    </button>
                    <button 
                      @click="handleLeave"
                      class="w-full rpg-btn-small bg-slate-800 border-slate-700 text-red-400 hover:bg-slate-700 hover:text-red-300 font-bold uppercase text-xs"
                    >
                      Leave
                    </button>
                  </template>
                </template>
                <template v-else>
                  <router-link :to="'/matches/' + match.id" class="block">
                    <button 
                      class="w-full rpg-btn-small font-bold uppercase transition-all duration-300"
                      :class="rewardsStore.hasRewardForSource(match.id) 
                        ? 'bg-amber-600 border-amber-800 text-slate-900 animate-pulse hover:bg-amber-500 shadow-[0_0_15px_rgba(245,158,11,0.4)]' 
                        : 'bg-slate-800 border-slate-600 text-slate-300 hover:bg-slate-700'"
                    >
                      <template v-if="rewardsStore.hasRewardForSource(match.id)">
                        ✨ CLAIM LOOT
                      </template>
                      <template v-else>
                        {{ match.status === 'running' ? 'Watch' : 'Results' }}
                      </template>
                    </button>
                  </router-link>
                </template>
              </div>
            </div>

            <!-- Match Preview Tooltip -->
            <transition name="fade">
              <div v-if="hoveredMatchId === match.id && match.status === 'lobby'" class="absolute inset-0 bg-slate-900/95 z-20 p-4 flex flex-col">
                <h4 class="text-amber-400 font-bold uppercase text-sm mb-3">Match Preview</h4>
                <div class="flex-1 space-y-2 overflow-y-auto custom-scrollbar">
                  <div v-for="reg in match.registrations" :key="reg.fighterId" class="bg-slate-800 p-2 rounded text-xs">
                    <div class="flex items-center gap-2">
                      <img :src="PIXEL_ASSETS.ICON_FIGHTER" class="w-6 h-6 pixelated" />
                      <div>
                        <div class="text-slate-200 font-bold">{{ reg.fighterName || 'Unknown' }}</div>
                        <div class="text-slate-500">Lvl {{ reg.fighterLevel || 1 }} • PWR {{ reg.fighterPower || 0 }}</div>
                      </div>
                    </div>
                  </div>
                  <div v-if="!match.registrations?.length" class="text-slate-500 text-xs text-center py-4">
                    No fighters yet
                  </div>
                </div>
              </div>
            </transition>
          </div>
        </div>
      </div>
    </div>

    <!-- Mobile FAB -->
    <button
      v-if="!currentMatchId"
      @click="handleQuickJoin"
      :disabled="isQuickJoining || roster.fighters.length === 0"
      class="fixed bottom-6 right-6 z-50 md:hidden w-14 h-14 rounded-full bg-blue-600 border-4 border-blue-800 shadow-xl flex items-center justify-center text-white animate-bounce-slow"
      data-testid="mobile-fab"
    >
      <span v-if="isQuickJoining" class="animate-spin text-xl">⚡</span>
      <span v-else class="text-2xl">⚔️</span>
    </button>

    <!-- Create Modal -->
    <BaseModal :show="showCreate" @close="showCreate = false">
      <template #title>
        <div class="flex items-center gap-3 text-amber-500">
          <img :src="PIXEL_ASSETS.ICON_SCROLL" class="w-6 h-6 pixelated" />
          <span class="uppercase font-black text-xl tracking-wide">Draft Contract</span>
        </div>
      </template>
      <form @submit.prevent="handleCreate" class="space-y-6 font-mono text-slate-200">
        <div class="space-y-4">
          <div class="p-4 bg-slate-900 border-2 border-slate-700">
            <label class="flex items-center justify-between cursor-pointer group">
              <div>
                <span class="font-bold text-amber-100 uppercase tracking-wide group-hover:text-amber-400">Private Event</span>
                <p class="text-[10px] text-slate-500 uppercase mt-1">Requires invitation code</p>
              </div>
              <input type="checkbox" v-model="options.isPrivate" class="accent-amber-600 w-5 h-5 cursor-pointer" />
            </label>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="bg-slate-900 p-3 border-2 border-slate-700">
              <label class="text-[10px] font-bold uppercase text-slate-500 block mb-2">Enemy Bots</label>
              <input type="number" v-model.number="options.botCount" min="0" max="10" 
                class="w-full bg-black border border-slate-700 p-2 text-center text-amber-400 font-bold focus:outline-none focus:border-amber-500 rounded" />
            </div>
            <div class="bg-slate-900 p-3 border-2 border-slate-700">
              <label class="text-[10px] font-bold uppercase text-slate-500 block mb-2">Difficulty Lvl</label>
              <input type="number" v-model.number="options.botPowerlevel" min="1" max="100" 
                class="w-full bg-black border border-slate-700 p-2 text-center text-amber-400 font-bold focus:outline-none focus:border-amber-500 rounded" />
            </div>
          </div>

          <div class="p-4 bg-slate-900 border-2 border-emerald-900/50">
            <label class="flex items-center justify-between cursor-pointer group">
              <div>
                <span class="font-bold text-emerald-100 uppercase tracking-wide group-hover:text-emerald-400">Auto-Start</span>
                <p class="text-[10px] text-slate-500 uppercase mt-1">Begin battle immediately when ready</p>
              </div>
              <input type="checkbox" v-model="options.autoStart" class="accent-emerald-600 w-5 h-5 cursor-pointer" />
            </label>
          </div>
        </div>

        <div class="flex justify-between gap-3 pt-4 border-t-2 border-slate-800 border-dashed">
          <button type="button" @click="showCreate = false" class="px-4 py-2 text-xs uppercase font-bold text-slate-500 hover:text-slate-300">Cancel</button>
          <button type="submit" :disabled="isCreating" class="rpg-btn bg-amber-600 border-amber-800 text-slate-900 hover:bg-amber-500 font-black uppercase tracking-wider px-8 py-2">
            Create
          </button>
        </div>
      </form>
    </BaseModal>

    <!-- Join Modal -->
    <BaseModal :show="showJoinModal" @close="showJoinModal = false">
      <template #title>
        <span class="uppercase font-black text-xl text-indigo-400">Select Hero</span>
      </template>
      <div class="space-y-4 font-mono text-slate-200">
        <div class="space-y-2 max-h-[400px] overflow-y-auto custom-scrollbar pr-2">
          <div
            v-for="f in roster.fighters"
            :key="f.id"
            class="relative flex items-center gap-4 p-3 border-2 cursor-pointer transition-all duration-200 group bg-slate-900"
            :class="selectedFighterId === f.id ? 'border-indigo-500 bg-indigo-900/20' : 'border-slate-800 hover:border-slate-600 hover:bg-slate-800'"
            @click="selectedFighterId = f.id"
          >
            <div class="w-12 h-12 bg-black border-2 border-slate-700 flex items-center justify-center">
              <img :src="PIXEL_ASSETS.ICON_FIGHTER" class="w-10 h-10 pixelated object-cover" />
            </div>
            <div class="flex-1">
              <h4 class="font-bold text-white uppercase tracking-wider" :class="selectedFighterId === f.id ? 'text-indigo-300' : ''">{{ f.name }}</h4>
              <div class="text-[10px] text-slate-500 uppercase mt-1">Lvl {{ f.level ?? 1 }} • {{ f.class ?? 'Warrior' }}</div>
            </div>
            <div v-if="selectedFighterId === f.id" class="text-indigo-400 font-bold text-lg">
              <<
            </div>
          </div>
        </div>

        <div class="flex justify-between gap-3 pt-6 border-t-2 border-slate-800 border-dashed">
          <button type="button" @click="showJoinModal = false" class="px-4 py-2 text-xs uppercase font-bold text-slate-500 hover:text-slate-300">Retreat</button>
          <button :disabled="!selectedFighterId || isJoining" @click="confirmJoin" class="rpg-btn bg-indigo-600 border-indigo-900 text-white hover:bg-indigo-500 font-black uppercase tracking-wider px-6 py-2">
            Ready
          </button>
        </div>
      </div>
    </BaseModal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, watch } from 'vue';
import { useRouter } from 'vue-router';
import BaseModal from '@/shared/ui/BaseModal.vue';
import Skeleton from '@/shared/ui/Skeleton.vue';
import EmptyState from '@/shared/ui/EmptyState.vue';
import { useAuthStore } from '@/features/auth/store';
import { useRosterStore } from '@/features/roster/store';
import { useRewardsStore } from '@/features/rewards/store';
import { PIXEL_ASSETS, getMatchStatusIcon, getMatchHeroImage } from '@/shared/utils/pixelAssets';
import type { Match, MatchStatus } from '@/features/matches/api';
import { 
  getMatches, 
  getCurrentMatch, 
  getOnlinePlayers, 
  quickJoinMatch, 
  joinMatch, 
  createMatch, 
  leaveMatch, 
  startMatch 
} from '@/features/matches/api';

const router = useRouter();
const auth = useAuthStore();
const roster = useRosterStore();
const rewardsStore = useRewardsStore();

interface MatchOption {
  isPrivate: boolean;
  botCount: number;
  botPowerlevel: number;
  maxPlayers: number;
  autoStart: boolean;
}

interface FilterOption {
  value: MatchStatus | 'all';
  label: string;
}

const filterOptions: FilterOption[] = [
  { value: 'all', label: 'All' },
  { value: 'lobby', label: 'Open' },
  { value: 'running', label: 'Live' },
  { value: 'completed', label: 'Done' },
];

const matchCounts = ref({ lobby: 0, running: 0, completed: 0 });

const matches = ref<Match[]>([]);
const isLoading = ref(false);
const search = ref('');
const browseStatus = ref<MatchStatus | 'all'>('all');

const showCreate = ref(false);
const showJoinModal = ref(false);
const joinTargetMatchId = ref<string | null>(null);
const selectedFighterId = ref<string | null>(null);
const isJoining = ref(false);
const isCreating = ref(false);

const currentMatchId = ref<string | null>(null);
const currentMatch = ref<Match | null>(null);
const myFighterIdInMatch = ref<string | null>(null);
const isStarting = ref(false);

const isQuickJoining = ref(false);
const onlinePlayers = ref(0);
const onlinePlayersPoll = ref<number | null>(null);

const statusMessage = ref('');
const statusTone = ref<'info' | 'success' | 'warning' | 'error'>('info');
const pollHandle = ref<number | null>(null);
const wsRef = ref<WebSocket | null>(null);
const hoveredMatchId = ref<string | null>(null);

const options = ref<MatchOption>({
  isPrivate: false,
  botCount: 1,
  botPowerlevel: 10,
  maxPlayers: 2,
  autoStart: true
});

const currentMatchStatus = computed(() => currentMatch.value?.status ?? (currentMatchId.value ? 'lobby' : null));

const currentMatchStatusLabel = computed(() => {
  const s = currentMatchStatus.value;
  if (s === 'lobby') return 'LOBBY ACTIVE';
  if (s === 'running') return 'BATTLE ENGAGED';
  if (s === 'completed') return 'QUEST COMPLETE';
  return 'ACTIVE';
});

const statusIconImg = computed(() => getMatchStatusIcon(currentMatchStatus.value || 'lobby'));

const statusTextClass = computed(() => {
  const tone = statusTone.value;
  if (tone === 'success') return 'text-emerald-400';
  if (tone === 'warning') return 'text-amber-400';
  if (tone === 'error') return 'text-red-400';
  return 'text-slate-400';
});

const filteredMatches = computed(() => {
  let list = matches.value;
  if (browseStatus.value !== 'all') {
    list = list.filter(m => m.status === browseStatus.value);
  }
  if (search.value.trim()) {
    const q = search.value.trim().toLowerCase();
    list = list.filter((m: Match) => (m.id || '').toLowerCase().includes(q));
  }
  return list;
});

const emptyStateTitle = computed(() => {
  if (browseStatus.value === 'lobby' || browseStatus.value === 'all') return 'NO OPEN CONTRACTS';
  if (browseStatus.value === 'running') return 'NO ACTIVE BATTLES';
  return 'ARCHIVES EMPTY';
});

const emptyStateMessage = computed(() => {
  if (browseStatus.value === 'lobby' || browseStatus.value === 'all') {
    return roster.fighters.length === 0 ? 'Recruit a hero first to join battles' : 'Create a new contract to start';
  }
  if (browseStatus.value === 'running') return 'No battles in progress';
  return 'Complete some matches to see history';
});

function getCountForStatus(status: string): number {
  if (status === 'all') return matches.value.length;
  return matchCounts.value[status as keyof typeof matchCounts.value] || 0;
}

function getMatchTitle(match: Match): string {
  if (match.status === 'lobby') return 'Open Arena';
  if (match.status === 'running') return 'Battle in Progress';
  return 'Match Complete';
}

function getMatchCardGradient(status: string): string {
  if (status === 'lobby') return 'bg-gradient-to-br from-slate-800 to-slate-900';
  if (status === 'running') return 'bg-gradient-to-br from-red-900/50 to-slate-900';
  return 'bg-gradient-to-br from-emerald-900/30 to-slate-900';
}

function getStatusBadgeClass(status: string): string {
  if (status === 'lobby') return 'bg-green-900/80 text-green-300 border-green-600';
  if (status === 'running') return 'bg-red-900/80 text-red-300 border-red-600 animate-pulse';
  return 'bg-slate-700 text-slate-300 border-slate-500';
}

function isMatchHot(match: Match): boolean {
  const regCount = match.registrations?.length || 0;
  const maxPlayers = match.options?.maxFightersPerUser || 2;
  return regCount >= maxPlayers - 1;
}

function setStatus(message: string, tone: 'info' | 'success' | 'warning' | 'error' = 'info') {
  statusMessage.value = message;
  statusTone.value = tone;
}

async function fetchMatches() {
  if (!auth.token) return;
  isLoading.value = true;
  try {
    const data = await getMatches(auth.token, 1, 100);
    matches.value = data?.items ?? [];
    
    matchCounts.value = {
      lobby: matches.value.filter(m => m.status === 'lobby').length,
      running: matches.value.filter(m => m.status === 'running').length,
      completed: matches.value.filter(m => m.status === 'completed').length,
    };
  } catch (e) {
    console.error(e);
    setStatus('Failed to load matches', 'error');
  } finally {
    isLoading.value = false;
  }
}

async function fetchCurrentMatch() {
  if (!auth.token || !currentMatchId.value) return;
  try {
    const m = await getCurrentMatch(auth.token);
    if (m) {
      currentMatch.value = m;
    } else if (currentMatchId.value) {
      clearCurrentMatch();
    }
  } catch (e) {
    console.error(e);
    if ((e as any)?.status === 404) clearCurrentMatch();
  }
}

function connectWebSocket() {
  if (!auth.token || !currentMatchId.value) return;
  
  // Robustly construct WebSocket URL by stripping trailing slashes and /api suffix
  let base = (import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:54321').replace(/\/$/, '');
  if (base.endsWith('/api')) {
    base = base.substring(0, base.length - 4);
  }
  base = base.replace(/^http/, 'ws');
  
  const url = `${base}/ws/match`;
  
  if (wsRef.value) disconnectWebSocket();

  try {
    const ws = new WebSocket(url);
    ws.onopen = () => {
      ws.send(JSON.stringify({ action: 'subscribe', matchId: currentMatchId.value }));
    };
    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.matchId === currentMatchId.value || data.type === 'matchStatus' || data.type === 'matchEnded' || data.type === 'lobbyUpdate') {
          fetchCurrentMatch();
          if (data.type === 'matchEnded') {
             setStatus('QUEST COMPLETE!', 'success');
             fetchMatches();
          }
        }
      } catch (_) {}
    };
    ws.onclose = () => { wsRef.value = null; };
    wsRef.value = ws;
  } catch (_) {}
}

function disconnectWebSocket() {
  if (wsRef.value) {
    try {
      wsRef.value.send(JSON.stringify({ action: 'unsubscribe' }));
      wsRef.value.close();
    } catch (_) {}
    wsRef.value = null;
  }
}

function startPolling() {
  if (pollHandle.value) return;
  pollHandle.value = window.setInterval(() => {
    fetchCurrentMatch();
  }, 5000);
}

function stopPolling() {
  if (pollHandle.value) {
    clearInterval(pollHandle.value);
    pollHandle.value = null;
  }
}

function clearCurrentMatch() {
  currentMatchId.value = null;
  currentMatch.value = null;
  myFighterIdInMatch.value = null;
  setStatus('', 'info');
  stopPolling();
  disconnectWebSocket();
}

function openJoinModal(matchId: string) {
  if (!roster.fighters.length) {
    setStatus('RECRUIT A HERO FIRST', 'warning');
    return;
  }
  joinTargetMatchId.value = matchId;
  selectedFighterId.value = roster.fighters[0]?.id ?? null;
  showJoinModal.value = true;
}

async function confirmJoin() {
  if (!auth.token || !joinTargetMatchId.value || !selectedFighterId.value) return;
  isJoining.value = true;
  try {
    await joinMatch(auth.token, joinTargetMatchId.value, selectedFighterId.value);
    currentMatchId.value = joinTargetMatchId.value;
    myFighterIdInMatch.value = selectedFighterId.value;
    currentMatch.value = null;
    setStatus('JOINED PARTY. PREPARE FOR BATTLE.', 'success');
    showJoinModal.value = false;
    joinTargetMatchId.value = null;
    await fetchCurrentMatch();
    await fetchMatches();
    connectWebSocket();
    startPolling();
  } catch (e: any) {
    setStatus(e?.message || 'FAILED TO JOIN.', 'error');
  } finally {
    isJoining.value = false;
  }
}

async function handleCreate() {
  if (!auth.token) return;
  isCreating.value = true;
  try {
    const match = await createMatch(auth.token, {
      isPrivate: options.value.isPrivate,
      botCount: options.value.botCount,
      botPowerlevel: options.value.botPowerlevel,
      autoStart: options.value.autoStart,
    });
    showCreate.value = false;
    currentMatchId.value = match.id;
    myFighterIdInMatch.value = roster.fighters[0]?.id ?? null;
    currentMatch.value = match;
    setStatus('CONTRACT SIGNED.', 'success');
    connectWebSocket();
    startPolling();
    await fetchMatches();
  } catch (e: any) {
    setStatus(e?.message || 'FAILED TO CREATE.', 'error');
  } finally {
    isCreating.value = false;
  }
}

async function handleLeave() {
  if (!auth.token || !currentMatchId.value) return;
  let fighterId = myFighterIdInMatch.value;
  if (!fighterId && currentMatch.value && roster.fighters.length) {
     const myFighterIds = roster.fighters.map(f => f.id);
     const reg = currentMatch.value.registrations?.find((r: any) => myFighterIds.includes(r.fighterId));
     if (reg) fighterId = reg.fighterId;
  }
  if (!fighterId) fighterId = roster.fighters[0]?.id;

  if (!fighterId) return;

  try {
    await leaveMatch(auth.token, currentMatchId.value, fighterId);
    clearCurrentMatch();
    setStatus('ABANDONED QUEST.', 'warning');
    await fetchMatches();
  } catch (e: any) {
    setStatus(e?.message || 'FAILED TO FLEE.', 'error');
  }
}

async function handleStart() {
  if (!auth.token || !currentMatchId.value) return;
  isStarting.value = true;
  setStatus('ENTERING ARENA...', 'info');
  try {
    await startMatch(auth.token, currentMatchId.value);
    setStatus('BATTLE STARTED!', 'success');
    await fetchCurrentMatch();
  } catch (e: any) {
    setStatus(e?.message || 'FAILED TO START.', 'error');
  } finally {
    isStarting.value = false;
  }
}

async function fetchOnlinePlayers() {
  if (!auth.token) return;
  try {
    const data = await getOnlinePlayers(auth.token);
    onlinePlayers.value = data?.onlinePlayers ?? 0;
  } catch (e) {
    console.error('Failed to fetch online players', e);
  }
}

function startOnlinePlayersPolling() {
  fetchOnlinePlayers();
  onlinePlayersPoll.value = window.setInterval(fetchOnlinePlayers, 30000);
}

function stopOnlinePlayersPolling() {
  if (onlinePlayersPoll.value) {
    clearInterval(onlinePlayersPoll.value);
    onlinePlayersPoll.value = null;
  }
}

async function handleQuickJoin() {
  if (!auth.token || roster.fighters.length === 0) return;
  isQuickJoining.value = true;
  setStatus('FINDING BATTLE...', 'info');

  try {
    const fighter = roster.fighters[0];
    const match = await quickJoinMatch(auth.token, fighter.id);

    currentMatchId.value = match.id;
    myFighterIdInMatch.value = fighter.id;
    currentMatch.value = match;
    setStatus('JOINED BATTLE! PREPARE FOR COMBAT.', 'success');
    connectWebSocket();
    startPolling();
    await fetchMatches();
  } catch (e: any) {
    setStatus(e?.message || 'NO OPEN LOBBIES AVAILABLE.', 'warning');
  } finally {
    isQuickJoining.value = false;
  }
}

watch(currentMatchId, (id) => {
  if (id) {
    fetchCurrentMatch();
    connectWebSocket();
    startPolling();
  } else {
    stopPolling();
    disconnectWebSocket();
    currentMatch.value = null;
  }
});

watch(browseStatus, () => fetchMatches());

onMounted(async () => {
  await roster.fetchFighters();
  await rewardsStore.fetchRewards();
  try {
    const active = await getCurrentMatch(auth.token ?? '').catch(() => null);
    if (active && active.id) {
        currentMatchId.value = active.id;
        currentMatch.value = active;
        setStatus('SESSION RESTORED.', 'info');
    }
  } catch (e) {
    console.error("Failed to restore session", e);
  }

  await fetchMatches();
  startOnlinePlayersPolling();
});

onUnmounted(() => {
  stopPolling();
  disconnectWebSocket();
  stopOnlinePlayersPolling();
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

.pixel-box-sm {
  border: 2px solid #334155;
  box-shadow: 2px 2px 0 #0f172a;
}

.status-badge {
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  padding: 2px 6px;
  border: 1px solid;
}

.animate-bounce-slow {
  animation: bounce 2s infinite;
}

@keyframes bounce {
  0%, 100% { transform: translateY(-5%); }
  50% { transform: translateY(5%); }
}

.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}

.custom-scrollbar::-webkit-scrollbar-track {
  background: #1e293b;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #475569;
  border-radius: 3px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.match-card-skeleton {
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}
</style>
