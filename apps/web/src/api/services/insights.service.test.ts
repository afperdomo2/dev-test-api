import { getInsightSummary } from './insights.service'
import apiClient from '@/api/client'

vi.mock('@/api/client', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
  },
}))

const mockedGet = vi.mocked(apiClient.get)

describe('insights.service', () => {
  beforeEach(() => vi.clearAllMocks())

  it('llama al endpoint con el periodo y devuelve el informe', async () => {
    const summary = { period: 'week', sessions: 2 }
    mockedGet.mockResolvedValue({ data: summary })

    const res = await getInsightSummary('week')

    expect(mockedGet).toHaveBeenCalledWith('/api/v1/insights/summary', {
      params: { period: 'week' },
    })
    expect(res).toEqual(summary)
  })
})
