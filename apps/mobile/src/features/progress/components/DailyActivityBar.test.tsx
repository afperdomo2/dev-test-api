import { render, screen } from '@testing-library/react-native'

import { DailyActivityBar } from './DailyActivityBar'

function styleOf(testID: string): Record<string, unknown> {
  const element = screen.getByTestId(testID)
  const style = element.props.style as
    | Record<string, unknown>
    | Array<Record<string, unknown> | undefined>
  return Array.isArray(style) ? Object.assign({}, ...style) : style
}

describe('DailyActivityBar', () => {
  it('etiqueta el día con su número de preguntas', async () => {
    await render(<DailyActivityBar date="2026-10-19" value={5} max={10} />)

    expect(screen.getByLabelText('2026-10-19: 5 preguntas')).toBeOnTheScreen()
  })

  it('la altura es proporcional al máximo', async () => {
    await render(<DailyActivityBar date="2026-10-19" value={5} max={10} />)

    // 5/10 * 80 = 40
    expect(styleOf('daily-bar-2026-10-19').height).toBe(40)
  })

  it('un día sin actividad deja la barra al mínimo', async () => {
    await render(<DailyActivityBar date="2026-10-18" value={0} max={10} />)

    expect(styleOf('daily-bar-2026-10-18').height).toBe(2)
  })
})
