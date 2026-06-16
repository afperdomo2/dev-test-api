<script setup lang="ts">
import { computed, ref } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { usersListOptions, deleteUserMutation } from '@/queries/users.queries'
import { useAppStore } from '@/stores/app.store'
import { usePagination } from '@/composables/usePagination'
import UserTable from '../components/UserTable.vue'
import type { User } from '@/types/user.types'

const appStore = useAppStore()
const queryClient = useQueryClient()
const { page, perPage } = usePagination()

const { data, isLoading } = useQuery(
  usersListOptions(() => page.value, () => perPage.value),
)

const deleteMut = useMutation(deleteUserMutation())
const deleteTarget = ref<User | null>(null)
const deleteDialog = ref(false)

function totalUsers(): number {
  return data.value?.meta?.total ?? 0
}

const userList = computed<Array<User>>(() => {
  return data.value?.data ?? []
})

function confirmDelete(user: User) {
  deleteTarget.value = user
  deleteDialog.value = true
}

async function executeDelete() {
  if (!deleteTarget.value) return
  try {
    await deleteMut.mutateAsync(deleteTarget.value.id)
    appStore.showSnackbar('Usuario eliminado correctamente')
    deleteDialog.value = false
    deleteTarget.value = null
    queryClient.invalidateQueries({ queryKey: ['users', 'list'] })
  } catch (err: unknown) {
    const detail =
      err && typeof err === 'object' && 'detail' in err
        ? (err as { detail: string }).detail
        : 'Error al eliminar usuario'
    appStore.showSnackbar(detail, 'error')
  }
}
</script>

<template>
  <v-container>
    <div class="d-flex align-center justify-space-between mb-4">
      <h1 class="text-h4">Usuarios</h1>
      <v-btn color="primary" prepend-icon="mdi-plus" to="/users/create">
        Nuevo usuario
      </v-btn>
    </div>

    <v-card>
      <v-card-text>
        <UserTable
          :users="userList"
          :loading="isLoading"
          @delete="confirmDelete"
        />

        <div
          v-if="totalUsers() > perPage"
          class="d-flex justify-center mt-4"
        >
          <v-pagination
            v-model="page"
            :length="Math.ceil(totalUsers() / perPage)"
            :total-visible="5"
            @update:model-value="page = $event"
          />
        </div>
      </v-card-text>
    </v-card>

    <v-dialog v-model="deleteDialog" max-width="420">
      <v-card>
        <v-card-title>Eliminar usuario</v-card-title>
        <v-card-text>
          ¿Estás seguro de eliminar a
          <strong>{{ deleteTarget?.email }}</strong
          >?
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="deleteDialog = false"> Cancelar </v-btn>
          <v-btn
            color="error"
            :loading="deleteMut.isPending.value"
            @click="executeDelete"
          >
            Eliminar
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>
