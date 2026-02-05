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
      <div class="mx-auto flex h-14 max-w-7xl items-center justify-between px-4">
        <!-- Logo -->
        <router-link to="/" class="flex items-center gap-3 group">
          <div class="pixel-logo-box bg-amber-600 p-1.5 group-hover:bg-amber-500 transition-colors">
            <img :src="PIXEL_ASSETS.ICON_LOGO" alt="EP" class="w-6 h-6 pixelated" />
          </div>
          <span class="text-lg font-bold text-amber-400 text-shadow-retro tracking-wider hidden sm:block">
            EMPOWERED<span class="text-white">PIXELS</span>
          </span>
        </router-link>

        <!-- Desktop Nav Links -->
        <div v-if="auth.token" class="hidden md:flex items-center gap-1">
          <router-link 
            v-for="item in navItems" 
            :key="item.path" 
            :to="item.path"
            class="nav-link"
            :class="{ 'nav-link-active': isActive(item.path) }"
          >
            <img :src="item.icon" alt="" class="w-4 h-4 pixelated" />
            <span>{{ item.name }}</span>
          </router-link>
        </div>

        <!-- Right Side -->
        <div class="flex items-center gap-3">
          <template v-if="auth.token">
            <!-- Gold Display -->
            <router-link to="/shop" class="gold-display-nav pixel-box-sm bg-amber-900/30 border-amber-600/50 px-2 py-1 flex items-center gap-1.5 hover:bg-amber-900/50 transition-colors">
              <img :src="PIXEL_ASSETS.ICON_GOLD" alt="Gold" class="w-4 h-4 pixelated" />
              <span class="text-xs font-bold text-amber-400">{{ formattedGold() }}</span>
            </router-link>
            <button @click="logout" class="footer-logout-btn group p-1.5" title="Logout">
              <div class="pixel-box-sm bg-slate-800/80 p-1 group-hover:bg-red-900/50 transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400 group-hover:text-red-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="square" stroke-linejoin="miter">
                  <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
                  <polyline points="16 17 21 12 16 7"></polyline>
                  <line x1="21" y1="12" x2="9" y2="12"></line>
                </svg>
              </div>
            </button>
            <!-- User Badge -->
            <div class="pixel-box-sm bg-slate-800/80 px-3 py-1.5 hidden sm:flex items-center gap-2">
              <div class="w-6 h-6 pixel-box-sm bg-indigo-900/50 border-indigo-500/50 flex items-center justify-center">
                <img :src="PIXEL_ASSETS.ICON_USER" alt="" class="w-4 h-4 pixelated" />
              </div>
              <span class="text-xs text-slate-300 font-bold uppercase">Commander</span>
            </div>
          </template>
          <template v-else>
            <router-link to="/login">
              <button class="rpg-btn-small text-xs">SIGN IN</button>
            </router-link>
            <router-link to="/register" class="hidden sm:block">
              <button class="rpg-btn text-xs py-1.5 px-3">GET STARTED</button>
            </router-link>
          </template>

          <!-- Mobile Menu Toggle -->
          <button 
            v-if="auth.token"
            @click="mobileMenuOpen = !mobileMenuOpen"
            class="md:hidden rpg-btn-small p-2"
          >
            <span class="text-lg">{{ mobileMenuOpen ? '✕' : '☰' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Mobile Menu -->
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-2"
    >
      <div 
        v-if="mobileMenuOpen && auth.token" 
        class="md:hidden pixel-box bg-slate-900/98 border-t-0 mx-2 mt-0"
      >
        <div class="p-2 space-y-1">
          <router-link 
            v-for="item in navItems" 
            :key="item.path" 
            :to="item.path"
            @click="mobileMenuOpen = false"
            class="mobile-nav-link"
            :class="{ 'mobile-nav-link-active': isActive(item.path) }"
          >
            <img :src="item.icon" alt="" class="w-5 h-5 pixelated" />
            <span>{{ item.name }}</span>
          </router-link>
        </div>
      </div>
    </Transition>
  </nav>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '@/features/auth/store';
import { useShopStore } from '@/features/shop/store';

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
  ICON_GOLD: 'https://vibemedia.space/icon_gold_coin_nav_8f7e6d.png?prompt=golden%20coin%20with%20shine%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
};

const auth = useAuthStore();
const shop = useShopStore();
const router = useRouter();
const route = useRoute();
const mobileMenuOpen = ref(false);

const navItems = [
  { name: 'Command', path: '/dashboard', icon: PIXEL_ASSETS.ICON_DASHBOARD },
  { name: 'Roster', path: '/roster', icon: PIXEL_ASSETS.ICON_ROSTER },
  { name: 'Battle', path: '/matches', icon: PIXEL_ASSETS.ICON_MATCHES },
  { name: 'Vault', path: '/inventory', icon: PIXEL_ASSETS.ICON_INVENTORY },
  { name: 'Leagues', path: '/leagues', icon: PIXEL_ASSETS.ICON_LEAGUES },
  { name: 'Shop', path: '/shop', icon: PIXEL_ASSETS.ICON_SHOP },
];

// Fetch gold balance when logged in
onMounted(() => {
  if (auth.token) {
    shop.fetchGoldBalance();
  }
});

watch(() => auth.token, (newToken) => {
  if (newToken) {
    shop.fetchGoldBalance();
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

.mobile-nav-link {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  font-size: 0.875rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #94a3b8;
  border: 2px solid transparent;
  transition: all 0.15s;
}

.mobile-nav-link:hover {
  color: #fef3c7;
  background: rgba(30, 41, 59, 0.5);
}

.mobile-nav-link-active {
  color: #fbbf24;
  background: rgba(146, 64, 14, 0.3);
  border-color: #92400e;
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
</style>
