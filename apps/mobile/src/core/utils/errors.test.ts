import { getErrorDetail, isNetworkError } from './errors'

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

describe('isNetworkError', () => {
  it('detecta el error de red de apiClient', () => {
    expect(
      isNetworkError({ title: 'Network Error', detail: 'No se pudo conectar con el servidor' }),
    ).toBe(true)
  })

  it('un error HTTP no es error de red', () => {
    expect(isNetworkError({ title: 'Not Found', status: 404, detail: 'No existe' })).toBe(false)
    expect(isNetworkError({ title: 'Internal Server Error', status: 500 })).toBe(false)
  })

  it('un Error genérico o un valor vacío no son error de red', () => {
    expect(isNetworkError(new Error('boom'))).toBe(false)
    expect(isNetworkError(null)).toBe(false)
    expect(isNetworkError(undefined)).toBe(false)
    expect(isNetworkError('Network Error')).toBe(false)
  })
})
