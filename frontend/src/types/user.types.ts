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
  isAdmin?: boolean
}

export interface UpdateUserRequest {
  password?: string
  isAdmin?: boolean
}
