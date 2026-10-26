import { render, screen } from '@testing-library/react-native'

import { TopicRing } from './TopicRing'
import { colors } from '@/core/theme/tokens'

function colorOf(text: string): unknown {
  const style = screen.getByText(text).props.style as
    | Record<string, unknown>
    | Array<Record<string, unknown> | undefined>
  const flat = Array.isArray(style) ? Object.assign({}, ...style) : style
  return flat.color
}

describe('TopicRing', () => {
  it('muestra el porcentaje y lo anuncia como barra de progreso', async () => {
    await render(<TopicRing name="Go" accuracy={0.82} attempts={10} correctCount={8} />)

    expect(screen.getByText('82%')).toBeOnTheScreen()
    expect(screen.getByText('Go')).toBeOnTheScreen()

    const ring = screen.getByLabelText('Go: 82% de acierto, 8 de 10 aciertos')
    expect(ring).toHaveAccessibilityValue({ min: 0, max: 100, now: 82 })
  })

  it('marca el tema débil con el color y la etiqueta, sin icono', async () => {
    await render(<TopicRing name="SQL" accuracy={0.4} isWeak />)

    expect(colorOf('40%')).toBe(colors.accent)
    expect(screen.queryByText('Débil')).toBeNull()
    expect(screen.getByLabelText('SQL: 40% de acierto, tema débil')).toBeOnTheScreen()
  })

  it('un tema sano usa el color de texto principal', async () => {
    await render(<TopicRing name="Vue" accuracy={0.67} />)

    expect(colorOf('67%')).toBe(colors.foreground)
    expect(screen.getByLabelText('Vue: 67% de acierto')).toBeOnTheScreen()
  })
})
