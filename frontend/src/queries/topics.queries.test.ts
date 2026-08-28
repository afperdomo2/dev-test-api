/* eslint-disable @typescript-eslint/no-non-null-assertion */
import {
  topicsListOptions,
  topicDetailOptions,
  createTopicMutation,
  updateTopicMutation,
  deleteTopicMutation,
} from './topics.queries'
import * as topicsService from '@/api/services/topics.service'

vi.mock('@/api/services/topics.service', () => ({
  listTopics: vi.fn(),
  getTopicById: vi.fn(),
  createTopic: vi.fn(),
  updateTopic: vi.fn(),
  deleteTopic: vi.fn(),
}))

describe('topics.queries', () => {
  beforeEach(() => vi.clearAllMocks())

  it('topicsListOptions calls service', async () => {
    vi.mocked(topicsService.listTopics).mockResolvedValue({ data: [] } as never)
    const opts = topicsListOptions(() => 1, () => 10, () => 'name', () => 'asc', () => 'q', () => true)
    await (opts as unknown as { queryFn: (ctx: unknown) => Promise<unknown> }).queryFn({} as never)
    expect(topicsService.listTopics).toHaveBeenCalledWith(1, 10, 'name', 'asc', 'q', true)
  })

  it('topicsListOptions handles undefined optionals', async () => {
    vi.mocked(topicsService.listTopics).mockResolvedValue({} as never)
    const opts = topicsListOptions(() => 1, () => 10)
    await (opts as unknown as { queryFn: (ctx: unknown) => Promise<unknown> }).queryFn({} as never)
    expect(topicsService.listTopics).toHaveBeenCalledWith(1, 10, undefined, undefined, undefined, undefined)
  })

  it('topicDetailOptions', async () => {
    vi.mocked(topicsService.getTopicById).mockResolvedValue({ id: '1' } as never)
    const opts = topicDetailOptions(() => '1')
    const res = await (opts as unknown as { queryFn: (ctx: unknown) => Promise<unknown> }).queryFn({} as never)
    expect(topicsService.getTopicById).toHaveBeenCalledWith('1')
    expect(res).toEqual({ id: '1' })
  })

  it('createTopicMutation', () => {
    const mut = createTopicMutation()
    expect(mut.mutationKey).toEqual(['topics', 'create'])
    expect(mut.mutationFn).toBe(topicsService.createTopic)
  })

  it('updateTopicMutation delegates', async () => {
    vi.mocked(topicsService.updateTopic).mockResolvedValue({ id: '1' } as never)
    const mut = updateTopicMutation()
    const res = await mut.mutationFn({ id: '1', data: { name: 'x' } } as never)
    expect(topicsService.updateTopic).toHaveBeenCalledWith('1', { name: 'x' })
    expect(res).toEqual({ id: '1' })
  })

  it('deleteTopicMutation', () => {
    const mut = deleteTopicMutation()
    expect(mut.mutationKey).toEqual(['topics', 'delete'])
  })
})
