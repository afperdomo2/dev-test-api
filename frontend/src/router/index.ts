import { createRouter, createWebHistory } from 'vue-router'
import type { Router } from 'vue-router'
import { routes } from './routes'
import { useAuthStore } from '@/stores/auth.store'

declare module 'vue-router' {
  interface RouteMeta {
    requiresAuth?: boolean
    requiresAdmin?: boolean
    layout?: string
  }
}

const router: Router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

router.beforeEach((to, _from) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }

  if (to.meta.requiresAdmin && !authStore.isAdmin) {
    return { path: '/' }
  }

  if (!to.meta.requiresAuth && authStore.isLoggedIn) {
    if (to.name === 'Login' || to.name === 'Setup') {
      return { path: '/' }
    }
  }
})

export default router
