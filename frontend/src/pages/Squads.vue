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
              :slot-index="Number(index)"
              :is-active="true"
            />
          </div>
        </div>

        <div class="squad-actions-view">
          <button
            @click="openEdit"
            class="ep-button ep-button-secondary"
            data-testid="edit-squad-btn"
          >
            Edit Squad
          </button>
        </div>

        <!-- Eligible Leagues Section (New Integration) -->
        <div class="eligible-leagues-section mt-10 pt-10 border-t border-slate-800">
           <div class="flex items-center justify-between mb-6">
              <h3 class="text-amber-500 font-bold uppercase tracking-widest text-lg">Eligible Competitions</h3>
              <span class="text-[10px] text-slate-500 font-mono">BASED ON CURRENT COMPOSITION</span>
           </div>
           
           <div v-if="eligibleLeagues.length > 0" class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div v-for="league in eligibleLeagues" :key="league.id" 
                   class="bg-black/40 border border-slate-800 p-4 rounded-lg flex items-center justify-between group hover:border-amber-500/50 transition-all">
                  <div class="flex items-center gap-4">
                      <div class="w-10 h-10 bg-slate-900 rounded-full flex items-center justify-center text-xl">🏆</div>
                      <div>
                          <h4 class="text-slate-200 font-bold text-sm">{{ league.name }}</h4>
                          <p class="text-[10px] text-slate-500 uppercase">{{ league.options?.tier || 'General' }} • Active</p>
                      </div>
                  </div>
                  <router-link :to="`/leagues?id=${league.id}`" 
                               class="text-[10px] bg-amber-600/10 hover:bg-amber-600 text-amber-500 hover:text-black border border-amber-500/50 px-3 py-1.5 rounded font-black transition-all">
                      REGISTER
                  </router-link>
              </div>
           </div>
           <div v-else class="text-center py-6 bg-slate-900/50 rounded-lg border border-dashed border-slate-800">
              <p class="text-slate-500 text-sm italic">Adjust your squad composition to unlock new league eligibilities.</p>
           </div>
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
      <button
        @click="openEdit"
        class="ep-button ep-button-primary"
        data-testid="create-squad-btn"
      >
        Create Squad
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useSquadStore } from '@/features/squads/store';
import { useFighterStore } from '@/features/roster/store';
import { useLeaguesStore } from '@/features/leagues/store';
import { Fighter } from '@/features/roster/api';
import SquadSlot from '@/features/squads/SquadSlot.vue';
import SquadManagement from '@/features/squads/SquadManagement.vue';

const router = useRouter();
const squadStore = useSquadStore();
const fighterStore = useFighterStore();
const leaguesStore = useLeaguesStore();

const loading = ref(true);
const error = ref<string | null>(null);
const squad = ref<any>(null);

const eligibleLeagues = computed(() => {
  if (!squad.value) return [];
  // For now, simplify and show all major leagues as "eligible" if squad exists
  // Real logic would check squad.power / level etc
  return leaguesStore.leagues;
});

onMounted(() => {
  leaguesStore.fetchLeagues();
});

function getFighterById(fighterId: string) {
  return (fighterStore.fighters as Fighter[]).find((f) => f.id === fighterId);
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
