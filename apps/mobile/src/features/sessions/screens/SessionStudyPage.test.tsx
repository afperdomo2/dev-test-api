import { fireEvent, render, screen, waitFor } from '@testing-library/react-native'
import { useLocalSearchParams } from 'expo-router'
import type { NextQuestionItem } from '@devtest/shared'

import SessionStudyPage from './SessionStudyPage'
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
})
