export interface ApiError {
  type: string
  title: string
  status: number
  detail: string
  instance: string
}

export interface ApiResponse<T> {
  data: T
}

export interface PaginatedMeta {
  total: number
  page: number
  perPage: number
}

export interface PaginatedResponse<T> {
  data: Array<T>
  meta: PaginatedMeta
}
