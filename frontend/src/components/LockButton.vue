<script setup lang="ts">
const props = defineProps<{ locked: boolean }>()
const emit = defineEmits<{ toggle: [] }>()
</script>

<template>
  <button
    type="button"
    @click="emit('toggle')"
    :class="props.locked ? 'text-primary' : 'text-outline'"
    class="hover:text-primary transition-colors"
    :title="props.locked ? 'Desbloquear' : 'Bloquear'"
  >
    <span
      class="material-symbols-outlined text-lg"
      :style="props.locked ? `font-variation-settings: 'FILL' 1` : `font-variation-settings: 'FILL' 0`"
    >
      {{ props.locked ? 'lock' : 'lock_open' }}
    </span>
  </button>
</template>
