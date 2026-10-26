import { render, screen } from '@testing-library/react-native'
import type { Badge } from '@devtest/shared'

import { toMcIcon } from '@/core/utils/icons'

import { BadgeCard } from './BadgeCard'

const earnedBadge: Badge = {
  id: 'first_session',
  name: 'Primera sesión',
  description: 'Completa tu primera sesión',
  icon: 'mdi-flag-checkered',
  earned: true,
  current: 1,
  target: 1,
}

const pendingBadge: Badge = {
  id: 'answers_100',
  name: '100 respuestas',
  description: 'Responde 100 preguntas',
  icon: 'mdi-check-all',
  earned: false,
  current: 12,
  target: 100,
}

describe('BadgeCard', () => {
  it('un logro conseguido se anuncia con icono y texto', async () => {
    await render(<BadgeCard badge={earnedBadge} />)

    expect(screen.getByLabelText('Primera sesión. Conseguido')).toBeOnTheScreen()
    expect(screen.getByText('Conseguido')).toBeOnTheScreen()
    expect(screen.getByText('flag-checkered')).toBeOnTheScreen()
  })

  it('un logro pendiente anuncia y pinta su progreso', async () => {
    await render(<BadgeCard badge={pendingBadge} />)

    expect(screen.getByLabelText('100 respuestas. 12 de 100')).toBeOnTheScreen()
    expect(screen.getByText('12/100')).toBeOnTheScreen()
    expect(screen.getByText('check-all')).toBeOnTheScreen()
  })

  it('mapea mdi-* y cae al trofeo si el icono no existe', () => {
    expect(toMcIcon('mdi-fire')).toBe('fire')
    expect(toMcIcon('fire')).toBe('fire')
    expect(toMcIcon('mdi-icono-inexistente')).toBe('trophy-outline')
  })
})
