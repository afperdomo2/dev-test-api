export function useFormErrors() {
  function extractFieldErrors(error: unknown): Record<string, string> {
    const fieldErrors: Record<string, string> = {}

    if (error && typeof error === 'object' && 'detail' in error) {
      const detail = (error as { detail: string }).detail
      if (!detail) return fieldErrors

      const parts = detail.split(/[;,]\s*/)
      for (const part of parts) {
        const [field, ...rest] = part.split(':')
        if (field && rest.length > 0) {
          fieldErrors[field.trim().toLowerCase()] = rest.join(':').trim()
        }
      }
    }

    return fieldErrors
  }

  return {
    extractFieldErrors,
  }
}
