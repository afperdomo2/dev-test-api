import { fireEvent, render, screen } from '@testing-library/react-native'
import type { InsightSummary } from '@devtest/shared'

import WeeklySummaryScreen from './WeeklySummaryScreen'
import { getInsightSummary } from '@/core/api/services/insights.service'

jest.mock('@/core/api/services/insights.service', () => ({
  getInsightSummary: jest.fn(),
}))

const summary: InsightSummary = {
  period: 'week',
  from: '2026-10-13T00:00:00Z',
  to: '2026-10-19T12:00:00Z',
  sessions: 2,
  questionsAnswered: 5,
  accuracy: 0.6,
  timeSpentMs: 12 * 60 * 1000,
  topTopics: [{ topicId: 't1', topicName: 'Go', attempts: 3, accuracy: 0.6667 }],
  weakTopics: [
    {
      topicId: 't2',
      topicName: 'SQL',
      attempts: 5,
      correctCount: 2,
      accuracy: 0.4,
      isWeak: true,
    },
  ],
  xpGained: 75,
  streakDays: 3,
  dailyActivity: [
    { date: '2026-10-18', questionsAnswered: 0 },
    { date: '2026-10-19', questionsAnswered: 5 },
  ],
}

describe('WeeklySummaryScreen', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    jest.mocked(getInsightSummary).mockResolvedValue(summary)
  })

  it('muestra los KPIs, los temas y la actividad', async () => {
    await render(<WeeklySummaryScreen />)

    expect(await screen.findByText('5')).toBeOnTheScreen()
    expect(screen.getByText('Sesiones')).toBeOnTheScreen()
    expect(screen.getByText('60%')).toBeOnTheScreen()
    expect(screen.getByText('12 min')).toBeOnTheScreen()
    expect(screen.getByText('75')).toBeOnTheScreen()
    expect(screen.getByText('3 días')).toBeOnTheScreen()

    expect(screen.getByText('Más practicados')).toBeOnTheScreen()
    expect(screen.getByText('Go')).toBeOnTheScreen()
    expect(screen.getByText('Temas débiles')).toBeOnTheScreen()
    expect(screen.getByText('SQL')).toBeOnTheScreen()
    expect(screen.getByLabelText('2026-10-19: 5 preguntas')).toBeOnTheScreen()
  })

  it('al cambiar a Mes vuelve a consultar', async () => {
    await render(<WeeklySummaryScreen />)
    await screen.findByText('Sesiones')

    await fireEvent.press(screen.getByLabelText('Mes'))

    expect(getInsightSummary).toHaveBeenCalledWith('month')
  })

  it('sin actividad muestra el aviso', async () => {
    jest.mocked(getInsightSummary).mockResolvedValue({
      ...summary,
      sessions: 0,
      questionsAnswered: 0,
      topTopics: [],
      weakTopics: [],
      dailyActivity: [],
    })

    await render(<WeeklySummaryScreen />)

    expect(await screen.findByText('No hay actividad en este periodo.')).toBeOnTheScreen()
    expect(screen.getByText('Sin datos en este periodo.')).toBeOnTheScreen()
    expect(screen.queryByText('Temas débiles')).toBeNull()
  })

  it('si falla la carga permite reintentar', async () => {
    jest
      .mocked(getInsightSummary)
      .mockRejectedValueOnce({ title: 'Network Error', status: 500, detail: 'Sin conexión' })
      .mockResolvedValueOnce(summary)

    await render(<WeeklySummaryScreen />)

    expect(await screen.findByText('Sin conexión')).toBeOnTheScreen()

    await fireEvent.press(screen.getByLabelText('Reintentar'))

    expect(await screen.findByText('Sesiones')).toBeOnTheScreen()
    expect(getInsightSummary).toHaveBeenCalledTimes(2)
  })
})
