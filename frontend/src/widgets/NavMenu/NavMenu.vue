<template>
  <!-- Mobile Menu Button (Hamburger) -->
  <button 
    v-if="auth.token"
    class="fixed top-4 left-4 z-50 md:hidden pixel-box-xs bg-slate-900/90 border-slate-700 p-2"
    @click="showDrawer = true"
    aria-label="Open menu"
  >
    <svg v-if="!showDrawer" xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-slate-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="square" stroke-linejoin="miter">
      <line x1="3" y1="6" x2="21" y2="6"></line>
      <line x1="3" y1="12" x2="21" y2="12"></line>
      <line x1="3" y1="18" x2="21" y2="18"></line>
    </svg>
    <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-slate-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="square" stroke-linejoin="miter">
      <line x1="18" y1="6" x2="6" y2="18"></line>
      <line x1="6" y1="6" x2="18" y2="18"></line>
    </svg>
  </button>

  <!-- Mobile Drawer Overlay -->
  <div 
    v-if="showDrawer" 
    class="fixed inset-0 bg-black/60 z-50 md:hidden"
    @click="showDrawer = false"
  ></div>

  <!-- Mobile Navigation Drawer -->
  <aside 
    v-if="auth.token"
    class="fixed top-0 left-0 h-full w-72 bg-slate-900/98 border-r-4 border-slate-700 z-50 transform transition-transform duration-300 md:hidden flex flex-col"
    :class="showDrawer ? 'translate-x-0' : '-translate-x-full'"
    :style="{ backgroundImage: `url('${PIXEL_ASSETS.BG_NAV}')`, backgroundSize: '64px 64px', imageRendering: 'pixelated' }"
  >
    <!-- Drawer Header -->
    <div class="p-4 border-b-2 border-slate-800 bg-slate-950/90">
      <div class="flex items-center justify-between">
        <div class="pixel-logo-box flex items-center gap-2 px-3 py-2 bg-slate-800/80 border border-slate-600">
          <img :src="PIXEL_ASSETS.ICON_LOGO" alt="EP" class="w-6 h-6 pixelated" />
          <span class="text-sm font-bold text-amber-400 tracking-wider">EMPOWERED<span class="text-white">PIXELS</span></span>
        </div>
        <button @click="showDrawer = false" class="text-slate-400 hover:text-white text-xl">
          ✕
        </button>
      </div>
    </div>

    <!-- Drawer Navigation Items -->
    <nav class="flex-1 overflow-y-auto p-4 custom-scrollbar">
      <div class="space-y-2">
        <router-link 
          v-for="item in allNavItems" 
          :key="item.path" 
          :to="item.path"
          class="flex items-center gap-3 px-4 py-3 pixel-box-sm transition-all"
          :class="isActive(item.path) ? 'bg-amber-600/20 border-amber-500/50' : 'bg-slate-800/60 border-slate-700/50 hover:bg-slate-700/60'"
          @click="showDrawer = false"
        >
          <img :src="item.icon" alt="" class="w-6 h-6 pixelated" />
          <span class="font-bold text-slate-200">{{ item.name }}</span>
        </router-link>
      </div>
    </nav>

    <!-- Drawer Footer -->
    <div class="p-4 border-t-2 border-slate-800 bg-slate-950/90 space-y-3">
      <!-- Daily Reward Button -->
      <button 
        class="w-full pixel-box-sm flex items-center justify-center gap-2 px-4 py-3 transition-all relative"
        :class="dailyStore.canClaim ? 'bg-amber-600/20 border-amber-500/50 animate-pulse' : 'bg-slate-800/60 border-slate-700/50'"
        title="Daily Reward"
        @click="showDailyModal = true"
      >
        <span class="text-xl">{{ dailyStore.nextReward?.icon || '🎁' }}</span>
        <span class="font-bold">Daily Reward</span>
        <span v-if="dailyStore.canClaim" class="absolute top-1 right-2 w-2.5 h-2.5 bg-red-500 rounded-full animate-ping"></span>
      </button>

      <!-- Logout Button -->
      <button class="w-full pixel-box-sm bg-slate-800/60 border-slate-700/50 p-3 hover:bg-red-900/30 transition-colors flex items-center justify-center gap-2" title="Sign Out" @click="logout">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-400 hover:text-red-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="square" stroke-linejoin="miter">
          <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
          <polyline points="16 17 21 12 16 7"></polyline>
          <line x1="21" y1="12" x2="9" y2="12"></line>
        </svg>
        <span class="font-bold text-slate-300">Sign Out</span>
      </button>
    </div>
  </aside>

  <!-- Desktop Bottom Navigation -->
  <nav v-if="auth.token" class="fixed bottom-0 left-0 right-0 z-40 hidden md:block">
    <!-- Main Bottom Navigation Bar -->
    <div 
      class="pixel-nav-bottom bg-slate-900/95 border-t-4 border-slate-700"
      :style="{ 
        backgroundImage: `url('${PIXEL_ASSETS.BG_NAV}')`,
        backgroundSize: '64px 64px',
        imageRendering: 'pixelated'
      }"
    >
      <div class="mx-auto px-1">
        <!-- Primary Navigation Items -->
        <div class="grid grid-cols-5 gap-0.5 h-16">
          <router-link 
            v-for="item in allNavItems" 
            :key="item.path" 
            :to="item.path"
            class="bottom-nav-item"
            :class="{ 'bottom-nav-item-active': isActive(item.path) }"
          >
            <img :src="item.icon" alt="" class="w-5 h-5 pixelated mb-1 mx-auto" />
            <span class="text-[8px] font-bold uppercase tracking-wider truncate w-full text-center leading-tight">{{ item.name }}</span>
          </router-link>
        </div>
      </div>

      <!-- User Controls Bar (Secondary) -->
      <div class="pixel-nav-bottom-secondary bg-slate-950/90 border-t-2 border-slate-800">
        <div class="mx-auto flex items-center justify-between h-10 px-2 max-w-7xl">
          <!-- Daily Reward -->
          <button 
            class="pixel-box-xs flex items-center gap-2 px-2 py-1 transition-all relative"
            :class="dailyStore.canClaim ? 'bg-amber-600/20 border-amber-500/50 animate-pulse' : 'bg-slate-800/80 border-slate-600/50'"
            title="Daily Reward"
            @click="showDailyModal = true"
          >
            <span class="text-sm">{{ dailyStore.nextReward?.icon || '🎁' }}</span>
            <span class="text-[10px] font-bold">Daily</span>
            <span v-if="dailyStore.canClaim" class="absolute -top-1 -right-1 w-2.5 h-2.5 bg-red-500 rounded-full animate-ping"></span>
          </button>

          <!-- Logout -->
          <button class="pixel-box-xs bg-slate-800/80 border-slate-600/50 p-1.5 hover:bg-red-900/30 transition-colors" title="Sign Out" data-testid="logout-btn" @click="logout">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400 hover:text-red-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="square" stroke-linejoin="miter">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
              <polyline points="16 17 21 12 16 7"></polyline>
              <line x1="21" y1="12" x2="9" y2="12"></line>
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Logo / Brand (Optional floating element) -->
    <div class="fixed left-1/2 -translate-x-1/2 -top-8 z-50 hidden md:block">
      <router-link to="/dashboard" class="pixel-logo-float flex items-center gap-2 px-3 py-1.5 bg-slate-900/90 border-2 border-slate-700 shadow-lg">
        <div class="pixel-logo-box bg-amber-600 p-1">
          <img :src="PIXEL_ASSETS.ICON_LOGO" alt="EP" class="w-4 h-4 pixelated" />
        </div>
        <span class="text-xs font-bold text-amber-400 tracking-wider">EMPOWERED<span class="text-white">PIXELS</span></span>
      </router-link>
    </div>
  </nav>

  <!-- Daily Reward Modal -->
  <DailyRewardModal :show="showDailyModal" @close="showDailyModal = false" />
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '@/features/auth/store';

import { useDailyStore } from '@/features/daily/store';
import DailyRewardModal from '@/features/daily/components/DailyRewardModal.vue';

const PIXEL_ASSETS = {
  BG_NAV: 'https://vibemedia.space/bg_nav_wood_1x2y3z_v1.png?prompt=dark%20wood%20plank%20texture%20seamless%20horizontal&style=pixel_game_asset&key=NOGON',
  ICON_LOGO: 'https://vibemedia.space/icon_ep_logo_4a5b6c_v1.png?prompt=pixel%20art%20sword%20and%20shield%20emblem%20icon&style=pixel_game_asset&key=NOGON',
  ICON_USER: 'https://vibemedia.space/icon_user_helm_7d8e9f_v1.png?prompt=knight%20helmet%20pixel%20art%20icon%20small&style=pixel_game_asset&key=NOGON',
  ICON_DASHBOARD: 'https://vibemedia.space/icon_nav_castle_0g1h2i_v1.png?prompt=castle%20tower%20pixel%20art%20icon%20small&style=pixel_game_asset&key=NOGON',
  ICON_ROSTER: 'https://vibemedia.space/icon_nav_knights_3j4k5l_v1.png?prompt=knights%20group%20pixel%20art%20icon%20small&style=pixel_game_asset&key=NOGON',
  ICON_MATCHES: 'https://vibemedia.space/icon_nav_swords_6m7n8o_v1.png?prompt=crossed%20swords%20pixel%20art%20icon%20small&style=pixel_game_asset&key=NOGON',
  ICON_INVENTORY: 'https://vibemedia.space/icon_nav_chest_9p0q1r_v1.png?prompt=treasure%20chest%20pixel%20art%20icon%20small&style=pixel_game_asset&key=NOGON',
  ICON_LEAGUES: 'https://vibemedia.space/icon_nav_trophy_2s3t4u_v1.png?prompt=golden%20trophy%20pixel%20art%20icon%20small&style=pixel_game_asset&key=NOGON',
  ICON_GOLD: 'https://vibemedia.space/icon_gold_coin_nav_8f7e6d.png?prompt=golden%20coin%20with%20shine%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
};

const auth = useAuthStore();
const dailyStore = useDailyStore();
const router = useRouter();
const route = useRoute();
const showDailyModal = ref(false);
const showDrawer = ref(false);

// All navigation items (5 primary)
const allNavItems = [
  { name: 'Command', path: '/dashboard', icon: PIXEL_ASSETS.ICON_DASHBOARD },
  { name: 'Roster', path: '/roster', icon: PIXEL_ASSETS.ICON_ROSTER },
  { name: 'Battle', path: '/matches', icon: PIXEL_ASSETS.ICON_MATCHES },
  { name: 'Vault', path: '/inventory', icon: PIXEL_ASSETS.ICON_INVENTORY },
  { name: 'Leagues', path: '/leagues', icon: PIXEL_ASSETS.ICON_LEAGUES },
];

onMounted(() => {
  if (auth.token) {
    dailyStore.fetchStatus();
  }
});

watch(() => auth.token, (newToken) => {
  if (newToken) {
    dailyStore.fetchStatus();
  }
});

const isActive = (path: string) => {
  return route.path === path || route.path.startsWith(path + '/');
};

const logout = () => {
  auth.logout();
  router.push('/login');
};
</script>

<style scoped>
.pixelated {
  image-rendering: pixelated;
}

.pixel-nav-bottom {
  box-shadow: 
    0 -4px 0 #0f172a,
    0 -6px 8px rgba(0, 0, 0, 0.4);
}

.pixel-nav-bottom-secondary {
  box-shadow: 0 -2px 0 #0f172a;
}

.pixel-logo-float {
  border-radius: 0.5rem;
}

.pixel-logo-box {
  border: 2px solid #78350f;
  box-shadow: 2px 2px 0 #451a03;
  border-radius: 0.25rem;
}

.pixel-box-xs {
  border: 2px solid #334155;
  box-shadow: 2px 2px 0 #0f172a;
  border-radius: 0.375rem;
}

/* Bottom Navigation Items */
.bottom-nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.25rem;
  height: 100%;
  min-height: 64px;
  color: #94a3b8;
  border: 2px solid transparent;
  transition: all 0.15s;
  border-radius: 0.375rem;
  margin: 0.5px;
}

.bottom-nav-item:hover {
  color: #fef3c7;
  background: rgba(30, 41, 59, 0.3);
}

.bottom-nav-item-active {
  color: #fbbf24;
  background: rgba(146, 64, 14, 0.15);
  border-color: #92400e;
}

/* Safe area for mobile bottom nav */
@supports (padding-bottom: env(safe-area-inset-bottom)) {
  .bottom-nav-item {
    padding-bottom: calc(env(safe-area-inset-bottom) / 4);
  }
}

/* Custom scrollbar for drawer */
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}

.custom-scrollbar::-webkit-scrollbar-track {
  background: rgba(15, 23, 42, 0.5);
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(71, 85, 105, 0.8);
  border-radius: 3px;
}

.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: rgba(100, 116, 139, 0.9);
}

/* Prevent body scroll when drawer open */
body.drawer-open {
  overflow: hidden;
}
</style>
