<script setup lang="ts">
import { useAuthStore } from '@/stores/useAuthStore'
import { useAuthAPI } from '@/composables/useAuthAPI'
import { useRoute } from 'vue-router'

const auth = useAuthStore()
const { logout } = useAuthAPI()
const route = useRoute()

const navItems = [
  { label: 'Forja',      icon: 'auto_fix_high',  to: '/forja',     wip: false, requiresAuth: false },
  { label: 'Biblioteca', icon: 'library_books', to: '/biblioteca', wip: false, requiresAuth: true  },
  { label: 'Crónicas',   icon: 'auto_stories',  to: null,  wip: true,  requiresAuth: false },
  { label: 'Bestiario',  icon: 'pets',           to: null,  wip: true,  requiresAuth: false },
  { label: 'Cartografía',icon: 'map',            to: null,  wip: true,  requiresAuth: false },
]

function isActive(to: string): boolean {
  return route.path === to || route.path.startsWith(to + '/')
}
</script>

<template>
  <aside class="fixed left-0 top-0 h-full flex-col pt-20 pb-8 bg-surface-container-low w-64 z-40 hidden lg:flex">
    <!-- Nav -->
    <nav class="flex-1">
      <ul class="space-y-1">
        <li v-for="item in navItems" :key="item.label">
          <!-- Normal: link accesible y no-WIP -->
          <RouterLink
            v-if="item.to && !item.wip && (!item.requiresAuth || auth.isAuthenticated)"
            :to="item.to"
            class="flex items-center gap-3 pl-5 py-3 border-l-4 transition-colors"
            :class="isActive(item.to)
              ? 'border-primary bg-surface-container-lowest text-primary'
              : 'border-transparent text-on-surface-variant hover:bg-surface-container hover:text-on-surface'"
          >
            <span class="material-symbols-outlined">{{ item.icon }}</span>
            <span class="font-label uppercase tracking-widest text-xs font-semibold">{{ item.label }}</span>
          </RouterLink>
          <!-- Requiere auth y no está autenticado: link a /auth con candado -->
          <RouterLink
            v-else-if="item.requiresAuth && !auth.isAuthenticated && !item.wip"
            to="/auth"
            class="flex items-center gap-3 pl-5 py-3 border-l-4 border-transparent text-on-surface-variant hover:bg-surface-container hover:text-on-surface transition-colors"
          >
            <span class="material-symbols-outlined">{{ item.icon }}</span>
            <span class="font-label uppercase tracking-widest text-xs font-semibold">{{ item.label }}</span>
            <span class="material-symbols-outlined text-sm ml-auto mr-4 opacity-50">lock</span>
          </RouterLink>
          <!-- WIP: deshabilitado -->
          <div
            v-else
            class="flex items-center gap-3 pl-5 py-3 text-outline opacity-50 cursor-not-allowed"
          >
            <span class="material-symbols-outlined">{{ item.icon }}</span>
            <span class="font-label uppercase tracking-widest text-xs font-semibold">{{ item.label }}</span>
            <span class="text-[9px] font-bold uppercase tracking-widest text-secondary ml-auto mr-4">WIP</span>
          </div>
        </li>
      </ul>
    </nav>

    <!-- Bottom actions -->
    <div class="px-5 mt-auto border-t border-outline-variant/20 pt-4 space-y-3">
      <div class="flex items-center gap-3 text-outline text-xs font-semibold uppercase tracking-widest px-1 cursor-not-allowed opacity-50">
        <span class="material-symbols-outlined text-sm">help_center</span>
        <span>Soporte</span>
      </div>
      <button
        v-if="auth.isAuthenticated"
        @click="logout()"
        class="flex items-center gap-3 text-secondary text-xs font-label font-bold uppercase tracking-widest px-1 hover:text-primary transition-colors"
      >
        <span class="material-symbols-outlined text-sm">logout</span>
        <span>Salir</span>
      </button>
    </div>
  </aside>
</template>
