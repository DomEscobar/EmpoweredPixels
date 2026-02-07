<template>
  <div class="constellation-container">
    <div class="constellation-header">
      <h1 class="title">Mastery Constellation</h1>
      <div class="subtitle">Assign Soul-Shards to unlock your fighter's true potential.</div>
    </div>
    
    <div class="constellation-grid">
      <!-- Background / Map Container -->
      <div class="map" :style="{ backgroundImage: `url(${bgUrl})` }">
        
        <!-- Weapon Branch -->
        <div class="branch physical">
          <div class="branch-label text-blue-400">Path of the Blade</div>
          <div class="nodes-row">
            <div class="node-wrapper" v-for="(node, index) in weaponNodes" :key="node.id">
              <div class="node" 
                :class="{ 
                  'node-end': node.isEnd, 
                  'node-unlocked': isUnlocked(node.id),
                  'node-locked': !isUnlocked(node.id) && !canUnlock(node.id, weaponNodes)
                }" 
                @click="unlockNode(node.id, weaponNodes)"
                :data-testid="'node-' + node.id">
                <img :src="swordIcon" :alt="node.title" />
                <div class="node-tooltip">
                  <strong>{{ node.title }}</strong><br/>
                  <span>{{ node.desc }}</span>
                  <div v-if="!isUnlocked(node.id)" class="text-xs mt-1" :class="canUnlock(node.id, weaponNodes) ? 'text-green-400' : 'text-red-400'">
                    {{ canUnlock(node.id, weaponNodes) ? 'Click to Unlock (1 Shard)' : 'Previous Node Required' }}
                  </div>
                </div>
              </div>
              <div v-if="index < weaponNodes.length - 1" class="line horizontal" :class="{ 'line-active': isUnlocked(weaponNodes[index+1].id) }"></div>
            </div>
          </div>
        </div>

        <!-- Support Branch -->
        <div class="branch support">
          <div class="branch-label text-emerald-400">Path of the Soul</div>
          <div class="nodes-row">
            <div class="node-wrapper" v-for="(node, index) in supportNodes" :key="node.id">
              <div class="node" 
                :class="{ 
                  'node-end': node.isEnd, 
                  'node-unlocked': isUnlocked(node.id),
                  'node-locked': !isUnlocked(node.id) && !canUnlock(node.id, supportNodes)
                }" 
                @click="unlockNode(node.id, supportNodes)"
                :data-testid="'node-' + node.id">
                <img :src="heartIcon" :alt="node.title" />
                <div class="node-tooltip">
                  <strong>{{ node.title }}</strong><br/>
                  <span>{{ node.desc }}</span>
                  <div v-if="!isUnlocked(node.id)" class="text-xs mt-1" :class="canUnlock(node.id, supportNodes) ? 'text-green-400' : 'text-red-400'">
                    {{ canUnlock(node.id, supportNodes) ? 'Click to Unlock (1 Shard)' : 'Previous Node Required' }}
                  </div>
                </div>
              </div>
              <div v-if="index < supportNodes.length - 1" class="line horizontal" :class="{ 'line-active': isUnlocked(supportNodes[index+1].id) }"></div>
            </div>
          </div>
        </div>

        <!-- Utility Branch -->
        <div class="branch utility">
          <div class="branch-label text-purple-400">Path of the Forge</div>
          <div class="nodes-row">
            <div class="node-wrapper" v-for="(node, index) in utilityNodes" :key="node.id">
              <div class="node" 
                :class="{ 
                  'node-end': node.isEnd, 
                  'node-unlocked': isUnlocked(node.id),
                  'node-locked': !isUnlocked(node.id) && !canUnlock(node.id, utilityNodes)
                }" 
                @click="unlockNode(node.id, utilityNodes)"
                :data-testid="'node-' + node.id">
                <img :src="cogIcon" :alt="node.title" />
                <div class="node-tooltip">
                  <strong>{{ node.title }}</strong><br/>
                  <span>{{ node.desc }}</span>
                  <div v-if="!isUnlocked(node.id)" class="text-xs mt-1" :class="canUnlock(node.id, utilityNodes) ? 'text-green-400' : 'text-red-400'">
                    {{ canUnlock(node.id, utilityNodes) ? 'Click to Unlock (1 Shard)' : 'Previous Node Required' }}
                  </div>
                </div>
              </div>
              <div v-if="index < utilityNodes.length - 1" class="line horizontal" :class="{ 'line-active': isUnlocked(utilityNodes[index+1].id) }"></div>
            </div>
          </div>
        </div>

      </div>
    </div>

    <div class="constellation-footer">
      <div class="shards-count" data-testid="shards-count">Available Soul-Shards: <span>{{ currentShards }}</span></div>
      <button class="reset-btn" @click="reset">Reset Constellation</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

const props = defineProps<{
  initialShards?: number;
}>();

const emit = defineEmits(['unlock', 'reset']);

const unlockedNodes = ref<string[]>([]);
const spentShards = computed(() => unlockedNodes.value.length);
const currentShards = computed(() => (props.initialShards ?? 5) - spentShards.value);

const bgUrl = "https://vibemedia.space/bg_mastery_constellation_e9a8b7c6.png?prompt=dark%20iron%20texture%20background%20with%20faint%20ethereal%20blue%20spirit%20smoke%20and%20stars&style=pixel_game_asset&key=NOGON";
const swordIcon = "https://vibemedia.space/icon_mastery_sword_7d6c5b4a.png?prompt=glowing%20blue%20sword%20icon%20in%20iron%20frame%20skill%20node&style=pixel_game_asset&key=NOGON";
const heartIcon = "https://vibemedia.space/icon_mastery_heart_3b2a1c0d.png?prompt=glowing%20white%20heart%20icon%20in%20iron%20frame%20skill%20node&style=pixel_game_asset&key=NOGON";
const cogIcon = "https://vibemedia.space/icon_mastery_cog_9f8e7d6c.png?prompt=glowing%20purple%20cog%20icon%20in%20iron%20frame%20skill%20node&style=pixel_game_asset&key=NOGON";

const weaponNodes = [
  { id: 'w1', title: 'Iron Grip', desc: '+10% Attack Speed', isEnd: false },
  { id: 'w2', title: 'Severing Strike', desc: 'Crits apply Bleed (5% HP/sec)', isEnd: false },
  { id: 'w3', title: 'Blood-Iron Hurricane', desc: 'Every 4th attack hits all adjacent enemies for 200% damage', isEnd: true },
];

const supportNodes = [
  { id: 's1', title: 'Ethereal Mend', desc: 'Basic attacks heal lowest HP ally for 5%', isEnd: false },
  { id: 's2', title: 'Spirit Shield', desc: 'Buffing an ally gives them a shield for 10% Max HP', isEnd: false },
  { id: 's3', title: 'Ascendance Pillar', desc: 'On death, remain as invulnerable spirit for 10s, healing allies for 15% HP/sec', isEnd: true },
];

const utilityNodes = [
  { id: 'u1', title: 'Heavy Impact', desc: '10% chance to Stun on hit for 1s', isEnd: false },
  { id: 'u2', title: 'Forge Smoke', desc: 'HP below 30% creates cloud (+50% Dodge)', isEnd: false },
  { id: 'u3', title: 'Molten Lockdown', desc: 'Activating a skill roots enemies for 2s and burns for 3s', isEnd: true },
];

const isUnlocked = (id: string) => unlockedNodes.value.includes(id);

const canUnlock = (nodeId: string, branch: any[]) => {
  if (isUnlocked(nodeId)) return false;
  if (currentShards.value <= 0) return false;
  
  const index = branch.findIndex(n => n.id === nodeId);
  if (index === 0) return true; // First node always unlockable
  
  return isUnlocked(branch[index - 1].id); // Must unlock previous
};

const unlockNode = (nodeId: string, branch: any[]) => {
  if (canUnlock(nodeId, branch)) {
    unlockedNodes.value.push(nodeId);
    emit('unlock', nodeId);
  }
};

const reset = () => {
  unlockedNodes.value = [];
  emit('reset');
};
</script>

<style scoped>
.constellation-container {
  background: #0d0d11;
  color: #e0e0e0;
  padding: 2rem;
  min-height: 100vh;
  font-family: 'Cinzel', serif;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.constellation-header {
  text-align: center;
}

.title {
  text-transform: uppercase;
  letter-spacing: 6px;
  color: #c0c0c0;
  text-shadow: 0 0 15px rgba(255, 255, 255, 0.2);
  margin-bottom: 0.5rem;
  font-size: 2.5rem;
}

.subtitle {
  color: #7a7a8a;
  font-style: italic;
  font-size: 1.1rem;
}

.map {
  width: 100%;
  min-height: 700px;
  background-size: cover;
  background-position: center;
  border: 2px solid #3a3a42;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  justify-content: space-around;
  padding: 4rem;
  position: relative;
  box-shadow: inset 0 0 100px rgba(0,0,0,0.9);
}

.branch {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.branch-label {
  text-transform: uppercase;
  letter-spacing: 2px;
  font-size: 0.9rem;
  opacity: 0.8;
}

.nodes-row {
  display: flex;
  align-items: center;
}

.node-wrapper {
  display: flex;
  align-items: center;
}

.node {
  width: 80px;
  height: 80px;
  background: radial-gradient(circle, #2a2a32 0%, #1a1a22 100%);
  border: 3px solid #3a3a42;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 0 20px rgba(0,0,0,0.5);
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

.node:hover {
  transform: scale(1.15) rotate(5deg);
  border-color: #4a9eff;
  box-shadow: 0 0 25px rgba(74, 158, 255, 0.4);
}

.node-unlocked {
  border-color: #4a9eff;
  background: radial-gradient(circle, #304a60 0%, #1a1a22 100%);
  box-shadow: 0 0 30px rgba(74, 158, 255, 0.6);
}

.node-locked {
  opacity: 0.5;
  filter: grayscale(1);
  cursor: not-allowed;
}

.node-end {
  border-color: #a0a040;
  width: 100px;
  height: 100px;
}

.node-end.node-unlocked {
  border-color: #ffd700;
  box-shadow: 0 0 40px rgba(255, 215, 0, 0.5);
}

.node img {
  width: 70%;
  height: 70%;
  image-rendering: pixelated;
  opacity: 0.9;
}

.node-tooltip {
  position: absolute;
  top: 110%;
  left: 50%;
  transform: translateX(-50%);
  background: #1a1a22e0;
  border: 1px solid #4a9eff;
  padding: 0.8rem;
  width: 200px;
  font-size: 0.85rem;
  z-index: 10;
  pointer-events: none;
  opacity: 0;
  transition: opacity 0.2s;
  box-shadow: 0 10px 20px rgba(0,0,0,0.5);
  backdrop-filter: blur(4px);
}

.node:hover .node-tooltip {
  opacity: 1;
}

.line {
  height: 6px;
  background: #3a3a42;
  width: 120px;
  opacity: 0.4;
  position: relative;
}

.line-active {
  background: linear-gradient(90deg, #3a3a42, #4a9eff, #3a3a42);
  opacity: 1;
}

.line-active::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  height: 100%;
  width: 30%;
  background: white;
  filter: blur(4px);
  animation: flow 2s infinite linear;
}

@keyframes flow {
  from { left: -30%; }
  to { left: 130%; }
}

.constellation-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  background: #1a1a22;
  border-top: 1px solid #3a3a42;
}

.shards-count {
  font-size: 1.2rem;
}

.shards-count span {
  color: #4a9eff;
  font-weight: bold;
}

.reset-btn {
  background: transparent;
  border: 1px solid #ff4a4a;
  color: #ff4a4a;
  padding: 0.6rem 1.5rem;
  text-transform: uppercase;
  cursor: pointer;
  transition: all 0.2s;
}

.reset-btn:hover {
  background: #ff4a4a;
  color: white;
}
</style>
