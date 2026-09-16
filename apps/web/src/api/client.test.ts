import apiClient from './client'
import MockAdapter from 'axios-mock-adapter'

vi.mock('@/utils/storage', async () => {
  const actual = await vi.importActual('@/utils/storage')
  return {
    ...(actual as object),
    getToken: vi.fn(() => null),
    removeToken: vi.fn(),
  }
})

import { getToken, removeToken } from '@/utils/storage'

describe('api/client', () => {
  let mock: MockAdapter

  beforeEach(() => {
    mock = new MockAdapter(apiClient)
    vi.clearAllMocks()
    vi.mocked(getToken).mockReturnValue(null)
  })

  afterEach(() => {
    mock.restore()
  })

  describe('request interceptor - Authorization', () => {
    it('adds Bearer token when present', async () => {
      vi.mocked(getToken).mockReturnValue('my-token')
      mock.onGet('/test').reply(200, { data: 'ok' })

      await apiClient.get('/test')

      const req = mock.history.get[0]
      expect(req).toBeDefined()
      expect(req?.headers?.Authorization).toBe('Bearer my-token')
    })

    it('does not add Authorization when no token', async () => {
      mock.onGet('/test').reply(200, { data: 'ok' })
      await apiClient.get('/test')
      const req = mock.history.get[0]
      expect(req).toBeDefined()
      expect(req?.headers?.Authorization).toBeUndefined()
    })
  })

  describe('response interceptor - unwrap', () => {
    it('unwraps { data: ... } without meta', async () => {
      mock.onGet('/unwrap').reply(200, { data: { id: 1, name: 'foo' } })
      const res = await apiClient.get('/unwrap')
      expect(res.data).toEqual({ id: 1, name: 'foo' })
    })

    it('does not unwrap when meta present (paginated)', async () => {
      const paginated = { data: [{ id: 1 }], meta: { total: 1, page: 1, perPage: 10 } }
      mock.onGet('/paginated').reply(200, paginated)
      const res = await apiClient.get('/paginated')
      expect(res.data).toEqual(paginated)
    })

    it('does not unwrap primitive or null', async () => {
      mock.onGet('/plain').reply(200, 'plain string')
      const res = await apiClient.get('/plain')
      expect(res.data).toBe('plain string')
    })
  })

  describe('response interceptor - error handling', () => {
    it('calls removeToken on 401', async () => {
      mock.onGet('/protected').reply(401, {
        type: 'about:blank',
        title: 'Unauthorized',
        status: 401,
        detail: 'No autorizado',
        instance: '/api/v1/profile',
      })

      await expect(apiClient.get('/protected')).rejects.toEqual(
        expect.objectContaining({ status: 401 }),
      )
      expect(removeToken).toHaveBeenCalled()
    })

    it('does not call removeToken on non-401', async () => {
      mock.onGet('/err').reply(400, {
        type: 'about:blank',
        title: 'Bad Request',
        status: 400,
        detail: 'bad',
        instance: '/err',
      })
      await expect(apiClient.get('/err')).rejects.toEqual(expect.objectContaining({ status: 400 }))
      expect(removeToken).not.toHaveBeenCalled()
    })

    it('rejects with ApiError from response', async () => {
      const apiError = {
        type: 'about:blank',
        title: 'Not Found',
        status: 404,
        detail: 'No encontrado',
        instance: '/missing',
      }
      mock.onGet('/missing').reply(404, apiError)

      await expect(apiClient.get('/missing')).rejects.toEqual(apiError)
    })

    it('rejects with fallback 500 when no response', async () => {
      mock.onGet('/network').networkError()

      await expect(apiClient.get('/network')).rejects.toEqual(
        expect.objectContaining({
          status: 500,
          title: 'Internal Server Error',
          detail: 'Ocurrió un error inesperado',
        }),
      )
    })

    it('fallback instance uses request url', async () => {
      mock.onGet('/fallback-url').networkError()
      await expect(apiClient.get('/fallback-url')).rejects.toEqual(
        expect.objectContaining({ instance: '/fallback-url' }),
      )
    })
  })
})
