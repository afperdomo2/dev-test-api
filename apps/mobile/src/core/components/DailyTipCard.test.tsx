import { fireEvent, render, screen } from '@testing-library/react-native'
import type { Tip } from '@devtest/shared'

import { DailyTipCard } from './DailyTipCard'

const tip: Tip = {
  id: 'tip1',
  userId: 'u1',
  topicId: 't1',
  topic: 'Go',
  content: 'Usa defer para cerrar recursos.',
  source: 'manual',
  createdAt: '2026-10-19T00:00:00Z',
  updatedAt: '2026-10-19T00:00:00Z',
}

describe('DailyTipCard', () => {
  it('muestra el título, el tema y el contenido', async () => {
    await render(<DailyTipCard tip={tip} onDismiss={jest.fn()} />)

    expect(screen.getByText('Tip del día')).toBeOnTheScreen()
    expect(screen.getByText('Go')).toBeOnTheScreen()
    expect(screen.getByText(/Usa defer para cerrar recursos/)).toBeOnTheScreen()
  })

  it('el botón de cerrar llama a onDismiss y tiene nombre accesible', async () => {
    const onDismiss = jest.fn()
    await render(<DailyTipCard tip={tip} onDismiss={onDismiss} />)

    await fireEvent.press(screen.getByLabelText('Cerrar tip del día'))

    expect(onDismiss).toHaveBeenCalledTimes(1)
  })
})
