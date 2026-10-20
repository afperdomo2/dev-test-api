import { fireEvent, render, screen, waitFor } from '@testing-library/react-native'
import { useLocalSearchParams } from 'expo-router'
import type { NextQuestionItem, Tip } from '@devtest/shared'

import SessionStudyPage from './SessionStudyPage'
import { clearPendingTip, setPendingTip } from '@/features/sessions/dailyTipStore'
import {
  getNextQuestion,
  getSessionById,
  getSessionSummary,
  submitAnswer,
} from '@/core/api/services/sessions.service'

jest.mock('@/core/api/services/sessions.service', () => ({
  getSessionById: jest.fn(),
  getSessionSummary: jest.fn(),
  getNextQuestion: jest.fn(),
  submitAnswer: jest.fn(),
  finishSession: jest.fn(),
}))

const question: NextQuestionItem = {
  id: 'q1',
  type: 'single_choice',
  content: 'Pregunta de prueba',
  difficulty: 'beginner',
  options: [
    { id: 'o1', content: 'Opción A', isCorrect: true },
    { id: 'o2', content: 'Opción B', isCorrect: false },
  ],
  codeChallenge: null,
  topics: [],
}

const mockedParams = useLocalSearchParams as unknown as jest.Mock

describe('SessionStudyPage', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    clearPendingTip()
    mockedParams.mockReturnValue({ id: 's1' })
    jest.mocked(getSessionById).mockResolvedValue({ id: 's1', status: 'in_progress' } as never)
    jest
      .mocked(getSessionSummary)
      .mockResolvedValue({ answerCount: 0, questionLimit: null } as never)
    jest.mocked(getNextQuestion).mockResolvedValue(question)
  })

  it('carga y muestra la pregunta', async () => {
    await render(<SessionStudyPage />)

    expect(await screen.findByText(/Pregunta de prueba/)).toBeOnTheScreen()
    expect(screen.getByLabelText('Opción A')).toBeOnTheScreen()
  })

  it('responder envía la respuesta al backend', async () => {
    jest
      .mocked(submitAnswer)
      .mockResolvedValue({ id: 'a1', isCorrect: true, explanation: '' } as never)

    await render(<SessionStudyPage />)
    await screen.findByText(/Pregunta de prueba/)

    await fireEvent.press(screen.getByLabelText('Opción A'))
    await fireEvent.press(screen.getByLabelText('Confirmar respuesta'))

    await waitFor(() =>
      expect(submitAnswer).toHaveBeenCalledWith('s1', {
        questionId: 'q1',
        selectedOptions: ['o1'],
        clientAnswerId: expect.any(String),
      }),
    )
  })

  it('un error de red no avanza y el reintento reutiliza el clientAnswerId', async () => {
    jest
      .mocked(submitAnswer)
      .mockRejectedValueOnce({
        title: 'Network Error',
        detail: 'No se pudo conectar con el servidor',
      })
      .mockResolvedValueOnce({ id: 'a1', isCorrect: true, explanation: '' } as never)

    await render(<SessionStudyPage />)
    await screen.findByText(/Pregunta de prueba/)

    await fireEvent.press(screen.getByLabelText('Opción A'))
    await fireEvent.press(screen.getByLabelText('Confirmar respuesta'))

    expect(
      await screen.findByText('Sin conexión. Revisa tu red e inténtalo de nuevo.'),
    ).toBeOnTheScreen()
    expect(screen.queryByLabelText('Siguiente pregunta')).toBeNull()
    expect(screen.getByLabelText('Confirmar respuesta')).toBeOnTheScreen()

    await fireEvent.press(screen.getByLabelText('Reintentar'))

    await waitFor(() => expect(submitAnswer).toHaveBeenCalledTimes(2))
    const [first, second] = jest.mocked(submitAnswer).mock.calls
    expect(second?.[1]?.clientAnswerId).toBe(first?.[1]?.clientAnswerId)
    expect(second?.[1]).toMatchObject({ questionId: 'q1', selectedOptions: ['o1'] })
    expect(screen.queryByText('Sin conexión. Revisa tu red e inténtalo de nuevo.')).toBeNull()
  })

  it('un error de API muestra el detalle y ofrece reintentar', async () => {
    jest
      .mocked(submitAnswer)
      .mockRejectedValueOnce({ title: 'Conflict', status: 409, detail: 'La sesión ya fue completada' })

    await render(<SessionStudyPage />)
    await screen.findByText(/Pregunta de prueba/)

    await fireEvent.press(screen.getByLabelText('Opción A'))
    await fireEvent.press(screen.getByLabelText('Confirmar respuesta'))

    expect(await screen.findByText('La sesión ya fue completada')).toBeOnTheScreen()
    expect(screen.getByLabelText('Reintentar')).toBeOnTheScreen()
  })

  it('un doble toque en confirmar envía una sola vez', async () => {
    let resolveSubmit: (value: unknown) => void = () => undefined
    jest.mocked(submitAnswer).mockReturnValue(
      new Promise((resolve) => {
        resolveSubmit = resolve
      }) as never,
    )

    await render(<SessionStudyPage />)
    await screen.findByText(/Pregunta de prueba/)

    await fireEvent.press(screen.getByLabelText('Opción A'))
    await fireEvent.press(screen.getByLabelText('Confirmar respuesta'))
    await fireEvent.press(screen.getByLabelText('Confirmar respuesta'))

    expect(submitAnswer).toHaveBeenCalledTimes(1)

    resolveSubmit({ id: 'a1', isCorrect: true })
    expect(await screen.findByLabelText('Siguiente pregunta')).toBeOnTheScreen()
  })

  it('el clientAnswerId cambia al pasar a otra pregunta', async () => {
    jest
      .mocked(getNextQuestion)
      .mockResolvedValueOnce(question)
      .mockResolvedValueOnce({ ...question, id: 'q2', content: 'Segunda pregunta' })
    jest.mocked(submitAnswer).mockResolvedValue({ id: 'a1', isCorrect: true } as never)

    await render(<SessionStudyPage />)
    await screen.findByText(/Pregunta de prueba/)

    await fireEvent.press(screen.getByLabelText('Opción A'))
    await fireEvent.press(screen.getByLabelText('Confirmar respuesta'))
    await screen.findByLabelText('Siguiente pregunta')

    await fireEvent.press(screen.getByLabelText('Siguiente pregunta'))
    expect(await screen.findByText(/Segunda pregunta/)).toBeOnTheScreen()

    await fireEvent.press(screen.getByLabelText('Opción A'))
    await fireEvent.press(screen.getByLabelText('Confirmar respuesta'))

    await waitFor(() => expect(submitAnswer).toHaveBeenCalledTimes(2))
    const [first, second] = jest.mocked(submitAnswer).mock.calls
    expect(second?.[1]?.clientAnswerId).not.toBe(first?.[1]?.clientAnswerId)
  })

  it('sin más preguntas pasa al estado terminado', async () => {
    jest.mocked(getNextQuestion).mockRejectedValue({ status: 404 })

    await render(<SessionStudyPage />)

    expect(await screen.findByText(/Respondiste/)).toBeOnTheScreen()
  })

  describe('tip del día', () => {
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

    it('muestra la tarjeta al empezar la sesión', async () => {
      setPendingTip('s1', tip)

      await render(<SessionStudyPage />)

      expect(await screen.findByText('Tip del día')).toBeOnTheScreen()
      expect(screen.getByText('Go')).toBeOnTheScreen()
    })

    it('sin tip no muestra la tarjeta', async () => {
      await render(<SessionStudyPage />)
      await screen.findByText(/Pregunta de prueba/)

      expect(screen.queryByText('Tip del día')).toBeNull()
    })

    it('al cerrarla deja de mostrarse', async () => {
      setPendingTip('s1', tip)

      await render(<SessionStudyPage />)
      await screen.findByText('Tip del día')

      await fireEvent.press(screen.getByLabelText('Cerrar tip del día'))

      expect(screen.queryByText('Tip del día')).toBeNull()
    })

    it('no se muestra durante el resultado de una respuesta', async () => {
      setPendingTip('s1', tip)
      jest
        .mocked(submitAnswer)
        .mockResolvedValue({ id: 'a1', isCorrect: true, explanation: '' } as never)

      await render(<SessionStudyPage />)
      await screen.findByText('Tip del día')

      await fireEvent.press(screen.getByLabelText('Opción A'))
      await fireEvent.press(screen.getByLabelText('Confirmar respuesta'))

      await waitFor(() => expect(screen.queryByText('Tip del día')).toBeNull())
    })
  })
})
