<script setup lang="ts">
import { computed, ref } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { usersListOptions, deleteUserMutation } from '@/queries/users.queries'
import { useAppStore } from '@/stores/app.store'
import { usePagination } from '@/composables/usePagination'
import ListPageHeader from '@/components/ListPageHeader.vue'
import PaginatedFooter from '@/components/PaginatedFooter.vue'
import UserTable from '../components/UserTable.vue'
import UserFormDialog from '../components/UserFormDialog.vue'
import type { User } from '@/types/user.types'

const appStore = useAppStore()
const queryClient = useQueryClient()
const { page, perPage, reset: resetPagination } = usePagination()

const { data, isLoading } = useQuery(
  usersListOptions(
    () => page.value,
    () => perPage.value,
  ),
)

const deleteMut = useMutation(deleteUserMutation())
const deleteTarget = ref<User | null>(null)
const deleteDialog = ref(false)

const dialogOpen = ref(false)
const editingUser = ref<User | null>(null)

const userList = computed<Array<User>>(() => {
  return data.value?.data ?? []
})

const totalUsers = computed(() => data.value?.meta?.total ?? 0)

function onPerPageChange(val: number) {
  perPage.value = val
  page.value = 1
}

function handleRefresh() {
  resetPagination()
  queryClient.invalidateQueries({ queryKey: ['users', 'list'] })
}

function openCreate() {
  editingUser.value = null
  dialogOpen.value = true
}

function openEdit(user: User) {
  editingUser.value = user
  dialogOpen.value = true
}

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
    <ListPageHeader
      title="Usuarios"
      create-label="Nuevo usuario"
      @refresh="handleRefresh"
      @create="openCreate"
    />

    <v-card>
      <v-card-text>
        <UserTable
          :users="userList"
          :loading="isLoading"
          :items-per-page="perPage"
          @edit="openEdit"
          @delete="confirmDelete"
        >
          <template #footer>
            <PaginatedFooter
              :page="page"
              :per-page="perPage"
              :total="totalUsers"
              :in-table="true"
              @update:page="page = $event"
              @update:per-page="onPerPageChange"
            />
          </template>
        </UserTable>
      </v-card-text>
    </v-card>

    <UserFormDialog v-model="dialogOpen" :user="editingUser" @saved="dialogOpen = false" />

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
          <v-btn color="error" :loading="deleteMut.isPending.value" @click="executeDelete">
            Eliminar
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>
