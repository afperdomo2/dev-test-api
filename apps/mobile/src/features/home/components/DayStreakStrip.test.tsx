import { render, screen } from '@testing-library/react-native'
import type { DailyActivityItem } from '@devtest/shared'

import { DayStreakStrip } from './DayStreakStrip'
import { colors } from '@/core/theme/tokens'

function styleOf(testID: string): Record<string, unknown> {
  // La tira oculta sus celdas a accesibilidad: el contenedor pulsable pone la etiqueta.
  const element = screen.getByTestId(testID, { includeHiddenElements: true })
  const style = element.props.style as
    | Record<string, unknown>
    | Array<Record<string, unknown> | undefined>
  return Array.isArray(style) ? Object.assign({}, ...style) : style
}

const days: Array<DailyActivityItem> = [
  { date: '2026-10-19', questionsAnswered: 4 },
  { date: '2026-10-20', questionsAnswered: 0 },
]

describe('DayStreakStrip', () => {
  it('rellena solo los días con actividad', async () => {
    await render(<DayStreakStrip days={days} />)

    expect(styleOf('day-cell-2026-10-19').backgroundColor).toBe(colors.primary)
    expect(styleOf('day-cell-2026-10-20').backgroundColor).toBe(colors.muted)
  })

  it('sin días no pinta nada', async () => {
    await render(<DayStreakStrip days={[]} />)

    expect(screen.queryByTestId('day-cell-2026-10-19')).toBeNull()
  })
})
