export function getErrorDetail(err: unknown, fallback: string): string {
  if (err && typeof err === 'object' && 'detail' in err) {
    return String((err as { detail: unknown }).detail)
  }
  return fallback
}
