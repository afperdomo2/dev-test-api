import { renderHook } from '@testing-library/react-native'

import { useScrollToTopOnChange } from './useScrollToTopOnChange'

describe('useScrollToTopOnChange', () => {
  it('lleva el scroll al inicio al cambiar de pregunta', async () => {
    const scrollTo = jest.fn()
    const ref = { current: { scrollTo } }

    const { rerender } = await renderHook(
      ({ questionId }: { questionId: string | null }) =>
        useScrollToTopOnChange(ref, questionId),
      { initialProps: { questionId: 'q1' } },
    )

    expect(scrollTo).toHaveBeenCalledTimes(1)
    expect(scrollTo).toHaveBeenLastCalledWith({ y: 0, animated: false })

    await rerender({ questionId: 'q2' })

    expect(scrollTo).toHaveBeenCalledTimes(2)
    expect(scrollTo).toHaveBeenLastCalledWith({ y: 0, animated: false })
  })

  it('no hace nada si todavía no hay pregunta', async () => {
    const scrollTo = jest.fn()
    const ref = { current: { scrollTo } }

    await renderHook(
      ({ questionId }: { questionId: string | null }) =>
        useScrollToTopOnChange(ref, questionId),
      { initialProps: { questionId: null } },
    )

    expect(scrollTo).not.toHaveBeenCalled()
  })

  it('no rompe si el ScrollView aún no tiene ref', async () => {
    const ref = { current: null }

    await renderHook(
      ({ questionId }: { questionId: string | null }) =>
        useScrollToTopOnChange(ref, questionId),
      { initialProps: { questionId: 'q1' } },
    )

    // Sin ref no hay scroll que hacer, pero tampoco debe lanzar.
    expect(ref.current).toBeNull()
  })
})
