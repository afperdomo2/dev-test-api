/* eslint-disable @typescript-eslint/no-non-null-assertion */
import {
  upcomingQuestionsOptions,
  savedQuestionsOptions,
  submitProgressAnswerMutation,
  toggleSaveMutation,
} from './progress.queries'
import * as progressService from '@/api/services/progress.service'

vi.mock('@/api/services/progress.service', () => ({
  getUpcomingReviews: vi.fn(),
  getSavedQuestions: vi.fn(),
  submitProgressAnswer: vi.fn(),
  toggleSaveQuestion: vi.fn(),
}))

describe('progress.queries', () => {
  beforeEach(() => vi.clearAllMocks())

  it('upcomingQuestionsOptions', async () => {
    vi.mocked(progressService.getUpcomingReviews).mockResolvedValue({ data: [] } as never)
    const opts = upcomingQuestionsOptions({ page: 1, perPage: 10 })
    expect(opts.queryKey).toEqual(['progress', 'upcoming', { page: 1, perPage: 10 }])
    await (opts as unknown as { queryFn: (ctx: unknown) => Promise<unknown> }).queryFn({} as never)
    expect(progressService.getUpcomingReviews).toHaveBeenCalledWith(1, 10)
  })

  it('savedQuestionsOptions', async () => {
    vi.mocked(progressService.getSavedQuestions).mockResolvedValue({ data: [] } as never)
    const opts = savedQuestionsOptions({ page: 2, perPage: 20 })
    await (opts as unknown as { queryFn: (ctx: unknown) => Promise<unknown> }).queryFn({} as never)
    expect(progressService.getSavedQuestions).toHaveBeenCalledWith(2, 20)
  })

  it('submitProgressAnswerMutation', async () => {
    vi.mocked(progressService.submitProgressAnswer).mockResolvedValue(undefined as never)
    const mut = submitProgressAnswerMutation()
    expect(mut.mutationKey).toEqual(['progress', 'answer'])
    await mut.mutationFn({ questionId: 'q1', isCorrect: true } as never)
    expect(progressService.submitProgressAnswer).toHaveBeenCalledWith('q1', true)
  })

  it('toggleSaveMutation', () => {
    const mut = toggleSaveMutation()
    expect(mut.mutationKey).toEqual(['progress', 'toggle-save'])
    expect(mut.mutationFn).toBe(progressService.toggleSaveQuestion)
  })
})
