<template>
  <div class="h-svh max-h-svh overflow-hidden flex flex-col bg-slate-950">
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

    <!-- Canvas Area -->
    <div class="flex-1 relative bg-slate-900">
      <canvas ref="canvasRef" class="w-full h-full"></canvas>
      
      <!-- Loading -->
      <div v-if="isLoading" class="absolute inset-0 flex items-center justify-center bg-slate-950/80">
        <div class="flex flex-col items-center gap-2">
          <div class="h-8 w-8 animate-spin border-2 border-amber-500 border-t-transparent rounded-full"></div>
          <p class="text-amber-200 text-sm font-bold uppercase">Loading...</p>
        </div>
      </div>

      <!-- Victory -->
      <div v-if="isVictory" class="absolute inset-0 flex items-center justify-center bg-black/60">
        <div class="border-2 border-amber-500 bg-slate-900 p-6 text-center">
          <p class="text-3xl font-black text-amber-500 uppercase mb-1">Victory</p>
          <p class="text-xs text-slate-400">Match Complete</p>
        </div>
      </div>
    </div>

    <!-- Controls -->
    <div class="shrink-0 bg-slate-900 border-t border-slate-700 p-2">
      <!-- Progress -->
      <div class="relative h-3 bg-slate-800 rounded cursor-pointer mb-2" @click="seekToPercent">
        <div class="absolute inset-y-0 left-0 bg-amber-600 rounded-l" :style="{ width: progressPercent + '%' }"></div>
      </div>

      <!-- Main Controls -->
      <div class="flex items-center justify-between">
        <!-- Left: Steps -->
        <div class="flex items-center gap-1">
          <button @click="stepRound(-1)" class="w-10 h-10 flex items-center justify-center bg-slate-800 border border-slate-600 text-slate-300 hover:text-white hover:border-amber-500 text-lg font-bold">
            ‹
          </button>
        </div>

        <!-- Center: Play -->
        <button @click="togglePlayback" class="w-14 h-14 flex items-center justify-center bg-amber-600 border-2 border-amber-400 text-slate-900 font-black hover:bg-amber-500 text-xl rounded">
          <span v-if="isPlaying">||</span>
          <span v-else>▶</span>
        </button>

        <!-- Right: Steps -->
        <div class="flex items-center gap-1">
          <button @click="stepRound(1)" class="w-10 h-10 flex items-center justify-center bg-slate-800 border border-slate-600 text-slate-300 hover:text-white hover:border-amber-500 text-lg font-bold">
            ›
          </button>
        </div>
      </div>

      <!-- Round Info -->
      <div class="flex items-center justify-center mt-2 gap-4">
        <span class="text-amber-400 font-mono font-bold text-sm">
          Round {{ currentRoundIndex + 1 }} / {{ rounds.length }}
        </span>
        <button @click="showLog = !showLog" class="text-xs text-slate-400 hover:text-amber-400">
          {{ showLog ? 'Hide Log' : 'Show Log' }}
        </button>
      </div>

      <!-- Config (expandable) -->
      <div v-if="showConfig" class="mt-2 pt-2 border-t border-slate-700 flex items-center justify-center gap-6">
        <div class="flex items-center gap-2">
          <span class="text-xs text-slate-500 uppercase font-bold">Speed</span>
          <input type="range" min="0.5" max="4" step="0.5" v-model.number="playbackSpeed" class="w-20 accent-amber-500" />
          <span class="text-xs text-amber-400 font-mono font-bold w-8">{{ playbackSpeed }}x</span>
        </div>
        <button @click="showConfig = false" class="text-xs text-slate-500">✕</button>
      </div>
      <div v-else class="flex justify-center mt-1">
        <button @click="showConfig = true" class="text-xs text-slate-600 hover:text-slate-400">⚙ Settings</button>
      </div>
    </div>

    <!-- Log Modal -->
    <transition name="slide-up">
      <div v-if="showLog" class="fixed inset-x-0 bottom-0 z-50 h-[40vh] bg-slate-900 border-t-4 border-amber-600 flex flex-col">
        <div class="flex justify-center py-2 cursor-pointer" @click="showLog = false">
          <div class="w-16 h-1.5 bg-slate-600 rounded-full"></div>
        </div>
        <div class="flex items-center justify-between px-3 py-1 border-b border-slate-700">
          <span class="text-amber-500 font-bold uppercase">Combat Log</span>
          <span class="text-xs text-slate-500 font-mono">{{ rounds.length }} Rounds</span>
        </div>
        <div class="flex-1 overflow-y-auto p-2 space-y-1">
          <div 
            v-for="(round, idx) in rounds" 
            :key="round.round"
            @click="selectRound(idx)"
            class="p-2 cursor-pointer border rounded"
            :class="idx === currentRoundIndex ? 'border-amber-500 bg-amber-900/20' : 'border-slate-700 hover:border-slate-500'"
          >
            <div class="flex items-center justify-between">
              <span class="font-bold text-amber-100">Round {{ round.round }}</span>
              <span class="text-xs text-slate-500">{{ round.ticks?.length || 0 }} events</span>
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
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { request } from '@/shared/api/http';
import { endpoints } from '@/shared/api/endpoints';
import { useAuthStore } from '@/features/auth/store';
import { useRosterStore } from '@/features/roster/store';

const auth = useAuthStore();
const roster = useRosterStore();
const route = useRoute();

const matchId = ref(route.params.id as string);
const matchStatus = ref<string | null>(null);
const rounds = ref<any[]>([]);
const isLoading = ref(false);
const showLog = ref(false);
const showConfig = ref(false);
const canvasRef = ref<HTMLCanvasElement | null>(null);

const selectedRoundIdx = ref(0);
const isPlaying = ref(false);
const playbackSpeed = ref(1.0);

let animationId: number | null = null;
let lastStepTime = 0;

const TILE_W = 48;
const TILE_H = 24;
const WORLD_SIZE = 20;

const camera = ref({ x: 0, y: 0, zoom: 1 });

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

const buildRoundStates = () => {
  const playerIds = new Set(roster.fighters.map(f => f.id));
  const states: any[] = [];
  let currentEntities: Record<string, any> = {};

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
    
    states.push({ entities: { ...entities } });
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
    for (const ent of Object.values(state.entities) as any[]) {
      const pos = toScreen(ent.x, ent.y, canvas.width, canvas.height);
      
      if (!ent.alive) {
        ctx.fillStyle = '#374151';
        ctx.fillRect(pos.x - 6*z, pos.y, 12*z, 8*z);
        continue;
      }

      // Shadow
      ctx.fillStyle = 'rgba(0,0,0,0.4)';
      ctx.beginPath();
      ctx.ellipse(pos.x, pos.y + 12*z, 10*z, 4*z, 0, 0, Math.PI * 2);
      ctx.fill();

      const bob = Math.sin(Date.now() / 200) * 2 * z;
      const colors = ent.isPlayer 
        ? { head: '#fbbf24', body: '#3b82f6' }
        : { head: '#4ade80', body: '#166534' };

      // Body
      ctx.fillStyle = colors.body;
      ctx.fillRect(pos.x - 8*z, pos.y - 18*z + bob, 16*z, 18*z);
      
      // Head
      ctx.fillStyle = colors.head;
      ctx.fillRect(pos.x - 6*z, pos.y - 28*z + bob, 12*z, 12*z);

      // HP Bar
      const hpPct = Math.max(0, ent.hp / ent.maxHp);
      ctx.fillStyle = '#000';
      ctx.fillRect(pos.x - 10*z, pos.y - 36*z + bob, 20*z, 4*z);
      ctx.fillStyle = hpPct > 0.5 ? '#10b981' : hpPct > 0.25 ? '#f59e0b' : '#ef4444';
      ctx.fillRect(pos.x - 10*z, pos.y - 36*z + bob, 20*z * hpPct, 4*z);
    }
  }

  // Round number overlay
  ctx.fillStyle = 'rgba(0,0,0,0.6)';
  ctx.fillRect(10, 10, 80, 28);
  ctx.fillStyle = '#fbbf24';
  ctx.font = `bold ${16*z}px monospace`;
  ctx.fillText(`R${rounds.value[selectedRoundIdx.value]?.round || 0}`, 20, 30);
};

const loop = () => {
  draw();
  
  if (isPlaying.value && rounds.value.length > 1) {
    const now = Date.now();
    const delay = 1200 / playbackSpeed.value;
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
    const [match, ticksData] = await Promise.all([
      request<{ status?: string }>(`${endpoints.match}/${matchId.value}`, { token: auth.token }),
      request<any[]>(`${endpoints.match}/${matchId.value}/roundticks`, { token: auth.token })
    ]);
    matchStatus.value = match?.status ?? null;
    rounds.value = ticksData || [];
  } catch (e) {
    console.error('Failed to load match:', e);
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
