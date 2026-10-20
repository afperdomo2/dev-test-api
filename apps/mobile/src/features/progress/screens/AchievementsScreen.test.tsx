import { fireEvent, render, screen } from '@testing-library/react-native'
import type { Badge, BadgesResponse } from '@devtest/shared'

import AchievementsScreen from './AchievementsScreen'
import { getBadges } from '@/core/api/services/gamification.service'

jest.mock('@/core/api/services/gamification.service', () => ({
  getBadges: jest.fn(),
}))

const badges: Badge[] = [
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
    id: 'streak_3',
    name: 'Racha de 3',
    description: '3 días seguidos de estudio',
    icon: 'mdi-fire',
    earned: false,
    current: 1,
    target: 3,
  },
  {
    id: 'streak_7',
    name: 'Racha de 7',
    description: '7 días seguidos de estudio',
    icon: 'mdi-fire',
    earned: false,
    current: 1,
    target: 7,
  },
  {
    id: 'streak_30',
    name: 'Racha de 30',
    description: '30 días seguidos de estudio',
    icon: 'mdi-fire',
    earned: false,
    current: 1,
    target: 30,
  },
  {
    id: 'perfect_session',
    name: 'Sesión perfecta',
    description: 'Termina una sesión con el 100% de acierto',
    icon: 'mdi-star',
    earned: false,
    current: 0,
    target: 1,
  },
  {
    id: 'answers_100',
    name: '100 respuestas',
    description: 'Responde 100 preguntas',
    icon: 'mdi-check-all',
    earned: false,
    current: 12,
    target: 100,
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
  {
    id: 'fast_solver',
    name: 'Resolvedor veloz',
    description: 'Sesión perfecta con menos de 10 s por pregunta',
    icon: 'mdi-lightning-bolt',
    earned: false,
    current: 0,
    target: 1,
  },
]

const response: BadgesResponse = { badges, earned: 1, total: 8 }

describe('AchievementsScreen', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    jest.mocked(getBadges).mockResolvedValue(response)
  })

  it('pinta los 8 logros y el contador', async () => {
    await render(<AchievementsScreen />)

    expect(await screen.findByText('1 de 8 logros')).toBeOnTheScreen()
    expect(screen.getByText('13%')).toBeOnTheScreen()

    for (const badge of badges) {
      expect(screen.getByText(badge.name)).toBeOnTheScreen()
    }

    expect(screen.getByText('Conseguido')).toBeOnTheScreen()
    expect(screen.getByText('12/100')).toBeOnTheScreen()
    expect(screen.getByLabelText('100 respuestas. 12 de 100')).toBeOnTheScreen()
  })

  it('sin logros muestra el estado vacío', async () => {
    jest.mocked(getBadges).mockResolvedValue({ badges: [], earned: 0, total: 0 })

    await render(<AchievementsScreen />)

    expect(await screen.findByText('Aún no hay logros disponibles')).toBeOnTheScreen()
    expect(screen.queryByText('0 de 0 logros')).toBeNull()
  })

  it('si falla la carga permite reintentar', async () => {
    jest
      .mocked(getBadges)
      .mockRejectedValueOnce({ title: 'Network Error', status: 500, detail: 'Sin conexión' })
      .mockResolvedValueOnce(response)

    await render(<AchievementsScreen />)

    expect(await screen.findByText('Sin conexión')).toBeOnTheScreen()

    await fireEvent.press(screen.getByLabelText('Reintentar'))

    expect(await screen.findByText('1 de 8 logros')).toBeOnTheScreen()
    expect(getBadges).toHaveBeenCalledTimes(2)
  })
})
