/**
 * Identidad visual de la app por entorno.
 *
 * `VITE_APP_NAME` y `VITE_APP_COLOR` se compilan en el bundle: cada entorno
 * (dev, pruebas, producción) debe construirse con sus propios valores. No se
 * persisten en localStorage ni en la DB, y no hay detección automática.
 */
export const DEFAULT_APP_NAME = 'DevTest'

/**
 * Colores admitidos por `VITE_APP_COLOR`. El valor de la variable es la clave
 * (p. ej. `amber`), no el hex. Mantener la lista corta y con buen contraste.
 */
export const APP_COLOR_PRESETS = {
  blue: '#1976D2',
  green: '#2E7D32',
  amber: '#F9A825',
  red: '#C62828',
  purple: '#6A1B9A',
  teal: '#00796B',
  slate: '#37474F',
} as const

export type AppColorName = keyof typeof APP_COLOR_PRESETS

export const DEFAULT_APP_COLOR: AppColorName = 'blue'

/** Un nombre vacío o ausente cae al default (`DevTest`). */
export function resolveAppName(raw?: string): string {
  return raw?.trim() || DEFAULT_APP_NAME
}

/**
 * Resuelve un preset de color. Si el valor no existe, cae al default (`blue`).
 * Acepta mayúsculas y espacios alrededor para tolerar errores de tipeo.
 */
export function resolveAppColor(raw?: string): string {
  const key = raw?.trim().toLowerCase()
  if (key && Object.prototype.hasOwnProperty.call(APP_COLOR_PRESETS, key)) {
    return APP_COLOR_PRESETS[key as AppColorName]
  }
  return APP_COLOR_PRESETS[DEFAULT_APP_COLOR]
}

export const APP_NAME = resolveAppName(import.meta.env.VITE_APP_NAME)
export const APP_COLOR = resolveAppColor(import.meta.env.VITE_APP_COLOR)
