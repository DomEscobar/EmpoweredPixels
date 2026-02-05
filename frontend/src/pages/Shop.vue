<template>
  <div class="shop-page min-h-screen p-4 md:p-6">
    <div class="max-w-7xl mx-auto space-y-8">
      <!-- Header -->
      <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <div>
          <h1 class="text-3xl font-bold flex items-center gap-3">
            <img :src="SHOP_ICON" alt="" class="w-10 h-10 pixelated" />
            <span class="text-amber-400">SHOP</span>
          </h1>
          <p class="text-slate-400 mt-1">Acquire gold and powerful equipment bundles</p>
        </div>
        <GoldDisplay />
      </div>

      <!-- Tabs -->
      <div class="flex gap-2 border-b-4 border-slate-800 pb-2">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          class="tab-btn px-4 py-2 font-bold text-sm uppercase tracking-wide rounded-t transition-all"
          :class="activeTab === tab.id ? 'bg-slate-800 text-amber-400 border-b-4 border-amber-500 -mb-[6px]' : 'text-slate-400 hover:text-white'"
          @click="activeTab = tab.id"
        >
          {{ tab.name }}
        </button>
      </div>

      <!-- Loading State -->
      <div v-if="shop.loading" class="flex items-center justify-center py-20">
        <div class="text-center">
          <div class="animate-spin text-4xl mb-4">⚙️</div>
          <p class="text-slate-400">Loading shop...</p>
        </div>
      </div>

      <!-- Error State -->
      <div v-else-if="shop.error" class="p-6 bg-red-900/20 border border-red-600/50 rounded-lg text-center">
        <p class="text-red-400 mb-4">{{ shop.error }}</p>
        <button class="rpg-btn" @click="shop.loadShopData()">Try Again</button>
      </div>

      <!-- Gold Packages Tab -->
      <div v-else-if="activeTab === 'gold'" class="space-y-6">
        <div class="flex items-center gap-2 text-sm text-slate-400">
          <span>💰</span>
          <span>Purchase gold to buy bundles and upgrades</span>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <GoldPackageCard
            v-for="pkg in shop.goldPackages"
            :key="pkg.id"
            :item="pkg"
            @purchase="openPurchaseModal"
          />
        </div>

        <div v-if="shop.goldPackages.length === 0" class="text-center py-12 text-slate-500">
          No gold packages available at the moment.
        </div>
      </div>

      <!-- Bundles Tab -->
      <div v-else-if="activeTab === 'bundles'" class="space-y-6">
        <div class="flex items-center gap-2 text-sm text-slate-400">
          <span>📦</span>
          <span>Equipment bundles with guaranteed rarity items</span>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <BundleCard
            v-for="bundle in shop.bundles"
            :key="bundle.id"
            :item="bundle"
            @purchase="openPurchaseModal"
          />
        </div>

        <div v-if="shop.bundles.length === 0" class="text-center py-12 text-slate-500">
          No bundles available at the moment.
        </div>
      </div>

      <!-- Purchase History (optional view) -->
      <div v-if="showHistory" class="mt-8 p-4 bg-slate-800/50 rounded-lg border border-slate-700">
        <h3 class="font-bold mb-4 flex items-center gap-2">
          <span>📜</span>
          <span>Recent Purchases</span>
        </h3>
        <div v-if="shop.transactions.length === 0" class="text-slate-500 text-sm">
          No purchase history yet.
        </div>
        <div v-else class="space-y-2">
          <div 
            v-for="tx in shop.transactions.slice(0, 5)" 
            :key="tx.id"
            class="flex items-center justify-between text-sm py-2 border-b border-slate-700 last:border-0"
          >
            <div>
              <span class="font-medium">{{ tx.item_name }}</span>
              <span class="text-slate-500 ml-2">{{ formatDate(tx.created) }}</span>
            </div>
            <div :class="tx.gold_change > 0 ? 'text-green-400' : 'text-amber-400'">
              {{ tx.gold_change > 0 ? '+' : '' }}{{ tx.gold_change.toLocaleString() }} Gold
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Purchase Modal -->
    <PurchaseModal
      :visible="purchaseModalVisible"
      :item="selectedItem"
      :loading="shop.purchaseLoading"
      @close="closePurchaseModal"
      @confirm="confirmPurchase"
    />

    <!-- Success Toast -->
    <Transition name="toast">
      <div 
        v-if="successMessage" 
        class="fixed bottom-4 right-4 p-4 bg-green-900/90 border border-green-500 rounded-lg shadow-xl z-50"
      >
        <div class="flex items-center gap-3">
          <span class="text-2xl">✅</span>
          <div>
            <p class="font-bold text-green-400">Purchase Successful!</p>
            <p class="text-sm text-green-300">{{ successMessage }}</p>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useShopStore } from '@/features/shop/store';
import { ShopItem } from '@/features/shop/api';
import GoldDisplay from '@/features/shop/components/GoldDisplay.vue';
import GoldPackageCard from '@/features/shop/components/GoldPackageCard.vue';
import BundleCard from '@/features/shop/components/BundleCard.vue';
import PurchaseModal from '@/features/shop/components/PurchaseModal.vue';

const SHOP_ICON = 'https://vibemedia.space/icon_shop_storefront_5m6n7o.png?prompt=medieval%20shop%20storefront%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON';

const shop = useShopStore();

const tabs = [
  { id: 'gold', name: 'Gold Packages' },
  { id: 'bundles', name: 'Equipment Bundles' },
];

const activeTab = ref<'gold' | 'bundles'>('gold');
const showHistory = ref(false);

const purchaseModalVisible = ref(false);
const selectedItem = ref<ShopItem | null>(null);
const successMessage = ref<string | null>(null);

onMounted(async () => {
  await shop.loadShopData();
});

function openPurchaseModal(item: ShopItem) {
  selectedItem.value = item;
  purchaseModalVisible.value = true;
}

function closePurchaseModal() {
  purchaseModalVisible.value = false;
  selectedItem.value = null;
}

async function confirmPurchase(item: ShopItem) {
  const success = await shop.purchase(item.id);
  if (success) {
    closePurchaseModal();
    successMessage.value = `You purchased ${item.name}!`;
    setTimeout(() => {
      successMessage.value = null;
    }, 3000);
    // Refresh data
    await shop.loadShopData();
  }
}

function formatDate(dateStr: string): string {
  const date = new Date(dateStr);
  return date.toLocaleDateString();
}
</script>

<style scoped>
.pixelated {
  image-rendering: pixelated;
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

.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
</style>
