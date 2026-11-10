import { flushPromises, mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createMemoryHistory } from 'vue-router'
import { QueryClient, VueQueryPlugin } from '@tanstack/vue-query'
import DefaultLayout from './DefaultLayout.vue'
import { useAuthStore } from '@/stores/auth.store'
import { getUnreadCount } from '@/api/services/notifications.service'
import { APP_NAME } from '@/config/app'

vi.mock('@/utils/storage', async () => {
  const actual = await vi.importActual('@/utils/storage')
  return {
    ...(actual as object),
    getToken: vi.fn(() => null),
    setToken: vi.fn(),
    removeToken: vi.fn(),
  }
})

vi.mock('@/api/services/notifications.service', () => ({
  getUnreadCount: vi.fn(() => Promise.resolve({ count: 0 })),
  getNotifications: vi.fn(),
  markNotificationRead: vi.fn(),
  markAllNotificationsRead: vi.fn(),
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
      { path: '/notifications', component: { template: '<div>notifications</div>' } },
    ],
  })

  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  })

  const wrapper = mount(DefaultLayout, {
    global: {
      plugins: [router, [VueQueryPlugin, { queryClient }]],
      stubs: {
        RouterView: { template: '<div><slot /></div>' },
        'router-view': { template: '<div><slot /></div>' },
      },
    },
  })
  return { wrapper, auth, router }
}

describe('DefaultLayout', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
    document.title = ''
    vi.mocked(getUnreadCount).mockResolvedValue({ count: 0 })
  })

  it('renders the configured app name', () => {
    const { wrapper } = mountLayout()
    expect(wrapper.text()).toContain(APP_NAME)
  })

  it('muestra el avatar del usuario en la app bar', () => {
    const { wrapper } = mountLayout(true)

    const identity = wrapper.find('.identity-activator')
    expect(identity.exists()).toBe(true)
    expect(identity.find('.v-avatar').exists()).toBe(true)
    expect(identity.text()).not.toContain('test@example.com')
  })

  it('expone email y rol en el aria-label y el title del botón de cuenta', () => {
    const { wrapper } = mountLayout(true)

    const identity = wrapper.find('.identity-activator')
    expect(identity.attributes('aria-label')).toBe('Cuenta de test@example.com (Administrador)')
    expect(identity.attributes('title')).toBe('Cuenta de test@example.com (Administrador)')
  })

  it('identifica al usuario en el sidebar', () => {
    const { wrapper } = mountLayout(true)

    expect(wrapper.find('.v-navigation-drawer').text()).toContain('test@example.com')
  })

  it('incluye rol y email en el título de la pestaña', () => {
    mountLayout(true)

    expect(document.title).toBe(`Admin · test@example.com — ${APP_NAME}`)
  })

  it('vuelve al nombre de la app al cerrar sesión', async () => {
    const { wrapper } = mountLayout(true)
    const vm = wrapper.vm as unknown as { logout: () => void }

    vm.logout()
    await wrapper.vm.$nextTick()

    expect(document.title).toBe(APP_NAME)
  })

  it('muestra la campana de notificaciones', () => {
    const { wrapper } = mountLayout()
    expect(wrapper.find('[aria-label="Notificaciones"]').exists()).toBe(true)
  })

  it('el bloque de identidad del sidebar lleva a Mi perfil', () => {
    const { wrapper } = mountLayout(true)

    const drawer = wrapper.find('.v-navigation-drawer')
    const profileLink = drawer.find('a[href="/profile"]')
    expect(profileLink.exists()).toBe(true)
    expect(profileLink.attributes('aria-label')).toContain('test@example.com')
    expect(profileLink.text()).toContain('test@example.com')
    expect(profileLink.text()).toContain('Administrador')
  })

  it('no duplica el acceso a Mi perfil en el sidebar', () => {
    const { wrapper } = mountLayout()

    const drawer = wrapper.find('.v-navigation-drawer')
    expect(drawer.findAll('a[href="/profile"]')).toHaveLength(1)
  })

  it('muestra el badge con el número de no leídas', async () => {
    vi.mocked(getUnreadCount).mockResolvedValue({ count: 3 })
    const { wrapper } = mountLayout()
    await flushPromises()

    expect(wrapper.find('.v-badge__badge').text()).toBe('3')
  })

  it('oculta el badge cuando no hay no leídas', async () => {
    const { wrapper } = mountLayout()
    await flushPromises()

    const badge = wrapper.find('.v-badge__badge')
    expect(badge.exists() ? badge.isVisible() : false).toBe(false)
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
    expect(titles).toContain('Logros')
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
    expect(titles).not.toContain('Logros')
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
    expect(estudioSection?.items.map((i) => i.title)).toEqual([
      'Sesiones',
      'Progreso',
      'Informes',
      'Logros',
    ])
    expect(estudioSection?.items.find((i) => i.title === 'Informes')?.to).toBe('/insights')
    expect(estudioSection?.items.find((i) => i.title === 'Logros')?.to).toBe('/achievements')

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

    expect(titles).toEqual(expect.arrayContaining(['Generación con IA', 'Configuración de IA']))
    expect(adminItems.find((i) => i.title === 'Generación con IA')?.to).toBe('/generation')
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
    expect(titles).not.toContain('Generación con IA')
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
