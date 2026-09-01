import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import DashboardPage from './DashboardPage.vue'
import { useAuthStore } from '@/stores/auth.store'

vi.mock('@/utils/storage', () => ({
  getToken: vi.fn(() => null),
  setToken: vi.fn(),
  removeToken: vi.fn(),
}))

describe('DashboardPage', () => {
  function mountDashboard(user: Record<string, unknown> | null = null) {
    setActivePinia(createPinia())
    const auth = useAuthStore()
    if (user) {
      auth.user = user as never
      auth.token = 'tok'
    }
    return mount(DashboardPage, {
      global: {
        stubs: { RouterView: true },
      },
    })
  }

  it('renders dashboard title', () => {
    const wrapper = mountDashboard()
    expect(wrapper.text()).toContain('Dashboard')
  })

  it('shows user email when logged', () => {
    const wrapper = mountDashboard({ email: 'a@b.com', isAdmin: false, id: '1', createdAt: '2026-08-28T12:00:00Z' })
    expect(wrapper.text()).toContain('a@b.com')
  })

  it('shows Administrador chip for admin', () => {
    const wrapper = mountDashboard({ email: 'a@b.com', isAdmin: true, id: '1', createdAt: '2026-08-28T12:00:00Z' })
    expect(wrapper.text()).toContain('Administrador')
  })

  it('shows Usuario chip for non-admin', () => {
    const wrapper = mountDashboard({ email: 'a@b.com', isAdmin: false, id: '1', createdAt: '2026-08-28T12:00:00Z' })
    expect(wrapper.text()).toContain('Usuario')
  })

  it('shows formattedDate when user has createdAt', () => {
    const wrapper = mountDashboard({ email: 'a@b.com', id: '1', createdAt: '2026-01-15T10:00:00Z' })
    // formattedDate computed should not be empty
    expect(wrapper.text()).toContain('Miembro desde')
  })

  it('shows ID', () => {
    const wrapper = mountDashboard({ email: 'a@b.com', id: 'user-123', createdAt: '2026-08-28T12:00:00Z' })
    expect(wrapper.text()).toContain('user-123')
  })
})
