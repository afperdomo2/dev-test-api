import axios from 'axios'
import type { AxiosError, InternalAxiosRequestConfig } from 'axios'
import type { ApiError } from '@/types/api.types'
import { getToken, removeToken } from '@/utils/storage'

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

apiClient.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = getToken()
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

apiClient.interceptors.response.use(
  (response) => {
    let data = response.data

    if (data && typeof data === 'object') {
      // unwrap envelope { data: ... } → inner value
      if ('data' in data && !('meta' in data)) {
        data = data.data
      }
    }

    return { ...response, data }
  },
  (error: AxiosError<ApiError>) => {
    if (error.response?.status === 401) {
      removeToken()
    }

    const apiError: ApiError = error.response?.data ?? {
      type: 'about:blank',
      title: 'Internal Server Error',
      status: 500,
      detail: 'Ocurrió un error inesperado',
      instance: error.config?.url ?? '',
    }

    return Promise.reject(apiError)
  },
)

export default apiClient
