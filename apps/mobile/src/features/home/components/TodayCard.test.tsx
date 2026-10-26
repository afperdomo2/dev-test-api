import { fireEvent, render, screen } from '@testing-library/react-native'
import type { ComponentProps } from 'react'
import type { DailyActivityItem, DailyGoalInfo, StreakInfo } from '@devtest/shared'

import { TodayCard } from './TodayCard'

const streak: StreakInfo = { current: 5, longest: 9, freezesAvailable: 2 }
const goal: DailyGoalInfo = { current: 30, target: 50, progress: 60, reached: false }
const week: Array<DailyActivityItem> = [
  { date: '2026-10-19', questionsAnswered: 4 },
  { date: '2026-10-20', questionsAnswered: 0 },
  { date: '2026-10-21', questionsAnswered: 2 },
  { date: '2026-10-22', questionsAnswered: 0 },
  { date: '2026-10-23', questionsAnswered: 1 },
  { date: '2026-10-24', questionsAnswered: 0 },
  { date: '2026-10-25', questionsAnswered: 6 },
]

function renderCard(overrides: Partial<ComponentProps<typeof TodayCard>> = {}) {
  return render(
    <TodayCard
      streak={streak}
      dailyGoal={goal}
      week={week}
      freezing={false}
      onFreeze={jest.fn()}
      onOpenActivity={jest.fn()}
      {...overrides}
    />,
  )
}

describe('TodayCard', () => {
  it('resume racha, semana y meta diaria', async () => {
    await renderCard()

    expect(screen.getByText('5')).toBeOnTheScreen()
    expect(screen.getByText('días')).toBeOnTheScreen()
    expect(screen.getByText('9')).toBeOnTheScreen()
    expect(screen.getByText('Meta diaria')).toBeOnTheScreen()
    expect(screen.getByText('30 / 50 XP')).toBeOnTheScreen()
  })

  it('la semana es pulsable y lleva a la actividad', async () => {
    const onOpenActivity = jest.fn()
    await renderCard({ onOpenActivity })

    await fireEvent.press(
      screen.getByLabelText('Ver tu actividad: practicaste 4 de los últimos 7 días'),
    )

    expect(onOpenActivity).toHaveBeenCalledTimes(1)
  })

  it('sin freezes disponibles el botón queda deshabilitado', async () => {
    const onFreeze = jest.fn()
    await renderCard({
      streak: { ...streak, freezesAvailable: 0 },
      onFreeze,
    })

    const button = screen.getByLabelText('Usar freeze para conservar racha')
    expect(button).toBeDisabled()

    await fireEvent.press(button)
    expect(onFreeze).not.toHaveBeenCalled()
  })

  it('sin datos de la semana oculta la tira de días', async () => {
    await renderCard({ week: null })

    expect(screen.queryByLabelText(/Ver tu actividad/)).toBeNull()
    expect(screen.getByText('Meta diaria')).toBeOnTheScreen()
  })
})
