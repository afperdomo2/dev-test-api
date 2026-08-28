import {
  listUsers,
  getUserById,
  createUser,
  updateUser,
  deleteUser,
  getProfile,
} from './users.service'
import apiClient from '@/api/client'

vi.mock('@/api/client', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}))

const mockedGet = vi.mocked(apiClient.get)
const mockedPost = vi.mocked(apiClient.post)
const mockedPut = vi.mocked(apiClient.put)
const mockedDelete = vi.mocked(apiClient.delete)

describe('users.service', () => {
  beforeEach(() => vi.clearAllMocks())

  it('listUsers gets with params', async () => {
    const paginated = { data: [], meta: { total: 0, page: 1, perPage: 10 } }
    mockedGet.mockResolvedValue({ data: paginated })
    const result = await listUsers(2, 20)
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/users', { params: { page: 2, perPage: 20 } })
    expect(result).toEqual(paginated)
  })

  it('getUserById', async () => {
    const user = { id: '1' }
    mockedGet.mockResolvedValue({ data: user })
    const result = await getUserById('1')
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/users/1')
    expect(result).toEqual(user)
  })

  it('createUser posts', async () => {
    const user = { id: '1' }
    const payload = { email: 'a@b.com', password: 'secret123' }
    mockedPost.mockResolvedValue({ data: user })
    const result = await createUser(payload as never)
    expect(mockedPost).toHaveBeenCalledWith('/api/v1/users', payload)
    expect(result).toEqual(user)
  })

  it('updateUser puts', async () => {
    const user = { id: '1' }
    mockedPut.mockResolvedValue({ data: user })
    const result = await updateUser('1', { isAdmin: true } as never)
    expect(mockedPut).toHaveBeenCalledWith('/api/v1/users/1', { isAdmin: true })
    expect(result).toEqual(user)
  })

  it('deleteUser deletes', async () => {
    mockedDelete.mockResolvedValue({})
    await deleteUser('1')
    expect(mockedDelete).toHaveBeenCalledWith('/api/v1/users/1')
  })

  it('getProfile', async () => {
    const user = { id: '1' }
    mockedGet.mockResolvedValue({ data: user })
    const result = await getProfile()
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/profile')
    expect(result).toEqual(user)
  })
})
