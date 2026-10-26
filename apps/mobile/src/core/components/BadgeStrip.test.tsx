import { render, screen } from '@testing-library/react-native'
import type { Badge, BadgesResponse } from '@devtest/shared'

import { BadgeStrip } from './BadgeStrip'

const badges: Array<Badge> = [
  {
    id: 'streak_7',
    name: 'Racha de 7',
    description: '7 días seguidos',
    icon: 'mdi-fire',
    earned: false,
    current: 1,
    target: 7,
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
    description: '3 días seguidos',
    icon: 'mdi-fire',
    earned: false,
    current: 1,
    target: 3,
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
    id: 'perfect_session',
    name: 'Sesión perfecta',
    description: 'Sesión sin fallos',
    icon: 'mdi-star',
    earned: false,
    current: 0,
    target: 1,
  },
]

const response: BadgesResponse = { badges, earned: 1, total: 6 }

describe('BadgeStrip', () => {
  it('en modo compacto pinta iconos, contador y el chip de los que sobran', async () => {
    await render(<BadgeStrip data={response} max={4} />)

    expect(screen.getByText('1 de 6')).toBeOnTheScreen()
    expect(screen.getByLabelText('Primera sesión, conseguido')).toBeOnTheScreen()
    expect(screen.getByLabelText('5 temas dominados')).toBeOnTheScreen()
    expect(screen.getByLabelText('2 logros más')).toBeOnTheScreen()
    expect(screen.queryByLabelText('17% de logros')).toBeNull()
  })

  it('en modo detallado añade el siguiente logro y la barra', async () => {
    await render(<BadgeStrip data={response} detailed />)

    expect(
      screen.getByText('1 de 6 · Siguiente: 5 temas dominados (2/5)'),
    ).toBeOnTheScreen()

    const bar = screen.getByLabelText('17% de logros')
    expect(bar).toHaveAccessibilityValue({ min: 0, max: 100, now: 17 })
  })

  it('sin pendientes no propone siguiente logro', async () => {
    const allEarned: BadgesResponse = {
      badges: [{ ...badges[2]! }],
      earned: 1,
      total: 1,
    }

    await render(<BadgeStrip data={allEarned} detailed />)

    expect(screen.getByText('1 de 1')).toBeOnTheScreen()
    expect(screen.queryByText(/Siguiente/)).toBeNull()
  })

  it('sin catálogo no pinta nada', async () => {
    const empty = await render(<BadgeStrip data={{ badges: [], earned: 0, total: 0 }} />)

    expect(empty.toJSON()).toBeNull()
  })
})
