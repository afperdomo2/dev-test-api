import { insightSummaryOptions } from './insights.queries'
import * as insightsService from '@/api/services/insights.service'

vi.mock('@/api/services/insights.service', () => ({
  getInsightSummary: vi.fn(),
}))

describe('insights.queries', () => {
  beforeEach(() => vi.clearAllMocks())

  it('insightSummaryOptions depende del periodo', async () => {
    vi.mocked(insightsService.getInsightSummary).mockResolvedValue({} as never)

    const periodFn = () => 'week' as const
    const opts = insightSummaryOptions(periodFn)

    expect(opts.queryKey).toEqual(['insights', 'summary', periodFn])

    await (opts as unknown as { queryFn: (ctx: unknown) => Promise<unknown> }).queryFn({} as never)
    expect(insightsService.getInsightSummary).toHaveBeenCalledWith('week')
  })
})
