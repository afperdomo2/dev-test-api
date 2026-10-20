import { apiClient } from '@/core/api/client'
import {
  createDailySession,
  finishSession,
  getNextQuestion,
  getSessionById,
  listSessions,
  submitAnswer,
} from './sessions.service'

jest.mock('@/core/api/client', () => ({
  apiClient: { get: jest.fn(), post: jest.fn(), put: jest.fn(), delete: jest.fn() },
}))

const mocked = apiClient as unknown as {
  get: jest.Mock
  post: jest.Mock
  put: jest.Mock
}

describe('sessions.service', () => {
  beforeEach(() => jest.clearAllMocks())

  it('listSessions envía paginación y estado', async () => {
    mocked.get.mockResolvedValue({ data: [], meta: { total: 0 } })

    await listSessions(2, 20, 'completed')

    expect(mocked.get).toHaveBeenCalledWith('/api/v1/sessions', {
      page: 2,
      perPage: 20,
      status: 'completed',
    })
  })

  it('getSessionById desempaqueta la sesión del detalle', async () => {
    const session = { id: 's1' }
    mocked.get.mockResolvedValue({ session, summary: null })

    await expect(getSessionById('s1')).resolves.toEqual(session)
    expect(mocked.get).toHaveBeenCalledWith('/api/v1/sessions/s1')
  })

  it('createDailySession hace POST sin body', async () => {
    mocked.post.mockResolvedValue({ session: { id: 's1' }, dueCount: 3 })

    await createDailySession()

    expect(mocked.post).toHaveBeenCalledWith('/api/v1/sessions/daily')
  })

  it('submitAnswer envía la respuesta', async () => {
    mocked.post.mockResolvedValue({ id: 'a1' })
    const body = { questionId: 'q1', selectedOptions: ['o1'] }

    await submitAnswer('s1', body)

    expect(mocked.post).toHaveBeenCalledWith('/api/v1/sessions/s1/answer', body)
  })

  it('getNextQuestion desempaqueta la pregunta', async () => {
    const question = { question: { id: 'q1' }, answeredCount: 0 }
    mocked.get.mockResolvedValue({ question })

    await expect(getNextQuestion('s1')).resolves.toEqual(question)
    expect(mocked.get).toHaveBeenCalledWith('/api/v1/sessions/s1/next')
  })

  it('finishSession hace PUT', async () => {
    mocked.put.mockResolvedValue({ id: 's1', status: 'completed' })

    await finishSession('s1')

    expect(mocked.put).toHaveBeenCalledWith('/api/v1/sessions/s1/finish')
  })
})
