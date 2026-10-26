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
 * Tiempo relativo compacto para listados. Ej: "Ahora", "Hace 12 min",
 * "Hace 3h", "Ayer", "5 oct", "5 oct 2025".
 */
export function formatRelativeTime(dateString: string): string {
  const date = new Date(dateString)
  if (Number.isNaN(date.getTime())) return ''

  const diffMs = Date.now() - date.getTime()
  const diffSec = Math.floor(diffMs / 1000)
  const diffMin = Math.floor(diffSec / 60)
  const diffH = Math.floor(diffMin / 60)
  const diffD = Math.floor(diffH / 24)

  if (diffMin < 1) return 'Ahora'
  if (diffMin < 60) return `Hace ${diffMin} min`
  if (diffH < 24) return `Hace ${diffH}h`
  if (diffD === 1) return 'Ayer'
  if (diffD < 7) return `Hace ${diffD} días`

  const month = date.toLocaleDateString('es-CO', { month: 'short' }).replace(/\.$/, '')
  const day = date.getDate()
  const year = date.getFullYear()
  const currentYear = new Date().getFullYear()

  if (year === currentYear) return `${day} ${month}`
  return `${day} ${month} ${year}`
}

/**
 * Duración legible entre dos fechas. Devuelve cadena vacía si no hay
 * `finishedAt` o si el rango es inválido. Ej: "34 min", "1 h 12 min".
 */
export function formatDuration(startedAt: string, finishedAt: string | null): string {
  if (!finishedAt) return ''
  const start = new Date(startedAt).getTime()
  const end = new Date(finishedAt).getTime()
  if (Number.isNaN(start) || Number.isNaN(end) || end <= start) return ''

  const totalMinutes = Math.round((end - start) / 60000)
  if (totalMinutes < 1) return 'menos de 1 min'
  const hours = Math.floor(totalMinutes / 60)
  const minutes = totalMinutes % 60
  if (hours === 0) return `${minutes} min`
  if (minutes === 0) return `${hours} h`
  return `${hours} h ${minutes} min`
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
