import { useFormErrors } from './useFormErrors'

describe('useFormErrors', () => {
  const { extractFieldErrors } = useFormErrors()

  it('parses semicolon-separated detail', () => {
    const err = { detail: 'email: ya existe; password: muy corta' }
    expect(extractFieldErrors(err)).toEqual({
      email: 'ya existe',
      password: 'muy corta',
    })
  })

  it('parses comma-separated detail', () => {
    const err = { detail: 'email: inválido, name: requerido' }
    expect(extractFieldErrors(err)).toEqual({
      email: 'inválido',
      name: 'requerido',
    })
  })

  it('lowercases field keys', () => {
    const err = { detail: 'Email: bad' }
    expect(extractFieldErrors(err)).toEqual({ email: 'bad' })
  })

  it('handles value containing colon', () => {
    const err = { detail: 'field: a:b:c' }
    expect(extractFieldErrors(err)).toEqual({ field: 'a:b:c' })
  })

  it('trims spaces', () => {
    const err = { detail: ' email :  spaced  ;  name :  ok ' }
    expect(extractFieldErrors(err)).toEqual({ email: 'spaced', name: 'ok' })
  })

  it('returns empty for detail empty', () => {
    expect(extractFieldErrors({ detail: '' })).toEqual({})
  })

  it('returns empty for missing detail', () => {
    expect(extractFieldErrors({})).toEqual({})
    expect(extractFieldErrors(null)).toEqual({})
    expect(extractFieldErrors(undefined)).toEqual({})
    expect(extractFieldErrors('string')).toEqual({})
  })

  it('ignores parts without colon', () => {
    const err = { detail: 'email ya existe; password: corta' }
    expect(extractFieldErrors(err)).toEqual({ password: 'corta' })
  })
})
