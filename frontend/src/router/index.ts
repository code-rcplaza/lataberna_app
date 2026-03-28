import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/useAuthStore'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/HomeView.vue'),
    },
    {
      path: '/forja',
      name: 'forja',
      component: () => import('@/views/GeneratorView.vue'),
    },
    {
      path: '/auth',
      name: 'auth',
      component: () => import('@/views/AuthView.vue'),
      meta: { public: true },
    },
    {
      path: '/auth/verify',
      name: 'auth-verify',
      component: () => import('@/views/VerifyView.vue'),
      meta: { public: true },
    },
    {
      path: '/biblioteca',
      name: 'biblioteca',
      component: () => import('@/views/BibliotecaView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/biblioteca/:id',
      name: 'biblioteca-detail',
      component: () => import('@/views/CharacterDetailView.vue'),
      meta: { requiresAuth: true },
    },
  ],
})

// Guard — solo bloquea rutas con meta.requiresAuth
// La generación (/) es pública. El guard protege rutas futuras como /biblioteca.
router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return '/auth'
  }
})

export default router
