<script setup lang="ts">
import { useAuthStore } from '@/stores/useAuthStore'
import { useAuthAPI } from '@/composables/useAuthAPI'

const auth = useAuthStore()
const { logout } = useAuthAPI()
</script>

<template>
  <header class="fixed top-0 w-full z-50 flex items-center justify-between px-8 py-3 bg-background shadow-[0_24px_24px_0_rgba(65,0,2,0.06)]">

    <!-- Brand -->
    <span class="text-primary font-headline italic text-2xl shrink-0">La Taberna RPG</span>

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
