<template>
  <div class="toast-container">
    <TransitionGroup name="toast">
      <div 
        v-for="toast in toasts" 
        :key="toast.id" 
        :class="['toast', toast.type]"
        @click="removeToast(toast.id)"
      >
        <span class="toast-icon">{{ toast.type === 'success' ? '✅' : '❌' }}</span>
        <span class="toast-message">{{ toast.message }}</span>
      </div>
    </TransitionGroup>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface Toast {
  id: number
  message: string
  type: 'success' | 'error'
}

const toasts = ref<Toast[]>([])
let nextId = 0

function addToast(message: string, type: 'success' | 'error' = 'success', duration = 3000) {
  const id = nextId++
  toasts.value.push({ id, message, type })
  
  if (duration > 0) {
    setTimeout(() => {
      removeToast(id)
    }, duration)
  }
}

function removeToast(id: number) {
  const index = toasts.value.findIndex(t => t.id === id)
  if (index !== -1) {
    toasts.value.splice(index, 1)
  }
}

defineExpose({ addToast })
</script>

<style scoped>
.toast-container {
  position: fixed;
  bottom: 1.5rem;
  right: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  z-index: 1000;
  pointer-events: none;
}

.toast {
  pointer-events: auto;
  padding: 1rem 1.25rem;
  background: #1f2937;
  border-radius: 12px;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  gap: 0.75rem;
  min-width: 280px;
  max-width: 400px;
  border-left: 4px solid #3b82f6;
  cursor: pointer;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.toast.success {
  border-left: 4px solid #10b981;
}

.toast.error {
  border-left: 4px solid #ef4444;
}

.toast-message {
  color: white;
  font-weight: 500;
  font-size: 0.875rem;
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.toast-leave-to {
  opacity: 0;
  transform: scale(0.9);
}
</style>
