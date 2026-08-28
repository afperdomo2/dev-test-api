/* eslint-disable @typescript-eslint/no-non-null-assertion */
import {
  questionsListOptions,
  questionDetailOptions,
  createQuestionMutation,
  updateQuestionMutation,
  deleteQuestionMutation,
} from './questions.queries'
import * as questionsService from '@/api/services/questions.service'

vi.mock('@/api/services/questions.service', () => ({
  listQuestions: vi.fn(),
  getQuestionById: vi.fn(),
  createQuestion: vi.fn(),
  updateQuestion: vi.fn(),
  deleteQuestion: vi.fn(),
}))

describe('questions.queries', () => {
  beforeEach(() => vi.clearAllMocks())

  it('questionsListOptions', async () => {
    vi.mocked(questionsService.listQuestions).mockResolvedValue({ data: [] } as never)
    const opts = questionsListOptions(() => 1, () => 10, () => ({ type: 'mcq' }))
    await (opts as unknown as { queryFn: (ctx: unknown) => Promise<unknown> }).queryFn({} as never)
    expect(questionsService.listQuestions).toHaveBeenCalledWith(1, 10, { type: 'mcq' })
  })

  it('questionDetailOptions', async () => {
    vi.mocked(questionsService.getQuestionById).mockResolvedValue({ id: '1' } as never)
    const opts = questionDetailOptions(() => '1')
    const res = await (opts as unknown as { queryFn: (ctx: unknown) => Promise<unknown> }).queryFn({} as never)
    expect(questionsService.getQuestionById).toHaveBeenCalledWith('1')
    expect(res).toEqual({ id: '1' })
  })

  it('createQuestionMutation', () => {
    const mut = createQuestionMutation()
    expect(mut.mutationKey).toEqual(['questions', 'create'])
    expect(mut.mutationFn).toBe(questionsService.createQuestion)
  })

  it('updateQuestionMutation', async () => {
    vi.mocked(questionsService.updateQuestion).mockResolvedValue({ id: '1' } as never)
    const mut = updateQuestionMutation()
    const res = await mut.mutationFn({ id: '1', data: { text: 'x' } as never } as never)
    expect(questionsService.updateQuestion).toHaveBeenCalledWith('1', { text: 'x' })
    expect(res).toEqual({ id: '1' })
  })

  it('deleteQuestionMutation', () => {
    const mut = deleteQuestionMutation()
    expect(mut.mutationKey).toEqual(['questions', 'delete'])
  })
})
