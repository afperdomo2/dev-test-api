import { createReport, listReports, reviewReport } from './moderation.service'
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

describe('moderation.service', () => {
  beforeEach(() => vi.clearAllMocks())

  it('createReport posts payload', async () => {
    const report = { id: 'r1' }
    mockedPost.mockResolvedValue({ data: report })
    const res = await createReport({ targetType: 'question', targetId: 'q1', reason: 'mal' })
    expect(mockedPost).toHaveBeenCalledWith('/api/v1/reports', {
      targetType: 'question',
      targetId: 'q1',
      reason: 'mal',
    })
    expect(res).toEqual(report)
  })

  it('listReports without status', async () => {
    const paginated = { data: [], meta: { total: 0 } }
    mockedGet.mockResolvedValue({ data: paginated })
    const res = await listReports(1, 10)
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/reports', {
      params: { page: 1, perPage: 10 },
    })
    expect(res).toEqual(paginated)
  })

  it('listReports with status filter', async () => {
    mockedGet.mockResolvedValue({ data: { data: [], meta: { total: 0 } } })
    await listReports(1, 10, 'pending')
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/reports', {
      params: { page: 1, perPage: 10, status: 'pending' },
    })
  })

  it('reviewReport puts unpublish flag', async () => {
    const report = { id: 'r1', status: 'reviewed' }
    mockedPut.mockResolvedValue({ data: report })
    const res = await reviewReport('r1', { unpublish: true })
    expect(mockedPut).toHaveBeenCalledWith('/api/v1/reports/r1/review', { unpublish: true })
    expect(res).toEqual(report)
  })
})
