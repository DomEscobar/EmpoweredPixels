<template>
  <div class="h-screen max-h-screen flex flex-col bg-slate-950">
    <!-- Header -->
    <header class="shrink-0 flex items-center justify-between px-3 py-2 bg-slate-900 border-b border-slate-700">
      <button @click="$router.push('/matches')" class="text-amber-400 hover:text-amber-300 font-bold text-sm">
        ← Back
      </button>
      <div class="flex items-center gap-2">
        <span v-if="matchStatus === 'running'" class="flex items-center gap-1">
          <span class="w-2 h-2 bg-red-500 rounded-full animate-pulse"></span>
          <span class="text-xs text-red-400 font-bold uppercase">LIVE</span>
        </span>
        <span class="text-xs text-slate-500 font-mono">#{{ matchId?.substring(0, 4) }}</span>
      </div>
    </header>

    <!-- Canvas - Full Height -->
    <div class="flex-1 relative bg-slate-900 overflow-hidden">
      <canvas ref="canvasRef" class="w-full h-full block"></canvas>
      
      <!-- Loading -->
      <div v-if="isLoading" class="absolute inset-0 flex items-center justify-center bg-slate-950/90">
        <div class="flex flex-col items-center gap-3">
          <div class="h-10 w-10 animate-spin border-3 border-amber-500 border-t-transparent rounded-full"></div>
          <p class="text-amber-200 font-bold uppercase">Loading match...</p>
        </div>
      </div>

      <!-- Victory -->
      <div v-if="isVictory" class="absolute inset-0 flex items-center justify-center bg-black/70 pointer-events-none">
        <div class="border-3 border-amber-500 bg-slate-900 p-8 text-center">
          <p class="text-4xl font-black text-amber-500 uppercase mb-2">Victory!</p>
          <p class="text-sm text-slate-400">Match Complete</p>
        </div>
      </div>

      <!-- Error / No Data -->
      <div v-if="!isLoading && rounds.length === 0" class="absolute inset-0 flex items-center justify-center">
        <div class="text-center p-4">
          <p class="text-slate-400 mb-2">No match data available</p>
          <p class="text-xs text-slate-600">Rounds: {{ rounds.length }}</p>
        </div>
      </div>
    </div>

    <!-- Controls -->
    <div class="shrink-0 bg-slate-900 border-t border-slate-700 p-3">
      <!-- Progress Bar -->
      <div class="relative h-4 bg-slate-800 rounded cursor-pointer mb-3" @click="seekToPercent">
        <div class="absolute inset-y-0 left-0 bg-amber-600 rounded-l transition-all" :style="{ width: progressPercent + '%' }"></div>
        <div class="absolute inset-y-0 w-1 bg-white shadow" :style="{ left: progressPercent + '%' }"></div>
      </div>

      <!-- Main Controls -->
      <div class="flex items-center justify-between">
        <!-- Left: Step Back -->
        <button @click="stepRound(-1)" class="w-12 h-12 flex items-center justify-center bg-slate-800 border-2 border-slate-600 text-slate-300 hover:text-white hover:border-amber-500 text-2xl font-bold rounded">
          ‹
        </button>

        <!-- Center: Play -->
        <button @click="togglePlayback" :disabled="rounds.length <= 1" class="w-16 h-16 flex items-center justify-center bg-amber-600 border-3 border-amber-400 text-slate-900 font-black hover:bg-amber-500 text-2xl rounded disabled:opacity-50 disabled:cursor-not-allowed">
          <span v-if="isPlaying" class="tracking-widest">||</span>
          <span v-else>▶</span>
        </button>

        <!-- Right: Step Forward -->
        <button @click="stepRound(1)" class="w-12 h-12 flex items-center justify-center bg-slate-800 border-2 border-slate-600 text-slate-300 hover:text-white hover:border-amber-500 text-2xl font-bold rounded">
          ›
        </button>
      </div>

      <!-- Round Info -->
      <div class="flex items-center justify-center mt-3 gap-4">
        <span class="text-amber-400 font-mono font-bold text-lg">
          Round {{ currentRoundIndex + 1 }} / {{ rounds.length }}
        </span>
        <button @click="showLog = !showLog" class="px-3 py-1 text-sm text-slate-400 hover:text-amber-400 border border-slate-700 hover:border-amber-500 rounded">
          {{ showLog ? 'Hide Log' : 'Show Log' }}
        </button>
      </div>

      <!-- Settings Toggle -->
      <div class="flex justify-center mt-2">
        <button @click="showConfig = !showConfig" class="text-xs text-slate-500 hover:text-slate-300">
          {{ showConfig ? '▼ Hide Settings' : '▲ Show Settings' }}
        </button>
      </div>

      <!-- Config Panel -->
      <div v-if="showConfig" class="mt-2 pt-2 border-t border-slate-700 flex items-center justify-center gap-8">
        <div class="flex items-center gap-2">
          <span class="text-xs text-slate-500 uppercase font-bold">Speed</span>
          <input type="range" min="0.5" max="4" step="0.5" v-model.number="playbackSpeed" class="w-24 accent-amber-500" />
          <span class="text-sm text-amber-400 font-mono font-bold w-10">{{ playbackSpeed }}x</span>
        </div>
      </div>
    </div>

    <!-- Log Modal -->
    <transition name="slide-up">
      <div v-if="showLog" class="fixed inset-x-0 bottom-0 z-50 h-[50vh] bg-slate-900 border-t-4 border-amber-600 flex flex-col">
        <div class="flex justify-center py-2 cursor-pointer" @click="showLog = false">
          <div class="w-20 h-2 bg-slate-600 rounded-full"></div>
        </div>
        <div class="flex items-center justify-between px-4 py-2 border-b border-slate-700">
          <span class="text-amber-500 font-bold uppercase">Combat Log</span>
          <span class="text-sm text-slate-500 font-mono">{{ rounds.length }} Rounds</span>
        </div>
        <div class="flex-1 overflow-y-auto p-2 space-y-1">
          <div 
            v-for="(round, idx) in rounds" 
            :key="round.round"
            @click="selectRound(idx)"
            class="p-3 cursor-pointer border-2 rounded transition-all"
            :class="idx === currentRoundIndex ? 'border-amber-500 bg-amber-900/30' : 'border-slate-700 hover:border-slate-500'"
          >
            <div class="flex items-center justify-between">
              <span class="font-bold text-amber-100 text-lg">Round {{ round.round }}</span>
              <span class="text-sm text-slate-500">{{ round.ticks?.length || 0 }} events</span>
            </div>
          </div>
        </div>
      </div>
    </transition>

    <!-- Backdrop -->
    <div v-if="showLog" @click="showLog = false" class="fixed inset-0 z-40 bg-black/50"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';
import { request } from '@/shared/api/http';
import { endpoints } from '@/shared/api/endpoints';
import { useAuthStore } from '@/features/auth/store';
import { useRosterStore } from '@/features/roster/store';

interface Tick {
  type: string;
  payload: Record<string, any>;
}

interface Round {
  round: number;
  ticks?: Tick[];
}

const auth = useAuthStore();
const roster = useRosterStore();
const route = useRoute();

const matchId = ref(route.params.id as string);
const matchStatus = ref<string | null>(null);
const rounds = ref<Round[]>([]);
const isLoading = ref(false);
const showLog = ref(false);
const showConfig = ref(false);
const canvasRef = ref<HTMLCanvasElement | null>(null);

const selectedRoundIdx = ref(0);
const isPlaying = ref(false);
const playbackSpeed = ref(1.0);

let animationId: number | null = null;
let lastStepTime = 0;

const TILE_W = 40;
const TILE_H = 20;
const WORLD_SIZE = 16;

const camera = ref({ x: 0, y: 0, zoom: 1.2 });

const currentRoundIndex = computed(() => selectedRoundIdx.value);

const progressPercent = computed(() => {
  if (!rounds.value.length) return 0;
  return ((selectedRoundIdx.value + 1) / rounds.value.length) * 100;
});

const isVictory = computed(() => {
  return matchStatus.value === 'completed' && selectedRoundIdx.value === rounds.value.length - 1;
});

const toIso = (x: number, y: number) => ({
  x: (x - y) * TILE_W * 0.5,
  y: (x + y) * TILE_H * 0.5
});

const toScreen = (x: number, y: number, w: number, h: number) => {
  const iso = toIso(x, y);
  const cx = w / 2;
  const cy = h / 3;
  return {
    x: cx + iso.x * camera.value.zoom + camera.value.x,
    y: cy + iso.y * camera.value.zoom + camera.value.y
  };
};

const stepRound = (dir: number) => {
  const newIdx = selectedRoundIdx.value + dir;
  if (newIdx >= 0 && newIdx < rounds.value.length) {
    selectedRoundIdx.value = newIdx;
  }
  isPlaying.value = false;
};

const selectRound = (idx: number) => {
  selectedRoundIdx.value = idx;
  showLog.value = false;
};

const togglePlayback = () => {
  if (rounds.value.length <= 1) return;
  if (isPlaying.value) {
    isPlaying.value = false;
  } else {
    if (selectedRoundIdx.value >= rounds.value.length - 1) {
      selectedRoundIdx.value = 0;
    }
    isPlaying.value = true;
    lastStepTime = Date.now();
  }
};

const seekToPercent = (e: MouseEvent) => {
  if (!rounds.value.length) return;
  const rect = (e.currentTarget as HTMLElement).getBoundingClientRect();
  const p = Math.max(0, Math.min(1, (e.clientX - rect.left) / rect.width));
  selectedRoundIdx.value = Math.floor(p * (rounds.value.length - 1));
};

interface Entity {
  id: string;
  x: number;
  y: number;
  hp: number;
  maxHp: number;
  alive: boolean;
  isPlayer: boolean;
}

const buildRoundStates = (): Record<string, Entity>[] => {
  const playerIds = new Set(roster.fighters.map(f => f.id));
  const states: Record<string, Entity>[] = [];
  let currentEntities: Record<string, Entity> = {};

  for (const round of rounds.value) {
    const entities = { ...currentEntities };
    
    for (const tick of round.ticks || []) {
      const p = tick.payload || {};
      
      if (tick.type === 'spawn') {
        entities[p.fighterId] = {
          id: p.fighterId,
          x: p.x ?? 0,
          y: p.y ?? 0,
          hp: p.hp ?? 100,
          maxHp: p.hp ?? 100,
          alive: true,
          isPlayer: playerIds.has(p.fighterId),
        };
      } else if (tick.type === 'move' && entities[p.fighterId]) {
        entities[p.fighterId].x = p.toX ?? entities[p.fighterId].x;
        entities[p.fighterId].y = p.toY ?? entities[p.fighterId].y;
      } else if (tick.type === 'attack' && entities[p.targetId]) {
        entities[p.targetId].hp = Math.max(0, entities[p.targetId].hp - (p.damage || 0));
      } else if (tick.type === 'died' && entities[p.fighterId]) {
        entities[p.fighterId].alive = false;
        entities[p.fighterId].hp = 0;
      }
    }
    
    states.push({ ...entities });
    currentEntities = entities;
  }
  
  return states;
};

const draw = () => {
  if (!canvasRef.value) return;
  const canvas = canvasRef.value;
  const ctx = canvas.getContext('2d');
  if (!ctx) return;

  const rect = canvas.getBoundingClientRect();
  if (canvas.width !== rect.width || canvas.height !== rect.height) {
    canvas.width = rect.width;
    canvas.height = rect.height;
  }

  // Background
  ctx.fillStyle = '#0f172a';
  ctx.fillRect(0, 0, canvas.width, canvas.height);

  // Grid
  const z = camera.value.zoom;
  for (let x = 0; x <= WORLD_SIZE; x++) {
    for (let y = 0; y <= WORLD_SIZE; y++) {
      const pos = toScreen(x, y, canvas.width, canvas.height);
      const tw = TILE_W * z;
      const th = TILE_H * z;
      
      ctx.beginPath();
      ctx.moveTo(pos.x, pos.y);
      ctx.lineTo(pos.x + tw/2, pos.y + th/2);
      ctx.lineTo(pos.x, pos.y + th);
      ctx.lineTo(pos.x - tw/2, pos.y + th/2);
      ctx.closePath();
      ctx.fillStyle = '#1e293b';
      ctx.fill();
      ctx.strokeStyle = '#334155';
      ctx.lineWidth = 1;
      ctx.stroke();
    }
  }

  // Entities
  const states = buildRoundStates();
  const state = states[selectedRoundIdx.value];
  if (state) {
    for (const ent of Object.values(state)) {
      const pos = toScreen(ent.x, ent.y, canvas.width, canvas.height);
      
      if (!ent.alive) {
        ctx.fillStyle = '#374151';
        ctx.fillRect(pos.x - 8*z, pos.y, 16*z, 10*z);
        continue;
      }

      // Shadow
      ctx.fillStyle = 'rgba(0,0,0,0.4)';
      ctx.beginPath();
      ctx.ellipse(pos.x, pos.y + 14*z, 12*z, 5*z, 0, 0, Math.PI * 2);
      ctx.fill();

      const bob = Math.sin(Date.now() / 200) * 2 * z;
      const colors = ent.isPlayer 
        ? { head: '#fbbf24', body: '#3b82f6' }
        : { head: '#4ade80', body: '#166534' };

      // Body
      ctx.fillStyle = colors.body;
      ctx.fillRect(pos.x - 10*z, pos.y - 20*z + bob, 20*z, 20*z);
      
      // Head
      ctx.fillStyle = colors.head;
      ctx.fillRect(pos.x - 8*z, pos.y - 32*z + bob, 16*z, 14*z);

      // HP Bar
      const hpPct = Math.max(0, ent.hp / ent.maxHp);
      ctx.fillStyle = '#000';
      ctx.fillRect(pos.x - 12*z, pos.y - 42*z + bob, 24*z, 5*z);
      ctx.fillStyle = hpPct > 0.5 ? '#10b981' : hpPct > 0.25 ? '#f59e0b' : '#ef4444';
      ctx.fillRect(pos.x - 12*z, pos.y - 42*z + bob, 24*z * hpPct, 5*z);
    }
  }

  // Round number overlay
  ctx.fillStyle = 'rgba(0,0,0,0.7)';
  ctx.fillRect(8, 8, 100, 36);
  ctx.fillStyle = '#fbbf24';
  ctx.font = `bold ${20*z}px monospace`;
  ctx.fillText(`R${rounds.value[selectedRoundIdx.value]?.round || 0}`, 16, 34);
};

const loop = () => {
  draw();
  
  if (isPlaying.value && rounds.value.length > 1) {
    const now = Date.now();
    const delay = 1500 / playbackSpeed.value;
    if (now - lastStepTime >= delay) {
      if (selectedRoundIdx.value < rounds.value.length - 1) {
        selectedRoundIdx.value++;
      } else {
        isPlaying.value = false;
      }
      lastStepTime = now;
    }
  }
  
  animationId = requestAnimationFrame(loop);
};

const fetchData = async () => {
  if (!auth.token || !matchId.value) return;
  isLoading.value = true;
  try {
    // Fetch match status
    const match = await request<{ status?: string }>(`${endpoints.match}/${matchId.value}`, { token: auth.token });
    matchStatus.value = match?.status ?? null;
    
    // Fetch round ticks
    const data = await request<any[]>(`${endpoints.match}/${matchId.value}/roundticks`, { token: auth.token });
    
    console.log('[MatchViewer] Raw data:', data);
    
    // Handle different data formats
    if (Array.isArray(data)) {
      // Check if it's already in the right format
      if (data.length > 0 && data[0] && 'round' in data[0]) {
        rounds.value = data;
      } else if (data.length > 0 && Array.isArray(data[0])) {
        // Might be nested arrays - flatten or use as-is
        rounds.value = data.map((item, idx) => ({
          round: idx,
          ticks: Array.isArray(item) ? item : [item]
        }));
      } else {
        // Try to use as single round
        rounds.value = [{ round: 0, ticks: data }];
      }
    } else if (data && typeof data === 'object') {
      const obj = data as Record<string, any>;
      // Handle object with rounds property
      if (Array.isArray(obj.rounds)) {
        rounds.value = obj.rounds;
      } else if (Array.isArray(obj.roundTicks)) {
        rounds.value = obj.roundTicks;
      } else {
        rounds.value = [{ round: 0, ticks: obj.roundTicks ? [obj.roundTicks] : [obj] }];
      }
    }
    
    console.log('[MatchViewer] Parsed rounds:', rounds.value.length, rounds.value);
  } catch (e) {
    console.error('[MatchViewer] Failed to load match:', e);
  } finally {
    isLoading.value = false;
  }
};

onMounted(async () => {
  await roster.fetchFighters();
  await fetchData();
  
  if (rounds.value.length) {
    selectedRoundIdx.value = matchStatus.value === 'completed' ? 0 : rounds.value.length - 1;
  }
  
  animationId = requestAnimationFrame(loop);
});

onUnmounted(() => {
  if (animationId) cancelAnimationFrame(animationId);
});
</script>

<style scoped>
.slide-up-enter-active, .slide-up-leave-active {
  transition: transform 0.3s ease;
}
.slide-up-enter-from, .slide-up-leave-to {
  transform: translateY(100%);
}
</style>
