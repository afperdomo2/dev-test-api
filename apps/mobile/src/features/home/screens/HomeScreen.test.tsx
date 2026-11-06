import { Alert } from 'react-native'
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react-native'
import type { BadgesResponse, OverviewResponse, Tip } from '@devtest/shared'

import HomeScreen from './HomeScreen'
import { useAuthStore } from '@/core/auth/auth.store'
import { freezeStreak, getBadges, getOverview } from '@/core/api/services/gamification.service'
import { getInsightSummary } from '@/core/api/services/insights.service'
import { createDailySession } from '@/core/api/services/sessions.service'
import { clearPendingTip, getPendingTip } from '@/features/sessions/dailyTipStore'
import { resetRouterMock, routerMock } from '@/__tests__/router'

jest.mock('@/core/api/services/gamification.service', () => ({
  getOverview: jest.fn(),
  getBadges: jest.fn(),
  freezeStreak: jest.fn(),
}))

jest.mock('@/core/api/services/insights.service', () => ({
  getInsightSummary: jest.fn(),
}))

jest.mock('@/core/api/services/sessions.service', () => ({
  createDailySession: jest.fn(),
}))

const overview = {
  totalXp: 120,
  level: 2,
  xpToLevel: 80,
  streak: { current: 3, longest: 5, freezesAvailable: 1 },
  dailyGoal: { current: 10, target: 50, progress: 20, reached: false },
  topicXp: [{ topicId: 't1', name: 'Go', xp: 30 }],
} as unknown as OverviewResponse

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
      id: 'streak_7',
      name: 'Racha de 7',
      description: '7 días seguidos',
      icon: 'mdi-fire',
      earned: false,
      current: 3,
      target: 7,
    },
  ],
  earned: 1,
  total: 2,
} as unknown as BadgesResponse

const week = {
  period: 'week',
  from: '2026-10-13T00:00:00Z',
  to: '2026-10-19T12:00:00Z',
  sessions: 1,
  questionsAnswered: 4,
  accuracy: 0.75,
  timeSpentMs: 60000,
  topTopics: [],
  weakTopics: [],
  xpGained: 10,
  streakDays: 3,
  dailyActivity: [
    { date: '2026-10-18', questionsAnswered: 2 },
    { date: '2026-10-19', questionsAnswered: 2 },
  ],
}

describe('HomeScreen', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    resetRouterMock()
    clearPendingTip()
    useAuthStore.setState({
      user: { id: 'u1', email: 'ana@test.com', isAdmin: false } as never,
      token: 't',
    })
    jest.mocked(getOverview).mockResolvedValue(overview)
    jest.mocked(getBadges).mockResolvedValue(badges)
    jest.mocked(getInsightSummary).mockResolvedValue(week as never)
  })

  it('resume el día: nivel, repaso, racha, meta y logros', async () => {
    await render(<HomeScreen />)

    expect(await screen.findByText('Hola, ana')).toBeOnTheScreen()
    expect(await screen.findByLabelText('Nivel 2. 80 XP para el nivel 3')).toBeOnTheScreen()
    expect(screen.getByLabelText('Repasar hoy')).toBeOnTheScreen()
    expect(screen.getByText('Tu día')).toBeOnTheScreen()
    expect(screen.getByText('Meta diaria')).toBeOnTheScreen()
    expect(await screen.findByText('Logros')).toBeOnTheScreen()
    expect(screen.getByText('1 de 2')).toBeOnTheScreen()
    expect(screen.queryByText('Temas')).toBeNull()
  })

  it('el widget de logros enlaza con el detalle', async () => {
    await render(<HomeScreen />)
    await screen.findByText('Logros')

    await fireEvent.press(screen.getByLabelText('Ver todos logros'))

    expect(routerMock.push).toHaveBeenCalledWith('/(tabs)/progreso/logros')
  })

  it('el CTA crea la sesión diaria y navega a ella', async () => {
    jest.mocked(createDailySession).mockResolvedValue({ session: { id: 's1' } } as never)

    await render(<HomeScreen />)
    await screen.findByText('Hola, ana')

    await fireEvent.press(screen.getByLabelText('Repasar hoy'))

    await waitFor(() => expect(createDailySession).toHaveBeenCalledTimes(1))
    expect(routerMock.push).toHaveBeenCalledWith('/(tabs)/sesiones/s1')
  })

  it('guarda el tip del día para la sesión creada', async () => {
    const tip: Tip = {
      id: 'tip1',
      userId: 'u1',
      topicId: 't1',
      topic: 'Go',
      content: 'Usa defer',
      source: 'manual',
      createdAt: '2026-10-19T00:00:00Z',
      updatedAt: '2026-10-19T00:00:00Z',
    }
    jest.mocked(createDailySession).mockResolvedValue({ session: { id: 's1' }, tip } as never)

    await render(<HomeScreen />)
    await screen.findByText('Hola, ana')

    await fireEvent.press(screen.getByLabelText('Repasar hoy'))

    await waitFor(() => expect(getPendingTip('s1')).toEqual(tip))
  })

  it('si falla la sesión diaria muestra el detalle', async () => {
    jest.mocked(createDailySession).mockRejectedValue({
      title: 'Network Error',
      status: 500,
      detail: 'No se pudo conectar con el servidor',
    })

    await render(<HomeScreen />)
    await screen.findByText('Hola, ana')

    await fireEvent.press(screen.getByLabelText('Repasar hoy'))

    expect(await screen.findByText('No se pudo conectar con el servidor')).toBeOnTheScreen()
    expect(routerMock.push).not.toHaveBeenCalled()
  })

  it('la semana del "Tu día" lleva a la actividad', async () => {
    await render(<HomeScreen />)
    await screen.findByText('Hola, ana')

    await fireEvent.press(screen.getByLabelText(/Ver tu actividad/))

    expect(routerMock.push).toHaveBeenCalledWith('/(tabs)/progreso/actividad')
  })

  it('el freeze pide confirmación y solo entonces congela', async () => {
    const alertSpy = jest.spyOn(Alert, 'alert')
    jest.mocked(freezeStreak).mockResolvedValue(overview as never)

    await render(<HomeScreen />)
    await screen.findByText('Hola, ana')

    await fireEvent.press(screen.getByLabelText('Usar freeze para conservar racha'))

    expect(alertSpy).toHaveBeenCalledTimes(1)
    const buttons = alertSpy.mock.calls[0]?.[2] as
      | Array<{ text?: string; onPress?: () => Promise<void> }>
      | undefined
    const confirm = buttons?.find((button) => button.text === 'Usar')
    expect(confirm).toBeDefined()
    await act(async () => {
      await confirm?.onPress?.()
    })

    expect(freezeStreak).toHaveBeenCalledTimes(1)
    alertSpy.mockRestore()
  })

  it('si falla el overview ofrece reintentar sin bloquear el repaso', async () => {
    jest
      .mocked(getOverview)
      .mockRejectedValueOnce({ title: 'Network Error', status: 500, detail: 'Sin conexión' })
      .mockResolvedValueOnce(overview)

    await render(<HomeScreen />)

    expect(await screen.findByText('No se pudo cargar tu progreso')).toBeOnTheScreen()
    expect(screen.getByLabelText('Repasar hoy')).toBeOnTheScreen()

    await fireEvent.press(screen.getByLabelText('Reintentar carga de progreso'))

    expect(await screen.findByText('Meta diaria')).toBeOnTheScreen()
  })
})
