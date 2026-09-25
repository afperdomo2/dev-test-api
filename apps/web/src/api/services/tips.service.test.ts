import {
  listTips,
  getTipById,
  createTip,
  updateTip,
  deleteTip,
  importTips,
  importTipsFromText,
} from './tips.service'
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

describe('tips.service', () => {
  beforeEach(() => vi.clearAllMocks())

  it('listTips without filters', async () => {
    const paginated = { data: [], meta: { total: 0 } }
    mockedGet.mockResolvedValue({ data: paginated })
    const res = await listTips(1, 10)
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/tips', {
      params: { page: 1, perPage: 10 },
    })
    expect(res).toEqual(paginated)
  })

  it('listTips with topicId filter', async () => {
    mockedGet.mockResolvedValue({ data: { data: [], meta: { total: 0 } } })
    await listTips(2, 20, { topicId: 'topic-1' })
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/tips', {
      params: { page: 2, perPage: 20, topicId: 'topic-1' },
    })
  })

  it('getTipById', async () => {
    const tip = { id: 't1' }
    mockedGet.mockResolvedValue({ data: tip })
    const res = await getTipById('t1')
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/tips/t1')
    expect(res).toEqual(tip)
  })

  it('createTip', async () => {
    mockedPost.mockResolvedValue({ data: { id: 't1' } })
    const res = await createTip({ topicId: 'topic-1', content: 'hola' })
    expect(mockedPost).toHaveBeenCalledWith('/api/v1/tips', {
      topicId: 'topic-1',
      content: 'hola',
    })
    expect(res).toEqual({ id: 't1' })
  })

  it('updateTip', async () => {
    mockedPut.mockResolvedValue({ data: { id: 't1' } })
    const res = await updateTip('t1', { content: 'nuevo' })
    expect(mockedPut).toHaveBeenCalledWith('/api/v1/tips/t1', { content: 'nuevo' })
    expect(res).toEqual({ id: 't1' })
  })

  it('deleteTip', async () => {
    mockedDelete.mockResolvedValue({})
    await deleteTip('t1')
    expect(mockedDelete).toHaveBeenCalledWith('/api/v1/tips/t1')
  })

  it('importTips posts FormData as multipart', async () => {
    const formData = new FormData()
    formData.append('file', new Blob(['a,b']))
    mockedPost.mockResolvedValue({ data: { total: 1, imported: 1, failed: 0, errors: [] } })
    const res = await importTips(formData)
    expect(mockedPost).toHaveBeenCalledWith('/api/v1/tips/import', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    expect(res.imported).toBe(1)
  })

  it('importTipsFromText wraps content in FormData', async () => {
    mockedPost.mockResolvedValue({ data: { total: 0, imported: 0, failed: 0, errors: [] } })
    await importTipsFromText('content,topic\nhola,go')
    const call = mockedPost.mock.calls[0]
    expect(call).toBeDefined()
    expect(call?.[2]).toEqual({ headers: { 'Content-Type': 'multipart/form-data' } })
    expect((call?.[1] as FormData).get('content')).toBe('content,topic\nhola,go')
  })
})
