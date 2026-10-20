import { apiClient } from '@/core/api/client'
import { getInsightSummary } from './insights.service'

jest.mock('@/core/api/client', () => ({
  apiClient: { get: jest.fn(), post: jest.fn(), put: jest.fn(), delete: jest.fn() },
}))

const mocked = apiClient as unknown as { get: jest.Mock }

describe('insights.service', () => {
  beforeEach(() => jest.clearAllMocks())

  it('getInsightSummary envía el periodo y devuelve el informe', async () => {
    const summary = { period: 'week', sessions: 2 }
    mocked.get.mockResolvedValue(summary)

    await expect(getInsightSummary('week')).resolves.toEqual(summary)
    expect(mocked.get).toHaveBeenCalledWith('/api/v1/insights/summary', { period: 'week' })
  })
})
