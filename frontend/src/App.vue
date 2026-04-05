<script setup lang="ts">
import { onMounted } from 'vue'
import TheTopBar from '@/components/TheTopBar.vue'
import TheSideNav from '@/components/TheSideNav.vue'
import TheBottomNav from '@/components/TheBottomNav.vue'
import { useAuthAPI } from '@/composables/useAuthAPI'
import { useLayoutStore } from '@/stores/useLayoutStore'

const { restoreSession } = useAuthAPI()
const layout = useLayoutStore()
onMounted(restoreSession)
</script>

<template>
  <TheTopBar />
  <TheSideNav />
  <TheBottomNav />
  <main
    class="pt-20 pb-14 lg:pb-0 min-h-screen transition-[padding-left] duration-200 ease-in-out"
    :class="layout.sidebarPinned ? 'lg:pl-64' : 'lg:pl-14'"
  >
    <router-view />
  </main>
</template>
