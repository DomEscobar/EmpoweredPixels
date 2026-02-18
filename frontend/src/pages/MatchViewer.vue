<template>
  <div class="fixed inset-0 bg-black overflow-hidden select-none font-mono">
    
    <!-- Layer 0: Full Canvas -->
    <canvas
      ref="canvasRef"
      class="fixed inset-0 w-full h-full object-contain pixelated cursor-move z-0"
      @mousedown="startDrag"
      @mousemove="onDrag"
      @mouseup="endDrag"
      @mouseleave="endDrag"
      @wheel.prevent="onWheel"
    ></canvas>

    <!-- Layer 1: Vignette -->
    <div class="fixed inset-0 pointer-events-none z-[1] opacity-50 bg-[radial-gradient(circle_at_center,transparent_30%,rgba(0,0,0,0.95)_100%)]"></div>

    <!-- Layer 2: Top Bar (Minimal) -->
    <div class="fixed top-0 left-0 right-0 z-20 p-3 flex justify-between items-center pointer-events-none">
      <div class="flex items-center gap-2 pointer-events-auto">
        <button @click="$router.push('/matches')" class="w-8 h-8 flex items-center justify-center bg-black/60 border border-amber-900/50 hover:bg-amber-900/60 text-amber-500 text-xs transition-all">
          ←
        </button>
        <div class="bg-black/60 border border-amber-900/50 px-2 py-1 text-[10px] text-amber-400/80 font-bold tracking-wider">
          {{ match?.id?.slice(0, 8) || 'LOADING' }}
        </div>
      </div>
      
      <div class="flex items-center gap-1 pointer-events-auto">
        <!-- Round indicator -->
        <div class="bg-black/60 border border-amber-900/50 px-3 py-1 flex items-center gap-2">
          <span class="text-amber-600 text-[10px]">RND</span>
          <span class="text-amber-400 font-black text-sm">{{ selectedRound }}<span class="text-amber-600 text-xs">/{{ orderedRounds[orderedRounds.length-1] || 0 }}</span></span>
        </div>
        
        <!-- Status badge -->
        <div v-if="matchStatus === 'completed'" class="bg-emerald-950/80 border border-emerald-700 px-2 py-1">
          <span class="text-emerald-400 text-[10px] font-bold">DONE</span>
        </div>
        <div v-else-if="isPlaying" class="bg-amber-950/80 border border-amber-700 px-2 py-1 flex items-center gap-1">
          <span class="w-1.5 h-1.5 bg-amber-500 rounded-full animate-pulse"></span>
          <span class="text-amber-400 text-[10px] font-bold">LIVE</span>
        </div>
      </div>
    </div>

    <!-- Center Round Display (fades when playing) -->
    <div v-if="rounds.length" class="fixed inset-0 z-10 flex items-center justify-center pointer-events-none">
      <div class="text-center transition-all duration-500" :class="isPlaying ? 'opacity-5 scale-90' : 'opacity-20 scale-100'">
        <span class="text-amber-500/30 font-black text-[12vw] leading-none tracking-[0.1em]">{{ selectedRound }}</span>
      </div>
    </div>

    <!-- No Data / Loading Overlay -->
    <div v-if="!rounds.length && matchStatus !== 'loading'" 
         class="fixed inset-0 z-20 flex items-center justify-center bg-black/80">
      <div class="text-center p-8">
        <div class="text-amber-500 text-4xl mb-4">⚔</div>
        <h2 class="text-amber-400 font-black text-xl uppercase tracking-wider mb-2">No Replay Data</h2>
        <p class="text-amber-600/70 text-sm mb-4">{{ matchStatus === 'lobby' ? 'Match has not started yet' : 'Battle replay unavailable' }}</p>
        <button @click="$router.push('/matches')" class="px-4 py-2 bg-amber-700 hover:bg-amber-600 text-black font-bold text-sm rounded">
          Back to Matches
        </button>
      </div>
    </div>

    <!-- Layer 4: Bottom Dock (Compacts to icons, expands on interaction) -->
    <div v-if="rounds.length" class="fixed bottom-0 left-0 right-0 z-30">
      
      <!-- Progress Bar (Always visible thin line) -->
      <div class="h-1 bg-black/80 cursor-pointer" @click="seekToPercent">
        <div class="h-full bg-gradient-to-r from-amber-700 to-amber-500 transition-all duration-100" 
             :style="{ width: progressPercent + '%' }"></div>
      </div>
      
      <!-- Main Dock Container -->
      <div class="bg-gradient-to-t from-black via-black/95 to-transparent pb-2 pt-3 px-4">
        
        <!-- Expanded Panel (slides up when needed) -->
        <Transition name="dock-expand">
          <div v-if="expandedDock" class="mb-3 bg-amber-950/20 border border-amber-900/30 rounded-lg p-3">
            <!-- Log Panel -->
            <div v-if="expandedDock === 'log'" class="h-32 overflow-y-auto custom-scrollbar">
              <div v-for="(log, idx) in combatLogs" :key="idx" class="text-[10px] text-amber-400/70 mb-1">
                <span class="text-amber-900">[{{ idx.toString().padStart(3, '0') }}]</span> {{ log }}
              </div>
              <div v-if="!combatLogs.length" class="text-amber-600/50 text-[10px] italic">No events</div>
            </div>
            
            <!-- Settings Panel -->
            <div v-else-if="expandedDock === 'settings'" class="flex gap-6 justify-center">
              <div class="flex flex-col items-center gap-1">
                <span class="text-[9px] text-amber-600 uppercase tracking-wider">Speed</span>
                <div class="flex gap-1">
                  <button v-for="s in [1, 2, 4]" :key="s" @click="playbackSpeed = s"
                          class="w-7 h-7 text-[10px] border transition-all"
                          :class="playbackSpeed === s ? 'bg-amber-600 border-amber-400 text-black font-bold' : 'border-amber-900/50 text-amber-500 hover:bg-amber-900/30'">
                    {{ s }}x
                  </button>
                </div>
              </div>
              <div class="flex flex-col items-center gap-1">
                <span class="text-[9px] text-amber-600 uppercase tracking-wider">Zoom</span>
                <div class="flex items-center gap-1">
                  <button @click="zoom(-0.2)" class="w-6 h-6 border border-amber-900/50 text-amber-500 hover:bg-amber-900/30 text-xs">-</button>
                  <span class="w-10 text-center text-[10px] text-amber-400 font-mono">{{ Math.round(camera.zoom * 100) }}</span>
                  <button @click="zoom(0.2)" class="w-6 h-6 border border-amber-900/50 text-amber-500 hover:bg-amber-900/30 text-xs">+</button>
                </div>
              </div>
            </div>
          </div>
        </Transition>

        <!-- Dock Controls (Icon Bar) -->
        <div class="flex items-end justify-center gap-1">
          
          <!-- Left: Secondary Actions -->
          <div class="flex gap-1 mr-4">
            <button @click="toggleDock('log')" 
                    class="dock-btn" 
                    :class="expandedDock === 'log' ? 'bg-amber-600 text-black' : 'text-amber-500'">
              <span class="text-xs">☰</span>
            </button>
            <button @click="toggleDock('settings')" 
                    class="dock-btn" 
                    :class="expandedDock === 'settings' ? 'bg-amber-600 text-black' : 'text-amber-500'">
              <span class="text-xs">⚙</span>
            </button>
          </div>

          <!-- Center: Playback -->
          <div class="flex items-center gap-2">
            <button @click="stepRound(-1)" class="dock-btn-secondary text-amber-600">
              ‹
            </button>
            
            <button @click="togglePlayback" 
                    class="w-12 h-12 rounded-full bg-amber-600 hover:bg-amber-500 flex items-center justify-center shadow-lg hover:shadow-amber-600/30 transition-all active:scale-95 border-2 border-amber-400">
              <span v-if="isPlaying" class="text-black text-xs font-bold">❚❚</span>
              <span v-else class="text-black text-xs font-bold ml-0.5">▶</span>
            </button>
            
            <button @click="stepRound(1)" class="dock-btn-secondary text-amber-600">
              ›
            </button>
          </div>

          <!-- Right: Info -->
          <div class="flex gap-1 ml-4">
            <div class="h-8 px-2 flex items-center bg-black/40 border border-amber-900/30 text-[9px] text-amber-600">
              {{ currentRoundData?.ticks?.length || 0 }} EV
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Victory Overlay -->
    <Transition name="fade">
      <div v-if="matchStatus === 'completed' && orderedRounds.length && selectedRound === orderedRounds[orderedRounds.length-1]" 
           class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 pointer-events-none">
        <div class="text-center p-8 border-4 border-amber-500 bg-black/90">
          <div class="text-amber-500 text-5xl mb-2">🏆</div>
          <h2 class="text-amber-500 font-black text-2xl tracking-[0.2em] uppercase">Victory</h2>
        </div>
      </div>
    </Transition>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getMatch, getMatchRoundTicks } from '../features/matches/api'
import { useAuthStore } from '@/features/auth/store'

const route = useRoute()
const router = useRouter()
const matchId = route.params.id as string
const auth = useAuthStore()

// Refs
const canvasRef = ref<HTMLCanvasElement | null>(null)
const match = ref<any>(null)
const matchStatus = ref('loading')
const rounds = ref<any[]>([])
const selectedRound = ref(0)
const isPlaying = ref(false)
const playbackSpeed = ref(1)
const lastTickTime = ref(0)
const camera = ref({ x: 0, y: 0, zoom: 1.5 })
const isDragging = ref(false)
const dragStart = ref({ x: 0, y: 0 })
const expandedDock = ref<'log' | 'settings' | null>(null)

// Computed
const orderedRounds = computed(() => {
  return rounds.value.map(r => r.round).sort((a, b) => a - b)
})

const progressPercent = computed(() => {
  if (!orderedRounds.value.length) return 0
  const max = orderedRounds.value[orderedRounds.value.length - 1]
  return max ? (selectedRound.value / max) * 100 : 0
})

const currentRoundData = computed(() => {
  return rounds.value.find(r => r.round === selectedRound.value)
})

const combatLogs = computed(() => {
  if (!currentRoundData.value) return []
  return currentRoundData.value.ticks?.slice(0, 50).map((t: any) => {
    if (t.type === 'spawn') return `Spawned at ${Math.round(t.payload?.x || 0)},${Math.round(t.payload?.y || 0)}`
    if (t.type === 'attack') return `${t.payload?.attackerId?.slice(0,4)} hits ${t.payload?.targetId?.slice(0,4)} for ${t.payload?.damage}`
    if (t.type === 'move') return `Moved to ${Math.round(t.payload?.x || 0)},${Math.round(t.payload?.y || 0)}`
    if (t.type === 'died') return `${t.payload?.fighterId?.slice(0,4)} died`
    return t.type
  }) || []
})

// Render
let animationId: number

const render = () => {
  if (!canvasRef.value) return
  const ctx = canvasRef.value.getContext('2d')
  if (!ctx) return

  const dpr = window.devicePixelRatio || 1
  if (canvasRef.value.width !== window.innerWidth * dpr || canvasRef.value.height !== window.innerHeight * dpr) {
    canvasRef.value.width = window.innerWidth * dpr
    canvasRef.value.height = window.innerHeight * dpr
  }

  const { width, height } = canvasRef.value
  ctx.clearRect(0, 0, width, height)
  
  ctx.save()
  ctx.translate(width / 2 + camera.value.x, height / 2 + camera.value.y)
  ctx.scale(camera.value.zoom * dpr, camera.value.zoom * dpr)

  // Grid
  ctx.strokeStyle = '#1a1515'
  ctx.lineWidth = 0.5
  for (let i = -500; i <= 500; i += 50) {
    ctx.beginPath(); ctx.moveTo(i, -500); ctx.lineTo(i, 500); ctx.stroke()
    ctx.beginPath(); ctx.moveTo(-500, i); ctx.lineTo(500, i); ctx.stroke()
  }

  // Entities
  if (currentRoundData.value?.ticks) {
    const entities = new Map()
    currentRoundData.value.ticks.forEach((tick: any) => {
      if (tick.type === 'spawn' && tick.payload) {
        entities.set(tick.payload.fighterId, { ...tick.payload, alive: true })
      }
      if (tick.type === 'died' && tick.payload?.fighterId) {
        const e = entities.get(tick.payload.fighterId)
        if (e) e.alive = false
      }
      if ((tick.type === 'attack' || tick.type === 'move') && tick.payload?.fighterId) {
        const e = entities.get(tick.payload.fighterId)
        if (e) Object.assign(e, tick.payload)
      }
    })

    entities.forEach((e: any) => {
      const isBot = e.fighterId?.includes('-bot-') || e.fighterId?.startsWith('bot-')
      const color = isBot ? '#dc2626' : '#f59e0b'
      
      // Body
      ctx.fillStyle = e.alive ? color : '#374151'
      ctx.fillRect(e.x - 6, e.y - 6, 12, 12)
      
      // HP
      const hpPct = Math.max(0, (e.hp || 100) / 200)
      ctx.fillStyle = '#000'
      ctx.fillRect(e.x - 10, e.y - 14, 20, 3)
      ctx.fillStyle = hpPct > 0.5 ? '#22c55e' : hpPct > 0.25 ? '#eab308' : '#ef4444'
      ctx.fillRect(e.x - 10, e.y - 14, 20 * hpPct, 3)
    })
  }

  ctx.restore()
  animationId = requestAnimationFrame(render)
}

// Playback
const stepRound = (delta: number) => {
  const max = orderedRounds.value[orderedRounds.value.length - 1] || 0
  const next = selectedRound.value + delta
  if (next >= 0 && next <= max) {
    selectedRound.value = next
  } else if (next > max) {
    isPlaying.value = false
  }
}

const togglePlayback = () => {
  isPlaying.value = !isPlaying.value
}

const seekToPercent = (e: MouseEvent) => {
  const rect = (e.target as HTMLElement).getBoundingClientRect()
  const pct = (e.clientX - rect.left) / rect.width
  const max = orderedRounds.value[orderedRounds.value.length - 1] || 0
  selectedRound.value = Math.round(max * pct)
}

const tick = (time: number) => {
  if (isPlaying.value) {
    const delta = time - lastTickTime.value
    const speedMs = 1000 / (playbackSpeed.value || 1)
    if (delta > speedMs) {
      stepRound(1)
      lastTickTime.value = time
    }
  }
  requestAnimationFrame(tick)
}

// Interaction
const zoom = (delta: number) => {
  camera.value.zoom = Math.max(0.3, Math.min(5, camera.value.zoom + delta))
}
const startDrag = (e: MouseEvent) => { isDragging.value = true; dragStart.value = { x: e.clientX - camera.value.x, y: e.clientY - camera.value.y } }
const onDrag = (e: MouseEvent) => { if (isDragging.value) { camera.value.x = e.clientX - dragStart.value.x; camera.value.y = e.clientY - dragStart.value.y } }
const endDrag = () => isDragging.value = false
const onWheel = (e: WheelEvent) => { zoom(e.deltaY * -0.001) }

const toggleDock = (panel: 'log' | 'settings') => {
  expandedDock.value = expandedDock.value === panel ? null : panel
}

// Init
onMounted(async () => {
  try {
    const token = auth.token || ''
    const [matchData, roundTicksData] = await Promise.all([
      getMatch(token, matchId),
      getMatchRoundTicks(token, matchId).catch(() => [])
    ])
    match.value = matchData
    rounds.value = roundTicksData || []
    matchStatus.value = matchData.status
    if (orderedRounds.value.length) selectedRound.value = orderedRounds.value[0]
    render()
    requestAnimationFrame(tick)
  } catch (err) {
    console.error('Failed to load match:', err)
  }
})

onUnmounted(() => {
  cancelAnimationFrame(animationId)
})
</script>

<style scoped>
.dock-btn {
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(0, 0, 0, 0.6);
  border: 1px solid rgba(146, 64, 14, 0.5);
  border-radius: 0.25rem;
  transition: all 0.2s;
  font-size: 0.75rem;
}
.dock-btn:hover {
  background-color: rgba(146, 64, 14, 0.6);
}
.dock-btn-secondary {
  width: 1.75rem;
  height: 1.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 0.25rem;
  transition: all 0.2s;
  font-size: 1.125rem;
  font-weight: bold;
}
.dock-btn-secondary:hover {
  background-color: rgba(146, 64, 14, 0.3);
}
.custom-scrollbar::-webkit-scrollbar { width: 3px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #451a03; border-radius: 3px; }

.dock-expand-enter-active, .dock-expand-leave-active { transition: all 0.25s ease; }
.dock-expand-enter-from, .dock-expand-leave-to { transform: translateY(10px); opacity: 0; }
.fade-enter-active, .fade-leave-active { transition: opacity 0.4s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
