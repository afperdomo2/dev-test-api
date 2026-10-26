import { fireEvent, render, screen } from '@testing-library/react-native'
import type { InsightSummary } from '@devtest/shared'

import ActivityScreen from './ActivityScreen'
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

describe('ActivityScreen', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    jest.mocked(getInsightSummary).mockResolvedValue(summary)
  })

  it('resume el periodo, los temas y la actividad', async () => {
    await render(<ActivityScreen />)

    expect(await screen.findByText('5')).toBeOnTheScreen()
    expect(screen.getByText('preguntas en 2 sesiones')).toBeOnTheScreen()
    expect(screen.getByText('60%')).toBeOnTheScreen()
    expect(screen.getByText('12 min')).toBeOnTheScreen()
    expect(screen.getByText('+75')).toBeOnTheScreen()

    expect(screen.getByText('Más practicados')).toBeOnTheScreen()
    expect(screen.getByText('Go')).toBeOnTheScreen()
    expect(screen.getByText('Puntos a reforzar')).toBeOnTheScreen()
    expect(screen.getByText('SQL')).toBeOnTheScreen()
    expect(screen.getByLabelText('2026-10-19: 5 preguntas')).toBeOnTheScreen()
  })

  it('al cambiar a Mes consulta el periodo y pinta el mapa de calor', async () => {
    await render(<ActivityScreen />)
    await screen.findByText('Más practicados')

    await fireEvent.press(screen.getByLabelText('Mes'))

    expect(getInsightSummary).toHaveBeenCalledWith('month')
    expect(await screen.findByTestId('heat-cell-2026-10-19')).toBeOnTheScreen()
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

    await render(<ActivityScreen />)

    expect(await screen.findByText('No hay actividad en este periodo.')).toBeOnTheScreen()
    expect(screen.queryByText('Más practicados')).toBeNull()
    expect(screen.queryByText('Puntos a reforzar')).toBeNull()
  })

  it('si falla la carga permite reintentar', async () => {
    jest
      .mocked(getInsightSummary)
      .mockRejectedValueOnce({ title: 'Network Error', status: 500, detail: 'Sin conexión' })
      .mockResolvedValueOnce(summary)

    await render(<ActivityScreen />)

    expect(await screen.findByText('Sin conexión')).toBeOnTheScreen()

    await fireEvent.press(screen.getByLabelText('Reintentar'))

    expect(await screen.findByText('Más practicados')).toBeOnTheScreen()
    expect(getInsightSummary).toHaveBeenCalledTimes(2)
  })
})
