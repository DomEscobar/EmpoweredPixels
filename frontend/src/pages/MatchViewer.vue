<template>
  <div class="fixed inset-0 bg-black overflow-hidden select-none font-mono">
    
    <!-- Layer 0: Gameplay Canvas (Die Spielwelt) -->
    <canvas
      ref="canvasRef"
      class="fixed inset-0 w-full h-full object-contain pixelated cursor-crosshair z-0"
      @mousedown="startDrag"
      @mousemove="onDrag"
      @mouseup="endDrag"
      @mouseleave="endDrag"
      @wheel.prevent="onWheel"
    ></canvas>

    <!-- Layer 1: Atmosphäre & Vignette -->
    <div class="fixed inset-0 pointer-events-none z-[1] opacity-60 bg-[radial-gradient(circle_at_center,transparent_0%,rgba(0,0,0,0.9)_100%)]"></div>
    <div class="fixed inset-0 pointer-events-none z-[2] opacity-5 bg-[linear-gradient(transparent_50%,rgba(0,0,0,0.5)_50%)] bg-[length:100%_4px]"></div>

    <!-- Layer 10: HUD (Head-Up Display) -->
    <div class="fixed inset-0 z-10 pointer-events-none flex flex-col justify-between p-4 sm:p-6">
      
      <!-- Top Section -->
      <div class="flex justify-between items-start pointer-events-auto">
        <button @click="$router.push('/matches')" class="rpg-btn-small p-3 bg-amber-900/60 border-amber-700 hover:bg-amber-800 transition-colors">
          <span class="i-hi-arrow-left w-6 h-6 text-amber-500"></span>
        </button>
        
        <div class="flex gap-2">
           <div v-if="match" class="bg-black/80 border-2 border-amber-900 px-3 py-2 text-amber-500/80 text-[10px] hidden md:block">
            MATCH_{{ match.id.slice(0, 8) }}
          </div>
          <button @click="showLog = !showLog" 
                  class="rpg-btn-small p-3 transition-all" 
                  :class="showLog ? 'bg-amber-600 border-amber-400' : 'bg-amber-900/60 border-amber-700 hover:bg-amber-800'">
            <span class="i-hi-clipboard-list w-6 h-6" :class="showLog ? 'text-black' : 'text-amber-500'"></span>
          </button>
          <button @click="showConfig = !showConfig" 
                  class="rpg-btn-small p-3 transition-all"
                  :class="showConfig ? 'bg-amber-600 border-amber-400' : 'bg-amber-900/60 border-amber-700 hover:bg-amber-800'">
            <span class="i-hi-cog w-6 h-6" :class="showConfig ? 'text-black' : 'text-amber-500'"></span>
          </button>
        </div>
      </div>

      <!-- Center Overlay (Round Info - Only visible when not playing) -->
      <div class="flex flex-col items-center justify-center transition-opacity duration-500 pointer-events-none" :class="isPlaying ? 'opacity-10' : 'opacity-100'">
        <span class="text-amber-500/20 font-black text-7xl sm:text-9xl tracking-[0.2em] italic uppercase">ROUND {{ selectedRound }}</span>
      </div>

      <!-- Bottom Section -->
      <div class="flex flex-col gap-4 pointer-events-auto">
        
        <!-- Log Modal (Slide-up, 30svh) -->
        <Transition name="slide-up">
          <div v-if="showLog" class="bg-black/95 border-b-0 border-4 border-amber-900 h-[30svh] overflow-y-auto p-4 pixel-box relative custom-scrollbar shadow-[0_-20px_50px_rgba(0,0,0,0.8)]">
            <div class="sticky top-0 right-0 left-0 bg-black/90 flex justify-between items-center mb-3 border-b border-amber-900/50 pb-2">
              <span class="text-amber-500 font-black text-xs uppercase tracking-widest">>> COMBAT_LOG</span>
              <span class="text-amber-900 text-[10px]">TICK_FEED: READY</span>
            </div>
            <div v-for="(log, idx) in combatLogs" :key="idx" class="mb-2 text-[12px] leading-snug">
              <span class="text-amber-900 mr-2">[{{ idx.toString().padStart(3, '0') }}]</span>
              <span class="text-amber-400/90">{{ log }}</span>
            </div>
          </div>
        </Transition>

        <!-- Config/Settings Overlay -->
        <Transition name="fade">
          <div v-if="showConfig" class="bg-amber-950/90 border-4 border-amber-600 p-6 pixel-box flex flex-wrap gap-8 justify-center items-center shadow-2xl">
             <div class="flex flex-col gap-2 items-center">
               <span class="text-amber-300 text-[10px] uppercase font-black tracking-widest">Playback Speed</span>
               <div class="flex gap-2">
                 <button v-for="s in [1, 2, 4]" :key="s" @click="playbackSpeed = s" 
                         class="w-10 h-10 flex items-center justify-center border-2 border-amber-600 text-sm font-bold transition-all" 
                         :class="playbackSpeed === s ? 'bg-amber-500 text-black shadow-[0_0_15px_#f59e0b]' : 'text-amber-600 hover:bg-amber-900'">
                   {{ s }}x
                 </button>
               </div>
             </div>
             <div class="flex flex-col gap-2 items-center">
               <span class="text-amber-300 text-[10px] uppercase font-black tracking-widest">Camera Scale</span>
               <div class="flex gap-2">
                 <button @click="zoom(-0.2)" class="w-10 h-10 border-2 border-amber-600 text-amber-500 font-bold hover:bg-amber-900 active:scale-90">-</button>
                 <div class="w-16 flex items-center justify-center font-mono text-amber-400 font-bold border-y-2 border-amber-900/50 text-xs">
                    {{ Math.round(camera.zoom * 100) }}%
                 </div>
                 <button @click="zoom(0.2)" class="w-10 h-10 border-2 border-amber-600 text-amber-500 font-bold hover:bg-amber-900 active:scale-90">+</button>
               </div>
             </div>
          </div>
        </Transition>

        <!-- Playback Controller (Minimalist Floating Bar) -->
        <div class="flex justify-center items-end gap-6 pb-2">
          
          <!-- Prev -->
          <button @click="stepRound(-1)" class="w-12 h-12 flex items-center justify-center text-amber-700 hover:text-amber-400 active:scale-75 transition-all">
            <span class="i-hi-rewind w-10 h-10"></span>
          </button>
          
          <!-- MAIN ACTION: PLAY/PAUSE -->
          <button @click="togglePlayback" 
                  class="w-24 h-24 rounded-full bg-amber-600 hover:bg-amber-500 flex items-center justify-center shadow-[0_0_40px_rgba(217,119,6,0.4)] active:scale-95 transition-all mb-2 border-4 border-amber-800">
            <span :class="isPlaying ? 'i-hi-pause' : 'i-hi-play'" class="w-12 h-12 text-black translate-x-0.5"></span>
          </button>

          <!-- Next -->
          <button @click="stepRound(1)" class="w-12 h-12 flex items-center justify-center text-amber-700 hover:text-amber-400 active:scale-75 transition-all">
            <span class="i-hi-fast-forward w-10 h-10"></span>
          </button>
        </div>

        <!-- Progress Timeline (Ultra-Thin) -->
        <div class="h-1.5 bg-black/40 border border-amber-900/30 rounded-full overflow-hidden mb-2 mx-12">
          <div class="h-full bg-amber-600 shadow-[0_0_15px_#d97706] transition-all duration-300" 
               :style="{ width: progressPercent + '%' }"></div>
        </div>

      </div>
    </div>

    <!-- Victory Overlay (Centered) -->
    <Transition name="fade">
      <div v-if="matchStatus === 'completed' && orderedRounds.length && selectedRound === orderedRounds[orderedRounds.length-1]" 
           class="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm pointer-events-none">
        <div class="text-center p-12 border-4 border-amber-500 bg-amber-950/30 pixel-box shadow-[0_0_100px_rgba(217,119,6,0.3)]">
          <h2 class="text-7xl font-black text-amber-500 tracking-[0.3em] uppercase italic text-shadow-retro mb-4">VICTORY</h2>
          <div class="h-1 bg-gradient-to-r from-transparent via-amber-600 to-transparent w-full mb-4"></div>
          <p class="text-amber-200 font-bold uppercase tracking-widest text-sm">Challenge Conquered</p>
        </div>
      </div>
    </Transition>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { matchesApi } from '../features/matches/api'

// Assets
const PIXEL_ASSETS = {
  ICON_TROPHY: 'https://api.iconify.design/pixelarticons:trophy.svg?color=%23f59e0b',
}

const route = useRoute()
const router = useRouter()
const matchId = route.params.id as string

// Refs
const canvasRef = ref<HTMLCanvasElement | null>(null)
const match = ref<any>(null)
const matchStatus = ref('loading')
const rounds = ref<any[]>([])
const selectedRound = ref(0)
const isPlaying = ref(false)
const showLog = ref(false)
const showConfig = ref(false)
const playbackSpeed = ref(1)
const lastTickTime = ref(0)
const camera = ref({ x: 0, y: 0, zoom: 1.5 })
const isDragging = ref(false)
const dragStart = ref({ x: 0, y: 0 })

// HUD UI refs
const progressPercent = computed(() => {
  if (!orderedRounds.value.length) return 0
  return (selectedRound.value / orderedRounds.value[orderedRounds.value.length - 1]) * 100
})

const orderedRounds = computed(() => {
  return rounds.value.map(r => r.round).sort((a, b) => a - b)
})

const currentRoundData = computed(() => {
  return rounds.value.find(r => r.round === selectedRound.value)
})

const combatLogs = computed(() => {
  if (!currentRoundData.value) return []
  return currentRoundData.value.ticks
    .map((t: any) => {
      if (t.type === 'spawn') return `Unit spawned at ${Math.round(t.payload.x)},${Math.round(t.payload.y)}`
      if (t.type === 'attack') return `Unit attacked for ${t.payload.damage} damage`
      if (t.type === 'move') return `Unit moved to ${Math.round(t.payload.x || 0)},${Math.round(t.payload.y || 0)}`
      return t.type
    })
})

// Combat Visualization Logic (Canvas)
let animationId: number

const render = () => {
  if (!canvasRef.value) return
  const ctx = canvasRef.value.getContext('2d')
  if (!ctx) return

  // Auto-resize
  if (canvasRef.value.width !== window.innerWidth * window.devicePixelRatio ||
      canvasRef.value.height !== window.innerHeight * window.devicePixelRatio) {
    canvasRef.value.width = window.innerWidth * window.devicePixelRatio
    canvasRef.value.height = window.innerHeight * window.devicePixelRatio
  }

  const { width, height } = canvasRef.value
  ctx.clearRect(0, 0, width, height)
  
  ctx.save()
  // Apply Camera
  ctx.translate(width / 2 + camera.value.x, height / 2 + camera.value.y)
  ctx.scale(camera.value.zoom * window.devicePixelRatio, camera.value.zoom * window.devicePixelRatio)

  // Draw Grid
  ctx.strokeStyle = '#1e1b1b'
  ctx.lineWidth = 0.5
  for (let i = -500; i <= 500; i += 50) {
    ctx.beginPath(); ctx.moveTo(i, -500); ctx.lineTo(i, 500); ctx.stroke()
    ctx.beginPath(); ctx.moveTo(-500, i); ctx.lineTo(500, i); ctx.stroke()
  }

  // Draw Entities
  if (currentRoundData.value) {
    currentRoundData.value.ticks.forEach((tick: any) => {
       if (tick.payload && typeof tick.payload.x === 'number') {
          const { x, y, hp, fighterId } = tick.payload
          // Entity Body
          ctx.fillStyle = fighterId?.includes('-bot-') ? '#991b1b' : '#d97706'
          ctx.fillRect(x - 5, y - 5, 10, 10)
          
          // HP Bar
          ctx.fillStyle = '#000'
          ctx.fillRect(x - 8, y - 12, 16, 2)
          ctx.fillStyle = '#22c55e'
          ctx.fillRect(x - 8, y - 12, 16 * (hp / 200), 2)
       }
    })
  }

  ctx.restore()
  animationId = requestAnimationFrame(render)
}

// Playback Logic
const stepRound = (delta: number) => {
  const next = selectedRound.value + delta
  if (next >= 0 && next <= orderedRounds.value[orderedRounds.value.length - 1]) {
    selectedRound.value = next
  } else if (next > orderedRounds.value[orderedRounds.value.length - 1]) {
    isPlaying.value = false
  }
}

const togglePlayback = () => {
  isPlaying.value = !isPlaying.value
}

// Tick Loop
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
  camera.value.zoom = Math.max(0.5, Math.min(5, camera.value.zoom + delta))
}
const startDrag = (e: MouseEvent) => { isDragging.value = true; dragStart.value = { x: e.clientX - camera.value.x, y: e.clientY - camera.value.y } }
const onDrag = (e: MouseEvent) => { if (isDragging.value) { camera.value.x = e.clientX - dragStart.value.x; camera.value.y = e.clientY - dragStart.value.y } }
const endDrag = () => isDragging.value = false
const onWheel = (e: WheelEvent) => { zoom(e.deltaY * -0.001) }

// Hooks
onMounted(async () => {
  try {
    const data = await matchesApi.getMatch(matchId)
    match.value = data
    rounds.value = data.roundTicks || []
    matchStatus.value = data.status
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
.pixelated { image-rendering: pixelated; }
.pixel-box { box-shadow: 6px 6px 0px 0px rgba(0,0,0,0.8); }
.text-shadow-retro { text-shadow: 4px 4px 0 #000; }

.rpg-btn-small {
  border-bottom: 4px solid rgba(0,0,0,0.3);
  box-shadow: 2px 2px 0px 0px rgba(0,0,0,0.2);
}
.rpg-btn-small:active {
  border-bottom-width: 0;
  transform: translateY(4px);
  box-shadow: none;
}

.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: #000; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #451a03; border-radius: 10px; }

.slide-up-enter-active, .slide-up-leave-active { transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1); }
.slide-up-enter-from, .slide-up-leave-to { transform: translateY(100%); opacity: 0; }

.fade-enter-active, .fade-leave-active { transition: opacity 0.3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
