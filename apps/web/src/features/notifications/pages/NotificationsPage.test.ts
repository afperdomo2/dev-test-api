import { flushPromises } from '@vue/test-utils'
import NotificationsPage from './NotificationsPage.vue'
import { createTestQueryClient, mountWithProviders } from '@/__tests__/utils/mount'
import {
  getNotifications,
  markAllNotificationsRead,
  markNotificationRead,
} from '@/api/services/notifications.service'
import type { Notification } from '@/types/notification.types'

vi.mock('@/api/services/notifications.service', () => ({
  getNotifications: vi.fn(),
  getUnreadCount: vi.fn(),
  markNotificationRead: vi.fn(),
  markAllNotificationsRead: vi.fn(),
}))

const items: Array<Notification> = [
  {
    id: 'n1',
    userId: 'u1',
    type: 'daily_review',
    title: 'Repaso de hoy',
    body: 'Tienes 3 preguntas para repasar',
    isRead: false,
    createdAt: new Date().toISOString(),
  },
  {
    id: 'n2',
    userId: 'u1',
    type: 'daily_review',
    title: 'Repaso de ayer',
    body: 'Tienes 2 preguntas para repasar',
    isRead: true,
    createdAt: new Date(Date.now() - 60 * 60 * 1000).toISOString(),
  },
]

const response = { data: items, meta: { total: 2, page: 1, perPage: 20 } }

async function mountPage() {
  const wrapper = await mountWithProviders(NotificationsPage)
  await flushPromises()
  return wrapper
}

describe('NotificationsPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(getNotifications).mockResolvedValue(response)
    vi.mocked(markNotificationRead).mockResolvedValue({ ...items[0], isRead: true } as Notification)
    vi.mocked(markAllNotificationsRead).mockResolvedValue({ updated: 1 })
  })

  it('pinta el historial y distingue las no leídas con texto', async () => {
    const wrapper = await mountPage()

    expect(wrapper.text()).toContain('Notificaciones')
    expect(wrapper.text()).toContain('Repaso de hoy')
    expect(wrapper.text()).toContain('Tienes 3 preguntas para repasar')
    expect(wrapper.text()).toContain('Sin leer')
    expect(wrapper.text().match(/Sin leer/g)?.length).toBe(1)
    expect(wrapper.text()).toContain('Marcar leída')
  })

  it('sin notificaciones muestra el estado vacío', async () => {
    vi.mocked(getNotifications).mockResolvedValue({
      data: [],
      meta: { total: 0, page: 1, perPage: 20 },
    })

    const wrapper = await mountPage()

    expect(wrapper.text()).toContain('No tienes notificaciones')
  })

  it('marcar una como leída llama al servicio', async () => {
    const wrapper = await mountPage()

    const vm = wrapper.vm as unknown as { handleMarkRead: (id: string) => Promise<void> }
    await vm.handleMarkRead('n1')
    await flushPromises()

    expect(markNotificationRead).toHaveBeenCalledWith('n1')
  })

  it('marcar todas invalida la caché de notificaciones', async () => {
    const queryClient = createTestQueryClient()
    const invalidateSpy = vi.spyOn(queryClient, 'invalidateQueries')
    const wrapper = await mountWithProviders(NotificationsPage, { queryClient })
    await flushPromises()

    const button = wrapper
      .findAll('button')
      .find((b) => b.text().includes('Marcar todas como leídas'))
    await button?.trigger('click')
    await flushPromises()

    expect(markAllNotificationsRead).toHaveBeenCalled()
    expect(invalidateSpy).toHaveBeenCalledWith({ queryKey: ['notifications'] })
  })

  it('si falla la carga permite reintentar', async () => {
    vi.mocked(getNotifications).mockRejectedValueOnce(new Error('boom'))

    const wrapper = await mountPage()

    expect(wrapper.text()).toContain('No se pudieron cargar las notificaciones')
    expect(wrapper.text()).toContain('Reintentar')
  })
})
