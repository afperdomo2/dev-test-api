function toCamelCase(key: string): string {
  return key.replace(/_([a-z])/g, (_, char: string) => char.toUpperCase())
}

export function snakeToCamel<T>(value: T): T {
  if (value === null || value === undefined) {
    return value
  }

  if (Array.isArray(value)) {
    return value.map((item) => snakeToCamel(item)) as unknown as T
  }

  if (typeof value === 'object' && !(value instanceof Date)) {
    const result: Record<string, unknown> = {}
    for (const key of Object.keys(value as Record<string, unknown>)) {
      const camelKey = toCamelCase(key)
      result[camelKey] = snakeToCamel((value as Record<string, unknown>)[key])
    }
    return result as unknown as T
  }

  return value
}
