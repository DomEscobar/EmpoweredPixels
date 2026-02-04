import { defineStore } from 'pinia';
import { ref } from 'vue';
import { useAuthStore } from '@/features/auth/store';
import * as api from './api';

export const useInventoryStore = defineStore('inventory', () => {
  const auth = useAuthStore();
  
  // State
  const particles = ref(0);
  const commonTokens = ref(0);
  const rareTokens = ref(0);
  const fabledTokens = ref(0);
  const mythicTokens = ref(0);
  
  const equipment = ref<api.Equipment[]>([]);
  const totalEquipment = ref(0);
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  // Actions
  async function fetchBalances() {
    if (!auth.token) return;
    try {
      const [p, c, r, f, m] = await Promise.all([
        api.getParticleBalance(auth.token),
        api.getTokenBalance(auth.token, 'common'),
        api.getTokenBalance(auth.token, 'rare'),
        api.getTokenBalance(auth.token, 'fabled'),
        api.getTokenBalance(auth.token, 'mythic')
      ]);
      
      particles.value = p;
      commonTokens.value = c;
      rareTokens.value = r;
      fabledTokens.value = f;
      mythicTokens.value = m;
    } catch (e) {
      console.error('Failed to fetch balances', e);
    }
  }

  async function fetchInventory(page = 1, pageSize = 20, filter?: api.InventoryPageRequest['filter'], sort?: api.InventoryPageRequest['sort']) {
    if (!auth.token) return;
    isLoading.value = true;
    error.value = null;
    try {
      const res = await api.getInventoryPage(auth.token, { page, pageSize, filter, sort });
      equipment.value = res.items;
      totalEquipment.value = res.totalCount;
    } catch (e: any) {
      error.value = e.message || 'Failed to fetch inventory';
    } finally {
      isLoading.value = false;
    }
  }

  async function enhanceItem(id: string) {
    if (!auth.token) return;
    try {
      const updated = await api.enhanceEquipment(auth.token, id);
      const index = equipment.value.findIndex(e => e.id === id);
      if (index !== -1) {
        equipment.value[index] = updated;
      }
      await fetchBalances(); // Update costs
      return updated;
    } catch (e: any) {
      throw e;
    }
  }

  async function salvageItem(id: string) {
    if (!auth.token) return;
    try {
      const result = await api.salvageEquipment(auth.token, id);
      equipment.value = equipment.value.filter(e => e.id !== id);
      totalEquipment.value--;
      await fetchBalances(); // Update particles
      return result;
    } catch (e: any) {
      throw e;
    }
  }

  async function toggleFavorite(id: string, isFavorite: boolean) {
    if (!auth.token) return;
    try {
      await api.setFavorite(auth.token, id, isFavorite);
      const item = equipment.value.find(e => e.id === id);
      if (item) {
        item.isFavorite = isFavorite;
      }
    } catch (e) {
      console.error('Failed to toggle favorite', e);
    }
  }

  return {
    particles,
    commonTokens,
    rareTokens,
    fabledTokens,
    mythicTokens,
    equipment,
    totalEquipment,
    isLoading,
    error,
    fetchBalances,
    fetchInventory,
    enhanceItem,
    salvageItem,
    toggleFavorite
  };
});
