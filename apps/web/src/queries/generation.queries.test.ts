import {
  createGenerationJobMutation,
  generationJobsListOptions,
  generationStatsOptions,
  settingsListOptions,
  updateSettingsMutation,
} from './generation.queries'
import * as generationService from '@/api/services/generation.service'
import type { GenerationJob, GenerationStats } from '@/types/generation.types'

vi.mock('@/api/services/generation.service', () => ({
  createGenerationJob: vi.fn(),
  listGenerationJobs: vi.fn(),
  getGenerationStats: vi.fn(),
  listSettings: vi.fn(),
  updateSettings: vi.fn(),
}))

type QueryOptions = { queryFn: (ctx: never) => Promise<unknown> }

function job(status: GenerationJob['status']): GenerationJob {
  return {
    id: 'j1',
    requestedBy: 'u1',
    topicIds: [],
    topics: [],
    difficulty: 'beginner',
    quantity: 5,
    status,
    generated: 0,
    createdAt: '2026-01-01T00:00:00Z',
  }
}

function stats(overrides: Partial<GenerationStats> = {}): GenerationStats {
  return {
    generatedToday: 0,
    reservedToday: 0,
    dailyCap: 200,
    perJobCap: 10,
    defaultQuantity: 5,
    pending: 0,
    running: 0,
    failed: 0,
    aiConfigured: true,
    ...overrides,
  }
}

describe('generation.queries', () => {
  beforeEach(() => vi.clearAllMocks())

  it('generationJobsListOptions calls the service', async () => {
    vi.mocked(generationService.listGenerationJobs).mockResolvedValue({ data: [] } as never)
    const opts = generationJobsListOptions(
      () => 1,
      () => 10,
      () => undefined,
    )
    await (opts as unknown as QueryOptions).queryFn({} as never)
    expect(generationService.listGenerationJobs).toHaveBeenCalledWith(1, 10, undefined)
  })

  it('generationJobsListOptions passes the status filter', async () => {
    vi.mocked(generationService.listGenerationJobs).mockResolvedValue({ data: [] } as never)
    const opts = generationJobsListOptions(
      () => 2,
      () => 25,
      () => 'failed',
    )
    await (opts as unknown as QueryOptions).queryFn({} as never)
    expect(generationService.listGenerationJobs).toHaveBeenCalledWith(2, 25, 'failed')
  })

  it('generationJobsListOptions polls while a job is in flight', () => {
    const opts = generationJobsListOptions(
      () => 1,
      () => 10,
      () => undefined,
    ) as unknown as {
      refetchInterval: (query: { state: { data?: unknown } }) => number | false
    }

    const pending = { state: { data: { data: [job('pending')] } } }
    const done = { state: { data: { data: [job('completed')] } } }

    expect(opts.refetchInterval(pending)).toBe(3000)
    expect(opts.refetchInterval(done)).toBe(false)
  })

  it('generationStatsOptions polls slowly when idle and fast when busy', () => {
    const opts = generationStatsOptions() as unknown as {
      refetchInterval: (query: { state: { data?: unknown } }) => number | false
    }

    expect(opts.refetchInterval({ state: { data: stats({ running: 1 }) } })).toBe(3000)
    expect(opts.refetchInterval({ state: { data: stats() } })).toBe(30_000)
  })

  it('createGenerationJobMutation keys on generation/create', async () => {
    vi.mocked(generationService.createGenerationJob).mockResolvedValue(job('pending'))
    const mut = createGenerationJobMutation()
    expect(mut.mutationKey).toEqual(['generation', 'create'])

    const res = await mut.mutationFn({
      topicIds: ['t1'],
      difficulty: 'beginner',
      quantity: 3,
    })
    expect(generationService.createGenerationJob).toHaveBeenCalledWith({
      topicIds: ['t1'],
      difficulty: 'beginner',
      quantity: 3,
    })
    expect(res.status).toBe('pending')
  })

  it('settingsListOptions calls the service', async () => {
    vi.mocked(generationService.listSettings).mockResolvedValue([])
    const opts = settingsListOptions()
    await (opts as unknown as QueryOptions).queryFn({} as never)
    expect(generationService.listSettings).toHaveBeenCalled()
  })

  it('updateSettingsMutation keys on settings/update', async () => {
    vi.mocked(generationService.updateSettings).mockResolvedValue([])
    const mut = updateSettingsMutation()
    expect(mut.mutationKey).toEqual(['settings', 'update'])

    await mut.mutationFn({ settings: [{ key: 'ai.max_per_job', value: '8' }] })
    expect(generationService.updateSettings).toHaveBeenCalledWith({
      settings: [{ key: 'ai.max_per_job', value: '8' }],
    })
  })
})
