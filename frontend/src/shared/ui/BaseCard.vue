<template>
  <div :class="['rounded-xl border border-slate-800 bg-slate-900/50 p-6 shadow-sm backdrop-blur-sm', className]">
    <div v-if="$slots.header" class="mb-4 border-b border-slate-800 pb-4">
      <slot name="header" />
    </div>
    <slot />
    <div v-if="$slots.footer" class="mt-4 border-t border-slate-800 pt-4">
      <slot name="footer" />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  className?: string;
}>();
</script>
