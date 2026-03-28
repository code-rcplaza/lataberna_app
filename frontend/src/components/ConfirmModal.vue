<script setup lang="ts">
defineProps<{
  open: boolean
  title: string
  message: string
  confirmLabel?: string
}>()

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') emit('cancel')
}
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="open"
        class="fixed inset-0 z-50 flex items-center justify-center px-6"
        @keydown="onKeydown"
      >
        <!-- Backdrop -->
        <div
          class="absolute inset-0 bg-background/80 backdrop-blur-sm"
          @click="emit('cancel')"
        />

        <!-- Dialog -->
        <div class="relative bg-surface-container w-full max-w-sm shadow-xl flex flex-col">
          <!-- Header -->
          <div class="flex items-center gap-3 px-6 pt-6 pb-4 border-b border-outline-variant/20">
            <span class="material-symbols-outlined text-error text-2xl">warning</span>
            <h2 class="font-headline text-on-surface text-xl font-bold">{{ title }}</h2>
          </div>

          <!-- Body -->
          <p class="px-6 py-5 text-on-surface-variant text-sm leading-relaxed">
            {{ message }}
          </p>

          <!-- Actions -->
          <div class="flex items-center justify-end gap-3 px-6 pb-6">
            <button
              @click="emit('cancel')"
              class="px-4 py-2 border border-outline-variant text-secondary font-label font-bold uppercase tracking-widest text-xs hover:border-outline hover:text-on-surface transition-colors"
            >
              Cancelar
            </button>
            <button
              @click="emit('confirm')"
              class="px-4 py-2 bg-error text-on-error font-label font-bold uppercase tracking-widest text-xs hover:opacity-90 transition-opacity"
            >
              {{ confirmLabel ?? 'Eliminar' }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.15s ease;
}
.modal-enter-active .relative,
.modal-leave-active .relative {
  transition: transform 0.15s ease, opacity 0.15s ease;
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
.modal-enter-from .relative,
.modal-leave-to .relative {
  transform: translateY(8px);
  opacity: 0;
}
</style>
