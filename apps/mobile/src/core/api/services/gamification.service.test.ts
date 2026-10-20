import { apiClient } from '@/core/api/client'
import { freezeStreak, getOverview, setDailyGoal } from './gamification.service'

jest.mock('@/core/api/client', () => ({
  apiClient: { get: jest.fn(), post: jest.fn(), put: jest.fn(), delete: jest.fn() },
}))

const mocked = apiClient as unknown as {
  get: jest.Mock
  post: jest.Mock
  put: jest.Mock
}

describe('gamification.service', () => {
  beforeEach(() => jest.clearAllMocks())

  it('getOverview', async () => {
    mocked.get.mockResolvedValue({ totalXp: 10 })

    await expect(getOverview()).resolves.toEqual({ totalXp: 10 })
    expect(mocked.get).toHaveBeenCalledWith('/api/v1/gamification/overview')
  })

  it('setDailyGoal hace PUT con el body', async () => {
    mocked.put.mockResolvedValue({ xpGained: 0 })
    const body = { dailyGoalXp: 100 }

    await setDailyGoal(body)

    expect(mocked.put).toHaveBeenCalledWith('/api/v1/gamification/daily-goal', body)
  })

  it('freezeStreak hace POST sin body', async () => {
    mocked.post.mockResolvedValue({ xpGained: 0 })

    await freezeStreak()

    expect(mocked.post).toHaveBeenCalledWith('/api/v1/gamification/streak/freeze')
  })
})
