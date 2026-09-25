import {
  tipsListOptions,
  tipDetailOptions,
  createTipMutation,
  updateTipMutation,
  deleteTipMutation,
  importTipsMutation,
  importTipsTextMutation,
} from './tips.queries'
import * as tipsService from '@/api/services/tips.service'

vi.mock('@/api/services/tips.service', () => ({
  listTips: vi.fn(),
  getTipById: vi.fn(),
  createTip: vi.fn(),
  updateTip: vi.fn(),
  deleteTip: vi.fn(),
  importTips: vi.fn(),
  importTipsFromText: vi.fn(),
}))

describe('tips.queries', () => {
  beforeEach(() => vi.clearAllMocks())

  it('tipsListOptions calls service with pagination and filters', async () => {
    vi.mocked(tipsService.listTips).mockResolvedValue({ data: [] } as never)
    const opts = tipsListOptions(
      () => 2,
      () => 20,
      () => ({ topicId: 'topic-1' }),
    )
    await (opts as unknown as { queryFn: (ctx: unknown) => Promise<unknown> }).queryFn({} as never)
    expect(tipsService.listTips).toHaveBeenCalledWith(2, 20, { topicId: 'topic-1' })
  })

  it('tipDetailOptions', async () => {
    vi.mocked(tipsService.getTipById).mockResolvedValue({ id: 't1' } as never)
    const opts = tipDetailOptions(() => 't1')
    const res = await (opts as unknown as { queryFn: (ctx: unknown) => Promise<unknown> }).queryFn(
      {} as never,
    )
    expect(tipsService.getTipById).toHaveBeenCalledWith('t1')
    expect(res).toEqual({ id: 't1' })
  })

  it('createTipMutation', () => {
    const mut = createTipMutation()
    expect(mut.mutationKey).toEqual(['tips', 'create'])
    expect(mut.mutationFn).toBe(tipsService.createTip)
  })

  it('updateTipMutation delegates id and data', async () => {
    vi.mocked(tipsService.updateTip).mockResolvedValue({ id: 't1' } as never)
    const mut = updateTipMutation()
    const res = await mut.mutationFn({ id: 't1', data: { content: 'nuevo' } })
    expect(tipsService.updateTip).toHaveBeenCalledWith('t1', { content: 'nuevo' })
    expect(res).toEqual({ id: 't1' })
  })

  it('deleteTipMutation', () => {
    const mut = deleteTipMutation()
    expect(mut.mutationKey).toEqual(['tips', 'delete'])
    expect(mut.mutationFn).toBe(tipsService.deleteTip)
  })

  it('importTipsMutation', async () => {
    vi.mocked(tipsService.importTips).mockResolvedValue({
      total: 1,
      imported: 1,
      failed: 0,
      errors: [],
    })
    const mut = importTipsMutation()
    const formData = new FormData()
    const res = await mut.mutationFn(formData)
    expect(tipsService.importTips).toHaveBeenCalledWith(formData)
    expect(res.imported).toBe(1)
  })

  it('importTipsTextMutation', async () => {
    vi.mocked(tipsService.importTipsFromText).mockResolvedValue({
      total: 0,
      imported: 0,
      failed: 0,
      errors: [],
    })
    const mut = importTipsTextMutation()
    await mut.mutationFn('content,topic')
    expect(tipsService.importTipsFromText).toHaveBeenCalledWith('content,topic')
  })
})
