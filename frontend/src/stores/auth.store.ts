import type { User } from '@/types/user.types'
import { getToken, removeToken, setToken } from '@/utils/storage'
import { getProfile } from '@/api/services/users.service'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(getToken())
  const user = ref<User | null>(null)
  const loading = ref(false)

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.is_admin ?? false)

  function setSession(newToken: string, newUser: User) {
    token.value = newToken
    user.value = newUser
    setToken(newToken)
  }

  function clearSession() {
    token.value = null
    user.value = null
    removeToken()
  }

  function setUser(newUser: User) {
    user.value = newUser
  }

  async function initSession(): Promise<void> {
    if (!token.value || user.value) return

    loading.value = true
    try {
      const profile = await getProfile()
      user.value = profile
    } catch {
      clearSession()
    } finally {
      loading.value = false
    }
  }

  return {
    token,
    user,
    loading,
    isLoggedIn,
    isAdmin,
    setSession,
    clearSession,
    setUser,
    initSession,
  }
})
