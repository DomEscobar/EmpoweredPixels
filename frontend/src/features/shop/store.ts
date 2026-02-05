import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { useAuthStore } from "@/features/auth/store";
import * as api from "./api";

export const useShopStore = defineStore("shop", () => {
  const auth = useAuthStore();

  // State
  const goldPackages = ref<api.ShopItem[]>([]);
  const bundles = ref<api.ShopItem[]>([]);
  const playerGold = ref<api.PlayerGold | null>(null);
  const transactions = ref<api.Transaction[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const purchaseLoading = ref(false);
  const purchaseResult = ref<api.PurchaseResponse | null>(null);

  // Getters
  const goldBalance = computed(() => playerGold.value?.balance ?? 0);
  const hasGold = computed(() => goldBalance.value > 0);

  // Actions
  async function fetchGoldPackages() {
    if (!auth.token) return;
    try {
      loading.value = true;
      error.value = null;
      goldPackages.value = await api.getGoldPackages(auth.token);
    } catch (e) {
      error.value = e instanceof Error ? e.message : "Failed to load gold packages";
      console.error("Failed to fetch gold packages:", e);
    } finally {
      loading.value = false;
    }
  }

  async function fetchBundles() {
    if (!auth.token) return;
    try {
      loading.value = true;
      error.value = null;
      bundles.value = await api.getBundles(auth.token);
    } catch (e) {
      error.value = e instanceof Error ? e.message : "Failed to load bundles";
      console.error("Failed to fetch bundles:", e);
    } finally {
      loading.value = false;
    }
  }

  async function fetchPlayerGold() {
    if (!auth.token) return;
    try {
      playerGold.value = await api.getPlayerGold(auth.token);
    } catch (e) {
      console.error("Failed to fetch player gold:", e);
      // Initialize with 0 on error
      playerGold.value = { balance: 0, lifetime_earned: 0, lifetime_spent: 0 };
    }
  }

  async function fetchTransactions() {
    if (!auth.token) return;
    try {
      loading.value = true;
      error.value = null;
      transactions.value = await api.getTransactions(auth.token);
    } catch (e) {
      error.value = e instanceof Error ? e.message : "Failed to load transactions";
      console.error("Failed to fetch transactions:", e);
    } finally {
      loading.value = false;
    }
  }

  async function purchase(itemId: string): Promise<boolean> {
    if (!auth.token) return false;
    try {
      purchaseLoading.value = true;
      purchaseResult.value = null;
      error.value = null;

      const result = await api.purchaseItem(auth.token, itemId);
      purchaseResult.value = result;

      // Update local gold balance
      if (playerGold.value) {
        playerGold.value.balance = result.new_balance;
      }

      return result.success;
    } catch (e) {
      error.value = e instanceof Error ? e.message : "Purchase failed";
      console.error("Purchase failed:", e);
      return false;
    } finally {
      purchaseLoading.value = false;
    }
  }

  async function loadShopData() {
    await Promise.all([
      fetchGoldPackages(),
      fetchBundles(),
      fetchPlayerGold(),
    ]);
  }

  function canAfford(item: api.ShopItem): boolean {
    if (item.price_currency !== 'gold') return true; // USD purchases handled separately
    return goldBalance.value >= item.price_amount;
  }

  function clearError() {
    error.value = null;
  }

  function clearPurchaseResult() {
    purchaseResult.value = null;
  }

  return {
    // State
    goldPackages,
    bundles,
    playerGold,
    transactions,
    loading,
    error,
    purchaseLoading,
    purchaseResult,

    // Getters
    goldBalance,
    hasGold,

    // Actions
    fetchGoldPackages,
    fetchBundles,
    fetchPlayerGold,
    fetchTransactions,
    purchase,
    loadShopData,
    canAfford,
    clearError,
    clearPurchaseResult,
  };
});
