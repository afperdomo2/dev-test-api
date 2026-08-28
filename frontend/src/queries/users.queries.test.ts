/* eslint-disable @typescript-eslint/no-non-null-assertion */
import {
  usersListOptions,
  userDetailOptions,
  createUserMutation,
  updateUserMutation,
  deleteUserMutation,
  profileQueryOptions,
} from './users.queries'
import * as usersService from '@/api/services/users.service'

vi.mock('@/api/services/users.service', () => ({
  listUsers: vi.fn(),
  getUserById: vi.fn(),
  createUser: vi.fn(),
  updateUser: vi.fn(),
  deleteUser: vi.fn(),
  getProfile: vi.fn(),
}))

describe('users.queries', () => {
  beforeEach(() => vi.clearAllMocks())

  it('usersListOptions', async () => {
    const data = { data: [] }
    vi.mocked(usersService.listUsers).mockResolvedValue(data as never)
    const opts = usersListOptions(() => 2, () => 10)
    expect(opts.queryKey).toEqual(['users', 'list', expect.any(Function), expect.any(Function)])
    const result = await (opts as unknown as { queryFn: (ctx: unknown) => Promise<unknown> }).queryFn({} as never)
    expect(usersService.listUsers).toHaveBeenCalledWith(2, 10)
    expect(result).toEqual(data)
  })

  it('userDetailOptions', async () => {
    vi.mocked(usersService.getUserById).mockResolvedValue({ id: '1' } as never)
    const opts = userDetailOptions(() => '1')
    const result = await (opts as unknown as { queryFn: (ctx: unknown) => Promise<unknown> }).queryFn({} as never)
    expect(usersService.getUserById).toHaveBeenCalledWith('1')
    expect(result).toEqual({ id: '1' })
  })

  it('createUserMutation', async () => {
    vi.mocked(usersService.createUser).mockResolvedValue({ id: '1' } as never)
    const mut = createUserMutation()
    expect(mut.mutationKey).toEqual(['users', 'create'])
    const res = await mut.mutationFn({ email: 'a@b.com' } as never)
    expect(usersService.createUser).toHaveBeenCalled()
    expect(res).toEqual({ id: '1' })
  })

  it('updateUserMutation', async () => {
    vi.mocked(usersService.updateUser).mockResolvedValue({ id: '1' } as never)
    const mut = updateUserMutation()
    expect(mut.mutationKey).toEqual(['users', 'update'])
    const res = await mut.mutationFn({ id: '1', data: { isAdmin: true } } as never)
    expect(usersService.updateUser).toHaveBeenCalledWith('1', { isAdmin: true })
    expect(res).toEqual({ id: '1' })
  })

  it('deleteUserMutation', async () => {
    vi.mocked(usersService.deleteUser).mockResolvedValue(undefined as never)
    const mut = deleteUserMutation()
    expect(mut.mutationKey).toEqual(['users', 'delete'])
    await mut.mutationFn('1' as never)
    expect(usersService.deleteUser).toHaveBeenCalledWith('1')
  })

  it('profileQueryOptions', () => {
    const opts = profileQueryOptions()
    expect(opts.queryKey).toEqual(['users', 'profile'])
    expect(opts.staleTime).toBe(5 * 60 * 1000)
  })
})
