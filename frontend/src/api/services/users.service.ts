import apiClient from '@/api/client'
import type { User, CreateUserRequest, UpdateUserRequest } from '@/types/user.types'
import type { PaginatedResponse } from '@/types/api.types'

export async function listUsers(page: number, perPage: number): Promise<PaginatedResponse<User>> {
  const res = await apiClient.get<PaginatedResponse<User>>('/api/v1/users', {
    params: { page, perPage },
  })
  return res.data
}

export async function getUserById(id: string): Promise<User> {
  const res = await apiClient.get<User>(`/api/v1/users/${id}`)
  return res.data
}

export async function createUser(data: CreateUserRequest): Promise<User> {
  const res = await apiClient.post<User>('/api/v1/users', data)
  return res.data
}

export async function updateUser(id: string, data: UpdateUserRequest): Promise<User> {
  const res = await apiClient.put<User>(`/api/v1/users/${id}`, data)
  return res.data
}

export async function deleteUser(id: string): Promise<void> {
  await apiClient.delete(`/api/v1/users/${id}`)
}

export async function getProfile(): Promise<User> {
  const res = await apiClient.get<User>('/api/v1/profile')
  return res.data
}
