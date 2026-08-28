import { getToken, setToken, removeToken } from './storage'

describe('storage utils', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('getToken returns null initially', () => {
    expect(getToken()).toBeNull()
  })

  it('setToken + getToken round-trip', () => {
    setToken('abc123')
    expect(getToken()).toBe('abc123')
    expect(localStorage.getItem('auth_token')).toBe('abc123')
  })

  it('overwrites existing token', () => {
    setToken('first')
    setToken('second')
    expect(getToken()).toBe('second')
  })

  it('removeToken clears', () => {
    setToken('tok')
    removeToken()
    expect(getToken()).toBeNull()
    expect(localStorage.getItem('auth_token')).toBeNull()
  })

  it('removeToken is idempotent when empty', () => {
    expect(() => removeToken()).not.toThrow()
    expect(getToken()).toBeNull()
  })
})
