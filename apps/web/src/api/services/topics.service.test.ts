import { listTopics, getTopicById, createTopic, updateTopic, deleteTopic } from './topics.service'
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

describe('topics.service', () => {
  beforeEach(() => vi.clearAllMocks())

  it('listTopics with only required params', async () => {
    const paginated = { data: [] }
    mockedGet.mockResolvedValue({ data: paginated })
    await listTopics(1, 10)
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/topics', { params: { page: 1, perPage: 10 } })
  })

  it('listTopics with all optional params', async () => {
    mockedGet.mockResolvedValue({ data: {} })
    await listTopics(1, 10, 'name', 'asc', 'search term', 'backend', true)
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/topics', {
      params: {
        page: 1,
        perPage: 10,
        sortBy: 'name',
        sortOrder: 'asc',
        search: 'search term',
        category: 'backend',
        myOnly: true,
      },
    })
  })

  it('listTopics omits falsy optionals', async () => {
    mockedGet.mockResolvedValue({ data: {} })
    await listTopics(1, 10, undefined, undefined, '', '', false)
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/topics', { params: { page: 1, perPage: 10 } })
  })

  it('getTopicById', async () => {
    mockedGet.mockResolvedValue({ data: { id: '1' } })
    const res = await getTopicById('1')
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/topics/1')
    expect(res).toEqual({ id: '1' })
  })

  it('createTopic', async () => {
    mockedPost.mockResolvedValue({ data: { id: '1' } })
    const res = await createTopic({ name: 't' } as never)
    expect(mockedPost).toHaveBeenCalledWith('/api/v1/topics', { name: 't' })
    expect(res).toEqual({ id: '1' })
  })

  it('updateTopic', async () => {
    mockedPut.mockResolvedValue({ data: { id: '1' } })
    const res = await updateTopic('1', { name: 'x' } as never)
    expect(mockedPut).toHaveBeenCalledWith('/api/v1/topics/1', { name: 'x' })
    expect(res).toEqual({ id: '1' })
  })

  it('deleteTopic', async () => {
    mockedDelete.mockResolvedValue({})
    await deleteTopic('1')
    expect(mockedDelete).toHaveBeenCalledWith('/api/v1/topics/1')
  })
})
