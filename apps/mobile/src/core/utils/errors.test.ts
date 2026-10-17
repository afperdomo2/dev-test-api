import { getErrorDetail } from './errors'

describe('getErrorDetail', () => {
  it('returns the detail when the error carries one', () => {
    expect(getErrorDetail({ detail: 'Algo salió mal' }, 'fallback')).toBe('Algo salió mal')
  })

  it('stringifies non-string details', () => {
    expect(getErrorDetail({ detail: 404 }, 'fallback')).toBe('404')
  })

  it('falls back for nullish and primitive errors', () => {
    expect(getErrorDetail(null, 'fallback')).toBe('fallback')
    expect(getErrorDetail(undefined, 'fallback')).toBe('fallback')
    expect(getErrorDetail('boom', 'fallback')).toBe('fallback')
  })

  it('falls back when the object has no detail key', () => {
    expect(getErrorDetail({ title: 'Not Found' }, 'fallback')).toBe('fallback')
  })
})
