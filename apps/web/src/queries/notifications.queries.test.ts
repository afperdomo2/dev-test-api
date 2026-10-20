import {
  markAllReadMutation,
  markReadMutation,
  notificationsOptions,
  unreadCountOptions,
} from './notifications.queries'
import * as notificationsService from '@/api/services/notifications.service'

vi.mock('@/api/services/notifications.service', () => ({
  getNotifications: vi.fn(),
  getUnreadCount: vi.fn(),
  markNotificationRead: vi.fn(),
  markAllNotificationsRead: vi.fn(),
}))

describe('notifications.queries', () => {
  beforeEach(() => vi.clearAllMocks())

  it('notificationsOptions incluye página y perPage en la key', async () => {
    vi.mocked(notificationsService.getNotifications).mockResolvedValue({ data: [] } as never)
    const page = () => 2
    const perPage = () => 10
    const opts = notificationsOptions(page, perPage)

    expect(opts.queryKey).toEqual(['notifications', 'list', page, perPage])

    await (opts as unknown as { queryFn: (ctx: unknown) => Promise<unknown> }).queryFn({} as never)
    expect(notificationsService.getNotifications).toHaveBeenCalledWith(2, 10)
  })

  it('unreadCountOptions sondea cada 60 s y al enfocar', async () => {
    vi.mocked(notificationsService.getUnreadCount).mockResolvedValue({ count: 2 })
    const opts = unreadCountOptions()

    expect(opts.queryKey).toEqual(['notifications', 'unread-count'])
    expect(opts.refetchInterval).toBe(60 * 1000)
    expect(opts.refetchOnWindowFocus).toBe(true)
  })

  it('markReadMutation usa el servicio', () => {
    const mut = markReadMutation()
    expect(mut.mutationKey).toEqual(['notifications', 'mark-read'])

    vi.mocked(notificationsService.markNotificationRead).mockResolvedValue({ id: 'n1' } as never)
    mut.mutationFn({ id: 'n1' })
    expect(notificationsService.markNotificationRead).toHaveBeenCalledWith('n1')
  })

  it('markAllReadMutation usa el servicio', () => {
    const mut = markAllReadMutation()
    expect(mut.mutationKey).toEqual(['notifications', 'mark-all-read'])
    expect(mut.mutationFn).toBe(notificationsService.markAllNotificationsRead)
  })
})
