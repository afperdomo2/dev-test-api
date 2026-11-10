import { fireEvent, render, screen, within } from '@testing-library/react-native'
import type { NextQuestionItem, SessionAnswer } from '@devtest/shared'

import { QuestionCard } from './QuestionCard'

const question: NextQuestionItem = {
  id: 'q1',
  type: 'single_choice',
  content: '¿Cuál es la capital de Francia?',
  difficulty: 'beginner',
  // El backend nunca manda `isCorrect` en /next: el fixture debe reflejarlo.
  options: [
    { id: 'o1', content: 'París' },
    { id: 'o2', content: 'Lyon' },
  ],
  codeChallenge: null,
  topics: ['Geografía'],
}

const resultBase: SessionAnswer = {
  id: 'a1',
  questionId: 'q1',
  answerText: '',
  selectedOptions: [],
  isCorrect: true,
  aiFeedback: '',
  explanation: '',
  responseTimeMs: 0,
}

describe('QuestionCard', () => {
  it('muestra el contenido, los temas y las opciones', async () => {
    await render(
      <QuestionCard
        question={question}
        mode="answer"
        selected={[]}
        lastResult={null}
        onToggleOption={jest.fn()}
      />,
    )

    expect(screen.getByText(/capital de Francia/)).toBeOnTheScreen()
    expect(screen.getByText('Geografía')).toBeOnTheScreen()
    expect(screen.getByLabelText('París')).toBeOnTheScreen()
    expect(screen.getByLabelText('Lyon')).toBeOnTheScreen()
  })

  it('no rompe si el backend omite options/topics (contrato viejo)', async () => {
    const legacyQuestion = {
      ...question,
      options: undefined,
      topics: undefined,
    } as unknown as NextQuestionItem

    await render(
      <QuestionCard
        question={legacyQuestion}
        mode="answer"
        selected={[]}
        lastResult={null}
        onToggleOption={jest.fn()}
      />,
    )

    expect(screen.getByText(/capital de Francia/)).toBeOnTheScreen()
  })

  it('avisa de la opción pulsada (la selección la decide la pantalla)', async () => {
    const onToggleOption = jest.fn()
    await render(
      <QuestionCard
        question={question}
        mode="answer"
        selected={[]}
        lastResult={null}
        onToggleOption={onToggleOption}
      />,
    )

    await fireEvent.press(screen.getByLabelText('París'))

    expect(onToggleOption).toHaveBeenCalledWith('o1')
  })

  it('marca como seleccionadas las opciones que recibe', async () => {
    await render(
      <QuestionCard
        question={question}
        mode="answer"
        selected={['o2']}
        lastResult={null}
        onToggleOption={jest.fn()}
      />,
    )

    expect(screen.getByLabelText('París').props.accessibilityState).toEqual(
      expect.objectContaining({ selected: false }),
    )
    expect(screen.getByLabelText('Lyon').props.accessibilityState).toEqual(
      expect.objectContaining({ selected: true }),
    )
  })

  it('con disabled no permite togglear (mientras se envía)', async () => {
    const onToggleOption = jest.fn()
    await render(
      <QuestionCard
        question={question}
        mode="answer"
        selected={[]}
        lastResult={null}
        onToggleOption={onToggleOption}
        disabled
      />,
    )

    await fireEvent.press(screen.getByLabelText('París'))

    expect(onToggleOption).not.toHaveBeenCalled()
  })

  it('en modo resultado no permite togglear', async () => {
    const onToggleOption = jest.fn()
    await render(
      <QuestionCard
        question={question}
        mode="result"
        selected={['o1']}
        lastResult={resultBase}
        onToggleOption={onToggleOption}
      />,
    )

    await fireEvent.press(screen.getByLabelText('Lyon'))

    expect(onToggleOption).not.toHaveBeenCalled()
  })

  it('en resultado con acierto marca la opción correcta en verde', async () => {
    await render(
      <QuestionCard
        question={question}
        mode="result"
        selected={['o1']}
        lastResult={{
          ...resultBase,
          selectedOptions: ['o1'],
          isCorrect: true,
          correctOptions: ['o1'],
        }}
        onToggleOption={jest.fn()}
      />,
    )

    const correcta = screen.getByLabelText('París, correcta')
    expect(within(correcta).getByText('check-circle')).toBeOnTheScreen()
    expect(within(screen.getByLabelText('Lyon')).queryByText('close-circle')).toBeNull()
  })

  it('en resultado con fallo marca la correcta en verde y la elegida en rojo', async () => {
    await render(
      <QuestionCard
        question={question}
        mode="result"
        selected={['o2']}
        lastResult={{
          ...resultBase,
          selectedOptions: ['o2'],
          isCorrect: false,
          correctOptions: ['o1'],
        }}
        onToggleOption={jest.fn()}
      />,
    )

    expect(
      within(screen.getByLabelText('París, correcta')).getByText('check-circle'),
    ).toBeOnTheScreen()
    expect(
      within(screen.getByLabelText('Lyon, incorrecta')).getByText('close-circle'),
    ).toBeOnTheScreen()
  })

  it('sin correctOptions usa isCorrect + selección (backend viejo)', async () => {
    const { rerender } = await render(
      <QuestionCard
        question={question}
        mode="result"
        selected={['o1']}
        lastResult={{ ...resultBase, selectedOptions: ['o1'], isCorrect: true }}
        onToggleOption={jest.fn()}
      />,
    )

    expect(
      within(screen.getByLabelText('París, correcta')).getByText('check-circle'),
    ).toBeOnTheScreen()

    await rerender(
      <QuestionCard
        question={question}
        mode="result"
        selected={['o2']}
        lastResult={{ ...resultBase, selectedOptions: ['o2'], isCorrect: false }}
        onToggleOption={jest.fn()}
      />,
    )

    expect(
      within(screen.getByLabelText('Lyon, incorrecta')).getByText('close-circle'),
    ).toBeOnTheScreen()
    expect(within(screen.getByLabelText('París')).queryByText('check-circle')).toBeNull()
  })
})
