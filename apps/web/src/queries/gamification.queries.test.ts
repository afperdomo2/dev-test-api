import {
  overviewOptions,
  setDailyGoalMutation,
  freezeStreakMutation,
  badgesOptions,
} from './gamification.queries'
import * as gamificationService from '@/api/services/gamification.service'

vi.mock('@/api/services/gamification.service', () => ({
  getOverview: vi.fn(),
  setDailyGoal: vi.fn(),
  freezeStreak: vi.fn(),
  getBadges: vi.fn(),
}))

describe('gamification.queries', () => {
  beforeEach(() => vi.clearAllMocks())

  it('overviewOptions calls service', async () => {
    vi.mocked(gamificationService.getOverview).mockResolvedValue({ totalXp: 10 } as never)
    const opts = overviewOptions()
    const res = await (opts as unknown as { queryFn: (ctx: unknown) => Promise<unknown> }).queryFn(
      {} as never,
    )
    expect(gamificationService.getOverview).toHaveBeenCalled()
    expect(res).toEqual({ totalXp: 10 })
  })

  it('setDailyGoalMutation', () => {
    const mut = setDailyGoalMutation()
    expect(mut.mutationKey).toEqual(['gamification', 'daily-goal'])
    expect(mut.mutationFn).toBe(gamificationService.setDailyGoal)
  })

  it('freezeStreakMutation', () => {
    const mut = freezeStreakMutation()
    expect(mut.mutationKey).toEqual(['gamification', 'streak', 'freeze'])
    expect(mut.mutationFn).toBe(gamificationService.freezeStreak)
  })

  it('badgesOptions', async () => {
    vi.mocked(gamificationService.getBadges).mockResolvedValue({ badge: [] } as never)
    const opts = badgesOptions()
    expect(opts.queryKey).toEqual(['gamification', 'badges'])

    await (opts as unknown as { queryFn: (ctx: unknown) => Promise<unknown> }).queryFn({} as never)
    expect(gamificationService.getBadges).toHaveBeenCalled()
  })
})
