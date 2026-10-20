import type { Tip } from '@devtest/shared'

/**
 * El tip del día solo viaja en la respuesta de `POST /sessions/daily`
 * (no hay endpoint GET), así que se guarda en memoria entre Home y la
 * pantalla de la sesión. No se persiste: si la app se reinicia a mitad de
 * la sesión, el tip se pierde (decisión D-30).
 */
type PendingTip = {
  sessionId: string
  tip: Tip
}

let pending: PendingTip | null = null

export function setPendingTip(sessionId: string, tip?: Tip): void {
  pending = tip ? { sessionId, tip } : null
}

/** Devuelve el tip de esa sesión. NO lo consume: se puede leer varias veces. */
export function getPendingTip(sessionId: string): Tip | null {
  if (!pending || pending.sessionId !== sessionId) return null
  return pending.tip
}

export function clearPendingTip(): void {
  pending = null
}
