import * as SecureStore from 'expo-secure-store'

const DISMISSED_ON_KEY = 'daily_tip_dismissed_on'

/**
 * Fecha UTC del día en curso (`YYYY-MM-DD`). El backend elige el tip del día
 * con la misma referencia (UTC), así que ambos lados rotan a la vez.
 */
export function todayUTC(now: Date = new Date()): string {
  return now.toISOString().slice(0, 10)
}

/**
 * Día en que el usuario cerró el tip del día, o `null` si nunca lo cerró.
 * La preferencia es del dispositivo, no de la cuenta.
 */
export async function getTipDismissedOn(): Promise<string | null> {
  return SecureStore.getItemAsync(DISMISSED_ON_KEY)
}

/** Recuerda que el tip se cerró ese día; al día siguiente vuelve a mostrarse. */
export async function setTipDismissedOn(day: string): Promise<void> {
  await SecureStore.setItemAsync(DISMISSED_ON_KEY, day)
}

/** `true` si el tip de hoy ya se cerró. */
export async function isTipDismissedToday(now: Date = new Date()): Promise<boolean> {
  return (await getTipDismissedOn()) === todayUTC(now)
}
