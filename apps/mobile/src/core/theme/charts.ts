import { colors } from './tokens'

/**
 * Escala de intensidad del mapa de calor: 4 niveles (0 = sin actividad).
 * Deriva de `colors.primary` con alfa para no introducir hex nuevos.
 */
export const HEAT_COLORS = [
  colors.muted,
  `${colors.primary}33`,
  `${colors.primary}80`,
  colors.primary,
] as const

export type HeatLevel = 0 | 1 | 2 | 3

/** Nivel de intensidad de un día respecto al máximo del periodo. */
export function heatLevel(value: number, max: number): HeatLevel {
  if (value <= 0 || max <= 0) return 0

  const ratio = value / max
  if (ratio <= 1 / 3) return 1
  if (ratio <= 2 / 3) return 2
  return 3
}
