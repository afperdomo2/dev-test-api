import { formatDate, formatDateTime, formatScore } from './format'

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
})
