<script setup lang="ts">
import { ref } from 'vue'
import { useAuthAPI } from '@/composables/useAuthAPI'

const { requestMagicLink } = useAuthAPI()

const email = ref('')
const state = ref<'idle' | 'loading' | 'sent' | 'error'>('idle')
const errorMsg = ref('')

async function submit() {
  if (!email.value.trim()) return
  state.value = 'loading'
  try {
    await requestMagicLink(email.value.trim())
    state.value = 'sent'
  } catch (e: unknown) {
    errorMsg.value = e instanceof Error ? e.message : 'Error al enviar el enlace'
    state.value = 'error'
  }
}

function reset() {
  state.value = 'idle'
  email.value = ''
  errorMsg.value = ''
}
</script>

<template>
  <div class="min-h-screen bg-background flex items-center justify-center px-4">
    <div class="w-full max-w-md space-y-8">

      <!-- Brand -->
      <span class="font-headline italic text-primary text-2xl mb-8 block">La Taberna RPG</span>

      <!-- Estado: idle / loading / error — formulario de email -->
      <template v-if="state !== 'sent'">
        <h1 class="font-headline text-4xl font-bold text-on-surface leading-tight">
          Accedé con tu correo
        </h1>

        <form @submit.prevent="submit" class="space-y-4">
          <input
            v-model="email"
            type="email"
            placeholder="tu@correo.com"
            autocomplete="email"
            :disabled="state === 'loading'"
            class="w-full border border-outline-variant/50 bg-surface-container-lowest px-4 py-3 text-sm font-body focus:outline-none focus:border-primary transition-colors disabled:opacity-60"
          />

          <p v-if="state === 'error'" class="text-error text-sm font-body">
            {{ errorMsg }}
          </p>

          <button
            type="submit"
            :disabled="state === 'loading'"
            class="w-full bg-primary text-on-primary font-label font-bold uppercase tracking-widest py-4 text-sm hover:bg-primary-container transition-colors disabled:opacity-60"
          >
            <span v-if="state === 'loading'">Enviando…</span>
            <span v-else>Recibí mi enlace</span>
          </button>
        </form>

        <blockquote class="font-headline italic text-secondary text-sm leading-relaxed mt-8 border-l-2 border-primary-container/30 pl-4">
          "Las puertas del archivo se abren solo para quienes saben llamar."
        </blockquote>
      </template>

      <!-- Estado: sent — confirmación -->
      <template v-else>
        <h1 class="font-headline text-4xl font-bold text-on-surface leading-tight">
          Revisá tu correo
        </h1>

        <p class="font-body text-on-surface-variant text-sm leading-relaxed">
          Enviamos un enlace a<br />
          <strong class="text-on-surface">{{ email }}</strong>
        </p>

        <p class="font-body italic text-secondary text-xs leading-relaxed">
          En desarrollo, el enlace aparece en la consola del servidor.
        </p>

        <button
          @click="reset()"
          class="text-sm font-label uppercase tracking-widest text-primary hover:underline transition-colors"
        >
          Intentar con otro correo
        </button>
      </template>

    </div>
  </div>
</template>
