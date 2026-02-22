<template>
  <div class="home-page bg-slate-950 text-slate-100" data-testid="home-page">
    <section class="hero-section relative overflow-hidden min-h-screen flex flex-col items-center justify-center">
      <div class="absolute inset-0 bg-grid-pattern opacity-5"></div>
      <div class="hero-bg-gradient absolute top-0 left-0 w-full h-full"></div>
      
      <div class="floating-particles">
        <div v-for="i in 20" :key="i" class="particle" :style="{ '--delay': i * 0.3 + 's' }"></div>
      </div>

      <div class="hero-content relative z-10 text-center px-4 max-w-6xl mx-auto">
        <div class="hero-badge mb-6 animate-fade-in-up">
          <span class="inline-flex items-center gap-2 bg-gradient-to-r from-indigo-600 to-purple-600 text-white px-4 py-2 rounded-full text-sm font-bold uppercase tracking-wider border border-indigo-400/30 shadow-lg shadow-indigo-500/20">
            <span class="w-2 h-2 bg-green-400 rounded-full animate-pulse"></span>
            Deathmatch Alpha League Live
          </span>
        </div>

        <h1 class="hero-title animate-fade-in-up" style="--delay: 0.1s">
          <span class="block text-6xl md:text-8xl font-extrabold tracking-tight mb-2">
            <span class="text-transparent bg-clip-text bg-gradient-to-r from-amber-400 via-orange-400 to-red-400 animate-gradient-shift">
              EMPOWERED
            </span>
          </span>
          <span class="block text-5xl md:text-7xl font-extrabold tracking-tight text-white">
            <span class="text-transparent bg-clip-text bg-gradient-to-r from-cyan-400 via-blue-400 to-indigo-400 animate-gradient-shift">
              PIXELS
            </span>
          </span>
        </h1>

        <p class="hero-subtitle animate-fade-in-up text-xl md:text-2xl text-slate-400 max-w-3xl mx-auto mt-6 leading-relaxed" style="--delay: 0.2s">
          Command legendary voxel warriors, master elemental attunements, and dominate the arena in this strategic auto-battler. Your destiny awaits.
        </p>

        <div class="hero-cta mt-10 flex flex-col items-center gap-4 animate-fade-in-up" style="--delay: 0.3s">
          <template v-if="auth.token">
            <router-link to="/dashboard" data-testid="enter-dashboard-btn">
              <button class="rpg-btn-lg group relative overflow-hidden">
                <span class="relative z-10 flex items-center gap-2">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                  </svg>
                  Enter Command Center
                </span>
                <div class="absolute inset-0 bg-gradient-to-r from-amber-400 to-orange-500 opacity-0 group-hover:opacity-100 transition-opacity"></div>
              </button>
            </router-link>
          </template>
          <template v-else>
            <router-link to="/register" data-testid="get-started-btn">
              <button class="rpg-btn-lg group relative overflow-hidden">
                <span class="relative z-10 flex items-center gap-2">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                  </svg>
                  Start Your Journey
                </span>
                <div class="absolute inset-0 bg-gradient-to-r from-amber-400 to-orange-500 opacity-0 group-hover:opacity-100 transition-opacity"></div>
              </button>
            </router-link>
            <router-link to="/login" data-testid="watch-live-btn">
              <button class="rpg-btn-lg-secondary group flex items-center gap-2">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
                Watch Live Battles
              </button>
            </router-link>
          </template>
        </div>

        <div class="hero-stats mt-16 grid grid-cols-2 md:grid-cols-4 gap-6 animate-fade-in-up" style="--delay: 0.4s">
          <div v-for="stat in heroStats" :key="stat.label" class="stat-card">
            <div class="stat-value text-3xl md:text-4xl font-extrabold text-transparent bg-clip-text bg-gradient-to-br from-amber-300 to-amber-500">
              {{ stat.value }}
            </div>
            <div class="stat-label text-sm text-slate-500 uppercase tracking-wider font-semibold mt-1">
              {{ stat.label }}
            </div>
          </div>
        </div>
      </div>

      <div class="hero-fighters absolute inset-0 pointer-events-none">
        <div class="fighter-showcase-left">
          <VoxelFighter seed="hero-fire" attunement="Fire" :animate="true" class="w-64 h-64 md:w-96 md:h-96 opacity-30" />
        </div>
        <div class="fighter-showcase-right">
          <VoxelFighter seed="hero-water" attunement="Water" :animate="true" class="w-64 h-64 md:w-96 md:h-96 opacity-30" />
        </div>
      </div>
    </section>

    <section class="features-section py-24 px-4 relative">
      <div class="max-w-7xl mx-auto">
        <div class="text-center mb-16">
          <h2 class="text-4xl md:text-5xl font-extrabold text-white mb-4">
            <span class="text-transparent bg-clip-text bg-gradient-to-r from-cyan-400 to-blue-500">
              Master the Arena
            </span>
          </h2>
          <p class="text-xl text-slate-400 max-w-2xl mx-auto">
            Strategic depth meets visceral combat in every battle
          </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          <div v-for="feature in features" :key="feature.title" class="feature-card ep-card-iron">
            <div class="feature-icon-wrapper">
              <div :class="['feature-icon', feature.iconColor]">
                <component :is="feature.icon" class="w-8 h-8" />
              </div>
            </div>
            <h3 class="text-xl font-bold text-white mb-3">
              {{ feature.title }}
            </h3>
            <p class="text-slate-400 leading-relaxed">
              {{ feature.description }}
            </p>
          </div>
        </div>
      </div>
    </section>

    <section class="gameplay-section py-24 px-4 bg-gradient-to-b from-slate-950 to-slate-900">
      <div class="max-w-7xl mx-auto">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-16 items-center">
          <div class="gameplay-content">
            <h2 class="text-4xl md:text-5xl font-extrabold text-white mb-6">
              <span class="text-transparent bg-clip-text bg-gradient-to-r from-amber-400 to-orange-500">
                Auto-Battler
              </span>
              <br>
              Strategy Evolved
            </h2>
            <p class="text-lg text-slate-400 leading-relaxed mb-8">
              Watch your squad of voxel warriors engage in epic auto-battles. Configure loadouts, master elemental synergies, and outsmart your opponents through strategic positioning and timing.
            </p>
            <div class="gameplay-highlights space-y-4">
              <div v-for="highlight in gameplayHighlights" :key="highlight.text" class="highlight-item flex items-start gap-4">
                <div class="highlight-icon bg-indigo-500/20 border border-indigo-500/30 p-2 rounded-lg">
                  <svg class="w-5 h-5 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                </div>
                <div>
                  <h4 class="font-bold text-white mb-1">
                    {{ highlight.title }}
                  </h4>
                  <p class="text-slate-400 text-sm">
                    {{ highlight.text }}
                  </p>
                </div>
              </div>
            </div>
          </div>
          <div class="gameplay-visual relative">
            <div class="visual-container relative bg-slate-800/30 border border-slate-700 rounded-2xl overflow-hidden shadow-2xl">
              <div class="visual-header bg-slate-900/80 border-b border-slate-700 px-6 py-4 flex items-center justify-between">
                <div class="flex items-center gap-3">
                  <div class="w-3 h-3 bg-red-500 rounded-full"></div>
                  <div class="w-3 h-3 bg-yellow-500 rounded-full"></div>
                  <div class="w-3 h-3 bg-green-500 rounded-full"></div>
                </div>
                <span class="text-xs text-slate-500 font-mono">LIVE BATTLE</span>
              </div>
              <div class="visual-body h-80 relative flex items-center justify-center">
                <div class="battle-scene">
                  <VoxelFighter seed="scene1" attunement="Fire" :animate="true" :is-attacking="true" class="w-32 h-32" />
                  <div class="battle-vs">
                    VS
                  </div>
                  <VoxelFighter seed="scene2" attunement="Water" :animate="true" class="w-32 h-32" />
                </div>
              </div>
            </div>
            <div class="visual-glow absolute inset-0 bg-gradient-to-r from-red-500/20 via-blue-500/20 to-purple-500/20 blur-3xl -z-10"></div>
          </div>
        </div>
      </div>
    </section>

    <section class="leagues-section py-24 px-4">
      <div class="max-w-7xl mx-auto">
        <div class="text-center mb-16">
          <h2 class="text-4xl md:text-5xl font-extrabold text-white mb-4">
            <span class="text-transparent bg-clip-text bg-gradient-to-r from-amber-400 via-yellow-400 to-orange-400">
              Competitive
            </span>
            <br>
            League System
          </h2>
          <p class="text-xl text-slate-400 max-w-2xl mx-auto">
            Prove your worth in structured leagues and climb the ranks to glory
          </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div v-for="league in leagues" :key="league.name" class="league-card relative overflow-hidden">
            <div class="league-bg-gradient absolute inset-0" :style="league.gradient"></div>
            <div class="relative z-10 p-8 text-center">
              <div class="league-icon mb-6">
                <component :is="league.icon" class="w-16 h-16 mx-auto" :class="league.iconColor" />
              </div>
              <h3 class="text-2xl font-bold text-white mb-3">
                {{ league.name }}
              </h3>
              <p class="text-slate-300 text-sm mb-4">
                {{ league.description }}
              </p>
              <div class="league-stats flex justify-center gap-4 text-sm">
                <div>
                  <div class="font-bold text-white">
                    {{ league.players }}
                  </div>
                  <div class="text-slate-400">
                    Players
                  </div>
                </div>
                <div>
                  <div class="font-bold text-white">
                    {{ league.matches }}
                  </div>
                  <div class="text-slate-400">
                    Matches
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="cta-section py-24 px-4 relative overflow-hidden">
      <div class="absolute inset-0 bg-gradient-to-r from-indigo-900/50 via-purple-900/50 to-pink-900/50"></div>
      <div class="absolute inset-0 bg-grid-pattern opacity-10"></div>
      
      <div class="max-w-4xl mx-auto text-center relative z-10">
        <h2 class="text-4xl md:text-6xl font-extrabold text-white mb-6">
          Ready to
          <span class="text-transparent bg-clip-text bg-gradient-to-r from-amber-400 to-orange-500">
            Command?
          </span>
        </h2>
        <p class="text-xl text-slate-300 mb-10 max-w-2xl mx-auto">
          Join thousands of warriors already competing in the arena. Your legend begins with a single match.
        </p>
        
        <template v-if="!auth.token">
          <router-link to="/register" data-testid="final-cta-btn">
            <button class="rpg-btn-xl group relative overflow-hidden">
              <span class="relative z-10 flex items-center gap-3">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
                Begin Your Campaign
              </span>
              <div class="absolute inset-0 bg-gradient-to-r from-amber-400 via-orange-400 to-red-400 opacity-0 group-hover:opacity-100 transition-opacity"></div>
            </button>
          </router-link>
        </template>
        <template v-else>
          <router-link to="/dashboard" data-testid="dashboard-cta-btn">
            <button class="rpg-btn-xl group relative overflow-hidden">
              <span class="relative z-10 flex items-center gap-3">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
                Return to Command
              </span>
              <div class="absolute inset-0 bg-gradient-to-r from-amber-400 via-orange-400 to-red-400 opacity-0 group-hover:opacity-100 transition-opacity"></div>
            </button>
          </router-link>
        </template>
      </div>
    </section>

    <footer class="footer-section py-12 px-4 border-t border-slate-800">
      <div class="max-w-7xl mx-auto">
        <div class="grid grid-cols-2 md:grid-cols-4 gap-8 mb-12">
          <div v-for="column in footerColumns" :key="column.title">
            <h4 class="font-bold text-white mb-4">
              {{ column.title }}
            </h4>
            <ul class="space-y-2">
              <li v-for="link in column.links" :key="link.text">
                <router-link v-if="link.to" :to="link.to" class="text-slate-400 hover:text-amber-400 transition-colors text-sm">
                  {{ link.text }}
                </router-link>
                <a v-else :href="link.href" class="text-slate-400 hover:text-amber-400 transition-colors text-sm">
                  {{ link.text }}
                </a>
              </li>
            </ul>
          </div>
        </div>
        <div class="footer-bottom border-t border-slate-800 pt-8 flex flex-col md:flex-row items-center justify-between gap-4">
          <div class="flex items-center gap-2">
            <div class="w-8 h-8 bg-amber-600 rounded p-1">
              <img src="https://vibemedia.space/icon_ep_logo_4a5b6c_v1.png?prompt=pixel%20art%20sword%20and%20shield%20emblem%20icon&style=pixel_game_asset&key=NOGON" alt="EP" class="w-full h-full object-contain pixelated" />
            </div>
            <span class="text-white font-bold">EmpoweredPixels</span>
          </div>
          <p class="text-slate-500 text-sm">
            © 2024 EmpoweredPixels. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/features/auth/store';
import VoxelFighter from '@/shared/ui/VoxelFighter.vue';

const auth = useAuthStore();

const heroStats = [
  { value: '50K+', label: 'Warriors' },
  { value: '100K+', label: 'Battles' },
  { value: '4', label: 'Elements' },
  { value: '∞', label: 'Strategies' },
];

const features = [
  {
    title: 'Dynamic Squad Builder',
    description: 'Assemble your ultimate team of voxel fighters with unique appearances, abilities, and synergies.',
    iconColor: 'bg-red-500/20 text-red-400',
    icon: 'svg',
  },
  {
    title: 'Elemental Attunement',
    description: 'Master Fire, Water, Earth, Wind, and Lightning. Each element offers unique tactical advantages.',
    iconColor: 'bg-blue-500/20 text-blue-400',
    icon: 'svg',
  },
  {
    title: 'Strategic Auto-Battles',
    description: 'Configure loadouts and positioning, then watch your squad execute devastating combo chains.',
    iconColor: 'bg-purple-500/20 text-purple-400',
    icon: 'svg',
  },
  {
    title: 'Competitive Leagues',
    description: 'Climb the ranks in seasonal leagues, earn rewards, and prove your dominance.',
    iconColor: 'bg-amber-500/20 text-amber-400',
    icon: 'svg',
  },
  {
    title: 'Loot & Progression',
    description: 'Collect legendary weapons, armor, and artifacts. Upgrade your fighters to unlock their full potential.',
    iconColor: 'bg-green-500/20 text-green-400',
    icon: 'svg',
  },
  {
    title: 'Live Spectate Mode',
    description: 'Watch battles in real-time, learn from top players, and scout your competition.',
    iconColor: 'bg-cyan-500/20 text-cyan-400',
    icon: 'svg',
  },
];

const gameplayHighlights = [
  {
    title: 'Smart Positioning',
    text: 'Place your fighters strategically to maximize damage and minimize incoming threats.',
  },
  {
    title: 'Elemental Synergies',
    text: 'Combine elements to create devastating combo chains and tactical advantages.',
  },
  {
    title: 'Loadout Customization',
    text: 'Equip unique weapons and abilities to suit your playstyle and counter opponents.',
  },
];

const leagues = [
  {
    name: 'Deathmatch Alpha',
    description: 'High-octane 1v1 battles. Pure skill, no mercy.',
    icon: 'svg',
    iconColor: 'text-red-400',
    gradient: 'background: linear-gradient(135deg, rgba(239, 68, 68, 0.1) 0%, rgba(153, 27, 27, 0.1) 100%)',
    players: '12K',
    matches: '45K',
  },
  {
    name: 'Team Tactics',
    description: 'Coordinate with your squad in 3v3 tactical warfare.',
    icon: 'svg',
    iconColor: 'text-blue-400',
    gradient: 'background: linear-gradient(135deg, rgba(59, 130, 246, 0.1) 0%, rgba(30, 58, 138, 0.1) 100%)',
    players: '8K',
    matches: '32K',
  },
  {
    name: 'Championship',
    description: 'The ultimate test. Only the best survive.',
    icon: 'svg',
    iconColor: 'text-amber-400',
    gradient: 'background: linear-gradient(135deg, rgba(251, 191, 36, 0.1) 0%, rgba(180, 83, 9, 0.1) 100%)',
    players: '2K',
    matches: '8K',
  },
];

const footerColumns = [
  {
    title: 'Game',
    links: [
      { text: 'Battles', to: '/matches' as const, href: undefined },
      { text: 'Leagues', to: '/leagues' as const, href: undefined },
      { text: 'Leaderboards', to: '/leaderboard' as const, href: undefined },
    ],
  },
  {
    title: 'Community',
    links: [
      { text: 'Discord', to: undefined, href: '#' as const },
      { text: 'Twitter', to: undefined, href: '#' as const },
      { text: 'Forums', to: undefined, href: '#' as const },
    ],
  },
  {
    title: 'Support',
    links: [
      { text: 'Help Center', to: undefined, href: '#' as const },
      { text: 'Contact', to: undefined, href: '#' as const },
      { text: 'Status', to: undefined, href: '#' as const },
    ],
  },
  {
    title: 'Legal',
    links: [
      { text: 'Privacy', to: undefined, href: '#' as const },
      { text: 'Terms', to: undefined, href: '#' as const },
      { text: 'Cookies', to: undefined, href: '#' as const },
    ],
  },
];
</script>

<style scoped>
.home-page {
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
}

.bg-grid-pattern {
  background-image: 
    linear-gradient(rgba(148, 163, 184, 0.1) 1px, transparent 1px),
    linear-gradient(90deg, rgba(148, 163, 184, 0.1) 1px, transparent 1px);
  background-size: 50px 50px;
}

.hero-bg-gradient {
  background: radial-gradient(ellipse at center top, rgba(99, 102, 241, 0.15) 0%, transparent 50%),
              radial-gradient(ellipse at center bottom, rgba(234, 179, 8, 0.1) 0%, transparent 50%);
}

.floating-particles {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.particle {
  position: absolute;
  width: 4px;
  height: 4px;
  background: rgba(251, 191, 36, 0.6);
  border-radius: 50%;
  animation: float-up 8s ease-in-out infinite;
  animation-delay: var(--delay);
}

@keyframes float-up {
  0%, 100% {
    transform: translateY(100vh) translateX(0) scale(0);
    opacity: 0;
  }
  10% {
    opacity: 1;
    transform: translateY(80vh) translateX(20px) scale(1);
  }
  90% {
    opacity: 1;
    transform: translateY(20vh) translateX(-20px) scale(1);
  }
  100% {
    transform: translateY(0) translateX(0) scale(0);
    opacity: 0;
  }
}

.animate-fade-in-up {
  animation: fade-in-up 0.8s ease-out forwards;
  opacity: 0;
}

@keyframes fade-in-up {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-gradient-shift {
  background-size: 200% 200%;
  animation: gradient-shift 8s ease infinite;
}

@keyframes gradient-shift {
  0%, 100% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
}

.hero-badge {
  animation: fade-in-up 0.6s ease-out forwards;
  animation-delay: var(--delay);
}

.hero-title,
.hero-subtitle,
.hero-cta,
.hero-stats {
  animation: fade-in-up 0.8s ease-out forwards;
  animation-delay: var(--delay);
}

.rpg-btn-lg {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 1rem 2.5rem;
  font-size: 1.125rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: #fef3c7;
  background: linear-gradient(to bottom, #d97706, #b45309);
  border: 3px solid #92400e;
  box-shadow: 
    inset 0 1px 0 rgba(255, 255, 255, 0.2),
    0 4px 0 #78350f,
    0 6px 12px rgba(0, 0, 0, 0.4);
  cursor: pointer;
  transition: all 0.15s;
}

.rpg-btn-lg:hover {
  transform: translateY(-2px);
  box-shadow: 
    inset 0 1px 0 rgba(255, 255, 255, 0.2),
    0 6px 0 #78350f,
    0 8px 16px rgba(0, 0, 0, 0.5);
}

.rpg-btn-lg:active {
  transform: translateY(1px);
  box-shadow: 
    inset 0 1px 0 rgba(255, 255, 255, 0.2),
    0 2px 0 #78350f,
    0 4px 8px rgba(0, 0, 0, 0.4);
}

.rpg-btn-lg-secondary {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 1rem 2.5rem;
  font-size: 1.125rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: #e2e8f0;
  background: linear-gradient(to bottom, #475569, #334155);
  border: 2px solid #1e293b;
  box-shadow: 
    inset 0 1px 0 rgba(255, 255, 255, 0.1),
    0 3px 0 #0f172a;
  cursor: pointer;
  transition: all 0.15s;
}

.rpg-btn-lg-secondary:hover {
  background: linear-gradient(to bottom, #64748b, #475569);
  transform: translateY(-1px);
}

.rpg-btn-lg-secondary:active {
  transform: translateY(1px);
}

.rpg-btn-xl {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 1.25rem 3rem;
  font-size: 1.25rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: #fef3c7;
  background: linear-gradient(to bottom, #d97706, #b45309);
  border: 4px solid #92400e;
  box-shadow: 
    inset 0 1px 0 rgba(255, 255, 255, 0.2),
    0 6px 0 #78350f,
    0 8px 20px rgba(0, 0, 0, 0.5);
  cursor: pointer;
  transition: all 0.15s;
}

.rpg-btn-xl:hover {
  transform: translateY(-3px);
  box-shadow: 
    inset 0 1px 0 rgba(255, 255, 255, 0.2),
    0 8px 0 #78350f,
    0 12px 24px rgba(0, 0, 0, 0.6);
}

.rpg-btn-xl:active {
  transform: translateY(1px);
  box-shadow: 
    inset 0 1px 0 rgba(255, 255, 255, 0.2),
    0 2px 0 #78350f,
    0 4px 12px rgba(0, 0, 0, 0.4);
}

.stat-card {
  background: rgba(30, 41, 59, 0.5);
  border: 2px solid rgba(148, 163, 184, 0.2);
  border-radius: 0.75rem;
  padding: 1.5rem;
  text-align: center;
  backdrop-filter: blur(8px);
  transition: all 0.3s;
}

.stat-card:hover {
  border-color: rgba(251, 191, 36, 0.3);
  background: rgba(30, 41, 59, 0.7);
  transform: translateY(-4px);
}

.hero-fighters {
  z-index: 1;
}

.fighter-showcase-left,
.fighter-showcase-right {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  animation: breathe 6s ease-in-out infinite;
}

.fighter-showcase-left {
  left: 5%;
  animation-delay: 0s;
}

.fighter-showcase-right {
  right: 5%;
  animation-delay: 3s;
}

@keyframes breathe {
  0%, 100% {
    transform: translateY(-50%) scale(1);
  }
  50% {
    transform: translateY(-50%) scale(1.05);
  }
}

.feature-card {
  padding: 2rem;
  transition: all 0.3s;
  border: 2px solid rgba(148, 163, 184, 0.2);
}

.feature-card:hover {
  border-color: rgba(99, 102, 241, 0.4);
  transform: translateY(-4px);
  box-shadow: 0 20px 40px -12px rgba(99, 102, 241, 0.2);
}

.feature-icon-wrapper {
  display: flex;
  justify-content: center;
  margin-bottom: 1.5rem;
}

.feature-icon {
  width: 3.5rem;
  height: 3.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 0.75rem;
  border: 2px solid rgba(255, 255, 255, 0.1);
}

.highlight-item {
  transition: all 0.3s;
}

.highlight-item:hover {
  transform: translateX(8px);
}

.highlight-icon {
  transition: all 0.3s;
}

.highlight-item:hover .highlight-icon {
  background: rgba(99, 102, 241, 0.3);
  border-color: rgba(99, 102, 241, 0.5);
}

.visual-container {
  animation: float 6s ease-in-out infinite;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

.battle-scene {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 2rem;
}

.battle-vs {
  font-size: 2rem;
  font-weight: 900;
  color: rgba(255, 255, 255, 0.1);
  text-shadow: 0 0 20px rgba(255, 255, 255, 0.2);
}

.league-card {
  border: 2px solid rgba(148, 163, 184, 0.2);
  border-radius: 1rem;
  transition: all 0.3s;
  background: rgba(30, 41, 59, 0.5);
  backdrop-filter: blur(8px);
}

.league-card:hover {
  border-color: rgba(251, 191, 36, 0.4);
  transform: translateY(-4px);
  box-shadow: 0 20px 40px -12px rgba(251, 191, 36, 0.2);
}

.league-icon {
  transition: all 0.3s;
}

.league-card:hover .league-icon {
  transform: scale(1.1) rotate(5deg);
}

.pixelated {
  image-rendering: pixelated;
}

@media (max-width: 768px) {
  .fighter-showcase-left,
  .fighter-showcase-right {
    display: none;
  }
  
  .hero-title {
    font-size: 2.5rem;
  }
  
  .battle-scene {
    gap: 1rem;
  }
  
  .battle-scene > * {
    width: 100px;
    height: 100px;
  }
}
</style>
