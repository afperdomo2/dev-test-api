import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from './auth.store'
import * as usersService from '@/api/services/users.service'
import * as authService from '@/api/services/auth.service'

vi.mock('@/api/services/users.service', () => ({
  getProfile: vi.fn(),
}))

vi.mock('@/api/services/auth.service', () => ({
  getStatus: vi.fn(),
}))

vi.mock('@/utils/storage', async () => {
  const actual = await vi.importActual('@/utils/storage')
  return {
    ...(actual as object),
    getToken: vi.fn(() => null),
    setToken: vi.fn(),
    removeToken: vi.fn(),
  }
})

import { getToken, setToken, removeToken } from '@/utils/storage'

const mockUser = {
  id: '1',
  email: 'a@b.com',
  isAdmin: true,
  createdAt: '2026-01-01',
  updatedAt: '2026-01-01',
}

describe('auth.store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    localStorage.clear()
    vi.mocked(getToken).mockReturnValue(null)
  })

  it('initial state', () => {
    const store = useAuthStore()
    expect(store.token).toBeNull()
    expect(store.user).toBeNull()
    expect(store.isLoggedIn).toBe(false)
    expect(store.isAdmin).toBe(false)
    expect(store.needsSetup).toBeNull()
  })

  it('setSession sets token/user and persists', () => {
    const store = useAuthStore()
    store.setSession('tok123', mockUser as never)
    expect(store.token).toBe('tok123')
    expect(store.user).toEqual(mockUser)
    expect(store.isLoggedIn).toBe(true)
    expect(store.isAdmin).toBe(true)
    expect(setToken).toHaveBeenCalledWith('tok123')
  })

  it('clearSession nulls and removes', () => {
    const store = useAuthStore()
    store.setSession('tok', mockUser as never)
    store.clearSession()
    expect(store.token).toBeNull()
    expect(store.user).toBeNull()
    expect(removeToken).toHaveBeenCalled()
  })

  it('setUser updates user', () => {
    const store = useAuthStore()
    store.setUser(mockUser as never)
    expect(store.user).toEqual(mockUser)
    expect(store.isAdmin).toBe(true)
  })

  it('isLoggedIn computed from token', () => {
    const store = useAuthStore()
    expect(store.isLoggedIn).toBe(false)
    store.token = 'abc'
    expect(store.isLoggedIn).toBe(true)
  })

  it('isAdmin false when user not admin', () => {
    const store = useAuthStore()
    store.setUser({ ...mockUser, isAdmin: false } as never)
    expect(store.isAdmin).toBe(false)
  })

  describe('initSession', () => {
    it('returns early if no token', async () => {
      const store = useAuthStore()
      await store.initSession()
      expect(usersService.getProfile).not.toHaveBeenCalled()
    })

    it('returns early if user already set', async () => {
      vi.mocked(getToken).mockReturnValue('tok')
      setActivePinia(createPinia())
      const store = useAuthStore()
      // token initialized from getToken mock
      expect(store.token).toBe('tok')
      store.user = mockUser as never
      await store.initSession()
      expect(usersService.getProfile).not.toHaveBeenCalled()
    })

    it('fetches profile on success', async () => {
      vi.mocked(getToken).mockReturnValue('tok')
      setActivePinia(createPinia())
      const store = useAuthStore()
      vi.mocked(usersService.getProfile).mockResolvedValue(mockUser as never)
      await store.initSession()
      expect(store.user).toEqual(mockUser)
      expect(store.loading).toBe(false)
    })

    it('clears session on failure', async () => {
      vi.mocked(getToken).mockReturnValue('tok')
      setActivePinia(createPinia())
      const store = useAuthStore()
      vi.mocked(usersService.getProfile).mockRejectedValue(new Error('401'))
      await store.initSession()
      expect(store.token).toBeNull()
      expect(store.user).toBeNull()
      expect(removeToken).toHaveBeenCalled()
      expect(store.loading).toBe(false)
    })

    it('sets loading during fetch', async () => {
      vi.mocked(getToken).mockReturnValue('tok')
      setActivePinia(createPinia())
      const store = useAuthStore()
      let loadingDuring = false
      vi.mocked(usersService.getProfile).mockImplementation(async () => {
        loadingDuring = store.loading
        return mockUser as never
      })
      await store.initSession()
      expect(loadingDuring).toBe(true)
      expect(store.loading).toBe(false)
    })
  })

  describe('checkStatus', () => {
    it('sets needsSetup true when not initialized', async () => {
      const store = useAuthStore()
      vi.mocked(authService.getStatus).mockResolvedValue({ initialized: false } as never)
      await store.checkStatus()
      expect(store.needsSetup).toBe(true)
    })

    it('sets needsSetup false when initialized', async () => {
      const store = useAuthStore()
      vi.mocked(authService.getStatus).mockResolvedValue({ initialized: true } as never)
      await store.checkStatus()
      expect(store.needsSetup).toBe(false)
    })

    it('sets needsSetup false on error', async () => {
      const store = useAuthStore()
      vi.mocked(authService.getStatus).mockRejectedValue(new Error('fail'))
      await store.checkStatus()
      expect(store.needsSetup).toBe(false)
    })
  })
})
