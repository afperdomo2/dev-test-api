import { render, screen } from '@testing-library/react-native'

import { TopicMasteryRow } from './TopicMasteryRow'

const weakTopic = {
  topicName: 'Go',
  attempts: 5,
  accuracy: 0.4,
  isWeak: true,
  correctCount: 2,
}

describe('TopicMasteryRow', () => {
  it('muestra nombre, conteo y barra accesible', async () => {
    await render(<TopicMasteryRow {...weakTopic} />)

    expect(screen.getByText('Go')).toBeOnTheScreen()
    expect(screen.getByText('2/5 aciertos')).toBeOnTheScreen()

    const bar = screen.getByLabelText('Go: 40% de acierto')
    expect(bar).toHaveAccessibilityValue({ min: 0, max: 100, now: 40 })
  })

  it('muestra el chip "Débil" solo cuando isWeak', async () => {
    const weak = await render(<TopicMasteryRow {...weakTopic} />)
    expect(weak.getByText('Débil')).toBeOnTheScreen()

    const strong = await render(
      <TopicMasteryRow topicName="Vue" attempts={5} accuracy={0.8} correctCount={4} />,
    )
    expect(strong.queryByText('Débil')).toBeNull()
  })

  it('sin correctCount no pinta la línea de aciertos (uso en informes)', async () => {
    await render(<TopicMasteryRow topicName="Go" attempts={3} accuracy={0.6667} />)

    expect(screen.queryByText(/aciertos/)).toBeNull()
    expect(screen.getByLabelText('Go: 67% de acierto')).toBeOnTheScreen()
  })
})
