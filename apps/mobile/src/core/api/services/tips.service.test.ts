import { apiClient } from '@/core/api/client'
import { getTipOfTheDay } from './tips.service'

jest.mock('@/core/api/client', () => ({
  apiClient: { get: jest.fn(), post: jest.fn(), put: jest.fn(), delete: jest.fn() },
}))

const mocked = apiClient as unknown as { get: jest.Mock }

describe('tips.service', () => {
  beforeEach(() => jest.clearAllMocks())

  it('getTipOfTheDay pide el endpoint del día', async () => {
    const tip = { id: 't1', topic: 'Go', content: 'tip' }
    mocked.get.mockResolvedValue(tip)

    await expect(getTipOfTheDay()).resolves.toEqual(tip)
    expect(mocked.get).toHaveBeenCalledWith('/api/v1/tips/today')
  })
})
