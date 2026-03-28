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
  <header class="fixed top-0 w-full z-50 flex items-center justify-between px-8 py-3 bg-background shadow-[0_24px_24px_0_rgba(65,0,2,0.06)]">

    <!-- Brand -->
    <span class="text-primary font-headline italic text-2xl shrink-0">La Taberna RPG</span>

    <!-- Nav -->
    <nav class="hidden md:flex items-center gap-6">
      <RouterLink
        to="/forja"
        class="text-sm font-medium transition-colors"
        :class="isActive('/forja') ? 'text-primary font-bold underline underline-offset-8' : 'text-on-surface-variant hover:text-on-surface'"
      >Forja</RouterLink>
      <RouterLink
        v-if="auth.isAuthenticated"
        to="/biblioteca"
        class="text-sm font-medium transition-colors"
        :class="isActive('/biblioteca') ? 'text-primary font-bold underline underline-offset-8' : 'text-on-surface-variant hover:text-on-surface'"
      >Biblioteca</RouterLink>
      <span class="text-sm text-outline font-medium cursor-not-allowed opacity-40">Crónicas</span>
      <span class="text-sm text-outline font-medium cursor-not-allowed opacity-40">Bestiario</span>
    </nav>

    <!-- User block -->
    <div class="flex items-center gap-3 shrink-0">
      <template v-if="auth.isAuthenticated && auth.user">
        <div class="text-right">
          <p class="text-sm font-bold text-on-surface leading-tight">{{ auth.user.email.split('@')[0] }}</p>
          <p class="text-[10px] font-bold uppercase tracking-widest text-secondary leading-tight">Curador</p>
        </div>
        <button
          @click="logout()"
          class="w-9 h-9 rounded-full bg-surface-container flex items-center justify-center hover:bg-surface-container-high transition-colors"
          title="Cerrar sesión"
        >
          <span class="material-symbols-outlined text-secondary text-lg">person</span>
        </button>
      </template>
      <template v-else>
        <RouterLink
          to="/auth"
          class="text-xs font-bold uppercase tracking-widest text-secondary hover:text-primary transition-colors"
        >Iniciar sesión</RouterLink>
      </template>
    </div>

  </header>
</template>
