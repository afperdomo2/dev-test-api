import { render, screen } from '@testing-library/react-native'
import type { DailyActivityItem } from '@devtest/shared'

import { ActivityHeatmap } from './ActivityHeatmap'

const days: Array<DailyActivityItem> = [
  { date: '2026-10-19', questionsAnswered: 8 },
  { date: '2026-10-20', questionsAnswered: 0 },
  { date: '2026-10-21', questionsAnswered: 3 },
]

describe('ActivityHeatmap', () => {
  it('pinta una celda por día con su etiqueta accesible', async () => {
    await render(<ActivityHeatmap days={days} />)

    expect(screen.getByLabelText(/19: 8 preguntas$/)).toBeOnTheScreen()
    expect(screen.getByLabelText(/20: 0 preguntas$/)).toBeOnTheScreen()
    expect(screen.getByTestId('heat-cell-2026-10-21')).toBeOnTheScreen()
  })

  it('sin días no pinta nada', async () => {
    await render(<ActivityHeatmap days={[]} />)

    expect(screen.queryByTestId('heat-cell-2026-10-19')).toBeNull()
  })
})
