<template>
  <div class="squad-management overhaul-theme-container">
    <!-- Header Section -->
    <section class="squad-header">
      <div class="header-content">
        <h2 class="ep-header-gold">Squad Management</h2>
        <p class="header-subtitle">
          Select up to 3 fighters to form your combat team
        </p>
      </div>

      <!-- Squad Name -->
      <div class="squad-name-section">
        <label class="ep-header-gold" for="squad-name">Squad Name</label>
        <input
          id="squad-name"
          v-model="squadName"
          type="text"
          maxlength="50"
          class="ep-input"
          placeholder="Enter squad name..."
          :disabled="isSubmitting"
        />
      </div>

      <!-- Fighter Slots -->
      <div class="squad-slots">
        <div
          v-for="(slot, index) in 3"
          :key="slot"
          class="squad-slot-wrapper"
        >
          <div class="slot-header">
            <span class="slot-label">Slot {{ index + 1 }}</span>
            <span v-if="!fighters[slot - 1]" class="slot-status">Empty</span>
            <span v-else class="slot-status filled">
              {{ fighters[slot - 1].name }}
            </span>
          </div>
          <SquadSlot
            v-if="fighters[slot - 1]"
            :fighter="fighters[slot - 1]"
            :slot-index="slot - 1"
            :is-active="selectedSlot === slot - 1"
            class="slot-content"
          />
          <button
            v-else
            @click="selectSlot(slot - 1)"
            class="slot-select-btn"
            :disabled="selectedSlot !== null"
          >
            Select Fighter
          </button>
        </div>
      </div>

      <!-- Selected Fighter Display -->
      <div
        v-if="selectedSlot !== null && !fighters[selectedSlot]"
        class="selected-fighter-panel"
      >
        <h3 class="ep-header-gold">Select a Fighter</h3>
        <div class="fighter-list">
          <button
            v-for="fighter in availableFighters"
            :key="fighter.id"
            @click="addFighterToSlot(fighter)"
            class="fighter-select-btn"
          >
            <span class="fighter-name">{{ fighter.name }}</span>
            <span class="fighter-level">Lv.{{ fighter.level }}</span>
            <span class="fighter-power">⚡ {{ fighter.power }}</span>
          </button>
        </div>
      </div>
    </section>

    <!-- Actions -->
    <section class="squad-actions">
      <div class="actions-left">
        <button
          @click="clearSlot"
          class="ep-button ep-button-secondary"
          :disabled="selectedSlot === null || fighters[selectedSlot]"
        >
          Clear Slot
        </button>
      </div>

      <div class="actions-right">
        <button
          @click="resetAllSlots"
          class="ep-button ep-button-secondary"
          :disabled="membersCount === 0"
        >
          Reset All
        </button>

        <button
          @click="saveSquad"
          class="ep-button ep-button-primary"
          :disabled="isSubmitting || membersCount === 0"
        >
          <span v-if="isSubmitting">Saving...</span>
          <span v-else>{{ membersCount > 0 ? 'Save Squad' : 'Save Squad' }}</span>
        </button>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import SquadSlot from './SquadSlot.vue';
import { useSquadStore } from './store';
import { useFighterStore } from '../roster/store';

const router = useRouter();
const squadStore = useSquadStore();
const fighterStore = useFighterStore();

const squadName = ref('My Squad');
const selectedSlot = ref<number | null>(null);
const isSubmitting = ref(false);

const fighters = ref<any[]>(Array(3).fill(null));

const membersCount = computed(() => fighters.value.filter((f) => f !== null).length);

const availableFighters = computed(() => {
  return fighterStore.fighters.filter(
    (fighter) => !fighters.value.some((existing) => existing?.id === fighter.id)
  );
});

function selectSlot(index: number) {
  selectedSlot.value = index;
}

function addFighterToSlot(fighter: any) {
  if (selectedSlot.value !== null) {
    fighters.value[selectedSlot.value] = fighter;
    selectedSlot.value = null;
  }
}

function clearSlot() {
  if (selectedSlot.value !== null) {
    fighters.value[selectedSlot.value] = null;
    selectedSlot.value = null;
  }
}

function resetAllSlots() {
  fighters.value = Array(3).fill(null);
}

async function saveSquad() {
  isSubmitting.value = true;
  try {
    const fighterIds = fighters.value
      .filter((f) => f !== null)
      .map((f) => f.id);

    await squadStore.updateSquad(localStorage.getItem('token') || '', {
      name: squadName.value,
      fighterIds,
    });

    // Reload roster to show updated squad
    await fighterStore.fetchFighters();
    await squadStore.loadSquad(localStorage.getItem('token') || '');

    alert('Squad saved successfully!');
  } catch (error) {
    alert('Failed to save squad: ' + error);
  } finally {
    isSubmitting.value = false;
  }
}

// Load existing squad on mount
async function loadSquad() {
  try {
    const token = localStorage.getItem('token');
    if (token) {
      await squadStore.loadSquad(token);
      if (squadStore.squad) {
        squadName.value = squadStore.squad.name;
        squadStore.squad.members.forEach((member) => {
          fighters.value[member.slotIndex] = fighterStore.fighters.find(
            (f) => f.id === member.fighterId
          );
        });
      }
    }
  } catch (error) {
    console.error('Failed to load squad:', error);
  }
}

loadSquad();
</script>

<style scoped>
.squad-management {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px;
}

.squad-header {
  margin-bottom: 40px;
}

.header-content {
  text-align: center;
  margin-bottom: 30px;
}

.header-subtitle {
  color: #64748b;
  font-size: 14px;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  margin-top: 8px;
}

.squad-name-section {
  max-width: 600px;
  margin: 0 auto 40px auto;
}

.squad-name-section label {
  display: block;
  font-size: 12px;
  font-weight: bold;
  color: #fbbf24;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  margin-bottom: 8px;
}

.squad-name-section input {
  width: 100%;
  padding: 12px 16px;
  font-size: 16px;
  background: rgba(15, 23, 42, 0.9);
  border: 2px solid rgba(251, 191, 36, 0.3);
  border-radius: 8px;
  color: white;
  text-align: center;
}

.squad-name-section input:focus {
  outline: none;
  border-color: #fbbf24;
}

.squad-slots {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  margin-bottom: 40px;
}

.squad-slot-wrapper {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.slot-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 4px 8px;
  background: rgba(15, 23, 42, 0.6);
  border-radius: 4px;
  border-left: 2px solid rgba(251, 191, 36, 0.3);
}

.slot-label {
  font-size: 10px;
  font-weight: bold;
  color: #64748b;
  text-transform: uppercase;
}

.slot-status {
  font-size: 10px;
  color: #64748b;
}

.slot-status.filled {
  color: #10b981;
  font-weight: bold;
}

.slot-content {
  flex: 1;
}

.slot-select-btn {
  height: 200px;
  background: rgba(15, 23, 42, 0.6);
  border: 2px dashed rgba(251, 191, 36, 0.3);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.slot-select-btn:hover:not(:disabled) {
  border-color: #fbbf24;
  background: rgba(251, 191, 36, 0.1);
}

.slot-select-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.selected-fighter-panel {
  background: rgba(15, 23, 42, 0.8);
  border: 2px solid rgba(251, 191, 36, 0.3);
  border-radius: 12px;
  padding: 24px;
  margin-top: 30px;
}

.selected-fighter-panel h3 {
  font-size: 14px;
  font-weight: bold;
  color: #fbbf24;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  margin-bottom: 16px;
}

.fighter-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  max-height: 400px;
  overflow-y: auto;
}

.fighter-select-btn {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: rgba(30, 41, 59, 0.8);
  border: 2px solid rgba(251, 191, 36, 0.2);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.fighter-select-btn:hover {
  border-color: #fbbf24;
  background: rgba(251, 191, 36, 0.1);
}

.fighter-select-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.fighter-name {
  color: #e2e8f0;
  font-weight: bold;
}

.fighter-level {
  color: #94a3b8;
  font-size: 12px;
}

.fighter-power {
  color: #fbbf24;
  font-weight: bold;
}

.squad-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 0;
  border-top: 2px solid rgba(251, 191, 36, 0.2);
}

.actions-left,
.actions-right {
  display: flex;
  gap: 12px;
}

.actions-right {
  flex-direction: row-reverse;
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

.ep-button-secondary:hover:not(:disabled) {
  background: rgba(251, 191, 36, 0.2);
  border-color: #fbbf24;
}

.ep-button-primary {
  background: linear-gradient(135deg, #f59e0b, #d97706);
  color: black;
  font-weight: bold;
}

.ep-button-primary:hover:not(:disabled) {
  background: linear-gradient(135deg, #d97706, #b45309);
  box-shadow: 0 0 20px rgba(245, 158, 11, 0.4);
}

.ep-button-primary:disabled,
.ep-button-secondary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
