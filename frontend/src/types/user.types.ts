export interface User {
  id: string
  email: string
  is_admin: boolean
  created_at: string
  updated_at: string
}

export interface CreateUserRequest {
  email: string
  password: string
  is_admin?: boolean
}
