export function getErrorDetail(err: unknown, fallback: string): string {
  if (err && typeof err === 'object' && 'detail' in err) {
    return String((err as { detail: unknown }).detail)
  }
  return fallback
}

/**
 * `api/client.ts` lanza `{ title: 'Network Error' }` cuando `fetch` falla
 * (sin red, DNS, timeout). Cualquier respuesta HTTP, aunque sea 5xx, no es un
 * error de red.
 */
export function isNetworkError(err: unknown): boolean {
  return (
    err !== null &&
    typeof err === 'object' &&
    'title' in err &&
    (err as { title?: unknown }).title === 'Network Error'
  )
}
