import { apiClient } from '@/core/api/client'
import { getProfile, getStatus, login, setup } from './auth.service'

jest.mock('@/core/api/client', () => ({
  apiClient: { get: jest.fn(), post: jest.fn(), put: jest.fn(), delete: jest.fn() },
}))

const mocked = apiClient as unknown as {
  get: jest.Mock
  post: jest.Mock
}

describe('auth.service', () => {
  beforeEach(() => jest.clearAllMocks())

  it('login', async () => {
    const auth = { token: 't', user: { id: '1' } }
    mocked.post.mockResolvedValue(auth)

    await expect(login({ email: 'a@b.com', password: 'x' })).resolves.toEqual(auth)
    expect(mocked.post).toHaveBeenCalledWith('/api/v1/auth/login', {
      email: 'a@b.com',
      password: 'x',
    })
  })

  it('setup', async () => {
    mocked.post.mockResolvedValue({ token: 't', user: { id: '1' } })

    const body = { email: 'a@b.com', password: 'x' }
    await setup(body)

    expect(mocked.post).toHaveBeenCalledWith('/api/v1/auth/setup', body)
  })

  it('getStatus', async () => {
    mocked.get.mockResolvedValue({ initialized: true })

    await expect(getStatus()).resolves.toEqual({ initialized: true })
    expect(mocked.get).toHaveBeenCalledWith('/api/v1/auth/status')
  })

  it('getProfile', async () => {
    mocked.get.mockResolvedValue({ id: 'u1', email: 'a@b.com' })

    await getProfile()

    expect(mocked.get).toHaveBeenCalledWith('/api/v1/profile')
  })
})
