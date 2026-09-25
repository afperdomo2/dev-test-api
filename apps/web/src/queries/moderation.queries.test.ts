import {
  reportsListOptions,
  createReportMutation,
  reviewReportMutation,
} from './moderation.queries'
import * as moderationService from '@/api/services/moderation.service'

vi.mock('@/api/services/moderation.service', () => ({
  createReport: vi.fn(),
  listReports: vi.fn(),
  reviewReport: vi.fn(),
}))

describe('moderation.queries', () => {
  beforeEach(() => vi.clearAllMocks())

  it('reportsListOptions calls service without status', async () => {
    vi.mocked(moderationService.listReports).mockResolvedValue({ data: [] } as never)
    const opts = reportsListOptions(
      () => 1,
      () => 10,
      () => undefined,
    )
    await (opts as unknown as { queryFn: (ctx: unknown) => Promise<unknown> }).queryFn({} as never)
    expect(moderationService.listReports).toHaveBeenCalledWith(1, 10, undefined)
  })

  it('reportsListOptions passes status filter', async () => {
    vi.mocked(moderationService.listReports).mockResolvedValue({ data: [] } as never)
    const opts = reportsListOptions(
      () => 1,
      () => 10,
      () => 'pending',
    )
    await (opts as unknown as { queryFn: (ctx: unknown) => Promise<unknown> }).queryFn({} as never)
    expect(moderationService.listReports).toHaveBeenCalledWith(1, 10, 'pending')
  })

  it('createReportMutation', () => {
    const mut = createReportMutation()
    expect(mut.mutationKey).toEqual(['reports', 'create'])
    expect(mut.mutationFn).toBe(moderationService.createReport)
  })

  it('reviewReportMutation delegates id and data', async () => {
    vi.mocked(moderationService.reviewReport).mockResolvedValue({ id: 'r1' } as never)
    const mut = reviewReportMutation()
    const res = await mut.mutationFn({ id: 'r1', data: { unpublish: true } })
    expect(moderationService.reviewReport).toHaveBeenCalledWith('r1', { unpublish: true })
    expect(res).toEqual({ id: 'r1' })
  })
})
