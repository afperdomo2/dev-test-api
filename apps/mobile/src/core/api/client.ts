import type { ApiError } from '@devtest/shared'

import { API_BASE_URL } from '@/core/config'
import { getToken, removeToken } from '@/core/auth/storage'

type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE'

export type QueryParams = Record<string, string | number | boolean | undefined>

let unauthorizedHandler: (() => void) | null = null

export function onUnauthorized(handler: () => void): void {
  unauthorizedHandler = handler
}

function buildError(raw: unknown, status: number, detail: string, path: string): ApiError {
  if (raw && typeof raw === 'object' && 'title' in raw && 'status' in raw) {
    return raw as ApiError
  }
  return {
    type: 'about:blank',
    title: 'Internal Server Error',
    status,
    detail,
    instance: path,
  }
}

function unwrap<T>(raw: unknown): T {
  if (raw && typeof raw === 'object' && 'data' in raw && !('meta' in raw)) {
    return (raw as { data: T }).data
  }
  return raw as T
}

function buildUrl(path: string, params?: QueryParams): string {
  if (!params) return `${API_BASE_URL}${path}`
  const filtered = Object.entries(params).filter(
    (entry): entry is [string, string | number | boolean] => entry[1] !== undefined,
  )
  if (filtered.length === 0) return `${API_BASE_URL}${path}`
  const qs = new URLSearchParams(filtered.map(([k, v]) => [k, String(v)])).toString()
  return `${API_BASE_URL}${path}?${qs}`
}

async function request<T>(
  method: HttpMethod,
  path: string,
  body?: unknown,
  params?: QueryParams,
): Promise<T> {
  const token = await getToken()
  const headers: Record<string, string> = { 'Content-Type': 'application/json' }
  if (token) headers.Authorization = `Bearer ${token}`

  let res: Response
  try {
    res = await fetch(buildUrl(path, params), {
      method,
      headers,
      body: body === undefined ? undefined : JSON.stringify(body),
    })
  } catch {
    throw {
      type: 'about:blank',
      title: 'Network Error',
      status: 500,
      detail: 'No se pudo conectar con el servidor',
      instance: path,
    } satisfies ApiError
  }

  if (res.status === 401) {
    await removeToken()
    unauthorizedHandler?.()
  }

  const raw: unknown = await res.json().catch(() => null)

  if (!res.ok) {
    throw buildError(raw, res.status, 'Ocurrió un error inesperado', path)
  }

  return unwrap<T>(raw)
}

export const apiClient = {
  get: <T>(path: string, params?: QueryParams) => request<T>('GET', path, undefined, params),
  post: <T>(path: string, body?: unknown, params?: QueryParams) =>
    request<T>('POST', path, body, params),
  put: <T>(path: string, body?: unknown, params?: QueryParams) =>
    request<T>('PUT', path, body, params),
  delete: <T>(path: string) => request<T>('DELETE', path),
}
