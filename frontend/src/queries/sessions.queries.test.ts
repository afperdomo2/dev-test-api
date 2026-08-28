/* eslint-disable @typescript-eslint/no-non-null-assertion */
import {
  sessionsListOptions,
  sessionsInfiniteOptions,
  sessionDetailOptions,
  createSessionMutation,
  finishSessionMutation,
  deleteSessionMutation,
  submitAnswerMutation,
} from './sessions.queries'
import * as sessionsService from '@/api/services/sessions.service'

vi.mock('@/api/services/sessions.service', () => ({
  listSessions: vi.fn(),
  getSessionById: vi.fn(),
  getSessionDetail: vi.fn(),
  createSession: vi.fn(),
  finishSession: vi.fn(),
  deleteSession: vi.fn(),
  submitAnswer: vi.fn(),
  getNextQuestion: vi.fn(),
  getSessionSummary: vi.fn(),
}))

describe('sessions.queries', () => {
  beforeEach(() => vi.clearAllMocks())

  it('sessionsListOptions', async () => {
    vi.mocked(sessionsService.listSessions).mockResolvedValue({ data: [] } as never)
    const opts = sessionsListOptions(() => 1, () => 10)
    await (opts as unknown as { queryFn: (ctx: unknown) => Promise<unknown> }).queryFn({} as never)
    expect(sessionsService.listSessions).toHaveBeenCalledWith(1, 10)
  })

  it('sessionsInfiniteOptions getNextPageParam', async () => {
    vi.mocked(sessionsService.listSessions).mockResolvedValue({ data: [] } as never)
    const opts = sessionsInfiniteOptions(() => 'active')
    // queryFn with pageParam
    await (opts as unknown as { queryFn: (ctx: unknown) => Promise<unknown> }).queryFn({ pageParam: 2 } as never)
    expect(sessionsService.listSessions).toHaveBeenCalledWith(2, 20, 'active')

    // getNextPageParam logic
    const lastPage = { data: [{ id: '1' }], meta: { total: 3, page: 1, perPage: 20 } } as never
    const allPages = [{ data: [{ id: '1' }], meta: { total: 3 } } as never]
    const next = (opts as unknown as { getNextPageParam: (a:unknown,b:unknown)=>unknown }).getNextPageParam(lastPage, allPages as never)
    expect(next).toBe(2)

    const donePage = { data: [{ id: '1' }], meta: { total: 1 } } as never
    const done = (opts as unknown as { getNextPageParam: (a:unknown,b:unknown)=>unknown }).getNextPageParam(donePage, [donePage] as never)
    expect(done).toBeUndefined()
  })

  it('sessionDetailOptions', async () => {
    vi.mocked(sessionsService.getSessionById).mockResolvedValue({ id: '1' } as never)
    const opts = sessionDetailOptions(() => '1')
    const res = await (opts as unknown as { queryFn: (ctx: unknown) => Promise<unknown> }).queryFn({} as never)
    expect(sessionsService.getSessionById).toHaveBeenCalledWith('1')
    expect(res).toEqual({ id: '1' })
  })

  it('createSessionMutation', () => {
    const mut = createSessionMutation()
    expect(mut.mutationKey).toEqual(['sessions', 'create'])
    expect(mut.mutationFn).toBe(sessionsService.createSession)
  })

  it('finishSessionMutation', () => {
    const mut = finishSessionMutation()
    expect(mut.mutationKey).toEqual(['sessions', 'finish'])
  })

  it('deleteSessionMutation', () => {
    const mut = deleteSessionMutation()
    expect(mut.mutationKey).toEqual(['sessions', 'delete'])
  })

  it('submitAnswerMutation delegates', async () => {
    vi.mocked(sessionsService.submitAnswer).mockResolvedValue({ id: 'a1' } as never)
    const mut = submitAnswerMutation()
    const res = await mut.mutationFn({ sessionId: 's1', data: { answer: 'x' } as never } as never)
    expect(sessionsService.submitAnswer).toHaveBeenCalledWith('s1', { answer: 'x' })
    expect(res).toEqual({ id: 'a1' })
  })
})
