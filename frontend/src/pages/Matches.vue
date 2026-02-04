<template>
  <div class="retro-rpg min-h-screen p-4 md:p-8 font-mono text-amber-100" :style="{ backgroundImage: `url('${PIXEL_ASSETS.BG_DUNGEON}')` }">
    <!-- Overlay for atmosphere -->
    <div class="fixed inset-0 bg-slate-950/80 pointer-events-none z-0"></div>
    
    <div class="relative z-10 max-w-7xl mx-auto space-y-8">
      <!-- Title Section -->
      <header class="flex flex-col md:flex-row md:items-end justify-between gap-6 pb-6 border-b-4 border-amber-900/50">
        <div class="flex items-center gap-4">
           <div class="w-16 h-16 bg-slate-900 border-4 border-amber-600 flex items-center justify-center shadow-lg">
              <img :src="PIXEL_ASSETS.ICON_CHEST" class="w-10 h-10 pixelated" />
           </div>
           <div>
            <h1 class="text-4xl md:text-6xl font-black tracking-tight text-amber-500 uppercase text-shadow-retro leading-none">
              QUEST BOARD
            </h1>
            <p class="text-amber-200/60 mt-2 font-bold text-xs uppercase tracking-widest">
              <span class="text-red-500">>></span> Find your glory in the arena
            </p>
           </div>
        </div>
        
        <div v-if="!currentMatchId" class="flex gap-4">
          <button 
            @click="showCreate = true"
            class="rpg-btn bg-emerald-700 border-emerald-900 text-white hover:bg-emerald-600 group relative px-6 py-3 font-bold uppercase tracking-wider flex items-center gap-2"
          >
            <span class="absolute inset-0 border-2 border-white/20 pointer-events-none"></span>
            <img :src="PIXEL_ASSETS.ICON_SCROLL" class="w-5 h-5 pixelated" />
            <span>New Contract</span>
          </button>
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
        <div v-if="currentMatchId" class="rpg-panel border-4 border-amber-600 bg-slate-900 shadow-xl p-1 relative overflow-hidden">
          <!-- Decorative Corners -->
          <div class="absolute top-0 left-0 w-2 h-2 bg-amber-500"></div>
          <div class="absolute top-0 right-0 w-2 h-2 bg-amber-500"></div>
          <div class="absolute bottom-0 left-0 w-2 h-2 bg-amber-500"></div>
          <div class="absolute bottom-0 right-0 w-2 h-2 bg-amber-500"></div>

          <div class="bg-slate-900/90 p-4 md:p-6 flex flex-col md:flex-row items-center justify-between gap-6 border border-amber-900/30">
            <div class="flex items-center gap-6">
              <div class="relative">
                <div class="w-20 h-20 bg-slate-950 border-4 border-slate-700 flex items-center justify-center overflow-hidden">
                   <img :src="statusIconImg" class="w-16 h-16 pixelated object-contain animate-bounce-slow" />
                </div>
                <div class="absolute -bottom-3 -right-3 bg-slate-900 border-2 border-amber-600 text-amber-500 text-xs font-bold px-2 py-0.5">
                  #{{ currentMatchId.substring(0, 4) }}
                </div>
              </div>
              
              <div>
                <h3 class="font-black text-2xl text-amber-100 uppercase tracking-wide flex flex-col">
                  <span class="text-xs text-amber-500 font-bold mb-1">Current Objective</span>
                  {{ currentMatchStatusLabel }}
                </h3>
                <div class="mt-2 flex items-center gap-2 text-sm font-bold" :class="statusTextClass">
                  <span class="w-2 h-2 bg-current animate-pulse"></span>
                  {{ statusMessage || 'Awaiting input...' }}
                </div>
              </div>
            </div>

            <div class="flex flex-col gap-2 w-full md:w-auto min-w-[200px]">
               <!-- Lobby Actions -->
               <template v-if="currentMatchStatus === 'lobby'">
                  <div class="flex justify-between items-center bg-black/40 px-3 py-1 border border-slate-700 mb-2">
                    <span class="text-[10px] text-slate-400 uppercase">Party Size</span>
                    <span class="text-amber-400 font-bold">{{ currentMatch?.registrations?.length ?? 1 }}/{{ options.maxPlayers || 2 }}</span>
                  </div>
                  <div class="grid grid-cols-2 gap-2">
                    <button @click="handleLeave" class="rpg-btn-small bg-red-900/50 border-red-800 text-red-300 hover:bg-red-900 hover:text-white">
                      Flee
                    </button>
                    <button @click="handleStart" :disabled="isStarting" class="rpg-btn-small bg-amber-600 border-amber-800 text-slate-900 hover:bg-amber-500 font-black">
                      BEGIN
                    </button>
                  </div>
               </template>

               <!-- Running Actions -->
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

               <!-- Completed Actions -->
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
      <div class="space-y-6">
        <!-- Filters (Scroll Style) -->
        <div class="bg-amber-100/5 p-4 border-y-4 border-amber-900/50 flex flex-wrap items-center justify-between gap-4 backdrop-blur-sm">
          <div class="flex items-center gap-2">
            <template v-for="stat in ['lobby', 'running', 'completed']" :key="stat">
              <button
                @click="browseStatus = stat"
                class="px-4 py-2 text-xs font-bold uppercase tracking-wider transition-all border-2 relative overflow-hidden"
                :class="browseStatus === stat ? 'bg-amber-600 border-amber-400 text-slate-900 shadow-lg scale-105 z-10' : 'bg-slate-900 border-slate-700 text-slate-500 hover:text-amber-200 hover:border-amber-700'"
              >
                {{ stat }}
              </button>
            </template>
          </div>
          
          <div class="relative w-full md:w-64">
            <input
              v-model="search"
              placeholder="Search Archives..."
              class="w-full bg-slate-900 border-2 border-slate-700 p-2 pl-3 text-amber-100 placeholder-slate-600 focus:outline-none focus:border-amber-500 uppercase text-xs font-bold"
            />
            <div class="absolute right-2 top-2 pointer-events-none opacity-50">
               🔍
            </div>
          </div>
        </div>

        <!-- Grid -->
        <div v-if="isLoading" class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
           <div v-for="i in 4" :key="i" class="h-64 bg-slate-900/50 border-4 border-slate-800 animate-pulse flex items-center justify-center">
              <span class="text-slate-700 uppercase font-bold text-xs">Loading...</span>
           </div>
        </div>

        <div v-else-if="!filteredMatches.length" class="py-20 text-center border-4 border-dashed border-slate-800 bg-slate-900/20">
           <img :src="PIXEL_ASSETS.ICON_SCROLL" class="w-16 h-16 mx-auto opacity-20 pixelated grayscale mb-4" />
           <h3 class="text-slate-500 uppercase font-bold tracking-widest">No Quests Found</h3>
           <p class="text-slate-600 text-xs mt-2">{{ emptyMessage }}</p>
        </div>

        <div v-else class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
           <div 
            v-for="match in filteredMatches" 
            :key="match.id"
            class="group rpg-card bg-slate-900 border-4 border-slate-800 hover:border-amber-600/50 transition-all duration-200 flex flex-col relative overflow-hidden"
            :class="{ 'border-amber-500 ring-4 ring-amber-500/20 z-10': currentMatchId === match.id }"
           >
              <!-- Card Top Decoration -->
              <div class="h-1 bg-slate-800 group-hover:bg-amber-600/50 transition-colors"></div>
              
              <div class="p-4 flex flex-col h-full relative z-10">
                 <div class="flex justify-between items-start mb-4">
                    <div>
                       <span class="text-[10px] text-slate-500 uppercase font-bold tracking-wider block mb-1">Quest ID: {{ match.id.substring(0, 4) }}</span>
                       <h3 class="text-amber-100 font-bold uppercase leading-tight text-lg group-hover:text-amber-400 transition-colors">
                          {{ match.status === 'lobby' ? 'Open Arena' : match.status }}
                       </h3>
                    </div>
                    <div class="w-10 h-10 bg-slate-950 border-2 border-slate-700 flex items-center justify-center group-hover:border-amber-600 transition-colors">
                       <img :src="getMatchIconImg(match.status)" class="w-6 h-6 pixelated" />
                    </div>
                 </div>

                 <!-- Stats -->
                 <div class="space-y-2 mb-6 flex-1">
                    <div class="flex items-center justify-between text-xs bg-slate-950/50 p-2 border border-slate-800">
                       <span class="text-slate-500 uppercase">Fighters</span>
                       <span class="text-amber-200 font-bold">{{ match.registrations?.length || 0 }}</span>
                    </div>
                    <div class="flex items-center justify-between text-xs bg-slate-950/50 p-2 border border-slate-800">
                       <span class="text-slate-500 uppercase">Bots</span>
                       <span class="text-amber-200 font-bold">{{ match.options?.botCount || 0 }}</span>
                    </div>
                 </div>

                 <!-- Action -->
                 <div class="mt-auto space-y-2">
                    <template v-if="(match.status ?? 'lobby') === 'lobby'">
                       <button 
                        v-if="currentMatchId !== match.id"
                        @click="openJoinModal(match.id)"
                        class="w-full rpg-btn-small bg-indigo-900 border-indigo-700 text-indigo-200 hover:bg-indigo-800 hover:text-white font-bold uppercase"
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
                          <button class="w-full rpg-btn-small bg-slate-800 border-slate-600 text-slate-300 hover:bg-slate-700 font-bold uppercase">
                             {{ match.status === 'running' ? 'Watch' : 'Results' }}
                          </button>
                       </router-link>
                    </template>
                 </div>
              </div>

              <!-- BG Texture Overlay -->
              <div class="absolute inset-0 opacity-5 pointer-events-none" :style="{ backgroundImage: `url('${PIXEL_ASSETS.BG_DUNGEON}')`, backgroundSize: '64px' }"></div>
           </div>
        </div>
      </div>
    </div>

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
                <input type="checkbox" v-model="options.isPrivate" class="accent-amber-600 w-5 h-5 rounded-none cursor-pointer" />
             </label>
          </div>

          <div class="grid grid-cols-2 gap-4">
             <div class="bg-slate-900 p-3 border-2 border-slate-700">
                <label class="text-[10px] font-bold uppercase text-slate-500 block mb-2">Enemy Bots</label>
                <input type="number" v-model.number="options.botCount" min="0" max="10" 
                  class="w-full bg-black border border-slate-700 p-2 text-center text-amber-400 font-bold focus:outline-none focus:border-amber-500" />
             </div>
             <div class="bg-slate-900 p-3 border-2 border-slate-700">
                <label class="text-[10px] font-bold uppercase text-slate-500 block mb-2">Difficulty Lvl</label>
                <input type="number" v-model.number="options.botPowerlevel" min="1" max="100" 
                  class="w-full bg-black border border-slate-700 p-2 text-center text-amber-400 font-bold focus:outline-none focus:border-amber-500" />
             </div>
          </div>

          <div class="p-4 bg-slate-900 border-2 border-emerald-900/50">
             <label class="flex items-center justify-between cursor-pointer group">
                <div>
                   <span class="font-bold text-emerald-100 uppercase tracking-wide group-hover:text-emerald-400">Auto-Start</span>
                   <p class="text-[10px] text-slate-500 uppercase mt-1">Begin battle immediately when ready</p>
                </div>
                <input type="checkbox" v-model="options.autoStart" class="accent-emerald-600 w-5 h-5 rounded-none cursor-pointer" />
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
import BaseCard from '@/shared/ui/BaseCard.vue';
import BaseButton from '@/shared/ui/BaseButton.vue';
import BaseModal from '@/shared/ui/BaseModal.vue';
import { request } from '@/shared/api/http';
import { endpoints } from '@/shared/api/endpoints';
import { useAuthStore } from '@/features/auth/store';
import { useRosterStore } from '@/features/roster/store';

const PIXEL_ASSETS = {
  BG_DUNGEON: 'https://vibemedia.space/bg_dungeon_v2_99283.png?prompt=dark%20dungeon%20floor%20tile%20texture%20seamless&style=pixel_game_asset&key=NOGON',
  ICON_SWORDS: 'https://vibemedia.space/icon_swords_v2_11234.png?prompt=crossed%20steel%20swords%20pixel%20art&style=pixel_game_asset&key=NOGON',
  ICON_TROPHY: 'https://vibemedia.space/trophy_icon_4d5e6f_v1.png?prompt=golden%20trophy%20pixel%20art&style=pixel_game_asset&key=NOGON',
  ICON_SCROLL: 'https://vibemedia.space/icon_scroll_v2_66789.png?prompt=ancient%20magic%20scroll%20pixel%20art&style=pixel_game_asset&key=NOGON',
  ICON_CHEST: 'https://vibemedia.space/icon_chest_v2_55667.png?prompt=ancient%20wooden%20treasure%20chest%20with%20golden%20brass%20corners%20and%20mystical%20glow&style=pixel_game_asset&key=NOGON',
  ICON_POTION: 'https://vibemedia.space/icon_potion_v2_77890.png?prompt=red%20health%20potion%20pixel%20art&style=pixel_game_asset&key=NOGON',
  ICON_BOT: 'https://vibemedia.space/icon_skull_v2_55432.png?prompt=skull%20icon%20pixel%20art&style=pixel_game_asset&key=NOGON',
  ICON_FIGHTER: 'https://vibemedia.space/fighter_hooded_8p9q0r_v1.png?prompt=mystery%20hooded%20figure%20pixel%20art&style=pixel_game_asset&key=NOGON'
};

const auth = useAuthStore();
const roster = useRosterStore();

// State
const matches = ref<any[]>([]);
const isLoading = ref(false);
const search = ref('');
const browseStatus = ref('lobby');

// Modals
const showCreate = ref(false);
const showJoinModal = ref(false);
const joinTargetMatchId = ref<string | null>(null);
const selectedFighterId = ref<string | null>(null);
const isJoining = ref(false);
const isCreating = ref(false);

// Active Match
const currentMatchId = ref<string | null>(null);
const currentMatch = ref<any>(null);
const myFighterIdInMatch = ref<string | null>(null);
const isStarting = ref(false);

// Feedback
const statusMessage = ref('');
const statusTone = ref<'info' | 'success' | 'warning' | 'error'>('info');
const pollHandle = ref<number | null>(null);
const wsRef = ref<WebSocket | null>(null);

const options = ref({
  isPrivate: false,
  botCount: 1,
  botPowerlevel: 10,
  maxPlayers: 2,
  autoStart: true
});

// Computeds
const currentMatchStatus = computed(() => currentMatch.value?.status ?? (currentMatchId.value ? 'lobby' : null));

const currentMatchStatusLabel = computed(() => {
  const s = currentMatchStatus.value;
  if (s === 'lobby') return 'LOBBY ACTIVE';
  if (s === 'running') return 'BATTLE ENGAGED';
  if (s === 'completed') return 'QUEST COMPLETE';
  return 'ACTIVE';
});

const statusIconImg = computed(() => {
  const s = currentMatchStatus.value;
  if (s === 'running') return PIXEL_ASSETS.ICON_SWORDS;
  if (s === 'completed') return PIXEL_ASSETS.ICON_TROPHY;
  return PIXEL_ASSETS.ICON_SCROLL;
});

const statusTextClass = computed(() => {
  const tone = statusTone.value;
  if (tone === 'success') return 'text-emerald-400';
  if (tone === 'warning') return 'text-amber-400';
  if (tone === 'error') return 'text-red-400';
  return 'text-slate-400';
});

const filteredMatches = computed(() => {
  let list = matches.value;
  if (search.value.trim()) {
    const q = search.value.trim().toLowerCase();
    list = list.filter((m: any) => (m.id || '').toLowerCase().includes(q));
  }
  return list;
});

const emptyMessage = computed(() => {
  if (browseStatus.value === 'lobby') return 'NO OPEN CONTRACTS';
  if (browseStatus.value === 'running') return 'NO ACTIVE BATTLES';
  return 'ARCHIVES EMPTY';
});

function getMatchIconImg(status: string) {
  if (status === 'running') return PIXEL_ASSETS.ICON_SWORDS;
  if (status === 'completed') return PIXEL_ASSETS.ICON_TROPHY;
  return PIXEL_ASSETS.ICON_SCROLL;
}

function setStatus(message: string, tone: 'info' | 'success' | 'warning' | 'error' = 'info') {
  statusMessage.value = message;
  statusTone.value = tone;
}

// Actions
async function fetchMatches() {
  if (!auth.token) return;
  isLoading.value = true;
  try {
    const data = await request<any>(endpoints.match + '/browse', {
      method: 'POST',
      token: auth.token,
      body: { page: 1, pageSize: 50, status: browseStatus.value || undefined }
    });
    matches.value = data?.items ?? [];
  } catch (e) {
    console.error(e);
  } finally {
    isLoading.value = false;
  }
}

async function fetchCurrentMatch() {
  if (!auth.token) return;
  if (currentMatchId.value) {
    try {
        const m = await request<any>(`${endpoints.match}/${currentMatchId.value}`, { token: auth.token });
        currentMatch.value = m;
    } catch (e) {
        console.error(e);
        if ((e as any)?.status === 404) clearCurrentMatch();
    }
  }
}

function connectWebSocket() {
  if (!auth.token || !currentMatchId.value) return;
  const base = (import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:54321').replace(/^http/, 'ws');
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
  }, 4000);
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
    await request<void>(endpoints.match + '/join', {
      method: 'POST',
      token: auth.token,
      body: { matchId: joinTargetMatchId.value, fighterId: selectedFighterId.value }
    });
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
    const match = await request<any>(endpoints.match + '/create', {
      method: 'PUT',
      token: auth.token,
      body: options.value
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
    await request<void>(endpoints.match + '/leave', {
      method: 'POST',
      token: auth.token,
      body: { matchId: currentMatchId.value, fighterId }
    });
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
    await request<void>(`${endpoints.match}/${currentMatchId.value}/start`, {
      method: 'POST',
      token: auth.token
    });
    setStatus('BATTLE STARTED!', 'success');
    await fetchCurrentMatch();
  } catch (e: any) {
    setStatus(e?.message || 'FAILED TO START.', 'error');
  } finally {
    isStarting.value = false;
  }
}

// Watchers
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

// Lifecycle
onMounted(async () => {
  await roster.fetchFighters();
  
  try {
    const active = await request<any>(endpoints.match + '/current', { token: auth.token ?? undefined }).catch(() => null);
    if (active && active.id) {
        currentMatchId.value = active.id;
        currentMatch.value = active;
        setStatus('SESSION RESTORED.', 'info');
    }
  } catch (e) {
    console.error("Failed to restore session", e);
  }

  await fetchMatches();
});

onUnmounted(() => {
  stopPolling();
  disconnectWebSocket();
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

.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.animate-bounce-slow {
  animation: bounce 2s infinite;
}

@keyframes bounce {
  0%, 100% { transform: translateY(-5%); }
  50% { transform: translateY(5%); }
}
</style>
