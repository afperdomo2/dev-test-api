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
      { path: '/generation', component: { template: '<div>generation</div>' } },
      { path: '/generation/settings', component: { template: '<div>settings</div>' } },
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
    expect(titles).toContain('Tips')
    expect(titles).toContain('Informes')
    expect(titles).not.toContain('Moderación')
  })

  it('filters nav for admin (shows Usuarios, Preguntas e Importar)', () => {
    const { wrapper } = mountLayout(true)
    const vm = wrapper.vm as unknown as {
      filteredSections: Array<{ title?: string; items: Array<{ title: string }> }>
    }
    const titles = vm.filteredSections.flatMap((s) => s.items.map((n) => n.title))
    expect(titles).toContain('Usuarios')
    expect(titles).toContain('Preguntas')
    expect(titles).toContain('Importar')
    expect(titles).toContain('Temas')
    expect(titles).not.toContain('Sesiones')
    expect(titles).not.toContain('Progreso')
    expect(titles).not.toContain('Informes')
  })

  it('exposes Moderación as admin-only item in Administración, and Tips in Contenido', () => {
    const { wrapper } = mountLayout(true)
    const vm = wrapper.vm as unknown as {
      filteredSections: Array<{ title?: string; items: Array<{ title: string; to: string }> }>
    }
    const adminSection = vm.filteredSections.find((s) => s.title === 'Administración')
    const adminItems = adminSection?.items ?? []
    expect(adminItems.map((i) => i.title)).toEqual(
      expect.arrayContaining(['Usuarios', 'Moderación']),
    )
    expect(adminItems.find((i) => i.title === 'Moderación')?.to).toBe('/reports')
    const contenidoSection = vm.filteredSections.find((s) => s.title === 'Contenido')
    const contenidoItems = contenidoSection?.items ?? []
    expect(contenidoItems.find((i) => i.title === 'Tips')?.to).toBe('/tips')
  })

  it('groups nav into sections with correct titles', () => {
    const { wrapper } = mountLayout(false)
    const vm = wrapper.vm as unknown as {
      filteredSections: Array<{ title?: string; items: Array<{ title: string; to: string }> }>
    }
    const sectionTitles = vm.filteredSections.map((s) => s.title)
    expect(sectionTitles).toContain('Estudio')
    expect(sectionTitles).toContain('Contenido')
    expect(sectionTitles).not.toContain('Administración')

    const estudioSection = vm.filteredSections.find((s) => s.title === 'Estudio')
    expect(estudioSection?.items.map((i) => i.title)).toEqual(['Sesiones', 'Progreso', 'Informes'])
    expect(estudioSection?.items.find((i) => i.title === 'Informes')?.to).toBe('/insights')

    const contenidoSection = vm.filteredSections.find((s) => s.title === 'Contenido')
    expect(contenidoSection?.items.map((i) => i.title)).toEqual(
      expect.arrayContaining(['Preguntas', 'Importar', 'Temas', 'Estadísticas']),
    )
  })

  it('exposes the generation entries as admin-only', () => {
    const { wrapper } = mountLayout(true)
    const vm = wrapper.vm as unknown as {
      filteredSections: Array<{ title?: string; items: Array<{ title: string; to: string }> }>
    }
    const adminItems = vm.filteredSections.find((s) => s.title === 'Administración')?.items ?? []
    const titles = adminItems.map((i) => i.title)

    expect(titles).toEqual(expect.arrayContaining(['Generar preguntas', 'Configuración de IA']))
    expect(adminItems.find((i) => i.title === 'Generar preguntas')?.to).toBe('/generation')
    expect(adminItems.find((i) => i.title === 'Configuración de IA')?.to).toBe(
      '/generation/settings',
    )
  })

  it('does not expose the generation entries to non-admins', () => {
    const { wrapper } = mountLayout(false)
    const vm = wrapper.vm as unknown as {
      filteredSections: Array<{ title?: string; items: Array<{ title: string }> }>
    }
    const titles = vm.filteredSections.flatMap((s) => s.items.map((n) => n.title))
    expect(titles).not.toContain('Generar preguntas')
    expect(titles).not.toContain('Configuración de IA')
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
