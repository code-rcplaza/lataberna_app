<script setup lang="ts">
import { onMounted } from 'vue'
import TheTopBar from '@/components/TheTopBar.vue'
import TheSideNav from '@/components/TheSideNav.vue'
import { useAuthAPI } from '@/composables/useAuthAPI'

const { restoreSession } = useAuthAPI()
onMounted(restoreSession)
</script>

<template>
  <TheTopBar />
  <TheSideNav />
  <main class="lg:pl-64 pt-24 min-h-screen">
    <router-view />
  </main>
</template>
