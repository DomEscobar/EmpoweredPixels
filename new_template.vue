<template>
  <div class="relative w-full h-[100svh] overflow-hidden bg-slate-950 font-mono text-slate-200 select-none">
    
    <!-- Canvas (Background Layer) -->
    <div class="absolute inset-0 z-0 bg-[#050505]">
       <canvas
          ref="canvasRef"
          class="w-full h-full block cursor-grab active:cursor-grabbing pixelated"
          @mousedown="startDrag"
          @mousemove="onDrag"
          @mouseup="endDrag"
          @mouseleave="endDrag"
          @wheel.prevent="onWheel"
       ></canvas>
    </div>

    <!-- Global Vignette -->
    <div class="absolute inset-0 z-0 pointer-events-none bg-[radial-gradient(circle_at_center,transparent_0%,rgba(0,0,0,0.7)_100%)]"></div>

    <!-- Loading Overlay -->
    <div v-if="isLoading" class="absolute inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm">
        <div class="flex flex-col items-center gap-4">
            <div class="w-12 h-12 border-4 border-amber-500 border-t-transparent rounded-full animate-spin"></div>
            <div class="text-amber-500 font-bold tracking-widest animate-pulse text-shadow-retro">LOADING DUNGEON...</div>
        </div>
    </div>

    <!-- Match Status (Top Left) -->
    <div class="absolute top-4 left-4 z-10 flex flex-col gap-1 pointer-events-none">
        <div class="pointer-events-auto flex items-center gap-3">
             <button @click="$router.push('/matches')" class="px-3 py-1 bg-slate-900/80 border border-slate-700 text-slate-300 hover:text-white hover:border-amber-500/50 transition-colors text-xs font-bold uppercase tracking-wider backdrop-blur-md rounded">
                ← Exit
            </button>
            <h1 class="text-amber-500 font-black text-lg leading-none tracking-tight shadow-black drop-shadow-md text-shadow-retro">MATCH #{{ matchId?.substring(0,8) }}</h1>
        </div>
        <div class="flex items-center gap-2">
            <span v-if="matchStatus" class="px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider bg-slate-900/80 border border-slate-700 text-slate-400 rounded">
              {{ matchStatus }}
            </span>
             <div v-if="matchStatus === 'running'" class="flex items-center gap-1.5 px-2 py-0.5 bg-rose-900/30 border border-rose-500/30 rounded">
                <span class="animate-pulse w-1.5 h-1.5 bg-rose-500 rounded-full"></span>
                <span class="text-[10px] text-rose-300 font-bold uppercase">LIVE</span>
             </div>
        </div>
    </div>

    <!-- Round Counter (Top Center) -->
    <div class="absolute top-4 left-1/2 -translate-x-1/2 z-10 pointer-events-none">
        <div class="bg-slate-900/80 border-2 border-slate-700 px-6 py-2 backdrop-blur-md rounded-lg shadow-xl">
            <div class="text-[10px] text-slate-500 uppercase font-bold tracking-widest text-center mb-0.5">Round</div>
            <div class="text-3xl font-black text-amber-500 leading-none text-center text-shadow-retro">
                {{ selectedRound }} <span class="text-slate-600 text-lg">/ {{ orderedRounds.length ? orderedRounds[orderedRounds.length-1] : 0 }}</span>
            </div>
        </div>
    </div>

    <!-- Controls Bar (Bottom) -->
    <div class="absolute bottom-6 left-4 right-4 z-20 flex flex-col gap-3 group transition-opacity hover:opacity-100" :class="{'opacity-0': isPlaying && !showControlsHover}">
        
        <!-- Timeline Bar -->
        <div class="w-full h-3 bg-slate-900/80 border border-slate-700 rounded-full cursor-pointer relative overflow-hidden backdrop-blur-sm shadow-lg group/bar" @click="seekToPercent">
            <!-- Progress Fill -->
            <div class="absolute inset-y-0 left-0 bg-gradient-to-r from-amber-700 to-amber-500 transition-all duration-100" :style="{ width: progressPercent + '%' }"></div>
            <!-- Hover Tick -->
            <div class="absolute top-0 bottom-0 w-0.5 bg-white opacity-0 group-hover/bar:opacity-100 transition-opacity pointer-events-none mix-blend-overlay"></div>
        </div>

        <!-- Buttons Row -->
        <div class="flex items-center justify-between pointer-events-auto">
            
            <div class="flex items-center gap-2 sm:gap-4 bg-slate-900/90 border border-slate-700 p-2 rounded-xl backdrop-blur-md shadow-2xl">
                <!-- Play/Pause -->
                <button @click="togglePlayback" class="w-10 h-10 flex items-center justify-center bg-amber-600 rounded-lg text-slate-900 font-black hover:bg-amber-500 active:scale-95 transition-all border-b-4 border-amber-800 active:border-b-0 active:translate-y-1 shadow-lg">
                    <span v-if="isPlaying" class="text-lg">||</span>
                    <span v-else class="text-lg pl-0.5">▶</span>
                </button>
                
                <div class="w-px h-8 bg-slate-700 mx-1"></div>

                <!-- Step Controls -->
                <div class="flex items-center gap-1">
                    <button @click="stepRound(-1)" class="w-8 h-8 flex items-center justify-center text-slate-400 hover:text-white hover:bg-slate-800 rounded font-bold transition-colors">
                        ⏮
                    </button>
                    <button @click="stepRound(1)" class="w-8 h-8 flex items-center justify-center text-slate-400 hover:text-white hover:bg-slate-800 rounded font-bold transition-colors">
                        ⏭
                    </button>
                </div>

                <!-- Speed -->
                <div class="flex items-center gap-2 bg-slate-950 rounded px-2 py-1 ml-2 border border-slate-800">
                    <span class="text-[10px] font-bold text-slate-500">SPEED</span>
                    <input type="range" min="0.5" max="4" step="0.5" v-model.number="playbackSpeed" class="w-16 accent-amber-500 h-1.5 bg-slate-800 rounded-lg appearance-none cursor-pointer" />
                    <span class="text-[10px] font-mono font-bold text-amber-500 w-6 text-right">{{ playbackSpeed }}x</span>
                </div>
            </div>

            <!-- Log Toggle -->
            <button @click="showLogs = !showLogs" class="h-12 px-4 bg-slate-900/90 border border-slate-700 rounded-xl backdrop-blur-md shadow-2xl flex items-center gap-2 hover:bg-slate-800 transition-colors group/log" :class="{'!border-amber-500/50 !text-amber-400': showLogs}">
                <div class="relative">
                     <span class="text-xl">📜</span>
                     <div v-if="false" class="absolute -top-1 -right-1 w-3 h-3 bg-rose-500 rounded-full animate-bounce"></div> 
                </div>
                <span class="hidden sm:inline text-xs font-bold uppercase tracking-wider group-hover/log:text-white transition-colors">Battle Log</span>
            </button>

        </div>
    </div>

    <!-- Floating Combat Log (Drawer) -->
    <transition name="slide-up">
        <div v-if="showLogs" class="absolute bottom-24 right-4 z-30 w-full max-w-sm max-h-[40vh] bg-slate-900/95 border-2 border-slate-700 rounded-xl shadow-2xl backdrop-blur-xl flex flex-col font-mono text-xs overflow-hidden animate-in slide-in-from-bottom-4 fade-in duration-300">
            <div class="flex items-center justify-between p-3 border-b border-slate-700 bg-slate-950/50">
                <span class="text-amber-500 font-bold uppercase tracking-widest flex items-center gap-2">
                    <span>⚔️</span> Combat Events
                </span>
                <button @click="showLogs = false" class="text-slate-500 hover:text-white px-2">✕</button>
            </div>
            <div class="flex-1 overflow-y-auto p-3 space-y-2 scrollbar-thin scrollbar-thumb-slate-700 scrollbar-track-transparent" ref="logsContainer">
                 <div v-for="(log, i) in combatLogs" :key="i" class="p-2 rounded bg-slate-950/50 border-l-2 border-slate-700 hover:bg-slate-800/50 transition-colors">
                    <div class="opacity-80 break-words">{{ log }}</div>
                 </div>
                 <div v-if="combatLogs.length === 0" class="text-center py-6 text-slate-600 italic">No combat events recorded yet...</div>
            </div>
        </div>
    </transition>

    <!-- Post-Match Victory Overlay -->
    <div v-if="matchStatus === 'completed' && orderedRounds.length && selectedRound === orderedRounds[orderedRounds.length-1]" class="absolute inset-0 z-40 flex items-center justify-center bg-black/70 backdrop-blur-sm animate-in fade-in duration-700 pointer-events-auto">
         <div class="bg-slate-900 border-4 border-amber-600 p-8 rounded-2xl shadow-[0_0_100px_rgba(245,158,11,0.3)] max-w-lg w-full text-center relative overflow-hidden group">
             <!-- Rays Effect -->
             <div class="absolute inset-0 bg-[conic-gradient(from_0deg,transparent_0deg,rgba(245,158,11,0.1)_20deg,transparent_40deg)] animate-[spin_10s_linear_infinite] pointer-events-none"></div>
             
             <div class="relative z-10 flex flex-col gap-4">
                 <img v-if="isWinner" :src="PIXEL_ASSETS.ICON_TROPHY" class="w-24 h-24 mx-auto animate-bounce pixelated drop-shadow-xl" />
                 <img v-else :src="PIXEL_ASSETS.ICON_SKULL" class="w-20 h-20 mx-auto grayscale opacity-50 pixelated" />
                 
                 <h2 class="text-4xl font-black italic tracking-widest text-transparent bg-clip-text bg-gradient-to-b from-amber-300 to-amber-600 text-shadow-lg" :class="isWinner ? 'from-amber-300 to-amber-600' : 'from-slate-400 to-slate-700'">
                    {{ isWinner ? 'VICTORY' : 'DEFEAT' }}
                 </h2>

                 <div class="w-full h-px bg-gradient-to-r from-transparent via-amber-500/50 to-transparent my-2"></div>

                 <!-- Rewards -->
                 <div v-if="rewardsStore.rewardCount > 0" class="flex flex-col gap-3 bg-black/40 p-4 rounded-xl border border-amber-500/30">
                     <span class="text-[10px] uppercase font-bold text-amber-200/70 tracking-widest">Rewards Looted</span>
                     <div class="text-2xl font-black text-amber-400">+ {{ rewardsStore.rewardCount * 20 }} <span class="text-sm text-amber-600">Particles</span></div>
                 </div>

                 <button @click="claimAndExit" class="w-full py-4 mt-2 bg-gradient-to-r from-amber-600 to-amber-500 hover:from-amber-500 hover:to-amber-400 text-slate-900 font-black uppercase tracking-widest rounded-xl shadow-lg hover:shadow-amber-500/20 active:scale-95 transition-all text-sm border-t border-amber-300">
                    {{ rewardsStore.isLoading ? 'Claiming...' : 'Claim Rewards & Exit' }}
                 </button>
             </div>
         </div>
    </div>

  </div>
</template>
