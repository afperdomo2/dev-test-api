<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { createUserMutation } from '@/queries/users.queries'
import { useAppStore } from '@/stores/app.store'
import UserForm from '../components/UserForm.vue'
import type { CreateUserRequest } from '@/types/user.types'

const router = useRouter()
const appStore = useAppStore()
const queryClient = useQueryClient()

const createMut = useMutation(createUserMutation())

async function handleSubmit(data: CreateUserRequest) {
  try {
    await createMut.mutateAsync(data)
    appStore.showSnackbar('Usuario creado correctamente')
    queryClient.invalidateQueries({ queryKey: ['users', 'list'] })
    router.push('/users')
  } catch (err: unknown) {
    const detail =
      err && typeof err === 'object' && 'detail' in err
        ? (err as { detail: string }).detail
        : 'Error al crear usuario'
    appStore.showSnackbar(detail, 'error')
  }
}
</script>

<template>
  <v-container>
    <div class="d-flex align-center mb-4">
      <v-btn
        icon="mdi-arrow-left"
        variant="text"
        to="/users"
        class="mr-2"
      />
      <h1 class="text-h4">Crear usuario</h1>
    </div>

    <v-card max-width="560" class="pa-4">
      <UserForm @submit="handleSubmit" />
    </v-card>
  </v-container>
</template>
