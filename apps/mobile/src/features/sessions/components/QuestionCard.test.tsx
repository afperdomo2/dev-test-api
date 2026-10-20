import { fireEvent, render, screen } from '@testing-library/react-native'
import type { NextQuestionItem } from '@devtest/shared'

import { QuestionCard } from './QuestionCard'

const question: NextQuestionItem = {
  id: 'q1',
  type: 'single_choice',
  content: '¿Cuál es la capital de Francia?',
  difficulty: 'beginner',
  options: [
    { id: 'o1', content: 'París', isCorrect: true },
    { id: 'o2', content: 'Lyon', isCorrect: false },
  ],
  codeChallenge: null,
  topics: ['Geografía'],
}

const multipleQuestion: NextQuestionItem = {
  ...question,
  id: 'q2',
  type: 'multiple_choice',
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
})
