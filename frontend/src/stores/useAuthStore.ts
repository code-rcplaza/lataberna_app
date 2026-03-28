import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types/character'
import { useGeneratorHistoryStore } from './useGeneratorHistoryStore'

export const useAuthStore = defineStore('auth', () => {
  const sessionId = ref<string | null>(localStorage.getItem('sessionId'))
  const user = ref<User | null>(null)

  const isAuthenticated = computed(() => !!sessionId.value)

  function login(id: string, u: User) {
    sessionId.value = id
    user.value = u
    localStorage.setItem('sessionId', id)
  }

  function logout() {
    sessionId.value = null
    user.value = null
    localStorage.removeItem('sessionId')
    useGeneratorHistoryStore().clear()
  }

  return { sessionId, user, isAuthenticated, login, logout }
})
