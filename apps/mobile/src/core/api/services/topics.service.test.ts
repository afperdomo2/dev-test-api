import { apiClient } from '@/core/api/client'
import { listTopics } from './topics.service'

jest.mock('@/core/api/client', () => ({
  apiClient: { get: jest.fn(), post: jest.fn(), put: jest.fn(), delete: jest.fn() },
}))

const mocked = apiClient as unknown as { get: jest.Mock }

describe('topics.service', () => {
  beforeEach(() => jest.clearAllMocks())

  it('usa paginación por defecto y propaga hasQuestions', async () => {
    mocked.get.mockResolvedValue({ data: [], meta: { total: 0 } })

    await listTopics({ hasQuestions: true })

    expect(mocked.get).toHaveBeenCalledWith('/api/v1/topics', {
      page: 1,
      perPage: 100,
      search: undefined,
      hasQuestions: true,
    })
  })

  it('reenvía búsqueda y perPage', async () => {
    mocked.get.mockResolvedValue({ data: [], meta: { total: 0 } })

    await listTopics({ search: 'java', hasQuestions: true, perPage: 50 })

    expect(mocked.get).toHaveBeenCalledWith('/api/v1/topics', {
      page: 1,
      perPage: 50,
      search: 'java',
      hasQuestions: true,
    })
  })
})
