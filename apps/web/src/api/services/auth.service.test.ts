import { login, setup, getStatus } from './auth.service'
import apiClient from '@/api/client'

vi.mock('@/api/client', () => ({
  default: {
    post: vi.fn(),
    get: vi.fn(),
  },
}))

const mockedPost = vi.mocked(apiClient.post)
const mockedGet = vi.mocked(apiClient.get)

describe('auth.service', () => {
  beforeEach(() => vi.clearAllMocks())

  describe('login', () => {
    it('posts to /api/v1/auth/login and returns data', async () => {
      const payload = { email: 'a@b.com', password: 'secret123' }
      const response = { token: 'tok', user: { id: '1' } }
      mockedPost.mockResolvedValue({ data: response })

      const result = await login(payload as never)

      expect(mockedPost).toHaveBeenCalledWith('/api/v1/auth/login', payload)
      expect(result).toEqual(response)
    })
  })

  describe('setup', () => {
    it('posts to /api/v1/auth/setup', async () => {
      const payload = { email: 'a@b.com', password: 'secret123' }
      const response = { token: 'tok' }
      mockedPost.mockResolvedValue({ data: response })

      const result = await setup(payload as never)

      expect(mockedPost).toHaveBeenCalledWith('/api/v1/auth/setup', payload)
      expect(result).toEqual(response)
    })
  })

  describe('getStatus', () => {
    it('gets /api/v1/auth/status', async () => {
      const response = { initialized: true }
      mockedGet.mockResolvedValue({ data: response })

      const result = await getStatus()

      expect(mockedGet).toHaveBeenCalledWith('/api/v1/auth/status')
      expect(result).toEqual(response)
    })
  })
})
