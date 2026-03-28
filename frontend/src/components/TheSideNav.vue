<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/useAuthStore'
import { useAuthAPI } from '@/composables/useAuthAPI'
import { useLayoutStore } from '@/stores/useLayoutStore'
import { useRoute } from 'vue-router'

const auth = useAuthStore()
const { logout } = useAuthAPI()
const layout = useLayoutStore()
const route = useRoute()

const hovered = ref(false)
const isOpen = computed(() => layout.sidebarPinned || hovered.value)

const navItems = [
  { label: 'Forja',       icon: 'auto_fix_high', to: '/forja',     wip: false, requiresAuth: false },
  { label: 'Biblioteca',  icon: 'library_books', to: '/biblioteca', wip: false, requiresAuth: true  },
  { label: 'Crónicas',    icon: 'auto_stories',  to: null,          wip: true,  requiresAuth: false },
  { label: 'Bestiario',   icon: 'pets',          to: null,          wip: true,  requiresAuth: false },
  { label: 'Cartografía', icon: 'map',           to: null,          wip: true,  requiresAuth: false },
]

function isActive(to: string): boolean {
  return route.path === to || route.path.startsWith(to + '/')
}
</script>

<template>
  <aside
    class="fixed left-0 top-0 h-full flex-col pt-16 pb-6 bg-surface-container-low z-40 hidden lg:flex overflow-hidden transition-[width] duration-200 ease-in-out"
    :class="isOpen ? 'w-64' : 'w-14'"
    @mouseenter="hovered = true"
    @mouseleave="hovered = false"
  >
    <!-- Nav -->
    <nav class="flex-1">
      <ul class="space-y-1">
        <li v-for="item in navItems" :key="item.label">

          <!-- Activo y accesible -->
          <RouterLink
            v-if="item.to && !item.wip && (!item.requiresAuth || auth.isAuthenticated)"
            :to="item.to"
            :title="!isOpen ? item.label : undefined"
            class="flex items-center gap-3 pl-3 py-3 border-l-4 transition-colors"
            :class="isActive(item.to)
              ? 'border-primary bg-surface-container-lowest text-primary'
              : 'border-transparent text-on-surface-variant hover:bg-surface-container hover:text-on-surface'"
          >
            <span class="material-symbols-outlined shrink-0">{{ item.icon }}</span>
            <span
              class="font-label uppercase tracking-widest text-xs font-semibold whitespace-nowrap transition-opacity duration-150"
              :class="isOpen ? 'opacity-100' : 'opacity-0'"
            >{{ item.label }}</span>
          </RouterLink>

          <!-- Requiere auth -->
          <RouterLink
            v-else-if="item.requiresAuth && !auth.isAuthenticated && !item.wip"
            to="/auth"
            :title="!isOpen ? item.label : undefined"
            class="flex items-center gap-3 pl-3 py-3 border-l-4 border-transparent text-on-surface-variant hover:bg-surface-container hover:text-on-surface transition-colors"
          >
            <span class="material-symbols-outlined shrink-0">{{ item.icon }}</span>
            <span
              class="font-label uppercase tracking-widest text-xs font-semibold whitespace-nowrap transition-opacity duration-150"
              :class="isOpen ? 'opacity-100' : 'opacity-0'"
            >{{ item.label }}</span>
            <span v-if="isOpen" class="material-symbols-outlined text-sm ml-auto mr-4 opacity-50">lock</span>
          </RouterLink>

          <!-- WIP -->
          <div
            v-else
            :title="!isOpen ? item.label : undefined"
            class="flex items-center gap-3 pl-3 py-3 text-outline opacity-40 cursor-not-allowed"
          >
            <span class="material-symbols-outlined shrink-0">{{ item.icon }}</span>
            <span
              class="font-label uppercase tracking-widest text-xs font-semibold whitespace-nowrap transition-opacity duration-150"
              :class="isOpen ? 'opacity-100' : 'opacity-0'"
            >{{ item.label }}</span>
            <span v-if="isOpen" class="text-[9px] font-bold uppercase tracking-widest text-secondary ml-auto mr-4">WIP</span>
          </div>

        </li>
      </ul>
    </nav>

    <!-- Bottom actions -->
    <div class="mt-auto border-t border-outline-variant/20 pt-3 space-y-1">

      <!-- Logout -->
      <button
        v-if="auth.isAuthenticated"
        @click="logout()"
        :title="!isOpen ? 'Salir' : undefined"
        class="flex items-center gap-3 pl-3 py-2 w-full text-secondary text-xs font-label font-bold uppercase tracking-widest hover:text-primary transition-colors"
      >
        <span class="material-symbols-outlined text-sm shrink-0">logout</span>
        <span
          class="whitespace-nowrap transition-opacity duration-150"
          :class="isOpen ? 'opacity-100' : 'opacity-0'"
        >Salir</span>
      </button>

      <!-- Pin toggle -->
      <button
        @click="layout.togglePin()"
        :title="layout.sidebarPinned ? 'Desanclar sidebar' : 'Anclar sidebar'"
        class="flex items-center gap-3 pl-3 py-2 w-full text-outline hover:text-on-surface transition-colors"
      >
        <span class="material-symbols-outlined text-sm shrink-0">
          {{ layout.sidebarPinned ? 'keep' : 'keep_off' }}
        </span>
        <span
          class="text-xs font-label font-semibold uppercase tracking-widest whitespace-nowrap transition-opacity duration-150"
          :class="isOpen ? 'opacity-100' : 'opacity-0'"
        >{{ layout.sidebarPinned ? 'Anclado' : 'Anclar' }}</span>
      </button>

    </div>
  </aside>
</template>
