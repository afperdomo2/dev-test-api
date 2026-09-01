import {
  submitProgressAnswer,
  getUpcomingReviews,
  getSavedQuestions,
  toggleSaveQuestion,
} from './progress.service'
import apiClient from '@/api/client'

vi.mock('@/api/client', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
  },
}))

const mockedGet = vi.mocked(apiClient.get)
const mockedPost = vi.mocked(apiClient.post)

describe('progress.service', () => {
  beforeEach(() => vi.clearAllMocks())

  it('submitProgressAnswer posts isCorrect', async () => {
    mockedPost.mockResolvedValue({})
    await submitProgressAnswer('q1', true)
    expect(mockedPost).toHaveBeenCalledWith('/api/v1/progress/q1/answer', { isCorrect: true })
  })

  it('getUpcomingReviews', async () => {
    const paginated = { data: [] }
    mockedGet.mockResolvedValue({ data: paginated })
    const res = await getUpcomingReviews(1, 10)
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/progress/upcoming', {
      params: { page: 1, perPage: 10 },
    })
    expect(res).toEqual(paginated)
  })

  it('getSavedQuestions', async () => {
    const paginated = { data: [] }
    mockedGet.mockResolvedValue({ data: paginated })
    const res = await getSavedQuestions(2, 20)
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/progress/saved', {
      params: { page: 2, perPage: 20 },
    })
    expect(res).toEqual(paginated)
  })

  it('toggleSaveQuestion', async () => {
    mockedPost.mockResolvedValue({ data: { id: 'p1' } })
    const res = await toggleSaveQuestion('q1')
    expect(mockedPost).toHaveBeenCalledWith('/api/v1/progress/q1/toggle-save')
    expect(res).toEqual({ id: 'p1' })
  })
})
