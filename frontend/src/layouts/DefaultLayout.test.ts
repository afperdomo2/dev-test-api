import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createMemoryHistory } from 'vue-router'
import DefaultLayout from './DefaultLayout.vue'
import { useAuthStore } from '@/stores/auth.store'

vi.mock('@/utils/storage', () => ({
  getToken: vi.fn(() => null),
  setToken: vi.fn(),
  removeToken: vi.fn(),
}))

function mountLayout(isAdmin = false) {
  setActivePinia(createPinia())
  const auth = useAuthStore()
  auth.user = { email: 'test@example.com', isAdmin } as never
  auth.token = 'tok'

  const router = createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/', component: { template: '<div>home</div>' } },
      { path: '/login', component: { template: '<div>login</div>' } },
    ],
  })

  const wrapper = mount(DefaultLayout, {
    global: {
      plugins: [router],
      stubs: {
        RouterView: { template: '<div><slot /></div>' },
        'router-view': { template: '<div><slot /></div>' },
      },
    },
  })
  return { wrapper, auth, router }
}

describe('DefaultLayout', () => {
  it('renders DevTest branding', () => {
    const { wrapper } = mountLayout()
    expect(wrapper.text()).toContain('DevTest')
  })

  it('filters nav for non-admin (no Usuarios, shows Preguntas)', () => {
    const { wrapper } = mountLayout(false)
    const vm = wrapper.vm as unknown as {
      filteredSections: Array<{ title?: string; items: Array<{ title: string }> }>
    }
    const titles = vm.filteredSections.flatMap((s) => s.items.map((n) => n.title))
    expect(titles).toContain('Preguntas')
    expect(titles).toContain('Sesiones')
    expect(titles).not.toContain('Usuarios')
  })

  it('filters nav for admin (shows Usuarios, hides Preguntas)', () => {
    const { wrapper } = mountLayout(true)
    const vm = wrapper.vm as unknown as {
      filteredSections: Array<{ title?: string; items: Array<{ title: string }> }>
    }
    const titles = vm.filteredSections.flatMap((s) => s.items.map((n) => n.title))
    expect(titles).toContain('Usuarios')
    expect(titles).not.toContain('Preguntas')
  })

  it('groups nav into sections with correct titles', () => {
    const { wrapper } = mountLayout(false)
    const vm = wrapper.vm as unknown as {
      filteredSections: Array<{ title?: string; items: Array<{ title: string }> }>
    }
    const sectionTitles = vm.filteredSections.map((s) => s.title)
    // Top section has no title (Dashboard), then Preguntas and Estudio
    expect(sectionTitles).toContain('Preguntas')
    expect(sectionTitles).toContain('Estudio')
    expect(sectionTitles).not.toContain('Administración')
    const preguntasSection = vm.filteredSections.find((s) => s.title === 'Preguntas')
    expect(preguntasSection?.items.map((i) => i.title)).toEqual(
      expect.arrayContaining(['Preguntas', 'Importar', 'Temas', 'Estadísticas']),
    )
  })

  it('hides empty sections for admin', () => {
    const { wrapper } = mountLayout(true)
    const vm = wrapper.vm as unknown as {
      filteredSections: Array<{ title?: string; items: Array<{ title: string }> }>
    }
    const sectionTitles = vm.filteredSections.map((s) => s.title)
    expect(sectionTitles).not.toContain('Estudio')
    expect(sectionTitles).toContain('Administración')
  })

  it('shows user email and role', () => {
    const { wrapper } = mountLayout(true)
    const vm = wrapper.vm as unknown as { authStore: { user: { email: string }; isAdmin: boolean } }
    expect(vm.authStore.user.email).toBe('test@example.com')
    expect(vm.authStore.isAdmin).toBe(true)
  })

  it('logout clears session and pushes /login', async () => {
    const { wrapper, auth, router } = mountLayout()
    const pushSpy = vi.spyOn(router, 'push')
    const clearSpy = vi.spyOn(auth, 'clearSession')
    const vm = wrapper.vm as unknown as { logout: () => void }
    vm.logout()
    expect(clearSpy).toHaveBeenCalled()
    expect(pushSpy).toHaveBeenCalledWith('/login')
  })

  it('toggles theme and sidebar via appStore', async () => {
    const { wrapper } = mountLayout()
    const vm = wrapper.vm as unknown as {
      appStore: { theme: string; toggleTheme: () => void; toggleSidebar: () => void }
    }
    const before = vm.appStore.theme
    vm.appStore.toggleTheme()
    expect(vm.appStore.theme).not.toBe(before)
  })
})
