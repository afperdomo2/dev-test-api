import { fireEvent, render, screen } from '@testing-library/react-native'
import type {
  BadgesResponse,
  InsightSummary,
  OverviewResponse,
  TopicMasteryItem,
} from '@devtest/shared'

import ProgressDashboardScreen from './ProgressDashboardScreen'
import { getBadges, getOverview } from '@/core/api/services/gamification.service'
import { getInsightSummary } from '@/core/api/services/insights.service'
import { getTopicMastery } from '@/core/api/services/progress.service'
import { resetRouterMock, routerMock } from '@/__tests__/router'

jest.mock('@/core/api/services/gamification.service', () => ({
  getOverview: jest.fn(),
  getBadges: jest.fn(),
}))

jest.mock('@/core/api/services/insights.service', () => ({
  getInsightSummary: jest.fn(),
}))

jest.mock('@/core/api/services/progress.service', () => ({
  getTopicMastery: jest.fn(),
}))

const overview = {
  totalXp: 120,
  level: 2,
  xpToLevel: 80,
  streak: { current: 3, longest: 5, freezesAvailable: 1 },
  dailyGoal: { current: 10, target: 50, progress: 20, reached: false },
  topicXp: [],
} as unknown as OverviewResponse

const topics: Array<TopicMasteryItem> = [
  {
    topicId: 't1',
    topicName: 'Go',
    attempts: 5,
    correctCount: 2,
    accuracy: 0.4,
    isWeak: true,
  },
  {
    topicId: 't2',
    topicName: 'Vue',
    attempts: 5,
    correctCount: 4,
    accuracy: 0.8,
    isWeak: false,
  },
]

const summary = {
  period: 'week',
  from: '2026-10-13T00:00:00Z',
  to: '2026-10-19T12:00:00Z',
  sessions: 2,
  questionsAnswered: 5,
  accuracy: 0.6,
  timeSpentMs: 12 * 60 * 1000,
  topTopics: [],
  weakTopics: [],
  xpGained: 75,
  streakDays: 3,
  dailyActivity: [
    { date: '2026-10-18', questionsAnswered: 0 },
    { date: '2026-10-19', questionsAnswered: 5 },
  ],
} as unknown as InsightSummary

const badges = {
  badges: [
    {
      id: 'first_session',
      name: 'Primera sesión',
      description: 'Completa tu primera sesión',
      icon: 'mdi-flag-checkered',
      earned: true,
      current: 1,
      target: 1,
    },
    {
      id: 'topics_mastered_5',
      name: '5 temas dominados',
      description: 'Domina 5 temas',
      icon: 'mdi-school',
      earned: false,
      current: 2,
      target: 5,
    },
  ],
  earned: 1,
  total: 2,
} as unknown as BadgesResponse

describe('ProgressDashboardScreen', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    resetRouterMock()
    jest.mocked(getOverview).mockResolvedValue(overview)
    jest.mocked(getTopicMastery).mockResolvedValue(topics)
    jest.mocked(getInsightSummary).mockResolvedValue(summary)
    jest.mocked(getBadges).mockResolvedValue(badges)
  })

  it('muestra nivel, temas, actividad y logros en una sola vista', async () => {
    await render(<ProgressDashboardScreen />)

    expect(await screen.findByText('Go')).toBeOnTheScreen()
    expect(screen.getByText('Vue')).toBeOnTheScreen()
    expect(screen.getByText('Nv. 2')).toBeOnTheScreen()
    expect(screen.getByText('Temas')).toBeOnTheScreen()
    expect(screen.getByText('Actividad')).toBeOnTheScreen()
    expect(screen.getByText('Logros')).toBeOnTheScreen()
    expect(
      screen.getByText('5 preguntas · 60% de acierto · 12 min'),
    ).toBeOnTheScreen()
    expect(screen.getByText(/1 de 2 · Siguiente: 5 temas dominados/)).toBeOnTheScreen()
  })

  it('los accesos "Ver todos" navegan a cada detalle', async () => {
    await render(<ProgressDashboardScreen />)
    await screen.findByText('Go')

    await fireEvent.press(screen.getByLabelText('Ver todos temas'))
    await fireEvent.press(screen.getByLabelText('Ver detalle actividad'))
    await fireEvent.press(screen.getByLabelText('Ver todos logros'))

    expect(routerMock.push).toHaveBeenCalledWith('/(tabs)/progreso/temas')
    expect(routerMock.push).toHaveBeenCalledWith('/(tabs)/progreso/actividad')
    expect(routerMock.push).toHaveBeenCalledWith('/(tabs)/progreso/logros')
  })

  it('si fallan los temas el resto del dashboard sigue visible', async () => {
    jest
      .mocked(getTopicMastery)
      .mockRejectedValueOnce({ title: 'Network Error', status: 500, detail: 'Sin conexión' })

    await render(<ProgressDashboardScreen />)

    expect(await screen.findByText('Sin conexión')).toBeOnTheScreen()
    expect(screen.getByText(/5 preguntas · 60% de acierto/)).toBeOnTheScreen()
    expect(screen.getByLabelText('Reintentar temas')).toBeOnTheScreen()
  })

  it('sin temas responde con un mensaje, no con un hueco', async () => {
    jest.mocked(getTopicMastery).mockResolvedValue([])

    await render(<ProgressDashboardScreen />)

    expect(
      await screen.findByText('Responde preguntas para ver tu dominio por tema.'),
    ).toBeOnTheScreen()
  })

  it('sin actividad en la semana lo dice en la sección', async () => {
    jest.mocked(getInsightSummary).mockResolvedValue({
      ...summary,
      questionsAnswered: 0,
      sessions: 0,
      dailyActivity: [],
    })

    await render(<ProgressDashboardScreen />)

    expect(await screen.findByText('Aún no hay actividad esta semana.')).toBeOnTheScreen()
  })
})
