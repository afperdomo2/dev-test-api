import {
  createGenerationJob,
  getGenerationStats,
  listGenerationJobs,
  listSettings,
  updateSettings,
} from './generation.service'
import apiClient from '@/api/client'

vi.mock('@/api/client', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
  },
}))

const mockedGet = vi.mocked(apiClient.get)
const mockedPost = vi.mocked(apiClient.post)
const mockedPut = vi.mocked(apiClient.put)

describe('generation.service', () => {
  beforeEach(() => vi.clearAllMocks())

  it('createGenerationJob posts payload', async () => {
    const job = { id: 'j1' }
    mockedPost.mockResolvedValue({ data: job })
    const res = await createGenerationJob({
      topicIds: ['t1'],
      difficulty: 'beginner',
      quantity: 5,
    })
    expect(mockedPost).toHaveBeenCalledWith('/api/v1/generation', {
      topicIds: ['t1'],
      difficulty: 'beginner',
      quantity: 5,
    })
    expect(res).toEqual(job)
  })

  it('listGenerationJobs without status', async () => {
    const paginated = { data: [], meta: { total: 0 } }
    mockedGet.mockResolvedValue({ data: paginated })
    const res = await listGenerationJobs(1, 10)
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/generation', {
      params: { page: 1, perPage: 10 },
    })
    expect(res).toEqual(paginated)
  })

  it('listGenerationJobs with status filter', async () => {
    mockedGet.mockResolvedValue({ data: { data: [], meta: { total: 0 } } })
    await listGenerationJobs(2, 25, 'failed')
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/generation', {
      params: { page: 2, perPage: 25, status: 'failed' },
    })
  })

  it('getGenerationStats hits the stats subroute', async () => {
    const stats = { reservedToday: 3, dailyCap: 200 }
    mockedGet.mockResolvedValue({ data: stats })
    const res = await getGenerationStats()
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/generation/stats')
    expect(res).toEqual(stats)
  })

  it('listSettings hits the settings route', async () => {
    const settings = [{ key: 'ai.model', value: 'gpt-4o' }]
    mockedGet.mockResolvedValue({ data: settings })
    const res = await listSettings()
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/settings')
    expect(res).toEqual(settings)
  })

  it('updateSettings puts the wrapped payload', async () => {
    const settings = [{ key: 'ai.max_per_job', value: '8' }]
    mockedPut.mockResolvedValue({ data: settings })
    const res = await updateSettings({ settings: [{ key: 'ai.max_per_job', value: '8' }] })
    expect(mockedPut).toHaveBeenCalledWith('/api/v1/settings', {
      settings: [{ key: 'ai.max_per_job', value: '8' }],
    })
    expect(res).toEqual(settings)
  })
})
