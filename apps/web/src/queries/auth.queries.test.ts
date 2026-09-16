import { loginMutation, setupMutation } from './auth.queries'
import * as authService from '@/api/services/auth.service'

vi.mock('@/api/services/auth.service', () => ({
  login: vi.fn(),
  setup: vi.fn(),
  getStatus: vi.fn(),
}))

describe('auth.queries', () => {
  beforeEach(() => vi.clearAllMocks())

  it('loginMutation has correct key and calls login', async () => {
    const payload = { email: 'a@b.com', password: 'x' }
    const response = { token: 't' }
    vi.mocked(authService.login).mockResolvedValue(response as never)

    const mut = loginMutation()
    expect(mut.mutationKey).toEqual(['auth', 'login'])

    const result = await mut.mutationFn(payload as never)
    expect(authService.login).toHaveBeenCalledWith(payload)
    expect(result).toEqual(response)
  })

  it('setupMutation has correct key and calls setup', async () => {
    const payload = { email: 'a@b.com', password: 'x' }
    const response = { token: 't' }
    vi.mocked(authService.setup).mockResolvedValue(response as never)

    const mut = setupMutation()
    expect(mut.mutationKey).toEqual(['auth', 'setup'])

    const result = await mut.mutationFn(payload as never)
    expect(authService.setup).toHaveBeenCalledWith(payload)
    expect(result).toEqual(response)
  })
})
