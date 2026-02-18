<template>
  <div class="pixel-theme h-svh max-h-svh overflow-hidden flex flex-col font-mono text-slate-200" :style="{ backgroundImage: `url('${PIXEL_ASSETS.BG_DUNGEON}')` }">
    <div class="fixed inset-0 bg-slate-950/80 pointer-events-none z-0"></div>

    <!-- Header -->
    <header class="relative z-20 flex items-center justify-between shrink-0 px-2 py-1.5 bg-slate-900/90 border-b-2 border-amber-900/50">
      <div class="flex items-center gap-2">
        <button @click="$router.push('/matches')" class="text-amber-500 hover:text-amber-400 font-bold text-xs">
          ← BACK
        </button>
      </div>
      <div class="flex items-center gap-2">
        <span v-if="matchStatus === 'running'" class="flex items-center gap-1">
          <span class="w-1.5 h-1.5 bg-red-500 rounded-full animate-pulse"></span>
          <span class="text-[8px] text-red-400 font-bold uppercase">LIVE</span>
        </span>
        <span class="text-[8px] text-slate-500 font-mono">#{{ matchId?.substring(0, 4) }}</span>
      </div>
    </header>

    <!-- Canvas Area - Full Height -->
    <div class="relative z-10 flex-1 min-h-0">
      <canvas
        ref="canvasRef"
        class="w-full h-full object-cover cursor-move pixelated"
        @mousedown="startDrag"
        @mousemove="onDrag"
        @mouseup="endDrag"
        @mouseleave="endDrag"
        @wheel.prevent="onWheel"
      ></canvas>
      
      <!-- Vignette -->
      <div class="absolute inset-0 pointer-events-none bg-[radial-gradient(circle_at_center,transparent_40%,rgba(0,0,0,0.6)_100%)]"></div>

      <!-- Loading Overlay -->
      <div v-if="isLoading" class="absolute inset-0 z-10 flex items-center justify-center bg-slate-950/80">
        <div class="flex flex-col items-center gap-2">
          <div class="h-6 w-6 animate-spin border-2 border-amber-500 border-t-transparent rounded-full"></div>
          <p class="text-amber-200 text-xs font-bold uppercase">Loading...</p>
        </div>
      </div>

      <!-- Victory Overlay -->
      <div v-if="matchStatus === 'completed' && orderedRounds.length && selectedRound === orderedRounds[orderedRounds.length-1]" 
           class="absolute inset-0 z-10 flex items-center justify-center bg-black/60">
        <div class="border-2 border-amber-500 bg-slate-900 p-4 text-center">
          <img :src="PIXEL_ASSETS.ICON_TROPHY" class="w-8 h-8 mx-auto mb-2 pixelated" />
          <p class="text-xl font-black text-amber-500 uppercase">VICTORY</p>
        </div>
      </div>
    </div>

    <!-- Bottom Controls - Mobile Optimized TINY -->
    <div class="relative z-20 bg-slate-900/95 border-t-2 border-slate-700 px-1 py-1">
      <!-- Progress Bar -->
      <div class="relative h-2 bg-slate-950 border border-slate-700 cursor-pointer mb-1" @click="seekToPercent">
        <div class="absolute inset-y-0 left-0 bg-amber-600" :style="{ width: progressPercent + '%' }"></div>
        <div class="absolute inset-y-0 w-1 bg-white" :style="{ left: progressPercent + '%' }"></div>
      </div>

      <!-- Controls Row -->
      <div class="flex items-center justify-between gap-1">
        <!-- Left: TINY Playback -->
        <div class="flex items-center gap-0.5">
          <button @click="stepRound(-1)" class="w-6 h-6 flex items-center justify-center bg-slate-800 border border-slate-600 text-slate-400 hover:text-white text-[8px] font-bold">
            ◀
          </button>
          <button @click="togglePlayback" class="w-7 h-7 flex items-center justify-center bg-amber-600 border border-amber-800 text-slate-900 font-black hover:bg-amber-500 text-[8px]">
            <span v-if="isPlaying">||</span>
            <span v-else>▶</span>
          </button>
          <button @click="stepRound(1)" class="w-6 h-6 flex items-center justify-center bg-slate-800 border border-slate-600 text-slate-400 hover:text-white text-[8px] font-bold">
            ▶
          </button>
        </div>

        <!-- Center: Round Info -->
        <div class="text-[8px] text-amber-400 font-mono font-bold bg-slate-950 px-1.5 py-0.5 border border-slate-700">
          R{{ selectedRound }}/{{ orderedRounds.length || 0 }}
        </div>

        <!-- Right: TINY Actions -->
        <div class="flex items-center gap-0.5">
          <!-- Log Toggle -->
          <button @click="showLog = !showLog" class="w-6 h-6 flex items-center justify-center bg-slate-800 border border-slate-600 text-slate-400 hover:text-amber-400 text-[8px]">
            📜
          </button>
          <!-- Config Toggle -->
          <button @click="showConfig = !showConfig" class="w-6 h-6 flex items-center justify-center bg-slate-800 border border-slate-600 text-slate-400 hover:text-amber-400 text-[8px]">
            ⚙️
          </button>
        </div>
      </div>

      <!-- Hidden Config Panel -->
      <transition name="slide">
        <div v-if="showConfig" class="mt-1 pt-1 border-t border-slate-700 flex items-center justify-between gap-2">
          <div class="flex items-center gap-1">
            <span class="text-[7px] text-slate-500 uppercase font-bold">Speed</span>
            <input type="range" min="0.5" max="4" step="0.5" v-model.number="playbackSpeed" class="w-12 accent-amber-500 h-1 bg-slate-800" />
            <span class="text-[8px] text-amber-400 font-mono font-bold w-6">{{ playbackSpeed }}x</span>
          </div>
          <div class="flex items-center gap-1">
            <span class="text-[7px] text-slate-500 uppercase font-bold">Zoom</span>
            <button @click="zoom(-0.2)" class="w-5 h-5 flex items-center justify-center bg-slate-800 border border-slate-600 text-slate-400 hover:text-white font-bold text-[8px]">-</button>
            <span class="text-[8px] text-amber-400 font-mono font-bold w-8 text-center">{{ Math.round(camera.zoom * 100) }}%</span>
            <button @click="zoom(0.2)" class="w-5 h-5 flex items-center justify-center bg-slate-800 border border-slate-600 text-slate-400 hover:text-white font-bold text-[8px]">+</button>
          </div>
        </div>
      </transition>
    </div>

    <!-- Log Modal Overlay - 30svh -->
    <transition name="slide-up">
      <div v-if="showLog" class="fixed inset-x-0 bottom-0 z-30 h-[30svh] bg-slate-900 border-t-4 border-amber-600 flex flex-col shadow-2xl">
        <!-- Drag Handle -->
        <div class="flex justify-center py-1 cursor-pointer" @click="showLog = false">
          <div class="w-12 h-1 bg-slate-600 rounded-full"></div>
        </div>
        
        <!-- Header -->
        <div class="flex items-center justify-between px-2 py-1 border-b border-slate-700 shrink-0">
          <span class="text-amber-500 font-bold text-xs uppercase">Combat Log</span>
          <span class="text-[8px] text-slate-500 font-mono">{{ ticks.length }} Rounds</span>
        </div>
        
        <!-- Log Content -->
        <div class="flex-1 overflow-y-auto px-1 pb-1 space-y-1 custom-scrollbar">
          <div 
            v-for="round in ticks" 
            :key="round.round"
            @click="selectRound(round.round); showLog = false;"
            class="p-1.5 text-[9px] cursor-pointer border transition-all"
            :class="round.round === selectedRound ? 'border-amber-600 bg-amber-900/20' : 'border-slate-700 hover:border-slate-500'"
          >
            <div class="flex items-center justify-between mb-0.5">
              <span class="font-bold text-amber-100 uppercase">R{{ round.round.toString().padStart(2, '0') }}</span>
              <span class="text-slate-500">{{ round.ticks?.length || 0 }} events</span>
            </div>
            <div class="text-slate-400 space-y-0.5">
              <div v-for="(tick, idx) in (round.ticks || []).slice(0, 3)" :key="idx" class="truncate">
                <template v-if="tick.type === 'attack'">
                  <span class="text-emerald-400">{{ payloadValue(tick.payload, 'attackerId')?.substring(0,4) }}</span>
                  →
                  <span class="text-rose-400">{{ payloadValue(tick.payload, 'targetId')?.substring(0,4) }}</span>
                  <span class="text-slate-300"> -{{ payloadValue(tick.payload, 'damage') }}</span>
                </template>
                <template v-else-if="tick.type === 'died'">
                  <span class="text-rose-500">☠ {{ payloadValue(tick.payload, 'fighterId')?.substring(0,4) }}</span>
                </template>
              </div>
              <div v-if="(round.ticks?.length || 0) > 3" class="text-slate-600">+{{ (round.ticks?.length || 0) - 3 }} more</div>
            </div>
          </div>
        </div>
      </div>
    </transition>

    <!-- Log Backdrop -->
    <transition name="fade">
      <div v-if="showLog" @click="showLog = false" class="fixed inset-0 z-25 bg-black/50"></div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, reactive, computed } from 'vue';
import { useRoute } from 'vue-router';
import { request } from '@/shared/api/http';
import { endpoints } from '@/shared/api/endpoints';
import { useAuthStore } from '@/features/auth/store';
import { useRosterStore } from '@/features/roster/store';

const PIXEL_ASSETS = {
  BG_DUNGEON: 'https://vibemedia.space/bg_dungeon_v2_99283.png?prompt=dark%20dungeon%20floor%20tile%20texture%20seamless&style=pixel_game_asset&key=NOGON',
  ICON_TROPHY: 'https://vibemedia.space/trophy_icon_4d5e6f_v1.png?prompt=golden%20trophy%20pixel%20art&style=pixel_game_asset&key=NOGON',
};

const auth = useAuthStore();
const roster = useRosterStore();
const route = useRoute();

const matchId = ref(route.params.id as string);
const matchStatus = ref<string | null>(null);
const ticks = ref<any[]>([]);
const isLoading = ref(false);

const showLog = ref(false);
const showConfig = ref(false);

const canvasRef = ref<HTMLCanvasElement | null>(null);
const roundStateMap = ref<Record<number, any>>({});

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

const selectedRound = ref(0);
const isPlaying = ref(false);
const playbackSpeed = ref(1.0);
let segmentStart = 0;
let visualLoopHandle: number | null = null;
let livePollHandle: number | null = null;

const orderedRounds = computed(() => ticks.value.map((r) => r.round));

const progressPercent = computed(() => {
  if (!orderedRounds.value.length) return 0;
  const idx = orderedRounds.value.indexOf(selectedRound.value);
  return ((idx + 1) / orderedRounds.value.length) * 100;
});

const toIso = (x: number, y: number) => ({
  x: (x - y) * TILE_WIDTH * 0.5,
  y: (x + y) * TILE_HEIGHT * 0.5
});

const toScreen = (x: number, y: number, canvasW: number, canvasH: number) => {
  const iso = toIso(x, y);
  const centerX = canvasW / 2;
  const centerY = canvasH / 4;
  return {
    x: centerX + (iso.x * camera.zoom) + camera.x + camera.shakeX,
    y: centerY + (iso.y * camera.zoom) + camera.y + camera.shakeY
  };
};

const payloadValue = (payload: unknown, key: string) => {
  if (!payload || typeof payload !== 'object') return undefined;
  return (payload as Record<string, any>)[key];
};

const startDrag = (e: MouseEvent) => {
  camera.isDragging = true;
  camera.lastX = e.clientX;
  camera.lastY = e.clientY;
};

const onDrag = (e: MouseEvent) => {
  if (!camera.isDragging) return;
  camera.x += e.clientX - camera.lastX;
  camera.y += e.clientY - camera.lastY;
  camera.lastX = e.clientX;
  camera.lastY = e.clientY;
};

const endDrag = () => { camera.isDragging = false; };

const onWheel = (e: WheelEvent) => {
  const delta = -Math.sign(e.deltaY) * 0.2;
  camera.zoom = Math.max(0.3, Math.min(2.5, camera.zoom + delta));
};

const zoom = (delta: number) => {
  camera.zoom = Math.max(0.3, Math.min(2.5, camera.zoom + delta));
};

const stepRound = (dir: 1 | -1) => {
  const idx = orderedRounds.value.indexOf(selectedRound.value);
  if (idx === -1) return;
  const next = idx + dir;
  if (next >= 0 && next < orderedRounds.value.length) {
    selectedRound.value = orderedRounds.value[next];
  }
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
    segmentStart = 0;
  }
};

const seekToPercent = (e: MouseEvent) => {
  if (!orderedRounds.value.length) return;
  const rect = (e.currentTarget as HTMLElement).getBoundingClientRect();
  const p = Math.max(0, Math.min(1, (e.clientX - rect.left) / rect.width));
  const idx = Math.floor(p * (orderedRounds.value.length - 1));
  selectedRound.value = orderedRounds.value[idx];
};

// Drawing
const drawTile = (ctx: CanvasRenderingContext2D, x: number, y: number, canvasW: number, canvasH: number) => {
  const pos = toScreen(x, y, canvasW, canvasH);
  const w = TILE_WIDTH * camera.zoom;
  const h = TILE_HEIGHT * camera.zoom;
  
  ctx.beginPath();
  ctx.moveTo(pos.x, pos.y);
  ctx.lineTo(pos.x + w/2, pos.y + h/2);
  ctx.lineTo(pos.x, pos.y + h);
  ctx.lineTo(pos.x - w/2, pos.y + h/2);
  ctx.closePath();
  ctx.fillStyle = '#1e293b';
  ctx.fill();
  ctx.strokeStyle = '#334155';
  ctx.lineWidth = 1;
  ctx.stroke();
};

const drawEntity = (ctx: CanvasRenderingContext2D, entity: any, x: number, y: number, canvasW: number, canvasH: number) => {
  const pos = toScreen(x, y, canvasW, canvasH);
  const zoom = camera.zoom;
  
  if (!entity.alive) {
    ctx.fillStyle = '#334155';
    ctx.fillRect(pos.x - 5 * zoom, pos.y, 10 * zoom, 10 * zoom);
    return;
  }

  // Shadow
  ctx.fillStyle = 'rgba(0,0,0,0.5)';
  ctx.beginPath();
  ctx.ellipse(pos.x, pos.y + 14 * zoom, 10 * zoom, 5 * zoom, 0, 0, Math.PI * 2);
  ctx.fill();

  // Simple voxel character
  const colors = entity.isPlayer 
    ? { head: '#fbbf24', body: '#3b82f6' }
    : { head: '#4ade80', body: '#166534' };

  const bob = Math.sin(Date.now() / 300 + (entity.floatOffset || 0)) * 2 * zoom;
  
  // Body
  ctx.fillStyle = colors.body;
  ctx.fillRect(pos.x - 8 * zoom, pos.y - 20 * zoom + bob, 16 * zoom, 20 * zoom);
  
  // Head
  ctx.fillStyle = colors.head;
  ctx.fillRect(pos.x - 6 * zoom, pos.y - 32 * zoom + bob, 12 * zoom, 12 * zoom);

  // HP Bar
  const hpPercent = Math.max(0, entity.hp / entity.maxHp);
  ctx.fillStyle = '#000';
  ctx.fillRect(pos.x - 10 * zoom, pos.y - 40 * zoom + bob, 20 * zoom, 3 * zoom);
  ctx.fillStyle = hpPercent > 0.5 ? '#10b981' : hpPercent > 0.25 ? '#f59e0b' : '#ef4444';
  ctx.fillRect(pos.x - 10 * zoom, pos.y - 40 * zoom + bob, 20 * zoom * hpPercent, 3 * zoom);
};

const drawScene = () => {
  if (!canvasRef.value) return;
  const canvas = canvasRef.value;
  const ctx = canvas.getContext('2d');
  if (!ctx) return;

  const rect = canvas.getBoundingClientRect();
  if (canvas.width !== rect.width || canvas.height !== rect.height) {
    canvas.width = rect.width;
    canvas.height = rect.height;
  }

  ctx.clearRect(0, 0, canvas.width, canvas.height);

  // Grid
  for (let x = 0; x <= worldSize; x++) {
    for (let y = 0; y <= worldSize; y++) {
      drawTile(ctx, x, y, canvas.width, canvas.height);
    }
  }

  // Entities
  const state = roundStateMap.value[selectedRound.value];
  if (state) {
    Object.values(state.entities || {}).forEach((entity: any) => {
      drawEntity(ctx, entity, entity.x, entity.y, canvas.width, canvas.height);
    });
  }
};

const visualLoop = () => {
  drawScene();
  
  if (isPlaying.value && orderedRounds.value.length > 1) {
    if (!segmentStart) segmentStart = Date.now();
    const duration = 1500 / playbackSpeed.value;
    const elapsed = Date.now() - segmentStart;
    
    if (elapsed >= duration) {
      const idx = orderedRounds.value.indexOf(selectedRound.value);
      const nextIdx = idx + 1;
      if (nextIdx >= orderedRounds.value.length) {
        isPlaying.value = false;
        segmentStart = 0;
      } else {
        selectedRound.value = orderedRounds.value[nextIdx];
        segmentStart = Date.now();
      }
    }
  }
  
  visualLoopHandle = requestAnimationFrame(visualLoop);
};

const fetchLogs = async () => {
  if (!auth.token || !matchId.value) return;
  isLoading.value = true;
  try {
    const data = await request<any[]>(`${endpoints.match}/${matchId.value}/roundticks`, { token: auth.token });
    ticks.value = data || [];
    refreshRoundStates();
  } catch (e) { console.error(e); }
  finally { isLoading.value = false; }
};

const fetchMatchStatus = async () => {
  if (!auth.token || !matchId.value) return;
  try {
    const m = await request<{ status?: string }>(`${endpoints.match}/${matchId.value}`, { token: auth.token });
    matchStatus.value = m?.status ?? null;
  } catch (e) { matchStatus.value = null; }
};

const refreshRoundStates = () => {
  const playerIds = new Set(roster.fighters.map(f => f.id));
  const states: Record<number, any> = {};

  for (const round of ticks.value) {
    const entities: Record<string, any> = {};
    
    for (const tick of round.ticks || []) {
      const payload = tick.payload || {};
      
      if (tick.type === 'spawn') {
        entities[payload.fighterId] = {
          id: payload.fighterId,
          x: payload.x ?? 0,
          y: payload.y ?? 0,
          hp: payload.hp ?? 100,
          maxHp: payload.hp ?? 100,
          alive: true,
          isPlayer: playerIds.has(payload.fighterId),
          floatOffset: Math.random() * Math.PI * 2
        };
      } else if (tick.type === 'move') {
        if (entities[payload.fighterId]) {
          entities[payload.fighterId].x = payload.toX ?? entities[payload.fighterId].x;
          entities[payload.fighterId].y = payload.toY ?? entities[payload.fighterId].y;
        }
      } else if (tick.type === 'attack') {
        if (entities[payload.targetId]) {
          entities[payload.targetId].hp = Math.max(0, entities[payload.targetId].hp - (payload.damage || 0));
        }
      } else if (tick.type === 'died') {
        if (entities[payload.fighterId]) {
          entities[payload.fighterId].alive = false;
          entities[payload.fighterId].hp = 0;
        }
      }
    }
    
    states[round.round] = { entities };
  }
  
  roundStateMap.value = states;
};

const startLivePoll = () => {
  if (livePollHandle) return;
  livePollHandle = window.setInterval(async () => {
    await fetchMatchStatus();
    await fetchLogs();
    if (matchStatus.value === 'completed') stopLivePoll();
  }, 3000);
};

const stopLivePoll = () => {
  if (livePollHandle) {
    clearInterval(livePollHandle);
    livePollHandle = null;
  }
};

onMounted(async () => {
  await roster.fetchFighters();
  await fetchMatchStatus();
  await fetchLogs();
  
  if (orderedRounds.value.length) {
    selectedRound.value = matchStatus.value === 'completed' 
      ? orderedRounds.value[0]
      : orderedRounds.value[orderedRounds.value.length - 1];
  }
  
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

.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: #0f172a;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #475569;
  border-radius: 2px;
}

.slide-enter-active, .slide-leave-active {
  transition: all 0.2s ease;
}
.slide-enter-from, .slide-leave-to {
  opacity: 0;
  transform: translateY(10px);
}

.slide-up-enter-active, .slide-up-leave-active {
  transition: all 0.3s ease;
}
.slide-up-enter-from, .slide-up-leave-to {
  transform: translateY(100%);
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
