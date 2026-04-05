<script setup lang="ts">
defineProps<{ class?: string }>()
</script>
<template>
  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" :class="$props.class">
    <circle cx="4" cy="17" r="2" />
    <circle cx="20" cy="7" r="2" />
    <path d="M6 17 C 12 17 12 7 18 7" />
  </svg>
</template>
