import {
  getNotifications,
  getUnreadCount,
  markAllNotificationsRead,
  markNotificationRead,
} from './notifications.service'
import apiClient from '@/api/client'

vi.mock('@/api/client', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
  },
}))

const mockedGet = vi.mocked(apiClient.get)
const mockedPost = vi.mocked(apiClient.post)

const notification = {
  id: 'n1',
  userId: 'u1',
  type: 'daily_review' as const,
  title: 'Repaso de hoy',
  body: 'Tienes 3 preguntas para repasar',
  isRead: false,
  createdAt: '2026-10-20T12:00:00Z',
}

describe('notifications.service', () => {
  beforeEach(() => vi.clearAllMocks())

  it('getNotifications pide página y perPage y conserva el meta', async () => {
    const page = { data: [notification], meta: { total: 1, page: 2, perPage: 10 } }
    mockedGet.mockResolvedValue({ data: page })

    const res = await getNotifications(2, 10)

    expect(mockedGet).toHaveBeenCalledWith('/api/v1/notifications', {
      params: { page: 2, perPage: 10 },
    })
    expect(res).toEqual(page)
  })

  it('getUnreadCount', async () => {
    mockedGet.mockResolvedValue({ data: { count: 4 } })

    await expect(getUnreadCount()).resolves.toEqual({ count: 4 })
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/notifications/unread-count')
  })

  it('markNotificationRead hace POST al id', async () => {
    mockedPost.mockResolvedValue({ data: { ...notification, isRead: true } })

    await expect(markNotificationRead('n1')).resolves.toMatchObject({ isRead: true })
    expect(mockedPost).toHaveBeenCalledWith('/api/v1/notifications/n1/read')
  })

  it('markAllNotificationsRead hace POST sin body', async () => {
    mockedPost.mockResolvedValue({ data: { updated: 3 } })

    await expect(markAllNotificationsRead()).resolves.toEqual({ updated: 3 })
    expect(mockedPost).toHaveBeenCalledWith('/api/v1/notifications/read-all')
  })
})
