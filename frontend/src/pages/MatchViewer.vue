<template>
  <div class="fixed inset-0 bg-black overflow-hidden select-none font-mono">
    
    <!-- Layer 0: Gameplay Canvas (3D Isometric Board) -->
    <canvas
      ref="canvasRef"
      class="fixed inset-0 w-full h-full object-contain pixelated cursor-crosshair z-0"
      @mousedown="startDrag"
      @mousemove="onDrag"
      @mouseup="endDrag"
      @mouseleave="endDrag"
      @wheel.prevent="onWheel"
    ></canvas>

    <!-- Layer 1: Atmosphere & Scanlines -->
    <div class="fixed inset-0 pointer-events-none z-[1] opacity-50 bg-[radial-gradient(circle_at_center,transparent_0%,rgba(0,0,0,0.9)_100%)]"></div>
    <div class="fixed inset-0 pointer-events-none z-[2] opacity-5 bg-[linear-gradient(transparent_50%,rgba(0,0,0,0.5)_50%)] bg-[length:100%_4px]"></div>

    <!-- Layer 10: HUD -->
    <div class="fixed inset-0 z-10 pointer-events-none flex flex-col justify-between p-4 sm:p-6">
      
      <!-- Top: Back & Info -->
      <div class="flex justify-between items-start pointer-events-auto">
        <button @click="$router.push('/matches')" class="rpg-btn-small p-3 bg-amber-900/40 border-amber-600 hover:bg-amber-800 transition-colors">
          <span class="i-hi-arrow-left w-6 h-6 text-amber-500"></span>
        </button>
        
        <div class="flex flex-col items-end gap-2">
          <div v-if="match" class="bg-black/80 border-2 border-amber-900 px-3 py-1 text-amber-500/60 text-[10px] uppercase font-black">
            MATCH: {{ match.id.slice(0, 8) }}
          </div>
          <div class="flex gap-2">
            <button @click="expandedDock = expandedDock === 'log' ? null : 'log'" 
                    class="rpg-btn-small p-3 transition-all" 
                    :class="expandedDock === 'log' ? 'bg-amber-600 border-amber-400' : 'bg-amber-900/40 border-amber-600'">
              <span class="i-hi-clipboard-list w-6 h-6" :class="expandedDock === 'log' ? 'text-black' : 'text-amber-500'"></span>
            </button>
            <button @click="expandedDock = expandedDock === 'settings' ? null : 'settings'" 
                    class="rpg-btn-small p-3 transition-all"
                    :class="expandedDock === 'settings' ? 'bg-amber-600 border-amber-400' : 'bg-amber-900/40 border-amber-600'">
              <span class="i-hi-cog w-6 h-6" :class="expandedDock === 'settings' ? 'text-black' : 'text-amber-500'"></span>
            </button>
          </div>
        </div>
      </div>

      <!-- Center: Round Large Counter -->
      <div class="flex flex-col items-center justify-center transition-opacity duration-700 pointer-events-none" :class="isPlaying ? 'opacity-10' : 'opacity-100'">
        <span class="text-amber-500/10 font-black text-8xl sm:text-[12rem] tracking-[0.1em] italic uppercase select-none">RND {{ selectedRound }}</span>
      </div>

      <!-- Bottom: Controls & Log -->
      <div class="flex flex-col gap-4 pointer-events-auto max-w-2xl mx-auto w-full">
        
        <!-- Log Slide-up -->
        <Transition name="slide-up">
          <div v-if="expandedDock === 'log'" class="bg-black/95 border-b-0 border-4 border-amber-900 h-[25svh] overflow-y-auto p-4 pixel-box relative custom-scrollbar">
            <div class="sticky top-0 bg-black/90 flex justify-between items-center mb-3 border-b border-amber-900/50 pb-2">
              <span class="text-amber-500 font-bold text-[10px] uppercase tracking-widest">>> COMBAT_LOG</span>
              <span class="text-amber-900 text-[9px] uppercase">Tick Process Active</span>
            </div>
            <div v-for="(log, idx) in combatLogs" :key="idx" class="mb-1 text-[11px] leading-tight font-mono">
              <span class="text-amber-900 mr-2">[{{ idx.toString().padStart(3, '0') }}]</span>
              <span class="text-amber-400/80">{{ log }}</span>
            </div>
            <div v-if="!combatLogs.length" class="text-amber-900 text-[10px] text-center mt-4 uppercase">No actions in this round</div>
          </div>
        </Transition>

        <!-- Settings Slide-up -->
        <Transition name="slide-up">
          <div v-if="expandedDock === 'settings'" class="bg-amber-950/90 border-4 border-amber-600 p-6 pixel-box flex flex-wrap gap-8 justify-center items-center">
             <div class="flex flex-col gap-2 items-center">
               <span class="text-amber-300 text-[9px] uppercase font-black tracking-widest">Time Dilation</span>
               <div class="flex gap-1">
                 <button v-for="s in [1, 2, 4]" :key="s" @click="playbackSpeed = s" 
                         class="w-10 h-10 flex items-center justify-center border-2 border-amber-600 text-xs font-bold transition-all" 
                         :class="playbackSpeed === s ? 'bg-amber-500 text-black shadow-[0_0_15px_#f59e0b]' : 'text-amber-600 hover:bg-amber-900'">
                   {{ s }}x
                 </button>
               </div>
             </div>
             <div class="flex flex-col gap-2 items-center">
               <span class="text-amber-300 text-[9px] uppercase font-black tracking-widest">Camera Scale</span>
               <div class="flex gap-1">
                 <button @click="zoom(-0.2)" class="w-10 h-10 border-2 border-amber-600 text-amber-500 font-bold hover:bg-amber-900 active:scale-90">-</button>
                 <button @click="zoom(0.2)" class="w-10 h-10 border-2 border-amber-600 text-amber-500 font-bold hover:bg-amber-900 active:scale-90">+</button>
               </div>
             </div>
          </div>
        </Transition>

        <!-- Main Controller Row -->
        <div class="flex justify-center items-end gap-6">
          <button @click="stepRound(-1)" class="w-10 h-10 flex items-center justify-center text-amber-700 hover:text-amber-400 active:scale-75 transition-all">
            <span class="i-hi-rewind w-8 h-8"></span>
          </button>
          
          <button @click="togglePlayback" 
                  class="w-20 h-20 rounded-full bg-amber-600 hover:bg-amber-500 flex items-center justify-center shadow-[0_0_30px_rgba(217,119,6,0.5)] active:scale-95 transition-all border-4 border-amber-800">
            <span :class="isPlaying ? 'i-hi-pause' : 'i-hi-play'" class="w-10 h-10 text-black"></span>
          </button>

          <button @click="stepRound(1)" class="w-10 h-10 flex items-center justify-center text-amber-700 hover:text-amber-400 active:scale-75 transition-all">
            <span class="i-hi-fast-forward w-8 h-8"></span>
          </button>
        </div>

        <!-- Progress Timeline -->
        <div class="h-1 bg-amber-900/20 rounded-full overflow-hidden mb-2">
          <div class="h-full bg-amber-600 shadow-[0_0_10px_#d97706] transition-all duration-300" 
               :style="{ width: progressPercent + '%' }"></div>
        </div>

      </div>
    </div>

    <!-- Victory Overlay -->
    <Transition name="fade">
      <div v-if="matchStatus === 'completed' && orderedRounds.length && selectedRound === orderedRounds[orderedRounds.length-1]" 
           class="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-md pointer-events-none">
        <div class="text-center p-12 border-4 border-amber-500 bg-amber-950/40 pixel-box shadow-[0_0_100px_rgba(217,119,6,0.3)]">
          <h2 class="text-6xl font-black text-amber-500 tracking-[0.2em] uppercase italic text-shadow-retro mb-4">VICTORY</h2>
          <p class="text-amber-200 font-bold uppercase tracking-widest text-xs">Quest Status: Complete</p>
        </div>
      </div>
    </Transition>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getMatch } from '@/features/matches/api'
import { useAuthStore } from '@/features/auth/store'

const route = useRoute()
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
  if (!currentRoundData.value?.ticks) return []
  return currentRoundData.value.ticks.slice(0, 50).map((t: any) => {
    if (t.type === 'spawn') return `Unit deployed at ${Math.round(t.payload?.x || 0)},${Math.round(t.payload?.y || 0)}`
    if (t.type === 'attack') return `Unit hits for ${t.payload?.damage} dmg`
    if (t.type === 'move') return `Movement detected`
    if (t.type === 'died') return `Unit neutralized`
    return t.type
  })
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
  // Global 3D Isometry Tilt
  ctx.translate(width / 2, height / 2)
  ctx.scale(1, 0.6) // Perspective squeeze
  ctx.rotate(-Math.PI / 4) // Diamond tilt
  
  // Camera
  ctx.translate(camera.value.x, camera.value.y)
  ctx.scale(camera.value.zoom, camera.value.zoom)

  // Isometric Grid
  ctx.strokeStyle = '#221a1a'
  ctx.lineWidth = 1
  const gridSize = 1000
  for (let i = -gridSize; i <= gridSize; i += 50) {
    ctx.beginPath(); ctx.moveTo(i, -gridSize); ctx.lineTo(i, gridSize); ctx.stroke()
    ctx.beginPath(); ctx.moveTo(-gridSize, i); ctx.lineTo(gridSize, i); ctx.stroke()
  }

  // Draw Entities
  if (currentRoundData.value?.ticks) {
    const states = new Map()
    currentRoundData.value.ticks.forEach((tick: any) => {
       if (tick.payload?.fighterId) {
         states.set(tick.payload.fighterId, { ...tick.payload, alive: tick.type !== 'died' })
       }
    })

    states.forEach((e: any) => {
      const isBot = e.fighterId?.includes('-bot-') || e.fighterId?.startsWith('bot-')
      
      // Shadow (on the floor)
      ctx.fillStyle = 'rgba(0,0,0,0.4)'
      ctx.beginPath()
      ctx.ellipse(e.x, e.y, 8, 4, 0, 0, Math.PI * 2)
      ctx.fill()

      // Floating Body Upgrade
      ctx.save()
      ctx.translate(0, -12) // Hover height
      
      const bloom = isBot ? '#ef4444' : '#f59e0b'
      ctx.fillStyle = e.alive ? bloom : '#374151'
      ctx.shadowBlur = e.alive ? 15 : 0
      ctx.shadowColor = bloom
      ctx.fillRect(e.x - 7, e.y - 7, 14, 14)
      
      // HP Bar (Slightly tilted back)
      ctx.fillStyle = '#000'
      ctx.fillRect(e.x - 10, e.y - 18, 20, 3)
      ctx.fillStyle = '#22c55e'
      ctx.fillRect(e.x - 10, e.y - 18, 20 * (Math.max(0, e.hp || 100) / 200), 3)
      ctx.restore()
    })
  }

  ctx.restore()
  animationId = requestAnimationFrame(render)
}

// Playback Logic
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

const tickLoop = (time: number) => {
  if (isPlaying.value) {
    const delta = time - lastTickTime.value
    const speedMs = 1000 / (playbackSpeed.value || 1)
    if (delta > speedMs) {
      stepRound(1)
      lastTickTime.value = time
    }
  }
  requestAnimationFrame(tickLoop)
}

// Interaction
const zoom = (delta: number) => camera.value.zoom = Math.max(0.4, Math.min(4, camera.value.zoom + delta))
const startDrag = (e: MouseEvent) => { isDragging.value = true; dragStart.value = { x: e.clientX - camera.value.x, y: e.clientY - camera.value.y } }
const onDrag = (e: MouseEvent) => { if (isDragging.value) { camera.value.x = e.clientX - dragStart.value.x; camera.value.y = e.clientY - dragStart.value.y } }
const endDrag = () => isDragging.value = false
const onWheel = (e: WheelEvent) => zoom(e.deltaY * -0.001)

onMounted(async () => {
  try {
    if (!auth.token) return;
    const data = await getMatch(auth.token, matchId)
    match.value = data
    rounds.value = data.roundTicks || []
    matchStatus.value = data.status
    render()
    requestAnimationFrame(tickLoop)
  } catch (err) {
    console.error('Match failed to load:', err)
  }
})

onUnmounted(() => cancelAnimationFrame(animationId))
</script>

<style scoped>
.pixelated { image-rendering: pixelated; }
.pixel-box { box-shadow: 8px 8px 0px 0px rgba(0,0,0,0.8); }
.text-shadow-retro { text-shadow: 4px 4px 0 #000; }

.rpg-btn-small {
  border-bottom: 4px solid rgba(0,0,0,0.3);
  box-shadow: 3px 3px 0px 0px rgba(0,0,0,0.2);
}
.rpg-btn-small:active {
  border-bottom-width: 0;
  transform: translateY(4px);
  box-shadow: none;
}

.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #451a03; border-radius: 4px; }

.slide-up-enter-active, .slide-up-leave-active { transition: all 0.5s cubic-bezier(0.16, 1, 0.3, 1); }
.slide-up-enter-from, .slide-up-leave-to { transform: translateY(100%); opacity: 0; }

.fade-enter-active, .fade-leave-active { transition: opacity 0.5s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
