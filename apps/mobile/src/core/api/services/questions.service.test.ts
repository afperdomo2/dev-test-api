import { apiClient } from '@/core/api/client'
import { getQuestionById, listQuestions } from './questions.service'

jest.mock('@/core/api/client', () => ({
  apiClient: { get: jest.fn(), post: jest.fn(), put: jest.fn(), delete: jest.fn() },
}))

const mocked = apiClient as unknown as { get: jest.Mock }

describe('questions.service', () => {
  beforeEach(() => jest.clearAllMocks())

  it('listQuestions combina paginación y filtros', async () => {
    mocked.get.mockResolvedValue({ data: [], meta: { total: 0 } })

    await listQuestions(1, 20, { type: 'true_false', difficulty: 'beginner' })

    expect(mocked.get).toHaveBeenCalledWith('/api/v1/questions', {
      page: 1,
      perPage: 20,
      type: 'true_false',
      difficulty: 'beginner',
    })
  })

  it('listQuestions sin filtros envía solo paginación', async () => {
    mocked.get.mockResolvedValue({ data: [], meta: { total: 0 } })

    await listQuestions(1, 20)

    expect(mocked.get).toHaveBeenCalledWith('/api/v1/questions', { page: 1, perPage: 20 })
  })

  it('getQuestionById', async () => {
    mocked.get.mockResolvedValue({ id: 'q1' })

    await expect(getQuestionById('q1')).resolves.toEqual({ id: 'q1' })
    expect(mocked.get).toHaveBeenCalledWith('/api/v1/questions/q1')
  })
})
