import { create } from 'zustand'
import type { User } from '@devtest/shared'

import { onUnauthorized } from '@/core/api/client'
import {
  getProfile,
  getStatus,
  login as loginApi,
  setup as setupApi,
} from '@/core/api/services/auth.service'
import { getToken, removeToken, setToken } from '@/core/auth/storage'

type AuthState = {
  token: string | null
  user: User | null
  actionLoading: boolean
  needsSetup: boolean | null
  ready: boolean
  bootstrap: () => Promise<void>
  login: (email: string, password: string) => Promise<void>
  setup: (email: string, password: string) => Promise<void>
  logout: () => Promise<void>
  clearSession: () => Promise<void>
}

export const useAuthStore = create<AuthState>()((set) => ({
  token: null,
  user: null,
  actionLoading: false,
  needsSetup: null,
  ready: false,

  bootstrap: async () => {
    try {
      const status = await getStatus()
      const needsSetup = !status.initialized
      set({ needsSetup })
      if (!needsSetup) {
        const token = await getToken()
        if (token) {
          try {
            const user = await getProfile()
            set({ token, user })
          } catch {
            await removeToken()
            set({ token: null, user: null })
          }
        }
      }
    } catch {
      set({ needsSetup: null })
    } finally {
      set({ ready: true })
    }
  },

  login: async (email, password) => {
    set({ actionLoading: true })
    try {
      const res = await loginApi({ email, password })
      await setToken(res.token)
      set({ token: res.token, user: res.user })
    } finally {
      set({ actionLoading: false })
    }
  },

  setup: async (email, password) => {
    set({ actionLoading: true })
    try {
      await setupApi({ email, password })
      set({ needsSetup: false })
    } finally {
      set({ actionLoading: false })
    }
  },

  logout: async () => {
    await removeToken()
    set({ token: null, user: null })
  },

  clearSession: async () => {
    await removeToken()
    set({ token: null, user: null })
  },
}))

onUnauthorized(() => {
  void useAuthStore.getState().clearSession()
})

export function useIsLoggedIn(): boolean {
  return useAuthStore((s) => s.token !== null)
}

export function useIsAdmin(): boolean {
  return useAuthStore((s) => s.user?.isAdmin ?? false)
}
