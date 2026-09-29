const DEFAULT_API_BASE_URL = 'http://localhost:8080'

function readApiBaseUrl(): string {
  const env = process.env.EXPO_PUBLIC_API_BASE_URL?.trim()
  if (env) return env.replace(/\/+$/, '')
  return DEFAULT_API_BASE_URL
}

export const API_BASE_URL = readApiBaseUrl()
