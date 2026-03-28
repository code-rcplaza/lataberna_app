import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const useLayoutStore = defineStore('layout', () => {
  const sidebarPinned = ref(localStorage.getItem('sidebar_pinned') === 'true')

  watch(sidebarPinned, (val) => {
    localStorage.setItem('sidebar_pinned', String(val))
  })

  function togglePin() {
    sidebarPinned.value = !sidebarPinned.value
  }

  return { sidebarPinned, togglePin }
})
