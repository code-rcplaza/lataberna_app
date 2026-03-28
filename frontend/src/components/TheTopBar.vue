<script setup lang="ts">
import { useAuthStore } from '@/stores/useAuthStore'
import { useAuthAPI } from '@/composables/useAuthAPI'
import { useRoute } from 'vue-router'

const auth = useAuthStore()
const { logout } = useAuthAPI()
const route = useRoute()

function isActive(to: string): boolean {
  return route.path === to || route.path.startsWith(to + '/')
}
</script>

<template>
  <header class="fixed top-0 w-full z-50 flex justify-between items-center px-8 py-4 bg-background shadow-[0_24px_24px_0_rgba(65,0,2,0.06)]">
    <div class="flex items-center gap-8">
      <span class="text-primary font-headline italic text-2xl">La Taberna RPG</span>
      <nav class="hidden md:flex gap-6">
        <RouterLink
          to="/forja"
          class="font-medium transition-colors"
          :class="isActive('/forja') ? 'text-primary font-bold underline underline-offset-8' : 'text-on-surface-variant hover:text-on-surface'"
        >Forja</RouterLink>
        <RouterLink
          v-if="auth.isAuthenticated"
          to="/biblioteca"
          class="font-medium transition-colors"
          :class="isActive('/biblioteca')
            ? 'text-primary font-bold underline underline-offset-8'
            : 'text-on-surface-variant hover:text-on-surface'"
        >Biblioteca</RouterLink>
        <span class="text-outline font-medium cursor-not-allowed opacity-50">Crónicas</span>
        <span class="text-outline font-medium cursor-not-allowed opacity-50">Bestiario</span>
      </nav>
    </div>
    <div class="flex items-center gap-4">
      <button v-if="auth.isAuthenticated" @click="logout()" class="text-xs font-bold uppercase tracking-widest text-secondary hover:text-primary transition-colors">
        Salir
      </button>
    </div>
  </header>
</template>
