import { fireEvent, render, screen } from '@testing-library/react-native'
import type { TopicMasteryItem } from '@devtest/shared'

import TopicMasteryScreen from './TopicMasteryScreen'
import { getTopicMastery } from '@/core/api/services/progress.service'

jest.mock('@/core/api/services/progress.service', () => ({
  getTopicMastery: jest.fn(),
}))

const items: Array<TopicMasteryItem> = [
  {
    topicId: 't1',
    topicName: 'Go',
    attempts: 5,
    correctCount: 2,
    accuracy: 0.4,
    isWeak: true,
  },
  {
    topicId: 't2',
    topicName: 'Vue',
    attempts: 5,
    correctCount: 4,
    accuracy: 0.8,
    isWeak: false,
  },
]

describe('TopicMasteryScreen', () => {
  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('muestra una fila por tema', async () => {
    jest.mocked(getTopicMastery).mockResolvedValue(items)

    await render(<TopicMasteryScreen />)

    expect(await screen.findByText('Go')).toBeOnTheScreen()
    expect(screen.getByText('Vue')).toBeOnTheScreen()
    expect(screen.getByText('2/5 aciertos')).toBeOnTheScreen()
    expect(screen.getByText('Débil')).toBeOnTheScreen()
  })

  it('sin datos muestra el estado vacío', async () => {
    jest.mocked(getTopicMastery).mockResolvedValue([])

    await render(<TopicMasteryScreen />)

    expect(await screen.findByText('Aún no hay respuestas registradas')).toBeOnTheScreen()
  })

  it('si falla la carga muestra el error y permite reintentar', async () => {
    jest
      .mocked(getTopicMastery)
      .mockRejectedValueOnce({ title: 'Network Error', status: 500, detail: 'Sin conexión' })
      .mockResolvedValueOnce(items)

    await render(<TopicMasteryScreen />)

    expect(await screen.findByText('Sin conexión')).toBeOnTheScreen()

    await fireEvent.press(screen.getByLabelText('Reintentar'))

    expect(await screen.findByText('Go')).toBeOnTheScreen()
    expect(getTopicMastery).toHaveBeenCalledTimes(2)
  })
})
