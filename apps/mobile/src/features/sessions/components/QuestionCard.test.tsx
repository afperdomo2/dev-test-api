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

const multipleQuestion: NextQuestionItem = {
  ...question,
  id: 'q2',
  type: 'multiple_choice',
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
        selectedOptions={[]}
        lastResult={null}
        onAnswer={jest.fn()}
        onNext={jest.fn()}
      />,
    )

    expect(screen.getByText(/capital de Francia/)).toBeOnTheScreen()
    expect(screen.getByText('Geografía')).toBeOnTheScreen()
    expect(screen.getByLabelText('París')).toBeOnTheScreen()
    expect(screen.getByLabelText('Lyon')).toBeOnTheScreen()
  })

  it('envía la opción seleccionada al confirmar', async () => {
    const onAnswer = jest.fn()
    await render(
      <QuestionCard
        question={question}
        mode="answer"
        selectedOptions={[]}
        lastResult={null}
        onAnswer={onAnswer}
        onNext={jest.fn()}
      />,
    )

    await fireEvent.press(screen.getByLabelText('París'))
    await fireEvent.press(screen.getByLabelText('Confirmar respuesta'))

    expect(onAnswer).toHaveBeenCalledWith(['o1'])
  })

  it('en single_choice la nueva selección reemplaza a la anterior', async () => {
    const onAnswer = jest.fn()
    await render(
      <QuestionCard
        question={question}
        mode="answer"
        selectedOptions={[]}
        lastResult={null}
        onAnswer={onAnswer}
        onNext={jest.fn()}
      />,
    )

    await fireEvent.press(screen.getByLabelText('París'))
    await fireEvent.press(screen.getByLabelText('Lyon'))
    await fireEvent.press(screen.getByLabelText('Confirmar respuesta'))

    expect(onAnswer).toHaveBeenCalledWith(['o2'])
  })

  it('en multiple_choice acumula las selecciones', async () => {
    const onAnswer = jest.fn()
    await render(
      <QuestionCard
        question={multipleQuestion}
        mode="answer"
        selectedOptions={[]}
        lastResult={null}
        onAnswer={onAnswer}
        onNext={jest.fn()}
      />,
    )

    await fireEvent.press(screen.getByLabelText('París'))
    await fireEvent.press(screen.getByLabelText('Lyon'))
    await fireEvent.press(screen.getByLabelText('Confirmar respuesta'))

    expect(onAnswer).toHaveBeenCalledWith(['o1', 'o2'])
  })

  it('no envía nada si no hay selección', async () => {
    const onAnswer = jest.fn()
    await render(
      <QuestionCard
        question={question}
        mode="answer"
        selectedOptions={[]}
        lastResult={null}
        onAnswer={onAnswer}
        onNext={jest.fn()}
      />,
    )

    await fireEvent.press(screen.getByLabelText('Confirmar respuesta'))

    expect(onAnswer).not.toHaveBeenCalled()
  })

  it('en modo resultado muestra "Siguiente" y no permite responder', async () => {
    const onAnswer = jest.fn()
    const onNext = jest.fn()
    await render(
      <QuestionCard
        question={question}
        mode="result"
        selectedOptions={['o1']}
        lastResult={null}
        onAnswer={onAnswer}
        onNext={onNext}
      />,
    )

    await fireEvent.press(screen.getByLabelText('Siguiente pregunta'))
    expect(onNext).toHaveBeenCalledTimes(1)

    await fireEvent.press(screen.getByLabelText('Lyon'))
    expect(onAnswer).not.toHaveBeenCalled()
  })

  it('en resultado con acierto marca la opción correcta en verde', async () => {
    await render(
      <QuestionCard
        question={question}
        mode="result"
        selectedOptions={['o1']}
        lastResult={{
          ...resultBase,
          selectedOptions: ['o1'],
          isCorrect: true,
          correctOptions: ['o1'],
        }}
        onAnswer={jest.fn()}
        onNext={jest.fn()}
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
        selectedOptions={['o2']}
        lastResult={{
          ...resultBase,
          selectedOptions: ['o2'],
          isCorrect: false,
          correctOptions: ['o1'],
        }}
        onAnswer={jest.fn()}
        onNext={jest.fn()}
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
        selectedOptions={['o1']}
        lastResult={{ ...resultBase, selectedOptions: ['o1'], isCorrect: true }}
        onAnswer={jest.fn()}
        onNext={jest.fn()}
      />,
    )

    expect(
      within(screen.getByLabelText('París, correcta')).getByText('check-circle'),
    ).toBeOnTheScreen()

    await rerender(
      <QuestionCard
        question={question}
        mode="result"
        selectedOptions={['o2']}
        lastResult={{ ...resultBase, selectedOptions: ['o2'], isCorrect: false }}
        onAnswer={jest.fn()}
        onNext={jest.fn()}
      />,
    )

    expect(
      within(screen.getByLabelText('Lyon, incorrecta')).getByText('close-circle'),
    ).toBeOnTheScreen()
    expect(within(screen.getByLabelText('París')).queryByText('check-circle')).toBeNull()
  })

  it('al pasar del estado de envío a resultado el botón "Siguiente" sigue habilitado', async () => {
    const { rerender } = await render(
      <QuestionCard
        question={question}
        mode="answer"
        selectedOptions={[]}
        lastResult={null}
        onAnswer={jest.fn()}
        onNext={jest.fn()}
        submitting
      />,
    )

    expect(screen.getByLabelText('Confirmar respuesta')).toBeDisabled()

    await rerender(
      <QuestionCard
        question={question}
        mode="result"
        selectedOptions={['o1']}
        lastResult={null}
        onAnswer={jest.fn()}
        onNext={jest.fn()}
      />,
    )

    const next = screen.getByLabelText('Siguiente pregunta')
    expect(next).toBeEnabled()
    // Regresión Android (react-native#56777): si `disabled` pasa de `true` a
    // `undefined`, el host deja de ser touch target y solo responde el texto.
    expect(next.props.accessibilityState).toEqual(expect.objectContaining({ disabled: false }))
  })
})
