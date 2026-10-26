import { render, screen } from '@testing-library/react-native'
import type { Session } from '@devtest/shared'

import SessionsListPage from './SessionsListPage'
import { listSessions } from '@/core/api/services/sessions.service'

jest.mock('@/core/api/services/sessions.service', () => ({
  listSessions: jest.fn(),
  deleteSession: jest.fn(),
  finishSession: jest.fn(),
}))

jest.mock('@/features/sessions/components/SessionActionsSheet', () => ({
  SessionActionsSheet: () => null,
}))

function buildSession(overrides: Partial<Session> = {}): Session {
  return {
    id: 's1',
    userId: 'u1',
    name: 'Repaso de Go',
    status: 'in_progress',
    mode: 'review',
    difficulty: 'intermediate',
    questionLimit: 20,
    score: 0,
    startedAt: '2026-10-25T10:00:00.000Z',
    finishedAt: null,
    topics: ['Go', 'Concurrencia', 'Testing'],
    answerCount: 12,
    createdAt: '2026-10-25T10:00:00.000Z',
    updatedAt: '2026-10-25T10:00:00.000Z',
    ...overrides,
  }
}

function mockList(sessions: Array<Session>) {
  jest.mocked(listSessions).mockResolvedValue({
    data: sessions,
    meta: { total: sessions.length, page: 1, perPage: 20 },
  } as never)
}

describe('SessionsListPage', () => {
  beforeEach(() => jest.clearAllMocks())

  it('muestra el banner de continuar para la sesion en curso', async () => {
    mockList([buildSession()])

    await render(<SessionsListPage />)

    expect(await screen.findByText('Continuar: Repaso de Go')).toBeOnTheScreen()
    expect(screen.getByLabelText(/Continuar sesion Repaso de Go/)).toBeOnTheScreen()
  })

  it('muestra el progreso con el limite de preguntas', async () => {
    mockList([buildSession()])

    await render(<SessionsListPage />)

    expect(await screen.findByText('12/20 respuestas')).toBeOnTheScreen()
  })

  it('sin sesion en curso no muestra el banner', async () => {
    mockList([buildSession({ status: 'completed', score: 87.5, answerCount: 20 })])

    await render(<SessionsListPage />)

    expect(await screen.findByText('Repaso de Go')).toBeOnTheScreen()
    expect(screen.queryByText(/Continuar:/)).toBeNull()
    expect(screen.getByText('87.5 pts')).toBeOnTheScreen()
  })

  it('renderiza los cuatro filtros de estado', async () => {
    mockList([])

    await render(<SessionsListPage />)

    for (const filter of ['Todas', 'En progreso', 'Completada', 'Cancelada']) {
      expect(await screen.findByLabelText(`Filtrar por ${filter}`)).toBeOnTheScreen()
    }
  })
})
