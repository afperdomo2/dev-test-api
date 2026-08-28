import {
  listSessions,
  getSessionById,
  getSessionDetail,
  createSession,
  finishSession,
  getNextQuestion,
  submitAnswer,
  getSessionSummary,
  deleteSession,
} from './sessions.service'
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

describe('sessions.service', () => {
  beforeEach(() => vi.clearAllMocks())

  it('listSessions without status', async () => {
    mockedGet.mockResolvedValue({ data: { data: [] } })
    await listSessions(1, 10)
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/sessions', { params: { page: 1, perPage: 10 } })
  })

  it('listSessions with status', async () => {
    mockedGet.mockResolvedValue({ data: {} })
    await listSessions(1, 10, 'active')
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/sessions', { params: { page: 1, perPage: 10, status: 'active' } })
  })

  it('getSessionById extracts session', async () => {
    const session = { id: '1' }
    mockedGet.mockResolvedValue({ data: { session } })
    const res = await getSessionById('1')
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/sessions/1')
    expect(res).toEqual(session)
  })

  it('getSessionDetail returns detail', async () => {
    const detail = { session: { id: '1' } }
    mockedGet.mockResolvedValue({ data: detail })
    const res = await getSessionDetail('1')
    expect(res).toEqual(detail)
  })

  it('createSession', async () => {
    mockedPost.mockResolvedValue({ data: { id: '1' } })
    const res = await createSession({ topicIds: [] } as never)
    expect(mockedPost).toHaveBeenCalledWith('/api/v1/sessions', { topicIds: [] })
    expect(res).toEqual({ id: '1' })
  })

  it('finishSession', async () => {
    mockedPut.mockResolvedValue({ data: { id: '1' } })
    const res = await finishSession('1')
    expect(mockedPut).toHaveBeenCalledWith('/api/v1/sessions/1/finish')
    expect(res).toEqual({ id: '1' })
  })

  it('getNextQuestion extracts question', async () => {
    mockedGet.mockResolvedValue({ data: { question: { id: 'q1' } } })
    const res = await getNextQuestion('1')
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/sessions/1/next')
    expect(res).toEqual({ id: 'q1' })
  })

  it('submitAnswer', async () => {
    mockedPost.mockResolvedValue({ data: { id: 'a1' } })
    const res = await submitAnswer('1', { questionId: 'q1', answer: 'x' } as never)
    expect(mockedPost).toHaveBeenCalledWith('/api/v1/sessions/1/answer', { questionId: 'q1', answer: 'x' })
    expect(res).toEqual({ id: 'a1' })
  })

  it('getSessionSummary', async () => {
    mockedGet.mockResolvedValue({ data: { score: 90 } })
    const res = await getSessionSummary('1')
    expect(mockedGet).toHaveBeenCalledWith('/api/v1/sessions/1/summary')
    expect(res).toEqual({ score: 90 })
  })

  it('deleteSession', async () => {
    mockedDelete.mockResolvedValue({})
    await deleteSession('1')
    expect(mockedDelete).toHaveBeenCalledWith('/api/v1/sessions/1')
  })
})
