import { act, fireEvent, render, screen } from '@testing-library/react-native'
import type { Topic } from '@devtest/shared'

import { TopicPicker } from './TopicPicker'
import { listTopics } from '@/core/api/services/topics.service'

jest.mock('@/core/api/services/topics.service', () => ({
  listTopics: jest.fn(),
}))

const java: Topic = {
  id: 't1',
  slug: 'java',
  name: 'Java',
  category: 'lenguajes',
  isSystem: true,
  createdAt: '2026-10-01T00:00:00.000Z',
}

const spring: Topic = {
  id: 't2',
  slug: 'spring',
  name: 'Spring',
  category: 'backend',
  isSystem: true,
  createdAt: '2026-10-01T00:00:00.000Z',
}

function mockTopics(topics: Topic[]) {
  jest.mocked(listTopics).mockResolvedValue({
    data: topics,
    meta: { total: topics.length, page: 1, perPage: 50 },
  } as never)
}

describe('TopicPicker', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    mockTopics([java, spring])
  })

  it('muestra la pista y permite abrir el selector', async () => {
    const onChange = jest.fn()
    await render(<TopicPicker selected={[]} onChange={onChange} />)

    expect(screen.getByText('Elige al menos un tema para empezar.')).toBeOnTheScreen()

    await fireEvent.press(screen.getByLabelText('Añadir temas'))

    expect(await screen.findByText('Elegir temas')).toBeOnTheScreen()
    expect(await screen.findByLabelText('Java, lenguajes')).toBeOnTheScreen()
  })

  it('selecciona un tema y solo lo confirma al pulsar Listo', async () => {
    const onChange = jest.fn()
    await render(<TopicPicker selected={[]} onChange={onChange} />)

    await fireEvent.press(screen.getByLabelText('Añadir temas'))
    await fireEvent.press(await screen.findByLabelText('Java, lenguajes'))

    expect(onChange).not.toHaveBeenCalled()

    await fireEvent.press(screen.getByLabelText('Confirmar temas'))

    expect(onChange).toHaveBeenCalledWith([java])
  })

  it('respeta el máximo y deshabilita el resto', async () => {
    const onChange = jest.fn()
    await render(<TopicPicker selected={[]} onChange={onChange} max={1} />)

    await fireEvent.press(screen.getByLabelText('Añadir temas'))
    await fireEvent.press(await screen.findByLabelText('Java, lenguajes'))

    expect(await screen.findByText('Máximo 1 temas. Quita uno para añadir otro.')).toBeOnTheScreen()
    expect(screen.getByLabelText('Spring, backend')).toBeDisabled()
  })

  it('permite quitar un tema ya seleccionado desde el formulario', async () => {
    const onChange = jest.fn()
    await render(<TopicPicker selected={[java]} onChange={onChange} />)

    await fireEvent.press(screen.getByLabelText('Quitar Java'))

    expect(onChange).toHaveBeenCalledWith([])
  })

  it('busca en el backend con debounce', async () => {
    const onChange = jest.fn()
    await render(<TopicPicker selected={[]} onChange={onChange} />)
    await fireEvent.press(screen.getByLabelText('Añadir temas'))
    await screen.findByLabelText('Java, lenguajes')

    await fireEvent.changeText(screen.getByLabelText('Buscar tema'), 'spr')
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 450))
    })

    expect(jest.mocked(listTopics)).toHaveBeenLastCalledWith({
      search: 'spr',
      hasQuestions: true,
      perPage: 50,
    })
  })
})
