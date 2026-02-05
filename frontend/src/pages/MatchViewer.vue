<template>
  <div class="pixel-theme h-dvh max-h-dvh overflow-hidden flex flex-col gap-3 sm:gap-4 font-mono text-slate-200" :style="{ backgroundImage: `url('${PIXEL_ASSETS.BG_DUNGEON}')` }">
    <!-- Overlay for atmosphere -->
    <div class="fixed inset-0 bg-slate-950/80 pointer-events-none z-0"></div>

    <!-- Header / Nav -->
    <header class="relative z-10 flex flex-col sm:flex-row sm:items-center justify-between shrink-0 p-3 sm:p-4 border-b-4 border-amber-900/50 bg-slate-900/80 gap-3 sm:gap-4">
      <div class="flex items-center gap-3 sm:gap-4 min-w-0">
        <BaseButton variant="ghost" size="sm" @click="$router.push('/matches')" class-name="rpg-btn-small">
          ← BACK
        </BaseButton>
        <div>
          <h1 class="text-lg sm:text-2xl font-black tracking-tight text-amber-500 flex flex-wrap items-center gap-2 sm:gap-3 uppercase text-shadow-retro">
            BATTLE LOG #{{ matchId?.substring(0, 8) }}
            <span v-if="matchStatus" class="px-2 py-0.5 text-xs font-bold uppercase tracking-wider border-2" 
              :class="statusBadgeClass">
              {{ matchStatus }}
            </span>
          </h1>
        </div>
      </div>
      
      <!-- Live Indicator -->
      <div v-if="matchStatus === 'running'" class="flex items-center gap-2 px-3 py-1 bg-black/60 border border-amber-500/50">
         <span class="relative flex h-2 w-2">
            <span class="animate-ping absolute inline-flex h-full w-full bg-amber-400 opacity-75"></span>
            <span class="relative inline-flex h-2 w-2 bg-amber-500"></span>
         </span>
         <span class="text-xs font-bold text-amber-400 uppercase tracking-widest blink">LIVE FEED</span>
      </div>
    </header>

    <!-- Main Content Area -->
    <div class="relative z-10 flex-1 grid grid-cols-1 lg:grid-cols-3 gap-4 sm:gap-6 min-h-0 px-3 sm:px-4 pb-3 sm:pb-4 overflow-y-auto lg:overflow-hidden">
      
      <!-- Canvas / Viewport (2 cols) -->
      <div class="lg:col-span-2 flex flex-col gap-3 sm:gap-4 min-h-0 min-w-0">
        
        <!-- Viewer Container -->
        <div class="relative flex-1 min-h-[45vh] sm:min-h-[55vh] border-4 border-slate-700 bg-[#020617] overflow-hidden shadow-2xl flex flex-col group pixel-box">
          
          <!-- Loading Overlay -->
          <div v-if="isLoading" class="absolute inset-0 z-10 flex items-center justify-center bg-slate-950/80 backdrop-blur-sm">
             <div class="flex flex-col items-center gap-3">
               <div class="h-8 w-8 animate-spin border-4 border-amber-500 border-t-transparent"></div>
               <p class="text-amber-200 text-sm font-bold uppercase tracking-wider">Loading battle data...</p>
             </div>
          </div>

          <!-- Empty Overlay -->
          <div v-else-if="!ticks.length" class="absolute inset-0 z-10 flex items-center justify-center">
             <p class="text-slate-500 font-bold uppercase tracking-wider">Waiting for battle data...</p>
          </div>

          <!-- Canvas -->
          <canvas
            ref="canvasRef"
            class="w-full h-full object-cover cursor-move pixelated"
            @mousedown="startDrag"
            @mousemove="onDrag"
            @mouseup="endDrag"
            @mouseleave="endDrag"
            @wheel.prevent="onWheel"
          ></canvas>

          <!-- Vignette & Atmosphere Overlays -->
          <div class="absolute inset-0 pointer-events-none bg-[radial-gradient(circle_at_center,transparent_0%,rgba(0,0,0,0.8)_100%)]"></div>
          <!-- CRT Scanline Effect -->
          <div class="absolute inset-0 pointer-events-none opacity-10 bg-[linear-gradient(transparent_50%,rgba(0,0,0,0.5)_50%)] bg-size-[100%_4px]"></div>

          <!-- Victory / End Overlay -->
          <transition name="fade">
            <div v-if="matchStatus === 'completed' && orderedRounds.length && selectedRound === orderedRounds[orderedRounds.length-1]" 
                class="absolute inset-0 z-20 flex items-center justify-center bg-black/80 backdrop-blur-sm p-4">
              <div class="max-w-md w-full text-center border-4 border-amber-500 bg-slate-900 p-6 shadow-2xl pixel-box animate-in zoom-in duration-300">
                <img v-if="isWinner" :src="PIXEL_ASSETS.ICON_TROPHY" class="w-16 h-16 mx-auto mb-2 pixelated animate-bounce-slow" />
                <img v-else :src="PIXEL_ASSETS.ICON_SKULL" class="w-16 h-16 mx-auto mb-2 pixelated opacity-50" />
                
                <h2 v-if="isWinner" class="text-3xl font-black text-amber-500 tracking-widest uppercase text-shadow-retro">VICTORY</h2>
                <h2 v-else class="text-3xl font-black text-rose-500 tracking-widest uppercase text-shadow-retro">DEFEAT</h2>
                
                <!-- Rewards Preview Section -->
                <div class="mt-4 p-4 bg-slate-950/50 border-2 border-slate-800">
                  <h3 class="text-xs font-bold text-slate-500 uppercase mb-3 text-shadow-retro">Rewards available to claim</h3>
                  
                  <div v-if="rewardsStore.rewardCount > 0" class="flex flex-col gap-2">
                    <div class="flex items-center justify-between text-sm">
                      <span class="text-slate-400">Match Bonus:</span>
                      <span class="text-amber-400 font-bold">+ {{ rewardsStore.rewardCount * 20 }} Particles</span>
                    </div>
                    <button 
                      class="rpg-btn-small mt-2 w-full !bg-amber-600 !text-slate-900 border-amber-800 hover:!bg-amber-500 font-bold"
                      @click="claimAndExit"
                      :disabled="rewardsStore.isLoading"
                    >
                      {{ rewardsStore.isLoading ? 'CLAIMING...' : 'CLAIM & EXIT' }}
                    </button>
                  </div>
                  <div v-else class="py-4">
                    <p class="text-slate-500 text-xs italic">All rewards claimed or none earned.</p>
                    <button class="rpg-btn-small mt-3 w-full font-bold" @click="$router.push('/matches')">EXIT TO MAP</button>
                  </div>
                </div>
              </div>
            </div>
          </transition>

          <!-- Playback Controls Overlay (Bottom) -->
          <div class="absolute bottom-0 inset-x-0 bg-slate-900/90 border-t-4 border-slate-700 p-3 sm:p-4 flex flex-col gap-3 transition-opacity duration-300 opacity-0 group-hover:opacity-100">
             
             <!-- Progress Bar -->
             <div class="relative h-4 bg-slate-950 border-2 border-slate-700 cursor-pointer group/bar" @click="seekToPercent">
                <div class="absolute inset-y-0 left-0 bg-amber-600 transition-all duration-100" :style="{ width: progressPercent + '%' }"></div>
                <div class="absolute inset-y-0 w-2 bg-white opacity-0 group-hover/bar:opacity-100 transition-opacity" :style="{ left: progressPercent + '%' }"></div>
             </div>

             <!-- Controls Row -->
             <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between mt-1 gap-3">
                <div class="flex items-center gap-3 sm:gap-4">
                   <button @click="togglePlayback" class="rpg-btn-small bg-amber-600 border-amber-800 text-slate-900 hover:bg-amber-500 font-black h-10 w-10 flex items-center justify-center">
                      <span v-if="isPlaying">||</span>
                      <span v-else class="ml-1">></span>
                   </button>
                   <div class="flex items-center gap-2 bg-slate-950 px-3 py-1.5 border-2 border-slate-700">
                      <button @click="stepRound(-1)" class="p-1 text-slate-400 hover:text-white transition-colors active:scale-90 font-bold"><<</button>
                      <span class="text-xs font-mono text-amber-400 w-20 text-center font-bold">RND {{ selectedRound }}</span>
                      <button @click="stepRound(1)" class="p-1 text-slate-400 hover:text-white transition-colors active:scale-90 font-bold">>></button>
                   </div>
                </div>

                <div class="flex flex-wrap items-center gap-3 sm:gap-4">
                   <!-- Speed Control -->
                   <div class="flex items-center gap-2 text-xs text-slate-400 bg-slate-950 px-3 py-1.5 border-2 border-slate-700">
                      <span class="uppercase font-bold">SPD</span>
                      <input type="range" min="0.5" max="4" step="0.5" v-model.number="playbackSpeed" class="w-20 accent-amber-500 h-2 bg-slate-800 appearance-none cursor-pointer" />
                      <span class="w-8 text-right text-amber-200 font-mono font-bold">{{ playbackSpeed }}x</span>
                   </div>
                   <!-- Zoom Control -->
                   <div class="flex items-center gap-2 text-xs text-slate-400 bg-slate-950 px-3 py-1.5 border-2 border-slate-700">
                      <span class="uppercase font-bold">ZM</span>
                      <button @click="zoom(-0.1)" class="px-1.5 hover:text-white font-bold text-lg leading-none">-</button>
                      <span class="w-10 text-center font-mono font-bold">{{ Math.round(camera.zoom * 100) }}%</span>
                      <button @click="zoom(0.1)" class="px-1.5 hover:text-white font-bold text-lg leading-none">+</button>
                   </div>
                </div>
             </div>
          </div>
        </div>

        <!-- Legend (Inline below canvas) -->
        <div class="flex flex-wrap items-center gap-3 sm:gap-6 px-4 sm:px-5 py-3 border-2 border-slate-700 bg-slate-900/80 text-xs text-slate-400 shadow-lg">
           <div class="flex items-center gap-2.5">
              <img :src="PIXEL_ASSETS.ICON_SWORD" class="w-4 h-4 pixelated" />
              <span class="font-bold text-emerald-200 uppercase">Heroes</span>
           </div>
           <div class="flex items-center gap-2.5">
              <img :src="PIXEL_ASSETS.ICON_SKULL" class="w-4 h-4 pixelated" />
              <span class="font-bold text-rose-200 uppercase">Enemies</span>
           </div>
           <div class="flex-1 hidden sm:block"></div>
           <div class="flex flex-wrap items-center gap-4 sm:gap-6 font-mono text-[11px] tracking-wide">
              <div class="flex flex-col items-center leading-none gap-1">
                 <span class="text-slate-500 font-bold uppercase">TOTAL</span>
                 <span class="text-slate-200 font-bold text-sm">{{ selectedMeta.total }}</span>
              </div>
              <div class="w-px h-6 bg-slate-700"></div>
              <div class="flex flex-col items-center leading-none gap-1">
                 <span class="text-emerald-500/70 font-bold uppercase">ALIVE</span>
                 <span class="text-emerald-400 font-bold text-sm">{{ selectedMeta.alive }}</span>
              </div>
              <div class="w-px h-6 bg-slate-700"></div>
              <div class="flex flex-col items-center leading-none gap-1">
                 <span class="text-rose-500/70 font-bold uppercase">FALLEN</span>
                 <span class="text-rose-400 font-bold text-sm">{{ selectedMeta.fallen }}</span>
              </div>
           </div>
        </div>

      </div>

      <!-- Logs / Sidebar (1 col) -->
      <div class="border-4 border-slate-700 bg-slate-950 flex flex-col overflow-hidden min-h-0 shadow-xl pixel-box">
        <div class="p-4 border-b-4 border-slate-800 bg-slate-900 flex items-center justify-between">
           <h3 class="font-bold text-amber-500 text-sm uppercase tracking-wider flex items-center gap-2">
              <span class="w-2 h-2 bg-amber-500 animate-pulse"></span>
              Combat Log
           </h3>
           <span class="text-[10px] font-mono text-slate-400 bg-slate-950 px-2 py-1 border border-slate-700 font-bold">{{ ticks.length }} TICKS</span>
        </div>
        
        <div class="flex-1 overflow-y-auto p-2 space-y-2 custom-scrollbar scroll-smooth bg-[#020617]" ref="logsContainer">
           <div 
              v-for="round in ticks" 
              :key="round.round"
              :id="'log-round-' + round.round"
              @click="selectRound(round.round)"
              class="border-2 p-3 text-xs transition-all cursor-pointer group relative overflow-hidden"
              :class="round.round === selectedRound ? 'border-amber-600 bg-amber-900/20' : 'border-slate-800 hover:border-slate-600 hover:bg-slate-900'"
           >
              <!-- Active Indicator Bar -->
              <div v-if="round.round === selectedRound" class="absolute left-0 top-0 bottom-0 w-1 bg-amber-500"></div>

              <div class="flex items-center justify-between mb-2 opacity-80 group-hover:opacity-100 pl-2">
                 <span class="font-bold text-amber-100 font-mono tracking-tight uppercase">ROUND {{ round.round.toString().padStart(2, '0') }}</span>
                 <span v-if="round.ticks?.length" class="text-[10px] text-slate-500 font-bold">{{ round.ticks.length }} EVENTS</span>
              </div>
              
              <div class="space-y-1.5 pl-2">
                 <div v-for="(tick, idx) in round.ticks" :key="idx" class="pl-3 border-l-2 border-slate-800 py-0.5 relative">
                    <!-- Timeline Dot -->
                    <div class="absolute -left-[5px] top-2 w-[8px] h-[8px] bg-slate-800 border border-slate-600 rotate-45"></div>
                    
                    <template v-if="tick.type === 'attack'">
                       <div class="flex items-baseline gap-1.5 flex-wrap leading-relaxed">
                          <span class="text-emerald-400 font-bold">{{ formatFighterId(payloadValue(tick.payload, 'attackerId')) }}</span>
                          <span class="text-slate-600 text-[10px] uppercase font-bold tracking-wider">HIT</span>
                          <span class="text-rose-400 font-bold">{{ formatFighterId(payloadValue(tick.payload, 'targetId')) }}</span>
                          <span class="font-mono font-bold ml-1" :class="payloadValue(tick.payload, 'isCritical') ? 'text-amber-400 text-sm' : 'text-slate-300'">
                             -{{ payloadValue(tick.payload, 'damage') }}
                          </span>
                          <span v-if="payloadValue(tick.payload, 'isCritical')" class="text-[9px] text-slate-900 bg-amber-500 font-bold uppercase px-1">CRIT</span>
                       </div>
                    </template>
                    <template v-else-if="tick.type === 'died'">
                       <div class="flex items-center gap-2 p-1 bg-rose-950/30 border border-rose-900/50 mt-1">
                          <span class="text-rose-500 text-lg">☠</span>
                          <div class="flex flex-col">
                             <span class="text-rose-500 font-bold uppercase text-[10px] tracking-wide">SLAIN</span>
                             <span class="text-rose-300 font-bold">{{ formatFighterId(payloadValue(tick.payload, 'fighterId')) }}</span>
                          </div>
                       </div>
                    </template>
                    <template v-else-if="tick.type === 'spawn'">
                       <div class="text-emerald-500/60 italic text-[11px] font-bold">
                          ✨ {{ formatFighterId(payloadValue(tick.payload, 'fighterId')) }} APPEARED
                       </div>
                    </template>
                    <template v-else-if="tick.type === 'momentum'">
                       <div class="flex items-center gap-1.5 text-[11px]">
                          <span class="text-blue-400 font-bold">{{ formatFighterId(payloadValue(tick.payload, 'fighterId')) }}</span>
                          <span class="text-slate-500 uppercase">MOMENTUM</span>
                          <span class="font-mono font-bold" :class="payloadValue(tick.payload, 'momentum') > 50 ? 'text-amber-400' : 'text-blue-400'">
                             {{ payloadValue(tick.payload, 'momentum') }}%
                          </span>
                          <span v-if="payloadValue(tick.payload, 'consecutiveHits') > 1" class="text-amber-500 font-bold">
                             [{{ payloadValue(tick.payload, 'consecutiveHits') }}x]
                          </span>
                       </div>
                    </template>
                    <template v-else-if="tick.type === 'momentum_decay'">
                       <div class="flex items-center gap-1.5 text-[10px] opacity-70">
                          <span class="text-slate-400 font-bold">{{ formatFighterId(payloadValue(tick.payload, 'fighterId')) }}</span>
                          <span class="text-slate-600 uppercase">DECAY</span>
                          <span class="font-mono text-slate-500">{{ payloadValue(tick.payload, 'momentum') }}%</span>
                       </div>
                    </template>
                    <template v-else-if="tick.type === 'sunder'">
                       <div class="flex items-center gap-1.5 text-[11px]">
                          <span class="text-rose-400 font-bold">SUNDER</span>
                          <span class="text-slate-500">→</span>
                          <span class="text-rose-300 font-bold">{{ formatFighterId(payloadValue(tick.payload, 'targetId')) }}</span>
                          <span class="text-[9px] bg-rose-900/50 text-rose-400 px-1 border border-rose-700">
                             -{{ payloadValue(tick.payload, 'armorReduced') }} ARMOR
                          </span>
                       </div>
                    </template>
                    <template v-else-if="tick.type === 'flurry'">
                       <div class="flex items-center gap-1.5 p-1 bg-amber-950/30 border border-amber-900/50">
                          <span class="text-amber-500">⚡</span>
                          <span class="text-amber-400 font-bold">{{ formatFighterId(payloadValue(tick.payload, 'fighterId')) }}</span>
                          <span class="text-amber-300 uppercase text-[10px] tracking-wider font-bold">FLURRY!</span>
                          <span class="text-[9px] text-amber-500">+{{ payloadValue(tick.payload, 'attackSpeedBonus') }}% SPD</span>
                       </div>
                    </template>
                 </div>
              </div>
           </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick, computed, onUnmounted, reactive } from 'vue';
import { useRoute } from 'vue-router';
import { request } from '@/shared/api/http';
import { endpoints } from '@/shared/api/endpoints';
import { useAuthStore } from '@/features/auth/store';
import BaseButton from '@/shared/ui/BaseButton.vue';
import { useRosterStore } from '@/features/roster/store';
import { useRewardsStore } from '@/features/rewards/store';

// Pixel Assets Definition
const PIXEL_ASSETS = {
  BG_DUNGEON: 'https://vibemedia.space/bg_dungeon_v2_99283.png?prompt=dark%20dungeon%20floor%20tile%20texture%20seamless&style=pixel_game_asset&key=NOGON',
  CHAR_HERO: 'https://vibemedia.space/hero_knight_v2_88441.png?prompt=fantasy%20knight%20character%20sprite%20standing%20pixel%20art%20front%20view&style=pixel_game_asset&key=NOGON',
  CHAR_ENEMY: 'https://vibemedia.space/enemy_slime_v2_99223.png?prompt=cute%20slime%20monster%20sprite%20pixel%20art%20front%20view&style=pixel_game_asset&key=NOGON',
  ICON_SWORD: 'https://vibemedia.space/icon_sword_px_11223.png?prompt=pixel%20art%20sword%20icon&style=pixel_game_asset&key=NOGON',
  ICON_SKULL: 'https://vibemedia.space/icon_skull_px_99887.png?prompt=pixel%20art%20skull%20icon&style=pixel_game_asset&key=NOGON',
  ICON_TROPHY: 'https://vibemedia.space/trophy_icon_4d5e6f_v1.png?prompt=golden%20trophy%20pixel%20art&style=pixel_game_asset&key=NOGON',
  // Battlefield Textures (FX are now procedurally drawn)
  TEX_FLOOR: 'https://vibemedia.space/floor_cobble_7a8b9c_v1.png?prompt=cobblestone%20texture%20with%20worn%20stones%20and%20dirt%20between%20cracks&style=pixel_game_asset&key=NOGON',
  TEX_WALL: 'https://vibemedia.space/wall_stone_1x2y3z_v1.png?prompt=stone%20brick%20texture%20with%20weathered%20gray%20surface%20and%20moss%20details&style=pixel_game_asset&key=NOGON'
};

const auth = useAuthStore();
const roster = useRosterStore();
const rewardsStore = useRewardsStore();
const route = useRoute();
const matchId = ref(route.params.id as string);
const matchStatus = ref<string | null>(null);
const ticks = ref<any[]>([]);
const isLoading = ref(false);
let livePollHandle: number | null = null;
const canvasRef = ref<HTMLCanvasElement | null>(null);
const roundStateMap = ref<Record<number, RoundState>>({});

// Asset Preloading
const images: Record<string, HTMLImageElement> = {};
const assetsLoaded = ref(false);

const loadAssets = () => {
   let loadedCount = 0;
   const total = Object.keys(PIXEL_ASSETS).length;
   
   Object.entries(PIXEL_ASSETS).forEach(([key, url]) => {
      const img = new Image();
      img.src = url;
      img.onload = () => {
         loadedCount++;
         if (loadedCount === total) assetsLoaded.value = true;
      };
      images[key] = img;
   });
};

// Camera & View Settings
const TILE_WIDTH = 64;
const TILE_HEIGHT = 32;
const worldSize = 24; 
const camera = reactive({
  x: 0,
  y: 0,
  zoom: 1.0,
  shakeX: 0,
  shakeY: 0,
  isDragging: false,
  lastX: 0,
  lastY: 0
});

// Particles System
type Particle = {
   x: number;
   y: number;
   z: number; // Height
   vx: number;
   vy: number;
   vz: number;
   life: number;
   maxLife: number;
   color: string;
   size: number;
};
const particles = ref<Particle[]>([]);

const selectedRound = ref(0);
const isPlaying = ref(false);
const playbackSpeed = ref(1.0);
let playbackHandle: number | null = null;
let segmentStart = 0;
// We need a separate loop for purely visual updates (particles, shake decay) even if paused
let visualLoopHandle: number | null = null; 
const logsContainer = ref<HTMLElement | null>(null);

// Types
type EntityState = {
  id: string;
  x: number;
  y: number;
  hp: number;
  maxHp: number;
  alive: boolean;
  isPlayer: boolean;
  teamId?: string;
  // Visual state
  floatOffset: number; // For idle bobbing
  vx?: number;
  vy?: number;
  isMoving?: boolean;
  lastAngle?: number;
  // Combo-Momentum state
  momentum?: number;
  consecutiveHits?: number;
  sunderStacks?: number;
  flurryActive?: boolean;
};

type AttackEvent = {
  attackerId: string;
  targetId: string;
  isCritical?: boolean;
  isParried?: boolean;
  damage?: number;
};

type MomentumEvent = {
  fighterId: string;
  momentum: number;
  consecutiveHits: number;
  targetId?: string;
};

type SunderEvent = {
  targetId: string;
  stacks: number;
  armorReduced: number;
};

type FlurryEvent = {
  fighterId: string;
  attackSpeedBonus: number;
};

type RoundState = {
  round: number;
  entities: Record<string, EntityState>;
  events: {
    attacks: AttackEvent[];
    deaths: string[];
    momentum: MomentumEvent[];
    sunder: SunderEvent[];
    flurry: FlurryEvent[];
  };
};

// Utils
const payloadValue = (payload: unknown, key: string) => {
  if (!payload || typeof payload !== 'object') return undefined;
  return (payload as Record<string, any>)[key];
};

const formatFighterId = (id: any) => {
  if (typeof id !== 'string') return '--';
  return id.substring(0, 6);
};

const statusBadgeClass = computed(() => {
  if (matchStatus.value === 'running') return 'bg-amber-950 text-amber-500 border-amber-700';
  if (matchStatus.value === 'completed') return 'bg-slate-950 text-slate-500 border-slate-700';
  return 'bg-emerald-950 text-emerald-500 border-emerald-700';
});

const orderedRounds = computed(() => ticks.value.map((round) => round.round));

const progressPercent = computed(() => {
  if (!orderedRounds.value.length) return 0;
  const idx = orderedRounds.value.indexOf(selectedRound.value);
  return ((idx + 1) / orderedRounds.value.length) * 100;
});

const selectedMeta = computed(() => {
  const state = roundStateMap.value[selectedRound.value];
  if (!state) return { total: 0, alive: 0, fallen: 0 };
  const entities = Object.values(state.entities);
  return {
    total: entities.length,
    alive: entities.filter(e => e.alive).length,
    fallen: entities.filter(e => !e.alive).length
  };
});

// Isometric Math
const toIso = (x: number, y: number) => {
  const isoX = (x - y) * TILE_WIDTH * 0.5;
  const isoY = (x + y) * TILE_HEIGHT * 0.5;
  return { x: isoX, y: isoY };
};

const toScreen = (x: number, y: number, canvasW: number, canvasH: number) => {
  const iso = toIso(x, y);
  const centerX = canvasW / 2;
  const centerY = canvasH / 4; 
  
  return {
    x: centerX + (iso.x * camera.zoom) + camera.x + camera.shakeX,
    y: centerY + (iso.y * camera.zoom) + camera.y + camera.shakeY
  };
};

// Camera Shake
const addShake = (intensity: number) => {
   camera.shakeX = (Math.random() - 0.5) * intensity;
   camera.shakeY = (Math.random() - 0.5) * intensity;
};

// Canvas Interaction
const startDrag = (e: MouseEvent) => {
  camera.isDragging = true;
  camera.lastX = e.clientX;
  camera.lastY = e.clientY;
};

const onDrag = (e: MouseEvent) => {
  if (!camera.isDragging) return;
  const dx = e.clientX - camera.lastX;
  const dy = e.clientY - camera.lastY;
  camera.x += dx;
  camera.y += dy;
  camera.lastX = e.clientX;
  camera.lastY = e.clientY;
  // No need to call render here, visual loop handles it
};

const endDrag = () => { camera.isDragging = false; };
const onWheel = (e: WheelEvent) => {
  const delta = -Math.sign(e.deltaY) * 0.1;
  zoom(delta);
};
const zoom = (delta: number) => {
  camera.zoom = Math.max(0.2, Math.min(3.0, camera.zoom + delta));
};

// --- DRAWING LOGIC ---

// Tile Colors for "Dungeon Grid"
const gridColors = {
   top: '#1e293b', // slate-800
   side: '#0f172a', // slate-900
   stroke: '#334155', // slate-700
   highlight: '#475569' // slate-600
};

const drawTile = (ctx: CanvasRenderingContext2D, x: number, y: number, canvasW: number, canvasH: number, time: number) => {
  const pos = toScreen(x, y, canvasW, canvasH);
  const zoom = camera.zoom;
  const w = TILE_WIDTH * zoom;
  const h = TILE_HEIGHT * zoom;
  
  // Height variation for visual interest (noise based on coords)
  const noise = (Math.sin(x * 0.5) + Math.cos(y * 0.5)) * 2 * zoom;
  const tileY = pos.y - noise;

  // Draw Block (Sides first for depth)
  const depth = 8 * zoom;
  
  // Right Face
  ctx.beginPath();
  ctx.moveTo(pos.x, tileY + h);
  ctx.lineTo(pos.x + w/2, tileY + h/2);
  ctx.lineTo(pos.x + w/2, tileY + h/2 + depth);
  ctx.lineTo(pos.x, tileY + h + depth);
  ctx.fillStyle = gridColors.side;
  ctx.fill();
  
  // Left Face
  ctx.beginPath();
  ctx.moveTo(pos.x, tileY + h);
  ctx.lineTo(pos.x - w/2, tileY + h/2);
  ctx.lineTo(pos.x - w/2, tileY + h/2 + depth);
  ctx.lineTo(pos.x, tileY + h + depth);
  ctx.fillStyle = '#020617'; // Darker
  ctx.fill();

  // Top Face
  ctx.beginPath();
  ctx.moveTo(pos.x, tileY);
  ctx.lineTo(pos.x + w/2, tileY + h/2);
  ctx.lineTo(pos.x, tileY + h);
  ctx.lineTo(pos.x - w/2, tileY + h/2);
  ctx.closePath();

  ctx.save();
  ctx.clip();
  
  // Texture Mapping
  if (assetsLoaded.value && images['TEX_FLOOR']) {
     // Draw texture pattern (simple stretch for now to fit tile bounds)
     ctx.drawImage(images['TEX_FLOOR'], pos.x - w/2, tileY, w, h);
     
     // Overlay color to blend with theme
     ctx.fillStyle = 'rgba(30, 41, 59, 0.4)'; 
     ctx.fill();
  } else {
     ctx.fillStyle = gridColors.top;
     ctx.fill();
  }

  // Random "cracked tile" highlight
  if ((x + y * 7) % 5 === 0) {
     ctx.fillStyle = 'rgba(255, 255, 255, 0.1)';
     ctx.fill();
  }
  
  ctx.restore();
  
  ctx.strokeStyle = gridColors.stroke;
  ctx.lineWidth = 1 * zoom;
  ctx.stroke();
};

// --- 3D CHARACTER HELPERS ---

const rotatePt = (x: number, y: number, angle: number) => {
   const cos = Math.cos(angle);
   const sin = Math.sin(angle);
   return {
      x: x * cos - y * sin,
      y: x * sin + y * cos
   };
};

const shadeColor = (color: string, percent: number) => {
    const f = parseInt(color.slice(1), 16),
          t = percent < 0 ? 0 : 255,
          p = percent < 0 ? percent * -1 : percent,
          R = f >> 16,
          G = f >> 8 & 0x00FF,
          B = f & 0x0000FF;
    return "#" + (0x1000000 + (Math.round((t - R) * p) + R) * 0x10000 + (Math.round((t - G) * p) + G) * 0x100 + (Math.round((t - B) * p) + B)).toString(16).slice(1);
};

const drawIsoCuboid = (
  ctx: CanvasRenderingContext2D,
  cx: number, cy: number, cz: number, 
  w: number, d: number, h: number,   
  angle: number,
  color: string,
  zoom: number,
  canvasW: number, canvasH: number
) => {
  const hw = w / 2;
  const hd = d / 2;
  
  const rawCorners = [
     { x: hw, y: hd },   
     { x: -hw, y: hd },  
     { x: -hw, y: -hd }, 
     { x: hw, y: -hd }   
  ];
  
  const topWorld = rawCorners.map(p => {
     const rot = rotatePt(p.x, p.y, angle);
     return { x: cx + rot.x, y: cy + rot.y, z: cz + h };
  });
  
  const botWorld = rawCorners.map(p => {
     const rot = rotatePt(p.x, p.y, angle);
     return { x: cx + rot.x, y: cy + rot.y, z: cz };
  });

  const project = (p: {x:number, y:number, z:number}) => {
      const s = toScreen(p.x, p.y, canvasW, canvasH);
      return { x: s.x, y: s.y - p.z * TILE_HEIGHT * zoom };
  };

  const topScreen = topWorld.map(project);
  const botScreen = botWorld.map(project);
  
  ctx.fillStyle = shadeColor(color, 0.1); 
  ctx.beginPath();
  ctx.moveTo(topScreen[0].x, topScreen[0].y);
  for (let i = 1; i < 4; i++) ctx.lineTo(topScreen[i].x, topScreen[i].y);
  ctx.closePath();
  ctx.fill();
  ctx.strokeStyle = shadeColor(color, -0.2);
  ctx.lineWidth = 1;
  ctx.stroke();

  let maxYi = 0;
  let maxY = topScreen[0].y;
  for (let i = 1; i < 4; i++) {
     if (topScreen[i].y > maxY) {
        maxY = topScreen[i].y;
        maxYi = i;
     }
  }
  
  const drawSide = (idx1: number, idx2: number, shade: number) => {
     ctx.fillStyle = shadeColor(color, shade);
     ctx.beginPath();
     ctx.moveTo(topScreen[idx1].x, topScreen[idx1].y);
     ctx.lineTo(topScreen[idx2].x, topScreen[idx2].y);
     ctx.lineTo(botScreen[idx2].x, botScreen[idx2].y);
     ctx.lineTo(botScreen[idx1].x, botScreen[idx1].y);
     ctx.closePath();
     ctx.fill();
     ctx.stroke();
  };
  
  drawSide(maxYi, (maxYi + 1) % 4, -0.3);
  drawSide(maxYi, (maxYi + 3) % 4, -0.1);
};

const drawVoxelCharacter = (
  ctx: CanvasRenderingContext2D,
  x: number, y: number, floatZ: number,
  isPlayer: boolean,
  isMoving: boolean,
  isAttacking: boolean,
  angle: number,
  zoom: number,
  time: number,
  canvasW: number, canvasH: number,
  isHit?: boolean
) => {
   const baseColors = isPlayer ? {
      head: '#fbbf24', body: '#3b82f6', leg: '#1e3a8a', arm: '#93c5fd'
   } : {
      head: '#4ade80', body: '#166534', leg: '#14532d', arm: '#86efac'
   };
   
   if (isHit) {
      baseColors.head = '#ffffff'; baseColors.body = '#ffffff';
      baseColors.leg = '#ffffff'; baseColors.arm = '#ffffff';
   }
   
   const speed = 0.01;
   const limbCycle = isMoving ? Math.sin(time * speed) : 0;
   const legAngle = limbCycle * 0.5;
   
   let lArmAngle = -limbCycle * 0.5;
   let rArmAngle = limbCycle * 0.5;
   
   if (isAttacking) {
      // Attack Swing
      const attackCycle = Math.sin(time * 0.03); 
      rArmAngle = -Math.PI / 2 + attackCycle * 0.5; // Raised/Swing
   }
   
   const pz = 0 + floatZ * 0.01;
   const parts = [];
   
   const headS = 0.25;
   const bodyW = 0.25, bodyD = 0.12, bodyH = 0.3;
   const limbW = 0.08, limbD = 0.08, limbH = 0.3;
   
   const bodyZ = pz + limbH; 
   
   // Body
   parts.push({
      cx: 0, cy: 0, cz: bodyZ,
      w: bodyW, d: bodyD, h: bodyH,
      offsetVec: {x:0,y:0,z:0},
      color: baseColors.body
   });
   
   // Head
   parts.push({
      cx: 0, cy: 0, cz: bodyZ + bodyH,
      w: headS, d: headS, h: headS,
      offsetVec: {x:0, y:0, z: Math.sin(time * 0.002) * 0.02},
      color: baseColors.head
   });
   
   // Legs
   const legOffset = 0.06;
   parts.push({
      cx: 0, cy: legOffset, cz: pz,
      w: limbW, d: limbD, h: limbH,
      offsetVec: { x: Math.sin(legAngle)*0.1, y: 0, z: 0 },
      color: baseColors.leg
   });
   parts.push({
      cx: 0, cy: -legOffset, cz: pz,
      w: limbW, d: limbD, h: limbH,
      offsetVec: { x: Math.sin(-legAngle)*0.1, y: 0, z: 0 },
      color: baseColors.leg
   });
   
   // Arms
   const armZ = bodyZ + bodyH - 0.05;
   const armOffset = 0.16;
   
   // Left Arm
   parts.push({
      cx: 0, cy: armOffset, cz: armZ - limbH,
      w: limbW, d: limbD, h: limbH,
      offsetVec: { x: Math.sin(lArmAngle)*0.1, y: 0, z: Math.cos(lArmAngle)*0.1 - 0.1 },
      color: baseColors.arm
   });
   
   // Right Arm (Main Hand)
   parts.push({
      cx: 0, cy: -armOffset, cz: armZ - limbH,
      w: limbW, d: limbD, h: limbH,
      offsetVec: { x: Math.sin(rArmAngle)*0.2, y: 0, z: Math.cos(rArmAngle)*0.1 - 0.1 },
      color: baseColors.arm
   });

   // Weapon (Sword) if attacking
   if (isAttacking) {
       const swordLen = 0.4;
       // Attached to Right Arm
       const handX = Math.sin(rArmAngle)*0.2;
       const handZ = armZ - limbH + Math.cos(rArmAngle)*0.1 - 0.1;
       
       parts.push({
          cx: 0.1, cy: -armOffset - 0.05, cz: handZ, 
          w: swordLen, d: 0.05, h: 0.05,
          offsetVec: { x: handX + 0.2, y: 0, z: 0 },
          color: '#cbd5e1' // Blade
       });
   }

   const transformed = parts.map(p => {
      let lx = p.cx + p.offsetVec.x;
      let ly = p.cy + p.offsetVec.y;
      let lz = p.cz + p.offsetVec.z;
      
      const r = rotatePt(lx, ly, angle);
      
      const wx = x + r.x;
      const wy = y + r.y;
      
      const depth = wx + wy + lz * 0.1; 
      
      return { ...p, wx, wy, wz: lz, depth };
   });
   
   transformed.sort((a, b) => a.depth - b.depth);
   
   transformed.forEach(p => {
      drawIsoCuboid(ctx, p.wx, p.wy, p.wz, p.w, p.d, p.h, angle, p.color, zoom, canvasW, canvasH);
   });
};

const drawEntity = (
  ctx: CanvasRenderingContext2D, 
  entity: EntityState, 
  x: number, y: number, 
  canvasW: number, canvasH: number,
  time: number,
  visuals?: { hit?: boolean; attacking?: boolean; momentumEvent?: MomentumEvent; flurryActive?: boolean }
) => {
  const pos = toScreen(x, y, canvasW, canvasH);
  const zoom = camera.zoom;
  
  if (!entity.alive) {
     // Grave Marker / Corpse
     ctx.fillStyle = '#334155';
     const s = 10 * zoom;
     // Draw a pixelated skull or tombstone shape
     ctx.beginPath();
     ctx.rect(pos.x - s/2, pos.y, s, s);
     ctx.fill();
     return;
  }

  // Bobbing Animation
  const bob = Math.sin(time / 300 + entity.floatOffset) * 2 * zoom;
  const shakeX = visuals?.hit ? (Math.random() - 0.5) * 5 * zoom : 0;
  
  const charY = pos.y - 12 * zoom + bob;
  const centerX = pos.x + shakeX;

  // Shadow
  ctx.fillStyle = 'rgba(0,0,0,0.5)';
  ctx.beginPath();
  ctx.ellipse(pos.x, pos.y + 14 * zoom, 10 * zoom, 5 * zoom, 0, 0, Math.PI * 2);
  ctx.fill();

  // --- 3D VOXEL CHARACTER ---
  const isMoving = entity.isMoving || false;
  // Calculate angle from velocity or use last known or default
  let angle = 0;
  if (entity.vx != null || entity.vy != null) {
     angle = Math.atan2(entity.vy ?? 0, entity.vx ?? 0);
     // Adjust for iso view? atan2(y,x) gives angle in grid coords.
     // x=1,y=0 => 0 deg (SE). x=0,y=1 => 90 deg (SW).
     // This matches our rotation logic naturally.
  }
  
  drawVoxelCharacter(ctx, x, y, bob, entity.isPlayer, isMoving, visuals?.attacking || false, angle, zoom, time, canvasW, canvasH, visuals?.hit);

  // Health Bar (Retro Style)
  const hpPercent = Math.max(0, entity.hp / entity.maxHp);
  const barW = 24 * zoom;
  const barH = 4 * zoom;
  const barX = centerX - barW / 2;
  const barY = charY - 50 * zoom;

  // Border
  ctx.fillStyle = '#000';
  ctx.fillRect(barX - 1, barY - 1, barW + 2, barH + 2);
  
  // Fill
  ctx.fillStyle = hpPercent > 0.5 ? '#10b981' : (hpPercent > 0.25 ? '#f59e0b' : '#ef4444');
  ctx.fillRect(barX, barY, barW * hpPercent, barH);

  // --- MOMENTUM BAR ---
  const momentum = entity.momentum || 0;
  const maxMomentum = 100;
  const momentumPercent = momentum / maxMomentum;
  const momBarW = 24 * zoom;
  const momBarH = 3 * zoom;
  const momBarX = centerX - momBarW / 2;
  const momBarY = barY - momBarH - 3 * zoom;

  // Momentum bar background
  ctx.fillStyle = '#1e293b';
  ctx.fillRect(momBarX, momBarY, momBarW, momBarH);
  
  // Momentum bar fill (gradient from blue to purple to gold)
  let momentumColor = '#3b82f6'; // Blue
  if (momentum > 30) momentumColor = '#8b5cf6'; // Purple
  if (momentum > 60) momentumColor = '#f59e0b'; // Gold
  if (momentum > 80) momentumColor = '#ef4444'; // Red
  
  ctx.fillStyle = momentumColor;
  ctx.fillRect(momBarX, momBarY, momBarW * momentumPercent, momBarH);
  
  // Momentum border
  ctx.strokeStyle = '#475569';
  ctx.lineWidth = 0.5;
  ctx.strokeRect(momBarX, momBarY, momBarW, momBarH);

  // --- FLURRY INDICATOR ---
  if (entity.flurryActive || visuals?.flurryActive) {
    // Pulsing glow effect around character
    const pulse = Math.sin(time / 100) * 0.3 + 0.7;
    ctx.save();
    ctx.globalAlpha = pulse * 0.5;
    ctx.strokeStyle = '#f59e0b'; // Amber
    ctx.lineWidth = 2 * zoom;
    ctx.beginPath();
    ctx.arc(centerX, charY - 20 * zoom, 15 * zoom, 0, Math.PI * 2);
    ctx.stroke();
    ctx.restore();
  }

  // --- SUNDER STACKS INDICATOR ---
  const sunderStacks = entity.sunderStacks || 0;
  if (sunderStacks > 0) {
    const stackSize = 4 * zoom;
    const startX = centerX - (sunderStacks * stackSize) / 2;
    for (let i = 0; i < sunderStacks; i++) {
      ctx.fillStyle = '#ef4444'; // Red
      ctx.fillRect(startX + i * (stackSize + 1), momBarY - stackSize - 2, stackSize, stackSize);
    }
  }

  // --- COMBO COUNTER ---
  const combo = entity.consecutiveHits || 0;
  if (combo > 1) {
    ctx.save();
    ctx.font = `bold ${Math.round(10 * zoom)}px 'Courier New', monospace`;
    ctx.textAlign = 'center';
    ctx.fillStyle = combo >= 5 ? '#ef4444' : (combo >= 3 ? '#f59e0b' : '#fbbf24');
    ctx.fillText(`${combo}x`, centerX + barW / 2 + 10 * zoom, barY + barH);
    ctx.restore();
  }
};

const drawParticles = (ctx: CanvasRenderingContext2D, canvasW: number, canvasH: number, dt: number) => {
   for (let i = particles.value.length - 1; i >= 0; i--) {
      const p = particles.value[i];
      p.life -= dt;
      if (p.life <= 0) {
         particles.value.splice(i, 1);
         continue;
      }
      
      // Physics
      p.x += p.vx * dt * 0.06;
      p.y += p.vy * dt * 0.06;
      p.z += p.vz * dt * 0.06;
      p.vz -= 0.02; // Gravity
      
      if (p.z < 0) { p.z = 0; p.vx *= 0.5; p.vy *= 0.5; }

      const pos = toScreen(p.x, p.y, canvasW, canvasH);
      const zoom = camera.zoom;
      const screenY = pos.y - p.z * 10 * zoom;
      
      ctx.fillStyle = p.color;
      // Pixel Particles (squares)
      const size = Math.max(2, p.size * zoom);
      ctx.fillRect(pos.x - size/2, screenY - size/2, size, size);
   }
};

const drawEffects = (
  ctx: CanvasRenderingContext2D,
  state: RoundState,
  canvasW: number,
  canvasH: number,
  progress: number
) => {
  const IMPACT_TIME = 0.5;

  for (const attack of state.events.attacks) {
    const attacker = state.entities[attack.attackerId];
    const target = state.entities[attack.targetId];
    if (!attacker || !target) continue;

    const from = toScreen(attacker.x, attacker.y, canvasW, canvasH);
    const to = toScreen(target.x, target.y, canvasW, canvasH);
    // Adjust to height
    from.y -= 30 * camera.zoom;
    to.y -= 30 * camera.zoom;

    const isCrit = attack.isCritical;

    // 1. Projectile / Travel Phase
    if (progress < IMPACT_TIME) {
      const t = progress / IMPACT_TIME; // 0 to 1
      
      // Calculate current projectile position
      const curX = from.x + (to.x - from.x) * t;
      const curY = from.y + (to.y - from.y) * t;
      // Add a slight arc
      const arcHeight = Math.sin(t * Math.PI) * 40 * camera.zoom;
      
      // Draw 3D Energy Orb with Gradient
      const r = 8 * camera.zoom;
      const px = curX;
      const py = curY - arcHeight;

      // Create radial gradient for 3D sphere effect
      const grad = ctx.createRadialGradient(px - r/3, py - r/3, r/6, px, py, r * 1.2);
      if (isCrit) {
         // Critical: Golden/Amber Orb
         grad.addColorStop(0, '#ffffff');
         grad.addColorStop(0.2, '#fef3c7');
         grad.addColorStop(0.5, '#fbbf24');
         grad.addColorStop(0.8, '#d97706');
         grad.addColorStop(1, 'rgba(180, 83, 9, 0)');
      } else {
         // Normal: Blue Energy Orb
         grad.addColorStop(0, '#ffffff');
         grad.addColorStop(0.2, '#dbeafe');
         grad.addColorStop(0.5, '#60a5fa');
         grad.addColorStop(0.8, '#2563eb');
         grad.addColorStop(1, 'rgba(30, 64, 175, 0)');
      }

      // Outer glow
      ctx.save();
      ctx.shadowColor = isCrit ? '#fbbf24' : '#60a5fa';
      ctx.shadowBlur = 15 * camera.zoom;
      ctx.fillStyle = grad;
      ctx.beginPath();
      ctx.arc(px, py, r, 0, Math.PI * 2);
      ctx.fill();
      ctx.restore();

      // Inner bright core
      const coreGrad = ctx.createRadialGradient(px, py, 0, px, py, r * 0.4);
      coreGrad.addColorStop(0, 'rgba(255, 255, 255, 0.9)');
      coreGrad.addColorStop(1, 'rgba(255, 255, 255, 0)');
      ctx.fillStyle = coreGrad;
      ctx.beginPath();
      ctx.arc(px, py, r * 0.4, 0, Math.PI * 2);
      ctx.fill();
      
      // Trail with gradient
      if (t > 0.1) {
         const trailGrad = ctx.createLinearGradient(from.x, from.y, px, py);
         if (isCrit) {
            trailGrad.addColorStop(0, 'rgba(251, 191, 36, 0)');
            trailGrad.addColorStop(0.5, 'rgba(251, 191, 36, 0.3)');
            trailGrad.addColorStop(1, 'rgba(251, 191, 36, 0.8)');
         } else {
            trailGrad.addColorStop(0, 'rgba(96, 165, 250, 0)');
            trailGrad.addColorStop(0.5, 'rgba(96, 165, 250, 0.3)');
            trailGrad.addColorStop(1, 'rgba(96, 165, 250, 0.8)');
         }
         ctx.strokeStyle = trailGrad;
         ctx.lineWidth = 3 * camera.zoom;
         ctx.lineCap = 'round';
         ctx.beginPath();
         ctx.moveTo(from.x, from.y);
         ctx.quadraticCurveTo(
            from.x + (to.x - from.x) * 0.5, 
            from.y + (to.y - from.y) * 0.5 - 40 * camera.zoom, 
            px, py
         );
         ctx.stroke();
      }
    }

    // 2. Impact Phase - Procedural Explosion/Burst
    if (progress >= IMPACT_TIME && progress < IMPACT_TIME + 0.3) {
      const fxProgress = (progress - IMPACT_TIME) / 0.3; // 0 to 1
      const impactX = to.x;
      const impactY = to.y - 20 * camera.zoom;
      
      // Expanding ring
      const ringRadius = (isCrit ? 30 : 20) * camera.zoom * (0.3 + fxProgress * 0.7);
      const ringAlpha = 1 - fxProgress;
      
      ctx.save();
      ctx.globalAlpha = ringAlpha;
      
      // Outer burst ring
      const ringGrad = ctx.createRadialGradient(impactX, impactY, ringRadius * 0.5, impactX, impactY, ringRadius);
      if (isCrit) {
         ringGrad.addColorStop(0, 'rgba(251, 191, 36, 0)');
         ringGrad.addColorStop(0.6, 'rgba(251, 191, 36, 0.8)');
         ringGrad.addColorStop(0.8, 'rgba(245, 158, 11, 0.6)');
         ringGrad.addColorStop(1, 'rgba(217, 119, 6, 0)');
      } else {
         ringGrad.addColorStop(0, 'rgba(96, 165, 250, 0)');
         ringGrad.addColorStop(0.6, 'rgba(96, 165, 250, 0.8)');
         ringGrad.addColorStop(0.8, 'rgba(59, 130, 246, 0.6)');
         ringGrad.addColorStop(1, 'rgba(37, 99, 235, 0)');
      }
      ctx.fillStyle = ringGrad;
      ctx.beginPath();
      ctx.arc(impactX, impactY, ringRadius, 0, Math.PI * 2);
      ctx.fill();
      
      // Central flash (bright core that fades fast)
      if (fxProgress < 0.5) {
         const flashAlpha = 1 - (fxProgress / 0.5);
         const flashSize = (isCrit ? 15 : 10) * camera.zoom * (1 - fxProgress * 0.5);
         
         ctx.globalAlpha = flashAlpha;
         ctx.shadowColor = isCrit ? '#fbbf24' : '#60a5fa';
         ctx.shadowBlur = 20 * camera.zoom;
         
         const flashGrad = ctx.createRadialGradient(impactX, impactY, 0, impactX, impactY, flashSize);
         flashGrad.addColorStop(0, '#ffffff');
         flashGrad.addColorStop(0.5, isCrit ? '#fef3c7' : '#dbeafe');
         flashGrad.addColorStop(1, 'rgba(255, 255, 255, 0)');
         
         ctx.fillStyle = flashGrad;
         ctx.beginPath();
         ctx.arc(impactX, impactY, flashSize, 0, Math.PI * 2);
         ctx.fill();
      }
      
      // Spark particles (simple rays)
      const sparkCount = isCrit ? 8 : 5;
      const sparkLength = (isCrit ? 25 : 15) * camera.zoom * (0.5 + fxProgress * 0.5);
      ctx.globalAlpha = ringAlpha * 0.8;
      ctx.strokeStyle = isCrit ? '#fef3c7' : '#dbeafe';
      ctx.lineWidth = 2 * camera.zoom;
      ctx.lineCap = 'round';
      
      for (let i = 0; i < sparkCount; i++) {
         const angle = (Math.PI * 2 / sparkCount) * i + fxProgress * 0.5;
         const innerR = ringRadius * 0.3;
         const outerR = innerR + sparkLength * (1 - fxProgress * 0.5);
         
         ctx.beginPath();
         ctx.moveTo(
            impactX + Math.cos(angle) * innerR,
            impactY + Math.sin(angle) * innerR
         );
         ctx.lineTo(
            impactX + Math.cos(angle) * outerR,
            impactY + Math.sin(angle) * outerR
         );
         ctx.stroke();
      }
      
      ctx.restore();
    }

    // 3. Floating Combat Text
    if (progress >= IMPACT_TIME) {
      const localProgress = (progress - IMPACT_TIME) / (1 - IMPACT_TIME); // 0 to 1
      const t = Math.min(1, localProgress * 3); 
      const ease = 1 - Math.pow(1 - t, 3);
      
      const floatY = ease * 60 * camera.zoom; 
      const alpha = Math.max(0, 1 - localProgress * 1.5);
      
      ctx.save();
      ctx.globalAlpha = alpha;
      const fontSize = Math.round((attack.isCritical ? 24 : 14) * camera.zoom);
      ctx.font = `bold ${fontSize}px 'Courier New', monospace`; 
      ctx.textAlign = 'center';
      
      const txt = `${attack.damage}`;
      const numX = to.x;
      const numY = to.y - 30 * camera.zoom - floatY;
      
      // Shadow/Outline
      ctx.lineWidth = 3;
      ctx.strokeStyle = 'black';
      ctx.strokeText(txt, numX, numY);
      
      // Fill
      ctx.fillStyle = attack.isCritical ? '#fbbf24' : '#fff';
      ctx.fillText(txt, numX, numY);
      
      if (attack.isCritical) {
        ctx.font = `bold ${Math.round(10 * camera.zoom)}px 'Courier New', monospace`;
        ctx.fillStyle = '#f59e0b';
        ctx.fillText('CRIT!', numX, numY - 20 * camera.zoom);
      }
      ctx.restore();
    }
  }
};

const drawScene = (time: number) => {
  if (!canvasRef.value) return;
  const canvas = canvasRef.value;
  const ctx = canvas.getContext('2d');
  if (!ctx) return;
  
  // Responsive resize
  const rect = canvas.getBoundingClientRect();
  if (canvas.width !== rect.width || canvas.height !== rect.height) {
    canvas.width = rect.width;
    canvas.height = rect.height;
  }

  // Camera Shake Decay
  camera.shakeX *= 0.9;
  camera.shakeY *= 0.9;
  if (Math.abs(camera.shakeX) < 0.5) camera.shakeX = 0;
  if (Math.abs(camera.shakeY) < 0.5) camera.shakeY = 0;

  // Determine current playhead state
  let currentRoundNum = selectedRound.value;
  let progress = 0;
  
  if (isPlaying.value && segmentStart > 0) {
     const duration = 1000 / playbackSpeed.value;
     progress = Math.min(1, (time - segmentStart) / duration);
     
     // Loop particle spawning logic (simplified)
     // In a real game loop this would be in update()
     if (Math.random() < 0.05) { // dust
       // spawn ambient particle
     }
  }

  const state = roundStateMap.value[currentRoundNum];
  if (!state) return;
  
  const nextRoundIndex = orderedRounds.value.indexOf(currentRoundNum) + 1;
  const nextRound = orderedRounds.value[nextRoundIndex];
  const nextState = nextRound ? roundStateMap.value[nextRound] : undefined;

  // Clear
  ctx.clearRect(0, 0, canvas.width, canvas.height);

  // 1. Draw Grid
  for (let x = 0; x <= worldSize; x++) {
    for (let y = 0; y <= worldSize; y++) {
       drawTile(ctx, x, y, canvas.width, canvas.height, time);
    }
  }

  // 2. Interpolate & Render Entities
  const interpolated = getInterpolatedEntities(state, nextState, progress);
  
  // Detect Hit States
  const hitEntities = new Set<string>();
  const attackingEntities = new Set<string>();
  const flurryEntities = new Set<string>();
  const entityMomentumEvents: Record<string, MomentumEvent> = {};
  
  if (state.events?.attacks) {
     const IMPACT_TIME = 0.5;
     if (progress >= IMPACT_TIME && progress < IMPACT_TIME + 0.2) {
        state.events.attacks.forEach(a => hitEntities.add(a.targetId));
     }
     // Attack animation window (start of round)
     if (progress < 0.4) {
        state.events.attacks.forEach(a => attackingEntities.add(a.attackerId));
     }
  }
  
  // Track momentum events
  if (state.events?.momentum) {
    state.events.momentum.forEach(m => {
      entityMomentumEvents[m.fighterId] = m;
    });
  }
  
  // Track flurry activations
  if (state.events?.flurry) {
    state.events.flurry.forEach(f => flurryEntities.add(f.fighterId));
  }

  const renderList = Object.values(interpolated).map(e => ({
    type: 'entity',
    yDepth: e.x + e.y, 
    entity: e,
    x: e.x,
    y: e.y,
    hit: hitEntities.has(e.id),
    attacking: attackingEntities.has(e.id),
    momentumEvent: entityMomentumEvents[e.id],
    flurryActive: flurryEntities.has(e.id)
  }));

  renderList.sort((a, b) => a.yDepth - b.yDepth);

  for (const item of renderList) {
    if (item.type === 'entity') {
      drawEntity(ctx, item.entity, item.x, item.y, canvas.width, canvas.height, time, { 
        hit: item.hit, 
        attacking: item.attacking,
        momentumEvent: item.momentumEvent,
        flurryActive: item.flurryActive
      });
    }
  }

  // 3. VFX
  drawParticles(ctx, canvas.width, canvas.height, 16); // assume ~60fps dt=16ms
  drawEffects(ctx, state, canvas.width, canvas.height, progress);
};


const getInterpolatedEntities = (state: RoundState, nextState: RoundState | undefined, progress: number) => {
  if (!nextState || progress <= 0) return state.entities;
  const blended: Record<string, EntityState> = {};
  const ids = new Set([...Object.keys(state.entities), ...Object.keys(nextState.entities)]);
  ids.forEach((id) => {
    const from = state.entities[id];
    const to = nextState.entities[id];
    if (from && to) {
      const dx = to.x - from.x;
      const dy = to.y - from.y;
      const isMoving = progress < 1.0 && (Math.abs(dx) > 0.001 || Math.abs(dy) > 0.001);
      
      blended[id] = {
        ...from,
        x: from.x + dx * progress,
        y: from.y + dy * progress,
        hp: from.hp + (to.hp - from.hp) * progress,
        alive: progress < 0.8 ? from.alive : to.alive,
        floatOffset: from.floatOffset,
        vx: dx,
        vy: dy,
        isMoving
      };
    } else if (to) {
      blended[id] = { ...to, isMoving: false, vx: 0, vy: 0 };
    } else if (from) {
      blended[id] = { ...from, isMoving: false, vx: 0, vy: 0 };
    }
  });
  return blended;
};

// --- ANIMATION LOOP ---
// Separated visual loop from logic loop to allow idle anims when paused
const visualLoop = (time: number) => {
   drawScene(time);
   
   // Handle playback logic here to sync everything
   if (isPlaying.value) {
      if (!segmentStart) segmentStart = time;
      const duration = 1000 / playbackSpeed.value;
      const progress = (time - segmentStart) / duration;
      
      // Check for screen shake trigger on Crit (approximate timing)
      if (progress > 0.1 && progress < 0.15) {
         const state = roundStateMap.value[selectedRound.value];
         if (state?.events.attacks.some(a => a.isCritical)) {
            addShake(2 * camera.zoom);
         }
      }

      if (progress >= 1) {
         const index = orderedRounds.value.indexOf(selectedRound.value);
         const nextIndex = index + 1;
         if (nextIndex >= orderedRounds.value.length) {
            isPlaying.value = false;
            segmentStart = 0;
         } else {
            selectedRound.value = orderedRounds.value[nextIndex];
            segmentStart = time;
            scrollToLog(selectedRound.value);
         }
      }
   }
   
   visualLoopHandle = requestAnimationFrame(visualLoop);
};


// Logic helpers
const renderSelectedRound = () => {
   // No-op, visual loop handles it
};

const stepRound = (direction: 1 | -1) => {
  const index = orderedRounds.value.indexOf(selectedRound.value);
  if (index === -1) return;
  const nextIndex = index + direction;
  if (nextIndex < 0 || nextIndex >= orderedRounds.value.length) return;
  selectedRound.value = orderedRounds.value[nextIndex];
  scrollToLog(selectedRound.value);
};

const scrollToLog = (round: number) => {
  nextTick(() => {
    const el = document.getElementById('log-round-' + round);
    if (el && logsContainer.value) {
       const top = el.offsetTop - logsContainer.value.offsetTop;
       logsContainer.value.scrollTo({ top: top - 100, behavior: 'smooth' });
    }
  });
};

const seekToPercent = (e: MouseEvent) => {
  if (!orderedRounds.value.length) return;
  const rect = (e.currentTarget as HTMLElement).getBoundingClientRect();
  const x = e.clientX - rect.left;
  const p = Math.max(0, Math.min(1, x / rect.width));
  const idx = Math.floor(p * (orderedRounds.value.length - 1));
  selectedRound.value = orderedRounds.value[idx];
  scrollToLog(selectedRound.value);
};

const selectRound = (round: number) => {
  selectedRound.value = round;
  isPlaying.value = false;
};

const togglePlayback = () => {
  if (orderedRounds.value.length <= 1) return;
  if (isPlaying.value) {
    isPlaying.value = false;
  } else {
    if (selectedRound.value === orderedRounds.value[orderedRounds.value.length - 1]) {
       selectedRound.value = orderedRounds.value[0];
    }
    isPlaying.value = true;
    segmentStart = 0; // Will be reset in loop
  }
};

const isWinner = computed(() => {
  if (!roster.fighters.length) return false;
  const playerIds = new Set(roster.fighters.map(f => f.id));
  const finalState = roundStateMap.value[orderedRounds.value[orderedRounds.value.length - 1]];
  if (!finalState) return false;
  
  // Check if any of our fighters are still alive in the final state
  return Object.values(finalState.entities).some(e => playerIds.has(e.id) && e.alive);
});

const claimAndExit = async () => {
   try {
     await rewardsStore.claimAll();
     window.location.href = '/matches';
   } catch (e) {
     console.error("Failed to claim rewards:", e);
   }
};

// Data Fetching
const fetchLogs = async () => {
  if (!auth.token || !matchId.value) return;
  isLoading.value = true;
  try {
    const data = await request<any[]>(`${endpoints.match}/${matchId.value}/roundticks`, { token: auth.token });
    const nextTicks = data || [];
    const prevCount = ticks.value.length;
    ticks.value = nextTicks;
    refreshRoundStates();
    
    // Auto-advance if live
    if (matchStatus.value === 'running' && nextTicks.length > prevCount && nextTicks.length > 0) {
      const lastRound = nextTicks[nextTicks.length - 1]?.round;
      if (lastRound != null) {
         selectedRound.value = lastRound;
         scrollToLog(lastRound);
      }
    }
  } catch (e) { console.error(e); } 
  finally { isLoading.value = false; }
};

const fetchMatchStatus = async () => {
  if (!auth.token || !matchId.value) return;
  try {
    const m = await request<{ status?: string }>(`${endpoints.match}/${matchId.value}`, { token: auth.token });
    matchStatus.value = m?.status ?? null;
    return m?.status;
  } catch (e) {
    matchStatus.value = null;
    return null;
  }
};

const refreshRoundStates = () => {
  const playerIds = new Set(roster.fighters.map((fighter) => fighter.id));
  const entityMap = new Map<string, EntityState>();
  const states: Record<number, RoundState> = {};

  for (const round of ticks.value) {
    const events = { 
      attacks: [] as AttackEvent[], 
      deaths: [] as string[],
      momentum: [] as MomentumEvent[],
      sunder: [] as SunderEvent[],
      flurry: [] as FlurryEvent[]
    };
    
    for (const tick of round.ticks || []) {
      const payload = tick.payload as Record<string, any> | undefined;
      if (!payload) continue;
      
      if (tick.type === 'spawn') {
        const id = payload.fighterId as string;
        entityMap.set(id, {
          id,
          x: payload.x ?? 0,
          y: payload.y ?? 0,
          hp: payload.hp ?? 100,
          maxHp: payload.hp ?? 100,
          alive: true,
          isPlayer: playerIds.has(id),
          floatOffset: Math.random() * Math.PI * 2, // Random idle phase
          momentum: 0,
          consecutiveHits: 0,
          sunderStacks: 0,
          flurryActive: false
        });
      } else if (tick.type === 'move') {
        const id = payload.fighterId as string;
        const existing = entityMap.get(id);
        if (existing) {
          existing.x = payload.toX ?? existing.x;
          existing.y = payload.toY ?? existing.y;
        }
      } else if (tick.type === 'attack') {
        events.attacks.push({
          attackerId: payload.attackerId,
          targetId: payload.targetId,
          isCritical: payload.isCritical,
          isParried: payload.isParried,
          damage: payload.damage
        });
        const target = entityMap.get(payload.targetId);
        if (target) {
          target.hp = Math.max(0, target.hp - (payload.damage || 0));
        }
      } else if (tick.type === 'died') {
        const id = payload.fighterId as string;
        const existing = entityMap.get(id);
        if (existing) {
          existing.alive = false;
          existing.hp = 0;
        }
        events.deaths.push(id);
      } else if (tick.type === 'momentum' || tick.type === 'momentum_decay') {
        const event: MomentumEvent = {
          fighterId: payload.fighterId,
          momentum: payload.momentum,
          consecutiveHits: payload.consecutiveHits,
          targetId: payload.targetId
        };
        events.momentum.push(event);
        
        // Update entity state
        const entity = entityMap.get(payload.fighterId);
        if (entity) {
          entity.momentum = payload.momentum;
          entity.consecutiveHits = payload.consecutiveHits;
        }
      } else if (tick.type === 'sunder') {
        const event: SunderEvent = {
          targetId: payload.targetId,
          stacks: payload.stacks,
          armorReduced: payload.armorReduced
        };
        events.sunder.push(event);
        
        // Update target's sunder stacks
        const target = entityMap.get(payload.targetId);
        if (target) {
          target.sunderStacks = payload.stacks;
        }
      } else if (tick.type === 'flurry') {
        const event: FlurryEvent = {
          fighterId: payload.fighterId,
          attackSpeedBonus: payload.attackSpeedBonus
        };
        events.flurry.push(event);
        
        // Activate flurry on entity
        const entity = entityMap.get(payload.fighterId);
        if (entity) {
          entity.flurryActive = true;
        }
      }
    }
    
    // Deep clone for snapshot
    const snapshot: Record<string, EntityState> = {};
    entityMap.forEach((value, key) => { snapshot[key] = { ...value }; });
    states[round.round] = { round: round.round, entities: snapshot, events };
  }
  roundStateMap.value = states;
};

// Lifecycle
function startLivePoll() {
  if (livePollHandle) return;
  livePollHandle = window.setInterval(async () => {
    const status = await fetchMatchStatus();
    await fetchLogs();
    if (status === 'completed') stopLivePoll();
  }, 2500);
}

function stopLivePoll() {
  if (livePollHandle) {
    clearInterval(livePollHandle);
    livePollHandle = null;
  }
}

onMounted(async () => {
  await roster.fetchFighters();
  await rewardsStore.fetchRewards();
  loadAssets(); // Start loading sprites
  await fetchMatchStatus();
  await fetchLogs();
  
  if (orderedRounds.value.length) {
    if (matchStatus.value === 'completed') {
       selectedRound.value = orderedRounds.value[0];
    } else {
       selectedRound.value = orderedRounds.value[orderedRounds.value.length - 1];
    }
  }
  
  // Start Visual Loop
  visualLoopHandle = requestAnimationFrame(visualLoop);
  
  if (matchStatus.value === 'running') startLivePoll();
});

onUnmounted(() => {
  if (visualLoopHandle) cancelAnimationFrame(visualLoopHandle);
  stopLivePoll();
});
</script>

<style scoped>
.pixel-theme {
  image-rendering: pixelated;
  background-repeat: repeat;
  background-size: 128px;
}

.pixelated {
  image-rendering: pixelated;
}

.pixel-box {
  box-shadow: 4px 4px 0px 0px rgba(0,0,0,0.5);
}

.text-shadow-retro {
  text-shadow: 3px 3px 0 #000, -1px -1px 0 #000;
}

.rpg-btn-small {
  border-bottom-width: 4px;
  border-right-width: 4px;
  transition: all 0.1s;
}

.rpg-btn-small:active {
  border-bottom-width: 0px;
  border-right-width: 0px;
  transform: translate(4px, 4px);
}

.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: #0f172a;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #334155;
  border: 1px solid #1e293b;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #475569;
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 0.5s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}

.blink {
  animation: blink 1s step-end infinite;
}
@keyframes blink {
  50% { opacity: 0; }
}

.animate-bounce-slow {
  animation: bounce 2s infinite;
}
@keyframes bounce {
  0%, 100% { transform: translateY(-5%); }
  50% { transform: translateY(5%); }
}
</style>
