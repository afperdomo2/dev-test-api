import { fireEvent, render, screen, waitFor } from '@testing-library/react-native'
import { useRouter } from 'expo-router'
import type { Topic } from '@devtest/shared'

import StudySessionScreen from './StudySessionScreen'
import { createSession } from '@/core/api/services/sessions.service'
import { listTopics } from '@/core/api/services/topics.service'

jest.mock('@/core/api/services/sessions.service', () => ({
  createSession: jest.fn(),
}))

jest.mock('@/core/api/services/topics.service', () => ({
  listTopics: jest.fn(),
}))

jest.mock('expo-router', () => ({
  useRouter: jest.fn(),
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

const replace = jest.fn()
const back = jest.fn()

beforeEach(() => {
  jest.clearAllMocks()
  jest.mocked(useRouter).mockReturnValue({ replace, push: jest.fn(), back } as never)
  jest.mocked(listTopics).mockResolvedValue({
    data: [java, spring],
    meta: { total: 2, page: 1, perPage: 50 },
  } as never)
  jest.mocked(createSession).mockResolvedValue({ id: 's1' } as never)
})

async function pickTopics(labels: string[]) {
  await fireEvent.press(screen.getByLabelText('Añadir temas'))
  for (const label of labels) {
    await fireEvent.press(await screen.findByLabelText(label))
  }
  await fireEvent.press(screen.getByLabelText('Confirmar temas'))
}

describe('StudySessionScreen', () => {
  it('no permite empezar sin temas', async () => {
    await render(<StudySessionScreen />)

    expect(screen.getByLabelText('Empezar sesión')).toBeDisabled()
    expect(createSession).not.toHaveBeenCalled()
  })

  it('crea una sesión review con el límite elegido', async () => {
    await render(<StudySessionScreen />)
    await pickTopics(['Java, lenguajes'])

    await fireEvent.press(screen.getByLabelText('10 preguntas'))
    await fireEvent.press(screen.getByLabelText('Empezar sesión'))

    await waitFor(() =>
      expect(createSession).toHaveBeenCalledWith({
        name: 'Estudio · Java',
        mode: 'review',
        difficulty: 'random',
        topicIds: ['t1'],
        questionLimit: 10,
      }),
    )
    expect(replace).toHaveBeenCalledWith('/(tabs)/sesiones/s1')
  })

  it('envía questionLimit undefined cuando se elige sin límite', async () => {
    await render(<StudySessionScreen />)
    await pickTopics(['Java, lenguajes', 'Spring, backend'])

    await fireEvent.press(screen.getByLabelText('Sin límite de preguntas'))
    await fireEvent.press(screen.getByLabelText('Empezar sesión'))

    await waitFor(() =>
      expect(createSession).toHaveBeenCalledWith(
        expect.objectContaining({
          name: 'Estudio · Java, Spring',
          topicIds: ['t1', 't2'],
          questionLimit: undefined,
        }),
      ),
    )
  })

  it('ofrece reintentar sin límite cuando el backend rechaza la cantidad', async () => {
    jest
      .mocked(createSession)
      .mockRejectedValueOnce({
        status: 422,
        detail: 'Solo hay 30 preguntas disponibles para estos temas, pero pediste 50.',
      })
      .mockResolvedValueOnce({ id: 's2' } as never)

    await render(<StudySessionScreen />)
    await pickTopics(['Java, lenguajes'])

    await fireEvent.press(screen.getByLabelText('Empezar sesión'))

    expect(
      await screen.findByText('Solo hay 30 preguntas disponibles para estos temas, pero pediste 50.'),
    ).toBeOnTheScreen()

    await fireEvent.press(screen.getByLabelText('Empezar sin límite de preguntas'))

    await waitFor(() => expect(createSession).toHaveBeenCalledTimes(2))
    expect(jest.mocked(createSession).mock.calls[1]?.[0]).toMatchObject({
      questionLimit: undefined,
    })
  })

  it('respeta el nombre editado al cambiar de temas', async () => {
    await render(<StudySessionScreen />)
    await pickTopics(['Java, lenguajes'])

    await fireEvent.changeText(screen.getByLabelText('Nombre de la sesión'), 'Mi repaso')
    await pickTopics(['Spring, backend'])

    expect(screen.getByLabelText('Nombre de la sesión').props.value).toBe('Mi repaso')

    await fireEvent.press(screen.getByLabelText('Empezar sesión'))

    await waitFor(() =>
      expect(createSession).toHaveBeenCalledWith(
        expect.objectContaining({ name: 'Mi repaso', topicIds: ['t1', 't2'] }),
      ),
    )
  })

  it('vuelve atrás desde el header', async () => {
    await render(<StudySessionScreen />)

    await fireEvent.press(screen.getByLabelText('Volver'))

    expect(back).toHaveBeenCalled()
  })
})
