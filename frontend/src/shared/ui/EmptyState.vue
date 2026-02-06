<template>
  <div class="empty-state">
    <div class="empty-icon">{{ icon }}</div>
    <h3 class="empty-title">{{ title }}</h3>
    <p class="empty-message">{{ message }}</p>
    <slot></slot>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  icon?: string
  title: string
  message: string
}>()
</script>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  text-align: center;
  background: rgba(255, 255, 255, 0.02);
  border-radius: 16px;
  border: 2px dashed rgba(255, 255, 255, 0.05);
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
  opacity: 0.5;
}

.empty-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: white;
  margin: 0 0 0.5rem 0;
}

.empty-message {
  color: #9ca3af;
  max-width: 300px;
  line-height: 1.5;
  margin: 0;
}
</style>
