import { formatRelativeTime } from './dates'

const NOW = new Date('2026-10-17T12:00:00.000Z')

const ago = (ms: number) => new Date(NOW.getTime() - ms).toISOString()

describe('formatRelativeTime', () => {
  beforeEach(() => {
    jest.useFakeTimers()
    jest.setSystemTime(NOW)
  })

  afterEach(() => {
    jest.useRealTimers()
  })

  it('returns an empty string for an invalid date', () => {
    expect(formatRelativeTime('not-a-date')).toBe('')
  })

  it('returns "Ahora" under a minute', () => {
    expect(formatRelativeTime(ago(30_000))).toBe('Ahora')
  })

  it('returns minutes under an hour', () => {
    expect(formatRelativeTime(ago(5 * 60_000))).toBe('Hace 5 min')
  })

  it('returns hours under a day', () => {
    expect(formatRelativeTime(ago(3 * 3_600_000))).toBe('Hace 3h')
  })

  it('returns "Ayer" for a single day', () => {
    expect(formatRelativeTime(ago(24 * 3_600_000))).toBe('Ayer')
  })

  it('returns days under a week', () => {
    expect(formatRelativeTime(ago(3 * 24 * 3_600_000))).toBe('Hace 3 días')
  })

  it('returns day + short month for older dates in the current year', () => {
    expect(formatRelativeTime(new Date('2026-03-05T12:00:00.000Z').toISOString())).toBe('5 mar')
  })

  it('includes the year for older dates', () => {
    expect(formatRelativeTime(new Date('2024-03-05T12:00:00.000Z').toISOString())).toBe(
      '5 mar 2024',
    )
  })
})
