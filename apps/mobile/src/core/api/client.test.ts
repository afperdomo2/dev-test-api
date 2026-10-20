import { apiClient, onUnauthorized } from './client'
import { getToken, removeToken } from '@/core/auth/storage'

jest.mock('@/core/auth/storage', () => ({
  getToken: jest.fn(),
  removeToken: jest.fn(),
}))

const mockGetToken = jest.mocked(getToken)
const mockRemoveToken = jest.mocked(removeToken)

function jsonResponse(body: unknown, status = 200): Response {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: jest.fn(async () => body),
  } as unknown as Response
}

describe('apiClient', () => {
  const fetchMock = jest.fn()

  /** Última llamada a fetch, sin non-null assertions. */
  function lastFetchCall(): { url: string; init: { headers: Record<string, string>; body?: string } } {
    const call = fetchMock.mock.calls[fetchMock.mock.calls.length - 1]
    if (!call) throw new Error('fetch no fue llamado')
    return { url: call[0] as string, init: call[1] as never }
  }

  beforeEach(() => {
    jest.clearAllMocks()
    globalThis.fetch = fetchMock as unknown as typeof fetch
    mockGetToken.mockResolvedValue('tok')
  })

  it('envía el Bearer cuando hay token', async () => {
    fetchMock.mockResolvedValue(jsonResponse({ data: { id: '1' } }))

    await apiClient.get('/api/v1/x')

    expect(lastFetchCall().init.headers.Authorization).toBe('Bearer tok')
  })

  it('no envía Authorization cuando no hay token', async () => {
    mockGetToken.mockResolvedValue(null)
    fetchMock.mockResolvedValue(jsonResponse({ data: null }))

    await apiClient.get('/api/v1/x')

    expect(lastFetchCall().init.headers.Authorization).toBeUndefined()
  })

  it('desempaqueta el envelope { data }', async () => {
    fetchMock.mockResolvedValue(jsonResponse({ data: { id: '42' } }))

    await expect(apiClient.get('/x')).resolves.toEqual({ id: '42' })
  })

  it('NO desempaqueta respuestas paginadas (llevan meta)', async () => {
    const paginated = { data: [{ id: '1' }], meta: { total: 1, page: 1, perPage: 20 } }
    fetchMock.mockResolvedValue(jsonResponse(paginated))

    await expect(apiClient.get('/x')).resolves.toEqual(paginated)
  })

  it('propaga el ApiError RFC 9457 del backend', async () => {
    const apiError = {
      type: 'about:blank',
      title: 'Not Found',
      status: 404,
      detail: 'No existe',
      instance: '/x',
    }
    fetchMock.mockResolvedValue(jsonResponse(apiError, 404))

    await expect(apiClient.get('/x')).rejects.toMatchObject(apiError)
  })

  it('construye un ApiError de fallback si el cuerpo no es JSON', async () => {
    fetchMock.mockResolvedValue({
      ok: false,
      status: 500,
      json: async () => {
        throw new Error('no es json')
      },
    } as unknown as Response)

    await expect(apiClient.get('/x')).rejects.toMatchObject({
      title: 'Internal Server Error',
      status: 500,
      instance: '/x',
    })
  })

  it('en 401 borra el token y dispara onUnauthorized', async () => {
    const handler = jest.fn()
    onUnauthorized(handler)
    fetchMock.mockResolvedValue(jsonResponse({ title: 'Unauthorized', status: 401 }, 401))

    await expect(apiClient.get('/x')).rejects.toBeDefined()

    expect(mockRemoveToken).toHaveBeenCalledTimes(1)
    expect(handler).toHaveBeenCalledTimes(1)
  })

  it('lanza Network Error cuando fetch rechaza (contrato de mobile-offline)', async () => {
    fetchMock.mockRejectedValue(new TypeError('Network request failed'))

    await expect(apiClient.get('/x')).rejects.toMatchObject({
      title: 'Network Error',
      status: 500,
      instance: '/x',
    })
  })

  it('serializa query params y omite los undefined', async () => {
    fetchMock.mockResolvedValue(jsonResponse({ data: [] }))

    await apiClient.get('/x', { page: 2, status: undefined, only: true })

    expect(lastFetchCall().url).toBe('http://localhost:8080/x?page=2&only=true')
  })
})
