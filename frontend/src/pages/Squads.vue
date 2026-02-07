<template>
  <div class="squad-page">
    <div v-if="loading" class="loading-container">
      <div class="loading-spinner"></div>
      <p>Loading Squad...</p>
    </div>

    <div v-else-if="error" class="error-container">
      <p>{{ error }}</p>
      <button @click="loadSquad" class="ep-button ep-button-primary">Retry</button>
    </div>

    <div v-else-if="squad" class="squad-content">
      <!-- View Mode -->
      <section v-if="squad.isActive" class="squad-view">
        <h2 class="ep-header-gold">{{ squad.name }}</h2>
        <p class="squad-status">Active Squad</p>

        <div class="squad-slots-display">
          <div
            v-for="(member, index) in squad.members"
            :key="member.fighterId"
            class="slot-display"
          >
            <SquadSlot
              :fighter="getFighterById(member.fighterId)"
              :slot-index="index"
              :is-active="true"
            />
          </div>
        </div>

        <div class="squad-actions-view">
          <button @click="openEdit" class="ep-button ep-button-secondary">
            Edit Squad
          </button>
        </div>
      </section>

      <!-- Edit Mode -->
      <div v-else class="squad-edit-mode">
        <SquadManagement />
      </div>
    </div>

    <!-- No Squad -->
    <div v-else class="no-squad">
      <h2 class="ep-header-gold">No Active Squad</h2>
      <p class="no-squad-text">
        Create a squad to build your combat team
      </p>
      <button @click="openEdit" class="ep-button ep-button-primary">
        Create Squad
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useSquadStore } from '@/features/squads/store';
import { useFighterStore } from '@/features/roster/store';
import SquadSlot from '@/features/squads/SquadSlot.vue';
import SquadManagement from '@/features/squads/SquadManagement.vue';

const router = useRouter();
const squadStore = useSquadStore();
const fighterStore = useFighterStore();

const loading = ref(true);
const error = ref<string | null>(null);
const squad = ref<any>(null);

function getFighterById(fighterId: string) {
  return fighterStore.fighters.find((f) => f.id === fighterId);
}

function openEdit() {
  router.push('/squads/edit');
}

async function loadSquad() {
  loading.value = true;
  error.value = null;
  try {
    const token = localStorage.getItem('token');
    if (token) {
      await squadStore.loadSquad(token);
      squad.value = squadStore.squad;
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load squad';
  } finally {
    loading.value = false;
  }
}

loadSquad();
</script>

<style scoped>
.squad-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #0f172a, #1e293b);
}

.loading-container,
.error-container,
.no-squad {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  text-align: center;
  padding: 20px;
}

.loading-spinner {
  width: 50px;
  height: 50px;
  border: 3px solid rgba(251, 191, 36, 0.3);
  border-top-color: #fbbf24;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.loading-container p {
  margin-top: 20px;
  color: #94a3b8;
  font-size: 14px;
  text-transform: uppercase;
  letter-spacing: 0.1em;
}

.error-container p {
  color: #ef4444;
  margin-bottom: 20px;
}

.no-squad {
  background: rgba(15, 23, 42, 0.9);
  border: 2px solid rgba(251, 191, 36, 0.3);
  border-radius: 16px;
  margin: 40px auto;
  max-width: 600px;
}

.no-squad-text {
  color: #64748b;
  font-size: 14px;
  margin: 20px 0;
}

.squad-content {
  padding: 40px 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.squad-view {
  background: rgba(15, 23, 42, 0.8);
  border: 2px solid rgba(251, 191, 36, 0.3);
  border-radius: 16px;
  padding: 40px;
  margin-bottom: 40px;
}

.squad-view h2 {
  font-size: 32px;
  font-weight: bold;
  color: #fbbf24;
  text-shadow: 0 0 20px rgba(251, 191, 36, 0.3);
  margin-bottom: 8px;
}

.squad-status {
  color: #10b981;
  font-size: 12px;
  font-weight: bold;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  margin-bottom: 30px;
}

.squad-slots-display {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  margin-bottom: 30px;
}

.squad-actions-view {
  display: flex;
  justify-content: flex-end;
}

.squad-edit-mode {
  background: rgba(15, 23, 42, 0.8);
  border: 2px solid rgba(251, 191, 36, 0.3);
  border-radius: 16px;
  padding: 40px;
  max-height: 80vh;
  overflow-y: auto;
}

.ep-button {
  padding: 12px 24px;
  font-size: 12px;
  font-weight: bold;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: none;
}

.ep-button-secondary {
  background: rgba(30, 41, 59, 0.9);
  color: #e2e8f0;
  border: 2px solid rgba(251, 191, 36, 0.3);
}

.ep-button-secondary:hover {
  background: rgba(251, 191, 36, 0.2);
  border-color: #fbbf24;
}

.ep-button-primary {
  background: linear-gradient(135deg, #f59e0b, #d97706);
  color: black;
  font-weight: bold;
}

.ep-button-primary:hover {
  background: linear-gradient(135deg, #d97706, #b45309);
  box-shadow: 0 0 20px rgba(245, 158, 11, 0.4);
}

@media (max-width: 768px) {
  .squad-slots-display {
    grid-template-columns: 1fr;
  }
}
</style>
