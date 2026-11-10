import { useLayoutEffect } from 'react'

type Scrollable = {
  scrollTo: (options: { y: number; animated: boolean }) => void
}

/**
 * Devuelve el scroll al inicio cuando cambia `resetKey`. Se usa para que la
 * siguiente pregunta aparezca desde arriba aunque el usuario haya pulsado
 * "Siguiente" desde el final de una pregunta larga.
 *
 * `useLayoutEffect` (no `useEffect`) para que ocurra antes de pintar y no se
 * vea un frame con el contenido nuevo en el offset viejo.
 */
export function useScrollToTopOnChange(
  ref: { current: Scrollable | null },
  resetKey: string | null | undefined,
): void {
  useLayoutEffect(() => {
    if (!resetKey) return
    ref.current?.scrollTo({ y: 0, animated: false })
  }, [ref, resetKey])
}
