<template>
  <div class="relative overflow-hidden" data-testid="home-page">
    <!-- Background Decor -->
    <div class="absolute -top-24 -left-24 h-96 w-96 rounded-full bg-indigo-600/10 blur-3xl"></div>
    <div class="absolute top-1/2 -right-24 h-64 w-64 rounded-full bg-emerald-600/10 blur-3xl"></div>

    <section class="relative flex min-h-[calc(100vh-160px)] flex-col items-center justify-center text-center">
      <div class="mb-6 rounded-full bg-indigo-500/10 px-4 py-1.5 text-sm font-semibold text-indigo-400 border border-indigo-500/20">
        New: Deathmatch Alpha League is Live
      </div>
      
      <h1 class="max-w-4xl text-5xl font-extrabold tracking-tight text-white sm:text-7xl">
        Master the Pixels, <br />
        <span class="bg-gradient-to-r from-indigo-400 to-emerald-400 bg-clip-text text-transparent">Rule the Battlefield</span>
      </h1>
      
      <p class="mt-8 max-w-2xl text-lg text-slate-400 leading-relaxed">
        Assemble your ultimate squad of pixelated fighters, equip them with legendary gear, and climb the ranks in competitive leagues. The arena awaits your command.
      </p>

      <div class="mt-10 flex flex-col items-center gap-4 sm:flex-row">
        <template v-if="auth.token">
          <router-link to="/dashboard" data-testid="enter-dashboard-btn">
            <BaseButton size="lg" class-name="px-10 text-lg shadow-lg shadow-indigo-500/20">Enter Dashboard</BaseButton>
          </router-link>
        </template>
        <template v-else>
          <router-link to="/register" data-testid="get-started-btn">
            <BaseButton size="lg" class-name="px-10 text-lg shadow-lg shadow-indigo-500/20">Get Started Free</BaseButton>
          </router-link>
          <router-link to="/login" data-testid="watch-live-btn">
            <BaseButton variant="outline" size="lg" class-name="px-10 text-lg border-slate-700">Watch Live Matches</BaseButton>
          </router-link>
        </template>
      </div>

      <div class="mt-20 grid grid-cols-2 gap-8 text-slate-500 sm:grid-cols-4 grayscale opacity-50">
        <div class="font-bold tracking-widest uppercase italic text-sm" data-testid="tag-empowered">Empowered</div>
        <div class="font-bold tracking-widest uppercase italic text-sm" data-testid="tag-competitive">Competitive</div>
        <div class="font-bold tracking-widest uppercase italic text-sm" data-testid="tag-strategic">Strategic</div>
        <div class="font-bold tracking-widest uppercase italic text-sm" data-testid="tag-realtime">Realtime</div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/features/auth/store';
import BaseButton from '@/shared/ui/BaseButton.vue';

const auth = useAuthStore();
</script>
