/** `accuracy` llega como 0..1 y la UI trabaja en porcentaje entero. */
export function toPercent(value: number): number {
  return Math.round(value * 100)
}

/** Duración legible a partir de milisegundos. Ej: "45 min", "3h 20m", "2h". */
export function formatDurationMs(ms: number): string {
  const totalMinutes = Math.round(ms / 60000)
  if (totalMinutes < 60) return `${totalMinutes} min`

  const hours = Math.floor(totalMinutes / 60)
  const minutes = totalMinutes % 60
  return minutes > 0 ? `${hours}h ${minutes}m` : `${hours}h`
}
