<template>
  <div class="min-h-screen bg-slate-950 text-slate-100 selection:bg-amber-500/30 font-mono flex flex-col">
    <NavMenu />
    <main class="flex-1 flex flex-col">
      <router-view v-slot="{ Component }">
        <transition 
          name="pixel-fade" 
          mode="out-in"
        >
          <component :is="Component" class="flex-1" />
        </transition>
      </router-view>
    </main>
  </div>
</template>

<script setup lang="ts">
import NavMenu from '@/widgets/NavMenu/NavMenu.vue';
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Press+Start+2P&display=swap');

:root {
  --font-pixel: 'Press Start 2P', monospace;
}

* {
  scrollbar-width: thin;
  scrollbar-color: #475569 #1e293b;
}

*::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}

*::-webkit-scrollbar-track {
  background: #1e293b;
  border: 2px solid #334155;
}

*::-webkit-scrollbar-thumb {
  background: linear-gradient(to bottom, #64748b, #475569);
  border: 2px solid #334155;
}

*::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(to bottom, #94a3b8, #64748b);
}

*::-webkit-scrollbar-corner {
  background: #1e293b;
}

.pixel-fade-enter-active,
.pixel-fade-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.pixel-fade-enter-from {
  opacity: 0;
  transform: translateY(4px);
}

.pixel-fade-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
