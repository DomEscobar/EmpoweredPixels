import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { attunementApi, elementConfig, calculateProgress } from './api';
import type { PlayerAttunements, AttunementWithBonuses, AggregatedBonuses, Attunement } from './types';

export const useAttunementStore = defineStore('attunement', () => {
  // State
  const attunements = ref<Attunement[]>([]);
  const bonuses = ref<AggregatedBonuses | null>(null);
  const loading = ref(false);
  const error = ref<string | null>(null);

  // Getters
  const totalLevel = computed(() => 
    attunements.value.reduce((sum, a) => sum + a.level, 0)
  );

  const getAttunementByElement = computed(() => (element: string) => {
    return attunements.value.find(a => a.element === element);
  });

  const getElementConfig = computed(() => (element: string) => {
    return elementConfig[element] || { name: element, icon: '❓', color: '#9ca3af', description: '' };
  });

  // Actions
  async function fetchAttunements() {
    loading.value = true;
    error.value = null;
    try {
      const result = await attunementApi.getAttunements();
      attunements.value = result.attunements;
    } catch (err: any) {
      error.value = err.message || 'Failed to load attunements';
    } finally {
      loading.value = false;
    }
  }

  async function fetchBonuses() {
    try {
      bonuses.value = await attunementApi.getBonuses();
    } catch (err: any) {
      console.error('Failed to fetch bonuses:', err);
    }
  }

  async function awardXP(element: string, source: string) {
    loading.value = true;
    error.value = null;
    try {
      const result = await attunementApi.awardXP(element, source);
      if (result.success) {
        await fetchAttunements();
        await fetchBonuses();
      }
      return result;
    } catch (err: any) {
      error.value = err.message || 'Failed to award XP';
      return null;
    } finally {
      loading.value = false;
    }
  }

  return {
    attunements,
    bonuses,
    loading,
    error,
    totalLevel,
    getAttunementByElement,
    getElementConfig,
    fetchAttunements,
    fetchBonuses,
    awardXP
  };
});
