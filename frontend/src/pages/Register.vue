<template>
  <div class="flex min-h-[60vh] items-center justify-center">
    <BaseCard class-name="w-full max-w-md">
      <template #header>
        <div class="text-center">
          <h1 class="text-2xl font-bold text-white">Join the Arena</h1>
          <p class="text-sm text-slate-400 mt-1">Begin your journey as a Pixel Commander.</p>
        </div>
      </template>
      
      <form class="space-y-4" @submit.prevent="submit">
        <div class="space-y-1.5">
          <span class="text-xs font-semibold uppercase tracking-wider text-slate-500">Username</span>
          <input 
            v-model="username" 
            class="w-full rounded-md border border-slate-800 bg-slate-950 p-2.5 text-white transition-colors focus:border-indigo-500 focus:outline-none"
            placeholder="GhostCommander"
            required
          />
        </div>

        <div class="space-y-1.5">
          <span class="text-xs font-semibold uppercase tracking-wider text-slate-500">Email Address</span>
          <input 
            v-model="email" 
            type="email"
            class="w-full rounded-md border border-slate-800 bg-slate-950 p-2.5 text-white transition-colors focus:border-indigo-500 focus:outline-none"
            placeholder="commander@arena.com"
            required
          />
        </div>
        
        <div class="space-y-1.5">
          <span class="text-xs font-semibold uppercase tracking-wider text-slate-500">Password</span>
          <input 
            v-model="password" 
            type="password" 
            class="w-full rounded-md border border-slate-800 bg-slate-950 p-2.5 text-white transition-colors focus:border-indigo-500 focus:outline-none"
            placeholder="••••••••"
            required
          />
        </div>

        <div v-if="auth.error" class="rounded-md bg-red-500/10 p-3 text-sm text-red-400 border border-red-500/20">
          {{ auth.error }}
        </div>

        <div v-if="success" class="rounded-md bg-emerald-500/10 p-3 text-sm text-emerald-400 border border-emerald-500/20 text-center">
          Recruitment complete! You can now <router-link to="/login" class="underline font-bold">Sign In</router-link>.
        </div>

        <BaseButton class-name="w-full mt-2" :disabled="auth.isLoading">
          {{ auth.isLoading ? "Recruiting..." : "Create Commander Account" }}
        </BaseButton>

        <p class="text-center text-sm text-slate-500">
          Already a Commander? 
          <router-link to="/login" class="text-indigo-400 hover:text-indigo-300 font-medium">Access HQ</router-link>
        </p>
      </form>
    </BaseCard>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useAuthStore } from "@/features/auth/store";
import BaseButton from "@/shared/ui/BaseButton.vue";
import BaseCard from "@/shared/ui/BaseCard.vue";

const auth = useAuthStore();

const username = ref("");
const email = ref("");
const password = ref("");
const success = ref(false);

const submit = async () => {
  success.value = false;
  try {
    await auth.register(username.value, email.value, password.value);
    success.value = true;
  } catch (e) {
    // Error state in store
  }
};
</script>
