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
      }),
    )
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
