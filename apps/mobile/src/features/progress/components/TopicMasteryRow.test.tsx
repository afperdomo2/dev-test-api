import { render, screen } from '@testing-library/react-native'
import type { TopicMasteryItem } from '@devtest/shared'

import { TopicMasteryRow } from './TopicMasteryRow'

const weakTopic: TopicMasteryItem = {
  topicId: 't1',
  topicName: 'Go',
  attempts: 5,
  correctCount: 2,
  accuracy: 0.4,
  isWeak: true,
}

const strongTopic: TopicMasteryItem = {
  topicId: 't2',
  topicName: 'Vue',
  attempts: 5,
  correctCount: 4,
  accuracy: 0.8,
  isWeak: false,
}

describe('TopicMasteryRow', () => {
  it('muestra nombre, conteo y barra accesible', async () => {
    await render(<TopicMasteryRow item={weakTopic} />)

    expect(screen.getByText('Go')).toBeOnTheScreen()
    expect(screen.getByText('2/5 aciertos')).toBeOnTheScreen()

    const bar = screen.getByLabelText('Go: 40% de acierto')
    expect(bar).toHaveAccessibilityValue({ min: 0, max: 100, now: 40 })
  })

  it('muestra el chip "Débil" solo en temas débiles', async () => {
    const weak = await render(<TopicMasteryRow item={weakTopic} />)
    expect(weak.getByText('Débil')).toBeOnTheScreen()

    const strong = await render(<TopicMasteryRow item={strongTopic} />)
    expect(strong.queryByText('Débil')).toBeNull()
  })
})
