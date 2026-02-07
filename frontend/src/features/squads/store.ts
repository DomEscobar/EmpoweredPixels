import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { Squad, SetSquadRequest } from './api';
import { getActiveSquad, setActiveSquad } from './api';

export const useSquadStore = defineStore('squads', () => {
  const squad = ref<Squad | null>(null);
  const loading = ref(false);
  const error = ref<string | null>(null);

  const isActive = computed(() => squad.value?.isActive || false);
  const memberCount = computed(() => squad.value?.members.length || 0);
  const hasFullSquad = computed(() => memberCount.value === 3);

  async function loadSquad(token: string) {
    loading.value = true;
    error.value = null;
    try {
      squad.value = await getActiveSquad(token);
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to load squad';
    } finally {
      loading.value = false;
    }
  }

  async function updateSquad(token: string, data: SetSquadRequest) {
    loading.value = true;
    error.value = null;
    try {
      squad.value = await setActiveSquad(token, data.name, data.fighterIds);
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to update squad';
      throw err;
    } finally {
      loading.value = false;
    }
  }

  function clearError() {
    error.value = null;
  }

  return {
    squad,
    loading,
    error,
    isActive,
    memberCount,
    hasFullSquad,
    loadSquad,
    updateSquad,
    clearError,
  };
});
