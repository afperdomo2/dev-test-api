import { getOverview, setDailyGoal, freezeStreak, getBadges } from './gamification.service'
import apiClient from '@/api/client'

vi.mock('@/api/client', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
  },
}))

const mockedGet = vi.mocked(apiClient.get)
const mockedPost = vi.mocked(apiClient.post)
const mockedPut = vi.mocked(apiClient.put)

describe('gamification.service', () => {
  beforeEach(() => vi.clearAllMocks())

  it('getOverview', async () => {
    const overview = { totalXp: 120, level: 2 }
    mockedGet.mockResolvedValue({ data: overview })
    const res = await getOverview()
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/gamification/overview')
    expect(res).toEqual(overview)
  })

  it('setDailyGoal puts dailyGoalXp', async () => {
    const reward = { totalXp: 120 }
    mockedPut.mockResolvedValue({ data: reward })
    const res = await setDailyGoal({ dailyGoalXp: 50 })
    expect(mockedPut).toHaveBeenCalledWith('/api/v1/gamification/daily-goal', { dailyGoalXp: 50 })
    expect(res).toEqual(reward)
  })

  it('freezeStreak posts without body', async () => {
    const reward = { totalXp: 120 }
    mockedPost.mockResolvedValue({ data: reward })
    const res = await freezeStreak()
    expect(mockedPost).toHaveBeenCalledWith('/api/v1/gamification/streak/freeze')
    expect(res).toEqual(reward)
  })

  it('getBadges', async () => {
    const badges = { badges: [{ id: 'first_session' }], earned: 1, total: 8 }
    mockedGet.mockResolvedValue({ data: badges })

    const res = await getBadges()

    expect(mockedGet).toHaveBeenCalledWith('/api/v1/gamification/badges')
    expect(res).toEqual(badges)
  })
})
