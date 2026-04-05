<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/useAuthStore'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const navItems = [
  { label: 'Forja',      icon: 'auto_fix_high', to: '/forja',      requiresAuth: false },
  { label: 'Biblioteca', icon: 'library_books', to: '/biblioteca',  requiresAuth: true  },
]

function isActive(to: string): boolean {
  return route.path === to || route.path.startsWith(to + '/')
}

function navigate(item: typeof navItems[number]) {
  if (item.requiresAuth && !auth.isAuthenticated) {
    router.push('/auth')
  } else {
    router.push(item.to)
  }
}

const authLabel = computed(() => auth.isAuthenticated ? 'Cuenta' : 'Ingresar')
const authIcon  = computed(() => auth.isAuthenticated ? 'person'  : 'login')
const isAuthActive = computed(() => route.path.startsWith('/auth'))
</script>

<template>
  <nav class="fixed bottom-0 left-0 right-0 z-50 lg:hidden bg-surface-container-low border-t border-outline-variant/20 flex items-stretch h-14">
    <button
      v-for="item in navItems"
      :key="item.label"
      @click="navigate(item)"
      class="flex-1 flex flex-col items-center justify-center gap-0.5 text-[10px] font-label font-bold uppercase tracking-widest transition-colors"
      :class="isActive(item.to)
        ? 'text-primary'
        : 'text-on-surface-variant hover:text-on-surface'"
    >
      <span class="material-symbols-outlined text-xl leading-none">{{ item.icon }}</span>
      <span>{{ item.label }}</span>
    </button>

    <!-- Auth -->
    <RouterLink
      to="/auth"
      class="flex-1 flex flex-col items-center justify-center gap-0.5 text-[10px] font-label font-bold uppercase tracking-widest transition-colors"
      :class="isAuthActive
        ? 'text-primary'
        : 'text-on-surface-variant hover:text-on-surface'"
    >
      <span class="material-symbols-outlined text-xl leading-none">{{ authIcon }}</span>
      <span>{{ authLabel }}</span>
    </RouterLink>
  </nav>
</template>
