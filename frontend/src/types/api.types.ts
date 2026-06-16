export interface ApiError {
  type: string
  title: string
  status: number
  detail: string
  instance: string
}

export interface ApiValidationError {
  field: string
  message: string
}

export interface ApiResponse<T> {
  data: T
}

export interface PaginatedMeta {
  total: number
  page: number
  per_page: number
}

export interface PaginatedResponse<T> {
  data: Array<T>
  meta: PaginatedMeta
}
