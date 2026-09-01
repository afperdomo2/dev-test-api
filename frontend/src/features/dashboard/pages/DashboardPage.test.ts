import { createPinia, setActivePinia } from 'pinia'
import DashboardPage from './DashboardPage.vue'
import { useAuthStore } from '@/stores/auth.store'
import { mountWithProviders } from '@/__tests__/utils/mount'
import type * as QuestionsService from '@/api/services/questions.service'

vi.mock('@/utils/storage', () => ({
  getToken: vi.fn(() => null),
  setToken: vi.fn(),
  removeToken: vi.fn(),
}))

vi.mock('@/api/services/questions.service', async () => {
  const actual = await vi.importActual<typeof QuestionsService>('@/api/services/questions.service')
  return {
    ...actual,
    getAiQuota: vi.fn(() => Promise.resolve({ dailyLimit: 20, usedToday: 5, remaining: 15 })),
    getImportQuota: vi.fn(() =>
      Promise.resolve({ dailyLimit: 200, usedToday: 10, remaining: 190 }),
    ),
  }
})

describe('DashboardPage', () => {
  async function mountDashboard(user: Record<string, unknown> | null = null) {
    const pinia = createPinia()
    setActivePinia(pinia)
    const auth = useAuthStore()
    if (user) {
      auth.user = user as never
      auth.token = 'tok'
    }
    return mountWithProviders(DashboardPage, { pinia })
  }

  it('renders dashboard title', async () => {
    const wrapper = await mountDashboard()
    expect(wrapper.text()).toContain('Dashboard')
  })

  it('shows user email when logged', async () => {
    const wrapper = await mountDashboard({
      email: 'a@b.com',
      isAdmin: false,
      id: '1',
      createdAt: '2026-08-28T12:00:00Z',
    })
    expect(wrapper.text()).toContain('a@b.com')
  })

  it('shows Administrador chip for admin', async () => {
    const wrapper = await mountDashboard({
      email: 'a@b.com',
      isAdmin: true,
      id: '1',
      createdAt: '2026-08-28T12:00:00Z',
    })
    expect(wrapper.text()).toContain('Administrador')
  })

  it('shows Usuario chip for non-admin', async () => {
    const wrapper = await mountDashboard({
      email: 'a@b.com',
      isAdmin: false,
      id: '1',
      createdAt: '2026-08-28T12:00:00Z',
    })
    expect(wrapper.text()).toContain('Usuario')
  })

  it('shows formattedDate when user has createdAt', async () => {
    const wrapper = await mountDashboard({
      email: 'a@b.com',
      id: '1',
      createdAt: '2026-01-15T10:00:00Z',
    })
    // formattedDate computed should not be empty
    expect(wrapper.text()).toContain('Miembro desde')
  })

  it('shows ID', async () => {
    const wrapper = await mountDashboard({
      email: 'a@b.com',
      id: 'user-123',
      createdAt: '2026-08-28T12:00:00Z',
    })
    expect(wrapper.text()).toContain('user-123')
  })
})
