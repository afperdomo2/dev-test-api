import { queryOptions } from '@tanstack/vue-query'
import * as usersService from '@/api/services/users.service'

export function usersListOptions(page: () => number, perPage: () => number) {
  return queryOptions({
    queryKey: ['users', 'list', page, perPage],
    queryFn: () => usersService.listUsers(page(), perPage()),
    staleTime: 30 * 1000,
  })
}

export function userDetailOptions(id: () => string) {
  return queryOptions({
    queryKey: ['users', 'detail', id],
    queryFn: () => usersService.getUserById(id()),
    staleTime: 60 * 1000,
  })
}

export function createUserMutation() {
  return {
    mutationKey: ['users', 'create'],
    mutationFn: usersService.createUser,
  }
}

export function deleteUserMutation() {
  return {
    mutationKey: ['users', 'delete'],
    mutationFn: usersService.deleteUser,
  }
}
