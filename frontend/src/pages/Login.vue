<template>
  <div class="flex min-h-[60vh] items-center justify-center">
    <BaseCard class-name="w-full max-w-md">
      <template #header>
        <div class="text-center">
          <h1 class="text-2xl font-bold text-white">Welcome Back</h1>
          <p class="text-sm text-slate-400 mt-1">Commander, your fighters await your orders.</p>
        </div>
      </template>
      
      <form class="space-y-4" @submit.prevent="submit">
        <div class="space-y-1.5">
          <span class="text-xs font-semibold uppercase tracking-wider text-slate-500">Username or email</span>
          <input 
            v-model="user" 
            class="w-full rounded-md border border-slate-800 bg-slate-950 p-2.5 text-white transition-colors focus:border-indigo-500 focus:outline-none"
            placeholder="commander@example.com"
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

        <BaseButton class-name="w-full mt-2" :disabled="auth.isLoading">
          {{ auth.isLoading ? "Authorizing..." : "Sign In" }}
        </BaseButton>

        <p class="text-center text-sm text-slate-500">
          Don't have an account? 
          <router-link to="/register" class="text-indigo-400 hover:text-indigo-300 font-medium">Join the Arena</router-link>
        </p>
      </form>
    </BaseCard>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "@/features/auth/store";
import BaseButton from "@/shared/ui/BaseButton.vue";
import BaseCard from "@/shared/ui/BaseCard.vue";

const auth = useAuthStore();
const router = useRouter();

const user = ref("");
const password = ref("");

const submit = async () => {
  try {
    await auth.login(user.value, password.value);
    await router.push({ name: "dashboard" });
  } catch (e) {
    // Error state in store
  }
};
</script>
