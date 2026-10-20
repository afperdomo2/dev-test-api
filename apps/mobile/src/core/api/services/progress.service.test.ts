import { apiClient } from '@/core/api/client'
import { getSavedQuestions, getTopicMastery, getUpcomingReviews, toggleSaveQuestion } from './progress.service'

jest.mock('@/core/api/client', () => ({
  apiClient: { get: jest.fn(), post: jest.fn(), put: jest.fn(), delete: jest.fn() },
}))

const mocked = apiClient as unknown as { get: jest.Mock; post: jest.Mock }

describe('progress.service', () => {
  beforeEach(() => jest.clearAllMocks())

  it('getUpcomingReviews', async () => {
    mocked.get.mockResolvedValue({ data: [], meta: { total: 0 } })

    await getUpcomingReviews(1, 10)

    expect(mocked.get).toHaveBeenCalledWith('/api/v1/progress/upcoming', {
      page: 1,
      perPage: 10,
    })
  })

  it('getSavedQuestions', async () => {
    mocked.get.mockResolvedValue({ data: [], meta: { total: 0 } })

    await getSavedQuestions(2, 20)

    expect(mocked.get).toHaveBeenCalledWith('/api/v1/progress/saved', {
      page: 2,
      perPage: 20,
    })
  })

  it('toggleSaveQuestion', async () => {
    mocked.post.mockResolvedValue({ questionId: 'q1', isSaved: true })

    await expect(toggleSaveQuestion('q1')).resolves.toEqual({
      questionId: 'q1',
      isSaved: true,
    })
    expect(mocked.post).toHaveBeenCalledWith('/api/v1/progress/q1/toggle-save')
  })

  it('getTopicMastery', async () => {
    const items = [{ topicId: 't1', topicName: 'Go', attempts: 5, correctCount: 2, accuracy: 0.4, isWeak: true }]
    mocked.get.mockResolvedValue(items)

    await expect(getTopicMastery()).resolves.toEqual(items)
    expect(mocked.get).toHaveBeenCalledWith('/api/v1/progress/topics-mastery')
  })
})
