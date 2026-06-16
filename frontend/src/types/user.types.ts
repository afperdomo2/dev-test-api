export interface User {
  id: string
  email: string
  isAdmin: boolean
  createdAt: string
  updatedAt: string
}

export interface CreateUserRequest {
  email: string
  password: string
  is_admin?: boolean
}

export interface UpdateUserRequest {
  password?: string
  is_admin?: boolean
}
