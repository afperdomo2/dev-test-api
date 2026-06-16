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

export function formatScore(score: number): string {
  return `${score.toFixed(1)}%`
}
