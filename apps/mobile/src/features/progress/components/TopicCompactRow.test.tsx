import { render, screen } from '@testing-library/react-native'

import { TopicCompactRow } from './TopicCompactRow'

describe('TopicCompactRow', () => {
  it('muestra nombre, respuestas y porcentaje', async () => {
    await render(<TopicCompactRow topicName="React" attempts={12} accuracy={0.82} />)

    expect(screen.getByText('React')).toBeOnTheScreen()
    expect(screen.getByText('12 resp · 82%')).toBeOnTheScreen()
    expect(screen.getByLabelText('React: 82% de acierto en 12 respuestas')).toBeOnTheScreen()
  })

  it('añade el chip "Débil" y los aciertos cuando corresponde', async () => {
    await render(
      <TopicCompactRow topicName="SQL" attempts={5} accuracy={0.4} isWeak correctCount={2} />,
    )

    expect(screen.getByText('Débil')).toBeOnTheScreen()
    expect(
      screen.getByLabelText('SQL: 40% de acierto en 5 respuestas, tema débil, 2 de 5 aciertos'),
    ).toHaveAccessibilityValue({ min: 0, max: 100, now: 40 })
  })
})
