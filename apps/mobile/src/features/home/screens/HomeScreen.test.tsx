import { Alert } from 'react-native'
import { fireEvent, render, screen, waitFor } from '@testing-library/react-native'
import type { OverviewResponse } from '@devtest/shared'

import HomeScreen from './HomeScreen'
import { useAuthStore } from '@/core/auth/auth.store'
import { freezeStreak, getOverview } from '@/core/api/services/gamification.service'
import { createDailySession } from '@/core/api/services/sessions.service'
import { resetRouterMock, routerMock } from '@/__tests__/router'

jest.mock('@/core/api/services/gamification.service', () => ({
  getOverview: jest.fn(),
  freezeStreak: jest.fn(),
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

describe('HomeScreen', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    resetRouterMock()
    useAuthStore.setState({
      user: { id: 'u1', email: 'ana@test.com', isAdmin: false } as never,
      token: 't',
    })
    jest.mocked(getOverview).mockResolvedValue(overview)
  })

  it('muestra saludo, secciones y el CTA de repaso', async () => {
    await render(<HomeScreen />)

    expect(await screen.findByText('Hola, ana')).toBeOnTheScreen()
    expect(screen.getByLabelText('Repasar hoy')).toBeOnTheScreen()
    expect(screen.getByText('Racha')).toBeOnTheScreen()
    expect(screen.getByText('Meta diaria')).toBeOnTheScreen()
    expect(screen.getByText('Progreso')).toBeOnTheScreen()
  })

  it('el CTA crea la sesión diaria y navega a ella', async () => {
    jest.mocked(createDailySession).mockResolvedValue({ session: { id: 's1' } } as never)

    await render(<HomeScreen />)
    await screen.findByText('Hola, ana')

    await fireEvent.press(screen.getByLabelText('Repasar hoy'))

    await waitFor(() => expect(createDailySession).toHaveBeenCalledTimes(1))
    expect(routerMock.push).toHaveBeenCalledWith('/(tabs)/sesiones/s1')
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
    await confirm?.onPress?.()

    expect(freezeStreak).toHaveBeenCalledTimes(1)
    alertSpy.mockRestore()
  })
})
