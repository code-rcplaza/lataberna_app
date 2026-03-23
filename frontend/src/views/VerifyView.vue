<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthAPI } from '@/composables/useAuthAPI'

const route = useRoute()
const { verifyMagicLink } = useAuthAPI()

const state = ref<'loading' | 'success' | 'error'>('loading')
const errorMsg = ref('')

onMounted(async () => {
  const token = route.query.token as string
  if (!token) {
    state.value = 'error'
    errorMsg.value = 'Token no encontrado en la URL.'
    return
  }
  try {
    await verifyMagicLink(token)
    state.value = 'success'
    // verifyMagicLink redirige a / automáticamente al tener éxito
  } catch (e: unknown) {
    state.value = 'error'
    errorMsg.value = e instanceof Error ? e.message : 'El enlace es inválido o ya fue usado.'
  }
})
</script>

<template>
  <div class="min-h-screen bg-background flex items-center justify-center px-4">
    <div class="w-full max-w-md space-y-8">

      <!-- Brand -->
      <span class="font-headline italic text-primary text-2xl mb-8 block">La Taberna RPG</span>

      <!-- Cargando -->
      <template v-if="state === 'loading'">
        <div class="flex flex-col items-start gap-6">
          <div class="w-8 h-8 border-2 border-primary border-t-transparent rounded-full animate-spin" />
          <p class="font-headline text-2xl text-on-surface">Verificando tu enlace…</p>
        </div>
      </template>

      <!-- Éxito -->
      <template v-else-if="state === 'success'">
        <h1 class="font-headline text-4xl font-bold text-on-surface leading-tight">
          ¡Acceso concedido!
        </h1>
        <p class="font-body text-on-surface-variant text-sm">
          Redirigiendo al archivo…
        </p>
      </template>

      <!-- Error -->
      <template v-else>
        <h1 class="font-headline text-4xl font-bold text-on-surface leading-tight">
          Enlace inválido
        </h1>
        <p class="font-body text-error text-sm leading-relaxed">
          {{ errorMsg }}
        </p>
        <RouterLink
          to="/auth"
          class="inline-block text-sm font-label uppercase tracking-widest text-primary hover:underline transition-colors"
        >
          Volver a intentarlo
        </RouterLink>
      </template>

    </div>
  </div>
</template>
