import * as SecureStore from 'expo-secure-store'

import {
  getTipDismissedOn,
  isTipDismissedToday,
  setTipDismissedOn,
  todayUTC,
} from './dailyTipDismissal'

jest.mock('expo-secure-store', () => ({
  getItemAsync: jest.fn(),
  setItemAsync: jest.fn(),
}))

const store = SecureStore as unknown as {
  getItemAsync: jest.Mock
  setItemAsync: jest.Mock
}

describe('dailyTipDismissal', () => {
  beforeEach(() => jest.clearAllMocks())

  it('todayUTC usa la fecha UTC', () => {
    expect(todayUTC(new Date('2026-11-09T23:30:00Z'))).toBe('2026-11-09')
    expect(todayUTC(new Date('2026-11-10T00:10:00Z'))).toBe('2026-11-10')
  })

  it('guarda y lee el día del cierre', async () => {
    store.getItemAsync.mockResolvedValue(null)
    await expect(getTipDismissedOn()).resolves.toBeNull()

    await setTipDismissedOn('2026-11-09')
    expect(store.setItemAsync).toHaveBeenCalledWith('daily_tip_dismissed_on', '2026-11-09')
  })

  it('solo está cerrado si la fecha guardada es hoy', async () => {
    const now = new Date('2026-11-09T12:00:00Z')

    store.getItemAsync.mockResolvedValue('2026-11-09')
    await expect(isTipDismissedToday(now)).resolves.toBe(true)

    store.getItemAsync.mockResolvedValue('2026-11-08')
    await expect(isTipDismissedToday(now)).resolves.toBe(false)
  })
})
