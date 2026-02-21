<template>
  <nav class="sticky top-0 z-60 w-full font-mono">
    <!-- Main Nav Bar -->
    <div 
      class="pixel-nav bg-slate-900/95 border-b-4 border-slate-700"
      :style="{ 
        backgroundImage: `url('${PIXEL_ASSETS.BG_NAV}')`,
        backgroundSize: '64px 64px',
        imageRendering: 'pixelated'
      }"
    >
      <div class="mx-auto flex h-14 max-w-7xl items-center px-3">
        <!-- Logo -->
        <router-link to="/" class="flex items-center gap-2 group">
          <div class="pixel-logo-box bg-amber-600 p-1.5 group-hover:bg-amber-500 transition-colors">
            <img :src="PIXEL_ASSETS.ICON_LOGO" alt="EP" class="w-5 h-5 pixelated" />
          </div>
          <span class="text-sm font-bold text-amber-400 text-shadow-retro tracking-wider hidden sm:block">
            EMPOWERED<span class="text-white">PIXELS</span>
          </span>
        </router-link>

        <!-- Desktop Nav Links -->
        <div v-if="auth.token" class="hidden md:flex items-center gap-1 ml-4">
          <router-link 
            v-for="item in allNavItems" 
            :key="item.path" 
            :to="item.path"
            class="nav-link"
            :class="{ 'nav-link-active': isActive(item.path) }"
          >
            <img :src="item.icon" alt="" class="w-4 h-4 pixelated" />
            <span>{{ item.name }}</span>
          </router-link>
        </div>

        <!-- Right Side - Mobile Quick Stats -->
        <div class="flex items-center gap-2 ml-auto">
          <template v-if="auth.token">
            <!-- Daily Reward Button - visible on mobile -->
            <button 
              @click="showDailyModal = true"
              class="pixel-box-sm px-2 py-1 flex items-center gap-1 transition-all relative"
              :class="dailyStore.canClaim ? 'bg-amber-600 border-amber-400 animate-pulse' : 'bg-slate-800/80 border-slate-600'"
              title="Daily Reward"
            >
              <span class="text-sm">{{ dailyStore.nextReward?.icon || '🎁' }}</span>
              <span v-if="dailyStore.canClaim" class="absolute -top-1 -right-1 w-3 h-3 bg-red-500 rounded-full animate-ping"></span>
            </button>

            <!-- Gold Display - Compact - visible on mobile -->
            <router-link to="/shop" class="pixel-box-sm bg-amber-900/30 border-amber-600/50 px-2 py-1 flex items-center gap-1.5 hover:bg-amber-900/50 transition-colors min-h-[44px]">
              <img :src="PIXEL_ASSETS.ICON_GOLD" alt="Gold" class="w-4 h-4 pixelated" />
              <span class="text-xs font-bold text-amber-400">{{ formattedGold() }}</span>
            </router-link>

            <!-- Menu Toggle (Mobile Only) -->
            <button 
              @click="mobileMenuOpen = !mobileMenuOpen"
              class="pixel-box-sm md:hidden flex items-center justify-center w-10 h-10 transition-all"
              :class="mobileMenuOpen ? 'bg-amber-600 border-amber-400' : 'bg-slate-800/80 border-slate-600'"
              aria-label="Toggle menu"
            >
              <span class="text-lg">{{ mobileMenuOpen ? '✕' : '☰' }}</span>
            </button>

            <!-- Desktop User Badge & Logout -->
            <div class="hidden md:flex items-center gap-2">
              <div class="pixel-box-sm bg-slate-800/80 px-2 py-1 flex items-center gap-1.5">
                <div class="w-5 h-5 pixel-box-sm bg-indigo-900/50 border-indigo-500/50 flex items-center justify-center">
                  <img :src="PIXEL_ASSETS.ICON_USER" alt="" class="w-3.5 h-3.5 pixelated" />
                </div>
                <span class="text-xs text-slate-300 font-bold uppercase">Commander</span>
              </div>
              <button @click="logout" class="pixel-box-sm bg-slate-800/80 p-1.5 hover:bg-red-900/50 transition-colors" title="Logout" data-testid="logout-btn">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400 hover:text-red-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="square" stroke-linejoin="miter">
                  <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
                  <polyline points="16 17 21 12 16 7"></polyline>
                  <line x1="21" y1="12" x2="9" y2="12"></line>
                </svg>
              </button>
            </div>
          </template>
          <template v-else>
            <router-link to="/login" class="md:hidden">
              <button class="rpg-btn-small text-xs py-1.5 px-2">SIGN IN</button>
            </router-link>
            <router-link to="/register" class="hidden sm:block">
              <button class="rpg-btn text-xs py-1.5 px-3">GET STARTED</button>
            </router-link>
          </template>
        </div>
      </div>
    </div>

    <!-- Mobile Full-Screen Menu -->
    <Transition
      enter-active-class="transition-opacity duration-300"
      leave-active-class="transition-opacity duration-200"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div 
        v-if="mobileMenuOpen && auth.token" 
        class="fixed inset-0 z-50 bg-slate-900/98 md:hidden"
      >
        <div class="flex flex-col h-full">
          <!-- Header -->
          <div class="pixel-box bg-slate-800/50 p-4 flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="pixel-box-sm bg-amber-600/20 border-amber-500/50 p-2">
                <img :src="PIXEL_ASSETS.ICON_LOGO" alt="" class="w-6 h-6 pixelated" />
              </div>
              <div>
                <h3 class="text-sm font-bold text-amber-400">MENU</h3>
                <p class="text-xs text-slate-400">More options</p>
              </div>
            </div>
            <button @click="mobileMenuOpen = false" class="pixel-box-sm w-10 h-10 flex items-center justify-center bg-slate-800/80 hover:bg-red-900/50 transition-colors">
              <span class="text-lg">✕</span>
            </button>
          </div>

          <!-- Nav Items - Secondary only (not in bottom bar) -->
          <div class="flex-1 overflow-y-auto p-4 custom-scrollbar">
            <div v-if="overlayItems.length > 0" class="grid grid-cols-2 gap-3 mb-4">
              <router-link 
                v-for="item in overlayItems" 
                :key="item.path" 
                :to="item.path"
                @click="mobileMenuOpen = false"
                class="mobile-nav-card"
                :class="{ 'mobile-nav-card-active': isActive(item.path) }"
              >
                <img :src="item.icon" alt="" class="w-8 h-8 pixelated mb-2 mx-auto" />
                <span class="text-sm font-bold text-center">{{ item.name }}</span>
              </router-link>
            </div>

            <!-- Empty state if all items in bottom nav -->
            <div v-else class="text-center py-8">
              <p class="text-slate-500 text-sm">All features available in bottom bar</p>
            </div>

            <!-- Logout Button -->
            <button 
              @click="logout(); mobileMenuOpen = false"
              class="mobile-nav-link-secondary mt-4 w-full text-red-400 border-red-500/30 hover:bg-red-900/20"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="square" stroke-linejoin="miter">
                <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
                <polyline points="16 17 21 12 16 7"></polyline>
                <line x1="21" y1="12" x2="9" y2="12"></line>
              </svg>
              <span>Sign Out</span>
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Bottom Navigation Bar (Mobile Only) - 5 Items -->
    <div v-if="auth.token" class="md:hidden fixed bottom-0 left-0 right-0 z-45 pixel-nav border-t-4 border-slate-700 bg-slate-900/95">
      <div class="grid grid-cols-5 h-16">
        <router-link 
          v-for="item in bottomNavItems" 
          :key="item.path" 
          :to="item.path"
          class="bottom-nav-item"
          :class="{ 'bottom-nav-item-active': isActive(item.path) }"
        >
          <img :src="item.icon" alt="" class="w-6 h-6 pixelated mb-1 mx-auto" />
          <span class="text-[10px] font-bold uppercase tracking-wider truncate w-full text-center">{{ item.name }}</span>
        </router-link>
      </div>
    </div>

    <!-- Daily Reward Modal -->
    <DailyRewardModal :show="showDailyModal" @close="showDailyModal = false" />
  </nav>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '@/features/auth/store';
import { useShopStore } from '@/features/shop/store';
import { useDailyStore } from '@/features/daily/store';
import DailyRewardModal from '@/features/daily/components/DailyRewardModal.vue';

const PIXEL_ASSETS = {
  BG_NAV: 'https://vibemedia.space/bg_nav_wood_1x2y3z_v1.png?prompt=dark%20wood%20plank%20texture%20seamless%20horizontal&style=pixel_game_asset&key=NOGON',
  ICON_LOGO: 'https://vibemedia.space/icon_ep_logo_4a5b6c_v1.png?prompt=pixel%20art%20sword%20and%20shield%20emblem%20icon&style=pixel_game_asset&key=NOGON',
  ICON_USER: 'https://vibemedia.space/icon_user_helm_7d8e9f_v1.png?prompt=knight%20helmet%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  ICON_DASHBOARD: 'https://vibemedia.space/icon_nav_castle_0g1h2i_v1.png?prompt=castle%20tower%20pixel%20art%20icon%20small&style=pixel_game_asset&key=NOGON',
  ICON_ROSTER: 'https://vibemedia.space/icon_nav_knights_3j4k5l_v1.png?prompt=knights%20group%20pixel%20art%20icon%20small&style=pixel_game_asset&key=NOGON',
  ICON_MATCHES: 'https://vibemedia.space/icon_nav_swords_6m7n8o_v1.png?prompt=crossed%20swords%20pixel%20art%20icon%20small&style=pixel_game_asset&key=NOGON',
  ICON_INVENTORY: 'https://vibemedia.space/icon_nav_chest_9p0q1r_v1.png?prompt=treasure%20chest%20pixel%20art%20icon%20small&style=pixel_game_asset&key=NOGON',
  ICON_LEAGUES: 'https://vibemedia.space/icon_nav_trophy_2s3t4u_v1.png?prompt=golden%20trophy%20pixel%20art%20icon%20small&style=pixel_game_asset&key=NOGON',
  ICON_SHOP: 'https://vibemedia.space/icon_nav_shop_5v6w7x_v1.png?prompt=gold%20coins%20pile%20pixel%20art%20icon%20small&style=pixel_game_asset&key=NOGON',
  ICON_ATTUNEMENT: 'https://vibemedia.space/icon_nav_crystal_1a2b3c.png?prompt=magic%20crystal%20glowing%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  ICON_LEADERBOARD: 'https://vibemedia.space/icon_nav_ranking_4d5e6f.png?prompt=golden%20crown%20ranking%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  ICON_GOLD: 'https://vibemedia.space/icon_gold_coin_nav_8f7e6d.png?prompt=golden%20coin%20with%20shine%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
};

const auth = useAuthStore();
const shop = useShopStore();
const dailyStore = useDailyStore();
const router = useRouter();
const route = useRoute();
const mobileMenuOpen = ref(false);
const showDailyModal = ref(false);

// All navigation items
const allNavItems = [
  { name: 'Command', path: '/dashboard', icon: PIXEL_ASSETS.ICON_DASHBOARD },
  { name: 'Roster', path: '/roster', icon: PIXEL_ASSETS.ICON_ROSTER },
  { name: 'Battle', path: '/matches', icon: PIXEL_ASSETS.ICON_MATCHES },
  { name: 'Vault', path: '/inventory', icon: PIXEL_ASSETS.ICON_INVENTORY },
  { name: 'Leagues', path: '/leagues', icon: PIXEL_ASSETS.ICON_LEAGUES },
  { name: 'Attune', path: '/attunement', icon: PIXEL_ASSETS.ICON_ATTUNEMENT },
  { name: 'Rankings', path: '/leaderboard', icon: PIXEL_ASSETS.ICON_LEADERBOARD },
  { name: 'Shop', path: '/shop', icon: PIXEL_ASSETS.ICON_SHOP },
];

// Items shown in bottom nav (5 primary items)
const bottomNavItems = computed(() => allNavItems.slice(0, 5));

// Items shown in overlay menu (remaining secondary items)
const overlayItems = computed(() => allNavItems.slice(5));

// Fetch gold balance and daily rewards when logged in
onMounted(() => {
  if (auth.token) {
    shop.fetchGoldBalance();
    dailyStore.fetchStatus();
  }
});

watch(() => auth.token, (newToken) => {
  if (newToken) {
    shop.fetchGoldBalance();
    dailyStore.fetchStatus();
  }
});

const formattedGold = () => {
  const balance = shop.goldBalance?.balance ?? 0;
  if (balance >= 10000) return `${(balance / 1000).toFixed(1)}K`;
  return balance.toLocaleString();
};

const isActive = (path: string) => {
  return route.path === path || route.path.startsWith(path + '/');
};

const logout = () => {
  auth.logout();
  mobileMenuOpen.value = false;
  router.push('/login');
};
</script>

<style scoped>
.pixelated {
  image-rendering: pixelated;
}

.text-shadow-retro {
  text-shadow: 1px 1px 0 #000, 2px 2px 0 rgba(0, 0, 0, 0.3);
}

.pixel-nav {
  box-shadow: 0 4px 0 #0f172a, 0 6px 8px rgba(0, 0, 0, 0.4);
}

.pixel-logo-box {
  border: 2px solid #78350f;
  box-shadow: 2px 2px 0 #451a03;
}

.pixel-box {
  border: 4px solid #1e293b;
  box-shadow: 
    inset 0 0 0 2px #334155,
    4px 4px 0 #0f172a;
}

.pixel-box-sm {
  border: 2px solid #334155;
  box-shadow: 2px 2px 0 #0f172a;
}

/* Mobile Bottom Navigation */
.bottom-nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  height: 100%;
  min-height: 64px;
  color: #94a3b8;
  border-bottom: 3px solid transparent;
  transition: all 0.15s;
}

.bottom-nav-item:hover {
  color: #fef3c7;
  background: rgba(30, 41, 59, 0.3);
}

.bottom-nav-item-active {
  color: #fbbf24;
  background: rgba(146, 64, 14, 0.15);
  border-bottom-color: #92400e;
}

/* Mobile Menu Cards */
.mobile-nav-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.75rem;
  border: 3px solid #334155;
  background: rgba(30, 41, 59, 0.5);
  transition: all 0.15s;
  min-height: 80px;
}

.mobile-nav-card:hover {
  background: rgba(30, 41, 59, 0.8);
  border-color: #92400e;
  transform: translateY(-2px);
}

.mobile-nav-card-active {
  background: rgba(146, 64, 14, 0.2);
  border-color: #92400e;
  box-shadow: inset 0 0 12px rgba(251, 191, 36, 0.15);
}

/* Mobile Menu Secondary Links */
.mobile-nav-link-secondary {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.875rem 1rem;
  font-size: 0.875rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #94a3b8;
  border: 2px solid transparent;
  border-radius: 0.5rem;
  transition: all 0.15s;
  min-height: 48px;
}

.mobile-nav-link-secondary:hover {
  color: #fef3c7;
  background: rgba(30, 41, 59, 0.6);
}

.mobile-nav-link-secondary-active {
  color: #fbbf24;
  background: rgba(146, 64, 14, 0.25);
  border-color: #92400e;
}

/* Desktop Nav */
.nav-link {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.5rem 0.75rem;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #94a3b8;
  border: 2px solid transparent;
  transition: all 0.15s;
}

.nav-link:hover {
  color: #fef3c7;
  background: rgba(30, 41, 59, 0.5);
  border-color: #334155;
}

.nav-link-active {
  color: #fbbf24;
  background: rgba(146, 64, 14, 0.2);
  border-color: #92400e;
  box-shadow: inset 0 0 8px rgba(251, 191, 36, 0.1);
}

.rpg-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.5rem 1rem;
  font-weight: 700;
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #fef3c7;
  background: linear-gradient(to bottom, #d97706, #b45309);
  border: 2px solid #92400e;
  box-shadow: 
    inset 0 1px 0 rgba(255, 255, 255, 0.2),
    0 3px 0 #78350f,
    0 4px 3px rgba(0, 0, 0, 0.3);
  cursor: pointer;
  transition: all 0.1s;
}

.rpg-btn:hover {
  background: linear-gradient(to bottom, #f59e0b, #d97706);
  transform: translateY(-1px);
}

.rpg-btn:active {
  transform: translateY(1px);
  box-shadow: 0 1px 0 #78350f;
}

.rpg-btn-small {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.375rem 0.625rem;
  font-weight: 600;
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #e2e8f0;
  background: linear-gradient(to bottom, #475569, #334155);
  border: 2px solid #1e293b;
  box-shadow: 
    inset 0 1px 0 rgba(255, 255, 255, 0.1),
    0 2px 0 #0f172a;
  cursor: pointer;
  transition: all 0.1s;
}

.rpg-btn-small:hover {
  background: linear-gradient(to bottom, #64748b, #475569);
}

.rpg-btn-small:active {
  transform: translateY(1px);
  box-shadow: 0 1px 0 #0f172a;
}

/* Safe area for mobile bottom nav */
@supports (padding-bottom: env(safe-area-inset-bottom)) {
  .bottom-nav-item {
    padding-bottom: env(safe-area-inset-bottom);
  }
}
</style>
