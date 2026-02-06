<template>
  <div class="skeleton" :class="{ circle }" :style="style"></div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  width?: string
  height?: string
  circle?: boolean
  borderRadius?: string
}>()

const style = computed(() => ({
  width: props.width || '100%',
  height: props.height || '1rem',
  borderRadius: props.borderRadius || (props.circle ? '50%' : '4px')
}))
</script>

<style scoped>
.skeleton {
  background: linear-gradient(
    90deg,
    rgba(255, 255, 255, 0.05) 25%,
    rgba(255, 255, 255, 0.1) 50%,
    rgba(255, 255, 255, 0.05) 75%
  );
  background-size: 200% 100%;
  animation: loading 1.5s infinite;
}

@keyframes loading {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}
</style>
