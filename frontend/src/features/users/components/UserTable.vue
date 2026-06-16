<script setup lang="ts">
import type { User } from '@/types/user.types'
import { formatDateTime } from '@/utils/format'

interface Props {
  users: Array<User>
  loading: boolean
  itemsPerPage: number
}

defineProps<Props>()

const emit = defineEmits<{
  edit: [user: User]
  delete: [user: User]
}>()

const headers = [
  { title: 'Email', key: 'email', sortable: false },
  { title: 'Rol', key: 'isAdmin', sortable: false, align: 'center' as const },
  { title: 'Creado', key: 'createdAt', sortable: false, align: 'center' as const },
  { title: 'Acciones', key: 'actions', sortable: false, align: 'center' as const },
]
</script>

<template>
  <v-data-table
    :headers="headers"
    :items="users"
    :items-per-page="itemsPerPage"
    :loading="loading"
    no-data-text="No hay usuarios"
    loading-text="Cargando usuarios..."
    hover
  >
    <template #[`item.isAdmin`]="{ item }">
      <div class="d-flex justify-center">
        <v-chip :color="item.isAdmin ? 'primary' : ''" size="small" variant="tonal">
          {{ item.isAdmin ? 'Admin' : 'Usuario' }}
        </v-chip>
      </div>
    </template>

    <template #[`item.createdAt`]="{ item }">
      <div class="d-flex justify-center">
        {{ formatDateTime(item.createdAt) }}
      </div>
    </template>

    <template #[`item.actions`]="{ item }">
      <div class="d-flex ga-1 justify-center">
        <v-btn
          icon="mdi-pencil"
          variant="text"
          size="small"
          color="primary"
          @click="emit('edit', item)"
        />
        <v-btn
          icon="mdi-delete"
          variant="text"
          size="small"
          color="error"
          @click="emit('delete', item)"
        />
      </div>
    </template>

    <template #bottom>
      <slot name="footer" />
    </template>
  </v-data-table>
</template>
