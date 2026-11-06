/* eslint-disable @typescript-eslint/no-non-null-assertion */
import {
  sessionsListOptions,
  sessionsResumeOptions,
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

  it('sessionsListOptions forwards pagination, filter and sort', async () => {
    vi.mocked(sessionsService.listSessions).mockResolvedValue({ data: [] } as never)
    const opts = sessionsListOptions(
      () => 2,
      () => 20,
      () => 'completed',
      () => 'score',
      () => 'desc',
    )
    await (opts as unknown as { queryFn: (ctx: unknown) => Promise<unknown> }).queryFn({} as never)
    expect(sessionsService.listSessions).toHaveBeenCalledWith(2, 20, 'completed', 'score', 'desc')
    const key = (opts as unknown as { queryKey: Array<unknown> }).queryKey
    expect(key.slice(0, 2)).toEqual(['sessions', 'list'])
    expect((key[2] as () => number)()).toBe(2)
    expect((key[3] as () => number)()).toBe(20)
    expect((key[4] as () => string)()).toBe('completed')
    expect((key[5] as () => string)()).toBe('score')
    expect((key[6] as () => string)()).toBe('desc')
  })

  it('sessionsResumeOptions asks for the latest in-progress session', async () => {
    vi.mocked(sessionsService.listSessions).mockResolvedValue({ data: [] } as never)
    const opts = sessionsResumeOptions()
    await (opts as unknown as { queryFn: (ctx: unknown) => Promise<unknown> }).queryFn({} as never)
    expect(sessionsService.listSessions).toHaveBeenCalledWith(
      1,
      1,
      'in_progress',
      'started_at',
      'desc',
    )
    expect((opts as unknown as { queryKey: Array<unknown> }).queryKey).toEqual([
      'sessions',
      'list',
      'resume',
    ])
  })

  it('sessionDetailOptions', async () => {
    vi.mocked(sessionsService.getSessionById).mockResolvedValue({ id: '1' } as never)
    const opts = sessionDetailOptions(() => '1')
    const res = await (opts as unknown as { queryFn: (ctx: unknown) => Promise<unknown> }).queryFn(
      {} as never,
    )
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
