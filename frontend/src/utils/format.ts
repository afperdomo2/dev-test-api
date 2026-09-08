export function formatDateTime(dateString: string): string {
  return new Intl.DateTimeFormat('es-CO', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(dateString))
}

export function formatDate(dateString: string): string {
  return new Intl.DateTimeFormat('es-CO', {
    dateStyle: 'medium',
  }).format(new Date(dateString))
}

export function formatShortDateTime(dateString: string): string {
  if (!dateString) return ''
  const d = new Date(dateString)
  if (Number.isNaN(d.getTime())) return dateString
  return d.toLocaleString('es-CO', {
    dateStyle: 'short',
    timeStyle: 'short',
  })
}

export function formatScore(score: number): string {
  return `${score.toFixed(1)}%`
}

/**
 * Tiempo relativo para humanos. Ej: "Se reinicia en 48 minutos",
 * "Se reinicia en 5 días 22 horas", "Se reinicia en 12 días 17 horas".
 * Pensado para `resetsAt` de OpenCode pero reutilizable en otros componentes.
 */
export function formatTimeUntil(dateString: string | undefined): string {
  if (!dateString) return ''
  const target = new Date(dateString)
  if (Number.isNaN(target.getTime())) return dateString
  const diff = target.getTime() - Date.now()
  if (diff <= 0) return 'Se reinicia ahora'
  if (diff < 60 * 1000) return 'Se reinicia en menos de un minuto'

  const totalMinutes = Math.floor(diff / 60000)
  const days = Math.floor(totalMinutes / 1440)
  const hours = Math.floor((totalMinutes % 1440) / 60)
  const minutes = totalMinutes % 60

  const dayLabel = days === 1 ? 'día' : 'días'
  const hourLabel = hours === 1 ? 'hora' : 'horas'
  const minuteLabel = minutes === 1 ? 'minuto' : 'minutos'

  if (days > 0) {
    if (hours > 0) return `Se reinicia en ${days} ${dayLabel} ${hours} ${hourLabel}`
    return `Se reinicia en ${days} ${dayLabel}`
  }
  if (hours > 0) {
    if (minutes > 0) return `Se reinicia en ${hours} ${hourLabel} ${minutes} ${minuteLabel}`
    return `Se reinicia en ${hours} ${hourLabel}`
  }
  return `Se reinicia en ${minutes} ${minuteLabel}`
}
