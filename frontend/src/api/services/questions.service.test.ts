import {
  listQuestions,
  getQuestionById,
  createQuestion,
  updateQuestion,
  deleteQuestion,
} from './questions.service'
import apiClient from '@/api/client'

vi.mock('@/api/client', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}))

const mockedGet = vi.mocked(apiClient.get)
const mockedPost = vi.mocked(apiClient.post)
const mockedPut = vi.mocked(apiClient.put)
const mockedDelete = vi.mocked(apiClient.delete)

describe('questions.service', () => {
  beforeEach(() => vi.clearAllMocks())

  it('listQuestions without filters', async () => {
    mockedGet.mockResolvedValue({ data: { data: [] } })
    await listQuestions(1, 10)
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/questions', {
      params: { page: 1, perPage: 10 },
    })
  })

  it('listQuestions with filters', async () => {
    mockedGet.mockResolvedValue({ data: {} })
    await listQuestions(1, 10, { type: 'mcq', difficulty: 'hard' })
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/questions', {
      params: { page: 1, perPage: 10, type: 'mcq', difficulty: 'hard' },
    })
  })

  it('getQuestionById', async () => {
    mockedGet.mockResolvedValue({ data: { id: '1' } })
    const res = await getQuestionById('1')
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/questions/1')
    expect(res).toEqual({ id: '1' })
  })

  it('createQuestion', async () => {
    mockedPost.mockResolvedValue({ data: { id: '1' } })
    const res = await createQuestion({ text: 'q' } as never)
    expect(mockedPost).toHaveBeenCalledWith('/api/v1/questions', { text: 'q' })
    expect(res).toEqual({ id: '1' })
  })

  it('updateQuestion', async () => {
    mockedPut.mockResolvedValue({ data: { id: '1' } })
    const res = await updateQuestion('1', { text: 'x' } as never)
    expect(mockedPut).toHaveBeenCalledWith('/api/v1/questions/1', { text: 'x' })
    expect(res).toEqual({ id: '1' })
  })

  it('deleteQuestion', async () => {
    mockedDelete.mockResolvedValue({})
    await deleteQuestion('1')
    expect(mockedDelete).toHaveBeenCalledWith('/api/v1/questions/1')
  })
})
