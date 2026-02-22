<template>
  <div class="p-6">
    <h1 class="text-3xl font-bold mb-6">
      Guilds
    </h1>
    
    <div v-if="loading" class="flex justify-center p-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-500"></div>
    </div>
    
    <div v-else>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="guild in guilds" :key="guild.id" class="bg-slate-900 border border-slate-800 p-6 rounded-lg">
          <h2 class="text-xl font-bold text-indigo-400">
            {{ guild.name }}
          </h2>
          <p class="text-slate-400 mt-2 text-sm h-12 overflow-hidden">
            {{ guild.description }}
          </p>
          <div class="mt-4 flex justify-between items-center">
            <span class="text-xs text-slate-500">Level {{ guild.level }}</span>
            <button 
              class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 rounded text-sm transition"
              @click="joinGuild(guild.id)"
            >
              Join
            </button>
          </div>
        </div>
      </div>
      
      <div v-if="guilds.length === 0" class="text-center p-12 bg-slate-900 rounded-lg border border-slate-800">
        <p class="text-slate-400">
          No guilds found. Be the first to start a legacy!
        </p>
        <button class="mt-4 px-6 py-2 bg-indigo-600 hover:bg-indigo-700 rounded transition font-bold">
          Create Guild
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';

interface Guild {
  id: string;
  name: string;
  description: string;
  level: number;
}

const guilds = ref<Guild[]>([]);
const loading = ref(true);

const fetchGuilds = async () => {
  try {
    const response = await fetch('/api/guilds');
    if (response.ok) {
      guilds.value = await response.json();
    }
  } catch (error) {
    console.error('Failed to fetch guilds:', error);
  } finally {
    loading.value = false;
  }
};

const joinGuild = async (id: string) => {
  try {
    const response = await fetch(`/api/guilds/${id}/join`, { method: 'POST' });
    if (response.ok) {
      alert('Join request submitted!');
    }
  } catch (error) {
    console.error('Failed to join guild:', error);
  }
};

onMounted(fetchGuilds);
</script>
