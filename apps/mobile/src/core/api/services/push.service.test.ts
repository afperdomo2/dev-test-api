import { apiClient } from '@/core/api/client'
import { registerToken, unregisterToken } from './push.service'

jest.mock('@/core/api/client', () => ({
  apiClient: { post: jest.fn() },
}))

const mockedPost = jest.mocked(apiClient.post)

describe('push.service', () => {
  beforeEach(() => jest.clearAllMocks())

  it('registerToken hace POST con token y plataforma', async () => {
    mockedPost.mockResolvedValue({ data: { registered: true } })

    await registerToken('ExponentPushToken[abc]', 'ios')

    expect(mockedPost).toHaveBeenCalledWith('/api/v1/push/tokens', {
      token: 'ExponentPushToken[abc]',
      platform: 'ios',
    })
  })

  it('unregisterToken hace POST con el token', async () => {
    mockedPost.mockResolvedValue({ data: { removed: true } })

    await unregisterToken('ExponentPushToken[abc]')

    expect(mockedPost).toHaveBeenCalledWith('/api/v1/push/tokens/unregister', {
      token: 'ExponentPushToken[abc]',
    })
  })
})
