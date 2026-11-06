import {
  formatDate,
  formatDateTime,
  formatDuration,
  formatRelativeTime,
  formatScore,
} from './format'

describe('format utils', () => {
  describe('formatScore', () => {
    it('formats with one decimal + %', () => {
      expect(formatScore(0)).toBe('0.0%')
      expect(formatScore(100)).toBe('100.0%')
      expect(formatScore(75.567)).toBe('75.6%')
      expect(formatScore(33.333)).toBe('33.3%')
    })
  })

  describe('formatDate / formatDateTime', () => {
    it('returns non-empty string for valid ISO date', () => {
      const iso = '2026-08-28T12:34:00.000Z'
      const d = formatDate(iso)
      const dt = formatDateTime(iso)
      expect(typeof d).toBe('string')
      expect(d.length).toBeGreaterThan(0)
      expect(typeof dt).toBe('string')
      expect(dt.length).toBeGreaterThan(0)
    })

    it('formatDate and formatDateTime differ (time included)', () => {
      const iso = '2026-01-15T10:20:30.000Z'
      const d = formatDate(iso)
      const dt = formatDateTime(iso)
      expect(dt).not.toBe(d)
    })
  })

  describe('formatRelativeTime', () => {
    beforeEach(() => {
      vi.useFakeTimers()
      vi.setSystemTime(new Date('2026-10-25T12:00:00.000Z'))
    })

    afterEach(() => {
      vi.useRealTimers()
    })

    it('handles invalid dates', () => {
      expect(formatRelativeTime('not-a-date')).toBe('')
    })

    it('formats minutes, hours, days and older dates', () => {
      expect(formatRelativeTime('2026-10-25T11:59:40.000Z')).toBe('Ahora')
      expect(formatRelativeTime('2026-10-25T11:45:00.000Z')).toBe('Hace 15 min')
      expect(formatRelativeTime('2026-10-25T09:00:00.000Z')).toBe('Hace 3h')
      expect(formatRelativeTime('2026-10-24T12:00:00.000Z')).toBe('Ayer')
      expect(formatRelativeTime('2026-10-22T12:00:00.000Z')).toBe('Hace 3 días')
      expect(formatRelativeTime('2026-09-15T12:00:00.000Z')).toBe('15 sept')
      expect(formatRelativeTime('2025-09-15T12:00:00.000Z')).toBe('15 sept 2025')
    })
  })

  describe('formatDuration', () => {
    it('returns empty string without finishedAt', () => {
      expect(formatDuration('2026-10-25T10:00:00.000Z', null)).toBe('')
    })

    it('returns empty string for invalid ranges', () => {
      expect(formatDuration('2026-10-25T10:00:00.000Z', '2026-10-25T09:00:00.000Z')).toBe('')
      expect(formatDuration('nope', '2026-10-25T09:00:00.000Z')).toBe('')
    })

    it('formats minutes and hours', () => {
      expect(formatDuration('2026-10-25T10:00:00.000Z', '2026-10-25T10:34:00.000Z')).toBe('34 min')
      expect(formatDuration('2026-10-25T10:00:00.000Z', '2026-10-25T11:12:00.000Z')).toBe(
        '1 h 12 min',
      )
      expect(formatDuration('2026-10-25T10:00:00.000Z', '2026-10-25T12:00:00.000Z')).toBe('2 h')
    })
  })
})
