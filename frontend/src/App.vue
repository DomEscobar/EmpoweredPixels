<!--
  EmpoweredPixels - Full-Stack Pixel Art Game
  
  APPLICATION PURPOSE:
  A comprehensive gaming platform featuring fighter management, combat simulation,
  competitive leagues, and an in-game economy. The system provides user authentication,
  roster management, matchmaking, real-time battle replays, seasonal progression,
  leaderboards, daily rewards, and a virtual shop for gold packages and bundles.
  
  FRONTEND ROLE:
  Root Vue 3 component that serves as the single-page application entry point.
  Manages global layout with persistent navigation and dynamic route rendering.
  Implements the "Ethereal Iron" theme with dark slate palette and gold accents.
  
  ARCHITECTURE:
  - Frontend: Vue 3 (Composition API), Vite, Tailwind CSS v4, Pinia (state)
  - Backend: Go (Clean Architecture), PostgreSQL, WebSocket real-time hub
  - UI Style: Pixel-art aesthetic with high-readability dark mode interface
  
  COMPONENT RESPONSIBILITIES:
  - Render NavMenu (persistent navigation)
  - Display router-view with transition effects (pixel-fade)
  - Provide footer spacing for mobile/desktop layouts
  - Set global typography (monospace font, custom scrollbars)
-->
<template>
  <div class="min-h-screen bg-slate-950 text-slate-100 selection:bg-amber-500/30 font-mono flex flex-col">
    <NavMenu />
    <main class="flex-1 flex flex-col pb-16 md:pb-20">
      <router-view v-slot="{ Component }">
        <transition 
          name="pixel-fade" 
          mode="out-in"
        >
          <component
            :is="Component"
            class="flex-1"
          />
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
